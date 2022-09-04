package jsonrpc2

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/websocket"
	"reflect"
	"time"
)

type tokens struct {
	createTime int64
}

var tokensMap = make(map[string]tokens)

type wsMapStruct struct {
	token string
}

var wsMap = make(map[*websocket.Conn]wsMapStruct)

type stubMapping map[string]interface{}

var stubStorage = stubMapping{}

type RpcData struct {
	Id     int                    `json:"id"`
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
	Ws     *websocket.Conn
}

func AddToken(token string) {
	tokensMap[token] = tokens{createTime: time.Now().Unix()}
}

func Registry(name string, fn interface{}) {
	stubStorage[name] = fn
	fmt.Println(stubStorage)
}

func rpcHandler(ws *websocket.Conn, reply string) {
	var rpc RpcData
	rpc.Ws = ws
	err := json.Unmarshal([]byte(reply), &rpc)
	if err != nil {
		return
	}
	// function, ok := stubStorage[rpc.Method]
	fmt.Println("method: ", rpc.Method)
	if _, ok := stubStorage[rpc.Method]; ok {
		err := call(rpc.Method, rpc)
		if err != nil {
			errCode := fmt.Sprintf(
				`{"jsonrpc": "2.0", "error": {"code": -32000, "message": "Server error"}, "id": %d}`,
				rpc.Id,
			)
			_, _ = ws.Write([]byte(errCode))
		}
	} else {
		_, _ = ws.Write([]byte(fmt.Sprintf(
			`{"jsonrpc": "2.0", "error": {"code": -32601, "message": "Method %s not found"}, "id": %d}`, rpc.Method, rpc.Id,
		)))
	}

}

func call(funcName string, params ...interface{}) (err error) {
	f := reflect.ValueOf(stubStorage[funcName])
	if len(params) != f.Type().NumIn() {
		err = errors.New("the number of params is out of index")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	// var res []reflect.Value
	_ = f.Call(in)
	// result = res[0].Interface()
	return
}

func Notify(rpc RpcData, method string, params string) {
	ret := []byte(fmt.Sprintf(`{"jsonrpc": "2.0", "method": "%s", "params": %s}`, method, params))
	_, _ = rpc.Ws.Write(ret)
}

func Return(rpc RpcData, data []byte) {
	_, _ = rpc.Ws.Write(
		[]byte(
			fmt.Sprintf(`{"jsonrpc": "2.0", "result": %s, "id": %d}`, string(data), rpc.Id),
		),
	)

}

func Serve(ws *websocket.Conn) {
	//add ws to map
	wsMap[ws] = wsMapStruct{token: "123123"}
	fmt.Println(wsMap)
	var firstCall bool = true
	for {
		var reply string
		if err := websocket.Message.Receive(ws, &reply); err != nil {
			break
		}
		if firstCall {
			if _, ok := tokensMap[reply]; ok {
				//delete record
				delete(tokensMap, reply)
				firstCall = false
			} else {
				break
			}
		} else {
			go rpcHandler(ws, reply)
		}

	}
	//del ws from map
	delete(wsMap, ws)
	fmt.Println(wsMap)
}
