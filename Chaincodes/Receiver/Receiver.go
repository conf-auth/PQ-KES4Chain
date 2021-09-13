
package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"strconv"
	"fmt"
	"time"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"crypto/sha256"
)

type Receiver struct{}

func AES_Decrypt(encryptedDate []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte(""), err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	privDate := make([]byte, len(encryptedDate))
	blockMode.CryptBlocks(privDate, encryptedDate)
	privDate = PKCS7UnPadding(privDate)
	return privDate, nil
}


func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}


func (t *Receiver) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success([]byte("Success invoke and not opter!!"))
}

func (t *Receiver) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fn, args := stub.GetFunctionAndParameters()
	if fn == "Dec_Sec_Data" {
		return t.Dec_Sec_Data(stub, args)
	}
	return shim.Error("Recevied unkown function invocation")
}

func (t *Receiver) Dec_Sec_Data(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	t0 := time.Now()
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	tMap, err := stub.GetTransient()
	if err != nil {
		return shim.Error(fmt.Sprintf("Could not retrieve transient, err %s", err))
	}
	KSBytes, ok := tMap["KS"]
	if !ok {
		return shim.Error(fmt.Sprintf("Expected transient KS"))
	}

	senderKeyid := args[0]
	queryArgs := [][]byte{[]byte("Get_Sender_Data"), []byte(senderKeyid)}
	response := stub.InvokeChaincode("Sender", queryArgs, "mychannel")
	if response.Status != shim.OK {
		return shim.Error(fmt.Sprintf("failed to query chaincode.got error :%s", response.Payload))
	}
	enConBytes := bytes.Split(response.Payload, []byte("-----"))
	
	if len(KSBytes) >= 64 {
		hash := sha256.New()
		hash.Write(KSBytes)
		KSBytes = hash.Sum(nil)
	}
	messageBytes, err := AES_Decrypt(enConBytes[0], KSBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	elapsed := time.Since(t0)
	runtime := strconv.FormatFloat(elapsed.Seconds(), 'E', -1, 64)
	return shim.Success([]byte("The receiver successfully receive the message!!!---message:" + string(messageBytes)+"----runtime:"+ runtime))
}

func main() {
	err1 := shim.Start(new(Receiver))
	if err1 != nil {
		fmt.Printf("error starting simple chaincode:%s \n", err1)
	}
}
