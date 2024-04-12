package internal

import (
	"barista/api/http"
	"barista/internal/modules"
	"barista/pkg/middlewares"
	"barista/pkg/models"
	"barista/pkg/repo"
	"barista/pkg/utils"
	"os"

	"github.com/spf13/cast"
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
	apiV1 := service.engine.Group("/api")

	UserRepo := repo.NewUserRepoImp(postgres)
	TokenRepo := repo.NewTokenRepoImp(postgres)
	UserHandler := modules.UserHandler{UserRepo: UserRepo, TokenRepo: TokenRepo, Postgres: postgres}
	userHttpHandler := http.User{Handler: &UserHandler}

	user := apiV1.Group("/user")
	user.Handle(string(models.POST), "signup", userHttpHandler.SignUp)
	user.Handle(string(models.POST), "login", userHttpHandler.Login)
	user.Handle(string(models.GET), "get-user", authMiddleware.IsAuthorized, userHttpHandler.GetUser)
	user.Handle(string(models.GET), "verify-email", userHttpHandler.VerifyEmail)
	user.Handle(string(models.POST), "forget-password", userHttpHandler.ForgetPassword)
	user.Handle(string(models.GET), "user-profile", userHttpHandler.UserProfile)

	cafeRepo := repo.NewCafeRepoImp(postgres)
	cafeHandler := modules.CafeHandler{CafeRepo: cafeRepo}
	cafeHttpHandler := http.Cafe{Handler: &cafeHandler}

	cafe := apiV1.Group("/cafe")
	cafe.Handle(string(models.POST), "create", cafeHttpHandler.Create)
	cafe.Handle(string(models.GET), "get-cafe", cafeHttpHandler.GetCafe)
	cafe.Handle(string(models.GET), "search-cafe", cafeHttpHandler.SearchCafe)

	service.Run(":8080")
}
