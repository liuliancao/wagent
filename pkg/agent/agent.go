package main

import (
	"context"
	"log"

	"github.com/smallnest/rpcx/client"
)

type CmdArgs struct {
	RunAs    string
	Cmd      string
	Async    bool
	Timeout  int
	Callback string
}

type CmdReply struct {
	Result string
	Error  string
	Code   int
}

func main() {
	d, _ := client.NewPeer2PeerDiscovery("tcp@"+"localhost:3210", "")
	//d, _ := etcd_client.NewEtcdDiscovery("/wagent", "Arith", []string{"192.168.10.204:2379"}, nil)

	xclient := client.NewXClient("CmdExecutor", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	xclient.Auth("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IjY5MmQzOTQ2YzcyOTI1OTJjZmMzM2Y4NmUyNzA1ZDhjIiwicGFzc3dvcmQiOiI2OTJkMzk0NmM3MjkyNTkyY2ZjMzNmODZlMjcwNWQ4YyIsImV4cCI6MTYyNzQ5NTQ3OCwiaXNzIjoiZ2luLWJsb2cifQ.lHHMOCyb3Mns2xBRFm3p8zqS-bbuMHSsP6uTjcPSY-s")
	defer xclient.Close()
	cmdArgs := &CmdArgs{
		RunAs:   "root",
		Cmd:     "echo hellogolang",
		Async:   false,
		Timeout: 5,
	}
	cmdReply := &CmdReply{}
	err := xclient.Call(context.Background(), "Execute", cmdArgs, cmdReply)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}
	log.Println(cmdReply)
	return
	cmdArgs = &CmdArgs{
		RunAs:    "root",
		Cmd:      "ls /home/liuliancao -l &>/dev/null && echo 'hah' || echo 'bad'",
		Async:    true,
		Timeout:  2,
		Callback: "http://www.baidu.com",
	}
	cmdReply1 := &CmdReply{}
	err = xclient.Call(context.Background(), "Execute", cmdArgs, cmdReply1)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}
	log.Printf("%s executed with user %s and async %v is\n\r %s\n\r error is %s\n\r", cmdArgs.Cmd, cmdArgs.RunAs, cmdArgs.Async, cmdReply1.Result, cmdReply1.Error)

}
