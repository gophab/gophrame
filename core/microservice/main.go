package microservice

import (
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var app naming_client.INamingClient

func Application() naming_client.INamingClient {
	return app
}

func Service(configs []constant.ServerConfig, serverName, group string, serverPort uint64) {
	client, _ := clients.CreateNamingClient(map[string]any{
		"serverConfigs": configs,
		"clientConfig": constant.ClientConfig{
			NamespaceId:  "phsc",
			TimeoutMs:    10 * 1000,
			BeatInterval: 5 * 1000,
			//CacheDir:            "data/nacos/cache",
			NotLoadCacheAtStart: true,
		},
	})
	ip := getIpAddr()
	RegisterServiceInstance(client, vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        serverPort,
		ServiceName: serverName,
		Weight:      1,
		GroupName:   group,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata: map[string]string{
			"preserved.heart.beat.interval": strconv.Itoa(1000 * 10), //25s
			"preserved.register.source":     "SPRING_CLOUD",
		},
	})

	//连接
	app = client
}

func getIpAddr() string {
	address, _ := net.InterfaceAddrs()
	for _, address := range address {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

func RegisterServiceInstance(client naming_client.INamingClient, param vo.RegisterInstanceParam) {
	success, _ := client.RegisterInstance(param)
	if success {
		log.Printf("[INFO] 服务名 [%s] 注册成功  address [%s:%d] \n", param.ServiceName, param.Ip, param.Port)
	} else {
		log.Fatalf("[ERROR] 服务名 [%s] 注册失败  address [%s:%d] \n", param.ServiceName, param.Ip, param.Port)
	}
	go func() {
		exitChan := make(chan os.Signal, 100)
		signal.Notify(exitChan, os.Interrupt, syscall.SIGTERM, os.Kill)
		<-exitChan
		log.Printf("[EXIT] 服务关闭 [%s]  address [%s:%d] \n", param.ServiceName, param.Ip, param.Port)
		_, _ = client.DeregisterInstance(vo.DeregisterInstanceParam{
			Ip:          param.Ip,
			Port:        param.Port,
			Cluster:     param.ClusterName,
			ServiceName: param.ServiceName,
			GroupName:   param.GroupName,
			Ephemeral:   true, //立刻删除服务
		})
		os.Exit(1)
	}()

}
