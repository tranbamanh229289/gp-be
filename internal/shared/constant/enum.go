package constant

type UserRole string
const (
	UserRoleUser UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

type TokenType string 
const (
	AccessToken TokenType = "access_token"
	RefreshToken TokenType = "refresh_token"
)