// conf.go
package conf

import (
	"bufio"
	log "code.google.com/p/log4go"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type confInfo struct {
	sid            string
	token          string
	appId          string
	httpServ       string
	callRest       string
	customRest     string
	partnerId      string
	partnerKey     string
	number         string
	fee            float64
	callfee        float64
	dbtype         string
	dbname         string
	dbuser         string
	dbpwd          string
	dbaddr         string
	recordPath     string
	retranCount    int64 /*重传次数*/
	retranInterval int64 /*重传间隔*/
}

var globalConf confInfo

func GetRecordPath() string {
	return globalConf.recordPath
}

func GetNumber() string {
	return globalConf.number
}

func GetAppid() string {
	return globalConf.appId
}

func GetCallRestUrl() string {
	return globalConf.callRest
}

func GetHttpServ() string {
	return globalConf.httpServ
}

func GetSid() string {
	return globalConf.sid
}

func GetToken() string {
	return globalConf.token
}

func GetPartnerKey() string {
	return globalConf.partnerKey
}

func GetPartnerId() string {
	return globalConf.partnerId
}

func GetCustomRest() string {
	return globalConf.customRest
}

func GetFee() float64 {
	return globalConf.fee
}

func GetCallFee() float64 {
	return globalConf.callfee
}

func GetDbUser() string {
	return globalConf.dbuser
}

func GetDbUserPwd() string {
	return globalConf.dbpwd
}

func GetDbServer() string {
	return globalConf.dbaddr
}

func GetDbName() string {
	return globalConf.dbname
}

func GetDbType() string {
	return globalConf.dbtype
}

func GetRetranCount() int64 {
	return globalConf.retranCount
}

func GetRetranInterval() int64 {
	/*设置默认时间间隔*/
	if globalConf.retranInterval == 0 {
		globalConf.retranInterval = 3
	}
	return globalConf.retranInterval
}

func parseConf(conf string) {

	if 0 == strings.Index(conf, "#") || len(conf) == 0 {
		return
	}

	info := strings.Split(conf, "=")

	if len(info[1]) > 0 {

		if info[0] == "token" {
			log.Info("set :", info[0], info[1])
			globalConf.token = info[1]
		} else if info[0] == "sid" {
			log.Info("set :", info[0], info[1])
			globalConf.sid = info[1]
		} else if info[0] == "appid" {
			log.Info("set :", info[0], info[1])
			globalConf.appId = info[1]
		} else if info[0] == "callRest" {
			globalConf.callRest = info[1]
			log.Info("set :", info[0], info[1])
		} else if info[0] == "server" {
			globalConf.httpServ = info[1]
			log.Info("set :", info[0], info[1])
		} else if info[0] == "customRest" {
			globalConf.customRest = info[1]
			log.Info("set :", info[0], info[1])
		} else if info[0] == "partner_id" {
			globalConf.partnerId = info[1]
			log.Info("set :", info[0], info[1])
		} else if info[0] == "partner_key" {
			globalConf.partnerKey = info[1]
			log.Info("set :", info[0], info[1])
		} else if info[0] == "number" {
			globalConf.number = info[1]
			log.Info("set :", info[0], info[1])
		} else if info[0] == "callfee" {
			globalConf.callfee, _ = strconv.ParseFloat(info[1], 64)
			log.Info("set :", info[0], info[1])
		} else if info[0] == "fee" {
			globalConf.fee, _ = strconv.ParseFloat(info[1], 64)
			log.Info("set :", info[0], info[1])
		} else if info[0] == "dbtype" {
			globalConf.dbtype = info[1]
			log.Info("set :", info[0], info[1])
		} else if info[0] == "dbname" {
			globalConf.dbname = info[1]
			log.Info("set :", info[0], info[1])
		} else if info[0] == "dbuser" {
			globalConf.dbuser = info[1]
			log.Info("set :", info[0], info[1])
		} else if info[0] == "dbpwd" {
			globalConf.dbpwd = info[1]
			log.Info("set :", info[0], info[1])
		} else if info[0] == "dbaddr" {
			globalConf.dbaddr = info[1]
			log.Info("set :", info[0], info[1])
		} else if info[0] == "recordpath" {
			globalConf.recordPath = info[1]
			log.Info("set :", info[0], info[1])
		} else if info[0] == "retran_count" {
			globalConf.retranCount, _ = strconv.ParseInt(info[1], 10, 64)
			if globalConf.retranCount > 10 {
				globalConf.retranCount = 10
			}
			log.Info("set:", info[0], globalConf.retranCount)
		} else if info[0] == "retran_interval" {
			globalConf.retranInterval, _ = strconv.ParseInt(info[1], 10, 64)
			if globalConf.retranInterval > 30 {
				globalConf.retranInterval = 30
			}
			log.Info("set:", info[0], globalConf.retranInterval)
		}
	}
}

func LoadGlobalConf(fileName string) {

	f, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, ("配置文件(%s)不存在:%s"), fileName, err.Error())
		os.Exit(1)
	}

	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		parseConf(line)
		if err != nil {
			if err == io.EOF {
				return
			}
			return
		}
	}
	if len(globalConf.token) == 0 {
		fmt.Fprintf(os.Stderr, "token is not set")
		os.Exit(1)
	}

	if len(globalConf.sid) == 0 {
		fmt.Fprintf(os.Stderr, "sid is not set")
		os.Exit(1)
	}

	if len(globalConf.appId) == 0 {
		fmt.Fprintf(os.Stderr, "appId is not set")
		os.Exit(1)
	}

	if len(globalConf.httpServ) == 0 {
		fmt.Fprintf(os.Stderr, "http server is not set")
		os.Exit(1)
	}

	if len(globalConf.customRest) == 0 {
		fmt.Fprintf(os.Stderr, "customRest  is not set")
		os.Exit(1)
	}

	if len(globalConf.partnerId) == 0 {
		fmt.Fprintf(os.Stderr, "partnerId  is not set")
		os.Exit(1)
	}

	if len(globalConf.partnerKey) == 0 {
		fmt.Fprintf(os.Stderr, "partnerKey  is not set")
		os.Exit(1)
	}
}
