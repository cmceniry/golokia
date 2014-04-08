package golokia

import "testing"

// Note that these test currently expect a jolokia java process to be
// running on 7025. Currently tested with a cassandra process2

func TestListDomains(t *testing.T) {
	domains, err := ListDomains("http://localhost:7025")
	if err != nil {
		t.Errorf("err(%s) returned", err)
	}
	if len(domains) != 11 {
		t.Errorf("ListDomains = %v, want %v : %#v", len(domains), 11, domains)
	}
}

func TestListBeans(t *testing.T) {
	beans, err := ListBeans("http://localhost:7025","java.lang")
	if err != nil {
		t.Errorf("err(%s) returned, err")
	}
	if len(beans) != 14 {
		t.Errorf("ListBeans(java.lang) = %v, want %v : %#v", len(beans), 14, beans)
	}
}

func TestListProperties(t *testing.T) {
	props, err := ListProperties("http://localhost:7025","java.lang","type=Threading")
	if err != nil {
		t.Errorf("err(%s), returned, err")
	}
	if len(props) != 16 {
		t.Errorf("ListProperties(type=Threading) = %v, want %v, : %#v", len(props), 16, props)
	}
}

func TestGetAttr(t *testing.T)  {
	val, err := GetAttr("http://localhost:7025","java.lang","type=Threading","PeakThreadCount")
	if err != nil {
		t.Errorf("err(%s), returned", err)
	}
	if val.(float64) <= 100.0 {
		t.Errorf("GetAttr(PeakThreadCount) = %v, want > 100", val)
	}
}
