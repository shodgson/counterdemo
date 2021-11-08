.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	go mod tidy
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/counter cmd/counter.go

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose --stage $(STAGE)

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
