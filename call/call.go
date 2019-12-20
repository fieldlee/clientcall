package call

import (
	pb "clientcall/proto"
	"clientcall/utils"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
)

func CallSync(mbody []byte,uuid string)(*pb.NetRspInfo,error){
	caddr := fmt.Sprintf("%v:%v",utils.ServerAddr,utils.ServerPort)
	conn, err := grpc.Dial(caddr, grpc.WithInsecure())
	if err != nil {
		return nil,err
	}
	defer conn.Close()
	c := pb.NewMidServiceClient(conn)

	ctx := context.Background()

	rsp, err := c.Sync(ctx,&pb.NetReqInfo{M_Body:mbody,Service:"",Uuid:uuid})
	//////////////////////异步处理，调用客户端的接口，异步发送
	if err != nil {
		return nil,err
	}
	return rsp,nil
}

func CallAsync(mbody []byte)(*pb.NetRspInfo,error){
	caddr := fmt.Sprintf("%v:%v",utils.ServerAddr,utils.ServerPort)
	conn, err := grpc.Dial(caddr, grpc.WithInsecure())
	if err != nil {
		return nil,err
	}
	defer conn.Close()
	c := pb.NewMidServiceClient(conn)
	ctx := context.Background()

	rsp, err := c.Async(ctx,&pb.NetReqInfo{M_Body:mbody,Service:""})

	//////////////////////异步处理 ， 调用客户端的接口，异步发送
	if err != nil {
		return nil,err
	}
	//if rsp.M_Err != nil {
	//	return nil,errors.New(string(rsp.M_Err))
	//}

	return rsp,nil
}

func CallBroadcast(mbody []byte,service string)(*pb.NetRspInfo,error){
	caddr := fmt.Sprintf("%v:%v",utils.ServerAddr,utils.ServerPort)
	conn, err := grpc.Dial(caddr, grpc.WithInsecure())
	if err != nil {
		return nil,err
	}
	defer conn.Close()
	c := pb.NewMidServiceClient(conn)
	ctx := context.Background()
	rsp, err := c.Broadcast(ctx,&pb.NetReqInfo{M_Body:mbody,Service:service})
	if err != nil {
		return nil,err
	}
	//if rsp.M_Err != nil {
	//	return nil,errors.New(string(rsp.M_Err))
	//}
	return rsp,nil
}

func Register(sequence string)error{
	caddr := fmt.Sprintf("%v:%v",utils.ServerAddr,utils.ServerPort)

	conn, err := grpc.Dial(caddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewMidServiceClient(conn)
	var ctx context.Context
	ctx = context.Background()
	registerReturn , err := c.Register(ctx,&pb.RegisterInfo{
		Sequence:sequence,
	})
	if err != nil {
		return err
	}
	if registerReturn.Success != true {
		return errors.New(fmt.Sprintf("Register error"))
	}
	return nil
}

func Publish(service string)error{
	caddr := fmt.Sprintf("%v:%v",utils.ServerAddr,utils.ServerPort)
	conn, err := grpc.Dial(caddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewMidServiceClient(conn)
	var ctx context.Context
	ctx = context.Background()
	publishReturn , err :=c.Publish(ctx,&pb.PublishInfo{
		Service:service,
	})
	if err != nil {
		return err
	}
	if publishReturn.Success != true {
		return errors.New(fmt.Sprintf("Publish error"))
	}
	return nil
}

func Subscribe(service string)error{
	caddr := fmt.Sprintf("%v:%v",utils.ServerAddr,utils.ServerPort)
	conn, err := grpc.Dial(caddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewMidServiceClient(conn)
	var ctx context.Context
	ctx = context.Background()
	subscribeInfo,err := c.Subscribe(ctx,&pb.SubscribeInfo{
		Service:service,
	})
	if err != nil {
		return err
	}
	if subscribeInfo.Success != true {
		return errors.New(fmt.Sprintf("Subscribe error"))
	}
	return nil
}