package main

import (
	"os"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/helper"
	"github.com/SawitProRecruitment/UserService/repository"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	var server generated.ServerInterface = newServer()

	generated.RegisterHandlers(e, server)
	e.Logger.Fatal(e.Start(":1323"))
}

func newServer() *handler.Server {
	dbDsn := os.Getenv("DATABASE_URL")

	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: dbDsn,
	})

	var helper helper.HelperInterface = helper.NewHelper()

	opts := handler.NewServerOptions{
		Repository: repo,
		Helper:     helper,
	}

	return handler.NewServer(opts)
}
