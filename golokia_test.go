package golokia

import (
	"testing"
)

var service = "http://localhost:8080"

// Note that these test currently expect a jolokia java process to be
// running on 7025. Currently tested with a cassandra process2

func TestListDomains(t *testing.T) {
	domains, err := ListDomains(service)
	if err != nil {
		t.Errorf("err(%s) returned", err)
	}
	if len(domains) != 45 {
		t.Errorf("ListDomains = %v, want %v : %#v", len(domains), 45, domains)
	}
}

func TestListBeans(t *testing.T) {
	beans, err := ListBeans(service, "java.lang")
	if err != nil {
		t.Errorf("err(%s) returned", err)
	}
	if len(beans) != 14 {
		t.Errorf("ListBeans(java.lang) = %v, want %v : %#v", len(beans), 14, beans)
	}
}

func TestListProperties(t *testing.T) {
	props, err := ListProperties(service, "java.lang", "type=Threading")
	if err != nil {
		t.Errorf("err(%s), returned", err)
	}
	if len(props) != 17 {
		t.Errorf("ListProperties(type=Threading) = %v, want %v, : %#v", len(props), 17, props)
	}
}

func TestGetAttr(t *testing.T) {
	val, err := GetAttr(service, "java.lang", "type=Threading", "PeakThreadCount")
	if err != nil {
		t.Errorf("err(%s), returned", err)
	}
	if val.(float64) <= 100.0 {
		t.Errorf("GetAttr(PeakThreadCount) = %v, want > 100", val)
	}
}

func TestExecOp(t *testing.T) {
	val, err := ExecOp(service, "org.mobicents.ss7", "layer=SCCP,name=SccpStack,type=Router", "getRule", "1")
	if err != nil {
		t.Errorf("err(err(%s), returned", err)
	}
	tmp := val.(map[string]interface{})
	if tmp["ruleType"] != "SOLITARY" {
		t.Errorf("err(err(%s), returned", err)
	}
}

func TestListOperations(t *testing.T) {
	props, err := ListOperations(service, "org.mobicents.ss7", "layer=SCCP,name=SccpStack,type=Router")
	if err != nil {
		t.Errorf("err(%s), returned", err)
	}
	if len(props) != 19 {
		t.Errorf("ListOperations(type=Router) = %v, want %v, : %#v", len(props), 19, props)
	}
}

func TestClientListDomains(t *testing.T) {
	client := NewClient("localhost", "8080")
	domains, err := client.ListDomains()
	if err != nil {
		t.Errorf("err(%s) returned", err)
	}
	if len(domains) != 45 {
		t.Errorf("ListDomains = %v, want %v : %#v", len(domains), 45, domains)
	}
}

func TestClientListBeans(t *testing.T) {
	client := NewClient("localhost", "8080")
	beans, err := client.ListBeans("java.lang")
	if err != nil {
		t.Errorf("err(%s) returned", err)
	}
	if len(beans) != 14 {
		t.Errorf("ListBeans(java.lang) = %v, want %v : %#v", len(beans), 14, beans)
	}
}

func TestClientListProperties(t *testing.T) {
	client := NewClient("localhost", "8080")
	props, err := client.ListProperties("java.lang", "type=Threading")
	if err != nil {
		t.Errorf("err(%s), returned", err)
	}
	if len(props) != 17 {
		t.Errorf("ListProperties(type=Threading) = %v, want %v, : %#v", len(props), 17, props)
	}
}

func TestClientGetAttr(t *testing.T) {
	client := NewClient("localhost", "8080")
	val, err := client.GetAttr("java.lang", "type=Threading", "PeakThreadCount")
	if err != nil {
		t.Errorf("err(%s), returned", err)
		return
	}
	if val.(float64) <= 100.0 {
		t.Errorf("GetAttr(PeakThreadCount) = %v, want > 100", val)
		return
	}
}
