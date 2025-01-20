package types

type GetVariantPlaylistParams struct {
	Name     string `uri:"name" binding:"required,printascii"`
	ClientID string `uri:"client_id" binding:"required,uuid"`
	Variant  string `uri:"variant" binding:"required,printascii"`
}
