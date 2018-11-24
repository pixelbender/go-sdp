package sdp

import (
	"strings"
	"testing"
	"time"
)

const (
	seminarDescr = `v=0
o=jdoe 2890844526 2890842807 IN IP4 10.47.16.5
s=SDP Seminar
i=A Seminar on the session description protocol
u=http://www.example.com/seminars/sdp.pdf
e=j.doe@example.com (Jane Doe)
p=+1 617 555-6011
c=IN IP4 224.2.17.12/127
b=AS:2000
z=3034423619 -1h 3042462419 0
a=recvonly
t=3034423619 3042462419
r=7d 1h 0 25h
m=audio 49170 RTP/AVP 0
m=video 51372 RTP/AVP 99 100
a=rtpmap:99 h263-1998/90000
a=rtpmap:100 H264/90000
a=rtcp-fb:100 ccm fir
a=rtcp-fb:100 nack
a=rtcp-fb:100 nack pli
a=fmtp:100 profile-level-id=42c01f;level-asymmetry-allowed=1
`
	readmeDescr = `v=0
o=alice 2890844526 2890844526 IN IP4 alice.example.org
s=Example
c=IN IP4 127.0.0.1
a=sendrecv
t=0 0
m=audio 10000 RTP/AVP 0 8
a=rtpmap:0 PCMU/8000
a=rtpmap:8 PCMA/8000
`
)

func BenchmarkDecode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := ParseString(seminarDescr)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeReader(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := NewDecoder(strings.NewReader(seminarDescr)).Decode()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncode(b *testing.B) {
	b.ReportAllocs()
	sess, err := ParseString(seminarDescr)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		new(Encoder).Encode(sess)
	}
}

func BenchmarkEncodeReuse(b *testing.B) {
	b.ReportAllocs()
	sess, err := ParseString(seminarDescr)
	if err != nil {
		b.Fatal(err)
	}
	e := new(Encoder)
	for i := 0; i < b.N; i++ {
		e.Encode(sess)
	}
}

func TestReadmeExample(t *testing.T) {
	t.Parallel()

	sess := &Session{
		Origin: &Origin{
			Username:       "alice",
			Address:        "alice.example.org",
			SessionID:      2890844526,
			SessionVersion: 2890844526,
		},
		Name: "Example",
		Connection: &Connection{
			Address: "127.0.0.1",
		},
		Media: []*Media{
			{
				Type:  "audio",
				Port:  10000,
				Proto: "RTP/AVP",
				Formats: []*Format{
					{Payload: 0, Name: "PCMU", ClockRate: 8000},
					{Payload: 8, Name: "PCMA", ClockRate: 8000},
				},
			},
		},
		Mode: ModeSendRecv,
	}

	expected, err := ParseString(readmeDescr)
	if err != nil {
		t.Fatal(err)
	}
	assert(t, expected, sess)
}

func TestSeminarExample(t *testing.T) {
	t.Parallel()

	layout := "2006-01-02 15:04:05 -0700 MST"
	start, _ := time.Parse(layout, "1996-02-27 15:26:59 +0000 UTC")
	stop, _ := time.Parse(layout, "1996-05-30 16:26:59 +0000 UTC")

	sess := &Session{
		Origin: &Origin{
			Username:       "jdoe",
			SessionID:      2890844526,
			SessionVersion: 2890842807,
			Address:        "10.47.16.5",
		},
		Name:        "SDP Seminar",
		Information: "A Seminar on the session description protocol",
		URI:         "http://www.example.com/seminars/sdp.pdf",
		Email:       []string{"j.doe@example.com (Jane Doe)"},
		Phone:       []string{"+1 617 555-6011"},
		Connection: &Connection{
			Address: "224.2.17.12",
			TTL:     127,
		},
		Bandwidth: Bandwidth{
			"AS": 2000,
		},
		Timing: &Timing{
			Start: start,
			Stop:  stop,
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
			{Time: start, Offset: -time.Hour},
			{Time: stop, Offset: 0},
		},
		Mode: ModeRecvOnly,
		Media: []*Media{
			{
				Type:  "audio",
				Port:  49170,
				Proto: "RTP/AVP",
				Formats: []*Format{
					{Payload: 0},
				},
			},
			{
				Type:  "video",
				Port:  51372,
				Proto: "RTP/AVP",
				Formats: []*Format{
					{Payload: 99, Name: "h263-1998", ClockRate: 90000},
					{Payload: 100, Name: "H264", ClockRate: 90000, Params: []string{
						"profile-level-id=42c01f;level-asymmetry-allowed=1",
					}, Feedback: []string{
						"ccm fir", "nack", "nack pli",
					}},
				},
			},
		},
	}

	expected, err := ParseString(seminarDescr)
	if err != nil {
		t.Fatal(err)
	}
	assert(t, expected, sess)
}

func assert(t *testing.T, expected, result *Session) {
	r := strings.Split(result.String(), "\r\n")
	e := strings.Split(expected.String(), "\r\n")
	for i, it := range r {
		if i < len(e) {
			if e[i] != it {
				t.Fatalf("result line %d: '%s', expected: '%s'", i, it, e[i])
			}
			continue
		}
		t.Fatalf("unexpected line %d: '%s'", i, it)
	}
	if len(r) != len(e) {
		t.Errorf("wrong number of lines %d, expected %d", len(r), len(e))
	}
}
