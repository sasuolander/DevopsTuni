package pkg

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

func HttpServ(databaseString string, database string) {

	http.HandleFunc("/", getRoot)

	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
