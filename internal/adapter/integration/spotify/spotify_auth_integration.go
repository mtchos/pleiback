package spotify

type AuthIntegration interface {
	Authorize() (AuthorizeResponse, error)
	GetToken(code string) (GetTokenResponse, error)
}

type AuthorizeResponse struct {
	Code  string
	State string
}

type GetTokenResponse struct {
	AccessToken  string
	ExpiresIn    int64
	RefreshToken string
}
