package internal

import (
	"barista/api/http"
	"barista/internal/modules"
	"barista/pkg/middlewares"
	"barista/pkg/models"
	"barista/pkg/repo"
	"barista/pkg/utils"
	"github.com/fergusstrange/embedded-postgres"
	"github.com/spf13/cast"
	"os"
)

func IsTest() bool {
	test, ok := os.LookupEnv("TEST")
	if !ok {
		return false
	}

	return cast.ToBool(test)
}

func Run() {
	if IsTest() {
		postgres := embeddedpostgres.NewDatabase()
		err := postgres.Start()
		if err != nil {
			panic(err)
		}
	}

	postgres := utils.NewPostgres(
		models.Postgres{
			Host:     "localhost",
			Port:     5432,
			UserName: "postgres",
			Password: "postgres",
			DbName:   "postgres",
		},
	)

	authMiddleware := middlewares.AuthMiddleware{}
	service := StartService()
	apiV1 := service.AddGroup("/api/v1/")

	UserRepo := repo.NewUserRepoImp(postgres)
	UserHandler := modules.UserHandler{UserRepo: UserRepo}

	service.AddStructRoutes(apiV1, http.Auth{
		Handler: &UserHandler,
	})

	service.AddStructRoutes(apiV1, http.User{}, authMiddleware.IsAuthorized())

	service.Run(":8080")
}
