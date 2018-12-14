package main

import (
	"context"
	"log"
	"net/http"
	"syscall"
	"time"

	"github.com/yjp19871013/RPiService/users/db"

	"github.com/yjp19871013/RPiService/users/router"

	"github.com/gin-gonic/gin"
	DEATH "gopkg.in/vrecan/death.v3"
)

func main() {
	db.InitDb()
	defer db.CloseDb()

	r := gin.Default()

	router.InitRouter(r)

	srv := &http.Server{
		Addr:    ":10001",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	death := DEATH.NewDeath(syscall.SIGINT, syscall.SIGTERM)
	death.WaitForDeath()
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}
