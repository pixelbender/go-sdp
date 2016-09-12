package sdp

import (
	"encoding/json"
	"testing"
)

func TestReadmeDecode(t *testing.T) {
	descr, err := Parse(`v=0
o=alice 2890844526 2890844526 IN IP4 127.0.0.1
s=
c=IN IP4 127.0.0.1
t=0 0
m=audio 49170 RTP/AVP 0 8
a=rtpmap:0 PCMU/8000
a=rtpmap:8 PCMA/8000`)

	if err != nil {
		t.Fatal("decode", err)
	}

	b, _ := json.Marshal(descr)
	t.Logf("%s", b)
}
