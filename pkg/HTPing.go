package pkg

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
)

type HTPing struct {
	verbose    bool
	listen     string
	static     string
	kubeconfig string
	kube bool
	name string
}

// Add command line flags
func (htp *HTPing) RegisterArgs() {
	//flag.StringVar(&k8r.configPath, "config", "config.yml", "path to configuration file")
	flag.BoolVar(&htp.verbose, "verbose", false, "enable verbose logging")
	flag.BoolVar(&htp.kube, "kube", true, "whether this runs in kubernetes or not")
	flag.StringVar(&htp.name, "name", "", "how should this cluster be called?")
	flag.StringVar(&htp.listen, "listen", ":8080", "port and IP to listen on")
	flag.StringVar(&htp.static, "static", "static", "directory with static assets")
	flag.StringVar(&htp.kubeconfig, "kubeconfig", "", "when set, use this file as (absolute path) kubeconfig instead of attempting in-cluster config")
}

func (htp *HTPing) registerHTTPHandlers() {
	// Static files
	fs := http.FileServer(http.Dir(htp.static))
	http.Handle("/", fs)

	// Static text ping (debugging etc)
	http.HandleFunc("/sping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET, OPTIONS")
		_, err := fmt.Fprintf(w, "Pong!\n")
		if err != nil {
			log.WithError(err).Warn("Error when writing back to client")
		}
	})

	if htp.kube {
		// "Kubernetes" ping
		k8Handler := InitK8SHandler(htp.kubeconfig, htp.name)
		http.HandleFunc("/ping", k8Handler.Handler)
	} else {
		// "Static" ping
		staticHandler := InitStaticHandler(htp.name)
		http.HandleFunc("/ping", staticHandler.Handler)
	}
}

func (htp *HTPing) Listen() {
	log.WithFields(log.Fields{
		"listen": htp.listen,
	}).Info("Starting HTTP server")

	if htp.name == "" {
		log.Fatal("Name field is empty!")
	}

	htp.registerHTTPHandlers()

	// HTTP/2 via cleartext needs to be enabled manually
	handler := http.DefaultServeMux
	h2s := &http2.Server{}
	h1s := &http.Server{
		Addr:    htp.listen,
		Handler: h2c.NewHandler(handler, h2s),
	}
	err := h1s.ListenAndServe()
	if err != nil {
		log.WithError(err).Fatal("Error in HTTP server")
	}
}
