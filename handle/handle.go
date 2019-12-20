package handle

import (
	"clientcall/model"
	pb "clientcall/proto"
	"clientcall/utils"
	"errors"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"math/rand"
	"time"
)

func randInt(min int , max int) uint32 {
	rand.Seed( time.Now().UTC().UnixNano())
	return uint32(min + rand.Intn(max-min))
}

func MarshalBody(body []byte,Compress,Encrypt int)([]byte,error){
	net_msgbody := &pb.Net_Pack_Min_Net_MsgBody{
		MLAsktype:uint64(10005),
		MLServerSequence:uint64(100),
		MLAskSequence:uint64(1014),
		MCMsgAckType:int32(2),
		MCMsgType:int32(2),
		MSSendCount:int32(1),
		MLExpireTime:uint32(50),
		MISendTimeApp:uint64(6),
		MLResult:int32(1),
		MLBack:uint64(1),
		MIDiscard:int32(0),
	}

	net_pack := &pb.Net_Pack{
		M_Msg:body,
		M_MsgBody:net_msgbody,
	}


	gj_net_pack := &pb.GJ_Net_Pack{
		M_Net_Pack: map[uint32]*pb.Net_Pack{randInt(10000,1000000):net_pack},
	}

	gjbody,err := proto.Marshal(gj_net_pack)

	if err != nil {
		return nil,err
	}

	if model.COMPRESS_TYPE(Compress) < model.Compression_no || model.COMPRESS_TYPE(Compress) >= model.CompressionWayMax{
		return nil,errors.New("compress way error")
	}

	if model.ENCRPTION_TYPE(Encrypt) < model.Encryption_No || model.ENCRPTION_TYPE(Encrypt) >= model.Encryption_Max{
		return nil,errors.New("encrypt way error")
	}else{
		if model.ENCRPTION_TYPE(Encrypt) == model.Encryption_RSA{
			if len(gjbody) > 440 {
				return nil,errors.New("body bytes length too long, cann't use rsa encrypt type")
			}
		}
	}

	fullbody := FullHead(gjbody,int(Compress),int(Encrypt))

	return fullbody,nil
}


func FullHead(inbody []byte,compress ,encryptType int)[]byte{
	headINfo := model.HeadInfo{
		Tag:model.HeadTag,
		Version:int16(model.HeadVersion),
		ClientType:int16(model.HeadClientType),
		HeadLength:int16(model.HeadLength),
		CompressWay:uint8(compress),
		Encryption:uint8(encryptType),
		Sig:uint8(model.HeadSig),
		Format:uint8(model.HeadFormat),
		NetFlag:uint8(model.HeadNetFlag),
		Back1:0,
		BufSize:int32(len(inbody)),
		UncompressedSize:int32(len(inbody)),
		Back2:0,
	}
	encryptByte := inbody
	//////加密
	switch model.ENCRPTION_TYPE(headINfo.Encryption) {
	case model.Encryption_AES:
		encryptByte,_ = utils.EncryptAes(inbody,[]byte(model.PassPass16))
	case model.Encryption_RSA:
		pubkey := utils.BytesToPublicKey(model.PubKeyByte)
		encryptByte = utils.EncryptWithPublicKey(inbody,pubkey)
	case model.Encryption_Des:
		encryptByte = utils.Encrypt3DES(inbody,[]byte(model.PassPass24))
	}
	///加密后的长度有所改变
	inbody = encryptByte
	headINfo.UncompressedSize = int32(len(inbody))
	/////压缩body bytes
	if model.COMPRESS_TYPE(headINfo.CompressWay) == model.Compression_zip {
		if zipbody,err := utils.ZipByte(inbody);err != nil {
			fmt.Println(err.Error())
		}else{
			inbody = zipbody
		}
	}

	/////修改压缩后的buffer长度
	headINfo.BufSize = int32(len(inbody))

	inbody = utils.JoinHeadAndBody(headINfo,inbody)

	return inbody
}
