package handles

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

type Host struct {
	domain map[string]string
	sync.RWMutex
}

func NewHost(filename string) *Host {
	h := Host{
		domain: map[string]string{},
	}
	h.InitHosts(filename)
	return &h
}

func (h *Host) InitHosts(filename string) {
	buf, err := os.Open("./conf/" + filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer buf.Close()
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		line := scanner.Text()

		// comments
		if strings.HasPrefix(line, "#") {
			continue
		}

		// pattern： sli[0] domain; sli[1 ip
		sli := strings.Split(line, " ")
		if len(sli) != 2 {
			continue
		}

		// 验证domain、ip
		if IsUsefulIp(sli[1]) == false {
			continue
		}

		if len(strings.Split(sli[1], "|")) == 2 {
			R.Create(sli[0])
		}

		h.Set(sli[0], sli[1])
	}
}

func (h *Host) Get(name string) (string, bool) {
	h.Lock()
	defer h.Unlock()
	url, ok := h.domain[name]
	return url, ok
}

func (h *Host) Set(name, url string) {
	h.Lock()
	defer h.Unlock()
	h.domain[name] = url
}

func (h *Host) Delete(name string) {
	h.Lock()
	h.Unlock()
	delete(h.domain, name)
}

func (h *Host) Refresh(filename string) {

}

func ParserUrl(domain, ip string) ([]string, bool) {
	// 新建时，已经验证过了
	// 这里不再做验证了
	var ips []string
	seg := strings.Split(ip, "|")
	if len(seg) > 2 {
		return ips, false
	}
	if R.Get(domain)%2 == 0 && len(seg) == 2 {
		return append(ips, seg[1]), true
	}
	nor := strings.Split(seg[0], "&")
	for _, i := range nor {
		ips = append(ips, i)
	}
	return ips, true
}

func isIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func IsUsefulIp(ip string) bool {
	seg := strings.Split(ip, "|")
	if len(seg) > 2 {
		return false
	}
	nor := strings.Split(seg[0], "&")
	for _, i := range nor {
		if isIP(i) == false {
			return false
		}
	}
	if len(seg) == 2 {
		if isIP(seg[1]) == false {
			return false
		}
	}
	return true
}

type Counter struct {
	mapCount map[string]int
	sync.RWMutex
}

func (c *Counter) Create(domain string) {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.mapCount[domain]; ok == false {
		c.mapCount[domain] = 1
	}
}

func (c *Counter) Get(domain string) int {
	c.Lock()
	defer c.Unlock()
	count, ok := c.mapCount[domain]
	if ok {
		c.mapCount[domain]++
		return count
	}
	return 1
}
