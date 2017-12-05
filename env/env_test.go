package env_test

import (
	"sort"
	"testing"

	"github.com/cultureamp/parameter-store-exec/env"
	"github.com/stretchr/testify/assert"
)

func TestNameTransform(t *testing.T) {
	for name, expected := range map[string]string{
		"ZERO":             "ZERO",
		"ONE_ONE":          "ONE_ONE",
		"two":              "TWO",
		"three-three":      "THREE_THREE",
		"FourFour":         "FOURFOUR",
		"five five":        "FIVE_FIVE",
		"Six!@#$%^&*()siX": "SIX__________SIX",
	} {
		assert.Equal(t, expected, env.TransformName(name))
	}
}

func TestAppendWithoutOverwrite(t *testing.T) {
	fakeKeys := env.FakeEnv{Keys: map[string]bool{
		"PATH":        true,
		"USER":        true,
		"ENVIRONMENT": true,
	}}
	orig := []string{
		"PATH=/bin",
		"USER=root",
		"ENVIRONMENT=production",
	}
	new := map[string]string{
		"ENVIRONMENT": "already set",
		"NEW_THING":   "shiny",
	}
	expected := []string{
		"PATH=/bin",
		"USER=root",
		"ENVIRONMENT=production",
		"NEW_THING=shiny",
	}
	updated := env.AppendWithoutOverwrite(orig, new, fakeKeys)
	sort.Strings(expected)
	sort.Strings(updated)
	assert.Equal(t, expected, updated)
}
