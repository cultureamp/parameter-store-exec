package exec_test

import (
	"testing"

	"github.com/cultureamp/parameter-store-exec/exec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecutableFromArgs(t *testing.T) {
	cmd, err := exec.ExecutableFromArgs(
		[]string{"parameter-store-exec", "echo", "hello", "world"},
		[]string{"a=1", "b=2"},
	)
	require.NoError(t, err)
	assert.Equal(t, "/bin/echo", cmd.Program)
	assert.Equal(t, []string{"echo", "hello", "world"}, cmd.Argv)
	assert.Equal(t, []string{"a=1", "b=2"}, cmd.Environ)
}

func TestExecutableFromArgsProgramOnly(t *testing.T) {
	cmd, err := exec.ExecutableFromArgs(
		[]string{"parameter-store-exec", "/bin/echo"},
		[]string{},
	)
	require.NoError(t, err)
	assert.Equal(t, "/bin/echo", cmd.Program)
	assert.Equal(t, []string{"/bin/echo"}, cmd.Argv)
}

func TestExecutableFromArgsMissingArg(t *testing.T) {
	_, err := exec.ExecutableFromArgs(
		[]string{"parameter-store-exec"},
		[]string{},
	)
	require.Error(t, err)
}

func TestExecutableFromArgsCommandNotInPath(t *testing.T) {
	_, err := exec.ExecutableFromArgs(
		[]string{"parameter-store-exec", "probably-not-a-real-command"},
		[]string{},
	)
	require.Error(t, err)
}
