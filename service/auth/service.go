package auth

import "github.com/devesh44/gin-poc/repository"

type AuthS struct {
	ar repository.AuthRepo
}

// type AuthService interface {
// 	Create()
// }

func NewAuthService() *AuthS {
	ar := repository.NewAuthRepo()
	return &AuthS{ar: ar}
}

func (a *AuthS) Create() {

}
