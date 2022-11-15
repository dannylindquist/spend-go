package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/dannylindquist/spend-go/sqlite"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		cancel()
	}()
	fmt.Println("welcome to go")
	sql := sqlite.NewDB("data.sqlite")
	if err := sql.Open(); err != nil {
		fmt.Printf("couldn't open db: %v\n", err)
	}
	defer sql.Close()

	<-ctx.Done()
}