package test

import (
	"testing"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/agiledragon/gomonkey/v2/test/fake"
	. "github.com/smartystreets/goconvey/convey"
)

/*
  compare with apply_func_var_seq_test.go
*/
func TestApplyFuncVarReturn(t *testing.T) {
	Convey("TestApplyFuncVarReturn", t, func() {

		Convey("declares the values to be returned", func() {
			info1 := "hello cpp"

			patches := ApplyFuncVarReturn(&fake.Marshal, []byte(info1), nil)
			defer patches.Reset()
			for i := 0; i < 10; i++ {
				bytes, err := fake.Marshal("")
				So(err, ShouldEqual, nil)
				So(string(bytes), ShouldEqual, info1)
			}

			info2 := "hello golang"
			patches.ApplyFuncVarReturn(&fake.Marshal, []byte(info2), nil)
			for i := 0; i < 10; i++ {
				bytes, err := fake.Marshal("")
				So(err, ShouldEqual, nil)
				So(string(bytes), ShouldEqual, info2)
			}
		})

	})
}
