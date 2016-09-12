# Golang: SDP Protocol

[![Build Status](https://travis-ci.org/pixelbender/go-sdp.svg)](https://travis-ci.org/pixelbender/go-sdp)
[![Coverage Status](https://coveralls.io/repos/github/pixelbender/go-sdp/badge.svg?branch=master)](https://coveralls.io/github/pixelbender/go-sdp?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/pixelbender/go-sdp)](https://goreportcard.com/report/github.com/pixelbender/go-sdp)
[![GoDoc](https://godoc.org/github.com/pixelbender/go-sdp?status.svg)](https://godoc.org/github.com/pixelbender/go-sdp)

## Features

- [x] SDP Decoder/Encoder
- [ ] SDP Answer/Offer Negotiation
- [ ] SDP Media Capabilities Negotiation

## Installation

```sh
go get github.com/pixelbender/go-sdp
```

## SDP Decoding

```go
package main

import (
	"github.com/pixelbender/go-sdp/sdp"
	"fmt"
)

func main() {
	desc, err := sdp.Parse(`v=0
o=alice 2890844526 2890844526 IN IP4 alice.example.org
s=
c=IN IP4 127.0.0.1
t=0 0
m=audio 10000 RTP/AVP 0 8
a=rtpmap:0 PCMU/8000
a=rtpmap:8 PCMA/8000`)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(desc)
	}
}
```

## SDP Encoding

```go
package main

import (
	"github.com/pixelbender/go-sdp/sdp"
	"fmt"
)

func main() {
	desc := &sdp.Description{
	    Origin: &sdp.Origin{ Username:"alice", Address:"alice.example.org" },
	    Session: "Example",
	    Connection: &sdp.Connection{ Address: "127.0.0.1" },
	    Media: []*sdp.Media{
	        &sdp.Media{
	            Type: "audio",
	            Port: 10000,
	            Proto: "RTP/AVP",
	            Formats: []*sdp.Format{
	                &sdp.Format{ Payload: 0, Codec: "PCMU", Clock: 8000 },
	                &sdp.Format{ Payload: 8, Codec: "PCMA", Clock: 8000 },
	            },
	            Attributes: []*sdp.Attribute{
	                &sdp.Attribute{Name:"sendonly"},
	            },
	        },
	    },
	}

	fmt.Println(desc.String())
}
```

## Specifications

- [RFC 5389: Session Description Protocol](https://tools.ietf.org/html/rfc4566)
- [RFC 3264: Offer/Answer Model with SDP](https://tools.ietf.org/html/rfc3264)
- [RFC 6871: SDP Media Capabilities Negotiation](https://tools.ietf.org/html/rfc6871)
