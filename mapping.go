package awssd

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/route53"
)

type Action func(name string, ip []IP) error

// +gen set
type IP string

func (i IP) ToResourceRecord() *route53.ResourceRecord {
	return &route53.ResourceRecord{
		Value: aws.String(string(i)),
	}
}

type Mapping map[string]IPSet

func NewMapping() Mapping {
	return make(map[string]IPSet)
}

func (m Mapping) Add(name, ip string) {

	if set, ok := m[name]; ok {
		set.Add(IP(ip))
	} else {
		m[name] = NewIPSet(IP(ip))
	}

}

func (m Mapping) Diff(state Mapping, upsert Action) (bool, error) {

	var action = false

	if reflect.DeepEqual(m, state) {
		return action, nil
	}

	for tk, tv := range m {

		if ov, present := state[tk]; !present {

			if err := upsert(tk, tv.ToSlice()); err != nil {
				return action, err
			}

			action = true

		} else if diff := tv.SymmetricDifference(ov); diff.Cardinality() > 0 {

			if err := upsert(tk, tv.ToSlice()); err != nil {
				return action, err
			}

			action = true

		}

	}

	return action, nil

}

func (m Mapping) String() string {

	var buf []string

	for k, v := range m {
		buf = append(buf, fmt.Sprintf("%s %q", k, v.ToSlice()))
	}

	return strings.Join(buf, ", ")

}
