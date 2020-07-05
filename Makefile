build:
	go build -o gvm cmd/main.go

test:
	go test ./... -v

fmt:
	go fmt ./...

release:
	go get github.com/mitchellh/gox
	$$GOPATH/bin/gox -os="linux windows freebsd" -output="dist/gvm.{{.OS}}.{{.Arch}}" ./cmd
	$$GOPATH/bin/gox -osarch="darwin/amd64" -output="dist/gvm.{{.OS}}.{{.Arch}}" ./cmd
	./release.sh


