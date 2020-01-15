package func_stuff

import (
	"fmt"

	"toolkit/worker/utils"
)

// --------------------------------------------
// The functions registered here should
// 1. be idempotent
// 2. return
//     - at least one return value
//     - the last one must be an error
// 3. have parameters in these type
//     - any type
// --------------------------------------------

var funcsMap = map[string]interface{}{
	utils.GetFunctionName(TestFoo): TestFoo,
}

func GetRegisteredFunc(registerdFuncName string) (interface{}, error) {
	registerdFunc, ok := funcsMap[registerdFuncName]
	if !ok {
		err := fmt.Errorf("[async worker func stuff]Currently, function %v is not registered", registerdFuncName)
		return nil, err
	}
	return registerdFunc, nil
}

func TestFoo(a int64) error {
	fmt.Println("-------testfoo res", "a", a)
	return nil
}
