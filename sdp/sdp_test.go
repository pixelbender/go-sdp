package sdp

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

var testVectors = []*testVector{
	{
		Name: "RFC4566 Example",
		Data: `v=0
o=jdoe 2890844526 2890842807 IN IP4 10.47.16.5
s=SDP Seminar
i=A Seminar on the session description protocol
u=http://www.example.com/seminars/sdp.pdf
e=j.doe@example.com (Jane Doe)
p=+1 617 555-6011
c=IN IP4 224.2.17.12/127
b=AS:2000
t=3034423619 3042462419
r=7d 1h 0 25h
z=3034423619 -1h 3042462419 0
a=recvonly
m=audio 49170 RTP/AVP 0
m=video 51372 RTP/AVP 99 100
a=rtpmap:99 h263-1998/90000
a=rtpmap:100 H264/90000
a=rtcp-fb:100 ccm fir
a=rtcp-fb:100 nack
a=rtcp-fb:100 nack pli
a=fmtp:100 profile-level-id=42c01f;level-asymmetry-allowed=1
`,
		Session: &Session{
			Origin: &Origin{
				Username:       "jdoe",
				SessionID:      2890844526,
				SessionVersion: 2890842807,
				Network:        NetworkInternet,
				Type:           TypeIPv4,
				Address:        "10.47.16.5",
			},
			Name:        "SDP Seminar",
			Information: "A Seminar on the session description protocol",
			URI:         "http://www.example.com/seminars/sdp.pdf",
			Email:       []string{"j.doe@example.com (Jane Doe)"},
			Phone:       []string{"+1 617 555-6011"},
			Connection: &Connection{
				Network: NetworkInternet,
				Type:    TypeIPv4,
				Address: "224.2.17.12",
				TTL:     127,
			},
			Bandwidth: []*Bandwidth{
				{"AS", 2000},
			},
			Timing: &Timing{
				Start: parseTime("1996-02-27 15:26:59 +0000 UTC"),
				Stop:  parseTime("1996-05-30 16:26:59 +0000 UTC"),
			},
			Repeat: []*Repeat{
				{
					Interval: time.Duration(604800) * time.Second,
					Duration: time.Duration(3600) * time.Second,
					Offsets: []time.Duration{
						time.Duration(0),
						time.Duration(90000) * time.Second,
					},
				},
			},
			TimeZone: []*TimeZone{
				{Time: parseTime("1996-02-27 15:26:59 +0000 UTC"), Offset: -time.Hour},
				{Time: parseTime("1996-05-30 16:26:59 +0000 UTC"), Offset: 0},
			},
			Mode: RecvOnly,
			Media: []*Media{
				{
					Type:  "audio",
					Port:  49170,
					Proto: "RTP/AVP",
					Format: []*Format{
						{Payload: 0},
					},
				},
				{
					Type:  "video",
					Port:  51372,
					Proto: "RTP/AVP",
					Format: []*Format{
						{Payload: 99, Name: "h263-1998", ClockRate: 90000},
						{Payload: 100, Name: "H264", ClockRate: 90000, Params: []string{
							"profile-level-id=42c01f;level-asymmetry-allowed=1",
						}, Feedback: []string{
							"ccm fir", "nack", "nack pli",
						}},
					},
				},
			},
		},
	},
	{
		Name: "Readme Example",
		Data: `v=0
o=alice 2890844526 2890844526 IN IP4 alice.example.org
s=Example
c=IN IP4 127.0.0.1
t=0 0
a=sendrecv
m=audio 10000 RTP/AVP 0 8
a=rtpmap:0 PCMU/8000
a=rtpmap:8 PCMA/8000
`,
		Session: &Session{
			Origin: &Origin{
				Username:       "alice",
				SessionID:      2890844526,
				SessionVersion: 2890844526,
				Network:        NetworkInternet,
				Type:           TypeIPv4,
				Address:        "alice.example.org",
			},
			Name: "Example",
			Connection: &Connection{
				Network: NetworkInternet,
				Type:    TypeIPv4,
				Address: "127.0.0.1",
			},
			Media: []*Media{
				{
					Type:  "audio",
					Port:  10000,
					Proto: "RTP/AVP",
					Format: []*Format{
						{Payload: 0, Name: "PCMU", ClockRate: 8000},
						{Payload: 8, Name: "PCMA", ClockRate: 8000},
					},
				},
			},
			Mode: SendRecv,
		},
	},
	{
		Name: "SCTP Example",
		Data: `v=0
o=- 0 2 IN IP4 127.0.0.1
s=-
c=IN IP4 127.0.0.1
t=0 0
m=application 10000 DTLS/SCTP 5000
a=sctpmap:5000 webrtc-datachannel 256
m=application 10000 UDP/DTLS/SCTP webrtc-datachannel
a=sctp-port:5000
`,
		Session: &Session{
			Origin: &Origin{
				Username:       "-",
				SessionID:      0,
				SessionVersion: 2,
				Network:        NetworkInternet,
				Type:           TypeIPv4,
				Address:        "127.0.0.1",
			},
			Name: "-",
			Connection: &Connection{
				Network: NetworkInternet,
				Type:    TypeIPv4,
				Address: "127.0.0.1",
			},
			Media: []*Media{
				{
					Type:        "application",
					Port:        10000,
					Proto:       "DTLS/SCTP",
					FormatDescr: "5000",
					Attributes: Attributes{
						{"sctpmap", "5000 webrtc-datachannel 256"},
					},
				},
				{
					Type:        "application",
					Port:        10000,
					Proto:       "UDP/DTLS/SCTP",
					FormatDescr: "webrtc-datachannel",
					Attributes: Attributes{
						{"sctp-port", "5000"},
					},
				},
			},
		},
	},
}

type testVector struct {
	Name string
	Data string
	*Session
}

func TestVectors(t *testing.T) {
	for _, v := range testVectors {
		v := v
		t.Run(v.Name, func(inner *testing.T) {
			t := &T{inner}
			sess, err := ParseString(v.Data)
			if err != nil {
				t.Fatal(err)
			}
			t.AssertAny("decoded", sess, v.Session)
			t.AssertAny("encoded", strings.Split(v.Session.String(), "\r\n"), strings.Split(v.Data, "\n"))
		})
	}
}

type T struct {
	*testing.T
}

func (t *T) AssertAny(name string, got, exp interface{}) {
	a := reflect.ValueOf(got)
	b := reflect.ValueOf(exp)
	r := a.Type()
	t.Assert(name+" type", r, b.Type())

	switch r.Kind() {
	case reflect.Ptr:
		t.Assert(name+" ptr", a.IsNil(), b.IsNil())
		if a.IsNil() {
			break
		}
		t.AssertAny(name, a.Elem().Interface(), b.Elem().Interface())
	case reflect.Struct:
		switch r {
		case timeType:
			t.Assert(name, got, exp)
		default:
			for i := 0; i < r.NumField(); i++ {
				t.AssertAny(
					fmt.Sprintf("%s %s", name, strings.ToLower(r.Field(i).Name)),
					a.Field(i).Interface(),
					b.Field(i).Interface(),
				)
			}
		}
	case reflect.Slice:
		for i := 0; i < a.Len(); i++ {
			n := fmt.Sprintf("%s at index %d", name, i)
			if i < b.Len() {
				t.AssertAny(
					n,
					a.Index(i).Interface(),
					b.Index(i).Interface(),
				)
			} else {
				t.Fatalf("unexpected %s: %s", n, dump(a.Index(i).Interface()))
			}
		}
	default:
		t.Assert(name, got, exp)
	}
}

func (t *T) Assert(name string, got, exp interface{}) {
	if reflect.DeepEqual(got, exp) {
		return
	}
	t.Fatalf("bad %s, got: %s, expected: %s", name, dump(got), dump(exp))
}

func BenchmarkDecode(b *testing.B) {
	for _, v := range testVectors {
		v := v
		b.Run(v.Name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if _, err := ParseString(v.Data); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkDecodeReader(b *testing.B) {
	for _, v := range testVectors {
		v := v
		b.Run(v.Name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if _, err := NewDecoder(strings.NewReader(v.Data)).Decode(); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkEncode(b *testing.B) {
	for _, v := range testVectors {
		v := v
		b.Run(v.Name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				e := NewEncoder(nil)
				if err := e.Encode(v.Session); err != nil {
					b.Fatal(err)
				}
				_ = e.Bytes()
			}
		})
	}
}

func BenchmarkEncodeReuse(b *testing.B) {
	for _, v := range testVectors {
		v := v
		b.Run(v.Name, func(b *testing.B) {
			b.ReportAllocs()
			e := NewEncoder(nil)
			for i := 0; i < b.N; i++ {
				if err := e.Encode(v.Session); err != nil {
					b.Fatal(err)
				}
				_ = e.Bytes()
			}
		})
	}
}

var timeType = reflect.TypeOf(time.Time{})

func parseTime(s string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05 -0700 MST", s)
	if err != nil {
		panic(err)
	}
	return t
}

func dump(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}
