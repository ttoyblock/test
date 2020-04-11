package ws

type StateType int32

const (
	STATE_INIT         StateType = 0
	STATE_TCP_READY    StateType = 1
	STATE_SYN_RECEIVED StateType = 2
	STATE_ESTABLISHED  StateType = 3
	STATE_CLOSED       StateType = 4
)

func (st StateType) String() (s string) {
	switch st {
	case STATE_INIT:
		s = "INIT"
	case STATE_TCP_READY:
		s = "TCP_READY"
	case STATE_SYN_RECEIVED:
		s = "SYN_RECEIVED"
	case STATE_ESTABLISHED:
		s = "ESTABLISHED"
	case STATE_CLOSED:
		s = "CLOSED"
	default:
		s = "UNKNOWN"
	}
	return
}
