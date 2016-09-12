package sdp

import "time"

// MediaType is the media type for an SDP session description.
const MediaType = "application/sdp"

// Description represents an SDP session description. RFC 4566 Section 5.
type Description struct {
	Version     int          // Protocol Version ("v=")
	Origin      *Origin      // Origin ("o=")
	Session     string       // Session Name ("s=")
	Information []string     // Session Information ("i=")
	URI         string       // URI ("u=")
	Email       []string     // Email Address ("e=")
	Phone       []string     // Phone Number ("p=")
	Connection  *Connection  // Connection Data ("c=")
	Bandwidth   Bandwidth    // Bandwidth ("b=")
	Timing      *Timing      // Timing ("t=")
	Attributes  []*Attribute // Attribute ("a=")
	Media       []*Media     // Media Descriptions ("m=")
}

// String returns the encoded session description according the SDP specification.
func (desc *Description) String() string {
	enc := NewEncoder()
	enc.Encode(desc)
	return enc.String()
}

// Parse parses text into a Description structure.
func Parse(text string) (*Description, error) {
	dec := NewDecoderString(text)
	return dec.Decode()
}

// Origin represents an originator of the session. RFC 4566 Section 5.2.
type Origin struct {
	Username       string
	SessionID      int64
	SessionVersion int64
	Network        string
	Type           string
	Address        string
}

// Media contains media description. RFC 4566 Section 5.14.
type Media struct {
	Type       string
	Port       int
	PortNum    int
	Proto      string
	Formats    []*Format
	Attributes []*Attribute
}

// Format is a media format description.
type Format struct {
	Payload  int
	Codec    string
	Clock    int
	Feedback []string
	Params   map[string]string
}

// Attribute represents an a session or media attribute. RFC 4566 Section 5.14.
type Attribute struct {
	Name, Value string
}

// Connection contains connection data. RFC 4566 Section 5.7.
type Connection struct {
	Network string
	Type    string
	Address string
}

// Bandwidth denotes the proposed bandwidth to be used by the session or media. RFC 4566 Section 5.8.
// The bandwidth is interpreted as kilobits per second by default.
type Bandwidth map[string]int64

// Timing specifies start and stop times for a session.
type Timing struct {
	Start   *time.Time
	Stop    *time.Time
	Repeats *Repeats // Repeat Times ("r=")
	// TODO: Add zone adjustments...
}

type Repeats struct {
	Interval time.Duration
	Duration time.Duration
	Offsets  []time.Duration
}
