package statlog

import "testing"

func Test_Error_StatUdpInit(t *testing.T){
	statusUdp := StatUdpInit("testEnv", "testHost")
	if statusUdp.err == nil {
		t.Error("testErrorHost fail")
	}
}

func Test_StatUdpInit(t *testing.T){
	statusUdp := StatUdpInit("prd.palmgroup.access", "172.17.40.21:8125")
	if statusUdp.err != nil {
		t.Error("statusInit fail", statusUdp.err)
	}
}

func Test_StatusInit(t *testing.T){
	StatUdpInit("prd.palmgroup.access", "172.17.40.21:8125")
	AccessSetByIp("key", "set")
}