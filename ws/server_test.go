package ws_test

import (
	"fmt"
	"testing"
	"toolkit/ws"
)

func Test_Server(t *testing.T) {
	s, err := ws.NewServer()
	if err != nil {
		t.Error(err)
	}
	go s.Run()

	// 启动线程，不断发消息
	// go func() {
	// 	var (
	// 		err error
	// 	)
	// 	for {
	// 		if err = conn.WriteMessage([]byte("heartbeat")); err != nil {
	// 			return
	// 		}
	// 		time.Sleep(1 * time.Second)
	// 	}
	// }()

	for {
		if _, ok := ws.Conns[0]; !ok {
			continue
		}
		data, err := ws.Conns[0].ReadMessage()
		if err != nil {
			fmt.Println(err)
		}
		data = append(data, []byte(" hehe")...)
		if err = ws.Conns[0].WriteMessage(data); err != nil {
			fmt.Println(err)
		}
	}
}
