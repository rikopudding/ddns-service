package dns

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

type DNSPodInstance struct {
	Name      string
	Cfg       *DNSPodConfig
	Client    *dnspod.Client
	CurrentIp string
	RecordId  uint64
}

type DNSPodConfig struct {
	SecretId  string
	SecretKey string
	Domain    string
	SubDomain string
}

func NewDNSPod(secretId, secretKey, domain, subDomain string) *DNSPodInstance {
	// headers := make(map[string]string)
	// headers["X-Auth-Email"] = email
	// headers["X-Auth-Key"] = authKey
	return &DNSPodInstance{
		Name: "dnspod",
		Cfg: &DNSPodConfig{
			SecretId:  secretId,
			SecretKey: secretKey,
			Domain:    domain,
			SubDomain: subDomain,
		},
	}
}

// type CFGetDNSListResp struct {
// 	Result []*CFDNSItem
// }

// type CFPostDNSResp struct {
// 	Result *CFDNSItem
// }

// type CFDNSItem struct {
// 	ID      string `json:"id"`
// 	Content string `json:"content"`
// }

// type CFCreateOrUpdateParams struct {
// 	Type    string `json:"type"`
// 	Name    string `json:"name"`
// 	Content string `json:"content"`
// 	TTL     int    `json:"ttl"`
// }

func (ins *DNSPodInstance) GetName() string {
	return ins.Name
}

func (ins *DNSPodInstance) GetFullDomain() string {
	return ins.Cfg.SubDomain + "." + ins.Cfg.Domain
}

func (ins *DNSPodInstance) Init() (err error) {

	credential := common.NewCredential(ins.Cfg.SecretId, ins.Cfg.SecretKey)

	client, err := dnspod.NewClient(credential, regions.Guangzhou, profile.NewClientProfile())
	if err != nil {
		return
	}

	ins.Client = client

	req := dnspod.NewDescribeRecordListRequest()
	req.Domain = common.StringPtr(ins.Cfg.Domain)

	resp, err := ins.Client.DescribeRecordList(req)
	if err != nil {
		return
	}

	for _, rec := range resp.Response.RecordList {
		if *rec.Name == ins.Cfg.SubDomain && *rec.Type == "A" {
			ins.updateRecordID(*rec.RecordId)
			ins.updateCurrentIP(*rec.Value)
			break
		}
	}

	// url := fmt.Sprintf("https://api.DNSPod.com/client/v4/zones/%s/dns_records?type=A&name=%s", cf.Cfg.Zone, cf.Cfg.Host)
	// body, err := httpclient.Get(url, cf.headers)

	// var respData CFGetDNSListResp
	// err = json.Unmarshal(body, &respData)
	// if err != nil {
	// 	return
	// }

	// // DNS record not found
	// if len(respData.Result) == 0 {
	// 	return nil
	// }

	// // ins.updateCFID(respData.Result[0].ID)
	// ins.updateCurrentIP(respData.Result[0].Content)
	return nil
}

func (ins *DNSPodInstance) Compare(ip string) bool {
	return ip == ins.CurrentIp
}

func (ins *DNSPodInstance) SetIP(ip string) error {

	if ins.RecordId == 0 {
		req := dnspod.NewCreateRecordRequest()
		req.Domain = &ins.Cfg.Domain
		req.SubDomain = &ins.Cfg.SubDomain
		req.RecordType = common.StringPtr("A")
		req.RecordLine = common.StringPtr("默认")
		req.Value = common.StringPtr(ip)
		resp, err := ins.Client.CreateRecord(req)
		if err != nil {
			return err
		}
		ins.updateRecordID(*resp.Response.RecordId)
	} else {
		req := dnspod.NewModifyRecordRequest()
		req.RecordId = &ins.RecordId
		req.Domain = &ins.Cfg.Domain
		req.SubDomain = &ins.Cfg.SubDomain
		req.RecordType = common.StringPtr("A")
		req.RecordLine = common.StringPtr("默认")
		req.Value = common.StringPtr(ip)
		_, err := ins.Client.ModifyRecord(req)
		if err != nil {
			return err
		}
	}
	ins.updateCurrentIP(ip)
	return nil
}

// func (cf *DNSPodInstance) SetIP(ip string) error {

// 	var reqBody = CFCreateOrUpdateParams{
// 		Type:    "A",
// 		Name:    cf.Cfg.Host,
// 		Content: ip,
// 		TTL:     1,
// 	}
// 	reqBodyJson, _ := json.Marshal(reqBody)
// 	if cf.CurrentIp == "" {
// 		url := fmt.Sprintf("https://api.DNSPod.com/client/v4/zones/%s/dns_records", cf.Cfg.Zone)
// 		body, err := httpclient.PostJson(url, cf.headers, reqBodyJson)
// 		fmt.Sprintln(string(body))
// 		if err != nil {
// 			return err
// 		}
// 		var respBody CFPostDNSResp
// 		err = json.Unmarshal(body, &respBody)
// 		if err != nil {
// 			return err
// 		}
// 		if respBody.Result.ID != "" {
// 			cf.updateCFID(respBody.Result.ID)
// 			cf.updateCurrentIP(respBody.Result.Content)
// 		}

// 	} else {
// 		url := fmt.Sprintf("https://api.DNSPod.com/client/v4/zones/%s/dns_records/%s", cf.Cfg.Zone, cf.CfID)
// 		body, err := httpclient.PutJson(url, cf.headers, reqBodyJson)
// 		if err != nil {
// 			return err
// 		}
// 		var respBody CFGetDNSListResp
// 		err = json.Unmarshal(body, &respBody)
// 		if err != nil {
// 			return err
// 		}
// 		if len(respBody.Result) == 0 {
// 			return nil
// 		}
// 		if respBody.Result[0].ID != "" {
// 			cf.updateCFID(respBody.Result[0].ID)
// 			cf.updateCurrentIP(respBody.Result[0].ID)
// 		}

// 	}

//		return nil
//	}

func (ins *DNSPodInstance) updateRecordID(recordId uint64) {
	ins.RecordId = recordId
}

func (ins *DNSPodInstance) updateCurrentIP(ip string) {
	ins.CurrentIp = ip
}

func (ins *DNSPodInstance) GetCachedIP() string {
	return ins.CurrentIp
}
