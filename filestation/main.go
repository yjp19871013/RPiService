package main

import (
	"context"
	"log"
	"net/http"
	"syscall"
	"time"

	"github.com/yjp19871013/RPiService/filestation/download_proxy"

	"github.com/yjp19871013/RPiService/filestation/db"
	"github.com/yjp19871013/RPiService/filestation/dto"
	"github.com/yjp19871013/RPiService/filestation/router"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	validator "gopkg.in/go-playground/validator.v8"
	DEATH "gopkg.in/vrecan/death.v3"
)

func main() {
	db.InitDb()
	defer db.CloseDb()

	download_proxy.StartProxy()
	defer download_proxy.StopProxy()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("url_validator", dto.UrlValidator)
		if err != nil {
			log.Println("err:", err)
		}
	}

	r := gin.Default()

	router.InitRouter(r)

	srv := &http.Server{
		Addr:    ":10002",
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
