package test

import (
    . "github.com/agiledragon/gomonkey"
    . "github.com/smartystreets/goconvey/convey"
    "testing"
    "os/exec"
    "fmt"
    "errors"
    "strings"
)

var (
    ErrActual = errors.New("actual")
    outputExpect = "xxx-vethName100-yyy"
)

func Exec(cmd string, args ...string) (string, error) {
    cmdpath, err := exec.LookPath(cmd)
    if err != nil {
        fmt.Errorf("exec.LookPath err: %v, cmd: %s", err, cmd)
        return "", errors.New("any")
    }

    var output []byte
    output, err = exec.Command(cmdpath, args...).CombinedOutput()
    if err != nil {
        fmt.Errorf("exec.Command.CombinedOutput err: %v, cmd: %s", err, cmd)
        return "", errors.New("any")
    }
    fmt.Println("CMD[", cmdpath, "]ARGS[", args, "]OUT[", string(output), "]")
    return string(output), nil
}

func Belong(points string, lines []string) bool {
    flag := false
    for _, line := range lines {
        flag = true
        for _, r := range points {
            if !strings.ContainsRune(line, r) {
                flag = false
                break
            }
        }
        if flag {
            return true
        }
    }
    return false
}

func TestApplyFunc(t *testing.T) {
    Convey("TestApplyFunc", t, func() {

        Convey("one func for succ", func() {
            patches := ApplyFunc(Exec, func(_ string, _ ...string) (string, error) {
                    return outputExpect, nil
                })
            defer patches.Reset()
            output, err := Exec("", "")
            So(err, ShouldEqual, nil)
            So(output, ShouldEqual, outputExpect)
        })

        Convey("one func for fail", func() {
            patches := ApplyFunc(Exec, func(_ string, _ ...string) (string, error) {
                return "", ErrActual
            })
            defer patches.Reset()
            output, err := Exec("", "")
            So(err, ShouldEqual, ErrActual)
            So(output, ShouldEqual, "")
        })

        Convey("two funcs", func() {
            patches := ApplyFunc(Exec, func(_ string, _ ...string) (string, error) {
                return outputExpect, nil
            })
            defer patches.Reset()
            patches.ApplyFunc(Belong, func(_ string, _ []string) bool {
                return true
            })
            output, err := Exec("", "")
            So(err, ShouldEqual, nil)
            So(output, ShouldEqual, outputExpect)
            flag := Belong("", nil)
            So(flag, ShouldBeTrue)
        })
    })
}
