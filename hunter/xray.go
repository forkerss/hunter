package hunter

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/YE-Kits/hunter/config"
)

// StartXray 启动xray
func StartXray(ctx context.Context, wg *sync.WaitGroup) error {
	var (
		xray *exec.Cmd
	)
	xray = exec.Command("bash", "-c", fmt.Sprintf("%s webscan --listen %s --webhook-output http://%s/webhook", config.Xray.Path, config.Xray.Listen, config.WebHook.Listen))
	xray.Stdout = os.Stdout
	xray.Stderr = os.Stdout
	if err := xray.Start(); err != nil {
		return err
	}

	wg.Add(1)
	go func(ctx context.Context, wg *sync.WaitGroup, cmd *exec.Cmd) {
		select {
		case <-ctx.Done():
			if cmd.ProcessState != nil {
				return
			}
			if err := cmd.Process.Kill(); err != nil {
				log.Printf("xray process kill: %s\n", err)
			}
		}
		defer wg.Done()
		cmd.Wait()
	}(ctx, wg, xray)
	return nil
}
