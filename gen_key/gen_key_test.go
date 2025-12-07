package gen_key

import (
	"encoding/base64"
	"strconv"
	"strings"
	"testing"
)

func TestGenKey(t *testing.T) {
	t.Log("Start generation key\n\n")

	fib := arrayInit()

	date := "20260101" // change date - 30 days from now
	offset := 54       // change offset key

	// convert string to base64
	b64 := base64.StdEncoding.EncodeToString([]byte(date))
	t.Logf("date after convert base64: %s", b64)

	b64 = strings.ReplaceAll(b64, "=", "")

	// reverse
	rb64 := Reverse(b64)
	t.Logf("date after reverse: %s", rb64)

	// generation sting
	f := GenString()
	t.Logf("rand string: %s", f)

	// file convert to base64
	strOffset := strconv.Itoa(offset) //

	file := f[:1] + strOffset + f[1:]

	var key string

	for i, v := range rb64 {
		key += file[len(key):len(strOffset)+offset+fib[i]-1] + string(v)
	}

	key += file[len(key):]

	t.Logf("len file: %d", len(file))
	t.Logf("len key: %d", len(key))

	t.Logf("file: %s", file)
	t.Logf("key:  %s", key)

	t.Log("\n\nFinish generation key")
}
