package server

import (
	"testing"
)

func Test_StartPconnSrv(t *testing.T) {
	if e := StartPconnSrv(); e != nil { //try a unit test on function
		t.Error("start server failed...") // 如果不是如预期的那么就报错
	} else {
		t.Log("testing passed") //记录一些你期望记录的信息
	}
}
