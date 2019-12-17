package ginlambda


import (
	"github.com/ksopin/notes/internal/app"
	"github.com/ksopin/notes/internal/http"
	"context"
	"github.com/ksopin/notes/pkg/db"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"

)

var ginLambda *ginadapter.GinLambda
var o sync.Once

func getGinLambda() *ginadapter.GinLambda {
	o.Do(func() {
		cfg := app.GetConfig(strings.Join([]string{"config/config", os.Getenv("APP_ENV"), "yml"}, "."))

		db.InitConnection(cfg.Db)
		r := http.New()
		http.InitWelcome(r, cfg.App)
		ginLambda = ginadapter.New(r)
	})

	return ginLambda
}

func Run() error {
	lambda.Start(Handler)
	return nil
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return getGinLambda().ProxyWithContext(ctx, req)
}
