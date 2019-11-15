### Worker 部署备忘

#### 当前情况

1. 已经基本实现了基于 NSQ 的异步 worker
2. 在本地依次开启以下服务，可以使 worker 正常工作
    - `nsqlookupd`
    - `nsqd --lookupd-tcp-address=127.0.0.1:4160`
    - `nsqadmin --lookupd-http-address=127.0.0.1:4161`
    - 进入 libs/worker 目录，`go run main.go`
3. worker 的 logger 使用的是 `github.com/Sirupsen/logrus`，由于 `logrus` 没有实现 `Output` 方法，故写了一个 adaptor 来实现 logrus logger 与 nsq-logger 的转换，之后可能用同样的方法来实现 Beggo 的 logger 而不是使用 `logrus`
4. 对于 worker 的 `Job` 处理逻辑部分，可能需要根据实际情况作出调整

#### 基本线上部署方案

1. 部署 2-3 台 nsqlookupd
2. 与每台 API 服务器一起，部署 N 个 nsqd 服务
3. 有以下资料可以参考:
    - [NSQ 文档](http://nsq.io/overview/quick_start.html)
    - [NSQ Go Client 文档](https://godoc.org/github.com/nsqio/go-nsq)
4. 需要对 `libs/worker/worker_config.go` 做相应的更改
5. 关键字:
    - RDY
    - MaxInFlight
    - Attempts
6. 应当在认真参考 `服务端` 和 `客户端` 文档之后，再进行部署、调优