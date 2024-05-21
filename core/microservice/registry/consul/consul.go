package consul

import (
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"

	"github.com/hashicorp/consul/api"
)

type ConsulRegistryClient struct {
	ConsulClient *api.Client
}

func (m *ConsulRegistryClient) GetServiceEntry(service string) (string, error) {
	serviceEntries, _, err := m.ConsulClient.Health().Service(service, "", true, nil)
	if err != nil {
		return "", err
	}

	if len(serviceEntries) == 0 {
		return "", errors.New("no service to connect")
	}

	serviceEntry := serviceEntries[rand.Intn(len(serviceEntries))]

	scheme := "http"
	if meta, ok := serviceEntry.Service.Meta["secure"]; ok {
		if secure, err := strconv.ParseBool(meta); err == nil && secure {
			scheme = "https"
		}
	}

	url := &url.URL{
		Scheme: scheme,
		Host:   fmt.Sprintf("%s:%d", serviceEntry.Service.Address, serviceEntry.Service.Port),
	}

	return url.String(), nil
}
