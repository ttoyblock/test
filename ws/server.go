package ws

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

var Conns = make(map[string]*UnitConn)

type Server struct {
	addr          string
	servId        uint32
	cidSecret     string
	enableTestEnv bool
	hostname      string

	ackTimeout       int
	readTimeout      int
	writeTimeout     int
	heartbeatTimeout int

	wg *sync.WaitGroup
}

func NewServer() (s *Server, err error) {
	hostname, _ := os.Hostname()
	s = &Server{
		addr:     ":7777",
		hostname: hostname,
		wg:       new(sync.WaitGroup),
	}

	return s, nil
}

func (s *Server) Run() {
	s.listen(s.addr)
}

func (s *Server) Stop() {
	s.wg.Wait()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) listen(addr string) error {
	s.wg.Add(1)
	defer s.wg.Done()

	http.HandleFunc("/v1/chat", serveWs)
	err := http.ListenAndServe(addr, nil)
	return err
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		wsConn *websocket.Conn
		conn   *UnitConn
	)

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", 405)
		return
	}

	dv := r.Header.Get("device")
	if dv == "" {
		http.Error(w, "device is empty", 400)
		return
	}

	if wsConn, err = upgrader.Upgrade(w, r, nil); err != nil {
		fmt.Printf("upgrade websocket, err: %v, req header: %v", err, r.Header)
		return
	}

	conn = InitUnitConn(wsConn)
	// TODO: store conn
	Conns[conn.device] = conn
	fmt.Println("open new conn", conn.addr)
}
