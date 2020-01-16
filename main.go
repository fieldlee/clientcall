package main
import "C"
import (

	"clientcall/handle"
	pb "clientcall/proto"
	"clientcall/utils"

	"context"
	"fmt"
	"github.com/pborman/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
	"unsafe"
)

var (
	Host = utils.Address
	Port = fmt.Sprintf("%d",utils.Port)
)

type MsgHandle struct {}

//export GenerateMemory
func GenerateMemory(n int)unsafe.Pointer{
	p := C.malloc(C.sizeof_char * C.ulong(n))
	return unsafe.Pointer(p)
}

func (m *MsgHandle)Call(ctx context.Context, info *pb.CallReqInfo) (*pb.CallRspInfo, error) {
	out := pb.CallRspInfo{}
	//fmt.Println("info.Clientip:",string(info.Clientip))
	//fmt.Println("info.Service:",string(info.Service))
	//fmt.Println("info.Uuid:",string(info.Uuid))
	//fmt.Println(fmt.Sprintf("info.M_Msg:%s",string(info.M_Body.M_Msg)))
	if info.Uuid == ""{
		out.M_Net_Rsp = []byte(fmt.Sprintf("call result is:%s",string(info.M_Body.M_Msg)))
		return &out,nil
	}else{
		go func(uid string,body []byte) {
			resultByte := []byte(fmt.Sprintf("async call result for:%s",string(body)))
			infoBody , err := handle.MarshalBody(resultByte,0,0)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			if err != nil {
				fmt.Println(fmt.Sprintf("MarshalBody error : %s",err.Error()))
				return
			}

			syncResult,err := handle.CallSync(infoBody,uid)

			if err != nil {
				fmt.Println(fmt.Sprintf("call.CallSync error : %s",err.Error()))
				return
			}

			if syncResult.M_Err != nil {
				fmt.Println(fmt.Sprintf("syncResult error : %s",syncResult.M_Err))
				return
			}
			return
		}(info.Uuid,info.M_Body.M_Msg)
		out.M_Net_Rsp = []byte{0}
		return &out,nil
	}
}


func (m *MsgHandle)AsyncAnswer(ctx context.Context, resultInfo *pb.CallReqInfo) (*pb.CallRspInfo, error) {
	out := pb.CallRspInfo{}
	//fmt.Println("info.Clientip:",string(resultInfo.Clientip))
	//fmt.Println("info.Service:",string(resultInfo.Service))
	//fmt.Println("info.Uuid:",string(resultInfo.Uuid))
	//fmt.Println("info.M_Msg:",string(resultInfo.M_Body.M_Msg))
	out.M_Net_Rsp = []byte(fmt.Sprintf("async answer result for:%s",string(resultInfo.Clientip)))
	body := utils.AsyncBody{
		Body:out.M_Net_Rsp,
		Uid:resultInfo.Uuid,
	}
	go func(body utils.AsyncBody) {
		fmt.Println("async add body")
		utils.Queue.Add(body)
	}(body)
	return &out,nil
}


func Register(seq string){
	err := handle.Register(seq)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func Publish(service string){

	err := handle.Publish(service)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func Subscribe(service string){
	err := handle.Subscribe(service)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func Broadcast(body []byte,service string){
	//压缩1 不压缩0
	//加密1,2,3  不加密0
	infoBody , err := handle.MarshalBody(body,0,0)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	broadResult,err := handle.CallBroadcast(infoBody,service)
	if err != nil {
		fmt.Println(err.Error())
	}

	if broadResult.M_Err != nil {
		return
	}

	fmt.Println(broadResult)
}

func Sync(body []byte,compress , encrpyt int,uuid string){
	infoBody , err := handle.MarshalBody(body,compress,encrpyt)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	syncResult,err := handle.CallSync(infoBody,uuid)
	if err != nil {
		fmt.Println("C++ call Sync err:",err.Error())
		return
	}
	fmt.Println(syncResult)
}

func Async(body []byte,compress , encrpyt int,cb func([]byte)int){
	infoBody , err := handle.MarshalBody(body,compress , encrpyt)
	if err != nil {
		fmt.Println(err.Error())
	}
	uid := uuid.New()

	asyncResult,err := handle.CallAsync(infoBody,uid)
	if err != nil {
		fmt.Println("C++ call Async err:",err.Error())
	}
	fmt.Println(fmt.Sprintf("====async:%s",asyncResult))
	////异步回调
	go func(uid string,cb func([]byte)int) {
		for {
			select {
			case <- time.After(500*time.Microsecond):
				if  utils.BodyList.Check(uid) {
					body := utils.BodyList.Remove(uid)

					fmt.Println(fmt.Sprintf("异步回调:result body:%s",body.Body))
					result := cb(body.Body)
					if result == 0 {
						fmt.Println("call async  function error")
					}else{
						fmt.Println("call async function success")
					}
					return
				}
			}
		}
	}(uid,cb)
}

func main(){
	listener, err := net.Listen("tcp", Host+":"+Port)
	if err != nil {
		log.Fatalln("faile listen at: " + Host + ":" + Port)
	} else {
		log.Println("server is listening at: " + Host + ":" + Port)
	}
	rpcServer := grpc.NewServer()
	pb.RegisterClientServiceServer(rpcServer,&MsgHandle{})
	reflection.Register(rpcServer)
	if err = rpcServer.Serve(listener); err != nil {
		log.Fatalln("failed serve at: " + Host + ":" + Port)
	}
}


