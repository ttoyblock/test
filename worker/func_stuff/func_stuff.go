package func_stuff

import (
	"fmt"
)

// --------------------------------------------
// The functions registered here should
// 1. be idempotent
// 2. have two return values
//     - the type of first value doesn't matter
//     - the second should be error
// 3. have parameters in these type
//     - bool
//     - int
//     - int8
//     - int16
//     - int32
//     - int64
//     - uint
//     - uint8
//     - uint16
//     - uint32
//     - uint64
//     - float32
//     - float64
//     - string
// --------------------------------------------

var funcsMap = map[string]interface{}{
// utils.GetFunctionName(sms.SendSmsWithYunpian):                 sms.SendSmsWithYunpian,
}

func GetRegisteredFunc(registerdFuncName string) (interface{}, error) {
	registerdFunc, ok := funcsMap[registerdFuncName]
	if !ok {
		err := fmt.Errorf("[async worker func stuff]Currently, function %v is not registered", registerdFuncName)
		return nil, err
	}
	return registerdFunc, nil
}
