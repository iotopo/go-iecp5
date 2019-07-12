package asdu

import (
	"encoding/binary"
	"time"
)

func CP56Time2a(t *time.Time, loc *time.Location) []byte {
	if loc == nil {
		loc = time.UTC
	}
	ts := t.In(loc)
	return []byte{byte(ts.Nanosecond() / int(time.Millisecond)), byte(ts.Second()), byte(ts.Minute()),
		byte(ts.Hour()), byte(ts.Month()), byte(ts.Year() - 2000)}
}

// 7个八位位组二进制时间，建议所有时标采用UTC
// 读7字节,返回一个值，当无效时返回nil
// The year is assumed to be in the 20th century.
// See IEC 60870-5-4 § 6.8 and IEC 60870-5-101 second edition § 7.2.6.18.
func ParseCP56Time2a(bytes []byte, loc *time.Location) *time.Time {
	if loc == nil {
		loc = time.UTC
	}

	x := int(bytes[0])
	x |= int(bytes[1]) << 8
	msec := x % 1000
	sec := (x / 1000)

	o := bytes[2]
	min := int(o & 63)
	if o > 127 {
		return nil
	}

	hour := int(bytes[3] & 31)
	day := int(bytes[4] & 31)
	month := time.Month(bytes[5] & 15)
	year := 2000 + int(bytes[6]&127)

	nsec := msec * int(time.Millisecond)
	val := time.Date(year, month, day, hour, min, sec, nsec, loc)
	return &val
}

func CP24Time2a(t *time.Time, loc *time.Location) []byte {
	if loc == nil {
		loc = time.UTC
	}
	ts := t.In(loc)
	return []byte{byte(ts.Nanosecond() / int(time.Millisecond)), byte(ts.Second()), byte(ts.Minute())}
}

// 3个八位位组二进制时间，建议所有时标采用UTC
// 读3字节,返回一个值，当无效时返回nil
// The moment is assumed to be in the recent present.
// See IEC 60870-5-4 § 6.8 and IEC 60870-5-101 second edition § 7.2.6.19.
func ParseCP24Time2a(bytes []byte, loc *time.Location) *time.Time {
	if loc == nil {
		loc = time.UTC
	}

	x := int(bytes[0])
	x |= int(bytes[1]) << 8
	msec := x % 1000
	sec := (x / 1000)

	o := bytes[2]
	min := int(o & 63)
	if o > 127 {
		return nil
	}

	now := time.Now()
	year, month, day := now.Date()
	hour, currentMin, _ := now.Clock()

	nsec := msec * int(time.Millisecond)
	val := time.Date(year, month, day, hour, min, sec, nsec, loc)

	// 5 minute rounding - 55 minute span
	if min > currentMin+5 {
		val = val.Add(-time.Hour)
	}

	return &val
}

func CP16Time2a(msec uint16) []byte {
	return []byte{byte(msec), byte(msec >> 8)}
}

func ParseCP16Time2a(b []byte) uint16 {
	return binary.LittleEndian.Uint16(b)
}
