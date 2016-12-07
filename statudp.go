package statlog

import (
	"net"
	"strings"
	"bytes"
)

type StatusUdp struct {
	StatEnv string
	StatUdpHost string
	StatusUdpConn net.Conn
	LocalIp string
	err error
}

var statusUdpStruct *StatusUdp

func StatUdpInit(statEnv, statUdpHost string) *StatusUdp{
	StatusUdpConn, err := net.Dial("udp", statUdpHost)
	if err != nil {
		return &StatusUdp{err:err}
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return &StatusUdp{err:err}
	}
	statusUdp := new(StatusUdp)
	statusUdp.StatusUdpConn = StatusUdpConn
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				statusUdp.LocalIp = strings.Replace(ipnet.IP.String(), ".", "-", -1)
				break
			}

		}
	}
	statusUdp.StatEnv = statEnv
	statusUdpStruct = statusUdp
	return statusUdp
}

func (this *StatusUdp)sentUdp(data string){
	this.StatusUdpConn.Write([]byte(data))
}

//key + set去重统计 比如所有用户去重统计user + uid
func AccessSetByIp(key, set string){
	statusUdpStruct.sentUdp(statusUdpStruct.BufferStringJoin(".", key, ":", set, "|s"))
}

//根据本机IP + key + set去重统计 比如本机的用户去重统计 ip + user + uid
func AccessSetByIP(key, set string) {
	data := statusUdpStruct.BufferStringJoin(".", statusUdpStruct.LocalIp, ".", key, ":", set, "|s")
	statusUdpStruct.sentUdp(data)
}

//根据api + key + set去重统计 比如api + user + uid
func ApiSet(apiName, key, set string) {
	data := statusUdpStruct.BufferStringJoin(".", apiName, ".", key, ":", set, "|s")
	statusUdpStruct.sentUdp(data)
}

func ApiSetByIP(apiName, key, set string) {
	data := statusUdpStruct.BufferStringJoin(".", statusUdpStruct.LocalIp, ".", apiName, ".", key, ":", set, "|s")
	statusUdpStruct.sentUdp(data)
}

//多条统计 总访问量  api访问量 和 api执行时间统计 毫秒
func MultCount(apiName, ms string) {
	var buffer bytes.Buffer
	buffer.WriteString(statusUdpStruct.BufferStringJoin(":1|c\n"))  //所有访问量统计
	buffer.WriteString(statusUdpStruct.BufferStringJoin(":", ms, "|ms\n"))  //所有访问量 + 执行时间统计
	buffer.WriteString(statusUdpStruct.BufferStringJoin(".", statusUdpStruct.LocalIp, ":1|c\n")) //根据本机ip统计
	buffer.WriteString(statusUdpStruct.BufferStringJoin(".", statusUdpStruct.LocalIp, ":", ms, "|ms\n"))  //根据本机ip + 执行时间统计
	buffer.WriteString(statusUdpStruct.BufferStringJoin(".", apiName, ":1|c\n"))  //根据api接口统计
	buffer.WriteString(statusUdpStruct.BufferStringJoin(".", apiName, ":", ms, "|ms\n"))  //根据api接口 + 执行时间统计
	buffer.WriteString(statusUdpStruct.BufferStringJoin(".", statusUdpStruct.LocalIp, ".", apiName, ":1|c\n"))  //根据本机ip + api接口统计
	buffer.WriteString(statusUdpStruct.BufferStringJoin(".", statusUdpStruct.LocalIp, ".", apiName, ":", ms, "|ms")) //根据本机ip + api接口 + 执行时间统计
	statusUdpStruct.sentUdp(buffer.String())
}