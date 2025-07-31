package query

import (
	"jcourse_go/internal/application/viewobject"
	"jcourse_go/internal/domain/auth"
	"jcourse_go/internal/domain/common"
	"jcourse_go/pkg/apperror"
)

type UserQueryService interface {
	GetUserInfo(commonCtx *common.CommonContext) (*viewobject.UserInfoVO, error)
	UpdateUserInfo(commonCtx *common.CommonContext, nickname string) error
}

type userQueryService struct {
	userRepo auth.UserRepository
}

func NewUserQueryService(userRepo auth.UserRepository) UserQueryService {
	return &userQueryService{
		userRepo: userRepo,
	}
}

func (s *userQueryService) GetUserInfo(commonCtx *common.CommonContext) (*viewobject.UserInfoVO, error) {
	user, err := s.userRepo.GetByID(commonCtx.Ctx, commonCtx.User.UserID)
	if err != nil {
		return nil, apperror.ErrDB.Wrap(err)
	}
	if user == nil {
		return nil, apperror.ErrUserNotFound
	}

	userInfoVO := viewobject.NewUserInfoVO(user)
	return &userInfoVO, nil
}

func (s *userQueryService) UpdateUserInfo(commonCtx *common.CommonContext, nickname string) error {
	user, err := s.userRepo.GetByID(commonCtx.Ctx, commonCtx.User.UserID)
	if err != nil {
		return apperror.ErrDB.Wrap(err)
	}
	if user == nil {
		return apperror.ErrUserNotFound
	}

	user.UpdateNickname(nickname)
	return s.userRepo.Update(commonCtx.Ctx, user)
}
