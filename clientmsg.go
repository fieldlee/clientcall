package main
import "C"
import (
	"clientcall/call"
	"clientcall/handle"
	pb "clientcall/proto"
	"clientcall/utils"
	"context"
	"errors"
	"fmt"
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
	fmt.Println("info.Clientip:",string(info.Clientip))
	fmt.Println("info.Service:",string(info.Service))
	fmt.Println("info.Uuid:",string(info.Uuid))
	fmt.Println("info.M_Msg:",string(info.M_Body.M_Msg))
	out.M_Net_Rsp = []byte(fmt.Sprintf("call result for:%s",string(info.Clientip)))
	return &out,nil
}

func (m *MsgHandle)AsyncCall(ctx context.Context, resultInfo *pb.CallReqInfo) (*pb.CallRspInfo, error) {
	out := pb.CallRspInfo{}

	if resultInfo.Uuid == ""{
		return &out,errors.New("the Async Uuid is empty")
	}
	fmt.Println("info.Clientip:",string(resultInfo.Clientip))
	fmt.Println("info.Service:",string(resultInfo.Service))
	fmt.Println("info.Uuid:",string(resultInfo.Uuid))
	fmt.Println("info.M_Msg:",string(resultInfo.M_Body.M_Msg))
	out.M_Net_Rsp = []byte(fmt.Sprintf("async call result for:%s",string(resultInfo.Clientip)))
	go func() {
		<- time.After(time.Second * 2)
		Sync([]byte("answer info"),0,0,resultInfo.Uuid)
	}()
	return &out,nil
}

func (m *MsgHandle)AsyncAnswer(ctx context.Context, resultInfo *pb.CallReqInfo) (*pb.CallRspInfo, error) {
	out := pb.CallRspInfo{}
	resultInfo.Uuid = ""
	resultInfo.Service = ""

	fmt.Println("info.Clientip:",string(resultInfo.Clientip))
	fmt.Println("info.Service:",string(resultInfo.Service))
	fmt.Println("info.Uuid:",string(resultInfo.Uuid))
	fmt.Println("info.M_Msg:",string(resultInfo.M_Body.M_Msg))
	out.M_Net_Rsp = []byte(fmt.Sprintf("async answer result for:%s",string(resultInfo.Clientip)))

	return &out,nil
}


func Register(seq string){
	err := call.Register(seq)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func Publish(service string){

	err := call.Publish(service)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func Subscribe(service string){
	err := call.Subscribe(service)
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
	broadResult,err := call.CallBroadcast(infoBody,service)
	if err != nil {
		fmt.Println(err.Error())
	}

	if broadResult.M_Err != nil {
		return
	}

	fmt.Println(broadResult)
}

func B(){
	broadResult := &pb.NetRspInfo{
		M_Err:nil,
		M_Net_Rsp: map[uint32]*pb.SendResultInfo{
			12:&pb.SendResultInfo{
				Key:12,
				CheckErr:nil,
				Result:[]byte("hello struct"),
			},
			13:&pb.SendResultInfo{
				Key:13,
				CheckErr:nil,
				Result:[]byte("hello struct2"),
			},
		},
	}

	if broadResult.M_Err != nil {
		return
	}


}

func Sync(body []byte,compress , encrpyt int,uuid string){
	infoBody , err := handle.MarshalBody(body,compress,encrpyt)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(len(infoBody))
	fmt.Println(infoBody)


	syncResult,err := call.CallSync(infoBody,uuid)
	if err != nil {
		fmt.Println("C++ call Sync err:",err.Error())
		return
	}
	fmt.Println(syncResult)
}

func Async(body []byte,compress , encrpyt int){
	infoBody , err := handle.MarshalBody(body,compress , encrpyt)
	if err != nil {
		fmt.Println(err.Error())
	}
	asyncResult,err := call.CallAsync(infoBody)
	if err != nil {
		fmt.Println("C++ call Async err:",err.Error())
	}
	fmt.Println(asyncResult)
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


