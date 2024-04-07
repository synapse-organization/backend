package internal

import (
	"barista/api/http"
	"barista/internal/modules"
	"barista/pkg/middlewares"
	"barista/pkg/models"
	"barista/pkg/repo"
	"barista/pkg/utils"
	"github.com/spf13/cast"
	"os"
)

func getPostgres() (string, int) {
	address, ok := os.LookupEnv("postgres_address")
	if !ok {
		return "localhost", 5432
	}

	port, ok := os.LookupEnv("postgres_port")
	if !ok {
		return "localhost", 5432
	}

	return address, cast.ToInt(port)
}

func Run() {
	address, port := getPostgres()
	postgres := utils.NewPostgres(
		models.Postgres{
			Host:     address,
			Port:     port,
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
