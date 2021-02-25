package hunter

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/YE-Kits/hunter/config"
	"github.com/YE-Kits/hunter/dingrobot"
)

var robot dingrobot.Roboter

type vulnDetail struct {
	Addr     string                 `json:"addr"`
	Payload  string                 `json:"payload"`
	SnapShot []interface{}          `json:"snapshot"`
	Extra    map[string]interface{} `json:"extra"`
}

type webVuln struct {
	Type string `json:"type"`
	Data struct {
		CreateTime int        `json:"create_time"`
		Detail     vulnDetail `json:"detail"`
		Plugin     string     `json:"plugin"`
		Target     struct {
			Params []string `json:"params"`
			URL    string   `json:"url"`
		}
	} `json:"data"`
}

// StartWebhook 启动webhook
func StartWebhook(ctx context.Context, wg *sync.WaitGroup) error {
	if !config.WebHook.Enable {
		return nil
	}
	if config.WebHook.DingRobot.Enable {
		robot = dingrobot.NewRobot(config.WebHook.DingRobot.WebHook)
		robot.SetSecret(config.WebHook.DingRobot.Secret)
	}
	http.HandleFunc("/webhook", webhook)
	s := http.Server{Addr: config.WebHook.Listen, Handler: http.DefaultServeMux}
	wg.Add(1)
	log.Println("webhook:", config.WebHook.Listen)
	go s.ListenAndServe()
	go func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()
		select {
		case <-ctx.Done():
			s.Shutdown(ctx)
		}
	}(ctx, wg)
	return nil
}

func webhook(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		return
	}
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		return
	}
	defer req.Body.Close()
	bufStr := string(buf)
	if strings.Contains(bufStr, "web_vuln") {
		vuln := webVuln{}
		if err = json.Unmarshal(buf, &vuln); err != nil {
			log.Println("DingRobot:", err)
			return
		}
		if err = pushMsgWithDingRobot(vuln); err != nil {
			log.Println("DingRobot:", err)
		}
	}
}

func pushMsgWithDingRobot(vuln webVuln) error {
	if !config.WebHook.DingRobot.Enable {
		return nil
	}
	var msg, extra, snapshot string

	for k, v := range vuln.Data.Detail.Extra {
		extra += fmt.Sprintf(" %s:\n", k)
		switch b := v.(type) {
		case map[string]string:
			for k1, v1 := range b {
				extra += fmt.Sprintf("  %s:%s\n", k1, v1)
			}
		case map[string]interface{}:
			for k1, v1 := range b {
				extra += fmt.Sprintf("  %s:%v\n", k1, v1)
			}
		case []string:
			for _, v1 := range b {
				extra += fmt.Sprintf("  %s\n", v1)
			}
		}
	}

	for _, v := range vuln.Data.Detail.SnapShot {
		switch b := v.(type) {
		case []string:
			for _, v1 := range b {
				snapshot += fmt.Sprintf(" %s\n", v1)
			}
		case []interface{}:
			for _, v1 := range b {
				snapshot += fmt.Sprintf(" %v\n", v1)
			}
		}
	}
	msg = fmt.Sprintf("plugin: '%s' \naddr: '%s' \nextra:--- \n%s\npayload: '%v' \nsnapshot:--- \n%s\n", vuln.Data.Plugin, vuln.Data.Detail.Addr, extra, vuln.Data.Detail.Payload, snapshot)
	return robot.SendText("Xray 漏洞信息 \n\n"+msg, nil, false)
}
