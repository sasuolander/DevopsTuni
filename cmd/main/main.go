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
		pkg.ServerStarterServ()
	case "ORIGMain":
		pkg.ServerStarterORIG()
	case "OBSEMain":
		pkg.ServerStarterObse()
	case "IMEDMain":
		pkg.ServerStarterImed()
	case "ApiGateWay":
		pkg.ServerStarterApiGateWay()
	case "Monitor":
		pkg.ServerStarterMonitor()

	default:
		fmt.Println("unknown mode")
		errors.New("unknown mode")
	}

}
