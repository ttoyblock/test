package main

import (
	"syscall"
	"toolkit/worker/consumer"
	"toolkit/xkafka"

	log "github.com/Sirupsen/logrus"
	"github.com/vrecan/death"
)

// // Initial env from env package
// func init() {
// 	env.InitRunEnv()
// }

// ------------------------------------------------------------
// Run the worker as a independent program
// Should handel exit/interrupt logics: disconnect and stop
// Refs.
// [https://github.com/augurysys/go-nsqworker/blob/master/nsqworker.go#L69]
// [https://github.com/augurysys/go-nsqworker/blob/ad966d87ac429edeef7f4782dc995e7a476489e5/tool/main.go#L42]
// ------------------------------------------------------------

func main() {
	d := death.NewDeath(syscall.SIGINT, syscall.SIGTERM)

	// new consumer
	conf := map[string]interface{}{
		"topic":     "wallet_pay_queue",
		"groupName": "wallet_consumer",
		"broker":    "47.95.14.131:9182",
	}
	csm := xkafka.NewConsumer()
	csm.Init(conf)

	// run
	csm.RunByCall(consumer.JobExecutor)

	d.WaitForDeathWithFunc(func() {
		log.Info("Trying to stop worker and disconnect...")
		csm.Close()
		log.Info("Successfuly stop worker.")
	})
}
