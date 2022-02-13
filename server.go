package main

import (
	"flag"
	"fmt"
	"log"
	"wagent/pkg/auth"
	"wagent/pkg/cmdb"
	"wagent/pkg/command"
	"wagent/pkg/health"
	"wagent/pkg/setting"
	"wagent/pkg/wcron"

	"github.com/smallnest/rpcx/server"
)

func init() {
	fmt.Printf("initing %s ...\n", "wagent")
	setting.Setup()
}

func main() {
	log.Println("starting init the wagent...")
	host := flag.String("host", setting.WagentSetting.Host, "wagent listen addr")
	port := flag.Int("port", setting.WagentSetting.Port, "wagent listen port")

	flag.Parse()
	addr := fmt.Sprintf("%s:%d", *host, *port)
	fmt.Printf("addr is %s\n", addr)

	// ensure sn_id exists
	h := cmdb.Host{}
	id := h.MustExistsID()
	if id == 0 {
		log.Fatal("cannot init sn_id of cmdb!")
	} else {
		log.Println("cmdb check successful with id ", id)
	}

	wcron.Run()

	s := server.NewServer()
	s.AuthFunc = auth.Auth
	//addRegistryPlugin(s)
	s.RegisterName("CmdExecutor", new(command.CmdExecutor), "")
	s.RegisterName("Health", new(health.Health), "")

	err := s.Serve("tcp", addr)
	if err != nil {
		panic(err)
	}
}
