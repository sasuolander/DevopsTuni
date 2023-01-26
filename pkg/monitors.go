package pkg

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"mina.fi/devopstuni/pkg/json"
	"mina.fi/devopstuni/pkg/utils"
	"net/http"
)

type RoutePing struct {
	endpoint EndpointPing
}
type EndpointPing struct {
	name string
	url  string
}

func ServerStarterMonitor() {
	http.HandleFunc("/", requestPingApi)
	utils.StartServer(utils.Properties.PORT)
}

func buildRoutesPing() []RoutePing {
	var httpOrigv = utils.Properties.HTTPORIGV
	var httpServ = utils.Properties.HTTPSERV
	var httpObse = utils.Properties.HTTPOBSE
	var httpImed = utils.Properties.HTTPIMED
	var httpDB = utils.Properties.MONGODBURL

	var routes = []RoutePing{
		{
			endpoint: EndpointPing{name: "ORIG", url: httpOrigv},
		},
		{
			endpoint: EndpointPing{name: "SERV", url: httpServ},
		},
		{
			endpoint: EndpointPing{name: "OBSE", url: httpObse},
		},
		{
			endpoint: EndpointPing{name: "IMED", url: httpImed},
		},
		{
			endpoint: EndpointPing{name: "DB", url: httpDB},
		},
	}

	return routes
}

func makeRequestPing(route string) (PingStatusOutput, bool) {
	client := &http.Client{}

	req, err := http.NewRequest(utils.GET, route+"/ping", nil)
	utils.FailOnError(err, "failed to create request")

	resp, err := client.Do(req)
	if resp == nil {
		return PingStatusOutput{}, false
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		println(resp.Status)
		return PingStatusOutput{}, false
	}
	body, err := io.ReadAll(resp.Body)
	var converted PingStatusOutput
	json.UnMarshallRequest(body, &converted)
	return converted, true
}

type PingOk struct {
	Service  string `json:"service"`
	Running  bool   `json:"running"`
	Startime string `json:"startime"`
}

func requestPingApi(w http.ResponseWriter, r *http.Request) {

	var pingOk []PingOk
	var routes = buildRoutesPing()

	for _, route := range routes {
		var resp PingStatusOutput
		var statusRequest bool
		if route.endpoint.name != "DB" {
			resp, statusRequest = makeRequestPing(route.endpoint.url)
		} else {
			clientOptions := options.Client().ApplyURI(route.endpoint.url)
			client, err := mongo.Connect(context.TODO(), clientOptions)
			if err != nil {
				log.Fatal(err)
			}
			// Check the connection
			err = client.Ping(context.TODO(), nil)
			if err != nil {
				statusRequest = false
			} else {
				statusRequest = true
				resp.Serverstarttime = ""
				resp.Info = "OK"
			}
		}

		if statusRequest == true {
			if resp.Info == "OK" {
				pingOk = append(pingOk, PingOk{Service: route.endpoint.name, Running: true, Startime: resp.Serverstarttime})
			} else {
				pingOk = append(pingOk, PingOk{Service: route.endpoint.name, Running: false, Startime: ""})
			}
		} else {
			pingOk = append(pingOk, PingOk{Service: route.endpoint.name, Running: false, Startime: ""})
		}
	}
	jsonData := json.MarshallRequest(pingOk)
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(jsonData)
	utils.FailOnError(err, "failed to write response")
}
