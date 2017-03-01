package statlog

import (
	"strings"
	"bytes"
	"gopkg.in/alexcesaro/statsd.v2"
)

type StatusUdp struct {
	StatEnv string
	StatusUdpConn	*statsd.Client
	LocalIpAndPort string
	err error
}

var statusUdpStruct *StatusUdp

func StatUdpInit(statEnv, statUdpHost, localIp, localPort string) *StatusUdp{
	StatusUdpConn, err := statsd.New(statsd.Address(statUdpHost))
	if err != nil {
		return &StatusUdp{err:err}
	}

	//获取本机ip地址
	statusUdp := new(StatusUdp)
	statusUdp.StatusUdpConn = StatusUdpConn

	var buffer bytes.Buffer
	ipFormat := strings.Replace(localIp, ".", "-", -1)
	buffer.WriteString(ipFormat)
	buffer.WriteString("-")
	buffer.WriteString(localPort)
	statusUdp.LocalIpAndPort = buffer.String()

	statusUdp.StatEnv = statEnv
	statusUdpStruct = statusUdp
	return statusUdp
}

func (this *StatusUdp)sentUdp(data string){
	this.StatusUdpConn.Increment(data)
}

//key + set去重统计 比如所有用户去重统计user + uid
func AccessSet(key, set string){
	statusUdpStruct.sentUdp(statusUdpStruct.BufferStringJoin(".", key, ":", set, "|s"))
}

//根据本机IP + key + set去重统计 比如本机的用户去重统计 ip + user + uid
func AccessSetByIP(key, set string) {
	data := statusUdpStruct.BufferStringJoin(".", statusUdpStruct.LocalIpAndPort, ".", key, ":", set, "|s")
	statusUdpStruct.sentUdp(data)
}

//根据api + key + set去重统计 比如api + user + uid
func ApiSet(apiName, key, set string) {
	data := statusUdpStruct.BufferStringJoin(".", apiName, ".", key, ":", set, "|s")
	statusUdpStruct.sentUdp(data)
}

//根据本机IP + api + key + set去重统计 比如本机 + api + user + uid
func ApiSetByIP(apiName, key, set string) {
	data := statusUdpStruct.BufferStringJoin(".", statusUdpStruct.LocalIpAndPort, ".", apiName, ".", key, ":", set, "|s")
	statusUdpStruct.sentUdp(data)
}

//多条统计 总访问量  api访问量 和 api执行时间统计 毫秒
func MultCount(apiName, ms string) {
	var buffer bytes.Buffer
	buffer.WriteString(statusUdpStruct.BufferStringJoin(":1|c\n"))  //所有访问量统计
	buffer.WriteString(statusUdpStruct.BufferStringJoin(":", ms, "|ms\n"))  //所有访问量 + 执行时间统计
	buffer.WriteString(statusUdpStruct.BufferStringJoin(".", statusUdpStruct.LocalIpAndPort, ":1|c\n")) //根据本机ip统计
	buffer.WriteString(statusUdpStruct.BufferStringJoin(".", statusUdpStruct.LocalIpAndPort, ":", ms, "|ms\n"))  //根据本机ip + 执行时间统计
	buffer.WriteString(statusUdpStruct.BufferStringJoin(".", apiName, ":1|c\n"))  //根据api接口统计
	buffer.WriteString(statusUdpStruct.BufferStringJoin(".", apiName, ":", ms, "|ms\n"))  //根据api接口 + 执行时间统计
	buffer.WriteString(statusUdpStruct.BufferStringJoin(".", statusUdpStruct.LocalIpAndPort, ".", apiName, ":1|c\n"))  //根据本机ip + api接口统计
	buffer.WriteString(statusUdpStruct.BufferStringJoin(".", statusUdpStruct.LocalIpAndPort, ".", apiName, ":", ms, "|ms")) //根据本机ip + api接口 + 执行时间统计
	statusUdpStruct.sentUdp(buffer.String())
}
