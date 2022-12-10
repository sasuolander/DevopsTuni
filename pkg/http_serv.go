package pkg

import (
	"context"
	"errors"
	"fmt"
	"mina.fi/devopstuni/pkg/externalConnection/database"
	"net/http"
	"os"
	"strings"
)

func HttpServ() {

	http.HandleFunc("/", getLogResponse)

	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func getLogResponse(w http.ResponseWriter, r *http.Request) {
	var client = database.Connection(Properties()["mongoDbURL"])
	var collectionName = "observLogCollection"
	var result = database.GetItems(client, Properties()["dbName"], collectionName)
	var stringArray []string

	for _, message := range result {
		stringArray = append(stringArray, message.Message+" \n")
	}

	fmt.Print(result)
	stringByte := strings.Join(stringArray, "\x20\x00")
	w.Header().Add("Content-Type", "text/plain;charset=utf-8")
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte(stringByte))
	if err != nil {
		fmt.Printf("error in server: %s\n", err)
		return
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			FailOnError(err, "Failed to disconnect")
			panic(err)
		}
	}()
	return
}
