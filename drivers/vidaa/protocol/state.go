package protocol

type ConnState int

const (
	StateDisconnected ConnState = iota
	StateConnecting
	StateTLSEstablished
	StateAuthRequired
	StatePairingPending
	StateAuthenticated
	StateReconnecting
	StateFailed
)

type FSM struct {
	s ConnState
}

func NewFSM() *FSM { return &FSM{s: StateDisconnected} }

func (f *FSM) State() ConnState { return f.s }

func (f *FSM) Transition(to ConnState) bool {
	switch f.s {
	case StateDisconnected:
		if to == StateConnecting {
			f.s = to
			return true
		}
	case StateConnecting:
		if to == StateTLSEstablished || to == StateFailed {
			f.s = to
			return true
		}
	case StateTLSEstablished:
		if to == StateAuthRequired || to == StateFailed {
			f.s = to
			return true
		}
	case StateAuthRequired:
		if to == StatePairingPending || to == StateAuthenticated || to == StateFailed {
			f.s = to
			return true
		}
	case StatePairingPending:
		if to == StateAuthenticated || to == StateFailed {
			f.s = to
			return true
		}
	case StateAuthenticated:
		if to == StateReconnecting || to == StateDisconnected || to == StateFailed {
			f.s = to
			return true
		}
	case StateReconnecting:
		if to == StateConnecting || to == StateFailed {
			f.s = to
			return true
		}
	default:
		return false
	}
	return false
}
