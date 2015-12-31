package monitor

import (
	"bufio"
	"log"
	"os"
)

type Coordinator struct {
	AgentChan chan Service
}

func NewCoordinator(agentChan chan Service) *Coordinator {
	c := &Coordinator{agentChan}
	return c
}

func (c *Coordinator) Start() {
	// @TODO
	// Fetch data from source in a loop and notify agent channel about new data
	// or notify agent channel about removing of data
	c.AgentChan <- NewService("https://axcoto.com", "1")
	c.AgentChan <- NewService("http://log.axcoto.com", "2")

	file, err := os.Open("./url")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := scanner.Text()
		c.AgentChan <- NewService(scanner.Text(), url)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
