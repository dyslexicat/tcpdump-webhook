package main

import (
  "encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"time"
  "github.com/dyslexicat/tcpdump-webhook/pkg/mutate"
  admissionv1 "k8s.io/api/admission/v1"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "hello %q", html.EscapeString(r.URL.Path))
}

func handleMutate(w http.ResponseWriter, r *http.Request) {
  admReview, err := mutate.AdmissionReviewFromRequest(r)

  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    log.Printf("error getting admission review from request: %v", err)
    return
  }

  admResp, err := mutate.AdmissionResponseFromReview(admReview)

  if err != nil {
    w.WriteHeader(400)
		w.Write([]byte(err.Error()))
    return
  }

  // the final response will be another admission review
  var admissionReviewResponse admissionv1.AdmissionReview
  admissionReviewResponse.Response = admResp
  admissionReviewResponse.SetGroupVersionKind(admReview.GroupVersionKind())
  admissionReviewResponse.Response.UID = admReview.Request.UID
   
  resp, err := json.Marshal(admissionReviewResponse)
  if err != nil {
    msg := fmt.Errorf("error marshaling response: %v", err)
    log.Println(msg)
    w.WriteHeader(500)
		w.Write([]byte(msg.Error()))
		return
  }

  w.Header().Set("Content-Type", "application/json")
  log.Printf("allowing pod as %v", string(resp))
	w.Write(resp)
}

func main() {
  log.Println("starting server...")

  mux := http.NewServeMux()

  mux.HandleFunc("/", handleRoot)
  mux.HandleFunc("/mutate", handleMutate)

  s := &http.Server{
    Addr:   ":8443",
    Handler: mux,
    ReadTimeout: 10 * time.Second,
    WriteTimeout: 10 * time.Second,
    MaxHeaderBytes: 1 << 20, // 1048576
  }

  log.Fatal(s.ListenAndServeTLS("./ssl/tcpdump.pem", "./ssl/tcpdump.key"))
}
