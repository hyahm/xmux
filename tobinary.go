package xmux

import (
	"encoding/binary"
	"regexp"
)

// Uint16ToBinaryString get the string of a uint16 number in binary format.
func Uint16ToBinaryString(i uint16) string {
	bs := make([]byte, 2)
	binary.BigEndian.PutUint16(bs, i)
	return BytesToBinaryString(bs)
}

// Uint32ToBinaryString get the string of a uint32 number in binary format.
func Uint32ToBinaryString(i uint32) string {
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, i)
	return BytesToBinaryString(bs)
}

// Uint64ToBinaryString get the string of a uint64 number in binary format.
func Uint64ToBinaryString(i uint64) string {
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, i)
	return BytesToBinaryString(bs)
}

// ToBinaryString get string in binary format according to input data.
// The input data can be diffrent kinds of basic data type.
func ToBinaryString(v interface{}) (s string) {

	switch v := v.(type) {
	case []byte:
		s = BytesToBinaryString(v)
	case int8:
		s = ByteToBinaryString(uint8(v))
	case uint8: // byte
		s = ByteToBinaryString(v)
	case int16:
		s = Uint16ToBinaryString(uint16(v))
	case uint16:
		s = Uint16ToBinaryString(v)
	case int32:
		s = Uint32ToBinaryString(uint32(v))
	case uint32:
		s = Uint32ToBinaryString(v)
	case uint:
		s = Uint64ToBinaryString(uint64(v))
	case int:
		s = Uint64ToBinaryString(uint64(v))
	case int64:
		s = Uint64ToBinaryString(uint64(v))
	case uint64:
		s = Uint64ToBinaryString(v)
		//TODO add float number support
	default:
		panic(ErrTypeUnsupport)

	}
	return

}

// ReadBinaryString read the string in binary format into input data.
func ReadBinaryString(s string, data interface{}) (err error) {

	bs := BinaryStringToBytes(s)
	switch data := data.(type) {

	case *int8:
		*data = int8(bs[0])
	case *uint8:
		*data = bs[0]
	case *int16:
		*data = int16(bytesToUint16(bs))
	case *uint16:
		*data = bytesToUint16(bs)
	case *int32:
		*data = int32(bytesToUint32(bs))
	case *uint32:
		*data = bytesToUint32(bs)
	case *int64:
		*data = int64(bytesToUint64(bs))
	case *uint64:
		*data = bytesToUint64(bs)
		//TODO add float number support
	default:
		err = ErrTypeUnsupport
	}
	return
}

func bytesToUint16(bs []byte) uint16 {
	bs = fillBytes(bs, 2)
	return binary.BigEndian.Uint16(bs)
}

func bytesToUint32(bs []byte) uint32 {
	bs = fillBytes(bs, 4)
	return binary.BigEndian.Uint32(bs)
}

func bytesToUint64(bs []byte) uint64 {
	bs = fillBytes(bs, 8)
	return binary.BigEndian.Uint64(bs)
}

// fillBytes fills byte slice with zero bytes ahead when its'
// length is not greater than n.
func fillBytes(bs []byte, n int) []byte {

	l := len(bs)
	if l >= n {
		return bs
	}

	nbs := make([]byte, n)
	n -= l // n zero bytes need to fill

	for i := 0; i < n; i++ {
		nbs[i] = byte(0)
	}

	copy(nbs[n:], bs)

	return nbs
}

// ToHexString get string in Hexadecimal format according to input data.
// The input data can be diffrent kinds of basic data type.
func ToHexString(v interface{}) (s string) {
	//TODO implements it

	return
}

// ToOctalString get string in octal format according to input data.
// The input data can be diffrent kinds of basic data type.
func ToOctalString(v interface{}) (s string) {
	//TODO implements it

	return
}

// ByteToBinaryString get the string in binary format of a byte or uint8.
func ByteToBinaryString(b byte) string {
	buf := make([]byte, 0, 8)
	buf = appendBinaryString(buf, b)
	return string(buf)
}

// BytesToBinaryString get the string in binary format of a []byte or []int8.
func BytesToBinaryString(bs []byte) string {
	l := len(bs)
	bl := l*8 + l + 1
	buf := make([]byte, 0, bl)
	buf = append(buf, lsb)
	for _, b := range bs {
		buf = appendBinaryString(buf, b)
		buf = append(buf, space)
	}
	buf[bl-1] = rsb
	return string(buf)
}

// regex for delete useless string which is going to be in binary format.
var rbDel = regexp.MustCompile(`[^01]`)

// BinaryStringToBytes get the binary bytes according to the
// input string which is in binary format.
func BinaryStringToBytes(s string) (bs []byte) {
	if len(s) == 0 {
		panic(ErrEmptyString)
	}

	s = rbDel.ReplaceAllString(s, "")
	l := len(s)
	if l == 0 {
		panic(ErrBadStringFormat)
	}

	mo := l % 8
	l /= 8
	if mo != 0 {
		l++
	}
	bs = make([]byte, 0, l)
	mo = 8 - mo
	var n uint8
	for i, b := range []byte(s) {
		m := (i + mo) % 8
		switch b {
		case one:
			n += uint8arr[m]
		}
		if m == 7 {
			bs = append(bs, n)
			n = 0
		}
	}
	return
}

// append bytes of string in binary format.
func appendBinaryString(bs []byte, b byte) []byte {
	var a byte
	for i := 0; i < 8; i++ {
		a = b
		b <<= 1
		b >>= 1
		switch a {
		case b:
			bs = append(bs, zero)
		default:
			bs = append(bs, one)
		}
		b <<= 1
	}
	return bs
}
