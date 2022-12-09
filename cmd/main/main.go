package main

import (
	"errors"
	"fmt"
	"mina.fi/devopstuni/pkg"
	"os"
)

func main() {
	fmt.Println("start")
	fmt.Println(os.Args[1:])

	switch os.Args[1:][0] {
	case "HttpServ":
		pkg.HttpServ()
	case "HttpServerORIG":
		pkg.HttpServerORIG("compse140.o", true, "guest1:guest1@localhost:5672")
	case "HttpServerOBSE":
		pkg.HttpServerOBSE("compse140.o-1", "compse140.i", "guest1:guest1@localhost:5672")
	case "HttpServerIMED":
		pkg.HttpServerIMED("compse140.i", "compse140.o-2", true, "guest1:guest1@localhost:5672")
	case "ApiGateWay":
		//pkg.ApiGateWay("compse140.i", "compse140.o")
	default:
		fmt.Println("unknown mode")
		errors.New("unknown mode")
	}

}
