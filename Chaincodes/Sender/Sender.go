
package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"time"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"crypto/sha256"
	"github.com/open-quantum-safe/liboqs-go/oqs"
	"strconv"
)

type Sender struct{}


func AES_Encrypt(privData []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte(""), err
	}
	blockSize := block.BlockSize()
	privData = PKCS7Padding(privData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	encryptedDate := make([]byte, len(privData))
	blockMode.CryptBlocks(encryptedDate, privData)
	return encryptedDate, nil
}

func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func (t *Sender) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success([]byte("Success invoke and not opter!!"))
}

func (t *Sender) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fn, args := stub.GetFunctionAndParameters()
	if fn == "Gen_Sender_Data" {
		return t.Gen_Sender_Data(stub, args)
	} else if fn == "Get_Sender_Data" {
		return t.Get_Sender_Data(stub, args)
	}

	return shim.Error("Recevied unkown function invocation")
}

func (t *Sender) Gen_Sender_Data(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	t0 := time.Now()
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// 1.We get the private message and KS through the getTransient function of shim
	// and the type of message and KS is []byte.
	tMap, err := stub.GetTransient()
	if err != nil {
		return shim.Error(fmt.Sprintf("Could not retrieve transient, err %s", err))
	}
	messageBytes, ok := tMap["MESSAGE"]
	if !ok {
		return shim.Error(fmt.Sprintf("Expected transient message"))
	}
	KSBytes, ok := tMap["KS"]
	if !ok {
		return shim.Error(fmt.Sprintf("Expected transient KS"))
	}

	if len(KSBytes) >= 64 {
		hash := sha256.New()
		hash.Write(KSBytes)
		KSBytes = hash.Sum(nil)
	}
	
	//2.1 We use KS to encrypt message through the AES function.
	encryptedMessage, err:= AES_Encrypt(messageBytes, KSBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	//2.2 We use EncapSecret function to generate CT1,SS1,CT2,SS2.
	senderKeyid := args[0]
	EsPKKeyid1 := args[1]
	EsPKKeyid2 := args[2]
	kemName := args[3]
	queryArgs := [][]byte{[]byte("Get_EA_PubKey"), []byte(EsPKKeyid1)}
	response1 := stub.InvokeChaincode("Escrow", queryArgs, "mychannel")
	if response1.Status != shim.OK {
		return shim.Error(fmt.Sprintf("failed to query chaincode.got error :%s", response1.Payload))
	}
	queryArgs = [][]byte{[]byte("Get_EA_PubKey"), []byte(EsPKKeyid2)}
	response2 := stub.InvokeChaincode("Escrow", queryArgs, "mychannel")
	if response2.Status != shim.OK {
		return shim.Error(fmt.Sprintf("failed to query chaincode.got error :%s", response2.Payload))
	}

	t1 := time.Now()
	server := oqs.KeyEncapsulation{}
	defer server.Clean() 
	if err := server.Init(kemName, nil); err != nil {
		return shim.Error(err.Error())
	}
	//response.payload: publickey
	CT1, sharedSecret1, err := server.EncapSecret(response1.Payload)
	if err != nil {
		return shim.Error(err.Error())
	}
	CT2, sharedSecret2, err := server.EncapSecret(response2.Payload)
	if err != nil {
		return shim.Error(err.Error())
	}
	elapsed1 := time.Since(t1)
	encapSecretTime := strconv.FormatFloat(elapsed1.Seconds(), 'E', -1, 64)

	// XOR to get sharedSecret
	sharedSecret := make([]byte, len(sharedSecret1))
	for i:=0; i<len(sharedSecret); i++ {
		sharedSecret[i] = sharedSecret1[i] ^ sharedSecret2[i]
	}
	if len(sharedSecret) >= 64 {
		hash := sha256.New()
		hash.Write(sharedSecret)
		sharedSecret = hash.Sum(nil)
	}

	// 2.3 We use sharedSecret to encrypt message through the AES function.
	encrytedMessage2, err := AES_Encrypt(messageBytes, sharedSecret)
	if err != nil {
		return shim.Error(err.Error())
	}

	combineBytes := bytes.Join([][]byte{encryptedMessage, CT1, CT2, encrytedMessage2}, []byte("-----"))
	err = stub.PutState(senderKeyid, combineBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	elapsed := time.Since(t0)
	alltime := strconv.FormatFloat(elapsed.Seconds(), 'E', -1, 64)
	return shim.Success([]byte("Put encryptedMessages / cipherTexts  successfully!!!------ALLtime:" + alltime + "-------EncapSecretTime"+encapSecretTime))
}

func (t *Sender) Get_Sender_Data(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	keyid := args[0]
	result, err := stub.GetState(keyid)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(result)
}

func main() {
	err1 := shim.Start(new(Sender))
	if err1 != nil {
		fmt.Printf("error starting simple chaincode:%s \n", err1)
	}
}
