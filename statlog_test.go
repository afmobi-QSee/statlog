package statlog

import (
	"testing"
	"os"
)

type model struct {
	model1 string
	model2 string
}

func Test_StatLogInit(t *testing.T){
	err := StatLogInit("statLogPath", "statLogName", "prd.palmgroup.access", "172.17.40.21:8125", "127.0.0.1", "8080")
	if err!=nil {
		t.Error("statLogInit fail", err)
	}
	os.RemoveAll("statLogPath")
}

func Test_Stat(t *testing.T){
	StatLogInit("statLogPath", "statLogName", "prd.palmgroup.access", "172.17.40.21:8125", "127.0.0.1", "8080")
	err := Stat("modelName", &model{model1:"1",  model2:"2"})
	if err != nil {
		t.Error("stat test fail", err)
	}
	os.RemoveAll("statLogPath")
}