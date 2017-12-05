package paramstore

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/stretchr/testify/require"
)

type Service struct {
	client ssmiface.SSMAPI
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

func New(client ssmiface.SSMAPI) Service {
	return Service{client: client}
}

func (svc Service) GetParametersByPath(path string) (map[string]string, error) {
	input := &ssm.GetParametersByPathInput{Path: aws.String(path), WithDecryption: aws.Bool(true)}
	result := map[string]string{}
	err := svc.client.GetParametersByPathPages(
		input,
		func(resp *ssm.GetParametersByPathOutput, lastPage bool) bool {
			for _, p := range resp.Parameters {
				result[*p.Name] = *p.Value
			}
			return true
		},
	)
	return result, err
}
