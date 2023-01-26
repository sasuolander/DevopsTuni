package pkg

import (
	"context"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"log"
	"mina.fi/devopstuni/pkg/externalConnection/database"
	"mina.fi/devopstuni/pkg/externalConnection/rabbitmq"
	"mina.fi/devopstuni/pkg/utils"
	"net/http"
	"os"
	"strconv"
	"time"
)

type ServerState struct {
	Status string
}

type DesiredState struct {
	Status string
}

var stateGlobal ServerState
var desiredState DesiredState

func origMain() {
	var client = database.Connection(utils.Properties.MONGODBURL)
	ch := rabbitmq.OpenConnectionAndChannel(utils.Properties.RABBITMQ)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var exchanges = utils.Properties.EXCHANGES
	var condition, _ = strconv.ParseBool(utils.Properties.QUEUECREATED)
	if !condition {
		rabbitmq.CreateExchange(exchanges, ch)
		queue1 := rabbitmq.CreateQueue(utils.Properties.FANOUTQUEUE1, ch)
		queue2 := rabbitmq.CreateQueue(utils.Properties.FANOUTQUEUE2, ch)
		rabbitmq.ExchangeBindingToQueue(ch, exchanges, queue1.Name)
		rabbitmq.ExchangeBindingToQueue(ch, exchanges, queue2.Name)
		log.Print(queue1.Name)
		log.Print(queue2.Name)
	}

	var i = 1

	for true { // infinity event loop

		if stateGlobal.Status == utils.States.Init {
			stateGlobal.Status = utils.States.Running
			insert(client, utils.States.Running)
		}

		if desiredState.Status == utils.States.Paused {
			if stateGlobal.Status != desiredState.Status {
				stateGlobal.Status = utils.States.Paused
				insert(client, utils.States.Paused)
			}
		} else if desiredState.Status == utils.States.Running {
			if stateGlobal.Status != desiredState.Status {
				stateGlobal.Status = utils.States.Running
				insert(client, utils.States.Running)
			}
			mainActon(i, ch, ctx)
			i++
		} else if desiredState.Status == utils.States.Shutdown {
			if stateGlobal.Status != desiredState.Status {
				insert(client, utils.States.Shutdown)
			}
			os.Exit(0)
		} else {
			mainActon(i, ch, ctx)
			i++
		}
		time.Sleep(3 * time.Second) // each iteration await 3 second before starting next iteration
	}
}

func insert(client *mongo.Client, status string) {
	database.InsertItem(client, utils.Properties.DBNAME, utils.CreateLogEntry(status), utils.Properties.STATECOLLECTION)
}

func mainActon(i int, ch *amqp091.Channel, ctx context.Context) {
	body := "MSG_" + strconv.Itoa(i)
	rabbitmq.PublicEventExchange(ch, ctx, utils.Properties.EXCHANGES, body)
	log.Printf(" [x] Sent %s", body)
	log.Print("sleep")
}

func httpServerORIGServer() {
	log.Print("starting server")

	startTime := fmt.Sprintf("%s", time.Now().Format(time.RFC3339))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/state" && r.Method == utils.PUT {
			bodyBytes, err := io.ReadAll(r.Body)
			utils.FailOnError(err, "Failed to disconnect")

			bodyString := string(bodyBytes)
			if !validateState(bodyString) {
				log.Panicf("invalid state: %s", bodyString)
				return
			}
			desiredState.Status = bodyString
		}

		if r.URL.Path == "/state" && r.Method == utils.GET {
			w.Header().Add("Content-Type", "text/plain;charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(stateGlobal.Status))
			utils.FailOnError(err, "Error in server")
		}

		if r.Method == utils.GET && r.URL.Path == "/ping" {
			pigApi(w, r, startTime)
		}

	})

	utils.StartServer(utils.Properties.PORT)
}

func validateState(state string) bool {
	fmt.Print("state: " + utils.States.Paused)
	if state == utils.States.Init {
		return true
	} else if state == utils.States.Running {
		return true
	} else if state == utils.States.Paused {

		return true
	} else if state == utils.States.Shutdown {
		return true
	} else {
		return false
	}
}

func ServerStarterORIG() {
	stateGlobal.Status = utils.States.Init
	var client = database.Connection(utils.Properties.MONGODBURL)
	database.InsertItem(client, utils.Properties.DBNAME, utils.CreateLogEntry(utils.States.Init), utils.Properties.STATECOLLECTION)

	go origMain()
	httpServerORIGServer()
}
