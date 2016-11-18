// main.go
package main

import (
	callList "callManager"
	log "code.google.com/p/log4go"
	"conf"
	"fmt"
	"ljClient"
	"runtime"
	"server"
	"time"
)

const (
	version  = "V1.1.0.1-20160530"
	logConf  = "../conf/logconfig.xml"
	confPath = "../conf/app.conf"
)

func initLogger() {

	log.LoadConfiguration(logConf)
	log.Info("version : %s", version)

	log.Info("Current time is : %s", time.Now().Format("15:04:05 MST 2006/01/02"))

	return
}

func test() bool {
	var info ljClient.NumberResp
	err := ljClient.GetServiceNum("20160429222758646193-98e3ffaef6ddec12", "18898739887", "4009205045-8889", &info)
	if nil != err {
		fmt.Println("get service num err!", err)
	}

	fmt.Println("get info:", info)

	//ljClient.PostBill("20160429222758646193-98e3ffaef6ddec12")

	return true
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	conf.LoadGlobalConf(confPath)

	callList.CallInit()

	initLogger()

	//mysql.DbInit()

	fmt.Println("Hello World!")

	server.Start()
}
