package message

import (
	_ "unsafe"
)

type Tracer interface {
	TraceId() string
}

type tracer struct {
	id string
}

func newTracer(id ...string) *tracer {
	if len(id) > 0 {
		return &tracer{id: id[0]}
	}

	return &tracer{id: nextId()}
}

func (t *tracer) TraceId() string {
	return t.id
}

//go:linkname nextId github.com/alomerry/cat-go/cat.nextId
func nextId() string
