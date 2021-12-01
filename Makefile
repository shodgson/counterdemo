.PHONY: build clean deploy gomodgen

# Default parameters
STAGE = dev
FUNCTION = counter

build: gomodgen
	export GO111MODULE=on
	cd back && go mod tidy
	cd back && env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/counter_api cmd/counter_api.go

clean:
	rm -rf ./back/bin ./back/vendor ./back/go.sum

deploy: clean build
	sls deploy --verbose --stage $(STAGE)

deployfunction: clean build
	sls deploy function --function $(FUNCTION) --stage $(STAGE)

gomodgen:
	cd back && chmod u+x gomod.sh
	cd back && ./gomod.sh

test:
	go test ./... -v

integrationtest:
	go test ./... -v -tags=integration
