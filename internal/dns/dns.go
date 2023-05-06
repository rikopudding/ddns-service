package dns

import "ddns-service/config"

type DNSInstance interface {
	GetName() string
	GetFullDomain() string
	Init() error
	Compare(ip string) bool
	SetIP(ip string) error
	GetCachedIP() string
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

	for _, dnspodCfg := range config.Section.DnsApi.Dnspod {
		instance := NewDNSPod(dnspodCfg.SecretId, dnspodCfg.SecretKey, dnspodCfg.Domain, dnspodCfg.SubDomain)
		err := instance.Init()
		if err != nil {
			continue
		}
		Instances = append(Instances, instance)
	}

}
