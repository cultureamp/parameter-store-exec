package env

import (
	"os"
	"regexp"
	"strings"
)

var transformPattern *regexp.Regexp

func init() {
	transformPattern = regexp.MustCompile("[^A-Z_]")
}

type EnvChecker interface {
	Exists(key string) bool
}

type Env struct{}

func (e Env) Exists(key string) bool {
	_, ok := os.LookupEnv(key)
	return ok
}

type FakeEnv struct {
	Keys map[string]bool
}

func (f FakeEnv) Exists(key string) bool {
	return f.Keys[key]
}

func TransformName(name string) string {
	return transformPattern.ReplaceAllLiteralString(
		strings.ToUpper(name),
		"_",
	)
}

func AppendWithoutOverwrite(orig []string, new map[string]string, env EnvChecker) []string {
	updated := orig
	for k, v := range new {
		if !env.Exists(k) {
			updated = append(updated, k+"="+v)
		}
	}
	return updated
}
