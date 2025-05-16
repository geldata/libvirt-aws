package server

import (
	"fmt"
	"net/http"

	"github.com/geldata/libvirt-aws/awsapi"
)

func (s *AWSEmulatorServer) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		action := r.FormValue("Action")
		if action == "" {
			http.Error(w, "The action parameter is required", http.StatusBadRequest)
			return
		}
		switch action {
		case "DescribeAvailabilityZones":
			awsapi.DescribeAvailabilityZones(s.Region, w, r)
			return
		default:
			http.Error(w, fmt.Sprintf("The action %s is not valid for this web service.", action), http.StatusBadRequest)
			return
		}
	})
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
}
