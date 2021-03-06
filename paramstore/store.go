package paramstore

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

// Service provides a minimal interface to SSM Parameter Store,
// and should be
type Service struct {
	Client ssmiface.SSMAPI
}

func (svc Service) GetParametersByPath(path string) (map[string]string, error) {
	input := &ssm.GetParametersByPathInput{
		Path:           aws.String(path),
		Recursive:      aws.Bool(true),
		WithDecryption: aws.Bool(true),
	}
	result := map[string]string{}
	err := svc.Client.GetParametersByPathPages(
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
