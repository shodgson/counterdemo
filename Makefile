.PHONY: build clean deploy gomodgen

# Default parameters
STAGE = dev # Set to 'production' for live server
FUNCTION = counter
BACKEND = back
FRONTEND = front


buildfront: 
	cd $(FRONTEND) && npm run build

deployfront:
	aws s3 sync ./$(FRONTEND)/dist s3://$(BUCKETNAME) --profile $(AWS_PROFILE)

buildback: gomodgen
	export GO111MODULE=on
	cd $(BACKEND) && go mod tidy
	cd $(BACKEND) && env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/counter_api cmd/counter/counter_api.go
	cd $(BACKEND) && env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/preSignUp cmd/cognito/preSignUp.go


build: buildback buildfront

clean:
	rm -rf ./$(BACKEND)/bin ./$(BACKEND)/vendor ./$(BACKEND)/go.sum ./$(FRONTEND)/dist

deploy: clean buildback
	sls deploy --verbose --stage $(STAGE)

deployfunction: buildback
	sls deploy function --function $(FUNCTION) --stage $(STAGE)

gomodgen:
	cd $(BACKEND) && chmod u+x gomod.sh
	cd $(BACKEND) && ./gomod.sh

test: clean buildback
	sh -ac ' . ./.env; cd $(BACKEND) && go test ./... -v'

integrationtest: clean buildback
	sh -ac ' . ./.env; cd $(BACKEND) && go test ./... -v -tags=integration'
