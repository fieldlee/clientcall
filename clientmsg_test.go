package main

import "testing"

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
	Sync(body,0,1,"")
}

func TestAsync(t *testing.T) {
	body := []byte("2+2=?")
	Async(body,1,1)
}

func TestBroadcast(t *testing.T) {
	body := []byte("3+3=?")
	Broadcast(body,"test_service")
}
