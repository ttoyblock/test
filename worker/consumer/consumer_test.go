package consumer

import (
	"encoding/json"
	"testing"

	"github.com/Shopify/sarama"

	"toolkit/server/post_server/env"
	"toolkit/worker/func_stuff"
	"toolkit/worker/job"
	"toolkit/worker/utils"

	"gitee.com/ikongjix/go_common/mysql_db"
	"github.com/gojuukaze/YTask/v2/util"
)

func TestJobExecutor(t *testing.T) {
	j := job.Job{
		FuncName: utils.GetFunctionName(func_stuff.TestFoo),
		Retry:    false,
		MaxRetry: 0,
	}

	ctx := env.MyContext{UltraxDBConf: &mysql_db.DbConfArr{Master: mysql_db.DbConf{DbDsn: "ahhahahahaaha"}}}
	args := []interface{}{ctx, 199}
	j.Args, _ = util.GoVarsToYJsonSlice(args...)

	bs, _ := json.Marshal(j)
	msg := &sarama.ConsumerMessage{Value: bs}

	ok, err := JobExecutor(msg)
	if err != nil {
		t.Error(err)
	}
	t.Log(ok)
}
