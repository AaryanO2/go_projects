package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

type resp struct {
	Domain      string `json:"domain"`
	HasMX       bool   `json:"hasMX"`
	HasSPF      bool   `json:"hasSPF"`
	SpfRecord   string `json:"spfRecord"`
	HasDMARC    bool   `json:"hasDMARC"`
	DmarcRecord string `json:"dmarcRecord"`
}

type req struct {
	Domain string `json:"domain"`
}

func main() {
	http.HandleFunc("/", checker)
	fmt.Println("Serving HTTP on port 9010....")
	err := http.ListenAndServe("0.0.0.0:9010", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v\n", err)
	}
}

func checker(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error parsing body.", http.StatusBadRequest)
		return
	}

	var x req

	err = json.Unmarshal(body, &x)
	if err != nil {
		http.Error(w, "Error unmarshalling body.", http.StatusBadRequest)
		return
	}

	response := checkDomain(x.Domain)
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatalf("Error encoding response: %v\n", err)
		http.Error(w, "Error encoding response.", http.StatusInternalServerError)
		return
	}
}

func checkDomain(domain string) resp {
	var response resp
	response.Domain = domain

	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else if len(mxRecords) > 0 {
		response.HasMX = true
	}

	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		for _, record := range txtRecords {
			if strings.HasPrefix(record, "v=spf1") {
				response.HasSPF = true
				response.SpfRecord = record
				break
			}
		}
	}

	dmarc_records, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		for _, record := range dmarc_records {
			if strings.HasPrefix(record, "v=DMARC1") {
				response.HasDMARC = true
				response.DmarcRecord = record
				break
			}
		}
	}

	fmt.Printf("Response: %v\n", response)
	return response
}
