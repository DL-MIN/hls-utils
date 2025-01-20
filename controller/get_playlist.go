package controller

import (
	"github.com/gin-gonic/gin"
	problems "github.com/spacecafe/gobox/gin-problems"
	"hls-utils/types"
)

func GetPlaylist(streamManager *types.StreamManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var params types.GetPlaylistParams
		var err error

		if err = ctx.ShouldBindUri(&params); err != nil {
			problems.ProblemBadRequest.WithError(err).Abort(ctx)
			return
		}

		ctx.Header("Content-Type", "application/vnd.apple.mpegurl")
		if err = streamManager.GetPlaylist(ctx.Writer, params.Name); err != nil {
			types.ProblemNoSuchStream.Abort(ctx)
		}
	}
}
