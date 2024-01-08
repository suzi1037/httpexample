package main

import (
	"context"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type Handler struct {
	Content string
}

func (h Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	hostName, _ := os.Hostname()
	log.Println(req.RemoteAddr, req.Method, req.RequestURI)
	if req.URL.Path == "/slow" {
		query := req.URL.Query()
		timeoutSeconds := query.Get("t")
		duration := time.Duration(5)
		if timeoutSeconds != "" {
			tsInt, err := strconv.Atoi(timeoutSeconds)
			if err != nil {
				io.WriteString(resp, "slow error params")
				return
			}
			duration = time.Duration(tsInt)
		}
		time.Sleep(time.Second * duration)
		io.WriteString(resp, "slow")

	} else {
		io.WriteString(resp, h.Content)
	}
	io.WriteString(resp, "\n")
	io.WriteString(resp, hostName)
	io.WriteString(resp, "\n")
}

type Conf struct {
	Ip      string `yaml:"ip"`
	Port    int    `yaml:"port"`
	Content string `yaml:"content"`
	Timeout int    `yaml:"timeout"`
}

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	//read config file
	var conf Conf

	yamlFile, err := os.ReadFile("./conf/xx.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		panic(err)
	}

	handler := Handler{
		Content: conf.Content,
	}

	socketAddr := fmt.Sprintf("%s:%d", conf.Ip, conf.Port)
	log.Println("start server", socketAddr)

	server := http.Server{
		Addr:    socketAddr,
		Handler: handler,
	}

	go func() {
		err = server.ListenAndServe()
		if !errors.Is(http.ErrServerClosed, err) {
			panic(err)
		}
	}()

	sig := <-c
	log.Println("recevice signal:", sig)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(conf.Timeout))
	defer func() {
		cancel()
	}()

	if err = server.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown failed:", err)
	}

	log.Println("Server shutdown grace done")
}
