
all: test vet lint

build:
	go build -v -o bin/edgr ./cmd

lint:
	golint -set_exit_status ./...

test:
	go test -v ./...

vet:
	go vet ./...

coverage:
	# go currently cannot create coverage profiles when testing multiple packages, so we test each package
	# independently. This issue should be fixed in Go 1.10 (https://github.com/golang/go/issues/6909).
	go list ./... | xargs -n1 -I {} -P 4 go test -covermode=count -coverprofile=../../../{}/profile.coverprofile {}

clean:
	find . -name \*.coverprofile -delete
	rm -rf bin

db:
	docker run -p 5432:5432 -d postgres:9.6

pack:
	~/go/bin/go-bindata -ignore=\\.DS_Store  -o cmd/migrations.go migrations/...