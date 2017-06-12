package statlog

import (
	"github.com/ibbd-dev/go-async-log"
	"os"
	"time"
	"strconv"
	"encoding/json"
)

var StatLog *asyncLog.LogFile

func StatLogInit(statLogPath, statLogName, statEnv, statUdpHost, localIp, localPort string) error{
	statusUdp := StatUdpInit(statEnv, statUdpHost, localIp, localPort)

	if statusUdp.err != nil {
		return statusUdp.err
	}

	_, err := os.Stat(statLogPath)

	if os.IsNotExist(err){
		os.MkdirAll(statLogPath, 0755)
	}

	StatLog = asyncLog.NewLevelLog(statLogPath + statLogName, asyncLog.LevelOff)
	StatLog.SetFlags(asyncLog.NoFlag)
	StatLog.SetRotate(asyncLog.RotateDate)

	return nil
}

func Stat(modelName string, model interface{}) error{
	modelJson, err := json.Marshal(model)
	if err != nil {
		return err
	}
	StatLog.Write(strconv.FormatInt(time.Now().Unix(), 10) + "	" + modelName + "	" + string(modelJson))
	return nil
}
