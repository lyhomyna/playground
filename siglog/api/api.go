package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type SiglogServer struct {
    http *httpServer 
}

func (s *SiglogServer) Run(ctx context.Context) error {
    s.http = &httpServer{}
    ctx, cancel := context.WithCancel(ctx)

    var httpEc = make(chan error, 1)
    go func() {
	err := s.http.Run(ctx)
	if err != nil {
	    err = fmt.Errorf("HTTP server error. %w", err)
	}
	httpEc <- err
    }()
    
    err := <- httpEc 
    
    cancel()
    return err
}

type httpServer struct {
    http *http.Server
}

func (s *httpServer) Run(ctx context.Context) error {
    handler := NewHttpServer()

    s.http = &http.Server{
	Addr: "localhost:8080",
	Handler: handler,
	ReadHeaderTimeout: 5 * time.Second, // mitigate risk of Slowloris Attack
    }
    
    log.Println("HTTP server is running on port :8080")
    if err := s.http.ListenAndServe(); err != http.ErrServerClosed {
	return err
    }

    return nil
}

