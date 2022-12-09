package pkg

import (
	"fmt"
	"log"

	"strconv"
	"time"
)

func HttpServerOBSE(queueName1 string, queueName2 string, connectionString string) {
	var counter = 0

	ch1 := openConnectionAndChannel(connectionString)
	ch2 := openConnectionAndChannel(connectionString)
	/*var filename = "/tmp/log.txt"
	_, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}*/

	for true {
		var item1 = consumeEvent(ch1, queueName1)
		counter += 1
		var count1 = counter
		message1 := fmt.Sprintf("%s %s %s from %s", time.Now().Format(time.RFC3339), strconv.Itoa(count1), item1, queueName1)
		log.Print(queueName1)
		log.Print(message1)
		//writer(filename, message1)
		var item2 = consumeEvent(ch2, queueName2)
		counter += 1
		var count2 = counter
		message2 := fmt.Sprintf("%s %s %s from %s", time.Now().Format(time.RFC3339), strconv.Itoa(count2), item2, queueName2)
		//writer(filename, message1)
		log.Print(queueName2)
		log.Print(message2)
	}
}
