package test

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nbnet/side-chain/core/types"
	tmClient "github.com/tendermint/tendermint/rpc/client/http"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"
)

var (
	LargeTxDataPath string
	TxHash          string
)

func init() {
	flag.StringVar(&LargeTxDataPath, "ltdp", "", "description of large tx data path")
	flag.StringVar(&TxHash, "th", "", "tx hash")
}

func TestTx(t *testing.T) {
	mintTx := types.Tx{
		Ty:        types.Mint,
		Signature: "0xf9be5d1ae521c1688f7260bb3a9b725b776c818139ce777458c91b6ae85bddbe6428797a6ecc8b3416e87ff8bcc5340f1b37c341efc8edaa15d621ade973c1cb1b",
		Body: types.MintBody{

			Nonce:   0,
			Amount:  "0x1",
			Address: "0x9F8C645f2D0b2159767Bd6E0839DE4BE49e823DE",
		},
	}

	j, err := json.Marshal(&mintTx)
	if err != nil {
		panic(err)
	}
	println("json: ", string(j))

	var tx types.Tx
	err = json.Unmarshal(j, &tx)
	fmt.Printf("%+v\n", tx)

	blobTx := types.Tx{
		Ty:        types.Blob,
		Signature: "0xf9be5d1ae521c1688f7260bb3a9b725b776c818139ce777458c91b6ae85bddbe6428797a6ecc8b3416e87ff8bcc5340f1b37c341efc8edaa15d621ade973c1cb1b",
		Body: types.BlobBody{

			Data: "0xf9be5d1ae521c1688f7",
		},
	}

	j, err = json.Marshal(&blobTx)
	if err != nil {
		panic(err)
	}
	println("json: ", string(j))

	var blob types.Tx
	err = json.Unmarshal(j, &blob)
	fmt.Printf("%+v\n", blob)
}

func TestSendTx(t *testing.T) {
	url := "http://127.0.0.1:26657"
	blobTx := types.Tx{
		Ty:        types.Blob,
		Signature: "0xf9be5d1ae521c1688f7260bb3a9b725b776c818139ce777458c91b6ae85bddbe6428797a6ecc8b3416e87ff8bcc5340f1b37c341efc8edaa15d621ade973c1cb1b",
		Body: types.BlobBody{

			Data: "0xf9be5d1ae521c1688f7",
		},
	}

	j, err := json.Marshal(&blobTx)
	if err != nil {
		panic(err)
	}

	tx := base64.StdEncoding.EncodeToString(j)

	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "broadcast_tx_sync",
		"params": map[string]interface{}{
			"tx": tx,
		},
		"id": 1,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// {"id":1,"jsonrpc":"2.0","method":"broadcast_tx_sync","params":{"tx":"eyJzaWduYXR1cmUiOiIweGY5YmU1ZDFhZTUyMWMxNjg4ZjcyNjBiYjNhOWI3MjViNzc2YzgxODEzOWNlNzc3NDU4YzkxYjZhZTg1YmRkYmU2NDI4Nzk3YTZlY2M4YjM0MTZlODdmZjhiY2M1MzQwZjFiMzdjMzQxZWZjOGVkYWExNWQ2MjFhZGU5NzNjMWNiMWIiLCJib2R5Ijp7InR5cGUiOiJibG9iIiwiZGF0YSI6IjB4ZjliZTVkMWFlNTIxYzE2ODhmNyJ9fQ=="}}
	fmt.Println("Req:", string(jsonData))

	// 发送 POST 请求
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return
	}
	defer response.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// 打印响应状态和内容
	fmt.Println("Response Status:", response.Status)

	// {"jsonrpc":"2.0","id":1,"result":{"code":0,"data":"","log":"","codespace":"","hash":"84D2D0412EF0F270133F91D87AC86A624DE85BB2B1DCD7A0800333CC851D7791"}}
	fmt.Println("Response Body:", string(body))
}

func TestSendLargeTx(t *testing.T) {

	t.Log("flag value: ", LargeTxDataPath)

	b, err := os.ReadFile(LargeTxDataPath)
	if err != nil {
		panic(err)
	}

	t.Log("data len: ", len(b))

	url := "http://127.0.0.1:26657"
	blobTx := types.Tx{
		Ty:        types.Blob,
		Signature: "0xf9be5d1ae521c1688f7260bb3a9b725b776c818139ce777458c91b6ae85bddbe6428797a6ecc8b3416e87ff8bcc5340f1b37c341efc8edaa15d621ade973c1cb1b",
		Body: types.BlobBody{
			Data:    string(b),
			Address: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		},
	}

	j, err := json.Marshal(&blobTx)
	if err != nil {
		panic(err)
	}

	tx := base64.StdEncoding.EncodeToString(j)

	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "broadcast_tx_sync",
		"params": map[string]interface{}{
			"tx": tx,
		},
		"id": 1,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// 发送 POST 请求
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return
	}
	defer response.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// 打印响应状态和内容
	fmt.Println("Response Status:", response.Status)

	// {"jsonrpc":"2.0","id":1,"result":{"code":0,"data":"","log":"","codespace":"","hash":"84D2D0412EF0F270133F91D87AC86A624DE85BB2B1DCD7A0800333CC851D7791"}}
	fmt.Println("Response Body:", string(body))
}

func TestGzipSendLargeTx(t *testing.T) {

	t.Log("flag value: ", LargeTxDataPath)

	b, err := os.ReadFile(LargeTxDataPath)
	if err != nil {
		panic(err)
	}

	t.Log("data len: ", len(b))
	t.Log("data mb : ", len(b)/(1024*1024))

	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	if _, err := writer.Write(b); err != nil {
		panic(err)
	}

	compressData := hex.EncodeToString(buf.Bytes())
	compressDataLen := len(compressData)
	t.Log("compress len: ", compressDataLen)
	t.Log("compress mb : ", compressDataLen/(1024*1024))

	url := "http://127.0.0.1:26657"
	blobTx := types.Tx{
		Ty:        types.Blob,
		Signature: "",
		Body: types.BlobBody{
			Data:    compressData,
			Address: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		},
	}

	j, err := json.Marshal(&blobTx)
	if err != nil {
		panic(err)
	}

	tx := base64.StdEncoding.EncodeToString(j)

	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "broadcast_tx_sync",
		"params": map[string]interface{}{
			"tx": tx,
		},
		"id": 1,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// 发送 POST 请求
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return
	}
	defer response.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// 打印响应状态和内容
	fmt.Println("Response Status:", response.Status)

	// {"jsonrpc":"2.0","id":1,"result":{"code":0,"data":"","log":"","codespace":"","hash":"84D2D0412EF0F270133F91D87AC86A624DE85BB2B1DCD7A0800333CC851D7791"}}
	fmt.Println("Response Body:", string(body))
}

func TestCheckTx(t *testing.T) {
	t.Log("flag value: ", LargeTxDataPath)

	b, err := os.ReadFile(LargeTxDataPath)
	if err != nil {
		panic(err)
	}

	t.Log("data len: ", len(b))

	url := "http://127.0.0.1:26657"
	blobTx := types.Tx{
		Ty:        types.Blob,
		Signature: "0xf9be5d1ae521c1688f7260bb3a9b725b776c818139ce777458c91b6ae85bddbe6428797a6ecc8b3416e87ff8bcc5340f1b37c341efc8edaa15d621ade973c1cb1b",
		Body: types.BlobBody{
			Data:    string(b),
			Address: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		},
	}

	j, err := json.Marshal(&blobTx)
	if err != nil {
		panic(err)
	}

	tx := base64.StdEncoding.EncodeToString(j)

	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "check_tx",
		"params": map[string]interface{}{
			"tx": tx,
		},
		"id": 1,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// 发送 POST 请求
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return
	}
	defer response.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// 打印响应状态和内容
	fmt.Println("Response Status:", response.Status)

	// {"jsonrpc":"2.0","id":1,"result":{"code":0,"data":"","log":"","codespace":"","hash":"84D2D0412EF0F270133F91D87AC86A624DE85BB2B1DCD7A0800333CC851D7791"}}
	fmt.Println("Response Body:", string(body))
}

func TestRandBlobTx(t *testing.T) {
	maxSize := 104857600/2 - 1000 // 100MB

	data := make([]byte, maxSize)
	_, err := rand.Read(data)
	if err != nil {
		panic(err)
	}

	hexData := hex.EncodeToString(data)
	hexDataLen := len(hexData)
	t.Log("hex data len: ", hexDataLen)
	t.Log("hex data mb: ", hexDataLen/(1024*1024))
	if hexDataLen > maxSize*2 {
		t.Log("data len too large ")
		return
	}

	url := "http://127.0.0.1:26657"
	blobTx := types.Tx{
		Ty:        types.Blob,
		Signature: "",
		Body: types.BlobBody{
			Data:    hexData,
			Address: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		},
	}

	j, err := json.Marshal(&blobTx)
	if err != nil {
		panic(err)
	}

	tx := base64.StdEncoding.EncodeToString(j)

	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "broadcast_tx_sync",
		"params": map[string]interface{}{
			"tx": tx,
		},
		"id": 1,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// 发送 POST 请求
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return
	}
	defer response.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// 打印响应状态和内容
	fmt.Println("Response Status:", response.Status)

	// {"jsonrpc":"2.0","id":1,"result":{"code":0,"data":"","log":"","codespace":"","hash":"84D2D0412EF0F270133F91D87AC86A624DE85BB2B1DCD7A0800333CC851D7791"}}
	fmt.Println("Response Body:", string(body))
}

func TestGzipRandBlobTx(t *testing.T) {

	rand.Seed(time.Now().UnixNano())

	//min := 104857600/2 - 1000 // 100MB
	//max := 209715200/2 - 1000 // 300MB

	max := 104857600/2 - 1000 // 50MB
	min := 54857600/2 - 1000  // 50MB

	t.Log("max: ", max)
	t.Log("min: ", min)

	randomNum := rand.Intn(max-min) + min
	t.Log("random num:   ", randomNum)

	data := make([]byte, randomNum)
	_, err := rand.Read(data)
	if err != nil {
		panic(err)
	}

	hexData := hex.EncodeToString(data)
	hexDataLen := len(hexData)
	t.Log("hex data len: ", hexDataLen)
	t.Log("hex data mb: ", hexDataLen/(1024*1024))

	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	if _, err := writer.Write([]byte(hexData)); err != nil {
		panic(err)
	}

	compressData := hex.EncodeToString(buf.Bytes())
	compressDataLen := len(compressData)
	t.Log("compress len: ", compressDataLen)
	t.Log("compress mb : ", compressDataLen/(1024*1024))

	url := "http://127.0.0.1:26657"
	blobTx := types.Tx{
		Ty:        types.Blob,
		Signature: "",
		Body: types.BlobBody{
			Data:    compressData,
			Address: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		},
	}

	j, err := json.Marshal(&blobTx)
	if err != nil {
		panic(err)
	}

	tx := base64.StdEncoding.EncodeToString(j)

	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "broadcast_tx_sync",
		"params": map[string]interface{}{
			"tx": tx,
		},
		"id": 1,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// 发送 POST 请求
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return
	}
	defer response.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// 打印响应状态和内容
	fmt.Println("Response Status:", response.Status)

	// {"jsonrpc":"2.0","id":1,"result":{"code":0,"data":"","log":"","codespace":"","hash":"84D2D0412EF0F270133F91D87AC86A624DE85BB2B1DCD7A0800333CC851D7791"}}
	fmt.Println("Response Body:", string(body))
}

func TestRpcBlobTx(t *testing.T) {
	rpcUrl := "http://127.0.0.1:7074/blob"

	t.Log("flag value: ", LargeTxDataPath)

	data, err := os.ReadFile(LargeTxDataPath)
	if err != nil {
		panic(err)
	}

	dataMd5 := md5.Sum(data)
	t.Log("data hex md5: ", hex.EncodeToString(dataMd5[:]))

	req := types.BlobBody{
		Data:    hex.EncodeToString(data),
		Address: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// 发送 POST 请求
	response, err := http.Post(rpcUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return
	}
	defer response.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// 打印响应状态和内容
	fmt.Println("Response Status:", response.Status)

	// {"jsonrpc":"2.0","id":1,"result":{"code":0,"data":"","log":"","codespace":"","hash":"84D2D0412EF0F270133F91D87AC86A624DE85BB2B1DCD7A0800333CC851D7791"}}
	fmt.Println("Response Body:", string(body))

}

func TestHashDataToFile(t *testing.T) {

	url := "http://34.210.245.20:26657"
	toFilePath := "./compress2.txt"

	t.Log("flag value: ", TxHash)

	txHash, err := hex.DecodeString(TxHash)
	if err != nil {
		panic(err)
	}

	tdClient, err := tmClient.New(url, "/websocket")
	if err != nil {
		panic(err)
	}

	tx, err := tdClient.Tx(context.Background(), txHash, false)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(toFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	_, err = file.Write(tx.Tx)
	if err != nil {
		panic(err)
	}

}

func Test1(t *testing.T) {
	data := []byte("I[2024-10-30|18:06:19.325] service stop                                 module=p2p peer=95a31e2a033a9ed8a2101e9f38722e8eec0af142@127.0.0.1:47492 msg=\"Stopping MConnection service\" impl=MConn{127.0.0.1:47492}")

	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	if _, err := writer.Write(data); err != nil {
		panic(err)
	}
	writer.Close()

	fmt.Println("data: ", hex.EncodeToString(buf.Bytes()))
}

func Test2(t *testing.T) {
	dataPath := "/home/cloud/vscode/data"

	data, err := os.ReadFile(dataPath)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	if _, err := writer.Write(data); err != nil {
		panic(err)
	}
	writer.Close()

	err = os.WriteFile("./data.gz", buf.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
	//return

	compressedData := buf.Bytes()
	reader, err := gzip.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	var decompressedData bytes.Buffer
	if _, err := io.Copy(&decompressedData, reader); err != nil {
		panic(err)
	}

	err = os.WriteFile("./data2", decompressedData.Bytes(), 0644)
	if err != nil {
		panic(err)
	}

	dataMd5 := md5.Sum(data)
	fmt.Println("data md5: ", hex.EncodeToString(dataMd5[:]))
	data2Md5 := md5.Sum(decompressedData.Bytes())
	fmt.Println("data2 md5: ", hex.EncodeToString(data2Md5[:]))
}

func Test3(t *testing.T) {
	data, err := os.ReadFile("/home/cloud/yunyc12345/side-chain/compress.txt")
	if err != nil {
		panic(err)
	}

	dataMd5 := md5.Sum(data)
	t.Log("data hex md5: ", hex.EncodeToString(dataMd5[:]))

	data1, err := os.ReadFile("/home/cloud/yunyc12345/side-chain/core/types/test/compress.txt")
	if err != nil {
		panic(err)
	}

	dataMd51 := md5.Sum(data1)
	t.Log("data hex md5: ", hex.EncodeToString(dataMd51[:]))

}
