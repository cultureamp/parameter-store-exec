package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/cultureamp/parameter-store-exec/paramstore"
	"github.com/pkg/errors"
)

const (
	pathEnv = "PARAMETER_STORE_EXEC_PATH"
)

var transformPattern = regexp.MustCompile("[^A-Z0-9_]")

var Version = "dev"

func main() {
	log.SetOutput(os.Stderr)
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Printf("parameter-store-exec %s\n", Version)
		return nil
	}

	ctx := context.Background()

	argv, err := argvForExec(os.Args)
	if err != nil {
		return err
	}

	program, err := exec.LookPath(argv[0])
	if err != nil {
		log.Fatal(err)
	}

	environ := os.Environ()

	if path := os.Getenv(pathEnv); path != "" {
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return err
		}
		client := ssm.NewFromConfig(cfg)

		store := paramstore.Service{
			Client: client,
		}
		params, err := store.GetParametersByPath(ctx, path)
		if err != nil {
			return err
		}
		for name, v := range params {
			envKey := paramToEnv(name, path)
			if _, present := os.LookupEnv(envKey); present {
				log.Printf("%s => %s already set", name, envKey)
			} else {
				environ = append(environ, envKey+"="+v)
				log.Printf("%s => %s", name, envKey)
			}
		}
	} else {
		log.Println(pathEnv, "not set")
	}

	err = syscall.Exec(program, argv, environ)
	if err != nil {
		return err
	}

	return nil
}

// argvForExec takes the os.Args from parameter-store-exec,
// ensures there's a sub-command specified,
// and returns a new argv ready for syscall.Exec().
func argvForExec(osargs []string) ([]string, error) {
	argc := len(osargs)
	switch argc {
	case 0:
		return nil, errors.New("empty osargs")
	case 1:
		return nil, errors.New(osargs[0] + " expected program as first argument")
	default:
		return osargs[1:argc], nil
	}
}

// paramToEnv takes a SSM Parameter Store key name like /foo/bar/api-key,
// strips the specified path prefix e.g. /foo,
// and returns an environment-friendly name like BAR_API_KEY.
func paramToEnv(name, path string) string {
	pathStripped := strings.TrimPrefix(strings.TrimPrefix(name, path), "/")
	upper := strings.ToUpper(pathStripped)
	return transformPattern.ReplaceAllLiteralString(upper, "_")
}
