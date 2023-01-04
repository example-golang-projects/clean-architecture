package user

import (
	"fmt"
	"github.com/example-golang-projects/clean-architecture/cmd/server/user/config"
	"github.com/example-golang-projects/clean-architecture/services/user/dependency"
	"github.com/example-golang-projects/clean-architecture/services/user/modules/role/handler"
	"net/http"
)

type UserService struct {
	roleHandler handler.RoleHandler
}

func (s *UserService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.InitUserRouter().ServeHTTP(w, r)
}

func NewUserService(cfg config.Config) *UserService {
	dep := dependency.InitUserDependency(cfg)
	return &UserService{
		roleHandler: handler.NewRoleHandler(dep),
	}
}

func RunUserService(cfg config.Config) {
	server := http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.ServicePort),
		Handler: NewUserService(cfg),
	}
	server.ListenAndServe()
}
