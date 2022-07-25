package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aimzeter/autonotif"
	"github.com/aimzeter/autonotif/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log.Println("INFO | read config")
	cfg, err := config.NewConfig(os.Getenv("CONFIG_PATH"))
	if err != nil {
		log.Fatal("ERROR |", err)
	}

	log.Println("INFO | build dependencies and runner")
	runner, err := build(cfg)
	if err != nil {
		log.Fatal("ERROR |", err)
	}

	log.Println("INFO | serve http server")
	go serveHTTP(runner)

	log.Println("INFO | scheduler run")
	if err := schedule(cfg, runner); err != nil {
		log.Fatal("ERROR |", err)
	}

	log.Println("INFO | scheduler stop")
}

func build(cfg *config.Config) (*autonotif.Autonotif, error) {
	d, err := autonotif.BuildDependencies(cfg)
	if err != nil {
		return nil, err
	}

	return autonotif.BuildAutonotif(d), nil
}

func serveHTTP(runner *autonotif.Autonotif) {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/healthz", runner.HealthHandler)
	http.HandleFunc("/force-last-id", runner.ForceLastIDHandler)
	_ = http.ListenAndServe(":8080", nil)
}

func schedule(cfg *config.Config, runner *autonotif.Autonotif) error {
	terminated := make(chan bool)
	go awaitTermination(terminated)

	for {
		timer := time.NewTimer(time.Duration(cfg.Base.SchedulerPeriod) * time.Second)
		select {
		case <-terminated:
			_ = runner.Terminate()
			return nil
		case <-timer.C:
			_ = runner.Run()
		}
		timer.Stop()
	}
}

func awaitTermination(terminated chan<- bool) {
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	<-sign
	terminated <- true
}
