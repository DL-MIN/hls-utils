package controller

import (
	"path"

	"github.com/gin-gonic/gin"
	problems "github.com/spacecafe/gobox/gin-problems"
	"hls-utils/types"
)

func GetVariantPlaylist(config *types.Config, streamManager *types.StreamManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var params types.GetVariantPlaylistParams
		if err := ctx.ShouldBindUri(&params); err != nil {
			problems.ProblemBadRequest.WithError(err).Abort(ctx)
			return
		}

		stream := streamManager.GetStream(params.Name)
		if stream == nil {
			types.ProblemNoSuchStream.Abort(ctx)
			return
		}

		stream.Statistics.Add(params.ClientID)

		filename := path.Join(config.LiveDirectory, path.Clean(params.Name), path.Clean(params.Variant), "index.m3u8")
		ctx.Header("Content-Type", "application/vnd.apple.mpegurl")
		ctx.File(filename)
	}
}
