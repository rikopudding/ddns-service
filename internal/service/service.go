package service

import (
	"ddns-service/internal/dns"
	ipgetter "ddns-service/internal/ip-getter"
	"fmt"
	"log"

	"github.com/robfig/cron/v3"
)

func Start() {

	ipgetter.GetIp()
	dns.Init()
	getAndSet()
	c := cron.New()
	c.AddFunc("@every 5m", func() {
		getAndSet()
	})
}

func getAndSet() {
	ip := ipgetter.GetIp()
	fmt.Println("currentIp: ",ip)
	if ip == "" {
		return
	}

	for _, instance := range dns.Instances {
		if !instance.Compare(ip) {
			var retryCount = 0
			for {
				if retryCount==3{
					break
				}
				err := instance.SetIP(ip)
				if err==nil{
					break
				}
				log.Println("dns record update failed: "+err.Error())
				retryCount++
			}
		}
	}
}
