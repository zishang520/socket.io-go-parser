package parser

type (
	PacketType byte

	Packet struct {
		Type        PacketType `json:"type" msgpack:"type"`
		Nsp         string     `json:"nsp" msgpack:"nsp"`
		Data        any        `json:"data,omitempty" msgpack:"data,omitempty"`
		Id          *uint64    `json:"id,omitempty" msgpack:"id,omitempty"`
		Attachments *uint64    `json:"attachments,omitempty" msgpack:"attachments,omitempty"`
	}
)

const (
	CONNECT       PacketType = '0'
	DISCONNECT    PacketType = '1'
	EVENT         PacketType = '2'
	ACK           PacketType = '3'
	CONNECT_ERROR PacketType = '4'
	BINARY_EVENT  PacketType = '5'
	BINARY_ACK    PacketType = '6'

	// Type of packet in doc of socket io: https://socket.io/docs/v4/socket-io-protocol/
	CONNECT_INT       = 0
	DISCONNECT_INT    = 1
	EVENT_INT         = 2
	ACK_INT           = 3
	CONNECT_ERROR_INT = 4
	BINARY_EVENT_INT  = 5
	BINARY_ACK_INT    = 6
)

func (t PacketType) Valid() bool {
	return t >= '0' && t <= '6'
}

func TransferType(t PacketType) PacketType {
	switch t {
	case CONNECT_INT:
		return CONNECT
	case DISCONNECT_INT:
		return DISCONNECT
	case EVENT_INT:
		return EVENT
	case ACK_INT:
		return ACK
	case CONNECT_ERROR_INT:
		return CONNECT_ERROR
	case BINARY_EVENT_INT:
		return BINARY_EVENT
	case BINARY_ACK_INT:
		return BINARY_ACK
	}
	return t
}

func (t PacketType) String() string {
	switch t {
	case CONNECT, CONNECT_INT:
		return "CONNECT"
	case DISCONNECT, DISCONNECT_INT:
		return "DISCONNECT"
	case EVENT, EVENT_INT:
		return "EVENT"
	case ACK, ACK_INT:
		return "ACK"
	case CONNECT_ERROR, CONNECT_ERROR_INT:
		return "CONNECT_ERROR"
	case BINARY_EVENT, BINARY_EVENT_INT:
		return "BINARY_EVENT"
	case BINARY_ACK, BINARY_ACK_INT:
		return "BINARY_ACK"
	}
	return "UNKNOWN"
}
