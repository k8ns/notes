package ginlambda

import (
	"context"
	"github.com/ksopin/notes/internal/app"
	"github.com/ksopin/notes/internal/http"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

var ginLambda *ginadapter.GinLambda
var o sync.Once
var cfg *app.Config

func getGinLambda() *ginadapter.GinLambda {
	o.Do(func() {
		r := http.New(cfg)
		ginLambda = ginadapter.New(r)
	})

	return ginLambda
}

func Run(cfgApp *app.Config) error {
	cfg = cfgApp
	lambda.Start(Handler)
	return nil
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return getGinLambda().ProxyWithContext(ctx, req)
}
