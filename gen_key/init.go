package gen_key

import (
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"
)

const layout = "20060102"

func Reverse(s string) string {
	size := len(s)
	buf := make([]byte, size)
	for start := 0; start < size; {
		r, n := utf8.DecodeRuneInString(s[start:])
		start += n
		utf8.EncodeRune(buf[size-start:], r)
	}
	return string(buf)
}

func GenString() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 510
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String() // Например "ExcbsVQs"

	return str
}

func arrayInit() []int {
	array := make([]int, 0)
	array = append(array, 3)
	array = append(array, 5)
	array = append(array, 8)
	array = append(array, 13)
	array = append(array, 21)
	array = append(array, 34)
	array = append(array, 55)
	array = append(array, 89)
	array = append(array, 144)
	array = append(array, 233)
	array = append(array, 377)
	array = append(array, 610)
	array = append(array, 987)

	return array
}
