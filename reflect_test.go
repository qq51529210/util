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

type CopyStructDst struct {
	F CopyStruct1
	B string
	G *CopyStruct1
	C *string
	H *CopyStruct1
	A string
	D *string
	CopyStruct1
	E int
}

func newCopyStructDst() *CopyStructDst {
	dst := new(CopyStructDst)
	dst.A = "dst.a"
	dst.B = "dst.a"
	c := "dst.c"
	dst.C = &c
	dst.E = 1
	dst.G = new(CopyStruct1)
	//
	return dst
}

type CopyStructSrc struct {
	A *string
	B string
	C string
	D *string
	E string
	F CopyStruct2
	G CopyStruct2
	H CopyStruct1
	CopyStruct1
	I CopyStruct1
}

func newCopyStructSrc() *CopyStructSrc {
	src := new(CopyStructSrc)
	a := "src.a"
	src.A = &a
	src.B = "src.b"
	src.C = "src.c"
	d := "src.d"
	src.D = &d
	src.E = "src.e"
	src.F.A = "src.f.a"
	fb := "src.f.b"
	src.F.B = &fb
	gb := "src.g.b"
	src.G.A = "src.g.a"
	src.G.B = &gb
	src.H.A = "src.h.a"
	src.H.B = "src.h.b"
	src.CopyStruct1.A = "src.cs1.a"
	src.CopyStruct1.B = "src.cs1.b"
	src.I.A = "src.i.a"
	src.I.B = "src.i.b"
	//
	return src
}

type CopyStruct1 struct {
	A string
	B string
}

type CopyStruct2 struct {
	A string
	B *string
}

func Test_CopyStruct(t *testing.T) {
	src := newCopyStructSrc()
	dst1 := newCopyStructDst()
	dst2 := newCopyStructDst()
	//
	CopyStruct(dst1, src)
	//
	if dst1.A != dst2.A ||
		dst1.B != dst2.B ||
		*dst1.C != *dst2.C ||
		dst1.D == nil || *dst1.D != *src.D ||
		dst1.E != dst2.E ||
		dst1.F.A != src.F.A || dst1.F.B != dst2.F.B ||
		dst1.G.A != src.G.A || dst1.G.B != dst2.G.B ||
		dst1.H != nil ||
		dst1.CopyStruct1.A != src.CopyStruct1.A || dst1.CopyStruct1.B != src.CopyStruct1.B {
		t.FailNow()
	}
}

func Test_CopyStructAll(t *testing.T) {
	src := newCopyStructSrc()
	dst1 := newCopyStructDst()
	dst2 := newCopyStructDst()
	//
	CopyStructAll(dst1, src)
	//
	if dst1.A != dst2.A ||
		dst1.B != src.B ||
		*dst1.C != *dst2.C ||
		dst1.D == nil || *dst1.D != *src.D ||
		dst1.E != dst2.E ||
		dst1.F.A != src.F.A || dst1.F.B != dst2.F.B ||
		dst1.G.A != src.G.A || dst1.G.B != dst2.G.B ||
		dst1.H != nil ||
		dst1.CopyStruct1.A != src.CopyStruct1.A || dst1.CopyStruct1.B != src.CopyStruct1.B {
		t.FailNow()
	}
}

type StructToMap1 struct {
	A int
	B *string
}

type StructToMap2 struct {
	A int
	B string
}

type StructToMap3 struct {
	A int
	B string
	StructToMap1
	C StructToMap2
	D *StructToMap1
	E *StructToMap2
	F *StructToMap2
}

func newStructToMap3() *StructToMap3 {
	p := new(StructToMap3)
	p.A = 1
	p.B = "b"
	p.StructToMap1.A = 2
	bb := "bb"
	p.StructToMap1.B = &bb
	p.C.A = 3
	p.C.B = "cb"
	p.D = new(StructToMap1)
	p.D.A = 4
	p.E = new(StructToMap2)
	p.E.A = 5
	p.E.B = "eb"
	//
	return p
}

func Test_StructToMap(t *testing.T) {
	s := newStructToMap3()
	m := StructToMap(s)
	//
	if v, ok := m["A"].(int); !ok || v != 1 {
		t.FailNow()
	}
	//
	if v, ok := m["B"].(string); !ok || v != "b" {
		t.FailNow()
	}
	//
	mm, ok := m["StructToMap1"].(map[string]any)
	if !ok {
		t.FailNow()
	}
	if v, ok := mm["A"].(int); !ok || v != 2 {
		t.FailNow()
	}
	if v, ok := mm["B"].(string); !ok || v != "bb" {
		t.FailNow()
	}
	//
	mm, ok = m["C"].(map[string]any)
	if !ok {
		t.FailNow()
	}
	if v, ok := mm["A"].(int); !ok || v != 3 {
		t.FailNow()
	}
	if v, ok := mm["B"].(string); !ok || v != "cb" {
		t.FailNow()
	}
	//
	mm, ok = m["D"].(map[string]any)
	if !ok {
		t.FailNow()
	}
	if v, ok := mm["A"].(int); !ok || v != 4 {
		t.FailNow()
	}
	if mm["B"] != nil {
		t.FailNow()
	}
	//
	mm, ok = m["E"].(map[string]any)
	if !ok {
		t.FailNow()
	}
	if v, ok := mm["A"].(int); !ok || v != 5 {
		t.FailNow()
	}
	if v, ok := mm["B"].(string); !ok || v != "eb" {
		t.FailNow()
	}
	//
	if m["F"] != nil {
		t.FailNow()
	}
}

func Test_StructToMapIgnore(t *testing.T) {
	s := newStructToMap3()
	m := StructToMapIgnore(s)
	//
	if v, ok := m["A"].(int); !ok || v != 1 {
		t.FailNow()
	}
	//
	if v, ok := m["B"].(string); !ok || v != "b" {
		t.FailNow()
	}
	//
	mm, ok := m["StructToMap1"].(map[string]any)
	if !ok {
		t.FailNow()
	}
	if v, ok := mm["A"].(int); !ok || v != 2 {
		t.FailNow()
	}
	if v, ok := mm["B"].(string); !ok || v != "bb" {
		t.FailNow()
	}
	//
	mm, ok = m["C"].(map[string]any)
	if !ok {
		t.FailNow()
	}
	if v, ok := mm["A"].(int); !ok || v != 3 {
		t.FailNow()
	}
	if v, ok := mm["B"].(string); !ok || v != "cb" {
		t.FailNow()
	}
	//
	mm, ok = m["D"].(map[string]any)
	if !ok {
		t.FailNow()
	}
	if v, ok := mm["A"].(int); !ok || v != 4 {
		t.FailNow()
	}
	if _, ok := mm["B"]; ok {
		t.FailNow()
	}
	//
	mm, ok = m["E"].(map[string]any)
	if !ok {
		t.FailNow()
	}
	if v, ok := mm["A"].(int); !ok || v != 5 {
		t.FailNow()
	}
	if v, ok := mm["B"].(string); !ok || v != "eb" {
		t.FailNow()
	}
	//
	if _, ok := mm["F"]; ok {
		t.FailNow()
	}
}
