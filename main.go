package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	// fmt.Println("----")
	// a := make([]int, 0)
	// for _, v := range a {
	// 	fmt.Println(v)
	// }
	json.Marshal()
	// type T struct {
	// 	A int
	// 	B string
	// }
	// t := T{23, "skidoo"}
	// s := reflect.ValueOf(&t).Elem()
	// typeOfT := s.Type()
	// for i := 0; i < s.NumField(); i++ {
	// 	f := s.Field(i)
	// 	fmt.Printf("%d: %s %s = %v\n", i, typeOfT.Field(i).Name, f.Type(), f.Interface())
	// }

	// router := httprouter.New()
	// router.GET("/", net.Index)
	// router.GET("/hello/:name", net.Hello)
	//
	// log.Fatal(http.ListenAndServe(":8080", router))

	// sli := []int{1, 5, 23, 6, 2, 65, 17, 123, 4, 9, 2, 234}
	// utils.ListSort(sli, false)
	// fmt.Println(sli)

	bs, _ := ioutil.ReadFile("/Users/libc/work/gocode/src/wallet-server-go/uid.txt")
	con := strings.Replace(string(bs), "\n", ",", -1)
	fmt.Println(con)

	// a := make([]*Asd, 0)
	// for i := 0; i < 4; i++ {
	// 	a = append(a, &Asd{A: i})
	// }
	// fmt.Printf("%p \n", a)
	// fmt.Printf("%p \n", &a)
	// for _, v := range a {
	// 	fmt.Println(v)
	// 	fmt.Printf("%p \n", v)
	// }
}

// type Asd struct {
// 	A int
// }
