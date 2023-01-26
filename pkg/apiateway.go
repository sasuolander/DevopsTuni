package pkg

import (
	"context"
	_ "encoding/json"
	"fmt"
	"io"
	"log"
	"mina.fi/devopstuni/pkg/externalConnection/database"
	"mina.fi/devopstuni/pkg/externalConnection/rabbitmq"
	"mina.fi/devopstuni/pkg/json"
	"mina.fi/devopstuni/pkg/utils"
	"net/http"
	"strings"
)

type Route struct {
	externalEndpoint Endpoint
	internalEndpoint Endpoint
	forward          bool // If true forward into internal endpoint by using methodForward, otherwise use method
	method           func(w http.ResponseWriter, r *http.Request)
	methodForward    func(w http.ResponseWriter, r *http.Request, destination Route)
}

type Endpoint struct {
	url    string
	method string
}

type QueueStatistic struct {
	Name                      string
	DeliveryRate              float64 `json:"delivery_rate"`
	MessagesDelivered         int     `json:"messages_delivered"`
	MessagesDeliveredRecently int     `json:"messages_delivered_recently"`
	MessagesPublished         int     `json:"messages_published"`
	MessagesPublishedRecently int     `json:"messages_published_recently"`
}

type NodeInfo struct {
	DiskFree   int    `json:"disk_free"`
	MemUsed    int    `json:"mem_used"`
	Name       string `json:"name"`
	Processors int    `json:"processors"`
	RunQueue   int    `json:"run_queue"`
}

func buildRoutes() []Route {
	var httpRabbitMq = utils.Properties.RABBITMQURL
	var httpOrigv = utils.Properties.HTTPORIGV
	var httpServ = utils.Properties.HTTPSERV
	var httServiceStats = utils.Properties.HTTSERVICESTATS

	var routes = []Route{

		{externalEndpoint: Endpoint{
			url:    "/message",
			method: utils.GET,
		},
			internalEndpoint: Endpoint{
				url:    httpServ + "/",
				method: utils.GET,
			},
			forward:       true,
			methodForward: forwarder,
		},
		{externalEndpoint: Endpoint{
			url:    "/state",
			method: utils.GET,
		},
			internalEndpoint: Endpoint{
				url:    httpOrigv + "/state",
				method: utils.GET,
			},
			forward:       true,
			methodForward: forwarder,
		},
		{externalEndpoint: Endpoint{
			url:    "/state",
			method: utils.PUT,
		},
			internalEndpoint: Endpoint{
				url:    httpOrigv + "/state",
				method: utils.PUT,
			},
			forward:       true,
			methodForward: forwarder,
		},
		{externalEndpoint: Endpoint{
			url:    "/run-log",
			method: utils.GET,
		},
			internalEndpoint: Endpoint{
				url:    "",
				method: utils.GET,
			},
			forward: false,
			method:  getLogs,
		},
		{externalEndpoint: Endpoint{
			url:    "/node-statistic",
			method: utils.GET,
		},
			internalEndpoint: Endpoint{
				url:    httpRabbitMq,
				method: utils.GET,
			},
			methodForward: forwarderNodeStats,
			forward:       true,
		},
		{externalEndpoint: Endpoint{
			url:    "/queue-statistic",
			method: utils.GET,
		},
			internalEndpoint: Endpoint{
				url:    httpRabbitMq,
				method: utils.GET,
			},
			methodForward: forwarderQueueStats,
			forward:       true,
		},
		{externalEndpoint: Endpoint{
			url:    "/serviceRunning",
			method: utils.GET,
		},
			internalEndpoint: Endpoint{
				url:    httServiceStats + "/",
				method: utils.GET,
			},
			forward:       true,
			methodForward: forwarder,
		},
	}
	return routes
}

func ServerStarterApiGateWay() {
	http.HandleFunc("/", handleRequest)
	utils.StartServer(utils.Properties.PORT)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var routes = buildRoutes()
	var rURL = r.URL.Path
	var rMethod = r.Method
	println("rURL: " + rURL)
	println("method: " + rMethod)
	for _, route := range routes { // router loop
		if rURL == route.externalEndpoint.url &&
			rMethod == route.externalEndpoint.method {
			if route.forward == true {
				route.methodForward(w, r, route)
			} else {
				route.method(w, r)
			}
		} else {
			//println("No route")
		}
	}
}

func forwarder(w http.ResponseWriter, r *http.Request, route Route) {
	var method = route.internalEndpoint.method
	var url = route.internalEndpoint.url
	requestForward(url, method, r.Body, w)

}

func forwarderQueueStats(w http.ResponseWriter, r *http.Request, route Route) {
	var queueStats []rabbitmq.Queue
	json.UnMarshallRequest(request("http://"+route.internalEndpoint.url+"/api/queues"), &queueStats)
	var queueData []QueueStatistic
	for _, stat := range queueStats {
		queueData = append(queueData, QueueStatistic{
			Name:                      stat.Name,
			MessagesPublishedRecently: stat.MessageStats.Publish,
			MessagesPublished:         stat.Messages,
			MessagesDelivered:         stat.MessageStats.Deliver,
			MessagesDeliveredRecently: stat.MessageStats.DeliverGet,
			DeliveryRate:              stat.MessageStats.PublishDetails.Rate,
		})
	}

	jsonData := json.MarshallRequest(queueData)
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(jsonData)

	utils.FailOnError(err, "failed to write response")
	return

}

func forwarderNodeStats(w http.ResponseWriter, r *http.Request, route Route) {
	var nodeInfo []NodeInfo
	json.UnMarshallRequest(request("http://"+route.internalEndpoint.url+"/api/nodes"), &nodeInfo)
	jsonData := json.MarshallRequest(nodeInfo)
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(jsonData)
	utils.FailOnError(err, "failed to write response")
	return
}

func requestForward(urlString string, method string, bodyInput io.Reader, w http.ResponseWriter) {
	client := &http.Client{}
	req, err := http.NewRequest(method, urlString, bodyInput)
	utils.FailOnError(err, "failed to create request")

	resp, err := client.Do(req)
	utils.FailOnError(err, "failed to send request")

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Panicf("Error wit code: %s", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	utils.FailOnError(err, "failed to read response")
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	utils.FailOnError(err, "failed to write response")
}

func request(urlString string) []byte {
	client := &http.Client{}

	req, err := http.NewRequest("GET", urlString, nil)
	utils.FailOnError(err, "failed to create request")

	req.Header.Set("Authorization", "Basic "+utils.Properties.RABBITMQTOKEN)

	resp, err := client.Do(req)

	utils.FailOnError(err, "failed to send request")

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		println(resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	utils.FailOnError(err, "failed to read response")

	return body
}

func getLogs(w http.ResponseWriter, r *http.Request) {
	var client = database.Connection(utils.Properties.MONGODBURL)
	var result = database.GetItems(client, utils.Properties.DBNAME, utils.Properties.STATECOLLECTION)
	w.Header().Add("Content-Type", "text/plain;charset=utf-8")
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte(strings.Join(database.ConvertToStringArray(result), "\x20\x00")))
	utils.FailOnError(err, fmt.Sprintf("error in server: %s\n", err))

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			utils.FailOnError(err, "Failed to disconnect")
		}
	}()
	return

}
