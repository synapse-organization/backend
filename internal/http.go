package internal

import (
	"barista/pkg/log"
	"barista/pkg/models"
	"barista/pkg/utils"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Service struct {
	engine *gin.Engine
}

func StartService() *Service {
	service := Service{}
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	engine.Use(service.CorMiddleware())
	service.engine = engine
	return &service
}
func (s *Service) AddGroup(path string) *gin.RouterGroup {
	return s.engine.Group(path)
}

func (s *Service) AddRoutes(group *gin.RouterGroup, route ...models.Route) {
	for _, r := range route {
		group.Handle(r.Method, r.Path, r.Function)
	}
}

func (s *Service) AddStructRoutes(group *gin.RouterGroup, route interface{}, middleware ...gin.HandlerFunc) {
	routeType := reflect.TypeOf(route)
	if routeType.Kind() == reflect.Struct {
		routeGroup := group.Group(strings.ToLower(routeType.Name()))

		routeGroup.Use(middleware...)

		log.GetLog(true).WithField("BasePath", routeGroup.BasePath()).Info("Registering group")
		for i := 0; i < routeType.NumMethod(); i++ {
			method := routeType.Method(i)
			if method.Type.NumIn() != 2 {
				log.GetLog().Errorf("Method %v has more than 2 arguments", method.Name)
				continue
			}
			if method.Type.NumOut() != 0 {
				log.GetLog().Errorf("Method %v has more than 0 return value", method.Name)
				continue
			}
			methodType, name := utils.SplitMethodPrefix(method.Name)
			if strings.ToUpper(methodType) == "SET" {
				methodType = "POST"
			}
			log.GetLog(true).WithFields(logrus.Fields{
				"Method": strings.ToUpper(methodType),
				"Group":  strings.ToLower(routeType.Name()),
				"Path":   routeGroup.BasePath() + "/" + strings.ToLower(name),
			}).Info("Registering route")
			routeGroup.Handle(strings.ToUpper(methodType), "/"+strings.ToLower(name), func(c *gin.Context) {
				method.Func.Call([]reflect.Value{reflect.ValueOf(route), reflect.ValueOf(c)})
			})
		}
		return
	} else {
		log.GetLog().Error("Route is not a struct")
	}
}

func (s *Service) CorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (s *Service) AddMiddleware(middleware ...gin.HandlerFunc) {
	s.engine.Use(middleware...)
}

func (s *Service) Run(address string) {
	start := s.engine.Group("/")
	s.AddRoutes(start, models.Route{
		Method: "GET",
		Path:   "/ping",
		Function: func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		},
	})

	log.GetLog(true).WithField("Address", address).Info("Start listening HTTP")
	err := s.engine.Run(address)
	if err != nil {
		log.GetLog().Fatal(err)
	}

}
