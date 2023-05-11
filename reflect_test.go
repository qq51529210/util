package util

import (
	"reflect"
	"testing"
)

type IsNilOrEmpty1 struct {
	A *int
	B int
}

type IsNilOrEmpty2 struct {
	IsNilOrEmpty1
	C int
	D *int
}

func Test_IsNilOrEmpty(t *testing.T) {
	n := 0
	var model IsNilOrEmpty2
	v := reflect.ValueOf(&model)
	if !IsNilOrEmpty(v) {
		t.FailNow()
	}
	//
	// model.C = 1
	model.D = &n
	if IsNilOrEmpty(v) {
		t.FailNow()
	}
}
