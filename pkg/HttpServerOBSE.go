package pkg

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func HttpServerOBSE(queueName1 string, queueName2 string) {
	var counter = 0

	ch1 := openConnectionAndChannel("guest:guest@localhost:5672")
	ch2 := openConnectionAndChannel("guest:guest@localhost:5672")
	var filename = "/tmp/log.txt"
	_, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	for true {
		var item1 = consumeEvent(ch1, queueName1)
		counter += 1
		var count1 = counter
		message1 := fmt.Sprintf("%s %s %s to %s", time.Now().Format(time.RFC3339), strconv.Itoa(count1), item1, "compse140.o")
		log.Print(message1)
		writer(filename, message1)
		var item2 = consumeEvent(ch2, queueName2)
		counter += 1
		var count2 = counter
		message2 := fmt.Sprintf("%s %s %s to %s", time.Now().Format(time.RFC3339), strconv.Itoa(count2), item2, "compse140.i")
		writer(filename, message1)
		log.Print(message2)
	}
}
