package main

import (
	"context"
	"fmt"
	"log"
	gohttp "net/http"
	"os"
	"os/signal"

	"github.com/rafaelrubbioli/fileapi/pkg/config"
	"github.com/rafaelrubbioli/fileapi/pkg/http"
	"github.com/rafaelrubbioli/fileapi/pkg/service"

	s3config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	ctx := context.Background()
	cfg, err := s3config.LoadDefaultConfig(ctx, s3config.WithRegion(config.AwsRegion))
	if err != nil {
		log.Fatalf("failed to load SDK configuration, %v", err)
	}

	client := s3.NewFromConfig(cfg)
	services := service.NewS3Service(client)

	handler, err := http.NewServer(services)
	if err != nil {
		log.Fatal(err)
	}

	server := gohttp.Server{
		Addr:    fmt.Sprintf(":%d", config.Port()),
		Handler: handler,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		log.Println("Shutting down")
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal("Shutdown: ", err)
		}
		close(idleConnsClosed)
	}()

	log.Println("Listening on port:", config.Port())
	log.Println(fmt.Sprintf("Explorer available at: %s/explorer", config.BaseURL()))
	if err := server.ListenAndServe(); err != gohttp.ErrServerClosed {
		log.Fatal("ListenAndServe: ", err)
	}

	<-idleConnsClosed
}
