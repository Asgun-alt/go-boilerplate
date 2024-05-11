package container

import (
	"go-boilerplate/config"
	"go-boilerplate/infrastructure/database"
	"go-boilerplate/infrastructure/jwt"
	"go-boilerplate/infrastructure/redis"
	"go-boilerplate/pkg/ratelimiter"

	authRepository "go-boilerplate/domain/auth/repository"
	authUser "go-boilerplate/domain/auth/usecase/user"
	userRepository "go-boilerplate/domain/user/repository"
	userUsecase "go-boilerplate/domain/user/usecase"
)

type Container struct {
	Config          *config.Config
	HTTPLimiter     *ratelimiter.RateLimiter
	UserUsecase     userUsecase.UserUsecase
	AuthUserUsecase authUser.Usecase
	JWT             jwt.IJWT
}

func New() *Container {
	cfg := config.New()

	db := database.New(&cfg.Infrastructure.Database)
	rdb := redis.New(cfg.Infrastructure.Redis)

	// init feature
	httpLimiter := ratelimiter.NewRateLimiter(
		cfg.Feature.RateLimiter.HTTPLimiter)
	ijwt := jwt.New(cfg.Credential.Jwt)

	// init service

	// init repository
	userRepository := userRepository.NewUserRepository(db)
	authRepository := authRepository.NewAuthRepository(rdb)

	// init usecase
	userUsecase := userUsecase.NewUserUsecase(userRepository)
	authUserUsecase := authUser.NewAuthUserUsecase(
		authRepository,
		cfg.Credential.Auth,
		cfg.Credential.Jwt,
		cfg.Credential.SMTP,
		userUsecase,
		ijwt,
	)

	return &Container{
		Config:          cfg,
		HTTPLimiter:     httpLimiter,
		UserUsecase:     userUsecase,
		AuthUserUsecase: authUserUsecase,
		JWT:             ijwt,
	}
}
