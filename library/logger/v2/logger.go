package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"myself_framwork/utils"

	"github.com/google/uuid"
)

const SEND_TIMEOUT = "10s"

type LoggingInterface interface {
	AddField(string, interface{}) LogField
	Info(string, ...LogField)
	Warn(string, ...LogField)
	Error(string, ...LogField)
	Critical(string, ...LogField)
}

type Logging struct {
	SpanID string
}

type LogObj struct {
	InstitutionID string  `json:"institutionID" bson:"institutionID"` //"PEDE" [from caller]
	ServiceName   string  `json:"serviceName" bson:"serviceName"`     //"DIFI", [from caller]
	LogType       string  `json:"logType" bson:"logType"`             //"Transaction Log", [from caller]
	Data          LogData `json:"data" bson:"data"`
}

type LogData struct {
	UniqueId       string `json:"uniqueId" bson:"uniqueId"`
	PublishTime    string `json:"publishTime" bson:"publishTime"`       //"2006-01-02 15:04:05", [from library]
	LogLevel       string `json:"logLevel" bson:"logLevel"`             //"Error", [from library - depends which logger function caller hits]
	TraceID        string `json:"traceID" bson:"traceID"`               //from library for now
	ActionTo       string `json:"actionTo" bson:"actionTo"`             //"service_payment", [from caller]
	ActionName     string `json:"actionName" bson:"actionName"`         //"payment confirmation", [from caller]
	EndPoint       string `json:"endPoint" bson:"endPoint"`             //"http://endpoint.id", [from caller]
	ErrorDesc      string `json:"errorDesc" bson:"errorDesc"`           //"write tcp 172.26.3.2:14198->13.228.23.160:8432: write: broken pipe",  [from caller]
	FileName       string `json:"fileName" bson:"fileName"`             //"service.go-- line : 11", [from library]
	AdditionalData string `json:"additionalData" bson:"additionalData"` //JSON
	RequestBody    string `json:"requestBody" bson:"requestBody"`       //JSON mask pin, phonenumber, name,
	RequestHeader  string `json:"requestHeader" bson:"requestHeader"`   //JSON
	ResponseBody   string `json:"responseBody" bson:"responseBody"`     //JSON
	ResponseHeader string `json:"responseHeader" bson:"responseHeader"` //JSON
	ResponseId     string `json:"responseId" bson:"responseId"`         //JSON
}

type LogField struct {
	Key   string
	Value interface{}
}

var (
	APIURL string
)

const (
	LOGGER_API_URL          = "LOGGER_API_URL"
	LOGGER_API_URL_CONSTANT = "0.0.0.0:4000"
)

func init() {
	APIURL = utils.GetEnv(LOGGER_API_URL, LOGGER_API_URL_CONSTANT)
}

func InitLog() LoggingInterface {
	uid := fmt.Sprintf("%v", uuid.New())
	span := strings.ReplaceAll(uid, "-", "")

	return &Logging{
		SpanID: span,
	}
}

func InitLogs(req *http.Request) LoggingInterface {
	uid := fmt.Sprintf("%v", uuid.New())
	span := strings.ReplaceAll(uid, "-", "")

	return &Logging{SpanID: span}
}

func (log *Logging) GetSpanID() string {
	return log.SpanID
}

func (log *Logging) AddField(fieldname string, value interface{}) LogField {
	return LogField{Key: fieldname, Value: value}
}

func (log *Logging) Info(message string, fields ...LogField) {
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	runtime := getCaller()
	msg := make(chan []byte)
	defer close(msg)

	go log.processMessage("info", message, runtime, timeNow, fields, msg)
	go log.sendAndPrintMessage(<-msg)
}

func (log *Logging) Warn(message string, fields ...LogField) {
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	runtime := getCaller()
	msg := make(chan []byte)
	defer close(msg)

	go log.processMessage("warning", message, runtime, timeNow, fields, msg)
	go log.sendAndPrintMessage(<-msg)

}

func (log *Logging) Error(message string, fields ...LogField) {
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	runtime := getCaller()
	msg := make(chan []byte)
	defer close(msg)

	go log.processMessage("error", message, runtime, timeNow, fields, msg)
	go log.sendAndPrintMessage(<-msg)
}

func (log *Logging) Critical(message string, fields ...LogField) {
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	runtime := getCaller()
	msg := make(chan []byte)
	defer close(msg)

	go log.processMessage("critical", message, runtime, timeNow, fields, msg)
	go log.sendAndPrintMessage(<-msg)
}

func (l *Logging) processMessage(level string, message string, caller string, ts string, fields []LogField, msg chan []byte) {
	logObj := LogObj{
		Data: LogData{
			UniqueId:    l.SpanID,
			PublishTime: ts,
			LogLevel:    level,
			FileName:    caller,
		},
	}

	//Transformation of Log Key into LogObj struct
	for _, field := range fields {
		switch field.Key {
		case "InstitutionID", "Institutionid", "InstitutionId", "institutionID", "institutionId":
			logObj.InstitutionID = fmt.Sprintf("%+v", field.Value)
		case "ServiceName", "serviceName", "servicename", "Servicename":
			logObj.ServiceName = fmt.Sprintf("%+v", field.Value)
		case "LogType", "Logtype", "logType", "logtype":
			logObj.LogType = fmt.Sprintf("%+v", field.Value)
		case "PublishTime", "publishTime", "Publishtime", "publishtime":
			logObj.Data.PublishTime = fmt.Sprintf("%+v", field.Value)
		case "LogLevel", "Loglevel", "logLevel", "loglevel":
			logObj.Data.LogLevel = fmt.Sprintf("%+v", field.Value)
		case "TraceID", "traceID", "TraceId", "traceId":
			ns, _ := time.Parse("2006-01-02 15:04:05", ts)
			logObj.Data.TraceID = fmt.Sprintf("%+v", ns.Nanosecond())
		case "ActionTo", "Actionto", "actionTo", "actionao":
			logObj.Data.ActionTo = fmt.Sprintf("%+v", field.Value)
		case "ActionName", "actionName", "Actionname", "actionname":
			logObj.Data.ActionName = fmt.Sprintf("%+v", field.Value)
		case "EndPoint", "endPoint", "Endpoint", "endpoint":
			logObj.Data.EndPoint = fmt.Sprintf("%+v", field.Value)
		case "ErrorDesc", "errorDesc", "Errordesc", "errordesc":
			logObj.Data.ErrorDesc = fmt.Sprintf("%+v", field.Value)
		case "FileName", "fileName", "Filename", "filename":
			logObj.Data.FileName = fmt.Sprintf("%+v", field.Value)
		//Masking Required in case below
		case "AdditionalData", "additionalData", "Additionaldata", "additionaldata":
			switch field.Value.(type) {
			case []byte:
				byteData := field.Value.([]byte)
				logObj.Data.AdditionalData = string(byteData)
			case string:
				data := field.Value.(string)
				logObj.Data.AdditionalData = data
			default:
				if field.Value != nil {
					jsonByte, _ := json.Marshal(field.Value)
					logObj.Data.AdditionalData = string(jsonByte)
				} else {
					logObj.Data.AdditionalData = ""
				}
			}
		case "RequestBody", "requestBody", "Requestbody", "requestbody":
			//Mask "pin", "phonenumber", "nama" by regex the string(byte)
			keys := []string{"\"pin\":", "\"nama\":", "\"phonenumber\":"}
			switch field.Value.(type) {
			case []byte:
				byteData := field.Value.([]byte)
				logObj.Data.RequestBody = maskValues(string(byteData), keys)
			case string:
				data := field.Value.(string)
				logObj.Data.RequestBody = maskValues(data, keys)
			default:
				if field.Value != nil {
					jsonByte, _ := json.Marshal(field.Value)
					logObj.Data.RequestBody = maskValues(string(jsonByte), keys)
				} else {
					logObj.Data.RequestBody = ""
				}
			}
		case "RequestHeader", "requestHeader", "Requestheader", "requestheader":
			switch field.Value.(type) {
			case []byte:
				byteData := field.Value.([]byte)
				logObj.Data.RequestHeader = string(byteData)
			case string:
				data := field.Value.(string)
				logObj.Data.RequestHeader = data
			default:
				if field.Value != nil {
					jsonByte, _ := json.Marshal(field.Value)
					logObj.Data.RequestHeader = string(jsonByte)
				} else {
					logObj.Data.RequestHeader = ""
				}
			}
		case "ResponseBody", "responseBody", "Responsebody", "responsebody":
			switch field.Value.(type) {
			case []byte:
				byteData := field.Value.([]byte)
				logObj.Data.ResponseBody = string(byteData)
			case string:
				data := field.Value.(string)
				logObj.Data.ResponseBody = data
			default:
				if field.Value != nil {
					jsonByte, _ := json.Marshal(field.Value)
					logObj.Data.ResponseBody = string(jsonByte)
				} else {
					logObj.Data.ResponseBody = ""
				}
			}
		case "ResponseHeader", "Responseheader", "responseHeader", "responseheader":
			switch field.Value.(type) {
			case []byte:
				byteData := field.Value.([]byte)
				logObj.Data.ResponseHeader = string(byteData)
			case string:
				data := field.Value.(string)
				logObj.Data.ResponseHeader = data
			default:
				if field.Value != nil {
					jsonByte, _ := json.Marshal(field.Value)
					logObj.Data.ResponseHeader = string(jsonByte)
				} else {
					logObj.Data.ResponseHeader = ""
				}
			}
		case "ResponseId", "responseID", "ResponseID", "responseId":
			switch field.Value.(type) {
			case []byte:
				byteData := field.Value.([]byte)
				logObj.Data.ResponseId = string(byteData)
			case string:
				data := field.Value.(string)
				logObj.Data.ResponseId = data
			default:
				if field.Value != nil {
					jsonByte, _ := json.Marshal(field.Value)
					logObj.Data.ResponseId = string(jsonByte)
				} else {
					logObj.Data.ResponseId = ""
				}
			}
		}
	}

	//Publish Log to File
	go l.publishToFile(logObj)

	//Marshal log to publish in kafka
	jsonMap, err := json.Marshal(logObj)
	if err != nil {
		log.Printf("Error Marshal. Caller:%s, level:%s, Error:%+v\n", caller, level, err)
	}
	msg <- jsonMap
}

func getCaller() string {
	//format "service.go-- line : 11", [from library]
	_, filePath, line, ok := runtime.Caller(2)
	if !ok {
		return "file un-recoverable"
	}
	files := strings.Split(filePath, "/")
	return fmt.Sprintf("%s/%s-- line : %d", files[len(files)-2], files[len(files)-1], line)
}

//publishToFile Publish obj into txt file
func (l *Logging) publishToFile(obj LogObj) {
	//Create, Check path
	dir := "./logFiles"
	err := os.MkdirAll(dir, 0755) //0755 = chmod permision
	if err != nil {
		log.Println("MkdirAll Err:", err)
		return
	}

	//Check Filename which based from month & date
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Println("LoadLocation Err:", err)
		return
	}
	fileName := time.Now().In(location).Format("Jan2")
	flag := os.O_APPEND | os.O_CREATE | os.O_WRONLY               //flag = append, create if not found, write chmod file
	file, err := os.OpenFile(dir+"/"+fileName+".log", flag, 0644) //0644 = chmod permision
	if err != nil {
		log.Println("OpenFile Err:", err)
		return
	}
	defer file.Close()

	//Write to file
	msg := fmt.Sprintf("%s:\t%+v\n", obj.Data.PublishTime, obj)
	_, err = file.Write([]byte(msg))
	if err != nil {
		log.Println("Write Err:", err)
		return
	}
}

//sendAndPrintMessage Publish to Kafka, If failed, add to MongoDB or log to file
func (log *Logging) sendAndPrintMessage(msg []byte) {
	if len(msg) < 1 {
		//handle if msg is nil (err Marshalling)
		return
	}
	logObj := LogObj{}
	json.Unmarshal(msg, &logObj)
	log.sendToAPI(logObj)
}

//sendToAPI Send to API --> publish to kafka or mongoDB
func (l *Logging) sendToAPI(logObj interface{}) {
	//Publish to mongoDB
}

//maskValues mask certain values from param and scrambles it
func maskValues(strData string, keys []string) string {
	for _, key := range keys {
		re, err := regexp.Compile(`(?i)` + key + `\W?("[a-zA-Z0-9]+")`)
		if err != nil {
			return string(strData)
		}

		//Find index of masking key
		index := re.FindStringSubmatchIndex(strData) //should be in index[2] as start index[3] as end
		if len(index) < 1 {                          //doesn't match
			continue
		}

		//take the value from that key
		value := strData[index[2]:index[3]]

		//scramble the value
		scrambleStr := ScrambleString(value)

		//patch strData so that ..., "key": scrambled value, ...
		strData = strData[:index[2]] + "\"" + scrambleStr + "\"" + strData[index[3]:]
	}

	return strData
}

//ScrambleString change mask value to xxxx
func ScrambleString(str string) string {
	return "xxxx"
}
