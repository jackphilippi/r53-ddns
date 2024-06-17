package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/aws/smithy-go"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("usage: r53 <zone_id> <hostname> <ip>")
		fmt.Println("got: " + os.Args[1] + " " + os.Args[2] + " " + os.Args[3])
		os.Exit(64)
	}

	fmt.Println("Setting DDNS for " + os.Args[2] + " to " + os.Args[3] + " at " + time.Now().Format(time.RFC3339))

	cfg, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config, %v", err))
	}

	svc := route53.NewFromConfig(cfg)

	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action: types.ChangeActionUpsert,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: aws.String(os.Args[2]),
						Type: types.RRTypeA,
						ResourceRecords: []types.ResourceRecord{
							{
								Value: aws.String(os.Args[3]),
							},
						},
						TTL: aws.Int64(300),
					},
				},
			},
		},
		HostedZoneId: aws.String(os.Args[1]),
	}

	result, err := svc.ChangeResourceRecordSets(context.Background(), input)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			log.Printf("code: %s, message: %s, fault: %s", ae.ErrorCode(), ae.ErrorMessage(), ae.ErrorFault().String())
		} else {
			fmt.Fprintf(os.Stderr, err.Error())
		}

		os.Exit(1)
	}

	fmt.Println(result)
}
