package util

import (
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
	if !IsNilOrEmpty(&model) {
		t.FailNow()
	}
	//
	// model.C = 1
	model.D = &n
	if IsNilOrEmpty(&model) {
		t.FailNow()
	}
}

type CopyValue1 struct {
	A int
	B int
	C *int
	D *int
	E int
	F CopyValue3
	G *CopyValue3
	H *CopyValue3
	CopyValue3
}

type CopyValue2 struct {
	A *int
	B int
	C int
	D *int
	E string
	F CopyValue4
	G CopyValue4
	H CopyValue3
	CopyValue3
}

type CopyValue3 struct {
	A string
	B int
}

type CopyValue4 struct {
	A string
	B *int
}

func Test_CopyStruct(t *testing.T) {
	var dst CopyValue1
	var src CopyValue2
	a2 := 12
	src.A = &a2
	src.B = 22
	src.C = 32
	d2 := 42
	src.D = &d2
	src.E = "e2"
	//
	src.F.A = "fa12"
	b2 := 222
	src.F.B = &b2
	src.G.A = "ga"
	src.G.B = &b2
	src.CopyValue3.A = "fa22"
	src.CopyValue3.B = 2222
	dst.G = new(CopyValue3)
	//
	CopyStruct(&dst, &src)
	//
	if dst.A != 0 ||
		dst.B != src.B ||
		dst.C != nil ||
		dst.D == nil || *dst.D != *src.D ||
		dst.E != 0 ||
		dst.F.A != src.F.A || dst.F.B != 0 ||
		dst.G.A != src.G.A || dst.G.B != 0 ||
		dst.H != nil ||
		dst.CopyValue3.A != src.CopyValue3.A || dst.CopyValue3.B != src.CopyValue3.B {
		t.FailNow()
	}

}
