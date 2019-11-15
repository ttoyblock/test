package job

type Job struct {
	FuncName string
	Args     []FuncArg
	Retry    bool
	MaxRetry int
}

type FuncArg struct {
	Type  string
	Value interface{}
}
