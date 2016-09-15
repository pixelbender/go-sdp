package sdp

import (
	"bufio"
	"io"
	"strconv"
	"time"
)

type reader interface {
	ReadLine() (string, error)
}

// A Decoder reads and decodes SDP description from an input stream or buffer.
type Decoder struct {
	r   reader
	p   []string
	err error
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: &bufferedReader{buf: bufio.NewReader(r)}}
}

// NewDecoderString returns a new decoder that reads from v.
func NewDecoderString(v string) *Decoder {
	return &Decoder{r: &stringReader{buf: v}}
}

// Decode reads and decodes a SDP description from its input.
func (dec *Decoder) Decode() (*Description, error) {
	desc := &Description{}
	for {
		v, err := dec.r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if !dec.split(v, '=', 2, true) {
			return nil, dec.err
		}
		k := dec.p[0]
		if len(k) != 1 {
			return nil, decodeError(v)
		}
		if err = dec.decodeLine(desc, k[0], dec.p[1]); err != nil {
			return nil, err
		}
	}
	if err := dec.decodeSessionAttributes(desc); err != nil {
		return nil, err
	}
	for _, m := range desc.Media {
		if err := dec.decodeMediaAttributes(m); err != nil {
			return nil, err
		}
	}
	return desc, nil
}

func (dec *Decoder) decodeSessionAttributes(desc *Description) error {
	n := 0
	for _, it := range desc.Attributes {
		switch it.Name {
		case ModeSendRecv, ModeRecvOnly, ModeSendOnly, ModeInactive:
			desc.Mode = it.Name
		default:
			desc.Attributes[n] = it
			n++
		}
	}
	desc.Attributes = desc.Attributes[:n]
	return nil
}

func (dec *Decoder) decodeMediaAttributes(m *Media) error {
	n := 0
	for _, it := range m.Attributes {
		switch it.Name {
		case ModeSendRecv, ModeRecvOnly, ModeSendOnly, ModeInactive:
			m.Mode = it.Name
		case "rtpmap":
			dec.decodeMediaMap(m, it.Value)
		case "fmtp":
			dec.decodeMediaParams(m, it.Value)
		default:
			m.Attributes[n] = it
			n++
		}
	}
	m.Attributes = m.Attributes[:n]
	return nil
}

func (dec *Decoder) decodeMediaMap(m *Media, v string) error {
	if !dec.split(v, ' ', 2, true) {
		return dec.err
	}
	p, err := strconv.Atoi(dec.p[0])
	if err != nil {
		return err
	}
	if !dec.split(dec.p[1], '/', 2, true) {
		return dec.err
	}
	f := dec.touchMediaFormat(m, p)
	f.Codec = dec.p[0]
	if dec.split(dec.p[1], '/', 2, false) {
		if f.Channels, err = strconv.Atoi(dec.p[1]); err != nil {
			return err
		}
	}
	if f.Clock, err = strconv.Atoi(dec.p[0]); err != nil {
		return err
	}
	return nil
}

func (dec *Decoder) decodeMediaParams(m *Media, v string) error {
	if !dec.split(v, ' ', 2, true) {
		return dec.err
	}
	p, err := strconv.Atoi(dec.p[0])
	if err != nil {
		return err
	}
	f := dec.touchMediaFormat(m, p)
	f.Params = append(f.Params, dec.p[1])
	return nil
}

func (dec *Decoder) touchMediaFormat(m *Media, p int) *Format {
	if m.Formats == nil {
		m.Formats = make(map[int]*Format)
	}
	f, ok := m.Formats[p]
	if !ok {
		f = &Format{Payload: p}
		m.Formats[p] = f
	}
	return f
}

func (dec *Decoder) split(v string, sep rune, n int, required bool) bool {
	p := dec.p[:0]
	off := 0
	for i, it := range v {
		if it == sep {
			p = append(p, v[off:i])
			off = i + 1
		}
		if len(p)+1 == n {
			dec.p = append(p, v[off:])
			return true
		}
	}
	if required {
		dec.err = decodeError(v)
	} else {
		dec.p = append(p, v[off:])
	}
	return false
}

func (dec *Decoder) parseInt(v string) (int64, error) {
	return strconv.ParseInt(v, 10, 64)
}

func (dec *Decoder) parseTime(v string) (t time.Time, err error) {
	if v == "0" {
		return
	}
	var ts int64
	if ts, err = dec.parseInt(v); err != nil {
		return
	}
	t = ntpEpoch.Add(time.Second * time.Duration(ts))
	return
}

func (dec *Decoder) parseDuration(v string) (d time.Duration, err error) {
	mul := int64(1)
	if n := len(v) - 1; n >= 0 {
		switch v[n] {
		case 'd':
			mul, v = 86400, v[:n]
		case 'h':
			mul, v = 3600, v[:n]
		case 'm':
			mul, v = 60, v[:n]
		case 's':
			v = v[:n]
		}
	}
	var sec int64
	if sec, err = dec.parseInt(v); err != nil {
		return
	}
	d = time.Duration(sec*mul) * time.Second
	return
}

func (dec *Decoder) decodeLine(desc *Description, k byte, v string) (err error) {
	if n := len(desc.Media); n > 0 && k != 'm' {
		return dec.decodeMediaDesc(desc.Media[n-1], k, v)
	}
	return dec.decodeSessionDesc(desc, k, v)
}

func (dec *Decoder) decodeSessionDesc(desc *Description, k byte, v string) (err error) {
	switch k {
	case 'v':
		desc.Version, err = strconv.Atoi(v)
	case 'o':
		desc.Origin, err = dec.decodeOrigin(v)
	case 's':
		desc.Session = v
	case 'i':
		desc.Information = v
	case 'u':
		desc.URI = v
	case 'e':
		desc.Email = append(desc.Email, v)
	case 'p':
		desc.Phone = append(desc.Phone, v)
	case 'c':
		desc.Connection, err = dec.decodeConn(v)
	case 'b':
		if desc.Bandwidth == nil {
			desc.Bandwidth = make(map[string]int)
		}
		err = dec.decodeBandwidth(desc.Bandwidth, v)
	case 't':
		desc.Timing, err = dec.decodeTiming(v)
	case 'r':
		if desc.Timing != nil {
			desc.Timing.Repeat, err = dec.decodeRepeats(v)
		}
	case 'z':
		err = dec.decodeTimezones(desc, v)
	case 'k':
		desc.Key, err = dec.decodeKey(v)
	case 'a':
		var attr *Attribute
		if attr, err = dec.decodeAttr(v); err == nil {
			desc.Attributes = append(desc.Attributes, attr)
		}
	case 'm':
		var m *Media
		if m, err = dec.decodeMedia(v); err == nil {
			desc.Media = append(desc.Media, m)
		}
	default:
		err = decodeError(k)
	}
	return
}

func (dec *Decoder) decodeMediaDesc(m *Media, k byte, v string) (err error) {
	switch k {
	case 'i':
		m.Information = v
	case 'c':
		m.Connection, err = dec.decodeConn(v)
	case 'b':
		if m.Bandwidth == nil {
			m.Bandwidth = make(map[string]int)
		}
		err = dec.decodeBandwidth(m.Bandwidth, v)
	case 'k':
		m.Key, err = dec.decodeKey(v)
	case 'a':
		var attr *Attribute
		if attr, err = dec.decodeAttr(v); err == nil {
			m.Attributes = append(m.Attributes, attr)
		}
	default:
		err = decodeError(k)
	}
	return
}

func (dec *Decoder) decodeAttr(v string) (*Attribute, error) {
	if dec.split(v, ':', 2, false) {
		return &Attribute{Name: dec.p[0], Value: dec.p[1]}, nil
	}
	return &Attribute{Name: dec.p[0]}, nil
}

func (dec *Decoder) decodeKey(v string) (*Key, error) {
	if dec.split(v, ':', 2, false) {
		return &Key{Type: dec.p[0], Value: dec.p[1]}, nil
	}
	return &Key{Type: dec.p[0]}, nil
}

func (dec *Decoder) decodeOrigin(v string) (o *Origin, err error) {
	if !dec.split(v, ' ', 6, true) {
		return nil, dec.err
	}
	o = &Origin{
		Username: dec.p[0],
		Network:  dec.p[3],
		Type:     dec.p[4],
		Address:  dec.p[5],
	}
	if o.SessionID, err = dec.parseInt(dec.p[1]); err != nil {
		return nil, err
	}
	if o.SessionVersion, err = dec.parseInt(dec.p[2]); err != nil {
		return nil, err
	}
	return
}

func (dec *Decoder) decodeConn(v string) (c *Connection, err error) {
	if !dec.split(v, ' ', 3, true) {
		return nil, dec.err
	}
	c = &Connection{
		Network: dec.p[0],
		Type:    dec.p[1],
		Address: dec.p[2],
	}
	return
}

func (dec *Decoder) decodeBandwidth(b map[string]int, v string) (err error) {
	if !dec.split(v, ':', 2, true) {
		return dec.err
	}
	b[dec.p[0]], err = strconv.Atoi(dec.p[1])
	return
}

func (dec *Decoder) decodeTiming(v string) (t *Timing, err error) {
	if !dec.split(v, ' ', 2, true) {
		return nil, dec.err
	}
	t = &Timing{}
	t.Start, err = dec.parseTime(dec.p[0])
	if err == nil {
		t.Stop, err = dec.parseTime(dec.p[1])
	}
	return
}

func (dec *Decoder) decodeRepeats(v string) (r *Repeat, err error) {
	if !dec.split(v, ' ', 3, true) {
		return nil, dec.err
	}
	r = &Repeat{}
	r.Interval, err = dec.parseDuration(dec.p[0])
	if err == nil {
		r.Duration, err = dec.parseDuration(dec.p[1])
	}
	if err == nil {
		var d time.Duration
		dec.split(dec.p[2], ' ', 255, false)
		for _, it := range dec.p {
			if d, err = dec.parseDuration(it); err != nil {
				return
			}
			r.Offsets = append(r.Offsets, d)
		}
	}
	return
}

func (dec *Decoder) decodeTimezones(desc *Description, v string) (err error) {
	dec.split(v, ' ', 255, false)
	for i, n := 0, len(dec.p)-1; i < n; i += 2 {
		z := &TimeZone{}
		if z.Time, err = dec.parseTime(dec.p[i]); err != nil {
			return
		}
		if z.Offset, err = dec.parseDuration(dec.p[i+1]); err != nil {
			return
		}
		desc.TimeZones = append(desc.TimeZones, z)
	}
	return
}

func (dec *Decoder) decodeMedia(v string) (m *Media, err error) {
	if !dec.split(v, ' ', 4, true) {
		return nil, dec.err
	}
	m = &Media{
		Type:  dec.p[0],
		Proto: dec.p[2],
	}
	var formats = dec.p[3]
	if dec.split(dec.p[1], '/', 2, false) {
		m.PortNum, err = strconv.Atoi(dec.p[1])
	}
	if err == nil {
		m.Port, err = strconv.Atoi(dec.p[0])
	}
	if err == nil {
		dec.split(formats, ' ', 255, false)
		for _, it := range dec.p {
			p, err := strconv.Atoi(it)
			if err != nil {
				return nil, err
			}
			dec.touchMediaFormat(m, p)
		}
	}
	return
}

type bufferedReader struct {
	buf *bufio.Reader
}

func (r *bufferedReader) ReadLine() (string, error) {
	ln, p, err := r.buf.ReadLine()
	if p {
		err = decodeError("line is too large")
	}
	if err != nil {
		return "", err
	}
	return string(ln), nil
}

type stringReader struct {
	buf string
}

func (r *stringReader) ReadLine() (string, error) {
	n := len(r.buf)
	if n == 0 {
		return "", io.EOF
	}
	i, j := 0, 0
	for j < n {
		if c := r.buf[j]; c == '\n' {
			break
		} else if c == '\r' {
			j++
		} else {
			j++
			i = j
		}
	}
	v := r.buf[:i]
	if n > j {
		r.buf = r.buf[j+1:]
	} else {
		r.buf = ""
	}
	return v, nil
}

type decodeError string

func (d decodeError) Error() string {
	return "sdp: decode error '" + string(d) + "'"
}

var ntpEpoch = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
