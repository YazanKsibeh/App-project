package gen_key

import (
	"encoding/base64"
	"strconv"
	"testing"
	"time"
)

const key = "M54pwYCXjoQyxEmq0MeW2DdVph0AhjdwXYvZEWZWkBux2j0860nOPP3ysCcaD0sMCiC6336sOT8dAKuIPD6OLxFNQD1TZuCHSw2iKAluOc1zqeN8gSfxXJ0w6ao21YZvLsc9FuZvJ6NGGDiAyMCHLM3bTglm6vLIAHGNI7ncmmFR8j6oLOtjQf5PUEhhgeCpeh43KwqAmZdBlR2KbrigMHg7yzOvzbZCCrQUFQTlSrIONDYTHBcfv8gXnJWUgHHWOcOLtXCSdv7DkbXTcnUz0KoNKQS9HCg9jrH5PyZqD3duYiZPRD09uQxspmXy9o2RwTcTVLEK7ulpP5ctZEk0DEGNB3kOrFzJlKMqqzues7Zr8xt4ZWkGgXkPA46FpPvbXLjQKweZIepDaQI7tRKJkdt9DqZ1WjNUJ5da56Xk2xZpmOEIM9iZZ2H79CSmtwQZdnc8sXungkyQ7oA4KycFNtDZjhf5hlvmYy8Uhg1eF1u08DSHRhZw0PyFfEzYeTBJ"

func TestCheckKey(t *testing.T) {
	t.Log("Start generation key\n\n")

	array := arrayInit()

	// get offset
	var offset int64
	var str string

	// if number from 10 to 99
	offset, err := strconv.ParseInt(key[1:3], 10, 64)
	if err != nil {
		t.Fatalf("unable to convert string to int: %s", err.Error())
	}

	t.Logf("offset: %d", offset)

	for _, v := range array {
		if len(key) >= v {
			str += string(key[int(offset)+v+1])
		}
	}

	t.Logf("date format base64: %s", str)

	// check len date
	if len(str) > 11 {
		for {
			str = str[:len(str)-1]
			if len(str) == 11 {
				break
			}
		}
	}

	str = Reverse(str)
	str = str + "="

	t.Logf("date format base64 reverse: %s", str)

	b64, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		t.Fatalf("unable to umarshal string: %s", err.Error())
	}

	date := string(b64)

	t.Logf("date: %s", date)

	tm, err := time.Parse(layout, date)
	if err != nil {
		t.Fatalf("unable convert to time: %s", err.Error())
	}

	t.Logf("time from key: %+v", tm)
	t.Logf("time now: %+v", time.Now())

	if time.Now().After(tm) {
		t.Fatalf("license expired")
	}

	t.Logf("good job")

	t.Log("\n\nFinish generation key")
}
