package main

import (
	"fmt"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	Register("10005")

}

func TestPublish(t *testing.T) {
	Publish("test_service")
}

func TestSubscribe(t *testing.T) {
	Subscribe("test_service")
}

func TestSync(t *testing.T) {

	for i:=0 ; i<10;i++{
		body := []byte(fmt.Sprintf("%d+%d",i,i))
		<- time.After(time.Microsecond*30)
		Sync(body,0,1,"")
		<- time.After(time.Microsecond*25)
	}

	//Sync(body,0,1,"")
}

func TestAsync(t *testing.T) {
	cd := func(body []byte)int{
		fmt.Println(fmt.Sprintf("callback：%s",body))
		fmt.Println("====callback=====")
		return 1
	}

	for i:=0 ; i<1000;i++ {
		body := []byte(fmt.Sprintf("%d+%d",i,i))
		<- time.After(time.Microsecond*30)
		Async(body, 1, 1,cd)
	}
}

func TestBroadcast(t *testing.T) {
	body := []byte("3+3=?")
	Broadcast(body,"test_service")
}
