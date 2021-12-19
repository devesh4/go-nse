package repository

import (
	"github.com/devesh44/gin-poc/config"
	"gorm.io/gorm"
)

type authrepo struct {
	db *gorm.DB
}

type AuthRepo interface {
	Create()
}

func NewAuthRepo() AuthRepo {
	d := config.MainConf.DB
	return &authrepo{db: d}
}

func (a *authrepo) Create() {

}
