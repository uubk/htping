package pkg

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"net/http"

	// All auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

type K8SHandler struct {
	version string
	host    string
	clusterName string
}

func InitK8SHandler(kubeconfig string, clusterName string) *K8SHandler {
	obj := &K8SHandler{
		clusterName: clusterName,
	}

	var config *rest.Config
	var err error
	if kubeconfig == "" {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	if err != nil {
		log.WithError(err).Fatal("Kubernetes config load failed")
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.WithError(err).Fatal("Kubernetes connection failed")
	}
	obj.host = config.Host

	version, err := clientset.ServerVersion()
	if err != nil {
		log.WithError(err).Warn("Couldn't get server version")
	}
	obj.version = version.String()

	return obj
}

func (h *K8SHandler) Handler(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	w.Header().Add("Content-Type", "application/javascript")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET, OPTIONS")

	// Function from StaticHandler.go
	results := GetBaseInfo(r)

	results["name"] = h.clusterName
	results["version"] = h.version
	results["master"] = h.host

	_, err := fmt.Fprintf(w, "window.clusters[\"" + h.clusterName + "\"] = ")
	if err != nil {
		log.WithError(err).Warn("Error when writing back to client")
	}
	err = encoder.Encode(results)
	if err != nil {
		log.WithError(err).Warn("Error in response serializer")
	}
}
