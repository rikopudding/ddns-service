package dns

import (
	httpclient "ddns-service/pkgs/http-client"
	"encoding/json"
	"fmt"
)

type CloudflareInstance struct {
	Cfg       *CloudflareConfig
	headers   map[string]string
	CfID      string
	CurrentIp string
}

type CloudflareConfig struct {
	Email   string
	Zone    string
	AuthKey string
	Host    string
}

func NewCloudflare(email, zone, authKey, host string) *CloudflareInstance {
	headers := make(map[string]string)
	headers["X-Auth-Email"] = email
	headers["X-Auth-Key"] = authKey
	return &CloudflareInstance{
		Cfg: &CloudflareConfig{
			Email:   email,
			Zone:    zone,
			AuthKey: authKey,
			Host:    host,
		},
		headers: headers,
	}
}

type CFGetDNSListResp struct {
	Result []*CFDNSItem
}

type CFPostDNSResp struct {
	Result *CFDNSItem
}

type CFDNSItem struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

type CFCreateOrUpdateParams struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

func (cf *CloudflareInstance) Init() (err error) {

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records?type=A&name=%s", cf.Cfg.Zone, cf.Cfg.Host)
	body, err := httpclient.Get(url, cf.headers)


	var respData CFGetDNSListResp
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return
	}

	// DNS record not found
	if len(respData.Result) == 0 {
		return nil
	}

	cf.updateCFID(respData.Result[0].ID)
	cf.updateCurrentIP(respData.Result[0].Content)
	return nil
}

func (cf *CloudflareInstance) Compare(ip string) bool {
	return ip == cf.CurrentIp
}

func (cf *CloudflareInstance) SetIP(ip string) error {

	var reqBody = CFCreateOrUpdateParams{
		Type:    "A",
		Name:    cf.Cfg.Host,
		Content: ip,
		TTL:     1,
	}
	reqBodyJson, _ := json.Marshal(reqBody)
	if cf.CurrentIp == "" {
		url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", cf.Cfg.Zone)
		body, err := httpclient.PostJson(url, cf.headers, reqBodyJson)
		fmt.Sprintln(string(body))
		if err != nil {
			return err
		}
		var respBody CFPostDNSResp
		err = json.Unmarshal(body, &respBody)
		if err != nil {
			return err
		}
		if respBody.Result.ID != "" {
			cf.updateCFID(respBody.Result.ID)
			cf.updateCurrentIP(respBody.Result.Content)
		}
		
	} else {
		url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", cf.Cfg.Zone, cf.CfID)
		body, err := httpclient.PutJson(url, cf.headers, reqBodyJson)
		if err != nil {
			return err
		}
		var respBody CFGetDNSListResp
		err = json.Unmarshal(body, &respBody)
		if err != nil {
			return err
		}
		if len(respBody.Result)==0{
			return nil
		}
		if respBody.Result[0].ID != "" {
			cf.updateCFID(respBody.Result[0].ID)
			cf.updateCurrentIP(respBody.Result[0].ID)
		}
		
	}

	return nil
}
func (cf *CloudflareInstance) updateCurrentIP(ip string) {
	cf.CurrentIp = ip
}

func (cf *CloudflareInstance) updateCFID(cfId string) {
	cf.CfID = cfId
}
