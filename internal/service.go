package internal

import (
	"barista/api/http"
	"barista/internal/modules"
	"barista/pkg/middlewares"
	"barista/pkg/models"
	"barista/pkg/repo"
	"barista/pkg/utils"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func getMongo() (string, int) {
	address, ok := os.LookupEnv("mongo_address")
	if !ok {
		return "localhost", 27017
	}

	port, ok := os.LookupEnv("mongo_port")
	if !ok {
		return "localhost", 27017
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

	address, port = getMongo()
	mongoDb := utils.ConnectDB(
		models.Mongo{
			Host:     address,
			Port:     port,
			UserName: "root",
			Password: "root",
		},
	)
	mongoDbOpt := options.GridFSBucket().SetName("image-server")

	authMiddleware := middlewares.AuthMiddleware{Postgres: postgres}

	service := StartService()
	apiV1 := service.engine.Group("/api")

	UserRepo := repo.NewUserRepoImp(postgres)
	TokenRepo := repo.NewTokenRepoImp(postgres)
	UserHandler := modules.UserHandler{UserRepo: UserRepo, TokenRepo: TokenRepo, Postgres: postgres}
	userHttpHandler := http.User{Handler: &UserHandler}

	user := apiV1.Group("/user")
	user.Handle(string(models.POST), "signup", userHttpHandler.SignUp)
	user.Handle(string(models.POST), "login", userHttpHandler.Login)
	user.Handle(string(models.POST), "logout", authMiddleware.IsAuthorized, userHttpHandler.Logout)
	user.Handle(string(models.GET), "get-user", authMiddleware.IsAuthorized, userHttpHandler.GetUser)
	user.Handle(string(models.GET), "verify-email", userHttpHandler.VerifyEmail)
	user.Handle(string(models.POST), "forget-password", userHttpHandler.ForgetPassword)
	user.Handle(string(models.GET), "user-profile", authMiddleware.IsAuthorized, userHttpHandler.UserProfile)
	user.Handle(string(models.PUT), "edit-profile", authMiddleware.IsAuthorized, userHttpHandler.EditProfile)

	cafeRepo := repo.NewCafeRepoImp(postgres)
	imageRepo := repo.NewImageRepoImp(postgres)
	ratingRepo := repo.NewRatingsRepoImp(postgres)
	cafeHandler := modules.CafeHandler{CafeRepo: cafeRepo, Rating: ratingRepo, ImageRepo: imageRepo}
	cafeHttpHandler := http.Cafe{Handler: &cafeHandler}

	cafe := apiV1.Group("/cafe")
	cafe.Handle(string(models.POST), "create", cafeHttpHandler.Create)
	cafe.Handle(string(models.GET), "get-cafe", cafeHttpHandler.GetCafe)
	cafe.Handle(string(models.POST), "search-cafe", cafeHttpHandler.SearchCafe)

	imageHandler := http.ImageHandler{MongoDb: mongoDb, MongoOpt: mongoDbOpt, ImageRepo: imageRepo}
	image := apiV1.Group("/image")
	image.Handle(string(models.POST), "upload", imageHandler.UploadImage)
	image.Handle(string(models.GET), "download", imageHandler.DownloadImage)
	image.Handle(string(models.POST), "submit", imageHandler.SubmitImage)

	service.Run(":8080")
}
