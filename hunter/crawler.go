package hunter

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/YE-Kits/hunter/config"
)

// GenTask 生成任务
func GenTask(ctx context.Context, wg *sync.WaitGroup) error {
	f, err := os.Open(config.Target.File)
	if err != nil {
		return err
	}
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	for _, target := range strings.Split(string(buf), "\n") {
		if err = CrawlerScan(ctx, target, wg); err != nil {
			return err
		}
	}
	return nil
}

// CrawlerScan 启动爬虫扫描
func CrawlerScan(ctx context.Context, target string, wg *sync.WaitGroup) error {
	if config.Crawler.Radium.Enable {
		return startRadium(ctx, target, wg)
	}
	return startCrawlergo(ctx, target, wg)
}

func startRadium(ctx context.Context, target string, wg *sync.WaitGroup) error {
	var (
		cmd *exec.Cmd
	)
	cmd = exec.Command("bash", "-c", fmt.Sprintf("%s -t %s -http-proxy %s", config.Crawler.Radium.Path, target, config.Xray.Listen))
	if err := cmd.Start(); err != nil {
		return err
	}

	wg.Add(1)
	go func(wg *sync.WaitGroup, cmd *exec.Cmd) {
		defer wg.Done()
		cmd.Wait()
	}(wg, cmd)
	go func(ctx context.Context, cmd *exec.Cmd) {
		select {
		case <-ctx.Done():
			if cmd.ProcessState != nil {
				return
			}
			if err := cmd.Process.Kill(); err != nil {
				log.Printf("Rad process kill: %s\n", err)
			}
		}
	}(ctx, cmd)
	return nil
}

func startCrawlergo(ctx context.Context, target string, wg *sync.WaitGroup) error {
	return nil
}
