package main

import (
	"GitHub.com/mhthrh/HtmlElementsFinder/API"
	"GitHub.com/mhthrh/HtmlElementsFinder/Helper/LogUtil"
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	var ip string
	var port int
	flag.StringVar(&ip, "ip", "localhost", "What is listener IP address")
	flag.IntVar(&port, "port", 9999, "number of lines to read from the file")
	flag.Parse()

	l := LogUtil.NewLogger()

	server := http.Server{
		Addr:         fmt.Sprintf("%s:%d", ip, port),
		Handler:      &API.RequestHandler{},
		ErrorLog:     log.New(LogUtil.LogrusErrorWriter{}, "", 0),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  3600 * time.Second,
	}

	go func() {
		l.Println("Starting server on  %s:%d", ip, port)
		err := server.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	log.Println("Got signal:", <-c)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(ctx)

}
