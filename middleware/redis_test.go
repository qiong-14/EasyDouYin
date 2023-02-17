package middleware

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
)

func init() {
	InitRedis()
}

func TestSet(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			t.Error(err)
		}
	}()

	for i := 0; i < 10; i++ {
		res, err := set(0, strconv.Itoa(i), strconv.Itoa(i))
		if err != nil {
			t.Error(res)
		}
		fmt.Println(res)
	}
	for i := 0; i < 10; i++ {
		res, err := get(0, strconv.Itoa(i))
		if err != nil {
			t.Error(res)
		}
		if res != strconv.Itoa(i) {
			panic(errors.New("数据错误"))
		}
		fmt.Println(res)
	}

	res, err := flushDB(0)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)
}
