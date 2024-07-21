package parser

import (
	"bytes"
	"errors"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/mitchellh/mapstructure"
	"github.com/zishang520/engine.io-go-parser/types"
)

type Placeholder struct {
	Placeholder bool `json:"_placeholder" mapstructure:"_placeholder" msgpack:"_placeholder"`
	Num         int  `json:"num" mapstructure:"num" msgpack:"num"`
}

func init() {
	jsoniter.RegisterTypeEncoderFunc("types.BytesBuffer", func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
		bb := ((*types.BytesBuffer)(ptr))

		bufList := stream.Attachment.([]types.BufferInterface)
		_placeholder := &Placeholder{Placeholder: true, Num: len(bufList)}
		stream.WriteVal(_placeholder)
		stream.Attachment = append(bufList, bb)
	}, nil)

	jsoniter.RegisterTypeEncoderFunc("[]byte", func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
		bb := types.NewBytesBuffer(nil)
		barr := ((*[]byte)(ptr))
		bb.Write(*barr)

		bufList := stream.Attachment.([]types.BufferInterface)
		_placeholder := &Placeholder{Placeholder: true, Num: len(bufList)}
		stream.WriteVal(_placeholder)
		stream.Attachment = append(bufList, bb)
	}, nil)
}

// Replaces every io.Reader | []byte in packet with a numbered placeholder.
func DeconstructPacket(packet *Packet) (pack *Packet, buffers []types.BufferInterface) {
	pack = packet

	// Run the serialization now, replacing any bytebuffers/[]byte found along the way with placeholders
	buf := &bytes.Buffer{}
	ns := jsoniter.NewStream(jsoniter.ConfigDefault, buf, buf.Cap())
	ns.Attachment = buffers
	ns.WriteVal(pack.Data)
	buffers = ns.Attachment.([]types.BufferInterface)
	ns.Flush()
	pack.Data = buf.String()

	attachments := uint64(len(buffers))
	pack.Attachments = &attachments // number of binary 'attachments'
	return pack, buffers
}

// Reconstructs a binary packet from its placeholder packet and buffers
func ReconstructPacket(packet *Packet, buffers []types.BufferInterface) (*Packet, error) {
	data, err := _reconstructPacket(packet.Data, &buffers)
	if err != nil {
		return nil, err
	}
	packet.Data = data
	packet.Attachments = nil // Attachments are no longer needed
	return packet, nil
}

func _reconstructPacket(data any, buffers *[]types.BufferInterface) (any, error) {
	switch d := data.(type) {
	case nil:
		return nil, nil
	case []any:
		newData := make([]any, 0, len(d))
		for _, v := range d {
			_data, err := _reconstructPacket(v, buffers)
			if err != nil {
				return nil, err
			}
			newData = append(newData, _data)
		}
		return newData, nil
	case map[string]any:
		var _placeholder Placeholder
		if mapstructure.Decode(d, &_placeholder) == nil && _placeholder.Placeholder {
			if _placeholder.Num >= 0 && _placeholder.Num < len(*buffers) {
				return (*buffers)[_placeholder.Num], nil // appropriate buffer (should be natural order anyway)
			}
			return nil, errors.New("illegal attachments")
		}
		newData := make(map[string]any, len(d))
		for k, v := range d {
			_data, err := _reconstructPacket(v, buffers)
			if err != nil {
				return nil, err
			}
			newData[k] = _data
		}
		return newData, nil
	default:
		return data, nil
	}
}
