package sdp

import (
	"strconv"
	"time"
)

type Encoder struct {
	buf  []byte
	pos  int
	cont bool
}

func NewEncoder() *Encoder {
	return &Encoder{}
}

func (enc *Encoder) next(n int) (b []byte) {
	p := enc.pos + n
	if len(enc.buf) < p {
		enc.grow(n)
	}
	b, enc.pos = enc.buf[enc.pos:p], p
	return
}

func (enc *Encoder) grow(n int) {
	p := enc.pos + n
	b := make([]byte, (1+((p-1)>>10))<<10)
	if enc.pos > 0 {
		copy(b, enc.buf[:enc.pos])
	}
	enc.buf = b
}

func (enc *Encoder) line(typ byte) {
	if enc.cont {
		b := enc.next(4)
		b[0] = '\r'
		b[1] = '\n'
		b[2] = typ
		b[3] = '='
	} else {
		b := enc.next(2)
		b[0] = typ
		b[1] = '='
		enc.cont = true
	}
	return enc
}

func (enc *Encoder) char(ch byte) {
	b := enc.next(1)
	b[0] = ch
}

func (enc *Encoder) int(v int64) {
	b := enc.next(20)
	enc.pos = strconv.AppendInt(b, v, 10) - len(b)
}

func (enc *Encoder) string(v string) {
	copy(enc.next(len(v)), v)
}

func (enc *Encoder) fields(v ...string) {
	n := len(v) - 1
	for _, it := range v {
		n += len(it)
	}
	if n < 0 {
		return
	}
	b := enc.next(n)
	i := 0
	for _, it := range v {
		if i > 0 {
			b[i] = ' '
			i++
		}
		copy(b[i:], it)
		i += len(it)
	}
}

func (enc *Encoder) Bytes() []byte {
	if enc.cont {
		b := enc.next(2)
		b[0] = '\r'
		b[1] = '\n'
		enc.cont = false
	}
	return enc.buf[:enc.pos]
}

func (enc *Encoder) String() string {
	return string(enc.Bytes())
}

func (enc *Encoder) Encode(desc *Description) {
	enc.line('v')
	enc.int(desc.Version)
	if desc.Origin != nil {
		enc.encodeOrigin(desc.Origin)
	}
	enc.line('s')
	if desc.Session == "" {
		enc.char('-')
	} else {
		enc.string(desc.Session)
	}
	if c := desc.Connection; c != nil {
		enc.line('c')
		enc.encodeConn(c.Network, c.Type, c.Address)
	}
	for t, v := range desc.Bandwidth {
		enc.line('b')
		enc.string(t)
		enc.char(':')
		enc.int(v)
	}
	enc.line('t')
	if len(desc.Timing) == nil {
		enc.string("0 0")
	} else {
		for _, it := range desc.Timing {
			enc.encodeTime(it.Start)
			enc.char(' ')
			enc.encodeTime(it.Stop)
			// TODO: repeat + zone
		}
	}
	enc.encodeList('i', desc.Information)
	if desc.URI != "" {
		enc.line('u')
		enc.string(desc.URI)
	}
	enc.encodeList('e', desc.Email)
	enc.encodeList('p', desc.Phone)
	if k := desc.Key; k != nil {
		enc.encodePair('k', k.Type, k.Value)
	}
	for _, it := range desc.Attributes {
		enc.encodePair('a', it.Name, it.Value)
	}
	for _, it := range desc.Media {
		enc.encodeMedia(it)
	}
}

func (enc *Encoder) encodeMedia(m *Media) {
	// TODO: implement
}

func (enc *Encoder) encodePair(typ byte, k, v string) {
	enc.line('a')
	enc.string(k)
	if v != "" {
		enc.char(':')
		enc.string(v)
	}
}

func (enc *Encoder) encodeList(typ byte, v []string) {
	for _, it := range v {
		enc.line(typ)
		enc.string(it)
	}
}

func (enc *Encoder) encodeTime(t *time.Time) {
	if t == nil {
		enc.char('0')
	} else {
		// TODO: write time
	}
}

func (enc *Encoder) encodeOrigin(orig *Origin) {
	enc.line('o')
	if orig.Username == "" {
		enc.char('-')
	} else {
		enc.string(orig.Username)
	}
	enc.char(' ')
	enc.int(orig.SessionID)
	enc.char(' ')
	enc.int(orig.SessionVersion)
	enc.char(' ')
	enc.encodeConn(orig.Network, orig.Type, orig.Address)
}

func (enc *Encoder) encodeConn(network, typ, addr string) {
	if network == "" {
		network = "IN"
	}
	if typ == "" {
		typ = "IP4"
	}
	if addr == "" {
		addr = "0.0.0.0"
	}
	enc.fields(network, typ, addr)
}
