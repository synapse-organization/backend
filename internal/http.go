package internal

import (
	"barista/pkg/log"
	"barista/pkg/models"
	"barista/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"reflect"
	"strings"
)

type Service struct {
	engine *gin.Engine
}

func StartService() *Service {
	gin.SetMode(gin.ReleaseMode)
	return &Service{
		engine: gin.New(),
	}
}
func (s *Service) AddGroup(path string) *gin.RouterGroup {
	return s.engine.Group(path)
}

func (s *Service) AddRoutes(group *gin.RouterGroup, route ...models.Route) {
	for _, r := range route {
		group.Handle(r.Method, r.Path, r.Function)
	}
}

func (s *Service) AddStructRoutes(group *gin.RouterGroup, route interface{}) {
	routeType := reflect.TypeOf(route)
	if routeType.Kind() == reflect.Struct {
		routeGroup := group.Group(strings.ToLower(routeType.Name()))

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
