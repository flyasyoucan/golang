// sqlClient project sqlClient.go
package sqlClient

import (
	Call "callManager"
	log "code.google.com/p/log4go"
	"conf"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var dbHandle *sql.DB
var fileHandle *sql.DB

const (
	DB_BILL_NAME_STR  = "sid,app_sid,partner_call_id,answer_time,bill_duration,callee_num,callee_show_num,caller_num,caller_show_num,cost,call_duration,end_time,sound_url,result,start_time"
	DB_BILL_VALUE_STR = "?,?,?,?,?,?,?,?,?,?,?,?,?,?,?"
)

func DbInit() bool {

	dbinfo := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8", conf.GetDbUser(), conf.GetDbUserPwd(), conf.GetDbServer())
	db, err := sql.Open(conf.GetDbType(), dbinfo)
	//defer db.Close()

	if err != nil {
		log.Error("Open database error: %s\n", err)
		return false
	}

	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(100)

	err = db.Ping()
	if err != nil {
		log.Error("connect sql:", err)
		return false
	} else {
		log.Debug("connect sql success!", conf.GetDbServer())
	}

	dbHandle = db

	return true
}

func LianjiaFileDbInit() bool {

	dbinfo := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8", "ucp_ipcc", "sh25p2yFlAQdJdBU", "10.10.89.134:3307")

	db, err := sql.Open(conf.GetDbType(), dbinfo)
	if err != nil {
		log.Error("Open voice file database error: %s\n", err)
		return false
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(50)

	err = db.Ping()
	if err != nil {
		log.Error("connect voice file sql:", err)
		return false
	}

	fmt.Println("file handle init success!", dbinfo)

	fileHandle = db
	return true
}

func GetMainMenu(accessNumber string, appid string) (string, error) {

	var (
		voice string
	)

	if nil == dbHandle {
		log.Error("db client is not ready!")
		return "", errors.New("DB closed")
	}

	log.Debug("select accese number:", accessNumber)

	err := dbHandle.QueryRow("select fileName from ipcc_customer.tb_ipcc_lianjia_number_voice where number = ?", accessNumber).Scan(&voice)
	if err != nil {
		log.Error("get voice from sql failed:", err)
		return "", errors.New("selectErr")
	}

	if "" == voice {
		return "", errors.New("notSet")
	}

	return voice, nil
}

func GetIvrVoiceMenu(accessNumber string, appid string) (string, error) {

	var (
		number400 string
		voice     string
	)

	if nil == dbHandle {
		log.Error("db client is not ready!")
		return "", errors.New("DB closed")
	}

	log.Debug("select accese number:", accessNumber)

	err := dbHandle.QueryRow("select number from ipcc_customer.tb_ipcc_lianjia_number_appid where nbr = ?", accessNumber).Scan(&number400)
	if err != nil {
		log.Error("get 400 numbers from sql failed:", err)
		return "", errors.New("selectErr")
	}

	err = dbHandle.QueryRow("select voice from ipcc_customer.tb_ipcc_lianjia_number_voice where number = ?", number400).Scan(&voice)
	if err != nil {
		log.Error("get voice from sql failed:", err)
		return "", errors.New("selectErr")
	}

	return voice, nil
}

func GetIvrVoiceName(fileId string) string {

	var filePath string

	if nil == fileHandle {
		log.Error("db client is not ready!")

		return ""
	}

	log.Debug("select file:", fileId)

	err := fileHandle.QueryRow("SELECT remote_path FROM ucpaas.tb_ucpaas_ivr_ring WHERE id = ?", fileId).Scan(&filePath)
	if err != nil {
		log.Error("get voice path from sql failed:", err)
		return ""
	}

	return filePath
}

func Log2sql(callid string, eventType string, content string) bool {

	if nil == dbHandle {
		log.Error("db client is not ready!")
		return false
	}

	stmt, err := dbHandle.Prepare("INSERT INTO ipcc_customer.tb_ipcc_lianjia_log(callId,eventType,content) VALUES(?,?,?)")

	if err != nil {
		log.Error("prepare failed:", err)
		return false
	}

	defer stmt.Close()

	stmt.Exec(callid, eventType, content)
	if err != nil {
		log.Error("excute sql failed:", err)
		return false
	}

	return true
}

func Bill2Sql(call *Call.CallInfo) bool {

	if nil == dbHandle {
		log.Error("db client is not ready!")
		return false
	}
	//just for test
	dbPrepareStr := fmt.Sprintf("INSERT INTO %s.tb_ipcc_test_lianjia_bill_log(%s) VALUES(%s)", conf.GetDbName(), DB_BILL_NAME_STR, DB_BILL_VALUE_STR)
	//dbPrepareStr := fmt.Sprintf("INSERT INTO %s.tb_ipcc_lianjia_bill_log(%s) VALUES(%s)", conf.GetDbName(), DB_BILL_NAME_STR, DB_BILL_VALUE_STR)
	stmt, err := dbHandle.Prepare(dbPrepareStr)

	if err != nil {
		log.Error("prepare failed:", err)
		return false
	}

	defer stmt.Close()
	_, err = stmt.Exec(conf.GetSid(), conf.GetAppid(), call.GetCallId(), call.GetAnswerTime(),
		call.GetBillDuration(), call.GetCallee(), call.GetCalleeHideNum(), call.GetCaller(),
		call.GetCallerHideNum(), call.GetCost(), call.GetDuration(), call.GetEndTime(),
		call.GetRecord(), call.GetResult(), call.GetStartTime())
	if err != nil {
		log.Error("excute sql failed:", err)
		return false
	}

	log.Debug("write to sql success")

	return true
}

/* 链家项目 获取接入号对应的400号 */
func GetbindNumber(appid string, numbers map[string]string) error {

	var (
		nbr    string
		number string
	)

	if nil == dbHandle {
		log.Error("db client is not ready!")
		return errors.New("DB closed")
	}

	sqlRow, err := dbHandle.Query("select number,nbr from ipcc_customer.tb_ipcc_lianjia_number_appid")
	if err != nil {
		log.Error("get 400 numbers from sql failed:", err)
		return errors.New("selectErr")
	}

	log.Debug("get number:", nbr, "-", sqlRow)
	defer sqlRow.Close()

	for sqlRow.Next() {
		err := sqlRow.Scan(&number, &nbr)
		if err != nil {
			log.Error(err)
			continue
		}

		log.Debug("get number:", nbr, number)
		numbers[nbr] = number

	}

	return nil
}

/* 国信项目，从数据库中查找录音文件所在目录 */
func GetRecordDsid(callid string) string {

	var dsid string

	if nil == dbHandle {
		log.Error("db client is not ready!")
		return ""
	}

	log.Debug("select file:", callid)

	sqlcmd := fmt.Sprintf("SELECT dsid FROM %s.tb_ucpaas_bill_log WHERE call_id = ?", conf.GetDbName())

	err := dbHandle.QueryRow(sqlcmd, callid).Scan(&dsid)
	if err != nil {
		log.Error("get voice path from sql failed:", err)
		return ""
	}

	return dsid
}

func dbUnInit() {
	dbHandle.Close()
	fileHandle.Close()
}
