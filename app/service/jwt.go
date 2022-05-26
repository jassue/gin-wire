package service

import (
    "context"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt"
    "github.com/jassue/gin-wire/app/compo"
    "github.com/jassue/gin-wire/app/domain"
    cErr "github.com/jassue/gin-wire/app/pkg/error"
    "github.com/jassue/gin-wire/config"
    "go.uber.org/zap"
    "time"
)

type JwtService struct {
    conf *config.Configuration
    log *zap.Logger
    jRepo JwtRepo
    uS *UserService
    lockBuilder *compo.LockBuilder
}

func NewJwtService(conf *config.Configuration, log *zap.Logger, jRepo JwtRepo, uS *UserService, lb *compo.LockBuilder) *JwtService {
    return &JwtService{
        conf: conf,
        log:  log,
        jRepo: jRepo,
        uS: uS,
        lockBuilder: lb,
    }
}

type JwtRepo interface {
    JoinBlackList(ctx context.Context, tokenStr string, joinUnix int64, expires time.Duration) error
    GetBlackJoinUnix(ctx context.Context, tokenStr string) (int64, error)
}

func (s *JwtService) CreateToken(GuardName string, user domain.JwtUser) (*domain.TokenOutPut, *jwt.Token, error) {
    token := jwt.NewWithClaims(
        jwt.SigningMethodHS256,
        domain.CustomClaims{
            StandardClaims: jwt.StandardClaims{
                ExpiresAt: time.Now().Unix() + s.conf.Jwt.JwtTtl,
                Id:        user.GetUid(),
                Issuer:    GuardName,
                NotBefore: time.Now().Unix() - 1000,
            },
        },
    )

    tokenStr, err := token.SignedString([]byte(s.conf.Jwt.Secret))
    if err != nil {
        return nil, nil, cErr.BadRequest("create token error:" + err.Error())
    }

    return &domain.TokenOutPut{
        AccessToken: tokenStr,
        ExpiresIn:   int(s.conf.Jwt.JwtTtl),
        TokenType:   domain.TokenType,
    }, token, nil
}

func (s *JwtService) JoinBlackList(ctx *gin.Context, token *jwt.Token) error {
    nowUnix := time.Now().Unix()
    timer := time.Duration(token.Claims.(*domain.CustomClaims).ExpiresAt - nowUnix) * time.Second

    if err := s.jRepo.JoinBlackList(ctx, token.Raw, nowUnix, timer); err != nil {
        s.log.Error(err.Error())
        return cErr.BadRequest("登出失败")
    }

    return nil
}

func (s *JwtService) IsInBlacklist(ctx *gin.Context, tokenStr string) bool {
    joinUnix, err := s.jRepo.GetBlackJoinUnix(ctx, tokenStr)
    if err != nil {
        return false
    }

    if time.Now().Unix()-joinUnix < s.conf.Jwt.JwtBlacklistGracePeriod {
        return false
    }
    return true
}

func (s *JwtService) GetUserInfo(ctx *gin.Context, guardName, id string) (domain.JwtUser, error) {
    switch guardName {
    case domain.AppGuardName:
        return s.uS.GetUserInfo(ctx, id)
    default:
        return nil, cErr.BadRequest("guard " + guardName +" does not exist")
    }
}

func (s *JwtService) RefreshToken(ctx *gin.Context, guardName string, token *jwt.Token) (*domain.TokenOutPut, error) {
    idStr := token.Claims.(*domain.CustomClaims).Id

    lock := s.lockBuilder.NewLock(ctx, "refresh_token_lock:" + idStr, s.conf.Jwt.JwtBlacklistGracePeriod)
    if lock.Get() {
        user, err := s.GetUserInfo(ctx, guardName, idStr)
        if err != nil {
            s.log.Error(err.Error())
            lock.Release()
            return nil, err
        }

        tokenData, _, err := s.CreateToken(guardName, user)
        if err != nil {
            lock.Release()
            return nil, err
        }

        err = s.JoinBlackList(ctx, token)
        if err != nil {
            lock.Release()
            return nil, err
        }

        return tokenData, nil
    }

    return nil, cErr.BadRequest("系统繁忙")
}