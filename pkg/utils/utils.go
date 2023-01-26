package utils

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func envOrString(env string, parameter string) string {
	if len(env) != 0 {
		return env
	} else {
		return parameter
	}
}

var Properties = struct {
	QUEUECREATED        string
	MAINQUEUE1          string
	MAINQUEUE2          string
	EXCHANGES           string
	FANOUTQUEUE1        string
	FANOUTQUEUE2        string
	RABBITMQ            string
	RABBITMQURL         string
	RABBITMQTOKEN       string
	MONGODBURL          string
	DBNAME              string
	PORT                string
	OBSERVLOGCOLLECTION string
	STATECOLLECTION     string
	HTTPSERV            string
	HTTPORIGV           string
	HTTPOBSE            string
	HTTPIMED            string
	HTTSERVICESTATS     string
}{
	QUEUECREATED:        envOrString(os.Getenv("queueCreated"), "true"),
	MAINQUEUE2:          envOrString(os.Getenv("mainqueue2"), "compse140.i"),
	EXCHANGES:           envOrString(os.Getenv("exchanges"), "compse140.o-ex"),
	FANOUTQUEUE1:        envOrString(os.Getenv("fanoutqueue1"), "compse140.o"),
	FANOUTQUEUE2:        envOrString(os.Getenv("fanoutqueue2"), "compse140.o-2"),
	RABBITMQ:            envOrString(os.Getenv("rabbitmq"), "guest1:guest1@localhost:5672"),
	RABBITMQURL:         envOrString(os.Getenv("rabbitmqUrl"), "localhost:15672"),
	RABBITMQTOKEN:       envOrString(os.Getenv("rabbitmqToken"), "Z3Vlc3QxOmd1ZXN0MQ=="),
	MONGODBURL:          envOrString(os.Getenv("mongoDbURL"), "mongodb://localhost:27017/test"),
	DBNAME:              envOrString(os.Getenv("dbName"), "test"),
	PORT:                envOrString(os.Getenv("port"), "3333"),
	OBSERVLOGCOLLECTION: "observLogCollection",
	STATECOLLECTION:     "stateLogCollection",
	HTTPSERV:            envOrString(os.Getenv("httpServ"), "localhost:3333"),
	HTTPORIGV:           envOrString(os.Getenv("httpOrigv"), "localhost:3334"),
	HTTPOBSE:            envOrString(os.Getenv("httpObse"), "localhost:3335"),
	HTTPIMED:            envOrString(os.Getenv("httpImed"), "localhost:3336"),
	HTTSERVICESTATS:     envOrString(os.Getenv("httServiceStats"), "localhost:3337"),
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func CreateLogEntry(state string) string {
	return fmt.Sprintf("%s: %s", time.Now().Format(time.RFC3339), state)
}

func StartServer(port string) {
	err := http.ListenAndServe(":"+port, nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

const (
	GET     string = "GET"
	HEAD    string = "HEAD"
	POST    string = "POST"
	DELETE  string = "DELETE"
	CONNECT string = "CONNECT"
	OPTIONS string = "OPTIONS"
	TRACE   string = "TRACE"
	PATCH   string = "TRACE"
	PUT     string = "PUT"
)

const (
	Init     string = "INIT"
	Paused   string = "PAUSED"
	Running  string = "RUNNING"
	Shutdown string = "SHUTDOWN"
)

var States = struct {
	Init     string
	Paused   string
	Running  string
	Shutdown string
}{
	Init:     Init,
	Paused:   Paused,
	Running:  Running,
	Shutdown: Shutdown,
}
