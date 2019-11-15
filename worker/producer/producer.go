package producer

import (
	"encoding/json"
	"fmt"
	"reflect"
	"toolkit/worker/config"
	"toolkit/worker/job"
	"toolkit/worker/utils"
	"toolkit/xkafka"
)

// AsyncExecute ... 异步执行，默认为本机队列
func AsyncExecute(function interface{}, retry bool, maxRetry int, topic string, args ...interface{}) error {
	return AsyncExecuteWithCertainProducer(function, retry, maxRetry, topic, args...)
}

// AsyncExecuteWithoutError ... 异步执行，没有错误返回
func AsyncExecuteWithoutError(function interface{}, args ...interface{}) {
	err := AsyncExecute(function, false, 1, "", args...)
	if err != nil {
		fmt.Println(err)
	}
}

// AsyncExecuteWithCertainProducer ... 推到某个队列步执行
// Producer: push async job to kafka in web application
func AsyncExecuteWithCertainProducer(function interface{}, retry bool, maxRetry int, topic string, args ...interface{}) error {
	// TODO validate the param length and type for function
	var err error

	if topic == "" {
		topic = config.WorkerTopic
	}

	functionName := utils.GetFunctionName(function)
	argList := make([]job.FuncArg, len(args))
	for i, a := range args {
		encodedArg, err := utils.ArgEncoder(a)
		if err != nil {
			err = fmt.Errorf("Producer ArgEncode error: %v", err)
			return err
		}
		argList[i] = encodedArg
	}

	job := job.Job{
		FuncName: functionName,
		Args:     argList,
		Retry:    retry,
		MaxRetry: maxRetry,
	}

	marshaledJob, err := json.Marshal(job)
	if err != nil {
		err = fmt.Errorf("Producer marshal error: %v", err)
		return err
	}

	// push job
	p := xkafka.NewProducer()
	p.Init(map[string]interface{}{})
	err = p.PutMessage(marshaledJob)
	if err != nil {
		e := fmt.Errorf("Producer publish error: %v", err)
		fmt.Println(e)

		// sync call
		in := make([]reflect.Value, len(args))
		for k, arg := range args {
			in[k] = reflect.ValueOf(arg)
		}
		f := reflect.ValueOf(function)
		res := f.Call(in)
		outLength := len(res)
		if outLength != 2 {
			panic(fmt.Errorf("Invalid async func: %s", functionName))
		}

		r, re := res[0], res[1]
		fmt.Println(fmt.Sprintf("sync call: %s \nret: %v\nerror:%v\n", functionName, r, re))
		return err
	}

	return nil
}
