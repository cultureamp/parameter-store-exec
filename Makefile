package = github.com/cultureamp/parameter-store-exec

.PHONY: release
release: parameter-store-exec-darwin-amd64.gz parameter-store-exec-linux-amd64.gz

%.gz: %
	gzip $<

%-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -o $@ $(package)

%-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o $@ $(package)
