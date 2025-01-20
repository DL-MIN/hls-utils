package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spacecafe/gobox/gin-authentication"
	"github.com/spacecafe/gobox/gin-problems"
	"hls-utils/types"
)

func Authenticate(config *types.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request types.AuthenticateRequest
		var err error

		if err = ctx.ShouldBind(&request); err != nil {
			problems.ProblemBadRequest.WithError(err).Abort(ctx)
			return
		}

		token := config.Streams[request.Name]
		if token == "" {
			types.ProblemNoSuchStream.Abort(ctx)
			return
		}
		if authentication.ComparePasswords([]byte(token), []byte(request.Auth)) {
			ctx.AbortWithStatus(http.StatusOK)

			go func() {
				if config.NotificationEndpoint != "" {
					err := (&types.Notification{Name: request.Name}).Send(config.NotificationEndpoint)
					if err != nil {
						config.HTTPServer.Logger.Warn(err)
					}
				}
			}()
			return
		}

		problems.ProblemUnauthorized.Abort(ctx)
	}
}
