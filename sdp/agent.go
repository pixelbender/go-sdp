package sdp

const (
	TypeOffer = "offer"
)

// Agent represents an SDP offer/answer exchange agent.
type Agent struct {
	local  *descr
	remote *descr
}

func (s *Agent) SetLocalDescription(typ string, offer *Description) {
}

func (s *Agent) SetRemoteDescription(typ string, offer *Description) {
}

type descr struct {
	*Description
	typ string
}
