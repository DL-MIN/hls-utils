package types

type AuthenticateRequest struct {
	Call string `form:"call" binding:"required,eq=publish"`
	Name string `form:"name" binding:"required,printascii"`
	Auth string `form:"auth" binding:"required,printascii"`
}
