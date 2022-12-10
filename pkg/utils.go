package pkg

import (
	"log"
	"os"
)

func envOrString(env string, parameter string) string {
	if len(env) != 0 {
		return env
	} else {
		return parameter
	}
}

func Properties() map[string]string {
	properties := make(map[string]string)
	properties["mainqueue1"] = envOrString(os.Getenv("mainqueue1"), "compse140.o")
	properties["mainqueue2"] = envOrString(os.Getenv("mainqueue2"), "compse140.i")
	properties["rabbitmq"] = envOrString(os.Getenv("rabbitmq"), "guest1:guest1@localhost:5672")
	properties["fanoutqueue1"] = envOrString(os.Getenv("fanoutqueue1"), "compse140.o-1")
	properties["fanoutqueue2"] = envOrString(os.Getenv("fanoutqueue2"), "compse140.o-2")
	properties["mongoDbURL"] = envOrString(os.Getenv("mongoDbURL"), "mongodb://localhost:27017/test")
	properties["dbName"] = envOrString(os.Getenv("dbName"), "test")
	properties["exchanges"] = envOrString(os.Getenv("exchanges"), "compse140.o-ex")
	properties["queueCreated"] = envOrString(os.Getenv("queueCreated"), "true")
	return properties
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
