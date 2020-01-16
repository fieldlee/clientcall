package utils

import (
	"sync"
)

type BodyPool struct {
	sync.RWMutex
	BodyList []AsyncBody
}

type AsyncBody struct {
	Body []byte
	Uid string
}

type AsyncPool struct {
	Queue chan AsyncBody
}

var BodyList BodyPool
var Queue AsyncPool

func init()  {
	BodyList = BodyPool{
		BodyList:make([]AsyncBody,0),
	}

	Queue = AsyncPool{
		Queue:make(chan AsyncBody, 1024),
	}

	for{
		select {
		case body := <-Queue.Queue:
			BodyList.Append(body)
		}
	}

}
func (b *BodyPool)Append(body AsyncBody){
	b.Lock()
	b.BodyList = append(b.BodyList,body)
	b.Unlock()
}

func (b BodyPool)Check(uid string)bool{
	for _,body := range b.BodyList{
		if body.Uid != uid {
			return true
		}
	}
	return false
}

func (b *BodyPool)Remove(uid string)AsyncBody{
	b.Lock()
	defer b.Unlock()
	var targetBody AsyncBody
	targe := b.BodyList[:0]
	for _,body := range b.BodyList{
		if body.Uid != uid {
			targe = append(targe,body)
		}else{
			targetBody = body
		}
	}
	b.BodyList = targe

	return targetBody
}


func (q *AsyncPool) Add(body AsyncBody){
	q.Queue <- body
}
