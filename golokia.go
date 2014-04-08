package golokia

import (
	"encoding/json"
	"net/http"
	"errors"
	"sort"
)

type JolokiaResponse struct {
	Status    uint32
	Timestamp uint32
	Request   map[string]interface{}
	Value     map[string]interface{}
}

type JolokiaReadResponse struct {
	Status    uint32
	Timestamp uint32
	Request   map[string]interface{}
	Value     interface{}
}

func get(url string) (*JolokiaResponse, error) {
	resp, respErr := http.Get(url)
	if respErr != nil {
		return nil, respErr
	}
        defer resp.Body.Close()
	var respJ JolokiaResponse
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&respJ); err != nil {
		return nil, err
	}
	return &respJ, nil
}

func getAttr(url string) (*JolokiaReadResponse, error) {
	resp, respErr := http.Get(url)
	if respErr != nil {
		return nil, respErr
	}
        defer resp.Body.Close()
	var respJ JolokiaReadResponse
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&respJ); err != nil {
		return nil, err
	}
	return &respJ, nil
}

func ListDomains(service string) ([]string, error) {
	resp, err := get(service + "/jolokia/list?maxDepth=1")
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0, len(resp.Value))
	for key, _ := range resp.Value {
		ret = append(ret, key)
	}
	sort.Strings(ret)
	return ret, nil
}

func ListBeans(service, domain string) ([]string, error) {
	resp, err := get(service + "/jolokia/list/" + domain + "?maxDepth=1")
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0, len(resp.Value))
	for key, _ := range resp.Value {
		ret = append(ret, key)
	}
	sort.Strings(ret)
	return ret, nil
}

func ListProperties(service, domain, bean string) ([]string, error) {
	resp, err := get(service + "/jolokia/list/" + domain + "/" + bean + "?maxDepth=2")
	if err != nil {
		return nil, err
	}
	if _, ok := resp.Value["attr"]; !ok {
		return nil, errors.New("Invalid repsonse format - missing attr")
	}
	respItems := resp.Value["attr"].(map[string]interface{})
	ret := make([]string, 0, len(respItems))
	for key, _ := range respItems {
		ret = append(ret, key)
	}
	sort.Strings(ret)
	return ret, nil
}

func GetAttr(service, domain, bean, attr string) (interface{}, error) {
	resp, err := getAttr(service + "/jolokia/read/" + domain + ":" + bean + "/" + attr)
	if err != nil {
		return "", err
	}
	return resp.Value, nil
}
