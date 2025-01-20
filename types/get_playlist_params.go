package types

type GetPlaylistParams struct {
	Name string `uri:"name" binding:"required,printascii"`
}
