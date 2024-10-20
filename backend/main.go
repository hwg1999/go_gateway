package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hwg1999/go_gateway/backend/golang_common/lib"
	"github.com/hwg1999/go_gateway/backend/router"
)

func main() {
	lib.InitModule("./conf/dev/")
	defer lib.Destroy()
	router.HttpServerRun()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	router.HttpServerStop()
}
