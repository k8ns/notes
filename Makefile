
GO = docker run --rm -e GOOS=linux -e GOARCH=amd64 -v $(PWD):/usr/src/app -w /usr/src/app golang:1.13 go
R = docker run --rm --name notes --network dev -e APP_ENV=dev -e GOOS=linux -e GOARCH=amd64 -v $(PWD):/usr/src/app -w /usr/src/app golang:1.13 go

.PHONY: clean test build lambda_build lambda_artifact lambda_deploy deploy

build: clean test lambda_build lambda_artifact

run:
	$(R) run ./cmd/gin.go

test:
	@echo "Testing..."
	$(GO) test cmd/gin.go


lambda_build:
	@echo "Building..."
	$(GO) build -o build/out/notes_api cmd/lambda.go


lambda_artifact:
	zip -j build/out/main.zip build/out/notes_api
	zip build/out/main.zip config/*.yml data/key.pem


deploy:
	@echo "Deploying..."
	aws lambda update-function-code --function-name NotesAPI \
    --zip-file "fileb://build/out/main.zip"


deploy_test:
	curl https://finance-api.kvslab.icu/info


clean:
	rm build/out/notes_api || true
	rm build/out/main.zip || true