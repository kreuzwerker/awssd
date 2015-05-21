package awssd

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/ec2"
	"github.com/awslabs/aws-sdk-go/service/route53"
)

const (
	GlobalRegion = "us-east-1"
)

type Client struct {
	config *Config
	ec2c   *ec2.EC2
	r53c   *route53.Route53
	zoneId string
}

func NewClient(config *Config) *Client {

	ec2Config := *aws.DefaultConfig
	r53Config := *aws.DefaultConfig

	ec2Config.Region = config.Region

	return &Client{
		config: config,
		ec2c:   ec2.New(&ec2Config),
		r53c:   route53.New(&r53Config),
	}
}

func (c *Client) Diff() (bool, error) {

	if c.zoneId == "" {

		zoneId, err := c.getZoneId()

		if err != nil {
			return false, err
		}

		if Debug {
			log.Printf("[ DEBUG ] got zone id %s", zoneId)
		}

		c.zoneId = zoneId

	}

	zoneMap, err := c.getZoneMap()

	if err != nil {
		return false, err
	}

	if Debug {
		log.Printf("[ DEBUG ] got zone map %v", zoneMap)
	}

	ec2Map, err := c.getEc2Map()

	if err != nil {
		return false, err
	}

	if Debug {
		log.Printf("[ DEBUG ] got ec2 map %v", ec2Map)
	}

	return ec2Map.Diff(zoneMap, c.upsertRecordSet)

}

func (c *Client) upsertRecordSet(name string, ip []IP) error {

	var rr []*route53.ResourceRecord

	for _, e := range ip {
		rr = append(rr, e.ToResourceRecord())
	}

	ch := route53.Change{
		Action: aws.String("UPSERT"),
		ResourceRecordSet: &route53.ResourceRecordSet{
			Name:            aws.String(name),
			ResourceRecords: rr,
			TTL:             aws.Long(c.config.TTL),
			Type:            aws.String("A"),
		},
	}

	if c.config.DryRun {
		log.Printf("[ DEBUG ] dry-run, skipping update of record set for %s with %q", name, ip)
		return nil
	}

	cb := &route53.ChangeBatch{
		Changes: []*route53.Change{&ch},
		Comment: aws.String(fmt.Sprintf("awssd-addRecordSet %s: %q", name, ip)),
	}

	log.Printf("[  INFO ] updating %s with %q", name, ip)

	resp, err := c.r53c.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
		ChangeBatch:  cb,
		HostedZoneID: aws.String(c.zoneId),
	})

	if err == nil && Debug {
		id := resp.ChangeInfo.ID
		log.Printf("[ DEBUG ] changeset with id %q created", *id)
	}

	return err

}

func (c *Client) getEc2Map() (Mapping, error) {

	filters, err := c.config.ConvertFilter()

	if err != nil {
		return nil, err
	}

	resp, err := c.ec2c.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: filters,
	})

	if err != nil {
		return nil, err
	}

	var m = NewMapping()

	for _, r := range resp.Reservations {

		for _, e := range r.Instances {

			if *e.State.Name != "running" {
				continue
			}

			var ip *string

			if c.config.PreferPrivate {
				ip = e.PrivateIPAddress
			} else {

				if e.PublicIPAddress != nil {
					ip = e.PublicIPAddress
				} else {
					ip = e.PrivateIPAddress
				}

			}

			if value := findTagValue(c.config.GroupBy, e.Tags); value == nil {
				log.Printf("[ ERROR ] no value for group-by tag key %q for instance %q, skipping", c.config.GroupBy, *e.InstanceID)
			} else {
				m.Add(fmt.Sprintf("%s.%s.", *value, c.config.Domain), *ip)
			}

		}

	}

	return m, nil

}

func (c *Client) getZoneId() (string, error) {

	resp, err := c.r53c.ListHostedZones(&route53.ListHostedZonesInput{})

	if err != nil {
		return "", err
	}

	exp := regexp.MustCompile(`^/hostedzone/(.+)$`)

	for _, e := range resp.HostedZones {

		if *e.Name == fmt.Sprintf("%s.", c.config.Domain) {

			result := exp.FindStringSubmatch(*e.ID)

			if len(result) != 2 {
				log.Fatalf("cannot extract hosted zone id from ID '%v'", *e.ID)
			}

			return result[1], nil

		}

	}

	return "", fmt.Errorf("cannot find zone id for zone name '%s'", c.config.Domain)

}

func (c *Client) getZoneMap() (Mapping, error) {

	resp, err := c.r53c.ListResourceRecordSets(&route53.ListResourceRecordSetsInput{
		HostedZoneID: aws.String(c.zoneId),
	})

	if err != nil {
		return nil, err
	}

	var m = NewMapping()

	for _, e := range resp.ResourceRecordSets {

		if *e.Type == "A" && strings.HasSuffix(*e.Name, fmt.Sprintf(".%s.", c.config.Domain)) {

			for _, value := range e.ResourceRecords {
				m.Add(*e.Name, *value.Value)
			}

		}

	}

	return m, nil

}

func findTagValue(key string, tags []*ec2.Tag) *string {

	for _, e := range tags {

		k := *e.Key

		if key == k {
			return e.Value
		}

	}

	return nil

}
