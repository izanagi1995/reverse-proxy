package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type LoadBalancedReverseProxy struct {
	Port   string
	VHosts []*VirtualHost
}

type VirtualHost struct {
	Domain       string
	targetHosts  []url.URL
	currentIndex int
	updateMutex  *sync.Mutex
}

func NewVirtualHost(domain string, targets []string) (*VirtualHost, error) {
	urls := make([]url.URL, len(targets))
	for i, target := range targets {
		u, err := url.Parse(target)
		if err != nil {
			return nil, err
		}
		urls[i] = *u
	}
	vh := VirtualHost{
		Domain:       domain,
		targetHosts:  urls,
		currentIndex: 0,
		updateMutex:  &sync.Mutex{},
	}
	return &vh, nil
}

func (v *VirtualHost) updateAndGetIndex() int {
	v.updateMutex.Lock()
	v.currentIndex++
	if v.currentIndex == len(v.targetHosts) {
		v.currentIndex = 0
	}
	value := v.currentIndex
	v.updateMutex.Unlock()
	return value
}

func (v *VirtualHost) ValidateHost(host string) *url.URL {
	for _, target := range v.targetHosts {
		if target.Host == host {
			return &target
		}
	}
	return nil
}

func NewReverseProxy(virthost *VirtualHost) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		cookie, err := req.Cookie("RP-Target")
		if err != nil && err != http.ErrNoCookie {
			log.Println(err.Error())
		} else if err == nil {
			url := virthost.ValidateHost(cookie.Value)
			if url != nil {
				req.URL.Scheme = url.Scheme
				req.URL.Host = url.Host
				return
			}
		}
		nextTarget := virthost.targetHosts[virthost.updateAndGetIndex()]
		req.URL.Scheme = nextTarget.Scheme
		req.URL.Host = nextTarget.Host

	}
	reverseProxy := &httputil.ReverseProxy{Director: director}
	return reverseProxy
}

func (lb *LoadBalancedReverseProxy) BuildHttpServer() *http.Server {
	vhosts := make(map[string]*httputil.ReverseProxy, len(lb.VHosts))
	for _, vh := range lb.VHosts {
		vhosts[vh.Domain] = NewReverseProxy(vh)
	}
	handler := lbHandler{
		hosts: vhosts,
	}
	s := http.Server{Addr: lb.Port, Handler: &handler}
	return &s
}
