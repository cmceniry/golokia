package golokia

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

type JolokiaResponse struct {
	Status    uint32
	Timestamp uint32
	Request   map[string]interface{}
	Value     map[string]interface{}
	Error     string
}

type JolokiaReadResponse struct {
	Status    uint32
	Timestamp uint32
	Request   map[string]interface{}
	Value     interface{}
	Error     string
}

type Client struct {
	Service    string
	HttpClient *http.Client
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

func ListOperations(service, domain, bean string) ([]string, error) {
	resp, err := get(service + "/jolokia/list/" + domain + "/" + bean + "?maxDepth=2")
	if err != nil {
		return nil, err
	}
	if _, ok := resp.Value["op"]; !ok {
		return nil, errors.New("Invalid repsonse format - missing op")
	}
	respItems := resp.Value["op"].(map[string]interface{})
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

func ExecOp(service, domain, bean, operation string, value ...string) (interface{}, error) {
	resp, err := getAttr(service + "/jolokia/exec/" + domain + ":" + bean + "/" + operation + "/" + strings.Join(value, "/"))
	if err != nil {
		return "", err
	}
	return resp.Value, nil
}

func NewClient(host, port string) *Client {
	return &Client{
		"http://" + host + ":" + port,
		&http.Client{},
	}
}

func (c *Client) get(uripath string) (*JolokiaResponse, error) {
	resp, respErr := c.HttpClient.Get(c.Service + uripath)
	if respErr != nil {
		return nil, respErr
	}
	defer resp.Body.Close()
	var respJ JolokiaResponse
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&respJ); err != nil {
		return nil, err
	}
	if respJ.Status != 200 {
		return nil, fmt.Errorf("golokia: bad get response: %s", respJ.Error)
	}
	return &respJ, nil
}

func (c *Client) getAttr(uripath string) (*JolokiaReadResponse, error) {
	resp, respErr := c.HttpClient.Get(c.Service + uripath)
	if respErr != nil {
		return nil, respErr
	}
	defer resp.Body.Close()
	var respJ JolokiaReadResponse
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&respJ); err != nil {
		return nil, err
	}
	if respJ.Status != 200 {
		return nil, fmt.Errorf("golokia: bad get response: %s", respJ.Error)
	}
	return &respJ, nil
}

func (c *Client) ListDomains() ([]string, error) {
	resp, err := c.get("/jolokia/list?maxDepth=1")
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

func (c *Client) ListBeans(domain string) ([]string, error) {
	resp, err := c.get("/jolokia/list/" + domain + "?maxDepth=1")
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

func (c *Client) ListProperties(domain, bean string) ([]string, error) {
	resp, err := c.get("/jolokia/list/" + domain + "/" + bean + "?maxDepth=2")
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

func (c *Client) GetAttr(domain, bean, attr string) (interface{}, error) {
	resp, err := c.getAttr("/jolokia/read/" + domain + ":" + bean + "/" + attr)
	if err != nil {
		return "", err
	}
	return resp.Value, nil
}

func (c *Client) ExecOp(domain, bean, operation string, value ...string) (interface{}, error) {
	resp, err := c.getAttr("/jolokia/exec/" + domain + ":" + bean + "/" + operation + "/" + strings.Join(value, "/"))
	if err != nil {
		return "", err
	}
	return resp.Value, nil
}

func (c *Client) ListOperations(domain, bean string) ([]string, error) {
	resp, err := c.get("/jolokia/list/" + domain + "/" + bean + "?maxDepth=2")
	if err != nil {
		return nil, err
	}
	if _, ok := resp.Value["op"]; !ok {
		return nil, errors.New("Invalid repsonse format - missing op")
	}
	respItems := resp.Value["op"].(map[string]interface{})
	ret := make([]string, 0, len(respItems))
	for key, _ := range respItems {
		ret = append(ret, key)
	}
	sort.Strings(ret)
	return ret, nil
}
