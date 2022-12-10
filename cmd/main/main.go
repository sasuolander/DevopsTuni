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
	case "ORIGMain":
		pkg.ORIGMain()
	case "OBSEMain":
		pkg.OBSEMain()
	case "IMEDMain":
		pkg.IMEDMain()
	case "ApiGateWay":
		//pkg.ApiGateWay("compse140.i", "compse140.o")
	default:
		fmt.Println("unknown mode")
		errors.New("unknown mode")
	}

}
