package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/cultureamp/parameter-store-exec/paramstore"
	"github.com/pkg/errors"
)

const (
	pathEnv = "PARAMETER_STORE_EXEC_PATH"
	disableTranslation = "PARAMETER_STORE_EXEC_DISABLE_TRANSLATION"
)

var transformPattern *regexp.Regexp

var Version = "dev"

func init() {
	transformPattern = regexp.MustCompile("[^A-Z0-9_]")
}

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Printf("parameter-store-exec %s\n", Version)
		return
	}

	log.SetOutput(os.Stderr)

	argv, err := argvForExec(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	program, err := exec.LookPath(argv[0])
	if err != nil {
		log.Fatal(err)
	}

	environ := os.Environ()

	if path := os.Getenv(pathEnv); path != "" {
		store := paramstore.Service{
			Client: ssm.New(session.Must(session.NewSession(&aws.Config{}))),
		}
		params, err := store.GetParametersByPath(path)
		if err != nil {
			log.Fatal(err)
		}

		translationEnabled := isTranslationEnabled()

		for name, v := range params {
			envKey := paramToEnv(name, path, translationEnabled)
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

	syscall.Exec(program, argv, environ)
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
func paramToEnv(name, path string, translate bool) string {
	result := strings.TrimPrefix(strings.TrimPrefix(name, path), "/")
	if translate == true {
		result = strings.ToUpper(result)
		result = transformPattern.ReplaceAllLiteralString(result, "_")
	}
	return result
}

func isTranslationEnabled() bool {
	disableValue, exists := os.LookupEnv(disableTranslation)
	if !exists {
		return true
	}

	disable, err := strconv.ParseBool(disableValue)
	if err != nil {
		log.Printf("%s had an unparseable value, enabling translation by default", disableTranslation)
		disable = false
	}

	return !disable
}
