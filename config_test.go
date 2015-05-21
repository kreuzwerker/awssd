package awssd

import (
	"testing"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/assert"
)

func TestConvertFilters(t *testing.T) {

	assert := assert.New(t)

	var (
		a = aws.String("a")
		b = aws.String("b")
		c = aws.String("c")
		d = aws.String("d")
	)

	var tt = []struct {
		in  string
		out []*ec2.Filter
	}{
		{"", nil},
		{"a=b", []*ec2.Filter{
			&ec2.Filter{
				Name:   a,
				Values: []*string{b},
			},
		}},
		{"a=b,c=d", []*ec2.Filter{
			&ec2.Filter{
				Name:   a,
				Values: []*string{b},
			},
			&ec2.Filter{
				Name:   c,
				Values: []*string{d},
			},
		}},
	}

	for _, e := range tt {

		config := &Config{
			Filter: e.in,
		}

		filters, err := config.ConvertFilter()

		assert.NoError(err)
		assert.Equal(e.out, filters)

	}

}
