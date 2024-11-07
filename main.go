package main

import (
	"log"
	"math"
	"time"

	services "github.com/dhruv-assessment/load-balancer/service"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	log.Println("Load Balancer Started")

	for {
		noOfAppTierEC2, err := services.GetNoOfAppTierEC2()
		if err != nil {
			log.Fatalf("Unable to get no. of app tier instances: %v", err)
		}

		noOfMsgReqQueue, err := services.GetNoOfMessagesInRequestQueue()
		if err != nil {
			log.Fatalf("Unable to get no. of messages in request queue: %v", err)
		}

		if noOfMsgReqQueue > 0 && noOfMsgReqQueue > noOfAppTierEC2 {
			temp1 := 20 - noOfAppTierEC2
			if temp1 > 0 {
				temp2 := noOfMsgReqQueue - noOfAppTierEC2
				minNoOfIntances := math.Min(float64(temp1), float64(temp2))
				for i := 0; i < int(minNoOfIntances); i++ {
					services.CreateAppTierEC2(i + 1)
				}
			}
		}

		time.Sleep(time.Second * 1)
	}
}
