package internal

import (
	"barista/api/http"
	"barista/internal/modules"
	"barista/pkg/log"
	"barista/pkg/middlewares"
	"barista/pkg/models"
	"barista/pkg/repo"
	"barista/pkg/utils"
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
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

func getRedis() (string, int) {
	address, ok := os.LookupEnv("redis_address")
	if !ok {
		return "localhost", 6379
	}

	port, ok := os.LookupEnv("redis_port")
	if !ok {
		return "localhost", 6379
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

	address, port = getRedis()
	rdb := redis.NewClient(&redis.Options{
		Addr:     address + ":" + cast.ToString(port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	authMiddleware := middlewares.AuthMiddleware{Postgres: postgres}

	service := StartService()
	apiV1 := service.engine.Group("/api")

	userRepo := repo.NewUserRepoImp(postgres)
	tokenRepo := repo.NewTokenRepoImp(postgres)
	cafeRepo := repo.NewCafeRepoImp(postgres)
	paymentRepo := repo.NewTransactionImp(postgres)
	reservationRepo := repo.NewReservationRepoImp(postgres)
	UserHandler := modules.UserHandler{UserRepo: userRepo, TokenRepo: tokenRepo, ReservationRepo: reservationRepo, CafeRepo: cafeRepo, Postgres: postgres}
	userHttpHandler := http.User{Handler: &UserHandler}

	user := apiV1.Group("/user")
	user.Handle(string(models.POST), "signup", userHttpHandler.SignUp)
	user.Handle(string(models.POST), "login", userHttpHandler.Login)
	user.Handle(string(models.POST), "logout", authMiddleware.IsAuthorized, userHttpHandler.Logout)
	user.Handle(string(models.GET), "get-user", authMiddleware.IsAuthorized, userHttpHandler.GetUser)
	user.Handle(string(models.GET), "verify-email", userHttpHandler.VerifyEmail)
	user.Handle(string(models.POST), "forget-password", userHttpHandler.ForgetPassword)
	user.Handle(string(models.POST), "change-password", authMiddleware.IsAuthorized, userHttpHandler.ChangePassword)
	user.Handle(string(models.GET), "user-profile", authMiddleware.IsAuthorized, userHttpHandler.UserProfile)
	user.Handle(string(models.PATCH), "edit-profile", authMiddleware.IsAuthorized, userHttpHandler.EditProfile)
	user.Handle(string(models.POST), "manager-agreement", authMiddleware.IsAuthorized, userHttpHandler.ManagerAgreement)
	user.Handle(string(models.GET), "user-reservations", authMiddleware.IsAuthorized, userHttpHandler.UserReservations)

	imageRepo := repo.NewImageRepoImp(postgres)
	ratingRepo := repo.NewRatingsRepoImp(postgres)
	commentRepo := repo.NewCommentsRepoImp(postgres)
	eventRepo := repo.NewEventRepoImp(postgres)
	menuItemRepo := repo.NewMenuItemRepoImp(postgres)
	locationRepo := repo.NewLocationsRepoImp(postgres)
	favoriteRepo := repo.NewFavoritesRepoImp(postgres)

	cafeHandler := modules.CafeHandler{
		CafeRepo:        cafeRepo,
		Rating:          ratingRepo,
		CommentRepo:     commentRepo,
		ImageRepo:       imageRepo,
		EventRepo:       eventRepo,
		UserRepo:        userRepo,
		ReservationRepo: reservationRepo,
		MenuItemRepo:    menuItemRepo,
		PaymentRepo:     paymentRepo,
		FavoriteRepo:    favoriteRepo,
		LocationsRepo:   locationRepo,
		Redis:           rdb,
	}
	cafeHttpHandler := http.Cafe{Handler: &cafeHandler, Rating: ratingRepo, ImageRepo: imageRepo, FirstSearch: true}

	newTicker := time.NewTicker(1 * time.Minute)
	go func() {
		for {
			select {
			case <-newTicker.C:
				if err := rdb.FlushDB(context.Background()).Err(); err != nil {
					log.GetLog().Errorf("Unable to flush redis db. error: %v", err)
				}
				locations, err := locationRepo.FindAll(context.Background())
				if err != nil {
					log.GetLog().Errorf("Unable to get locations. error: %v", err)
				}

				for _, location := range locations {
					if err := rdb.GeoAdd(context.Background(), "locations", &redis.GeoLocation{
						Name:      cast.ToString(location.CafeID),
						Longitude: location.Lng,
						Latitude:  location.Lat,
					}).Err(); err != nil {
						log.GetLog().Errorf("Unable to add location to redis. error: %v", err)
					}
				}
			}
		}
	}()

	cafe := apiV1.Group("/cafe")
	cafe.Handle(string(models.POST), "create", authMiddleware.IsAuthorized, cafeHttpHandler.Create)
	cafe.Handle(string(models.POST), "search-cafe", cafeHttpHandler.SearchCafe)
	cafe.Handle(string(models.GET), "public-cafe", authMiddleware.OptionalAuth, cafeHttpHandler.PublicCafeProfile)
	cafe.Handle(string(models.POST), "add-comment", authMiddleware.IsAuthorized, cafeHttpHandler.AddComment)
	//cafe.Handle(string(models.GET), "get-comments", cafeHttpHandler.GetComments)
	cafe.Handle(string(models.POST), "create-event", authMiddleware.IsAuthorized, cafeHttpHandler.CreateEvent)
	cafe.Handle(string(models.POST), "add-menu-item", authMiddleware.IsAuthorized, cafeHttpHandler.AddMenuItem)
	cafe.Handle(string(models.GET), "home", cafeHttpHandler.Home)
	cafe.Handle(string(models.GET), "private-menu", authMiddleware.IsAuthorized, cafeHttpHandler.PrivateMenu)
	cafe.Handle(string(models.GET), "public-menu", cafeHttpHandler.PublicMenu)
	cafe.Handle(string(models.PATCH), "edit-menu-item", authMiddleware.IsAuthorized, cafeHttpHandler.EditMenuItem)
	cafe.Handle(string(models.DELETE), "delete-menu-item", authMiddleware.IsAuthorized, cafeHttpHandler.DeleteMenuItem)
	cafe.Handle(string(models.POST), "reserve-event", authMiddleware.IsAuthorized, cafeHttpHandler.ReserveEvent)
	cafe.Handle(string(models.GET), "private-cafe", authMiddleware.IsAuthorized, cafeHttpHandler.PrivateCafe)
	cafe.Handle(string(models.PATCH), "edit-cafe", authMiddleware.IsAuthorized, cafeHttpHandler.EditCafe)
	cafe.Handle(string(models.PATCH), "edit-event", authMiddleware.IsAuthorized, cafeHttpHandler.EditEvent)
	cafe.Handle(string(models.DELETE), "delete-event", authMiddleware.IsAuthorized, cafeHttpHandler.DeleteEvent)
	cafe.Handle(string(models.GET), "fully-booked-days", cafeHttpHandler.GetFullyBookedDays)
	cafe.Handle(string(models.GET), "time-slots", cafeHttpHandler.GetTimeSlots)
	cafe.Handle(string(models.POST), "reserve-cafe", authMiddleware.IsAuthorized, cafeHttpHandler.ReserveCafe)
	cafe.Handle(string(models.GET), "cafe-reservations", authMiddleware.IsAuthorized, cafeHttpHandler.GetCafeReservations)
	cafe.Handle(string(models.POST), "add-to-favorite", authMiddleware.IsAuthorized, cafeHttpHandler.AddToFavorite)
	cafe.Handle(string(models.DELETE), "remove-favorite", authMiddleware.IsAuthorized, cafeHttpHandler.RemoveFavorite)
	cafe.Handle(string(models.GET), "get-favorite-list", authMiddleware.IsAuthorized, cafeHttpHandler.GetFavoriteList)

	//rating
	cafe.Handle(string(models.POST), "add-rating", authMiddleware.IsAuthorized, cafeHttpHandler.AddRating)
	cafe.Handle(string(models.GET), "get-cafe-rating", authMiddleware.IsAuthorized, cafeHttpHandler.GetRating)

	// location
	cafe.Handle(string(models.POST), "get-nearest-cafes", cafeHttpHandler.GetNearestCafes)
	cafe.Handle(string(models.POST), "get-cafe-location", cafeHttpHandler.GetCafeLocation)
	cafe.Handle(string(models.POST), "set-location", cafeHttpHandler.SetCafeLocation)

	imageHandler := http.ImageHandler{MongoDb: mongoDb, MongoOpt: mongoDbOpt, ImageRepo: imageRepo}
	image := apiV1.Group("/image")
	image.Handle(string(models.POST), "upload", imageHandler.UploadImage)
	image.Handle(string(models.GET), "download", imageHandler.DownloadImage)
	image.Handle(string(models.POST), "submit", imageHandler.SubmitImage)

	paymentHandler := modules.PaymentHandler{PaymentRepo: paymentRepo, UserRepo: userRepo}
	paymentHttpHandler := http.Payment{Handler: &paymentHandler}
	payment := apiV1.Group("/payment")
	payment.Handle(string(models.POST), "transfer", authMiddleware.IsAuthorized, paymentHttpHandler.Transfer)
	payment.Handle(string(models.GET), "transactions-list", authMiddleware.IsAuthorized, paymentHttpHandler.TransactionsList)
	payment.Handle(string(models.POST), "deposit", authMiddleware.IsAuthorized, paymentHttpHandler.Deposit)
	payment.Handle(string(models.POST), "withdraw", authMiddleware.IsAuthorized, paymentHttpHandler.Withdraw)
	payment.Handle(string(models.GET), "balance", authMiddleware.IsAuthorized, paymentHttpHandler.Balance)

	publicHandler := http.PublicHandler{}
	public := apiV1.Group("/public")
	public.StaticFile("/province", "./assets/ostan.json")
	public.Handle(string(models.GET), "health", publicHandler.HealthCheck)
	public.Handle(string(models.GET), "/cities", publicHandler.GetCities)

	service.Run(":8080")
}
