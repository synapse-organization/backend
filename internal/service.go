package internal

import (
	"barista/api/http"
)

func Run() {
	service := StartService()

	apiV1 := service.AddGroup("/api/v1/")
	service.AddStructRoutes(apiV1, http.User{})

	service.Run(":8080")
}
