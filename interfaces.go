package main

import "io"

type ReaderProducer interface {
	Reader() (io.Reader, error)
}
