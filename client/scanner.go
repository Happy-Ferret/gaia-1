package client

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/notyim/gaia/types"
	"golang.org/x/net/http2"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"sort"
	//"strconv"
	"strings"
	"sync"
	"time"
)

// Only support 10k check atm
const MaxChecks = 10000
const maxRedirects = 10

type Scanner struct {
	Checks             []*types.Check
	CheckInterval      time.Duration
	totalCheck         int
	GaiaServerHost     string
	m                  *sync.Mutex
	CheckRegisterQueue chan string
	f                  *Flusher
}

type headers []string

func (h headers) String() string {
	var o []string
	for _, v := range h {
		o = append(o, "-H "+v)
	}
	return strings.Join(o, " ")
}

func (h *headers) Set(v string) error {
	*h = append(*h, v)
	return nil
}

func (h headers) Len() int      { return len(h) }
func (h headers) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h headers) Less(i, j int) bool {
	a, b := h[i], h[j]

	// server always sorts at the top
	if a == "Server" {
		return true
	}
	if b == "Server" {
		return false
	}

	endtoend := func(n string) bool {
		// https://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html#sec13.5.1
		switch n {
		case "Connection",
			"Keep-Alive",
			"Proxy-Authenticate",
			"Proxy-Authorization",
			"TE",
			"Trailers",
			"Transfer-Encoding",
			"Upgrade":
			return false
		default:
			return true
		}
	}

	x, y := endtoend(a), endtoend(b)
	if x == y {
		// both are of the same class
		return a < b
	}
	return x
}

func NewScanner(gaiaServerHost string) *Scanner {
	s := Scanner{
		Checks:             make([]*types.Check, MaxChecks),
		totalCheck:         0,
		CheckInterval:      15 * time.Second,
		GaiaServerHost:     gaiaServerHost,
		m:                  &sync.Mutex{},
		CheckRegisterQueue: make(chan string),
		f:                  NewFlusher(gaiaServerHost),
	}

	s.Sync()
	go s.Listen()
	go s.Monitor()

	return &s
}

func (s *Scanner) Start() {
	f := NewFlusher(s.GaiaServerHost)
	f.Start()
}

func (s *Scanner) Listen() {
	for c := range s.CheckRegisterQueue {
		check := decode(c)
		s.Checks[s.totalCheck] = check
		s.totalCheck++
	}
}

// Connect to Gaia server to get initial check list
// TODO: We will switch to a TCP server, it makes thing much simpler
func (s *Scanner) Sync() {
	resp, err := http.Get(fmt.Sprintf("%s/checks", s.GaiaServerHost))

	if err != nil {
		log.Fatalf("Cannot sync initialize check %v", err)
	}
	defer resp.Body.Close()
	log.Println("Got initial set of checks. Start decoding checks result set.")

	lineScanner := bufio.NewScanner(resp.Body)
	for lineScanner.Scan() {
		if line := lineScanner.Text(); line != "" {
			if check := decode(line); check != nil {
				log.Println("Found sync check", check)
				s.Checks[s.totalCheck] = check
				s.totalCheck++
			}
		}
	}
}

func (s *Scanner) AddCheck(check *types.Check) {
	s.Checks[s.totalCheck] = check
	s.totalCheck++
}

func (s *Scanner) RemoveCheck() {
}

// Monitor query stat of all checks in its internal data stucture and
// foward back to gaia server
func (s *Scanner) Monitor() {
	ticker := time.Tick(s.CheckInterval)
	for tick := range ticker {
		log.Println("Execute check at", tick)
		for _, c := range s.Checks {
			if c == nil {
				continue
			}

			go s.Execute(c)
		}
	}
}

func (s *Scanner) Execute(check *types.Check) {
	log.Println("Evaluate", check.URI)

	response := types.HTTPCheckResponse{
		CheckID: check.ID,
		Time:    make(map[string]time.Duration),
		Http:    make(map[string]string),
		Tcp:     make(map[string]string),
		Headers: make(map[string]string),
		Error:   false,
	}

	url := parseURL(check.URI)
	visit(url, &response)
	s.f.Write(&response)
}

func decode(s string) *types.Check {
	parts := strings.Split(s, ",")

	return &types.Check{
		ID:   parts[0],
		URI:  parts[1],
		Type: parts[2],
	}
}

func parseURL(uri string) *url.URL {
	if !strings.Contains(uri, "://") && !strings.HasPrefix(uri, "//") {
		uri = "//" + uri
	}

	url, err := url.Parse(uri)
	if err != nil {
		log.Fatalf("could not parse url %q: %v", uri, err)
	}

	if url.Scheme == "" {
		url.Scheme = "http"
		if !strings.HasSuffix(url.Host, ":80") {
			url.Scheme += "s"
		}
	}
	return url
}

func createBody(body string) io.Reader {
	return strings.NewReader(body)
}

func newRequest(method string, url *url.URL, body string) *http.Request {
	req, err := http.NewRequest(method, url.String(), createBody(body))
	if err != nil {
		log.Fatalf("unable to create request: %v", err)
	}
	//for _, h := range httpHeaders {
	//	k, v := headerKeyValue(h)
	//	if strings.EqualFold(k, "host") {
	//		req.Host = v
	//		continue
	//	}
	//	req.Header.Add(k, v)
	//}
	return req
}

func visit(url *url.URL, response *types.HTTPCheckResponse) {
	postBody := ""
	req := newRequest("GET", url, postBody)

	var t0, t1, t2, t3, t4 time.Time

	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) { t0 = time.Now() },
		DNSDone:  func(_ httptrace.DNSDoneInfo) { t1 = time.Now() },
		ConnectStart: func(_, _ string) {
			if t1.IsZero() {
				// connecting to IP
				t1 = time.Now()
			}
		},
		ConnectDone: func(net, addr string, err error) {
			if err != nil {
				response.Error = true
				response.ErrorMessage = fmt.Sprintf("unable to connect to host %v: %v", addr, err)
				return
			}
			t2 = time.Now()
		},
		GotConn:              func(_ httptrace.GotConnInfo) { t3 = time.Now() },
		GotFirstResponseByte: func() { t4 = time.Now() },
	}
	req = req.WithContext(httptrace.WithClientTrace(context.Background(), trace))

	tr := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	switch url.Scheme {
	case "https":
		host, _, err := net.SplitHostPort(req.Host)
		if err != nil {
			host = req.Host
		}

		tr.TLSClientConfig = &tls.Config{
			ServerName:         host,
			InsecureSkipVerify: false,
			Certificates:       nil,
		}

		// Because we create a custom TLSClientConfig, we have to opt-in to HTTP/2.
		// See https://github.com/golang/go/issues/14275
		err = http2.ConfigureTransport(tr)
		if err != nil {
			response.Error = true
			response.ErrorMessage = fmt.Sprintf("failed to prepare transport for HTTP/2: %v", err)
			return
		}
	}

	client := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// always refuse to follow redirects, visit does that
			// manually if required.
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		response.Error = true
		response.ErrorMessage = fmt.Sprintf("failed to read response: %v", err)
		return
	}

	bodyMsg := readResponseBody(req, resp)
	response.Body = bodyMsg
	resp.Body.Close()

	t5 := time.Now() // after read body
	if t0.IsZero() {
		// we skipped DNS
		t0 = t1
	}

	// print status line and headers
	//printf("\n%s%s%s\n", "HTTP", ("/"), color.CyanString("%d.%d %s", resp.ProtoMajor, resp.ProtoMinor, resp.Status))

	names := make([]string, 0, len(resp.Header))
	for k := range resp.Header {
		names = append(names, k)
	}
	sort.Sort(headers(names))
	for _, k := range names {
		//printf("%s %s\n", grayscale(14)(k+":"), color.CyanString(strings.Join(resp.Header[k], ",")))
		response.Headers[k] = strings.Join(resp.Header[k], ",")
	}

	switch url.Scheme {
	case "https":
		response.Time["DNSLookup"] = t1.Sub(t0)
		response.Time["TcpConnection"] = t2.Sub(t1)
		response.Time["TlsHandshake"] = t3.Sub(t2)
		response.Time["ServerProcessing"] = t4.Sub(t3)
		response.Time["ContentTransfer"] = t5.Sub(t4)
		response.Time["NameLookup"] = t1.Sub(t0)
		response.Time["Connect"] = t2.Sub(t0)
		response.Time["Pretransfer"] = t3.Sub(t0)
		response.Time["StartTransfer"] = t4.Sub(t0)
		response.Time["Total"] = t5.Sub(t0)
	case "http":
		response.Time["DNSLookup"] = t1.Sub(t0)
		response.Time["TcpConnection"] = t3.Sub(t1)
		response.Time["ServerProcessing"] = t4.Sub(t3)
		response.Time["ContentTransfer"] = t5.Sub(t4)
		response.Time["NameLookup"] = t1.Sub(t0)
		response.Time["Connect"] = t3.Sub(t0)
		response.Time["StartTransfer"] = t4.Sub(t0)
		response.Time["Total"] = t5.Sub(t0)
	}

	// Temp disable follow, TODO fix this
	followRedirects := false
	redirectsFollowed := 0
	if followRedirects && isRedirect(resp) {
		loc, err := resp.Location()
		if err != nil {
			if err == http.ErrNoLocation {
				// 30x but no Location to follow, give up.
				return
			}
			log.Fatalf("unable to follow redirect: %v", err)
		}

		redirectsFollowed++
		if redirectsFollowed > maxRedirects {
			log.Fatalf("maximum number of redirects (%d) followed", maxRedirects)
		}

		visit(loc, response)
	}
}

func isRedirect(resp *http.Response) bool {
	return resp.StatusCode > 299 && resp.StatusCode < 400
}

// readResponseBody consumes the body of the response.
// readResponseBody returns an informational message about the
// disposition of the response body's contents.
func readResponseBody(req *http.Request, resp *http.Response) string {
	if isRedirect(resp) || req.Method == http.MethodHead {
		return ""
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return string(body)
}

func printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf(format, a...)
}
