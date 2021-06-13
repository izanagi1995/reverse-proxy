package proxy

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

type ReverseProxyManager struct {
	servers   []*http.Server
	CloseChan chan os.Signal
	waitGroup *sync.WaitGroup
}

type DomainWithTargets map[string][]string

func NewReverseProxyManager(config map[string][]string) (*ReverseProxyManager, error) {
	byPorts := make(map[string]DomainWithTargets, 0)
	var err error
	for host, targets := range config {
		pUrl, errParse := url.Parse(host)
		if errParse != nil {
			err = errParse
			break
		}
		port := ":" + pUrl.Port()
		if value, ok := byPorts[port]; ok {
			value[pUrl.Host] = targets
		} else {
			byPorts[port] = DomainWithTargets{
				pUrl.Host: targets,
			}
		}
	}
	if err != nil {
		return nil, err
	}
	servers := make([]*http.Server, len(byPorts))
	i := 0
	for port, domains := range byPorts {
		lb := LoadBalancedReverseProxy{
			Port:   port,
			VHosts: make([]*VirtualHost, len(domains)),
		}
		j := 0
		for domain, targets := range domains {
			lb.VHosts[j], err = NewVirtualHost(domain, targets)
			if err != nil {
				return nil, err
			}
			j++
		}
		servers[i] = lb.BuildHttpServer()
	}
	return &ReverseProxyManager{
		servers:   servers,
		CloseChan: make(chan os.Signal),
	}, nil
}

func (r *ReverseProxyManager) Run() {
	r.waitGroup = new(sync.WaitGroup)
	r.waitGroup.Add(len(r.servers))
	for _, server := range r.servers {
		go func(s *http.Server) {
			s.ListenAndServe()
			r.waitGroup.Done()
		}(server)
	}

	go func() {
		<-r.CloseChan
		for _, server := range r.servers {
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			server.Shutdown(shutdownCtx)
		}
	}()

	r.waitGroup.Wait()
}
