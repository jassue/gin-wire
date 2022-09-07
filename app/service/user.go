package service

import (
    "context"
    "github.com/gin-gonic/gin"
    "github.com/jassue/gin-wire/app/domain"
    cErr "github.com/jassue/gin-wire/app/pkg/error"
    "github.com/jassue/gin-wire/app/pkg/request"
    "github.com/jassue/gin-wire/utils/hash"
    "strconv"
)

type UserRepo interface {
    FindByID(context.Context, uint64) (*domain.User, error)
    FindByMobile(context.Context, string) (*domain.User, error)
    Create(context.Context, *domain.User) (*domain.User, error)
}

type UserService struct {
    uRepo UserRepo
    tm Transaction
}

// NewUserService .
func NewUserService(uRepo UserRepo, tm Transaction) *UserService {
    return &UserService{uRepo: uRepo, tm: tm}
}

// Register 注册
func (s *UserService) Register(ctx *gin.Context, param *request.Register) (*domain.User, error) {
    user, _ := s.uRepo.FindByMobile(ctx, param.Mobile)
    if user != nil {
        return nil, cErr.BadRequest("手机号码已存在")
    }

    u, err := s.uRepo.Create(ctx, &domain.User{
        Name:     param.Name,
        Mobile:   param.Mobile,
        Password: hash.BcryptMake([]byte(param.Password)),
    })
    if err != nil {
        return nil, cErr.BadRequest("注册用户失败")
    }

    return u, nil
}

// Login 登录
func (s *UserService) Login(ctx *gin.Context, mobile, password string) (*domain.User, error) {
    u, err := s.uRepo.FindByMobile(ctx, mobile)
    if err != nil || !hash.BcryptMakeCheck([]byte(password), u.Password) {
        return nil, cErr.BadRequest("用户名不存在或密码错误")
    }

    return u, nil
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(ctx *gin.Context, idStr string) (*domain.User, error) {
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        return nil, cErr.NotFound("数据ID错误")
    }
    u, err := s.uRepo.FindByID(ctx, id)
    if err != nil {
        return nil, cErr.NotFound("数据不存在", cErr.USER_NOT_FOUND)
    }

    return u, nil
}