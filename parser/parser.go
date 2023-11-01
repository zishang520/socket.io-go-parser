package parser

import (
	"github.com/zishang520/engine.io-go-parser/types"
	"github.com/zishang520/engine.io/v2/events"
)

type (

	// A socket.io Encoder instance
	Encoder interface {
		Encode(*Packet) []types.BufferInterface
	}

	// A socket.io Decoder instance
	Decoder interface {
		events.EventEmitter

		Add(any) error
		Destroy()
	}

	Parser interface {
		// A socket.io Encoder instance
		Encoder() Encoder

		// A socket.io Decoder instance
		Decoder() Decoder
	}

	parser struct {
	}
)

func (p *parser) Encoder() Encoder {
	return NewEncoder()
}
func (p *parser) Decoder() Decoder {
	return NewDecoder()
}

func NewParser() Parser {
	return &parser{}
}
