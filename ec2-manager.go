package main

import (
	"os"
	"strings"
	"flag"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/Sirupsen/logrus"
)

type Filters struct {
	region string
	tag    string
	values []string
}

var log = logrus.New()

func main() {
	log.Out = os.Stdout
	action := flag.String("action", "", "Action to execute: start|stop")
	region := flag.String("region", "", "Region to look for instances")
	tag := flag.String("tag", "", "Tag to filter instances with")
	flag.Parse()

	filters := Filters{*region, *tag, flag.Args()}



	ignoreStates, actionFn := buildAction(action)
	if actionFn == nil {
		log.WithFields(logrus.Fields{
			"action": *action,
		}).Error("Invalid action. Only 'start' or 'stop' are supported")
		os.Exit(1)
	}
	checkArgs(filters)

	log.WithFields(logrus.Fields{
		"action": *action,
		"region": filters.region,
		"tag":   filters.tag,
		"values":   filters.values,
	}).Info("Starting ec2-manager")

	sess, err := session.NewSession()
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error creating session in AWS")
		panic(err)
	}

	svc := ec2.New(sess, &aws.Config{Region: aws.String(filters.region)})

	var nameFilter []*string
	for _, f := range flag.Args() {
		nameFilter = append(nameFilter, aws.String(strings.Join([]string{"*", f, "*"}, "")))
	}

	log.WithFields(logrus.Fields{
		"region": filters.region,
		"tag":   filters.tag,
	}).Info("Looking for instances. Filter:")

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:" + filters.tag),
				Values: nameFilter,
			},
		},
	}

	resp, err := svc.DescribeInstances(params)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error looking for instances")
		panic(err)
	}

	log.Info("Number of instances found: ", len(resp.Reservations))
	instanceIds := []*string{}
	for idx, _ := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			if !ignoreStates[*inst.State.Name] {
				instanceIds = append(instanceIds, inst.InstanceId)
			} else {
				log.Info("Ignoring instance '", *inst.InstanceId, "'. It's currently in state '", *inst.State.Name, "'")
			}
		}
	}

	if len(instanceIds) > 0 {
		actionFn(svc, instanceIds)
	}
	log.Info("Finishing ec2-manager. Number of instances processed: ", len(instanceIds))
}

func buildAction(action *string) (map[string]bool, func(*ec2.EC2, []*string) error) {
	switch *action {
	case "start":
		ignoreStates := map[string]bool{
			"terminated": true,
			"pending": true,
			"running": true,
		}
		return ignoreStates, func(svc *ec2.EC2, instanceIds []*string) error {
			resp, err := svc.StartInstances(
				&ec2.StartInstancesInput{
					InstanceIds: instanceIds,
				},
			)
			log.WithFields(logrus.Fields{
				"response": resp,
			}).Info	("Response starting instances")
			return err
		}
	case "stop":
		ignoreStates := map[string]bool{
			"terminated": true,
			"shutting-down": true,
			"stopping": true,
			"stopped": true,
		}
		return ignoreStates, func(svc *ec2.EC2, instanceIds []*string) error {
			resp, err := svc.StopInstances(
				&ec2.StopInstancesInput{
					InstanceIds: instanceIds,
				},
			)
			log.WithFields(logrus.Fields{
				"response": resp,
			}).Info("Response stopping instances")
			return err
		}
	default:
		return nil, nil
	}
}

func checkArgs(filters Filters) {
	if filters.region == "" ||
		filters.tag == "" || len(filters.values) == 0 {
		log.WithFields(logrus.Fields{
			"region": filters.region,
			"tag":   filters.tag,
			"":   filters.values,
		}).Error("Invalid args.")
		os.Exit(1)
	}
}
