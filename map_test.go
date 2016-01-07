package protobuf

import (
	"bytes"
	"reflect"
	"testing"
)

type Inner struct {
	Id   int32
	Name string
}

type MessageWithMap struct {
	NameMapping map[int32]string // = 1, required
	//ByteMapping   map[bool][]byte
	//StructMapping map[string]Inner
}

func TestMapFieldEncode(t *testing.T) {
	m := &MessageWithMap{
		NameMapping: map[int32]string{
			1: "Rob",
			4: "Ian",
			8: "Dave",
		},
	}

	b, err := Encode(m)
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}

	// b should be the concatenation of these three byte sequences in some order.
	parts := []string{
		"\n\a\b\x01\x12\x03Rob",
		"\n\a\b\x04\x12\x03Ian",
		"\n\b\b\x08\x12\x04Dave",
	}
	ok := false
	for i := range parts {
		for j := range parts {
			if j == i {
				continue
			}
			for k := range parts {
				if k == i || k == j {
					continue
				}
				try := parts[i] + parts[j] + parts[k]
				if bytes.Equal(b, []byte(try)) {
					ok = true
					break
				}
			}
		}
	}
	if !ok {
		t.Fatalf("Incorrect Encoding output.\n got %q\nwant %q (or a permutation of that)", b, parts[0]+parts[1]+parts[2])
	}
	t.Logf("FYI b: %q", b)

	//(new(Buffer)).DebugPrint("Dump of b", b)
}

func TestMapFieldRoundTrips(t *testing.T) {
	m := &MessageWithMap{
		NameMapping: map[int32]string{
			1: "Rob",
			4: "Ian",
			8: "Dave",
		},
		// MsgMapping: map[int64]*FloatingPoint{
		// 	0x7001: &FloatingPoint{F: Float64(2.0)},
		// },
		// ByteMapping: map[bool][]byte{
		// 	false: []byte("that's not right!"),
		// 	true:  []byte("aye, 'tis true!"),
		// },
	}
	b, err := Encode(m)
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	t.Logf("FYI b: %q", b)
	m2 := new(MessageWithMap)
	// First fix the encoding:
	if err := Decode(b, m2); err != nil {
		t.Fatalf("Decode: %v", err)
	}
	for _, pair := range [][2]interface{}{
		{m.NameMapping, m2.NameMapping},
	} {
		if !reflect.DeepEqual(pair[0], pair[1]) {
			t.Errorf("Map did not survive a round trip.\ninitial: %v\n  final: %v", pair[0], pair[1])
		}
	}
}

func TestMapFieldWithNil(t *testing.T) {
	// m := &MessageWithMap{
	// 	MsgMapping: map[int64]*FloatingPoint{
	// 		1: nil,
	// 	},
	// }
	// b, err := Marshal(m)
	// if err == nil {
	// 	t.Fatalf("Marshal of bad map should have failed, got these bytes: %v", b)
	// }
}
