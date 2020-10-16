package test

import (
	"encoding/json"
	"testing"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/agiledragon/gomonkey/v2/test/fake"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPatchPair(t *testing.T) {

	Convey("TestPatchPair", t, func() {

		Convey("TestPatchPair", func() {
			patchPairs := [][2]interface{}{
				{
					fake.Exec,
					func(_ string, _ ...string) (string, error) {
						return outputExpect, nil
					},
				},
				{
					json.Unmarshal,
					func(_ []byte, v interface{}) error {
						p := v.(*map[int]int)
						*p = make(map[int]int)
						(*p)[1] = 2
						(*p)[2] = 4
						return nil
					},
				},
			}
			patches := NewPatches()
			defer patches.Reset()
			for _, pair := range patchPairs {
				patches.ApplyFunc(pair[0], pair[1])
			}

			output, err := fake.Exec("", "")
			So(err, ShouldEqual, nil)
			So(output, ShouldEqual, outputExpect)

			var m map[int]int
			err = json.Unmarshal(nil, &m)
			So(err, ShouldEqual, nil)
			So(m[1], ShouldEqual, 2)
			So(m[2], ShouldEqual, 4)
		})

	})
}
