package handler

import (
	"github.com/SawitProRecruitment/UserService/helper"
	"github.com/SawitProRecruitment/UserService/repository"
)

type Server struct {
	Repository repository.RepositoryInterface
	Helper     helper.HelperInterface
}

type NewServerOptions struct {
	Repository repository.RepositoryInterface
	Helper     helper.HelperInterface
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		Repository: opts.Repository,
		Helper:     opts.Helper,
	}
}
