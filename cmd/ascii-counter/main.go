package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/stainour/test12/internal/counter"
	"github.com/stainour/test12/internal/reader"
)

func main() {
	if len(os.Args) < 2 {
		errorExit("please provide folder path as an argument")
	}

	path := os.Args[1]

	r, err := reader.NewFileReader(path)
	if err != nil {
		errorExit(err.Error())
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-c
		cancel()
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()

	counter := counter.NewASCIICounter(r)

	count, err := counter.Count(ctx)
	if err != nil {
		errorExit(fmt.Sprintf("error counting symbols: %s", err))
	}

	err = printHistogram(count)
	if err != nil {
		errorExit(fmt.Sprintf("error printing results: %s", err))
	}
}

func errorExit(m string) {
	println(m)
	os.Exit(1)
}
