package main

import (
	"io/ioutil"

	"os"

	"github.com/mholt/caddy"
	_ "github.com/mholt/caddy/caddyhttp"
	"io"
	"log"
	"net/http"
)

func init() {
	// configure default caddyfile
	caddy.SetDefaultCaddyfileLoader("default", caddy.LoaderFunc(defaultLoader))
}

func main() {
	caddy.AppName = "Sprocketplus"
	caddy.AppVersion = "1.2.3"

	// load caddyfile
	caddyfile, err := caddy.LoadCaddyfile("http")
	if err != nil {
		log.Fatal(err)
	}


	http.HandleFunc("/pushed", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello server push")
	})

	http.HandleFunc("/t", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello server push")
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if pusher, ok := w.(http.Pusher); ok {
			if err := pusher.Push("/pushed", nil); err != nil {
				log.Println("push failed")
			}
		}

		io.WriteString(w, "hello world")
	})

	http.ListenAndServeTLS(":1443", "localhost.cert", "localhost.key", nil)


	// start caddy server
	instance, err := caddy.Start(caddyfile)
	if err != nil {
		log.Fatal(err)
	}

	instance.Wait()

	}

// provide loader function
func defaultLoader(serverType string) (caddy.Input, error) {
	contents, err := ioutil.ReadFile(caddy.DefaultConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return caddy.CaddyfileInput{
		Contents:       contents,
		Filepath:       caddy.DefaultConfigFile,
		ServerTypeName: serverType,
	}, nil
}
