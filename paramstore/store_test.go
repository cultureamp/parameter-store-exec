package paramstore_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/cultureamp/parameter-store-exec/paramstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetParametersByPath(t *testing.T) {
	svc := paramstore.Service{Client: FakeClient{
		T:    t,
		Path: "/foo/bar",
		Pages: [][]*ssm.Parameter{
			{
				param("/foo/bar/one", "first"),
				param("/foo/bar/two", "second"),
			},
			{
				param("/foo/bar/three", "third"),
			},
		},
	}}
	params, err := svc.GetParametersByPath("/foo/bar")
	require.NoError(t, err)
	assert.Equal(t, map[string]string{
		"/foo/bar/one":   "first",
		"/foo/bar/two":   "second",
		"/foo/bar/three": "third",
	}, params)
}

func param(name, value string) *ssm.Parameter {
	return &ssm.Parameter{Name: &name, Value: &value}
}

type FakeClient struct {
	ssmiface.SSMAPI
	Path  string
	Pages [][]*ssm.Parameter
	T     *testing.T
}

func (f FakeClient) GetParametersByPathPages(input *ssm.GetParametersByPathInput, handler func(*ssm.GetParametersByPathOutput, bool) bool) error {
	require.Equal(f.T, *input.Path, f.Path)
	for i, page := range f.Pages {
		lastPage := (i == len(f.Pages)-1)
		handler(&ssm.GetParametersByPathOutput{Parameters: page}, lastPage)
	}
	return nil
}
