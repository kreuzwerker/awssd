package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kreuzwerker/awssd"
)

const (
	Empty         = ""
	RegionDefault = "eu-west-1"
)

var build string

func main() {

	var (
		domain        = flag.String("d", Empty, "Route53 domain name")
		dryRun        = flag.Bool("dry", false, "dry run, does not touch record sets")
		filter        = flag.String("f", Empty, "EC2 filter in key=value notation or multiple, comma-separated filters in the same notation")
		groupBy       = flag.String("g", Empty, "EC2 tag key to group instances by")
		help          = flag.Bool("h", false, "display usage informations")
		preferPrivate = flag.Bool("p", true, "Prefer private IP addresses (instead of public ones)")
		region        = flag.String("r", RegionDefault, fmt.Sprintf("AWS region for EC2 instances, default to %s", RegionDefault))
		ttl           = flag.Int64("t", 60, "TTL for DNS records, defaults to 60")
		verbose       = flag.Bool("v", false, "Be verbose")
	)

	flag.Parse()

	if *help || (*domain == Empty || *groupBy == Empty || *region == Empty) {
		printUsage()
	}

	awssd.Debug = *verbose

	config := &awssd.Config{
		Domain:        *domain,
		DryRun:        *dryRun,
		Filter:        *filter,
		GroupBy:       *groupBy,
		PreferPrivate: *preferPrivate,
		Region:        *region,
		TTL:           *ttl,
	}

	if *dryRun && *verbose {
		log.Printf("[ DEBUG ] Performing dry-run, won't touch any records")
	}

	client := awssd.NewClient(config)

	if changes, err := client.Diff(); err != nil {
		log.Fatalf("[ ERROR ] %v", err)
	} else if !changes {
		os.Exit(1)
	}

}

func printUsage() {

	fmt.Fprintf(os.Stderr, "Usage of %s (%s):\n", os.Args[0], build)
	flag.PrintDefaults()

	os.Exit(0)

}
