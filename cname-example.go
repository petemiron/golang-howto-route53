package main

import (
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

var name string
var target string
var TTL int64
var weight = int64(1)
var zoneId string

func init() {
	flag.StringVar(&name, "d", "", "domain name")
	flag.StringVar(&target, "t", "", "target of domain name")
	flag.StringVar(&zoneId, "z", "", "AWS Zone Id for domain")
	flag.Int64Var(&TTL, "ttl", int64(60), "TTL for DNS Cache")

}

func main() {
	flag.Parse()
	if name == "" || target == "" || zoneId == "" {
		fmt.Println(fmt.Errorf("Incomplete arguments: d: %s, t: %s, z: %s\n", name, target, zoneId))
		flag.PrintDefaults()
		return
	}
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := route53.New(sess)

	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{ // Required
			Changes: []*route53.Change{ // Required
				{ // Required
					Action: aws.String("UPSERT"), // Required
					ResourceRecordSet: &route53.ResourceRecordSet{ // Required
						Name: aws.String(name),    // Required
						Type: aws.String("CNAME"), // Required
						ResourceRecords: []*route53.ResourceRecord{
							{ // Required
								Value: aws.String(target), // Required
							},
						},
						TTL:           aws.Int64(TTL),
						Weight:        aws.Int64(weight),
						SetIdentifier: aws.String("Arbitrary Id describing this change set"),
					},
				},
			},
			Comment: aws.String("Sample update."),
		},
		HostedZoneId: aws.String(zoneId), // Required
	}
	resp, err := svc.ChangeResourceRecordSets(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}
