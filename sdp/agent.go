package sdp

// Agent represents an SDP offer/answer exchange agent.
type Agent struct {
	local  *descr
	remote *descr
}

// SetLocalDescription changes the local description.
func (s *Agent) SetLocalDescription(typ string, offer *Description) {
}

// SetRemoteDescription changes the remote description.
func (s *Agent) SetRemoteDescription(typ string, offer *Description) {
}

type descr struct {
	*Description
	typ string
}
