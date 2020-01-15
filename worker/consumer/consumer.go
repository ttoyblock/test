package consumer

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"toolkit/worker/func_stuff"
	"toolkit/worker/job"
	"toolkit/xkafka"

	"github.com/Shopify/sarama"
)

// ----------------------------------------------------
// Consumer: excute the async job in worker application
// ----------------------------------------------------
func LaunchConsumer(conf map[string]interface{}) {
	kafkaConsumer := xkafka.NewConsumer()
	kafkaConsumer.Init(conf)
	defer kafkaConsumer.Close()

	kafkaConsumer.RunByCall(func(msg *sarama.ConsumerMessage) bool {
		b, err := JobExecutor(msg)
		checkErr(err)
		return b
	})
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func JobExecutor(msg *sarama.ConsumerMessage) (ok bool, err error) {
	defer func() {
		// Recover from panic and set err.
		if e := recover(); e != nil {

			var perr error
			switch e := e.(type) {
			default:
				if realErr, succ := e.(error); succ {
					perr = realErr
				} else {
					perr = errors.New("panic")
				}
			case error:
				perr = e
			case string:
				perr = errors.New(e)
			}

			formatedErr := fmt.Errorf("[Panic] panic occurs while execute scheduler job with msg: %v", perr)
			log.Fatalln(formatedErr)
		}
	}()
	log.Println("JobExecutor NewMessage", string(msg.Value))

	var job job.Job
	err = json.Unmarshal(msg.Value, &job)
	if err != nil {
		err = fmt.Errorf("Consumer unmarshal error: %v, mark the message as consumed.", err)
		log.Fatalln(err)
		return
	}

	err = FuncExecutor(&job)
	if err != nil {
		log.Fatalln("FuncExecutor err", err)
	}
	return err == nil, err
}

// Ref. [https://github.com/RichardKnop/machinery/blob/bae667d798a33928db0349a5e2d098d20ce3b301/v1/worker.go]
func FuncExecutor(job *job.Job) error {
	registeredFunc, err := func_stuff.GetRegisteredFunc(job.FuncName)
	if err != nil {
		return err
	}

	funcValue := reflect.ValueOf(registeredFunc)

	var inStart = 0
	var inValue []reflect.Value

	inValue, err = utils.GetCallInArgs(funcValue, job.Args, inStart)
	if err != nil {
		err = fmt.Errorf("Decode args: %v for func: %s  args: %v", err, job.FuncName, job.Args)
		return err
	}

	outValue, err := TryCall(funcValue, inValue)
	if err != nil {
		err = fmt.Errorf("Worker try call registerd function: %s %v %v", job.FuncName, err, job.Args)
		return err
	}

	result, err := utils.GoValuesToYJsonSlice(outValue)
	if err != nil {
		err = fmt.Errorf("Decode result: %v for func: %s  args: %v", err, job.FuncName, job.Args)
		log.Fatalln(err)
	}
	log.Println("task result", result)

	return nil
}

// TryCall attempts to call the task with the supplied arguments.
//
// `err` is set in the return value in two cases:
// 1. The reflected function invocation panics (e.g. due to a mismatched
//    argument list).
// 2. The task func itself returns a non-nil error.
func TryCall(f reflect.Value, args []reflect.Value) (results []reflect.Value, err error) {
	//defer func() {
	//	// Recover from panic and set err.
	//	if e := recover(); e != nil {
	//		switch e := e.(type) {
	//		default:
	//			err = errors.New("Invoking task caused a panic")
	//
	//		case error:
	//			err = e
	//		case string:
	//			err = errors.New(e)
	//		}
	//	}
	//}()

	results = f.Call(args)

	if len(results) < 2 {
		log.Fatalln(fmt.Errorf("warning!!! wrong async func define"))
	}

	// If an error was returned by the task func, propagate it
	// to the caller via err.
	if !results[len(results)-1].IsNil() {
		return nil, results[1].Interface().(error)
	}
	return
}
