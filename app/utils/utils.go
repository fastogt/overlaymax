package utils

import (
	"errors"
	"net/http"
	"os/user"
	"path/filepath"

	"github.com/rs/cors"
)

func NewHttpServer(addr string, corsEnabled bool, handler http.Handler) *http.Server {
	if corsEnabled {
		options := cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
		}
		c := cors.New(options)
		corsHandler := c.Handler(handler)
		return &http.Server{Addr: addr, Handler: corsHandler}
	}

	return &http.Server{Addr: addr, Handler: handler}
}

func PreparePath(path string) (*string, error) {
	if len(path) == 0 || (path[0] != '~' && path[0] != '/') {
		return nil, errors.New("invalid input Path to file")
	}
	if path[0] == '/' {
		return &path, nil
	}
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	path = filepath.Join(usr.HomeDir, path[1:])
	return &path, nil
}
