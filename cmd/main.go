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

	userService := sqlite.NewUserService(sql)

	user, err := userService.CreateUser(ctx, "danny@fromdl.com", "afsjd3kasjdf8");
	if err == nil {
		fmt.Printf("user: %v\n", user)
	} else {
		fmt.Printf("err: %v\n", err)
	}

	<-ctx.Done()
}