package golokia

import "testing"
import "fmt"

// Note that these test currently expect a jolokia java process to be
// running on 7025. Currently tested with a cassandra process2

var (
	host       = "localhost"
	port       = "8080"
	jolokia    = "jolokia-war-1.2.3"
	targetHost = "localhost"
	targetPort = "9999"
)

func TestListDomains(t *testing.T) {
	domains, err := ListDomains("http://" + host + ":" + port)
	if err != nil {
		t.Errorf("err(%s) returned", err)
	}
	fmt.Println("Domains: ", domains)
}

func TestListBeans(t *testing.T) {
	beans, err := ListBeans("http://"+host+":"+port, "java.lang")
	if err != nil {
		t.Errorf("err(%s) returned, err")
	}
	fmt.Println("Beans: ", beans)
}

func TestListProperties(t *testing.T) {
	props, err := ListProperties("http://"+host+":"+port, "java.lang", "type=Threading")
	if err != nil {
		t.Errorf("err(%s), returned, err")
	}
	fmt.Println("Properties: ", props)
}

func TestGetAttr(t *testing.T) {
	val, err := GetAttr("http://"+host+":"+port, "java.lang", "type=Threading", "PeakThreadCount")
	if err != nil {
		t.Errorf("err(%s), returned", err)
	}
	fmt.Println("Value:", val)
}

func TestClientListDomains(t *testing.T) {
	client := NewJolokiaClient("http://" + host + ":" + port + "/" + jolokia)
	client.SetTarget(targetHost + ":" + targetPort)
	domains, err := client.ListDomains()
	if err != nil {
		t.Errorf("err(%s) returned", err)
	}
	fmt.Println("Domains: ", domains)
}

func TestClientListBeans(t *testing.T) {
	client := NewJolokiaClient("http://" + host + ":" + port + "/" + jolokia)
	client.SetTarget(targetHost + ":" + targetPort)
	beans, err := client.ListBeans("java.lang")
	if err != nil {
		t.Errorf("err(%s) returned", err)
	}
	fmt.Println("Beans: ", beans)
}

func TestClientListProperties(t *testing.T) {
	client := NewJolokiaClient("http://" + host + ":" + port + "/" + jolokia)
	client.SetTarget(targetHost + ":" + targetPort)
	props, err := client.ListProperties("java.lang", []string{"type=Threading"})
	if err != nil {
		t.Errorf("err(%s), returned", err)
	}
	fmt.Println("Properties: ", props)
}

func TestClientGetAttr(t *testing.T) {
	client := NewJolokiaClient("http://" + host + ":" + port + "/" + jolokia)
	client.SetTarget(targetHost + ":" + targetPort)
	val, err := client.GetAttr("java.lang", []string{"type=Threading"}, "PeakThreadCount")
	if err != nil {
		t.Errorf("err(%s), returned", err)
		return
	}
	fmt.Println("Value:", val)
}
