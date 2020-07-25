package test

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	. "github.com/henrylee2cn/gomonkey"
	"github.com/henrylee2cn/gomonkey/test/fake"
)

var (
	outputExpect = fake.OutputExpect
)

func TestApplyFunc(t *testing.T) {
	Convey("TestApplyFunc", t, func() {

		Convey("one func for succ", func() {
			patches := ApplyFunc(fake.Exec, func(_ string, _ ...string) (string, error) {
				return outputExpect, nil
			})
			defer patches.Reset()
			output, err := fake.Exec("", "")
			So(err, ShouldEqual, nil)
			So(output, ShouldEqual, outputExpect)
		})

		Convey("one func for fail", func() {
			patches := ApplyFunc(fake.Exec, func(_ string, _ ...string) (string, error) {
				return "", fake.ErrActual
			})
			defer patches.Reset()
			output, err := fake.Exec("", "")
			So(err, ShouldEqual, fake.ErrActual)
			So(output, ShouldEqual, "")
		})

		Convey("two funcs", func() {
			patches := ApplyFunc(fake.Exec, func(_ string, _ ...string) (string, error) {
				return outputExpect, nil
			})
			defer patches.Reset()
			patches.ApplyFunc(fake.Belong, func(_ string, _ []string) bool {
				return true
			})
			output, err := fake.Exec("", "")
			So(err, ShouldEqual, nil)
			So(output, ShouldEqual, outputExpect)
			flag := fake.Belong("", nil)
			So(flag, ShouldBeTrue)
		})

		Convey("input and output param", func() {
			patches := ApplyFunc(json.Unmarshal, func(data []byte, v interface{}) error {
				if data == nil {
					panic("input param is nil!")
				}
				p := v.(*map[int]int)
				*p = make(map[int]int)
				(*p)[1] = 2
				(*p)[2] = 4
				return nil
			})
			defer patches.Reset()
			var m map[int]int
			err := json.Unmarshal([]byte("123"), &m)
			So(err, ShouldEqual, nil)
			So(m[1], ShouldEqual, 2)
			So(m[2], ShouldEqual, 4)
		})
	})
}

func TestApplyFuncByTargetName(t *testing.T) {
	Convey("TestApplyFuncByTargetName", t, func() {

		Convey("one func for succ", func() {
			patches := ApplyFuncByTargetName("github.com/henrylee2cn/gomonkey/test/fake.Exec", func(_ string, _ ...string) (string, error) {
				return outputExpect, nil
			})
			defer patches.Reset()
			output, err := fake.Exec("", "")
			So(err, ShouldEqual, nil)
			So(output, ShouldEqual, outputExpect)
		})

		Convey("one func for fail", func() {
			patches := ApplyFuncByTargetName("github.com/henrylee2cn/gomonkey/test/fake.Exec", func(_ string, _ ...string) (string, error) {
				return "", fake.ErrActual
			})
			defer patches.Reset()
			output, err := fake.Exec("", "")
			So(err, ShouldEqual, fake.ErrActual)
			So(output, ShouldEqual, "")
		})

		Convey("two funcs", func() {
			patches := ApplyFuncByTargetName("github.com/henrylee2cn/gomonkey/test/fake.Exec", func(_ string, _ ...string) (string, error) {
				return outputExpect, nil
			})
			defer patches.Reset()
			patches.ApplyFunc(fake.Belong, func(_ string, _ []string) bool {
				return true
			})
			output, err := fake.Exec("", "")
			So(err, ShouldEqual, nil)
			So(output, ShouldEqual, outputExpect)
			flag := fake.Belong("", nil)
			So(flag, ShouldBeTrue)
		})

		Convey("input and output param", func() {
			patches := ApplyFuncByTargetName("encoding/json.Unmarshal", func(data []byte, v interface{}) error {
				if data == nil {
					panic("input param is nil!")
				}
				p := v.(*map[int]int)
				*p = make(map[int]int)
				(*p)[1] = 2
				(*p)[2] = 4
				return nil
			})
			defer patches.Reset()
			var m map[int]int
			err := json.Unmarshal([]byte("123"), &m)
			So(err, ShouldEqual, nil)
			So(m[1], ShouldEqual, 2)
			So(m[2], ShouldEqual, 4)
		})
	})
}

func TestApplyFuncByDoubleName(t *testing.T) {
	Convey("TestApplyFuncByDoubleName", t, func() {
		Convey("one func for succ", func() {
			patches := ApplyFuncByDoubleName(fake.Exec, "github.com/henrylee2cn/gomonkey/test/fake.execDouble1")
			defer patches.Reset()
			output, err := fake.Exec("", "")
			So(err, ShouldEqual, nil)
			So(output, ShouldEqual, outputExpect)
		})

		Convey("one func for fail", func() {
			patches := ApplyFuncByDoubleName(fake.Exec, "github.com/henrylee2cn/gomonkey/test/fake.execDouble2")
			defer patches.Reset()
			output, err := fake.Exec("", "")
			So(err, ShouldEqual, fake.ErrActual)
			So(output, ShouldEqual, "")
		})

		Convey("two funcs", func() {
			patches := ApplyFuncByDoubleName(fake.Exec, "github.com/henrylee2cn/gomonkey/test/fake.execDouble1")
			defer patches.Reset()
			patches.ApplyFunc(fake.Belong, func(_ string, _ []string) bool {
				return true
			})
			output, err := fake.Exec("", "")
			So(err, ShouldEqual, nil)
			So(output, ShouldEqual, outputExpect)
			flag := fake.Belong("", nil)
			So(flag, ShouldBeTrue)
		})

		Convey("input and output param", func() {
			patches := ApplyFuncByDoubleName(json.Unmarshal, "github.com/henrylee2cn/gomonkey/test/fake.unmarshalDouble")
			defer patches.Reset()
			var m map[int]int
			err := json.Unmarshal([]byte("123"), &m)
			So(err, ShouldEqual, nil)
			So(m[1], ShouldEqual, 2)
			So(m[2], ShouldEqual, 4)
		})
	})
}
