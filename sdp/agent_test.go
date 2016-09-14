package sdp

import (
	"testing"
)

// Offer/Answer Examples from RFC 4317

func TestAudioVideo1(t *testing.T) {
	//	offer, err := Parse(`v=0
	//o=alice 2890844526 2890844526 IN IP4 host.atlanta.example.com
	//s=
	//c=IN IP4 host.atlanta.example.com
	//t=0 0
	//m=audio 49170 RTP/AVP 0 8 97
	//a=rtpmap:0 PCMU/8000
	//a=rtpmap:8 PCMA/8000
	//a=rtpmap:97 iLBC/8000
	//m=video 51372 RTP/AVP 31 32
	//a=rtpmap:31 H261/90000
	//a=rtpmap:32 MPV/90000`)
	//
	//	if err != nil {
	//		t.Fail()
	//	}
	//
	//	expected := `v=0
	//o=bob 2808844564 2808844564 IN IP4 host.biloxi.example.com
	//s=
	//c=IN IP4 host.biloxi.example.com
	//t=0 0
	//m=audio 49174 RTP/AVP 0
	//a=rtpmap:0 PCMU/8000
	//m=video 49170 RTP/AVP 32
	//a=rtpmap:32 MPV/90000`
	//
	//	answer := NegotiateMedia(func(local *Media, remote *Media) *Media {
	//		for _, it := range remote.Formats {
	//			if it.Codec == "PCMU" {
	//				local.Port = 49174
	//				local.Proto = remote.Proto
	//				local.Formats = append(local.Formats, it)
	//			}
	//		}
	//		return f.Codec == "PCMU" || f.Codec == "MPV"
	//	})
	//
	//	a := strings.Split(answer.String(), "\r\n")
	//	b := strings.Split(expected, "\n")
	//	if len(b) > len(a) {
	//		t.Fail()
	//	}
	//	for i, exp := range b {
	//		if a[i] != exp {
	//			t.Fatalf("found '%s' expected '%s'", a[i], exp)
	//		}
	//	}
}
