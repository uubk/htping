package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	// All auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

type StaticHandler struct {
	clusterName string
}

func InitStaticHandler(clusterName string) *StaticHandler {
	obj := &StaticHandler{
		clusterName: clusterName,
	}

	return obj
}

func GetBaseInfo(r *http.Request) map[string]interface{} {
	return map[string]interface{}{
		"time":       time.Now(),
		"proto":      strconv.Itoa(r.ProtoMajor) + "." + strconv.Itoa(r.ProtoMinor),
		"remote":     r.RemoteAddr,
		"fwd_remote": r.Header.Get("X-Forwarded-For"),
		"fwd_uri":    r.Header.Get("X-Forwarded-Proto"),
		"fwd_port":   r.Header.Get("X-Forwarded-Port"),
		// Specific to our kubernetes
		"fwd_edge":       r.Header.Get("X-Kubernauts-Edge"),
		"fwd_alpn":       r.Header.Get("X-Kubernauts-ALPN"),
		"fwd_proto":      r.Header.Get("X-Kubernauts-HTTP"),
		"fwd_tlsver":     r.Header.Get("X-Kubernauts-TLSVersion"),
		"fwd_tlscipher":  r.Header.Get("X-Kubernauts-TLSCipher"),
		"fwd_tls13early": r.Header.Get("X-Kubernauts-TLS13Early"),
	}
}

func (h *StaticHandler) Handler(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	w.Header().Add("Content-Type", "application/javascript")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET, OPTIONS")

	results := GetBaseInfo(r)
	results["name"] = h.clusterName
	_, err := fmt.Fprintf(w, "window.clusters[\"" + h.clusterName + "\"] = ")
	if err != nil {
		log.WithError(err).Warn("Error when writing back to client")
	}
	err = encoder.Encode(results)
	if err != nil {
		log.WithError(err).Warn("Error in response serializer")
	}
}
