
GO = docker run --rm -e GOOS=linux -e GOARCH=amd64 -v $(PWD):/usr/src/app -w /usr/src/app golang:1.13 go
R = docker run --rm --name notes --network dev -e APP_ENV=dev -e GOOS=linux -e GOARCH=amd64 -v $(PWD):/usr/src/app -w /usr/src/app golang:1.13 go
LAMBDA_VERSION = 7

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

promote: set_env_prod publish_version point_prod_alias api_gateway_permissions set_env_stage
	#NotesAPI:${stageVariables.lambdaAlias}



set_env_prod:
	aws lambda update-function-configuration \
		--function-name NotesAPI \
		--environment Variables={APP_ENV=prod}

publish_version:
	aws lambda publish-version \
        --function-name "arn:aws:lambda:eu-west-1:720657473133:function:NotesAPI"

point_prod_alias:
	aws lambda update-alias \
        --function-name NotesAPI \
        --function-version ${LAMBDA_VERSION} \
        --name prod

api_gateway_permissions:
	aws lambda add-permission \
		--function-name "arn:aws:lambda:eu-west-1:720657473133:function:NotesAPI:prod" \
		--source-arn "arn:aws:execute-api:eu-west-1:720657473133:lccdn1r7ib/*/*/" \
		--principal apigateway.amazonaws.com \
		--statement-id ef133303-5ce3-4950-884f-abc0324fcf66 \
		--action lambda:InvokeFunction

	aws lambda add-permission \
		--function-name "arn:aws:lambda:eu-west-1:720657473133:function:NotesAPI:prod" \
		--source-arn "arn:aws:execute-api:eu-west-1:720657473133:lccdn1r7ib/*/*/*" \
		--principal apigateway.amazonaws.com \
		--statement-id a42da76e-38f2-498f-97cf-f3bbcc92f3fc \
		--action lambda:InvokeFunction

set_env_stage:
	aws lambda update-function-configuration \
		--function-name NotesAPI \
		--environment Variables={APP_ENV=stage}



deploy_test:
	curl https://finance-api.kvslab.icu/info


clean:
	rm build/out/notes_api || true
	rm build/out/main.zip || true