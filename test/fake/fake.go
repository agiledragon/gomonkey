package fake

import (
    "errors"
    "fmt"
    "os/exec"
    "strings"
)

var (
    ErrActual       = errors.New("actual")
    ErrElemExsit    = errors.New("elem already exist")
    ErrElemNotExsit = errors.New("elem not exist")
)

func Exec(cmd string, args ...string) (string, error) {
    cmdPath, err := exec.LookPath(cmd)
    if err != nil {
        fmt.Errorf("exec.LookPath err: %v, cmd: %s", err, cmd)
        return "", errors.New("any")
    }

    var output []byte
    output, err = exec.Command(cmdPath, args...).CombinedOutput()
    if err != nil {
        fmt.Errorf("exec.Command.CombinedOutput err: %v, cmd: %s", err, cmd)
        return "", errors.New("any")
    }
    fmt.Println("CMD[", cmdPath, "]ARGS[", args, "]OUT[", string(output), "]")
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


type Slice []int

func NewSlice() Slice {
    return make(Slice, 0)
}

func (this* Slice) Add(elem int) error {
    for _, v := range *this {
        if v == elem {
            fmt.Printf("Slice: Add elem: %v already exist\n", elem)
            return ErrElemExsit
        }
    }
    *this = append(*this, elem)
    fmt.Printf("Slice: Add elem: %v succ\n", elem)
    return nil
}

func (this* Slice) Remove(elem int) error {
    found := false
    for i, v := range *this {
        if v == elem {
            if i == len(*this) - 1 {
                *this = (*this)[:i]

            } else {
                *this = append((*this)[:i], (*this)[i+1:]...)
            }
            found = true
            break
        }
    }
    if !found {
        fmt.Printf("Slice: Remove elem: %v not exist\n", elem)
        return ErrElemNotExsit
    }
    fmt.Printf("Slice: Remove elem: %v succ\n", elem)
    return nil
}

func ReadLeaf(url string) (string, error) {
    output := fmt.Sprintf("%s, %s!", "Hello", "World")
    return output, nil
}

type Etcd struct {

}

func (this *Etcd) Retrieve(url string) (string, error) {
    output := fmt.Sprintf("%s, %s!", "Hello", "Etcd")
    return output, nil
}

var Marshal = func(v interface{}) ([]byte, error) {
    return nil, nil
}

type Db interface {
    Retrieve(url string) (string, error)
}

type Mysql struct {

}

func (this *Mysql) Retrieve(url string) (string, error) {
    output := fmt.Sprintf("%s, %s!", "Hello", "Mysql")
    return output, nil
}

func NewDb(style string) Db {
    if style == "etcd" {
        return new(Etcd)
    } else {
        return new(Mysql)
    }
}

type PrivateMethodStruct struct {

}

func (this *PrivateMethodStruct) ok() bool {
    return this != nil
}

func (this *PrivateMethodStruct) Happy() string {
    if this.ok() {
        return "happy"
    }
    return "unhappy"
}

func (this PrivateMethodStruct) haveEaten() bool {
    return this != PrivateMethodStruct{}
}

func (this PrivateMethodStruct) AreYouHungry() string {
    if this.haveEaten() {
        return "I am full"
    }

    return "I am hungry"
}