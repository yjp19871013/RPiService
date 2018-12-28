package main

import (
	"context"
	"log"
	"net/http"
	"syscall"
	"time"

	"github.com/yjp19871013/RPiService/settings"

	"github.com/yjp19871013/RPiService/api/filestation/download_proxy/proxy"
	"github.com/yjp19871013/RPiService/db"
	"github.com/yjp19871013/RPiService/router"

	"github.com/gin-gonic/gin"
	DEATH "gopkg.in/vrecan/death.v3"
)

func main() {
	db.InitDb()
	defer db.CloseDb()

	proxy.StartProxy()
	defer proxy.StopProxy()

	r := gin.Default()

	router.InitRouter(r)

	srv := &http.Server{
		Addr:    settings.ServerPort,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	death := DEATH.NewDeath(syscall.SIGINT, syscall.SIGTERM)
	_ = death.WaitForDeath()
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}
