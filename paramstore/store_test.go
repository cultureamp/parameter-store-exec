package paramstore_test

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/smithy-go/middleware"
	"github.com/cultureamp/parameter-store-exec/paramstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetParametersByPath(t *testing.T) {
	svc := paramstore.Service{Client: FakeClient{
		T:    t,
		Path: "/foo/bar",
		Pages: [][]types.Parameter{
			{
				param("/foo/bar/one", "first"),
				param("/foo/bar/two", "second"),
			},
			{
				param("/foo/bar/three", "third"),
			},
		},
	}}
	params, err := svc.GetParametersByPath(t.Context(), "/foo/bar")
	require.NoError(t, err)
	assert.Equal(t, map[string]string{
		"/foo/bar/one":   "first",
		"/foo/bar/two":   "second",
		"/foo/bar/three": "third",
	}, params)
}

func param(name, value string) types.Parameter {
	return types.Parameter{
		ARN:              nil,
		DataType:         nil,
		LastModifiedDate: nil,
		Name:             &name,
		Selector:         nil,
		SourceResult:     nil,
		Type:             "",
		Value:            &value,
		Version:          0,
	}
}

var _ ssm.GetParametersByPathAPIClient = FakeClient{Path: "", Pages: nil, T: nil}

type FakeClient struct {
	Path  string
	Pages [][]types.Parameter
	T     *testing.T
}

func (f FakeClient) GetParametersByPath(ctx context.Context, input *ssm.GetParametersByPathInput, opts ...func(*ssm.Options)) (*ssm.GetParametersByPathOutput, error) {
	require.Equal(f.T, *input.Path, f.Path)

	pageNum, err := getPage(input.NextToken)
	if err != nil {
		return nil, err
	}

	if pageNum >= len(f.Pages) {
		return nil, errors.New("invalid page")
	}

	nextPage := pageNum + 1
	output := &ssm.GetParametersByPathOutput{
		NextToken:      nil,
		Parameters:     f.Pages[pageNum],
		ResultMetadata: middleware.Metadata{},
	}

	if nextPage < len(f.Pages) {
		output.NextToken = aws.String(strconv.Itoa(nextPage))
	}

	return output, nil
}

func getPage(nextToken *string) (int, error) {
	pageNum := 0

	var err error

	if nextToken != nil {
		pageNum, err = strconv.Atoi(*nextToken)
		if err != nil {
			return 0, err
		}
	}

	return pageNum, nil
}
