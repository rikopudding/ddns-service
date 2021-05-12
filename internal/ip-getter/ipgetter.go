package ipgetter

import (
	"ddns-service/config"
	httpclient "ddns-service/pkgs/http-client"
	"regexp"
)

func GetIp() string {

	for _, url := range config.Section.IPGetter {

		body, err := httpclient.Get(url, nil)
		if err != nil {
			continue
		}
		r := regexp.MustCompile(`((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}`)
		find := r.Find(body)
		return string(find)
	}
	return ""
}
