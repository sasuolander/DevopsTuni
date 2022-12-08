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

func updateStates() {
	/*PAUSED = ORIG service is not sending messages
	RUNNING = ORIG service sends messages
	If the new state is equal to previous nothing happens.
		There are two special cases:
	INIT = everything (except log information for /run-log and /messages) is in the
	initial state and ORIG starts sending again,
		state is set to RUNNING
	SHUTDOWN = all containers are stopped*/
}

func pollStates() {
	/*PAUSED = ORIG service is not sending messages
	RUNNING = ORIG service sends messages
	If the new state is equal to previous nothing happens.
		There are two special cases:
	INIT = everything (except log information for /run-log and /messages) is in the
	initial state and ORIG starts sending again,
		state is set to RUNNING
	SHUTDOWN = all containers are stopped*/

	// keep log of state change
}

func ApiGateWay(routes []Route) {

	http.HandleFunc("/message", forwardRequest)
	http.HandleFunc("/state", forwardRequest)
	http.HandleFunc("/state", forwardRequest) // PUT (payload “INIT”, “PAUSED”, “RUNNING”, “SHUTDOWN”)
	http.HandleFunc("/run-log", forwardRequest)
	http.HandleFunc("/message", forwardRequest)
	http.HandleFunc("/node-statistic ", forwardRequest)
	http.HandleFunc("/queue-statistic", forwardRequest)

	err := http.ListenAndServe(":8083", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
