package grace

import (
	"context"
	"fmt"
	"github.com/eininst/flog"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const FIBER_CHILD_LOCK_FILE = "/tmp/fiber_child.lock"

type Config struct {
	Sig     os.Signal
	Timeout time.Duration
}

var DefaultConfig = Config{
	Sig:     syscall.SIGTERM,
	Timeout: time.Second * 10,
}
var pidMap map[int]int

func Listen(app *fiber.App, addr string, config ...Config) {
	cfg := DefaultConfig
	if len(config) > 0 {
		cfg = config[0]
	}
	go func() {
		if app.Config().Prefork {
			pidMap = make(map[int]int)
			app.Hooks().OnFork(func(i int) error {
				pidMap[i] = i
				return nil
			})
			_ = app.Listen(addr)
		} else {
			_ = app.Listen(addr)
		}
	}()
	listenSig(app, cfg)
}

func ListenTLS(app *fiber.App, addr string, certFile, keyFile string, config ...Config) {
	cfg := DefaultConfig
	if len(config) > 0 {
		cfg = config[0]
	}
	go func() {
		if app.Config().Prefork {
			pidMap = make(map[int]int)
			app.Hooks().OnFork(func(i int) error {
				pidMap[i] = i
				return nil
			})
			_ = app.ListenTLS(addr, certFile, keyFile)
		} else {
			_ = app.ListenTLS(addr, certFile, keyFile)
		}
	}()
	listenSig(app, cfg)
}

func listenSig(app *fiber.App, cfg Config) {
	if fiber.IsChild() {
		stop := make(chan int, 1)
		go func() {
			c := make(chan os.Signal, 1) // Create channel to signify a signal being sent
			signal.Notify(c, syscall.SIGTERM)
			for {
				sig := <-c
				switch sig {
				case syscall.SIGTERM:
					_ = app.Shutdown()
					fwrite(func(file *os.File) {
						_, _ = file.WriteString(fmt.Sprintf("%v\n", os.Getpid()))
					})
					flog.Infof("Child pid %s%v%s Shotdown", flog.Magenta, os.Getpid(), flog.Reset)
				case syscall.SIGINT:
					stop <- 1
					return
				}
			}
		}()
		<-stop
	} else {
		c := make(chan os.Signal, 1)
		signal.Notify(c, cfg.Sig)
		for {
			sig := <-c
			flog.Info("Received signal:", sig)
			switch sig {
			case cfg.Sig:
				f, _ := os.Create(FIBER_CHILD_LOCK_FILE)
				_ = f.Close()
				if app.Config().Prefork {
					stopChild(cfg.Timeout)
					_ = app.Shutdown()
				} else {
					stop(app, cfg.Timeout)
				}
				flog.Infof("Main pid %s%v%s Shotdown", flog.Magenta, os.Getpid(), flog.Reset)
				return
			}
		}
	}
}

func stop(app *fiber.App, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()
	chErr := make(chan error, 1)
	go func() {
		chErr <- app.Shutdown()
	}()

	for {
		select {
		case <-chErr:
			return
		case <-ctx.Done():
			return
		}
	}

}

func stopChild(timeout time.Duration) {
	for key := range pidMap {
		_ = syscall.Kill(key, syscall.SIGTERM)
	}
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	for {
		time.Sleep(time.Millisecond * 100)
		select {
		case <-ctx.Done():
			return
		default:
			content, err := os.ReadFile(FIBER_CHILD_LOCK_FILE)
			if err != nil {
				break
			}
			sarr := strings.Split(string(content), "\n")
			sax := map[int]int{}

			for _, id := range sarr {
				if v, err := strconv.Atoi(id); err == nil {
					sax[v] = v
				}
			}

			okCount := 0
			for id := range sax {
				if _, ok := pidMap[id]; ok {
					okCount += 1
				}
			}

			if okCount == len(pidMap) {
				cancel()
				return
			}
		}
	}
}

func fwrite(handler func(file *os.File)) {
	f, err := os.OpenFile(FIBER_CHILD_LOCK_FILE, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	defer func() { _ = f.Close() }()
	if err != nil {
		log.Println("create lock file failed", err)
		return
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_SH|syscall.LOCK_NB); err != nil {
		log.Println("add share lock in no block failed", err)
		return
	}
	handler(f)
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_UN); err != nil {
		log.Println("unlock share lock failed", err)
	}
	return
}
