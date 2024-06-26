package colfer

// Code generated by colf(1); DO NOT EDIT.
// The compiler used schema file structdef.colf.

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"time"
)

var intconv = binary.BigEndian

// Colfer configuration attributes
var (
	// ColferSizeMax is the upper limit for serial byte sizes.
	ColferSizeMax = 16 * 1024 * 1024
)

// ColferMax signals an upper limit breach.
type ColferMax string

// Error honors the error interface.
func (m ColferMax) Error() string { return string(m) }

// ColferError signals a data mismatch as as a byte index.
type ColferError int

// Error honors the error interface.
func (i ColferError) Error() string {
	return fmt.Sprintf("colfer: unknown header at byte %d", i)
}

// ColferTail signals data continuation as a byte index.
type ColferTail int

// Error honors the error interface.
func (i ColferTail) Error() string {
	return fmt.Sprintf("colfer: data continuation at byte %d", i)
}

type ColferA struct {
	Name string

	BirthDay time.Time

	Phone string

	Siblings int32

	Spouse bool

	Money float64
}

// MarshalTo encodes o as Colfer into buf and returns the number of bytes written.
// If the buffer is too small, MarshalTo will panic.
func (o *ColferA) MarshalTo(buf []byte) int {
	var i int

	if l := len(o.Name); l != 0 {
		buf[i] = 0
		i++
		x := uint(l)
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
		i += copy(buf[i:], o.Name)
	}

	if v := o.BirthDay; !v.IsZero() {
		s, ns := uint64(v.Unix()), uint32(v.Nanosecond())
		if s < 1<<32 {
			buf[i] = 1
			intconv.PutUint32(buf[i+1:], uint32(s))
			i += 5
		} else {
			buf[i] = 1 | 0x80
			intconv.PutUint64(buf[i+1:], s)
			i += 9
		}
		intconv.PutUint32(buf[i:], ns)
		i += 4
	}

	if l := len(o.Phone); l != 0 {
		buf[i] = 2
		i++
		x := uint(l)
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
		i += copy(buf[i:], o.Phone)
	}

	if v := o.Siblings; v != 0 {
		x := uint32(v)
		if v >= 0 {
			buf[i] = 3
		} else {
			x = ^x + 1
			buf[i] = 3 | 0x80
		}
		i++
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
	}

	if o.Spouse {
		buf[i] = 4
		i++
	}

	if v := o.Money; v != 0 {
		buf[i] = 5
		intconv.PutUint64(buf[i+1:], math.Float64bits(v))
		i += 9
	}

	buf[i] = 0x7f
	i++
	return i
}

// MarshalLen returns the Colfer serial byte size.
// The error return option is goserbench.ColferMax.
func (o *ColferA) MarshalLen() (int, error) {
	l := 1

	if x := len(o.Name); x != 0 {
		if x > ColferSizeMax {
			return 0, ColferMax(fmt.Sprintf("colfer: field goserbench.ColferA.Name exceeds %d bytes", ColferSizeMax))
		}
		for l += x + 2; x >= 0x80; l++ {
			x >>= 7
		}
	}

	if v := o.BirthDay; !v.IsZero() {
		if s := uint64(v.Unix()); s < 1<<32 {
			l += 9
		} else {
			l += 13
		}
	}

	if x := len(o.Phone); x != 0 {
		if x > ColferSizeMax {
			return 0, ColferMax(fmt.Sprintf("colfer: field goserbench.ColferA.Phone exceeds %d bytes", ColferSizeMax))
		}
		for l += x + 2; x >= 0x80; l++ {
			x >>= 7
		}
	}

	if v := o.Siblings; v != 0 {
		x := uint32(v)
		if v < 0 {
			x = ^x + 1
		}
		for l += 2; x >= 0x80; l++ {
			x >>= 7
		}
	}

	if o.Spouse {
		l++
	}

	if o.Money != 0 {
		l += 9
	}

	if l > ColferSizeMax {
		return l, ColferMax(fmt.Sprintf("colfer: struct goserbench.ColferA exceeds %d bytes", ColferSizeMax))
	}
	return l, nil
}

// MarshalBinary encodes o as Colfer conform encoding.BinaryMarshaler.
// The error return option is goserbench.ColferMax.
func (o *ColferA) MarshalBinary() (data []byte, err error) {
	l, err := o.MarshalLen()
	if err != nil {
		return nil, err
	}
	data = make([]byte, l)
	o.MarshalTo(data)
	return data, nil
}

// Unmarshal decodes data as Colfer and returns the number of bytes read.
// The error return options are io.EOF, goserbench.ColferError and goserbench.ColferMax.
func (o *ColferA) Unmarshal(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, io.EOF
	}
	header := data[0]
	i := 1

	if header == 0 {
		if i >= len(data) {
			goto eof
		}
		x := uint(data[i])
		i++

		if x >= 0x80 {
			x &= 0x7f
			for shift := uint(7); ; shift += 7 {
				if i >= len(data) {
					goto eof
				}
				b := uint(data[i])
				i++

				if b < 0x80 {
					x |= b << shift
					break
				}
				x |= (b & 0x7f) << shift
			}
		}

		if x > uint(ColferSizeMax) {
			return 0, ColferMax(fmt.Sprintf("colfer: goserbench.ColferA.Name size %d exceeds %d bytes", x, ColferSizeMax))
		}

		start := i
		i += int(x)
		if i >= len(data) {
			goto eof
		}
		o.Name = string(data[start:i])

		header = data[i]
		i++
	}

	if header == 1 {
		start := i
		i += 8
		if i >= len(data) {
			goto eof
		}
		o.BirthDay = time.Unix(int64(intconv.Uint32(data[start:])), int64(intconv.Uint32(data[start+4:]))).In(time.UTC)
		header = data[i]
		i++
	} else if header == 1|0x80 {
		start := i
		i += 12
		if i >= len(data) {
			goto eof
		}
		o.BirthDay = time.Unix(int64(intconv.Uint64(data[start:])), int64(intconv.Uint32(data[start+8:]))).In(time.UTC)
		header = data[i]
		i++
	}

	if header == 2 {
		if i >= len(data) {
			goto eof
		}
		x := uint(data[i])
		i++

		if x >= 0x80 {
			x &= 0x7f
			for shift := uint(7); ; shift += 7 {
				if i >= len(data) {
					goto eof
				}
				b := uint(data[i])
				i++

				if b < 0x80 {
					x |= b << shift
					break
				}
				x |= (b & 0x7f) << shift
			}
		}

		if x > uint(ColferSizeMax) {
			return 0, ColferMax(fmt.Sprintf("colfer: goserbench.ColferA.Phone size %d exceeds %d bytes", x, ColferSizeMax))
		}

		start := i
		i += int(x)
		if i >= len(data) {
			goto eof
		}
		o.Phone = string(data[start:i])

		header = data[i]
		i++
	}

	if header == 3 {
		if i+1 >= len(data) {
			i++
			goto eof
		}
		x := uint32(data[i])
		i++

		if x >= 0x80 {
			x &= 0x7f
			for shift := uint(7); ; shift += 7 {
				b := uint32(data[i])
				i++
				if i >= len(data) {
					goto eof
				}

				if b < 0x80 {
					x |= b << shift
					break
				}
				x |= (b & 0x7f) << shift
			}
		}
		o.Siblings = int32(x)

		header = data[i]
		i++
	} else if header == 3|0x80 {
		if i+1 >= len(data) {
			i++
			goto eof
		}
		x := uint32(data[i])
		i++

		if x >= 0x80 {
			x &= 0x7f
			for shift := uint(7); ; shift += 7 {
				b := uint32(data[i])
				i++
				if i >= len(data) {
					goto eof
				}

				if b < 0x80 {
					x |= b << shift
					break
				}
				x |= (b & 0x7f) << shift
			}
		}
		o.Siblings = int32(^x + 1)

		header = data[i]
		i++
	}

	if header == 4 {
		if i >= len(data) {
			goto eof
		}
		o.Spouse = true
		header = data[i]
		i++
	}

	if header == 5 {
		start := i
		i += 8
		if i >= len(data) {
			goto eof
		}
		o.Money = math.Float64frombits(intconv.Uint64(data[start:]))
		header = data[i]
		i++
	}

	if header != 0x7f {
		return 0, ColferError(i - 1)
	}
	if i < ColferSizeMax {
		return i, nil
	}
eof:
	if i >= ColferSizeMax {
		return 0, ColferMax(fmt.Sprintf("colfer: struct goserbench.ColferA size exceeds %d bytes", ColferSizeMax))
	}
	return 0, io.EOF
}

// UnmarshalBinary decodes data as Colfer conform encoding.BinaryUnmarshaler.
// The error return options are io.EOF, goserbench.ColferError, goserbench.ColferTail and goserbench.ColferMax.
func (o *ColferA) UnmarshalBinary(data []byte) error {
	i, err := o.Unmarshal(data)
	if i < len(data) && err == nil {
		return ColferTail(i)
	}
	return err
}
