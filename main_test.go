package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParamToEnvWithTranslation(t *testing.T) {
	for name, expected := range map[string]string{
		"/a/b/ZERO":             "ZERO",
		"/a/b/ONE_ONE":          "ONE_ONE",
		"/a/b/two":              "TWO",
		"/a/b/three-three":      "THREE_THREE",
		"/a/b/FourFour":         "FOURFOUR",
		"/a/b/five five":        "FIVE_FIVE",
		"/a/b/Six!@#$%^&*()siX": "SIX__________SIX",
		"/a/b/c/d/seven":        "C_D_SEVEN",
		"/a/b/eight8":           "EIGHT8",
	} {
		assert.Equal(t, expected, paramToEnv(name, "/a/b", true))
	}
}

func TestParamToEnvWithoutTranslation(t *testing.T) {
	for name, expected := range map[string]string{
		"/a/b/ZERO":             "ZERO",
		"/a/b/ONE_ONE":          "ONE_ONE",
		"/a/b/two":              "two",
		"/a/b/three-three":      "three-three",
		"/a/b/FourFour":         "FourFour",
		"/a/b/five five":        "five five",
		"/a/b/Six!@#$%^&*()siX": "Six!@#$%^&*()siX",
		"/a/b/c/d/seven":        "c/d/seven",
		"/a/b/eight8":           "eight8",
	} {
		assert.Equal(t, expected, paramToEnv(name, "/a/b", false))
	}
}

func TestArgvForExec(t *testing.T) {
	argv, err := argvForExec([]string{"parameter-store-exec", "echo", "hello"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"echo", "hello"}, argv)
}

func TestArgvForExecError(t *testing.T) {
	_, err := argvForExec([]string{"parameter-store-exec"})
	assert.Error(t, err)
}
