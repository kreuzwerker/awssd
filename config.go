package awssd

import (
	"fmt"
	"strings"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/gen/ec2"
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

func (c *Config) ConvertFilter() (filters []ec2.Filter, err error) {

	if c.Filter == "" {
		return filters, err
	}

	for _, e := range strings.Split(c.Filter, ",") {

		filter := strings.Split(e, "=")

		if len(filter) != 2 {
			return nil, fmt.Errorf("cannot create filter from filter specification '%v'", filter)
		}

		filters = append(filters, ec2.Filter{
			aws.String(filter[0]),
			filter[1:],
		})

	}

	return filters, err

}
