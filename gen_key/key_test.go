package gen_key

import (
	"encoding/base64"
	"strconv"
	"testing"
	"time"
)

func BenchmarkCheckKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Check()
	}
}

func Check() {
	fib := arrayInit()

	// get offset
	var offset int64
	var str string

	// if number from 10 to 99
	offset, err := strconv.ParseInt(key[1:3], 10, 64)
	if err != nil {
		// if number from 1 to 9
		offset, err = strconv.ParseInt(key[1:2], 10, 64)
		if err != nil {
			return
		}
	}

	for _, v := range fib {
		str += string(key[int(offset)+v+2])
	}

	// TODO Может ли быть base64 (20221231) больше чем 11 в base64?
	// check len date
	if len(str) > 11 {
		str = str[:len(str)-1]
	}

	str = Reverse(str)

	b64, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return
	}

	date := string(b64)

	tm, err := time.Parse(layout, date)
	if err != nil {
		return
	}

	if time.Now().After(tm) {
		return
	}
}
