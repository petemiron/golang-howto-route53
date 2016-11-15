# Overview
The goal of this how-to is to demonstrate adding a CNAME to a domain in Amazon's Route53 service.

# Prerequisites

## Create Hosted Zone
Setup a new hosted zone by going to the [Route53 home in your AWS account](https://console.aws.amazon.com/route53/)

Assuming you're registering a new domain, you'll get an email confirming your email address. Follow the steps, 
Amazon will notify you once it's ready.

Once you have your domain registered, you'll get a Hosted Zone Id. Keep it ready, you'll need it in a bit.

## Setup your AWS Credentials
The [AWS go-sdk](http://docs.aws.amazon.com/sdk-for-go/api/service/route53/) requires you to setup the credentials
in ~/.aws/credentials. 


### ~/.aws/credentials
```
[default]
aws_access_key_id = [your access key id]
aws_secret_access_key = [your secret access key] 
```

# Create your CNAME
This is actually a pretty straightforward common task, but I haven't found good documentation on it, so here goes.

I've stripped down the AWS sdk examples to just the bare minimum for setting up a CNAME record.

I'm using the UPSERT Action, which will automatically insert or update. Name is the domain name you're adding the CNAME
for. Target is the destination of the CNAME record. You can also add a TTL.

``` go
func createCNAME(svc *route53.Route53) {
...
	params := &route53.ChangeResourceRecordSetsInput{
	    ChangeBatch: &route53.ChangeBatch{ // Required
	        Changes: []*route53.Change{ // Required
	            { // Required
	                Action: aws.String("UPSERT"), // Required
	                ResourceRecordSet: &route53.ResourceRecordSet{ // Required
	                    Name: aws.String(name), // Required
	                    Type: aws.String("CNAME"),  // Required
	                    ResourceRecords: []*route53.ResourceRecord{
	                        { // Required
	                            Value: aws.String(target), // Required
	                        },
	                    },
	                    TTL:            aws.Int64(TTL),
	                    Weight:         aws.Int64(weight),
					    SetIdentifier:  aws.String("Arbitrary Id describing this change set"),
	                },
	            },
	        },
	        Comment: aws.String("Sample update."),
	    },
	    HostedZoneId: aws.String(zoneId), // Required
	}
	resp, err := svc.ChangeResourceRecordSets(params)
...
```

# List all domains
I can't figure out how to get this query to constrain results. If you've got a suggestion, let me know.

``` go
func listCNAMES(svc *route53.Route53) {
...
	listParams := &route53.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(zoneId), // Required
		...
	}
	respList, err := svc.ListResourceRecordSets(listParams)
...
}
```

# Run the example

``` bash
go run cname-example.go -d www.your-domain.com -t your-domain.com -z [Your Zone Id]
```


