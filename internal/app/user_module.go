package app

import (
	userhttp "golab/internal/modules/user/delivery/http"
	usersworker "golab/internal/modules/user/delivery/worker"
	userevents "golab/internal/modules/user/infra/events"
	usermailer "golab/internal/modules/user/infra/mailer"
	userspostgres "golab/internal/modules/user/infra/postgres"
	usersapp "golab/internal/modules/user/usecase"
	"golab/internal/platform/config"
)

type UserModule struct {
	HTTPHandler   *userhttp.Handler
	WorkerHandler *usersworker.Handler
}

func buildUserModule(infra *Infra, cfg config.Config) UserModule {
	userRepo := buildUserRepository(infra)
	userMailer := buildUserMailer(infra, cfg)
	userPublisher := buildUserPublisher(infra, cfg)
	userUseCase := buildUserUseCase(userRepo, userMailer, userPublisher)

	return buildUserHandlers(userUseCase, infra)
}

func buildUserRepository(infra *Infra) *userspostgres.UserRepository {
	return userspostgres.NewUserRepository(infra.DB)
}

func buildUserMailer(infra *Infra, cfg config.Config) *usermailer.Service {
	return usermailer.NewService(infra.SMTPClient, cfg.SMTPFrom, infra.Logger)
}

func buildUserPublisher(infra *Infra, cfg config.Config) *userevents.Publisher {
	return userevents.NewPublisher(infra.Rabbit, cfg.RabbitMQActivationRoutingKey)
}

func buildUserUseCase(
	userRepo *userspostgres.UserRepository,
	userMailer *usermailer.Service,
	userPublisher *userevents.Publisher,
) *usersapp.UserUseCase {
	return usersapp.NewUserUseCase(userRepo, userMailer, userPublisher)
}

func buildUserHandlers(userUseCase *usersapp.UserUseCase, infra *Infra) UserModule {
	return UserModule{
		HTTPHandler:   userhttp.NewHandler(userUseCase, infra.Logger),
		WorkerHandler: usersworker.NewHandler(userUseCase, infra.Logger),
	}
}
