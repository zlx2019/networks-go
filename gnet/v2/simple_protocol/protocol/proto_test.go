package protocol

import (
	"fmt"
	"testing"
)

func TestCodec(t *testing.T) {
	codec := SimpleCodec{}
	bytes, err := codec.Encode([]byte("Hello world!"))
	if err != nil {
		t.Fatal(err)
	}
	out, err := codec.Unpack(bytes)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(out))
}
