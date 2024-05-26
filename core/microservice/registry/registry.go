package registry

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/microservice/registry/config"

	"github.com/patrickmn/go-cache"
)

const (
	STATUS_UP             = "UP"             // Ready to receive traffic
	STATUS_DOWN           = "DOWN"           // Do not send traffic- healthcheck callback failed
	STATUS_STARTING       = "STARTING"       //Just about starting- initializations to be done - do not
	STATUS_OUT_OF_SERVICE = "OUT_OF_SERVICE" //Intentionally shutdown for traffic
	STATUS_UNKOWN         = "UNKNOWN"
)

type ServiceInfo struct {
	Name      string         `json:"name"`
	Instances []InstanceInfo `json:"instances"`
}

type PortInfo struct {
	Port    int         `json:"$,omitempty"`
	Enabled interface{} `json:"@enabled,omitempty"`
}

type DataCenterInfo struct {
	Class string `json:"@class,omitempty"`
	Name  string `json:"name,omitempty"`
}

type LeaseInfo struct {
	RenewalIntervalInSecs int   `json:"renewalIntervalInSecs,omitempty"`
	DurationInSecs        int   `json:"durationInSecs,omitempty"`
	RegistrationTimestamp int64 `json:"registrationTimestamp,omitempty"`
	LastRenewalTimestamp  int64 `json:"lastRenewalTimestamp,omitempty"`
	EvictionTimestamp     int64 `json:"evictionTimestamp,omitempty"`
	ServiceUpTimestamp    int64 `json:"serviceUpTimestamp,omitempty"`
}

type Metadata struct {
	ManagementPort string `json:"management.port,omitempty"`
	JmxPort        string `json:"jmx.port,omitempty"`
}

type InstanceInfo struct {
	ServiceName                   string         `json:"app,omitempty"`
	InstanceId                    string         `json:"instanceId,omitempty"`
	HostName                      string         `json:"hostName,omitempty"`
	IpAddr                        string         `json:"ipAddr,omitempty"`
	Port                          PortInfo       `json:"port,omitempty"`
	SecurePort                    PortInfo       `json:"securePort,omitempty"`
	Status                        string         `json:"status,omitempty"`
	OverriddenStatus              string         `json:"overriddenStatus,omitempty"`
	CountryId                     int            `json:"countryId,omitempty"`
	DataCenterInfo                DataCenterInfo `json:"dataCenterInfo,omitempty"`
	LeaseInfo                     LeaseInfo      `json:"leaseInfo,omitempty"`
	Metadata                      Metadata       `json:"metadata,omitempty"`
	HomePageUrl                   string         `json:"homePageUrl,omitempty"`
	StatusPageUrl                 string         `json:"statusPageUrl,omitempty"`
	HealthCheckUrl                string         `json:"healthCheckUrl,omitempty"`
	VipAddress                    string         `json:"vipAddress,omitempty"`
	SecureVipAddress              string         `json:"secureVipAddress,omitempty"`
	IsCoordinatingDiscoveryServer string         `json:"isCoordinatingDiscoveryServer,omitempty"`
	LastUpdatedTimestamp          string         `json:"lastUpdatedTimestamp,omitempty"`
	LastDirtyTimestamp            string         `json:"lastDirtyTimestamp,omitempty"`
	ActionType                    string         `json:"actionType,omitempty"`
}

// 服务发现
type DiscoveryClient interface {
	Register() (bool, error)
	Deregister() error

	GetServices() ([]ServiceInfo, error)
	GetService(serviceId string) (*ServiceInfo, error)

	GetInstances(serviceId string) ([]InstanceInfo, error)
	GetInstance(serviceId string) *InstanceInfo
	GetInstanceById(serviceId string, instanceId string) *InstanceInfo

	SendHeartBeat(instance *InstanceInfo, status string) (bool, error)
}

type AbstractDiscoveryClient struct {
}

func (c *AbstractDiscoveryClient) GetServices() ([]ServiceInfo, error) {
	return []ServiceInfo{}, nil
}

func (c *AbstractDiscoveryClient) GetService(serviceId string) (*ServiceInfo, error) {
	if services, err := c.GetServices(); err == nil {
		if len(services) > 0 {
			for i, si := range services {
				if si.Name == serviceId {
					return &services[i], nil
				}
			}
		}
		return nil, nil
	} else {
		return nil, err
	}
}

func (c *AbstractDiscoveryClient) GetInstances(serviceId string) ([]InstanceInfo, error) {
	if service, err := c.GetService(serviceId); err == nil {
		return service.Instances, nil
	} else {
		return nil, err
	}
}

func (c *AbstractDiscoveryClient) GetInstance(serviceId string) *InstanceInfo {
	if instances, err := c.GetInstances(serviceId); err == nil && len(instances) > 0 {
		activeInstances := []*InstanceInfo{}
		for i, ii := range instances {
			if ii.Status == "UP" {
				activeInstances = append(activeInstances, &instances[i])
			}
		}
		if len(activeInstances) > 0 {
			return activeInstances[rand.Intn(len(activeInstances))]
		}
	}
	return nil
}

func (c *AbstractDiscoveryClient) GetInstanceById(serviceId string, instanceId string) *InstanceInfo {
	if instances, err := c.GetInstances(serviceId); err == nil && len(instances) > 0 {
		for _, instance := range instances {
			if instance.InstanceId == instanceId {
				return &instance
			}
		}
	}
	return nil
}

func (*AbstractDiscoveryClient) SendHeartBeat() (bool, error) {
	return true, nil
}

// type DiscoveryClient interface {
// 	DiscoveryStatusUpdate(status string) (bool, error)
// 	GetApplications(regions ...string) (*ApplicationList, error)
// 	GetApplication(appName string) (*ApplicationInfo, error)
// 	GetInstance(appName, instanceId string) (*InstanceInfo, error)
// 	GetInstanceById(instanceId string) (*InstanceInfo, error)
// 	Shutdown()
// }

type Client interface {
	GetServiceEntry(service string) (string, error)
}

type RegistryClient struct {
	InstanceInfo
	discoveryClient DiscoveryClient `inject:"discoveryClient"`
	services        *cache.Cache
	wg              sync.WaitGroup
	closeChan       chan struct{}
}

func NewRegistryClient() *RegistryClient {
	currentTimeStr := fmt.Sprintf("%d", time.Now().UnixNano()/1000000)
	return &RegistryClient{
		InstanceInfo: InstanceInfo{
			ServiceName: config.Setting.ServiceName,
			InstanceId:  "",
			HostName:    "",
			IpAddr:      getPreferIP(),
			Status:      STATUS_UP,
			Port: PortInfo{
				Port:    config.Setting.Port,
				Enabled: true,
			},
			OverriddenStatus:     STATUS_UP,
			VipAddress:           config.Setting.ServiceName,
			SecureVipAddress:     config.Setting.ServiceName,
			LastUpdatedTimestamp: currentTimeStr,
			LastDirtyTimestamp:   currentTimeStr,
		},
		services:  cache.New(time.Minute*5, time.Minute),
		closeChan: make(chan struct{}),
	}
}

func (s *RegistryClient) Init() {
	if config.Setting.EnableAutoRegister {
		// 1. 自动注册
		s.Register()
	}

	// 2. 定时拉取服务
	go s.loadServicesTask()

	// 3. 和RegistryServer保持心跳
	go s.heartbeatTask()
}

func (s *RegistryClient) loadServicesTask() {
	ticker := time.NewTicker(time.Minute)
	s.wg.Add(1)
	for terminated := false; !terminated; {
		select {
		case <-ticker.C:
			_, err := s.loadServices()
			if err != nil {
				logger.Warn("RegistryClientStub load services error, ", err.Error())
			}
		case <-s.closeChan:
			terminated = true
		}
	}
	s.wg.Done()
}

func (s *RegistryClient) heartbeatTask() {
	ticker := time.NewTicker(time.Minute)
	s.wg.Add(1)
	for terminated := false; !terminated; {
		select {
		case <-ticker.C:
			_, err := s.sendHeartBeat()
			if err != nil {
				logger.Warn("RegistryClientStub send heartbeat error, ", err.Error())
			}
		case <-s.closeChan:
			terminated = true
		}
	}
	s.wg.Done()
}

// send heartbeat to eureka service
func (s *RegistryClient) sendHeartBeat() (success bool, err error) {
	success, err = s.discoveryClient.SendHeartBeat(&s.InstanceInfo, s.OverriddenStatus)
	if err != nil {
		return false, err
	}

	if !success {
		s.Status = "DIRTY"
		s.OverriddenStatus = "UP"
		s.LastDirtyTimestamp = fmt.Sprintf("%d", time.Now().UnixNano()/1000000)

		// try register
		if success, _ = s.Register(); success {
			s.Status = "UP"
		}
	} else {
		s.Status = s.OverriddenStatus
	}
	return true, nil
}

func (s *RegistryClient) loadServices() (bool, error) {
	// clear expired items
	s.services.DeleteExpired()

	// load current services
	if services, err := s.discoveryClient.GetServices(); err == nil {
		for _, si := range services {
			if instances, err := s.discoveryClient.GetInstances(si.Name); err == nil {
				si.Instances = instances
			}
			s.services.SetDefault(si.Name, &si)
		}
		return true, nil
	} else {
		return false, err
	}
}

func (s *RegistryClient) GetServices() []ServiceInfo {
	result := make([]ServiceInfo, 0)
	for _, item := range s.services.Items() {
		result = append(result, *item.Object.(*ServiceInfo))
	}
	return result
}

func (s *RegistryClient) GetService(name string) *ServiceInfo {
	if service, b := s.services.Get(name); b {
		return service.(*ServiceInfo)
	}

	// 从注册中心获取服务
	if service, err := s.discoveryClient.GetService(name); err == nil {
		if instances, err := s.discoveryClient.GetInstances(service.Name); err == nil {
			service.Instances = instances
		}
		s.services.SetDefault(name, service)
		return service
	}
	return nil
}

func (s *RegistryClient) GetInstances(serviceName string) []InstanceInfo {
	if si := s.GetService(serviceName); si != nil {
		return si.Instances
	}
	return nil
}

func (s *RegistryClient) GetInstance(serviceName string) *InstanceInfo {
	if si := s.GetService(serviceName); si != nil {
		if len(si.Instances) > 0 {
			return &si.Instances[rand.Intn(len(si.Instances))]
		}
	}
	return nil
}

func (s *RegistryClient) GetInstanceById(serviceName string, instanceId string) *InstanceInfo {
	if si := s.GetService(serviceName); si != nil {
		if len(si.Instances) > 0 {
			for _, ii := range si.Instances {
				if ii.InstanceId == instanceId {
					return &ii
				}
			}
		}
	}
	return nil
}

func (s *RegistryClient) Register() (bool, error) {
	// 微服务应用启动时调用此进行注册
	return s.discoveryClient.Register()
}

func (s *RegistryClient) Deregister() error {
	// 微服务应用启动时调用此进行注册
	return s.discoveryClient.Deregister()
}

func (s *RegistryClient) GetSerivceEntry(serviceName string) (string, error) {
	if si := s.GetService(serviceName); si != nil {
		if len(si.Instances) > 0 {
			instanceInfo := si.Instances[rand.Intn(len(si.Instances))]
			return fmt.Sprintf("http://%s:%d", instanceInfo.HostName, instanceInfo.Port.Port), nil
		} else {
			return "", errors.New("no instance")
		}
	}
	return "", nil
}

func (s *RegistryClient) Shutdown() {
	close(s.closeChan)
	s.Deregister()
	s.wg.Wait()
}

func getPreferIP() string {
	if config.Setting.PreferIP != "" {
		return config.Setting.PreferIP
	} else {
		return getLocalIP()
	}
}

// get local ip address
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for i := range addrs {
		// address type  not a loopback address
		if ipNet, ok := addrs[i].(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}

	return ""
}
