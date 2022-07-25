package main

import (
	log "github.com/sirupsen/logrus"
	"gitlab.gwdg.de/fe/digis/database-api/src/api"
)

func main() {
	echoAPI := api.InitializeAPI()
	log.Fatal(echoAPI.Start(":80"))
}
