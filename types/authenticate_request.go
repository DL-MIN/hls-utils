package types

// AuthenticateRequest represents the request structure for authentication.
type AuthenticateRequest struct {
	// Call specifies the type of call being made. It must be exactly "publish".
	Call string `form:"call" binding:"required,eq=publish"`

	// Name is a field representing the name associated with the request.
	// It must contain only printable ASCII characters and is required.
	Name string `form:"name" binding:"required,printascii"`

	// Auth is a field used for authentication purposes.
	// Similar to Name, it must also contain only printable ASCII characters and is required.
	Auth string `form:"auth" binding:"required,printascii"`
}
