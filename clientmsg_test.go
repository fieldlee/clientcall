package main

import (
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
	body := []byte("1+1=?")
	for i:=0 ; i<10000;i++{
		<- time.After(time.Microsecond*30)
		go Sync(body,0,1,"")
		<- time.After(time.Microsecond*25)
	}

}

func TestAsync(t *testing.T) {
	body := []byte("2+2=?")
	for i:=0 ; i<100000;i++ {
		<- time.After(time.Microsecond*30)
		Async(body, 1, 1)
	}
}

func TestBroadcast(t *testing.T) {
	body := []byte("3+3=?")
	Broadcast(body,"test_service")
}
