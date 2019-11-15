package consumer

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	log "github.com/Sirupsen/logrus"

	"toolkit/worker/func_stuff"
	"toolkit/worker/job"
	"toolkit/worker/utils"

	"github.com/Shopify/sarama"
)

// ----------------------------------------------------
// Consumer: excute the async job in worker application
// ----------------------------------------------------
func JobExecutor(message *sarama.ConsumerMessage) error {
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
			fmt.Println(formatedErr, formatedErr)
		}
	}()

	var err error
	// Job unmarshaling
	var job job.Job
	err = json.Unmarshal(message.Value, &job)
	// Job unmarshaling failed, mark the correspond message as finish, and log err
	if err != nil {
		formatedErr := fmt.Errorf("Consumer unmarshal error: %v, mark the message as consumed.", err)
		log.Error(formatedErr)
		// message.Finish()
		return nil
	}

	// Will not retry job while message attmpts exceed job's MaxRetry
	// if int(message.Version) > job.MaxRetry {
	// 	message.Finish()
	// 	return nil
	// }

	// Execute the function
	err = FuncExecutor(&job)
	if err == nil {
		// message.Finish()
	} else {
		// Will not retry job while job markd Retry field with false
		if !job.Retry {
			formatedErr := fmt.Errorf("Consumer func executor error: %v, this is a non-retry message, mark the message as consumed.", err)
			log.Error(formatedErr)
			// message.Finish()
			return nil
		}

		// TODO: 重试机制
		// message.Requeue(-1)

		formatedErr := fmt.Errorf("Consumer func executor error: %v, requeue message instantly.", err)
		log.Error(formatedErr)
	}

	return err
}

// Ref. [https://github.com/RichardKnop/machinery/blob/bae667d798a33928db0349a5e2d098d20ce3b301/v1/worker.go]
func FuncExecutor(job *job.Job) error {
	registeredFunc, err := func_stuff.GetRegisteredFunc(job.FuncName)
	if err != nil {
		return err
	}

	reflectedFunc := reflect.ValueOf(registeredFunc)
	reflectedArgs, err := utils.GetDecodedArgs(job.Args)
	if err != nil {
		e := fmt.Errorf("Decode args: %v for func: %s  args: %v", err, job.FuncName, job.Args)
		return e
	}

	_, err = TryCall(reflectedFunc, reflectedArgs)
	if err != nil {
		e := fmt.Errorf("Worker try call registerd function: %s %v %v", job.FuncName, err, job.Args)
		return e
	}

	return nil
}

// TryCall attempts to call the task with the supplied arguments.
//
// `err` is set in the return value in two cases:
// 1. The reflected function invocation panics (e.g. due to a mismatched
//    argument list).
// 2. The task func itself returns a non-nil error.
func TryCall(f reflect.Value, args []reflect.Value) (results []reflect.Value, err error) {
	defer func() {
		// Recover from panic and set err.
		if e := recover(); e != nil {
			switch e := e.(type) {
			default:
				err = errors.New("Invoking task caused a panic")

			case error:
				err = e
			case string:
				err = errors.New(e)
			}
		}
	}()

	results = f.Call(args)

	if len(results) < 2 {
		// sc_raven.CaptureError(fmt.Errorf("warning!!! wrong async func define"))
	}

	// If an error was returned by the task func, propagate it
	// to the caller via err.
	if !results[1].IsNil() {
		return nil, results[1].Interface().(error)
	}

	return results, err
}
