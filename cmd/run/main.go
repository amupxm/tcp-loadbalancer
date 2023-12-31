package main

import (
	"net"
	"time"

	config "github.com/amupxm/tcp-loadbalancer/cfg"
	"github.com/amupxm/tcp-loadbalancer/internal/client"
	"github.com/amupxm/tcp-loadbalancer/internal/loadbalancer"
	"github.com/amupxm/tcp-loadbalancer/internal/server"
	"github.com/amupxm/tcp-loadbalancer/pkg/logger"
)

type (
	Config struct {
		Servers []servers
		Clients []clients
	}
	servers struct {
		Host string
		Port int
	}
	clients struct {
		Host            string
		Port            int
		MessageInterval time.Duration
	}
)

func main() {

	cfg := config.GetConfig()

	lb := loadbalancer.NewLoadBalancer("localhost", 8080, "localhost", 8080, logger.NewLogger("run/loadbalancer"))

	for _, srv := range cfg.Config.Servers {
		lb.RegisterServer(srv.Host, srv.Port)
		go func(srv config.Server) {
			instance := server.NewServer(srv.Host, srv.Port, logger.NewLogger("run/server"), TCPHandler)
			instance.Listen()
		}(srv)
	}

	time.Sleep(1 * time.Second)
	go lb.StartAndListen()
	time.Sleep(1 * time.Second)
	for _, cl := range cfg.Config.Clients {
		go func(cl config.Client) {
			instance := client.NewClient(logger.NewLogger("run/client"), cl.Host, cl.Port)
			instance.Connect()
			time.Sleep(1 * time.Second)

			_, err := instance.SendMessage([]byte("hello\n"))
			if err != nil {
				logger.NewLogger("run/client").Error(err)
			}

		}(cl)
	}
	select {}
	// TODO : implement graceful shutdown
}

func TCPHandler(n *net.TCPConn, message string) bool {
	// TODO : implement handler for accepting more messages
	log := logger.NewLogger("run/server/tcpHandler")
	log.Infof("message received 2", message)
	server.ResponseString(n, message)

	return false
}
