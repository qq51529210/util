package util

import (
	"fmt"
	"testing"
)

func Test_BitWriter(t *testing.T) {
	w := new(BitWriter)
	w.Write8(123, 2)
	for _, b := range w.Raw() {
		fmt.Printf("%08b", b)
	}
	w.Write8(1, 3)
	for _, b := range w.Raw() {
		fmt.Printf("%08b", b)
	}
	w.Write8(211, 7)
	for _, b := range w.Raw() {
		fmt.Printf("%08b", b)
	}
	w.Write8(23, 5)
	for _, b := range w.Bytes() {
		fmt.Printf("%08b", b)
	}
	fmt.Println(w.Bytes())
}
