package util

import "testing"

func Test_HTTPQuery(t *testing.T) {
	s := struct {
		A int    `query:"aa"`
		B string `query:"bb"`
	}{
		A: 12,
		B: "abc",
	}
	q := HTTPQuery(&s, nil)
	if q.Get("aa") != "12" ||
		q.Get("bb") != "abc" {
		t.FailNow()
	}
}
