package sdp

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestDecode(t *testing.T) {
	t.Parallel()

	sdp := `v=0
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
k=prompt
a=recvonly
m=audio 49170 RTP/AVP 0
m=video 51372 RTP/AVP 99
a=rtpmap:99 h263-1998/90000`

	layout := "2006-01-02 15:04:05 -0700 MST"
	start, _ := time.Parse(layout, "1996-02-27 15:26:59 +0000 UTC")
	stop, _ := time.Parse(layout, "1996-05-30 16:26:59 +0000 UTC")

	for _, dec := range []*Decoder{
		NewDecoder(strings.NewReader(sdp)),
		NewDecoderString(sdp),
	} {

		desc, err := dec.Decode()
		if err != nil {
			t.Fatal(err)
			return
		}
		audio, video := desc.Media[0], desc.Media[1]

		for _, it := range []struct{ found, expected interface{} }{
			{desc.Version, 0},
			{desc.Origin, &Origin{
				Username:       "jdoe",
				SessionID:      2890844526,
				SessionVersion: 2890842807,
				Network:        "IN",
				Type:           "IP4",
				Address:        "10.47.16.5",
			}},
			{desc.Session, "SDP Seminar"},
			{desc.Information, "A Seminar on the session description protocol"},
			{desc.URI, "http://www.example.com/seminars/sdp.pdf"},
			{desc.Email, []string{"j.doe@example.com (Jane Doe)"}},
			{desc.Phone, []string{"+1 617 555-6011"}},
			{desc.Bandwidth["AS"], 2000},
			{desc.Mode, ModeRecvOnly},
			{desc.Connection, &Connection{
				Network: "IN",
				Type:    "IP4",
				Address: "224.2.17.12/127",
			}},
			{desc.Timing.Start, start},
			{desc.Timing.Stop, stop},
			{desc.Timing.Repeat, &Repeat{
				Interval: time.Duration(604800) * time.Second,
				Duration: time.Duration(3600) * time.Second,
				Offsets: []time.Duration{
					time.Duration(0),
					time.Duration(90000) * time.Second,
				},
			}},
			{desc.TimeZones, []*TimeZone{
				{Time: start, Offset: -time.Hour},
				{Time: stop, Offset: 0},
			}},
			{desc.Key, &Key{Type: "prompt"}},
			{audio.Type, "audio"},
			{audio.Port, 49170},
			{audio.Proto, "RTP/AVP"},
			{audio.Formats[0], &Format{Payload: 0}},
			{video.Type, "video"},
			{video.Port, 51372},
			{video.Proto, "RTP/AVP"},
			{video.Formats[99], &Format{Payload: 99, Codec: "h263-1998", Clock: 90000}},
		} {
			if !reflect.DeepEqual(it.found, it.expected) {
				t.Fatalf("found %+v expected %+v", it.found, it.expected)
				return
			}
		}
	}
}

func TestReadmeExampleDecode(t *testing.T) {
	t.Parallel()
	desc, err := Parse(`v=0
o=alice 2890844526 2890844526 IN IP4 alice.example.org
s=
c=IN IP4 127.0.0.1
t=0 0
a=sendrecv
m=audio 10000 RTP/AVP 0 8
a=rtpmap:0 PCMU/8000
a=rtpmap:8 PCMA/8000`)
	if err != nil {
		t.Fatal(err)
	}
	if desc.Media[0].Formats[0].Codec != "PCMU" {
		t.Fail()
	}
}
