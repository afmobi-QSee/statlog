package statlog

import "testing"

func Test_BufferStringJoin(t *testing.T){
	statusUdp := &StatusUdp{StatEnv:"testEnv"}
	result := statusUdp.BufferStringJoin("test")
	if result != "testEnvtest" {
		t.Error("bufferstringJoin test fail")
	}
}
