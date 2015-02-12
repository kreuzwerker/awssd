package awssd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSymetricDifference(t *testing.T) {

	var (
		ec2 = NewIPSet("a", "b", "c")
		r53 = NewIPSet("a", "c", "d")
	)

	assert.Equal(t, NewIPSet("b", "d"), ec2.SymmetricDifference(r53))

}
