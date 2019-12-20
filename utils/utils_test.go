package utils

import (
	"clientmsg/model"
	"fmt"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {

	i:=0
	for i < 1 {
		pr,pb:=GenerateKeyPair(4096)
		fmt.Println(PrivateKeyToBytes(pr))
		fmt.Println(PublicKeyToBytes(pb))
		i++
	}
}

func TestEncryptWithPublicKey(t *testing.T) {
	data := []byte("hello sfadsfassdfasdsfadsfassdfasdssfadsfadsfassdfasdfasdfasdfasdfdasfadsfworldfasdfasdfadsfadsfadsfasdfa的说法都是法师打发的算法第三方打到谁发的顺丰sdfasdfasdfasdfasdfdasfadsfworldfasdfasdfadsfadsfadsfasdfa的说法都是法师打发的算法第三方打到谁发的顺丰sdfasdfasdfasdfasdfdasfadsfworldfasdfasdfadsfadsfadsfasdfa的说法都是法师打发的算法第三方打到谁发的顺丰sdfasdfasdfasdfasdfdasfadsf")
	fmt.Println(len(data))
	//pri,pub:=GenerateKeyPair(4096)

	pubkey := BytesToPublicKey(model.PubKeyByte)
	body := EncryptWithPublicKey(data,pubkey)

	fmt.Println(body)

	prikey := BytesToPrivateKey(model.PriKeyByte)
	decode := DecryptWithPrivateKey(body,prikey)
	fmt.Println(decode)
	fmt.Println(string(decode))
}


func TestEncryptAes(t *testing.T) {
	data := []byte("hello sfadsfassdfasdsfadsfassdfasdssfadsfadsfassdfasdfasdfasdfasdfdasfadsfworldfasdfasdfadsfadsfadsfasdfa的说法都是法师打发的算法第三方打到谁发的顺丰sdfasdfasdfasdfasdfdasfadsfworldfasdfasdfadsfadsfadsfasdfa的说法都是法师打发的算法第三方打到谁发的顺丰sdfasdfasdfasdfasdfdasfadsfworldfasdfasdfadsfadsfadsfasdfa的说法都是法师打发的算法第三方打到谁发的顺丰sdfasdfasdfasdfasdfdasfadsf")
	fmt.Println(len(data))

	encryptBytes,err := EncryptAes(data,[]byte(model.PassPass16))
	if err != nil {
		fmt.Println(err)
	}

	decryptBytes,err := DecryptAes(encryptBytes,[]byte(model.PassPass16))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(decryptBytes))
}

func TestEncrypt3DES(t *testing.T) {
	data := []byte("hello sfadsfassdfasdsfadsfassdfasdssfadsfadsfassdfasdfasdfasdfasdfdasfadsfworldfasdfasdfadsfadsfadsfasdfa的说法都是法师打发的算法第三方打到谁发的顺丰sdfasdfasdfasdfasdfdasfadsfworldfasdfasdfadsfadsfadsfasdfa的说法都是法师打发的算法第三方打到谁发的顺丰sdfasdfasdfasdfasdfdasfadsfworldfasdfasdfadsfadsfadsfasdfa的说法都是法师打发的算法第三方打到谁发的顺丰sdfasdfasdfasdfasdfdasfadsf")
	fmt.Println(len(data))

	encryptBytes := Encrypt3DES(data,[]byte(model.PassPass24))

	decryptBytes := Decrypt3DES(encryptBytes,[]byte(model.PassPass24))

	fmt.Println(string(decryptBytes))
}