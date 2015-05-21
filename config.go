package awssd

import (
	"fmt"
	"strings"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/ec2"
)

var Debug bool

type Config struct {
	Domain        string
	DryRun        bool
	Filter        string
	GroupBy       string
	PreferPrivate bool
	Region        string
	TTL           int64
}

func (c *Config) ConvertFilter() (filters []*ec2.Filter, err error) {

	if c.Filter == "" {
		return filters, err
	}

	for _, e := range strings.Split(c.Filter, ",") {

		filter := strings.Split(e, "=")

		if len(filter) != 2 {
			return nil, fmt.Errorf("cannot create filter from filter specification '%v'", filter)
		}

		var values []*string

		for _, ee := range filter[1:] {
			values = append(values, &ee)
		}

		filters = append(filters, &ec2.Filter{
			Name:   aws.String(filter[0]),
			Values: values,
		})

	}

	return filters, err

}
