package paramstore

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// Service provides a minimal interface to SSM Parameter Store,
// giving us greater flexibility in a test context.
type Service struct {
	Client ssm.GetParametersByPathAPIClient
}

// GetParametersByPath retrieves all parameters matching a given path prefix,
// returning them as name-value pairs in a map.
func (svc Service) GetParametersByPath(ctx context.Context, path string) (map[string]string, error) {
	input := &ssm.GetParametersByPathInput{
		Path:           aws.String(path),
		Recursive:      aws.Bool(true),
		WithDecryption: aws.Bool(true),
	}
	result := map[string]string{}
	paginator := ssm.NewGetParametersByPathPaginator(svc.Client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, p := range page.Parameters {
			result[*p.Name] = *p.Value
		}
	}
	return result, nil
}
