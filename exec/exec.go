package exec

import (
	"os/exec"
	"syscall"

	"github.com/pkg/errors"
)

type Executable struct {
	Program string
	Argv    []string // including program name
	Environ []string
}

// ExecutableFromArgs takes the original os.Args (including the
// parameter-store-exec program name) and a (modified) os.Environ, returning
// an Executable.
func ExecutableFromArgs(argv, environ []string) (Executable, error) {
	exe := Executable{}
	argc := len(argv)

	if argc == 0 {
		return exe, errors.New("empty argv")
	} else if argc == 1 {
		return exe, errors.New(argv[0] + " expected program as first argument")
	}

	program, err := exec.LookPath(argv[1])
	if err != nil {
		return exe, errors.Wrap(err, "searching PATH")
	}

	exe.Program = program
	exe.Argv = argv[1:argc] // includes program name
	exe.Environ = environ

	return exe, nil
}

func (exe Executable) Exec() {
	syscall.Exec(exe.Program, exe.Argv, exe.Environ)
}
