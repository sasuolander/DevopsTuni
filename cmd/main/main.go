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
		pkg.HttpServerORIG("compse140.o")
	case "HttpServerOBSE":
		pkg.HttpServerOBSE("compse140.o", "compse140.i")
	case "HttpServerIMED":
		pkg.HttpServerIMED("compse140.i", "compse140.o")
	case "ApiGateWay":
		pkg.ApiGateWay("compse140.i", "compse140.o")
	default:
		fmt.Println("unknown mode")
		errors.New("unknown mode")
	}

}
