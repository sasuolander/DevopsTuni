package pkg

import (
	"context"
	"fmt"
	"mina.fi/devopstuni/pkg/externalConnection/database"
	"mina.fi/devopstuni/pkg/utils"
	"net/http"
	"strings"
	"time"
)

func ServerStarterServ() {
	startTime := fmt.Sprintf("%s", time.Now().Format(time.RFC3339))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == utils.GET && r.URL.Path == "/ping" {
			pigApi(w, r, startTime)
		}

		if r.URL.Path == "/" {
			getLogResponse(w, r)
		}

	})
	utils.StartServer(utils.Properties.PORT)
}

func getLogResponse(w http.ResponseWriter, r *http.Request) {
	var client = database.Connection(utils.Properties.MONGODBURL)
	var collectionName = utils.Properties.OBSERVLOGCOLLECTION
	var result = database.GetItems(client, utils.Properties.DBNAME, collectionName)

	w.Header().Add("Content-Type", "text/plain;charset=utf-8")
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte(strings.Join(database.ConvertToStringArray(result), "\x20\x00")))
	utils.FailOnError(err, "Error in server")

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			utils.FailOnError(err, "Failed to disconnect")
		}
	}()
	return
}
