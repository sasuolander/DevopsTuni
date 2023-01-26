package pkg

import (
	"mina.fi/devopstuni/pkg/json"
	"mina.fi/devopstuni/pkg/utils"
	"net/http"
)

type PingStatusOutput struct {
	Info            string `json:"info"`
	Serverstarttime string `json:"serverstarttime"`
}

func pigApi(w http.ResponseWriter, r *http.Request, time string) {

	if r.Method == utils.GET && r.URL.Path == "/ping" {
		w.Header().Add("Content-Type", "text/plain;charset=utf-8")
		message := PingStatusOutput{
			Info:            "OK",
			Serverstarttime: time,
		}

		_, err := w.Write(json.MarshallRequest(message))
		utils.FailOnError(err, "failed to write response")
	}
	return
}
