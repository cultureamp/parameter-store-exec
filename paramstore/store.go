package paramstore

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

type Service struct {
	client ssmiface.SSMAPI
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
