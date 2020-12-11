package main

import (
	"context"
	"fmt"
	"log"

	"gitlab.com/microser/go-affise/affise"
)

func main() {
	client := affise.NewClient(
		affise.WithAPIKey("apikey"),
		affise.WithEndpoint("https://api-rocketcompany.affise.com"),
	)

	ctx := context.Background()
	offer, _, err := client.Offers.GetByID(ctx, 1)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println(offer)
}
