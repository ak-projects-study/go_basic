package arrint

import (
	"fmt"
	"strings"
)

// ArrInt описывает слайс целых чисел типа int.
type ArrInt []int

// Add складывает поэлементно два массива и возвращает результирующий слайс.
func Add(a, b ArrInt) ArrInt {
	length := len(a)
	if length-len(b) > 0 {
		length = len(b)
	}
	c := make(ArrInt, length)
	for i := 0; i < length; i++ {
		c[i] = a[i] + b[i]
	}
	return c
}

// Метод String преобразует ArrInt в строку и возвращает её.
func (a ArrInt) String() string {
	out := make([]string,
		len(a))

	for i, v := range a {
		out[i] = fmt.Sprintf(`<%d>`, v)
	}
	return fmt.Sprintf(`[%s]`, strings.Join(out, ` `))
}
