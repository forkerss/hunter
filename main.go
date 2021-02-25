package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/YE-Kits/hunter/config"
	"github.com/YE-Kits/hunter/hunter"
)

var (
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
	sig    <-chan os.Signal
)

// registerSignal 注册系统信号
func registerSignal(sigs ...os.Signal) <-chan os.Signal {
	s := make(chan os.Signal, 1)
	signal.Notify(s, sigs...)
	return s
}

func main() {
	var (
		err error
	)
	if err = config.SetFromFile("hunter.yaml"); err != nil {
		log.Fatalln(err)
	}

	wg = &sync.WaitGroup{}
	ctx, cancel = context.WithCancel(context.Background())
	defer func() {
		cancel()
		wg.Wait()
		log.Println("bye")
	}()

	if err = hunter.StartWebhook(ctx, wg); err != nil {
		log.Fatalln(err)
	}
	if err = hunter.StartXray(ctx, wg); err != nil {
		log.Println(err)
		return
	}

	if err = hunter.GenTask(ctx, wg); err != nil {
		log.Println(err)
		return
	}

	log.Println("start...")
	sig = registerSignal(syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
