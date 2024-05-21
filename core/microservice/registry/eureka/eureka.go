package eureka

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
)

var (
	headerAccept      = "Accept"
	headerContentType = "Content-Type"
)

func NewEurekaClient(serviceUrlList []string, username, password string, contentIsJson bool) *EurekaClient {
	headerMap := make(map[string]string)
	if contentIsJson {
		headerMap[headerAccept] = applicationJsonValue
		headerMap[headerContentType] = applicationJsonValue
	} else {
		headerMap[headerAccept] = applicationXmlValue
		headerMap[headerContentType] = applicationXmlValue
	}

	for i := range serviceUrlList {
		serviceUrlList[i] = strings.TrimRight(serviceUrlList[i], "/")
	}

	contentType := "application/json"
	if !contentIsJson {
		contentType = "application/xml"
	}

	return &EurekaClient{
		contextType: contentType,
		serviceUrl:  serviceUrlList[0],
		serviceUrls: serviceUrlList,
		username:    username,
		password:    password,
		headers:     headerMap,
	}
}

type EurekaClient struct {
	contextType        string
	serviceUrl         string
	urlIndex           int
	serviceUrls        []string
	username, password string
	headers            map[string]string
}

// register to eureka server
func (hc *EurekaClient) Register(info *InstanceInfo) (success bool, err error) {
	var body []byte
	if body, err = json.Marshal(info); err != nil {
		return
	}

	if result := hc.HttpRequest().POST("/apps/{appName}", RequestParameters{"appName": info.Instance.App}).BODY(string(body)).Do(); result.Error != nil {
		return false, result.Error
	} else if result.StatusCode == http.StatusOK || result.StatusCode == http.StatusNoContent {
		success = true
	}

	return
}

// deregister from eureka server
func (hc *EurekaClient) Deregister(appName, instanceId string) (bool, error) {
	if result := hc.HttpRequest().DELETE("/apps/{appName}/{instanceId}", RequestParameters{"instanceId": instanceId}).Do(); result.Error != nil {
		return false, result.Error
	} else {
		return result.StatusCode == http.StatusOK, nil
	}
}

// send heartbeat to eureka
func (hc *EurekaClient) SendHeartBeat(appName, instanceId, status, lastDirtyTimestamp, overriddenStatus string) (int, error) {
	// /apps/{appName}/{instanceId}?status={status}&lastDirtyTimestamp={lastDirtyTimestamp}
	if result := hc.HttpRequest().PUT(
		"/apps/{appName}/{instanceId}?status={status}&lastDirtyTimestamp={lastDirtyTimestamp}&overriddenStatus={overriddenStatus}",
		RequestParameters{
			"appName":            appName,
			"instanceId":         instanceId,
			"status":             status,
			"lastDirtyTimestamp": lastDirtyTimestamp,
			"overriddenStatus":   overriddenStatus,
		}).Do(); result.Error != nil {
		return 0, result.Error
	} else {
		return result.StatusCode, nil
	}
}

// update eureka client status
func (hc *EurekaClient) StatusUpdate(appName, instanceId, status, lastDirtyTimestamp string) (bool, error) {
	if result := hc.HttpRequest().PUT(
		"/apps/{appName}/{instanceId}/status?value={status}&lastDirtyTimestamp={lastDirtyTimestamp}",
		RequestParameters{
			"appName":            appName,
			"instanceId":         instanceId,
			"status":             status,
			"lastDirtyTimestamp": lastDirtyTimestamp,
		}); result.Error != nil {
		return false, result.Error
	} else {
		return result.StatusCode == http.StatusOK, nil
	}
}

// if regions is nil, get all application's information
// if regions is not nil,get application's information by regions
func (hc *EurekaClient) GetApplications(regions ...string) (*ApplicationList, error) {
	var rgs string
	if len(regions) > 0 {
		rgs = strings.Join(regions, ",")
	}

	var result ApplicationList
	if _, err := hc.HttpRequest().GET("/apps/?regions={regions}", RequestParameters{"regions": rgs}).Fetch(&result); err != nil {
		return nil, err
	} else {
		return &result, nil
	}
}

// get application's information by appName
func (hc *EurekaClient) GetApplication(appName string) (*ApplicationInfo, error) {
	var result ApplicationInfo
	if _, err := hc.HttpRequest().GET("/apps/{appName}", RequestParameters{"appName": appName}).Fetch(&result); err != nil {
		return nil, err
	} else {
		return &result, nil
	}
}

// get instance's information by appName and instanceId
func (hc *EurekaClient) GetInstance(appName, instanceId string) (*InstanceInfo, error) {
	var result InstanceInfo
	if _, err := hc.HttpRequest().GET("/apps/{appName}/{instanceId}", RequestParameters{"appName": appName, "instanceId": instanceId}).Fetch(&result); err != nil {
		return nil, err
	} else {
		return &result, nil
	}
}

// get instance's information by instanceId
func (hc *EurekaClient) GetInstanceById(instanceId string) (*InstanceInfo, error) {
	var result InstanceInfo
	if _, err := hc.HttpRequest().GET("/instances/{instanceId}", RequestParameters{"instanceId": instanceId}).Fetch(&result); err != nil {
		return nil, err
	} else {
		return &result, nil
	}
}

func (hc *EurekaClient) ServiceUrl() string {
	if len(hc.serviceUrls) > 1 {
		hc.urlIndex = (hc.urlIndex + 1) % len(hc.serviceUrls)
		hc.serviceUrl = hc.serviceUrls[hc.urlIndex]
	}
	return hc.serviceUrl
}

func (hc *EurekaClient) HttpRequest() *HttpRequest {
	return &HttpRequest{
		ContentType: "application/json",
		Base:        hc.ServiceUrl(),
		Headers:     hc.headers,
		Username:    hc.username,
		Password:    hc.password,
	}
}

type EurekaClientStub struct {
	*WrapperInstanceInfo
	client    *EurekaClient
	once      sync.Once
	wg        sync.WaitGroup
	closeChan chan struct{}
	config    *EurekaClientConfig
}

func (d *EurekaClientStub) Client() *EurekaClient {
	if d.client == nil {
		d.once.Do(func() {
			d.client = &EurekaClient{
				serviceUrl:  d.config.ServiceUrls[0],
				serviceUrls: d.config.ServiceUrls,
				username:    d.config.Username,
				password:    d.config.Password,
			}
		})
	}
	return d.client
}

func (d *EurekaClientStub) Shutdown() {
	close(d.closeChan)
	_, _ = d.Client().Deregister(d.App, d.InstanceId)
	d.wg.Wait()
}
