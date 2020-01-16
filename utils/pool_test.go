package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestAsyncPool_Add(t *testing.T) {
	body1:= AsyncBody{
		Uid:"123",
		Body:[]byte("123"),
	}
	body2:= AsyncBody{
		Uid:"234",
		Body:[]byte("234"),
	}
	body3:= AsyncBody{
		Uid:"345",
		Body:[]byte("345"),
	}
	go func() {
		<- time.After(3*time.Second)
		Queue.Add(body1)
		<- time.After(2*time.Second)
		Queue.Add(body2)
		<- time.After(4*time.Second)
		Queue.Add(body3)
	}()

}

func TestAsyncPool_Get(t *testing.T) {
	body1 := Queue.Get("123")
	fmt.Println(fmt.Sprintf("%s",body1))
	body2 := Queue.Get("234")
	fmt.Println(fmt.Sprintf("%s",body2))
	body3 := Queue.Get("345")
	fmt.Println(fmt.Sprintf("%s",body3))
}

func TestAll(t *testing.T){
	TestAsyncPool_Get(t)
	<- time.After(2*time.Second)
	TestAsyncPool_Add(t)
}