package pkg

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type Route struct {
	externalEndpoint Endpoint
	internalEndpoint Endpoint
}

type Endpoint struct {
	url url.URL
}

// make sure http connection is kept alive so we can return response also using same connection
func forwardRequest(w http.ResponseWriter, r *http.Request) {

}

func ApiGateWay(routes []Route) {

	http.HandleFunc("/message", forwardRequest)
	http.HandleFunc("/state", forwardRequest)
	http.HandleFunc("/state", forwardRequest) // PUT (payload “INIT”, “PAUSED”, “RUNNING”, “SHUTDOWN”)
	http.HandleFunc("/run-log", forwardRequest)
	http.HandleFunc("/message", forwardRequest)

	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
