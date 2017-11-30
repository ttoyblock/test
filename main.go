package main

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"toolkit/net"

	"github.com/julienschmidt/httprouter"
)

func main() {
	fmt.Println("----")
	a := make([]int, 0)
	for _, v := range a {
		fmt.Println(v)
	}

	type T struct {
		A int
		B string
	}
	t := T{23, "skidoo"}
	s := reflect.ValueOf(&t).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%d: %s %s = %v\n", i, typeOfT.Field(i).Name, f.Type(), f.Interface())
	}

	router := httprouter.New()
	router.GET("/", net.Index)
	router.GET("/hello/:name", net.Hello)

	log.Fatal(http.ListenAndServe(":8080", router))
}
