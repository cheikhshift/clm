package main

import (
	"flag"
	"sync"
	"net"
	"io"
	"log"
	"fmt"
	"github.com/cheikhshift/gos/core"
	"os"
	"io/ioutil"
)






// Data structures to manage 
// web server instances
type ServerHost struct {
	Lock  *sync.RWMutex
	Cache map[int]int
}

func NewCache() ServerHost {
	return ServerHost{Lock: new(sync.RWMutex), Cache: make(map[int]int)}
}


var (
	BinApp string 
	Limit int
	PostStart int = 60000
	Host ServerHost
	TargetIP string
	ServerWait int
	IpFormat string = "%s:%v"
	bspath = "./launcher.sh"
)

func GetServerAvailable() (string,int) {
	var index int
	Host.Lock.Lock()
	defer Host.Lock.Unlock()
	for index,concount := range Host.Cache {
		if concount < Limit {
			Host.Cache[index] += 1
			return fmt.Sprintf(IpFormat, TargetIP, (PostStart + index) ),index
		}
	}
	index = len(Host.Cache)
	
	//run bash
	targetport := PostStart + index

	os.Setenv("PORT", fmt.Sprintf("%v", targetport) )
	core.RunCmd(fmt.Sprintf("sh %s", bspath))
	Host.Cache[index] = 1
	return fmt.Sprintf(IpFormat, TargetIP, targetport ),index
}

func main() {
	serverbin := flag.String("app", "", "Path to binary of web server.")
	port := flag.String("port", "8080", "Port to listen on.")
	maxcon := flag.Int("max", 100, "Maximum number of connections per instance. clm will scale as needed.")
	iptocomp := flag.String("ip","127.0.0.1", "IPv4 address of machine.")
	waitperiod := flag.Int("wait",10, "Time to wait for server to start.")

	flag.Parse()

	TargetIP = *iptocomp
	BinApp = *serverbin
	Limit = *maxcon
	ServerWait = *waitperiod
			shscript := fmt.Sprintf(`#!/bin/bash  
cmd="%s"
eval "${cmd}" &>clm.log &disown
sleep %v
exit 0`, BinApp,ServerWait)
	ioutil.WriteFile(bspath, []byte(shscript), 0777)
	Host = NewCache()
	
	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", *port) )
	if err != nil {
		panic(err)
	}



	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	ipaddr, indx := GetServerAvailable()
	proxy, err := net.Dial("tcp",  ipaddr )
	if err != nil {
		Host.Lock.Lock()
		defer Host.Lock.Unlock()
		delete(Host.Cache, indx)
		handleRequest(conn)
		return
	}
	
	go copyIO(conn, proxy,-5)
	go copyIO(proxy, conn,indx)
}

func copyIO(src, dest net.Conn, index int) {
	defer src.Close()
	defer dest.Close()
	io.Copy(src, dest)
	if index != -5 {
		Host.Lock.Lock()
		defer Host.Lock.Unlock()
		Host.Cache[index] -=  1
		fmt.Println(Host)
	}


}



