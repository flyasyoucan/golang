// callManager project callManager.go
package callManager

import (
	log "code.google.com/p/log4go"
	"errors"
	"strings"
	"time"
)

const (
	TimeType = "2006-01-02 15:04:05"
)

type CallInfo struct {
	callStatus    int
	fee           int /* 费率 */
	callId        string
	caller        string
	nbr           string /*呼入接入号*/
	callerHideNum string
	callee        string
	calleeHideNum string
	callStart     string /* 呼叫坐席时刻 YYYY-MM-DD HH:MM:SS */
	answerTime    string /* 坐席接通时刻，YYYY-MM-DD HH:MM:SS。若通话不成功，此字段值应传0000-00-00 00:00:00 */
	endTime       string /* 通话结束时刻，格式同上*/
	record        string /* 录音文件 */
	duration      string /* 通话时长 */
	billDuration  string /* 计费时长 */
	cost          string
	result        string /* 通话结果 */
}

type callJson struct {
	CallStatus    int    `json:"CallStatus"`
	CallId        string `json:"CallId"`
	Caller        string `json:"Caller"`
	CallerHideNum string `json:"CallerHideNum"`
	Callee        string `json:"Callee"`
	CalleeHideNum string `json:"CalleeHideNum"`
	CallStart     string `json:"CallStart"`
	AnswerTime    string `json:"AnswerTime"`
	EndTime       string `json:"EndTime"`
	Record        string `json:"Record"`
	Duration      string `json:"Duration"`
	BillDuration  string `json:"BillDuration"`
	Cost          string `json:"Cost"`
	Result        string `json:"Result"`
}

func (p *CallInfo) GetCallId() string {
	return p.callId
}

func (p *CallInfo) GetCaller() string {
	return p.caller
}

func (p *CallInfo) GetRecord() string {
	return p.record
}

func (p *CallInfo) GetResult() string {
	return p.result
}

func (p *CallInfo) GetCost() string {
	return p.cost
}

func (p *CallInfo) GetDuration() string {
	return p.duration
}

func (p *CallInfo) GetBillDuration() string {
	return p.billDuration
}

func (p *CallInfo) GetCallee() string {
	return p.callee
}

func (p *CallInfo) GetCalleeHideNum() string {
	return p.calleeHideNum
}

func (p *CallInfo) GetCallerHideNum() string {
	return p.callerHideNum
}

func (p *CallInfo) GetStartTime() string {
	return p.callStart
}

func (p *CallInfo) GetAnswerTime() string {
	return p.answerTime
}

func (p *CallInfo) GetEndTime() string {
	return p.endTime
}

func (p *CallInfo) GetNbr() string {
	return p.nbr
}

var CallList map[string]CallInfo

func GetCalleeHideNumber(callid string) string {
	if val, ok := CallList[callid]; ok {
		return val.GetCalleeHideNum()
	} else {
		return ""
	}
}

func InNewCall(callId string, caller string, callee string, nbr string) {
	var newCall CallInfo

	//删除前缀带023的手机号码的前缀
	if strings.HasPrefix(caller, "023177") {
		log.Error("wrong caller number:", caller)
		caller = caller[3:]
	}

	newCall.callId = callId
	newCall.caller = caller
	newCall.callee = callee
	newCall.nbr = nbr
	CallList[callId] = newCall
}

func UpdateCallee(callid string, callee string, calleeHide string, callerHide string) {

	log.Debug("update callee show number:", calleeHide)

	if val, ok := CallList[callid]; ok {
		val.callee = callee
		val.calleeHideNum = calleeHide
		val.callerHideNum = callerHide
		CallList[callid] = val
	} else {
		log.Error("Can not find call:", callid)
	}
}

func UpdateCallStartTime(callid string) {
	if val, ok := CallList[callid]; ok {
		val.callStart = time.Now().Format(TimeType)
		CallList[callid] = val
	} else {
		log.Error("Can not find call:", callid)
	}
}

func UpdateCallEndTime(callid string, record string) {

	if val, ok := CallList[callid]; ok {
		val.endTime = time.Now().Format(TimeType)
		if len(record) > 0 {
			val.record = record
		}
		val.callStatus++
		CallList[callid] = val
	} else {
		log.Error("Can not find call:", callid)
	}
}

func UpdateCallResult(callid string, result string) {

	log.Debug("UpdateCallResult ,update result:", callid, result)
	if val, ok := CallList[callid]; ok {

		if len(val.result) > 0 {
			return
		}

		val.result = result
		CallList[callid] = val

	} else {
		log.Error("Can not find call:", callid)
	}
}

func UpdateCallAnswerTime(callid string, result string) {
	if val, ok := CallList[callid]; ok {
		val.answerTime = time.Now().Format(TimeType)
		val.callStatus++
		val.result = result
		CallList[callid] = val
	} else {
		log.Error("Can not find call:", callid)
	}
}

func UpdateCallBill(callId string, callTime string, totalTime string, cost string) {

	if val, ok := CallList[callId]; ok {
		val.duration = callTime
		val.billDuration = totalTime
		val.cost = cost
		val.callStatus++
		CallList[callId] = val
	} else {
		log.Error("Can not find call:", callId)
	}
}

func FindCall(callid string) (CallInfo, error) {

	var val CallInfo
	var ok bool

	if val, ok = CallList[callid]; ok {
		return val, nil
	} else {
		return val, errors.New("cannot find the call")
	}
}

func DelCall(callid string) {
	delete(CallList, callid)
}

func CallInit() {
	CallList = make(map[string]CallInfo)
	testData()
}

func testData() {
	var testCall CallInfo
	//YYYY-MM-DD HH:MM:SS
	testCall.answerTime = "2016-05-03 16:11:12"
	testCall.billDuration = "10"
	testCall.callee = "18589033693"
	testCall.calleeHideNum = "18888888888"
	testCall.caller = "18898739887"
	testCall.callerHideNum = "18898739888"
	testCall.callId = "20160429222758646193-98e3ffaef6ddec12"
	testCall.callStart = "2016-05-03 16:11:10"

	testCall.cost = "0.56"
	testCall.duration = "15"
	testCall.endTime = "2016-05-03 16:11:22"
	testCall.record = "http://"
	testCall.result = "ANSWERED"

	CallList["20160429222758646193-98e3ffaef6ddec12"] = testCall
}
