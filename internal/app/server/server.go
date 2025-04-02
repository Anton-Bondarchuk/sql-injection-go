package server

import (
	"errors"

	"github.com/gin-gonic/gin"
)


type options struct {
	port *int
	host *string
}

type Option func(options *options) error

func WithPort(port int) Option {
	return func(options *options) error {
		if port < 0 {
			return errors.New("port should be positive")		
		}

		options.port = &port; 

		return nil
	}
}

func WithHost(host string) Option {
	return func(options *options) error {
		if host == "" {
			return errors.New("host can't be empty string")		
		}

		options.host = &host;

		return nil
	}
}


type HttpServer struct {
	Srv *gin.Engine
	Host string
	Port int
}


func New(srv *gin.Engine, opts ...Option) (*HttpServer, error) {
	var options options
	for _, opt := range opts {
		err := opt(&options)
		if err != nil {
			return &HttpServer{}, err
		}
	}

	var port int = 8080
	if options.port == nil {
		options.port = &port
	}

	var defaultHost string = "0.0.0.0"
	if options.host == nil {
		options.host = &defaultHost;
	}

	return &HttpServer{
		Host: *options.host,
		Port: *options.port,
	}, nil
}