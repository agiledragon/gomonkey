package test

import (
	. "github.com/agiledragon/gomonkey"
	"github.com/agiledragon/gomonkey/test/fake"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPbBuilderFunc(t *testing.T) {
	Convey("TestPbBuilderFunc", t, func() {

		Convey("first dsl", func() {
			patches := NewPatches()
			defer patches.Reset()
			patchBuilder := NewPatchBuilder(patches)

			patchBuilder.
				Func(fake.Belong).
				Stubs().
				With(Eq("zxl"), Any()).
				Will(Return(true)).
				End()

			flag := fake.Belong("zxl", []string{})
			So(flag, ShouldBeTrue)

			defer func() {
				if p := recover(); p != nil {
					str, ok := p.(string)
					So(ok, ShouldBeTrue)
					So(str, ShouldEqual, "input paras ddd is not matched")
				}
			}()
			fake.Belong("ddd", []string{"abc"})
		})

	})
}
