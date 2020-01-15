package job

import "github.com/fatih/structs"

type Job struct {
	FuncName string
	Args     []string
	Retry    bool
	MaxRetry int
}

func (job *Job) ToMap() map[string]interface{} {
	return structs.Map(job)
}
