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

// GenCrawlerTask 生成爬虫任务
func GenCrawlerTask(ctx context.Context, wg *sync.WaitGroup) error {
	var rad *exec.Cmd
	f, err := os.Open(config.Target.File)
	if err != nil {
		return err
	}
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	for _, target := range strings.Split(string(buf), "\n") {
		target = strings.TrimSpace(target)
		log.Println("开始爬取:", target)
		// 堵塞式的运行 rad 避免浪费过多资源
		rad = exec.Command("bash", "-c", fmt.Sprintf("%s -t %s -http-proxy %s", config.Crawler.Radium.Path, target, config.Xray.Listen))
		if err := rad.Start(); err != nil {
			return err
		}
		wg.Add(1)
		go func(ctx context.Context, wg *sync.WaitGroup, cmd *exec.Cmd) {
			select {
			case <-ctx.Done():
				// 已经停止了
				if cmd.ProcessState != nil {
					return
				}
				if err := cmd.Process.Kill(); err != nil {
					log.Printf("Rad process kill: %s\n", err)
				}
			}
		}(ctx, wg, rad)
		rad.Wait()
		wg.Done()
	}
	return nil
}
