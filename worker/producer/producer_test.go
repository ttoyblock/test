package producer

import (
	"testing"
	"time"

	"gitee.com/ikongjix/go_common/dlog"

	"toolkit/worker/func_stuff"
)

func TestAsyncExecute(t *testing.T) {
	dlog.SetLog(dlog.SetLogConf{
		Prefix: "worker",
	})

	err := AsyncExecute(func_stuff.TestFunc, "", time.Now().Unix())
	if err != nil {
		t.Error(err)
	}
}
