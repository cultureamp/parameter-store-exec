# Installing parameter-store-exec on an m1 Macbook

1. Change directory into `/usr/local/share` with the command `cd /usr/local/share`

2. Clone the github repo to your local machine. 

    ```
    sudo git clone https://github.com/cultureamp/parameter-store-exec.git
    ```

3. Change ownership of the cloned directory to your user with the command `sudo chown -R username:admin parameter-store-exec`

4. Change to the cloned directory with `cd parameter-store-exec`

5. Install go, if not already installed using homebrew by running the command `brew install go`

6. Initialize the go.mod file by running `go mod init paramstore`

7. Install all other required packages with `go mod tidy`

8. Mark the current directory as safe in your global git config with the command `git config --global --add safe.directory /usr/local/share/parameter-store-exec`

9. Compile the binary with the command `GOOS=linux GOARCH=arm64 go build -ldflags="-X main.Version=$(git describe --tags --candidates=1 --dirty) -s -w" -o $@ github.com/cultureamp/parameter-store-exec`

10. Create a wrapper for the binary you just compiled by creating `/usr/local/bin/parameter-store-exec` with the following contents

    ```bash
    #!/usr/bin/env bash

    SOURCE_DIR='/usr/local/share/parameter-store-exec'
    EXECUTABLE='github.com/cultureamp/parameter-store-exec'

    cd $SOURCE_DIR
    go run $EXECUTABLE $@

    ```

## Testing your install

1. In ALKS get your AWS keys and load them as environment variables

2. Set the parameter-store-exec environment variables

    ```
    export AWS_REGION=us-east-1
    export PARAMETER_STORE_EXEC_PATH=/path/to/your/parameters
    ```

3. Run the command `parameter-store-exec env`