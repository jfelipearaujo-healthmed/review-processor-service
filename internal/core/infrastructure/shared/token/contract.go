package token

type TokenService interface {
	CreateJwtToken(userID uint) (string, error)
}
