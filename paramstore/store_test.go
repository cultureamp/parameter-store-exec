package paramstore_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/cultureamp/parameter-store-exec/paramstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeClient struct {
	ssmiface.SSMAPI
	path  string
	pages [][]*ssm.Parameter
	t     *testing.T
}

func TestGetParametersByPath(t *testing.T) {
	svc := paramstore.New(fakeClient{
		t:    t,
		path: "/foo/bar",
		pages: [][]*ssm.Parameter{
			{
				param("/foo/bar/one", "first"),
				param("/foo/bar/two", "second"),
			},
			{
				param("/foo/bar/three", "third"),
			},
		},
	})
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

func (f fakeClient) GetParametersByPathPages(input *ssm.GetParametersByPathInput, handler func(*ssm.GetParametersByPathOutput, bool) bool) error {
	require.Equal(f.t, *input.Path, f.path)
	for i, page := range f.pages {
		lastPage := (i == len(f.pages)-1)
		handler(&ssm.GetParametersByPathOutput{Parameters: page}, lastPage)
	}
	return nil
}
