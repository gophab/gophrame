package eureka

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/wjshen/gophrame/core/microservice/registry"
	"github.com/wjshen/gophrame/core/microservice/registry/eureka/config"
)

func CreateEurekaDiscoveryClient() (*EurekaDiscoveryClient, error) {
	return nil, nil
}

type EurekaDiscoveryClient struct {
	registry.AbstractDiscoveryClient
	client *EurekaClient
	once   sync.Once
}

func (d *EurekaDiscoveryClient) Client() *EurekaClient {
	if d.client == nil {
		d.once.Do(func() {
			d.client = NewEurekaClient(
				config.Setting.ServiceUrls,
				config.Setting.Username,
				config.Setting.Password,
				true)
		})
	}
	return d.client
}

func (c *EurekaDiscoveryClient) Register(instance *registry.InstanceInfo) (bool, error) {
	return c.Client().Register(&InstanceInfo{
		Instance: wrapperInstanceInfo(instance),
	})
}

func (c *EurekaDiscoveryClient) Deregister(instance *registry.InstanceInfo) error {
	_, err := c.Client().Deregister(instance.ServiceName, instance.InstanceId)
	return err
}

func (c *EurekaDiscoveryClient) GetServices() ([]registry.ServiceInfo, error) {
	if list, err := c.Client().GetApplications(); err == nil {
		return mapServiceInfos(list.Applications.Applications), nil
	} else {
		return nil, err
	}
}

func (c *EurekaDiscoveryClient) GetService(serviceId string) (*registry.ServiceInfo, error) {
	if appInfo, err := c.Client().GetApplication(serviceId); err == nil && appInfo != nil {
		return mapServiceInfo(appInfo.Application), nil
	} else {
		return nil, err
	}
}

func (c *EurekaDiscoveryClient) GetInstances(serviceId string) ([]registry.InstanceInfo, error) {
	if appInfo, err := c.Client().GetApplication(serviceId); err == nil && appInfo != nil {
		return mapInstanceInfos(appInfo.Application.Instance), nil
	} else {
		return nil, err
	}
}

func (c *EurekaDiscoveryClient) GetInstance(serviceId string) *registry.InstanceInfo {
	if appInfo, err := c.Client().GetApplication(serviceId); err == nil && appInfo != nil {
		return mapInstanceInfo(&appInfo.Application.Instance[rand.Intn(len(appInfo.Application.Instance))])
	} else {
		return nil
	}
}

func (c *EurekaDiscoveryClient) GetInstanceById(serviceId string, instanceId string) *registry.InstanceInfo {
	if ii, err := c.Client().GetInstance(serviceId, instanceId); err == nil && ii != nil {
		return mapInstanceInfo(ii.Instance)
	}
	return nil
}

func (c *EurekaDiscoveryClient) SendHeartBeat(instance *registry.InstanceInfo, status string) (bool, error) {
	if statusCode, err := c.Client().SendHeartBeat(instance.ServiceName, instance.InstanceId, status, instance.LastDirtyTimestamp, instance.OverriddenStatus); err == nil {
		switch statusCode {
		case http.StatusOK, http.StatusNoContent:
			return true, nil
		case http.StatusNotFound:
			break
		}
		return false, nil
	} else {
		return false, err
	}
}

func wrapperInstanceInfo(instance *registry.InstanceInfo) *WrapperInstanceInfo {
	renewalIntervalInSecs := defaultLeaseRenewalInterval
	durationInSecs := defaultLeaseDuration
	if instance.LeaseInfo.RenewalIntervalInSecs > 0 {
		renewalIntervalInSecs = instance.LeaseInfo.RenewalIntervalInSecs
	}
	if instance.LeaseInfo.DurationInSecs > 0 {
		durationInSecs = instance.LeaseInfo.DurationInSecs
	}

	hostName := defaultHostName
	if instance.HostName != "" {
		hostName = instance.HostName
	}

	currentTimeStr := fmt.Sprintf("%d", time.Now().UnixNano()/1000000)

	return &WrapperInstanceInfo{
		InstanceId: instance.InstanceId,
		HostName:   hostName,
		App:        instance.ServiceName,
		IpAddr:     instance.IpAddr,
		Status:     statusUp,
		Port: &WrapperPort{
			Port:    instance.Port.Port,
			Enabled: true,
		},
		DataCenterInfo: &WrapperDataCenterInfo{
			Class: defaultDataCenterInfoClass,
			Name:  defaultDataCenterInfoName,
		},
		LeaseInfo: &WrapperLeaseInfo{
			RenewalIntervalInSecs: renewalIntervalInSecs,
			DurationInSecs:        durationInSecs,
		},
		Metadata: &WrapperMetadata{
			ManagementPort: fmt.Sprintf("%d", instance.Port.Port),
		},
		VipAddress:           instance.ServiceName,
		SecureVipAddress:     instance.ServiceName,
		LastUpdatedTimestamp: currentTimeStr,
		LastDirtyTimestamp:   currentTimeStr,
	}
}

func mapServiceInfo(app *WrapperApplicationInfo) *registry.ServiceInfo {
	return &registry.ServiceInfo{
		Name:      app.Name,
		Instances: mapInstanceInfos(app.Instance),
	}
}

func mapServiceInfos(services []WrapperApplicationInfo) []registry.ServiceInfo {
	var result []registry.ServiceInfo
	for _, ii := range services {
		result = append(result, *mapServiceInfo(&ii))
	}
	return result
}

func mapInstanceInfo(instance *WrapperInstanceInfo) *registry.InstanceInfo {
	return &registry.InstanceInfo{
		ServiceName: instance.App,
		InstanceId:  instance.InstanceId,
		HostName:    instance.HostName,
		IpAddr:      instance.IpAddr,
		Port: registry.PortInfo{
			Port:    instance.Port.Port,
			Enabled: instance.Port.Enabled,
		},
		Status:           instance.Status,
		OverriddenStatus: instance.OverriddenStatus,
		StatusPageUrl:    instance.StatusPageUrl,
		DataCenterInfo: registry.DataCenterInfo{
			Class: instance.DataCenterInfo.Class,
			Name:  instance.DataCenterInfo.Name,
		},
		LeaseInfo: registry.LeaseInfo{
			RenewalIntervalInSecs: instance.LeaseInfo.RenewalIntervalInSecs,
			DurationInSecs:        instance.LeaseInfo.DurationInSecs,
		},
		Metadata: registry.Metadata{
			ManagementPort: instance.Metadata.ManagementPort,
		},
		VipAddress:           instance.VipAddress,
		SecureVipAddress:     instance.SecureVipAddress,
		LastUpdatedTimestamp: instance.LastUpdatedTimestamp,
		LastDirtyTimestamp:   instance.LastDirtyTimestamp,
	}
}

func mapInstanceInfos(instances []WrapperInstanceInfo) []registry.InstanceInfo {
	var result []registry.InstanceInfo
	for _, ii := range instances {
		result = append(result, *mapInstanceInfo(&ii))
	}
	return result
}
