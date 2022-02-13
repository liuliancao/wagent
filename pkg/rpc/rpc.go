package rpc

import (
	"context"
	"fmt"
	"wagent/pkg/setting"

	"github.com/smallnest/rpcx/client"
)

type HttpArgs struct {
	Method string
	Data   string
	Url    string
}

type HttpReply struct {
	Result string
	Error  string
	Code   int
}

type HealthArgs struct {
	BasePath string
}

type HealthReply struct {
	Status uint
	Error  string
}

func CMDBRpc(method string, data string, url string) (result string, error string) {
	addr := fmt.Sprintf("%s:%d", setting.GuarderSetting.Host, setting.GuarderSetting.Port)
	dis, _ := client.NewPeer2PeerDiscovery("tcp@"+addr, "")

	xclient := client.NewXClient("HttpForwarder", client.Failtry, client.RandomSelect, dis, client.DefaultOption)

	defer xclient.Close()

	httpArgs := &HttpArgs{
		Method: method,
		Data:   data,
		Url:    url,
	}
	httpReply := &HttpReply{}

	err := xclient.Call(context.Background(), "Http"+method, httpArgs, httpReply)
	if err != nil {
		fmt.Println("call Http ", httpArgs, "error:", err.Error())
	}
	return httpReply.Result, httpReply.Error
}

func HealthRpc(addr string) (status int, err error) {
	// 0 down; 1 up
	dis, _ := client.NewPeer2PeerDiscovery("tcp@"+addr, "")

	xclient := client.NewXClient("Health", client.Failtry, client.RandomSelect, dis, client.DefaultOption)

	defer xclient.Close()

	healthArgs := &HealthArgs{}
	healthReply := &HealthReply{}
	err = xclient.Call(context.Background(), "Check", healthArgs, healthReply)
	if err != nil {
		return 0, err
	}
	return 1, nil
}
