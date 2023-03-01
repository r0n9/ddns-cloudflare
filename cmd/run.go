package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/r0n9/ddns-cloudflare/notify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"time"
)

type RespResults struct {
	Results  []Result `json:"result"`
	Success  bool     `json:"success"`
	Errors   []string `json:"errors"`
	Messages []string `json:"messages"`
}

type RecordResult struct {
	Result     Result   `json:"result"`
	CreatedOn  string   `json:"created_on"`
	ModifiedOn string   `json:"modified_on"`
	Success    bool     `json:"success"`
	Errors     []string `json:"errors"`
	Messages   []string `json:"messages"`
}

type Result struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type DnsRecord struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Ttl     int8   `json:"ttl"`
	Proxied bool   `json:"proxied"`
}

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "run the program once",
	Long:  `run the program once, the cloudflare is defined in config file`,
	Run: func(cmd *cobra.Command, args []string) {
		Init()

		externalIp := getExternalIp()
		if externalIp == "" {
			return
		}

		log.Printf("Date: %v\n", time.Now().Format("2006/01/02 15:04:05"))
		log.Printf("Public IP Address: %v\n", externalIp)

		names := Config.Domains

		needNotify := false
		for _, name := range names {
			updated, err := dnsUpdate(name, externalIp)
			if err != nil {
				log.Errorf(err.Error())
			} else {
				needNotify = needNotify || updated
			}
		}

		if Config.SendKey != "" && needNotify {
			err := notify.Send(Config.SendKey, externalIp)
			if err != nil {
				log.Errorf(err.Error())
			} else {
				log.Println("ServerChan notify success")
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(RunCmd)
}

func dnsUpdate(domainName, ip string) (bool, error) {

	zoneId := Config.ZoneId
	apiKey := Config.ApiKey
	mail := Config.Email

	createUrl := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zoneId)
	getHostIdUrl := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records?type=A&name=", zoneId)

	client := http.Client{}
	req, _ := http.NewRequest("GET", getHostIdUrl+domainName, nil)
	req.Header = http.Header{
		"X-Auth-Email": []string{mail},
		"Content-Type": []string{"application/json"},
		"X-Auth-Key":   []string{apiKey},
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("error: %v", err)
		return false, errors.New("failed to query dns records")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Errorf("error: call dns records api failed [status=%v], please check the config", resp.StatusCode)
		return false, errors.New("failed to query dns records")
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("error: %v", err)
		return false, errors.New("failed to query dns records")
	}
	var res RespResults
	err = json.Unmarshal(respBody, &res)
	if err != nil {
		log.Errorf("error: %v", err)
		return false, errors.New("failed to query dns records")
	}

	body := &DnsRecord{
		Type:    "A",
		Name:    domainName,
		Content: ip,
		Ttl:     1,
		Proxied: false,
	}
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)

	if len(res.Results) == 0 {
		log.Printf("Record not exist，creating new dns record：%v\n", domainName)
		updateReq, _ := http.NewRequest("POST", createUrl, payloadBuf)
		updateReq.Header = http.Header{
			"X-Auth-Email": []string{mail},
			"Content-Type": []string{"application/json"},
			"X-Auth-Key":   []string{apiKey},
		}
		updateResp, err := client.Do(updateReq)
		if err != nil {
			log.Errorf("error: %v", err)
			return false, errors.New("failed to create new dns record")
		}
		defer updateResp.Body.Close()
		updateRespBody, err := ioutil.ReadAll(updateResp.Body)
		if err != nil {
			log.Errorf("error: %v", err)
			return false, errors.New("failed to create new dns record")
		}
		var updateRes RecordResult
		err = json.Unmarshal(updateRespBody, &updateRes)
		if err != nil {
			log.Infof("error: %v", err)
			return false, errors.New("failed to create new dns record")
		}

		log.Printf("Created new dns record: %v => %v\n", updateRes.Result.Name, updateRes.Result.Content)
		return false, nil

	} else {
		recordId := res.Results[0].Id
		content := res.Results[0].Content
		if content == ip {
			log.Printf("Same dns record: %v => %v\n", domainName, content)
			return false, nil
		}

		log.Printf("Different dns record, start to update record_id: %v\n", recordId)

		updateReq, _ := http.NewRequest("PUT", createUrl+"/"+recordId, payloadBuf)
		updateReq.Header = http.Header{
			"X-Auth-Email": []string{mail},
			"Content-Type": []string{"application/json"},
			"X-Auth-Key":   []string{apiKey},
		}
		updateResp, err := client.Do(updateReq)
		if err != nil {
			log.Infof("error: %v", err)
			return false, errors.New("failed to update dns record")
		}
		defer updateResp.Body.Close()
		updateRespBody, err := ioutil.ReadAll(updateResp.Body)
		if err != nil {
			log.Infof("error: %v", err)
			return false, errors.New("failed to update dns record")
		}
		var updateRes RecordResult
		err = json.Unmarshal(updateRespBody, &updateRes)
		if err != nil {
			log.Infof("error: %v", err)
			return false, errors.New("failed to update dns record")
		}

		log.Printf("Updated dns record: %v => %v\n", updateRes.Result.Name, updateRes.Result.Content)
		return true, nil
	}
}

func getExternalIp() string {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		log.Errorf("error: %v", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Errorf("error: call http://myexternalip.com/raw failed [status=%v], please check the network", resp.StatusCode)
		return ""
	}

	content, _ := ioutil.ReadAll(resp.Body)
	return string(content)
}
