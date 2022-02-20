package test

import (
	"testing"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/agiledragon/gomonkey/v2/test/fake"
	. "github.com/smartystreets/goconvey/convey"
)

/*
   compare with apply_method_seq_test.go
*/

func TestApplyMethodReturn(t *testing.T) {
	e := &fake.Etcd{}
	Convey("TestApplyMethodReturn", t, func() {
		Convey("declares the values to be returned", func() {
			info1 := "hello cpp"
			patches := ApplyMethodReturn(e, "Retrieve", info1, nil)
			defer patches.Reset()
			for i := 0; i < 10; i++ {
				output1, err1 := e.Retrieve("")
				So(err1, ShouldEqual, nil)
				So(output1, ShouldEqual, info1)
			}

			patches.Reset() // if not reset will occur:patch has been existed
			info2 := "hello golang"
			patches.ApplyMethodReturn(e, "Retrieve", info2, nil)
			for i := 0; i < 10; i++ {
				output2, err2 := e.Retrieve("")
				So(err2, ShouldEqual, nil)
				So(output2, ShouldEqual, info2)
			}
		})
	})
}
