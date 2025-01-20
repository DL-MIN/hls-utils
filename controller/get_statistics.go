package controller

import (
	"github.com/gin-gonic/gin"
	problems "github.com/spacecafe/gobox/gin-problems"
	"hls-utils/types"
)

func GetStatistics(streamManager *types.StreamManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var params types.GetStatisticsParams
		var err error
		if err = ctx.ShouldBindUri(&params); err != nil {
			problems.ProblemBadRequest.WithError(err).Abort(ctx)
			return
		}

		if stream := streamManager.GetStream(params.Name); stream != nil {
			ctx.JSON(200, gin.H{"subscribers": stream.Statistics.Len()})
			return
		}

		types.ProblemNoSuchStream.Abort(ctx)
	}
}
