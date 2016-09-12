package sdp

type Session struct {
	Local  *Description
	Remote *Description
}

func (sess *Session) Negotiate(answer *Description) {

}

func Negotiate(local *Description, remote *Description) *Description {
	orig := local.Origin
	orig.SessionID++
	orig.SessionVersion++

	answer := &Description{
		Version:    remote.Version,
		Origin:     orig,
		Session:    remote.Session,
		Timing:     remote.Timing,
		Connection: local.Connection,
	}
	return answer
}

func NegotiateMedia(local *Media, remote *Media) *Media {
	return nil
}

//
//Version     int         // Protocol Version ("v=")
//Origin      *Origin     // Origin ("o=")
//Session     string      // Session Name ("s=")
//Information []string    // Session Information ("i=")
//URI         string      // URI ("u=")
//Email       []string    // Email Address ("e=")
//Phone       []string    // Phone Number ("p=")
//Connection  *Connection // Connection Data ("c=")
//Bandwidth   Bandwidth   // Bandwidth ("b=")
//Timing      []*Timing   // Timing ("t=")
//Key         *Key
//Attributes  Attributes
//Media       []*Media
