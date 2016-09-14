package sdp

import (
	"strings"
	"testing"
	"time"
)

func TestEncode(t *testing.T) {
	t.Parallel()

	layout := "2006-01-02 15:04:05 -0700 MST"
	start, _ := time.Parse(layout, "1996-02-27 15:26:59 +0000 UTC")
	stop, _ := time.Parse(layout, "1996-05-30 16:26:59 +0000 UTC")

	desc := &Description{
		Origin: &Origin{
			Username:       "jdoe",
			SessionID:      2890844526,
			SessionVersion: 2890842807,
			Network:        "IN",
			Type:           "IP4",
			Address:        "10.47.16.5",
		},
		Session:     "SDP Seminar",
		Information: "A Seminar on the session description protocol",
		URI:         "http://www.example.com/seminars/sdp.pdf",
		Email:       []string{"j.doe@example.com (Jane Doe)"},
		Phone:       []string{"+1 617 555-6011"},
		Connection: &Connection{
			Network: "IN",
			Type:    "IP4",
			Address: "224.2.17.12/127",
		},
		Bandwidth: map[string]int{
			"AS": 2000,
		},
		Timing: &Timing{
			Start: start,
			Stop:  stop,
			Repeat: &Repeat{
				Interval: time.Duration(604800) * time.Second,
				Duration: time.Duration(3600) * time.Second,
				Offsets: []time.Duration{
					time.Duration(0),
					time.Duration(90000) * time.Second,
				},
			},
		},
		TimeZones: []*TimeZone{
			{Time: start, Offset: -time.Hour},
			{Time: stop, Offset: 0},
		},
		Media: []*Media{
			{
				Type:  "audio",
				Port:  49170,
				Proto: "RTP/AVP",
				Formats: map[int]*Format{
					0: nil,
				},
			},
			{
				Type:  "video",
				Port:  51372,
				Proto: "RTP/AVP",
				Formats: map[int]*Format{
					99: {Payload: 99, Codec: "h263-1998", Clock: 90000},
				},
			},
		},
		Mode: ModeRecvOnly,
	}

	expected := `v=0
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
m=video 51372 RTP/AVP 99
a=rtpmap:99 h263-1998/90000`

	a := strings.Split(desc.String(), "\r\n")
	b := strings.Split(expected, "\n")
	if len(b) > len(a) {
		t.Fail()
	}
	for i, exp := range b {
		if a[i] != exp {
			t.Fatalf("found '%s' expected '%s'", a[i], exp)
		}
	}
}

func TestReadmeExampleEncode(t *testing.T) {
	t.Parallel()
	desc := &Description{
		Origin: &Origin{
			Username:       "alice",
			Address:        "alice.example.org",
			SessionID:      2890844526,
			SessionVersion: 2890844526,
		},
		Session:    "Example",
		Connection: &Connection{Address: "127.0.0.1"},
		Media: []*Media{
			{
				Type:  "audio",
				Port:  10000,
				Proto: "RTP/AVP",
				Formats: map[int]*Format{
					0: {Payload: 0, Codec: "PCMU", Clock: 8000},
					8: {Payload: 8, Codec: "PCMA", Clock: 8000},
				},
			},
		},
		Mode: ModeSendRecv,
	}
	if len(desc.String()) != 182 {
		t.Fail()
	}
}
