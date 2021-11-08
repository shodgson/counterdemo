package apilogger

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	lmdrouter "github.com/shodgson/lmdrouterv2"
)

func Logger(next lmdrouter.Handler) lmdrouter.Handler {
	return func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (
		res events.APIGatewayV2HTTPResponse,
		err error,
	) {
		// [LEVEL] [METHOD PATH QUERY] [CODE] EXTRA
		format := "[%s] [%s %s %s] [%d] %s"
		level := "INF"
		var code int
		var extra string

		res, err = next(ctx, req)
		if err != nil {
			level = "ERR"
			code = http.StatusInternalServerError
			extra = " " + err.Error()
		} else {
			code = res.StatusCode
			if code >= 400 {
				level = "ERR"
			}
		}

		log.Printf(format, level, req.RequestContext.HTTP.Method, req.RequestContext.HTTP.Path, req.QueryStringParameters, code, extra)

		return res, err
	}
}
