package flusher

type flusher interface {
	Write()
	WriteBatch()
}
