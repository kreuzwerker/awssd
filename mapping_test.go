package awssd

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiffError(t *testing.T) {

	var (
		ec2 = NewMapping()
		r53 = NewMapping()
	)

	// create some diff to act upon
	ec2.Add("a.api.example.com", "1.2.3.4")
	r53.Add("b.api.example.com", "1.2.3.4")

	err := fmt.Errorf("boom")

	assert.Equal(t, err, ec2.Diff(r53, func(host string, ips []IP) error {
		return err
	}))

}

func TestDiff(t *testing.T) {

	assert := assert.New(t)

	var (
		ec2  = NewMapping()
		r53  = NewMapping()
		diff = func() ([]string, error) {

			var a sort.StringSlice

			err := ec2.Diff(r53, func(host string, ips []IP) error {

				var s sort.StringSlice

				for _, e := range ips {
					s = append(s, string(e))
				}

				s.Sort()

				a = append(a, fmt.Sprintf("%s: %s", host, strings.Join(s, ", ")))

				return nil

			})

			a.Sort()

			return []string(a), err
		}
	)

	ec2.Add("a.api.example.com", "1.2.3.4")
	r53.Add("a.api.example.com", "1.2.3.4")

	// equal mapping
	assert.Nil(ec2.Diff(r53, nil))

	// change ec2
	ec2.Add("a.api.example.com", "1.2.3.5")

	add, _ := diff()

	assert.Len(add, 1)
	assert.Equal("a.api.example.com: 1.2.3.4, 1.2.3.5", add[0])

	// change zone
	r53.Add("b.api.example.com", "2.3.4.5")
	r53.Add("b.api.example.com", "1.2.3.7")
	r53.Add("c.api.example.com", "1.2.3.8")

	add, _ = diff()

	assert.Len(add, 1)
	assert.Equal("a.api.example.com: 1.2.3.4, 1.2.3.5", add[0])

	// change both
	ec2.Add("b.api.example.com", "1.2.3.6")
	r53.Add("b.api.example.com", "2.3.4.5")

	add, _ = diff()

	assert.Len(add, 2)
	assert.Equal("a.api.example.com: 1.2.3.4, 1.2.3.5", add[0])
	assert.Equal("b.api.example.com: 1.2.3.6", add[1])

}
