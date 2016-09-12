package sdp

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"time"
)

var ErrNotImplemented = errors.New("sdp: not implemented")

type reader interface {
	ReadString(delim byte) (string, error)
}

type Decoder struct {
	r   reader
	p   []string
	err error
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: bufio.NewReader(r)}
}

func NewDecoderString(v string) *Decoder {
	return &Decoder{r: &stringReader{buf: v}}
}

func (dec *Decoder) Decode() (desc *Description, err error) {
	var v string
	desc = &Description{}
	for err == nil {
		if v, err = dec.r.ReadString('\n'); err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}
		if n := len(v); n > 0 && v[n-1] == '\r' {
			v = v[:n-1]
		}
		if dec.split(v, '=', 2, true) {
			err = dec.decode(desc, dec.p[0], dec.p[1])
		} else {
			err = dec.err
		}
	}
	return
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

func (dec *Decoder) parseTime(v string) (t *time.Time, err error) {
	if v == "0" {
		return
	}
	var ts int64
	if ts, err = dec.parseInt(v); err != nil {
		return
	}
	r := ntpEpoch.Add(time.Second * time.Duration(ts))
	return &r, nil
}

func (dec *Decoder) decode(desc *Description, k, v string) (err error) {
	if len(k) != 1 {
		return decodeError(k)
	}
	switch k[0] {
	case 'v':
		desc.Version, err = strconv.Atoi(v)
	case 'o':
		if !dec.split(v, ' ', 6, true) {
			break
		}
		orig := &Origin{
			Username: dec.p[0],
			Network:  dec.p[3],
			Type:     dec.p[4],
			Address:  dec.p[5],
		}
		orig.SessionID, err = dec.parseInt(dec.p[1])
		if err == nil {
			orig.SessionVersion, err = dec.parseInt(dec.p[2])
		}
		if err == nil {
			desc.Origin = orig
		}
	case 's':
		desc.Session = v
	case 'i':
		desc.Information = append(desc.Information, v)
	case 'u':
		desc.URI = v
	case 'e':
		desc.Email = append(desc.Email, v)
	case 'p':
		desc.Phone = append(desc.Phone, v)
	case 'c':
		if !dec.split(v, ' ', 3, true) {
			break
		}
		desc.Connection = &Connection{
			Network: dec.p[0],
			Type:    dec.p[1],
			Address: dec.p[2],
		}
	case 'b':
		if !dec.split(v, ':', 2, true) {
			break
		}
		if desc.Bandwidth == nil {
			desc.Bandwidth = make(Bandwidth)
		}
		desc.Bandwidth[dec.p[0]], err = dec.parseInt(dec.p[1])
	case 't':
		if !dec.split(v, ' ', 2, true) {
			break
		}
		t := &Timing{}
		t.Start, err = dec.parseTime(dec.p[0])
		if err == nil {
			t.Stop, err = dec.parseTime(dec.p[1])
		}
		if err == nil {
			desc.Timing = append(desc.Timing, t)
		}
	case 'r':
		return ErrNotImplemented
	case 'z':
		return ErrNotImplemented
	case 'a':
		var attr *Attribute
		if dec.split(v, ':', 2, false) {
			attr = &stringAttr{n: dec.p[0], v: dec.p[1]}
		} else {
			attr = &stringAttr{n: dec.p[0]}
		}
		desc.Attributes = append(desc.Attributes, attr)
	case 'm':
		if !dec.split(v, ' ', 4, true) {
			break
		}
		m := &Media{
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
		if err != nil {
			break
		}
		dec.split(formats, ' ', 255, false)
		for _, it := range dec.p {
			f, err := strconv.Atoi(it)
			if err != nil {
				return err
			}
			m.Formats = append(m.Formats, f)
		}
		desc.Media = append(desc.Media, m)
	}
	if err == nil && dec.err != nil {
		err = dec.err
	}
	return
}

type stringReader struct {
	buf string
}

func (r *stringReader) ReadString(delim byte) (string, error) {
	i, n := 0, len(r.buf)
	if n == 0 {
		return "", io.EOF
	}
	for i < n && r.buf[i] != delim {
		i++
	}
	v := r.buf[:i]
	if len(r.buf) > i {
		r.buf = r.buf[i+1:]
	} else {
		r.buf = ""
	}
	return v, nil
}

type decodeError string

func (d decodeError) Error() string {
	return "sdp: decode error '" + string(d) + "'"
}

type stringAttr struct {
	n, v string
}

func (attr *stringAttr) Name() string {
	return attr.n
}

func (attr *stringAttr) Value() string {
	return attr.v
}

func (attr *stringAttr) String() string {
	return attr.n + ":" + attr.v
}

var ntpEpoch = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
