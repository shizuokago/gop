package gop_test

import (
	"testing"

	"github.com/shizuokago/gop"
)

func TestNewVersion(t *testing.T) {

	v1 := gop.NewVersion("1.2.3")
	if v1.IsError() {
		t.Errorf("%v not error", v1)
	}

	v2 := gop.NewVersion("1.2.4")
	if v2.IsError() {
		t.Errorf("%v not error", v2)
	}

	if !v1.Lt(v2) {
		t.Errorf("%v less then %v", v1, v2)
	}
	if !v1.Le(v2) {
		t.Errorf("%v less equals %v", v1, v2)
	}

	if v1.Gt(v2) {
		t.Errorf("%v not greater then %v", v1, v2)
	}
	if v1.Ge(v2) {
		t.Errorf("%v not greater then %v", v1, v2)
	}

	if v1.Eq(v2) {
		t.Errorf("%v not equals %v", v1, v2)
	}
	if v2.Eq(v1) {
		t.Errorf("%v not equals %v", v1, v2)
	}

	v3 := gop.NewVersion("1.2.4")
	if v3.IsError() {
		t.Errorf("%v not error", v3)
	}

	if !v2.Eq(v3) {
		t.Errorf("%v equals %v", v2, v3)
	}
	if !v3.Eq(v2) {
		t.Errorf("%v equals %v", v2, v3)
	}

	if !v3.Le(v2) {
		t.Errorf("%v less equals %v", v2, v3)
	}
	if !v3.Ge(v2) {
		t.Errorf("%v greater equals %v", v2, v3)
	}

	if !v1.IsVersion(1) {
		t.Errorf("%v is v1", v1)
	}

	if v1.IsVersion(2) {
		t.Errorf("%v is not v2", v1)
	}

	if v1.IsVersion(0) {
		t.Errorf("%v is not v0", v1)
	}

	v4 := gop.NewVersion("1.12.3")
	if v4.IsError() {
		t.Errorf("%v not error", v4)
	}

	if !v1.Lt(v4) {
		t.Errorf("%v less then %v", v1, v4)
	}
}

func TestBuild(t *testing.T) {

	v1 := gop.NewVersion("0.0.0-20180824175216-6c1c5e93cdc1")
	if v1.IsError() {
		t.Errorf("%v not error", v1)
	}
	v2 := gop.NewVersion("0.0.0-20190408220357-e5b8258f4918")
	if v2.IsError() {
		t.Errorf("%v not error", v2)
	}
	v3 := gop.NewVersion("0.0.0-20200226224502-204d844ad48d")
	if v3.IsError() {
		t.Errorf("%v not error", v3)
	}

	if !v1.Le(v2) {
		t.Errorf("%v less then %v", v1, v2)
	}
	if !v2.Le(v3) {
		t.Errorf("%v less then %v", v2, v3)
	}
	if !v1.Le(v3) {
		t.Errorf("%v less then %v", v1, v3)
	}

	v4 := gop.NewVersion("0.0.0-20200226224502-204d844ad48d")
	if !v3.Eq(v4) {
		t.Errorf("%v equals %v", v3, v4)
	}

	if !v1.EqCore(v2) {
		t.Errorf("core:%v equals %v", v1, v2)
	}
	if !v1.EqCore(v3) {
		t.Errorf("core:%v equals %v", v1, v3)
	}
	if !v1.EqCore(v4) {
		t.Errorf("core:%v equals %v", v1, v4)
	}
	if !v2.EqCore(v3) {
		t.Errorf("core:%v equals %v", v2, v3)
	}
	if !v2.EqCore(v4) {
		t.Errorf("core:%v equals %v", v2, v4)
	}
	if !v3.EqCore(v4) {
		t.Errorf("core:%v equals %v", v3, v4)
	}
}
