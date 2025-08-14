package main

import (
	"context"
	"exchangeapp/config"
	"exchangeapp/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	config.InitConfig()
	r := router.SetupRouter()
	port := config.Appconfig.App.Port
	if port == "" {
		port = ":8080"
	}
	srv := &http.Server{
		Addr: port,
		Handler: r,
	}
	// add gracefully shut down server
	go func ()  {
		if err := srv.ListenAndServe(); err != nil && err!=http.ErrServerClosed{
			log.Fatalf("Listen: %s\n",err)
		}
	}()
	quit := make(chan os.Signal,1)
	signal.Notify(quit,os.Interrupt)
	<-quit
	log.Println("Shutdown Sever ...")

	ctx, cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil{
		log.Fatal("Server Shutdown: ",err)
	}
	log.Println("Server exiting")
}
