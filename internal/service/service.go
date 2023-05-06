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
	c.AddFunc("@every 1m", func() {
		getAndSet()
	})
	c.Run()
}

func getAndSet() {
	ip := ipgetter.GetIp()
	// log.Println("currentIp: ", ip)
	if ip == "" {
		return
	}

	for _, instance := range dns.Instances {
		oldIP := instance.GetCachedIP()
		if !instance.Compare(ip) || oldIP == "" {
			fmt.Println(instance.GetName() + " domain:" + instance.GetFullDomain())
			fmt.Println("old " + oldIP + ", new " + ip + " --- CreateOrUpdate START")
			var retryCount = 0
			for {
				if retryCount == 3 {
					break
				}
				err := instance.SetIP(ip)
				if err == nil {
					break
				}
				log.Println("dns record update failed: " + err.Error())
				retryCount++
			}
			fmt.Println("old " + oldIP + ", new " + ip + " --- CreateOrUpdate END")
		}
	}
}
