package dns

import "ddns-service/config"

type DNSInstance interface {
	Init() error
	Compare(ip string) bool
	SetIP(ip string) error
}

var Instances []DNSInstance

func Init() {

	for _, cfCfg := range config.Section.DnsApi.Cloudflare {

		instance := NewCloudflare(cfCfg.Email, cfCfg.Zone, cfCfg.AuthKey, cfCfg.Host)
		err := instance.Init()
		if err != nil {
			continue
		}
		Instances = append(Instances, instance)
	}

}
