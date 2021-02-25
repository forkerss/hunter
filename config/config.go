package config

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type xray struct {
	Path   string `yaml:"path"`
	Listen string `yaml:"listen"`
}

type crawler struct {
	Chrome string `yaml:"chrome"`
	Radium struct {
		Enable bool   `yaml:"enable"`
		Path   string `yaml:"path"`
	} `yaml:"rad"`
	Crawlergo struct {
		Enable bool   `yaml:"enable"`
		Path   string `yaml:"path"`
	} `yaml:"crawlergo"`
}

type webhook struct {
	Enable    bool   `yaml:"enable"`
	Listen    string `yaml:"listen"`
	DingRobot struct {
		Enable  bool   `yaml:"enable"`
		WebHook string `yaml:"webhook"`
		Secret  string `yaml:"secret"`
	} `yaml:"dingrobot"`
}

type target struct {
	File string `yaml:"file"`
}

type config struct {
	Xray    xray    `yaml:"xray"`
	Crawler crawler `yaml:"crawler"`
	Webhook webhook `yaml:"webhook"`
	Target  target  `yaml:"target"`
}

var (
	conf config = config{}
	// Xray export xray config
	Xray xray = xray{}
	// Crawler export crawler config
	Crawler crawler = crawler{}
	// WebHook export webhook config
	WebHook webhook = webhook{}
	// Target export target config
	Target target = target{}
)

// SetFromFile 从文件设置 config
func SetFromFile(c string) error {
	var (
		f   *os.File
		buf []byte
		err error
	)
	if f, err = os.Open(c); err != nil {
		return err
	}
	if buf, err = ioutil.ReadAll(f); err != nil {
		return err
	}

	if err = yaml.Unmarshal(buf, &conf); err != nil {
		return err
	}

	Xray = conf.Xray
	Crawler = conf.Crawler
	WebHook = conf.Webhook
	Target = conf.Target
	return nil
}
