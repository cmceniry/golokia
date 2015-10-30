package golokia

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
)

/* Jolokia Client properties connecting a Jolokia agent */
type JolokiaClient struct {
	// The url of the jolokia agent
	url string
	// The user name for the jolokia agent
	user string
	// The password for the jolokia agent (required if user is specified)
	pass string
	// The target host in a jolokia proxy setup
	target string
	// The target host jmx username
	targetUser string
	// The target host jmx password (required if target user is specified)
	targetPass string
}

/* Jolokia Request properties */
type JolokiaRequest struct {
	// The type of the operation (LIST, READ)
	opType string
	// The domain (mbean name)
	mbean string
	// The bean selection key
	properties []string
	// The attribute name (property name)
	attribute string
	// The path to access node of json tree
	path string
}

/* Jolokia Response for Generic request */
type JolokiaResponse struct {
	// Status of the request
	Status uint32
	// Timestamp value
	Timestamp uint32
	// Request
	Request map[string]interface{}
	// Result value
	Value map[string]interface{}
	// Error
	Error string
}

/* Jolokia Response for Read request */
type JolokiaReadResponse struct {
	// Status of the request
	Status uint32
	// Timestamp value
	Timestamp uint32
	// Request
	Request map[string]interface{}
	// Result value
	Value interface{}
	// Error
	Error string
}

type Target struct {
	Url      string `json:"url"`
	Password string `json:"password"`
	User     string `json:"user"`
}

/* The wrapper structure for json request */
type RequestData struct {
	Type  string `json:"type"`
	Mbean string `json:"mbean"`
	Path  string `json:"path"`
	//Attribute string `json:"attribute"`  -- Attribute get added manually
	Target `json:"target"`
}

/* ENUM to be used */
const (
	// READ for read operation
	READ = "READ"
	// LIST for list operation
	LIST = "LIST"
)

/* Http request param */
type httpRequest struct {
	Url  string
	Body []byte
}

/* Http response param */
type httpResponse struct {
	Status string
	Body   []byte
}

/* Internal function to perform http get request and return JolokiaResponse*/
func executeGetRequest(url string) (*JolokiaResponse, error) {
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

/* Internal function to perform http get request for attribute and return JolokiaReadResponse*/
func executeGetAttrRequest(url string) (*JolokiaReadResponse, error) {
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

/* Return a list of domain (mbean) available for a specific jolokia agent for specified url */
func ListDomains(service string) ([]string, error) {
	resp, err := executeGetRequest(service + "/jolokia/list?maxDepth=1")
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

/* Return a list of Beans (mbean type) available for a specific jolokia agent for specified url and domain*/
func ListBeans(service, domain string) ([]string, error) {
	resp, err := executeGetRequest(service + "/jolokia/list/" + domain + "?maxDepth=1")
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

/* Return a list of properties (attributes) available for a specific jolokia agent for specified url, domain and bean */
func ListProperties(service, domain, bean string) ([]string, error) {
	resp, err := executeGetRequest(service + "/jolokia/list/" + domain + "/" + bean + "?maxDepth=2")
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

/* Return a attribute value available for a specific url, domain, bean and attribute */
func GetAttr(service, domain, bean, attr string) (interface{}, error) {
	resp, err := executeGetAttrRequest(service + "/jolokia/read/" + domain + ":" + bean + "/" + attr)
	if err != nil {
		return "", err
	}
	return resp.Value, nil
}

/* Internal Function to perform the Http post request to Jolokia Agent */
func performPostRequest(request *httpRequest) (*httpResponse, error) {

	var url = request.Url
	var req *http.Request
	var newReqErr error

	//fmt.Println("Request Url: ", url)
	//fmt.Println("Request Body: ", string(request.Body))

	if request.Body != nil {
		var jsonStr = request.Body
		req, newReqErr = http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		if newReqErr != nil {
			//fmt.Printf("Request Could not be prepared")
			return nil, newReqErr
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, newReqErr = http.NewRequest("POST", url, nil)
		if newReqErr != nil {
			//fmt.Printf("Request Could not be prepared")
			return nil, newReqErr
		}
	}
	client := &http.Client{}

	resp, reqErr := client.Do(req)
	if reqErr != nil {
		fmt.Printf("Request Could not be send")
		return nil, reqErr
	}
	defer resp.Body.Close()

	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))

	response := &httpResponse{}
	response.Status = resp.Status
	response.Body = body

	return response, nil
}

/* Internal function used to make the Json request string using the structure */
func wrapRequestData(opType, mbeanName string, properties []string, path, attribute, targetUrl, targerUser, targetPass string) ([]byte, error) {
	mbean := mbeanName
	if properties != nil {
		for pos := range properties {
			mbean = mbean + ":" + properties[pos]
		}
	}
	target := ""
	if targetUrl != "" {
		target = "service:jmx:rmi:///jndi/rmi://" + targetUrl + "/jmxrmi"
	}
	requestData := RequestData{Type: opType, Mbean: mbean, Path: path, Target: Target{Url: target, Password: targetPass, User: targerUser}}
	// Marshal to the json string
	jsbyte, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}
	jsString := string(jsbyte)
	if attribute != "" {
		last := len(jsString)
		jsString = jsString[:last-1]
		attributeString := "\"attribute\":\"" + attribute + "\""
		jsString = jsString + ", " + attributeString + "}"
	}
	jsbyte = []byte(jsString)
	return jsbyte, nil
}

/* Creates a Jolokia Client for a specific url e.g. "http://127.0.0.1:8080/jolokia" */
func NewJolokiaClient(jolokiaUrl string) *JolokiaClient {
	jolokiaClient := &JolokiaClient{url: jolokiaUrl}
	return jolokiaClient
}

/* Set access credential to Jolokia client */
func (jolokiaClient *JolokiaClient) SetCredential(userName string, pass string) {
	jolokiaClient.user = userName
	jolokiaClient.pass = pass
}

/* Set target host when jolokia agent working in proxy architecture e.g. "10.0.1.96:7911"
(see: https://jolokia.org/reference/html/proxy.html) */
func (jolokiaClient *JolokiaClient) SetTarget(targetHost string) {
	jolokiaClient.target = targetHost
}

/* Set target host access credential to Jolokia client */
func (jolokiaClient *JolokiaClient) SetTargetCredential(userName string, pass string) {
	jolokiaClient.targetUser = userName
	jolokiaClient.targetPass = pass
}

/* Creates a new Jolokia request */
func NewJolokiaRequest(requestType, mbeanName string, properties []string, attribute string) *JolokiaRequest {
	jolokiaRequest := &JolokiaRequest{opType: requestType, mbean: mbeanName, attribute: attribute}
	if properties != nil {
		jolokiaRequest.properties = make([]string, len(properties))
		for pos := range properties {
			jolokiaRequest.properties[pos] = properties[pos]
		}
	}
	return jolokiaRequest
}

/* Set path to access tree node of json response */
func (jolokiaRequest *JolokiaRequest) SetPath(path string) {
	jolokiaRequest.path = path
}

/* Executes a jolokia request using a jolokia client and return response */
func (jolokiaClient *JolokiaClient) executePostRequest(jolokiaRequest *JolokiaRequest, pattern string) (string, error) {
	jsonReq, wrapError := wrapRequestData(jolokiaRequest.opType, jolokiaRequest.mbean, jolokiaRequest.properties, jolokiaRequest.path, jolokiaRequest.attribute, jolokiaClient.target, jolokiaClient.targetUser, jolokiaClient.targetPass)
	if wrapError != nil {
		return "", fmt.Errorf("JSON Wrap Failed: %v", wrapError)
	}
	requestUrl := jolokiaClient.url
	if pattern != "" {
		requestUrl = requestUrl + "/?" + pattern
	}
	request := &httpRequest{Url: requestUrl, Body: jsonReq}
	response, httpErr := performPostRequest(request)
	if httpErr != nil {
		return "", fmt.Errorf("HTTP Request Failed: %v", httpErr)
	}
	jolokiaResponse := string(response.Body)
	return jolokiaResponse, nil
}

/* Executes a jolokia request using a jolokia client and return response */
func (jolokiaClient *JolokiaClient) ExecuteRequest(jolokiaRequest *JolokiaRequest, pattern string) (*JolokiaResponse, error) {
	resp, requestErr := jolokiaClient.executePostRequest(jolokiaRequest, pattern)
	if requestErr != nil {
		return nil, requestErr
	}
	var respJ JolokiaResponse
	decodeErr := json.Unmarshal([]byte(resp), &respJ)
	if decodeErr != nil {
		return nil, fmt.Errorf("Failed to decode Jolokia resp : %v", decodeErr)
	}
	return &respJ, nil
}

/* Executes a jolokia read request using a jolokia client and return response */
func (jolokiaClient *JolokiaClient) ExecuteReadRequest(jolokiaRequest *JolokiaRequest) (*JolokiaReadResponse, error) {
	resp, requestErr := jolokiaClient.executePostRequest(jolokiaRequest, "")
	if requestErr != nil {
		return nil, requestErr
	}
	var respJ JolokiaReadResponse
	decodeErr := json.Unmarshal([]byte(resp), &respJ)
	if decodeErr != nil {
		return nil, fmt.Errorf("Failed to decode Jolokia resp : %v", decodeErr)
	}
	return &respJ, nil
}

/* List the domains using a jolokia client and return response */
func (jolokiaClient *JolokiaClient) ListDomains() ([]string, error) {
	request := NewJolokiaRequest(LIST, "", nil, "")
	resp, requestErr := jolokiaClient.ExecuteRequest(request, "maxDepth=1")
	if requestErr != nil {
		return nil, requestErr
	}
	ret := make([]string, 0, len(resp.Value))
	for key, _ := range resp.Value {
		ret = append(ret, key)
	}
	sort.Strings(ret)
	return ret, nil
}

/* List the beans of a specific domain using a jolokia client and return response */
func (jolokiaClient *JolokiaClient) ListBeans(domain string) ([]string, error) {
	request := NewJolokiaRequest(LIST, "", nil, "")
	request.SetPath(domain)
	resp, requestErr := jolokiaClient.ExecuteRequest(request, "maxDepth=1")
	if requestErr != nil {
		return nil, requestErr
	}
	ret := make([]string, 0, len(resp.Value))
	for key, _ := range resp.Value {
		ret = append(ret, key)
	}
	sort.Strings(ret)
	return ret, nil
}

/* List the properties of a specific bean and domain using a jolokia client and return response */
func (jolokiaClient *JolokiaClient) ListProperties(domain string, properties []string) ([]string, error) {
	request := NewJolokiaRequest(READ, domain, properties, "")
	resp, requestErr := jolokiaClient.ExecuteRequest(request, "")
	if requestErr != nil {
		return nil, requestErr
	}
	ret := make([]string, 0, len(resp.Value))
	for key, _ := range resp.Value {
		ret = append(ret, key)
	}
	sort.Strings(ret)
	return ret, nil
}

/* Get a attribute value of a specific property of an bean and a domain using a jolokia client */
func (jolokiaClient *JolokiaClient) GetAttr(domain string, properties []string, attribute string) (interface{}, error) {
	request := NewJolokiaRequest(READ, domain, properties, attribute)
	resp, requestErr := jolokiaClient.ExecuteReadRequest(request)
	if requestErr != nil {
		return nil, requestErr
	}
	return resp.Value, nil
}
