package producer

import (
	"encoding/json"
	"fmt"
	"toolkit/worker/consts"
	"toolkit/worker/job"
	"toolkit/worker/utils"
	"toolkit/xkafka"

	"gitee.com/ikongjix/go_common/dlog"
)

var defaultProducer *xkafka.Producer

func getOrInitDefaultProducer(conf map[string]interface{}) *xkafka.Producer {
	if defaultProducer == nil {
		defaultProducer = xkafka.NewProducer()
		defaultProducer.Init(conf)
	}
	return defaultProducer
}

func AsyncExecute(function interface{}, topic string, args ...interface{}) (err error) {
	// TODO validate the param length and type for function
	if topic == "" {
		topic = consts.KafkaWorkerTopic
	}

	functionName := utils.GetFunctionName(function)
	argList, err := utils.GoVarsToYJsonSlice(args...)

	job := job.Job{
		FuncName: functionName,
		Args:     argList,
	}

	marshaledJob, err := json.Marshal(job)
	if err != nil {
		e := fmt.Errorf("Producer marshal error: %v", err)
		dlog.ERROR(e)
		return e
	}

	pd := getOrInitDefaultProducer(map[string]interface{}{"topic": ""})
	return pd.PutMessage(marshaledJob)
}
