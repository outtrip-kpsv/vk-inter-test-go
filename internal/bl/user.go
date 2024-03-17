package bl

import (
	"go.uber.org/zap"
	"time"
	"vk-inter-test-go/internal/db/repo"
	utilsJwt "vk-inter-test-go/internal/utils"
)

func (b *BL) CreateUser(user repo.User) (string, error) {
	user.Pass, _ = utilsJwt.HashPassword(user.Pass)
	err := b.Db.User.CreateUser(user)
	if err != nil {
		return "", err
	}
	token, err := utilsJwt.GenerateToken(time.Hour, user.Login)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (b *BL) AuthUser(user repo.User) (string, error) {
	dbUser, err := b.Db.User.GetUserByLogin(user.Login)
	if err != nil {
		return "", err
	}
	err = utilsJwt.VerifyPassword(dbUser.Pass, user.Pass)
	if err != nil {
		return "", err
	}
	token, err := utilsJwt.GenerateToken(time.Hour*5, user.Login)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (b *BL) CheckJwt(token string) bool {
	validateToken, err := utilsJwt.ValidateToken(token)

	if err != nil {
		return false
	}
	b.logger.Info("check jwt ok", zap.String("user: ", validateToken.(string)))
	return true
}

func (b *BL) CheckRole(name string, role string) bool {
	user, err := b.Db.User.GetUserByLogin(name)
	if err != nil {
		b.logger.Info("err", zap.Error(err))
		return false
	}
	roleDb, err := b.Db.Role.GetRoleById(user.RoleID)
	if err != nil {
		b.logger.Info("err", zap.Error(err))
		return false
	}
	return roleDb.Name == role
}
