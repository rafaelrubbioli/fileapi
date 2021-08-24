package resolver

import (
	"github.com/rafaelrubbioli/fileapi/pkg/graphql/gqlgen"
	"github.com/rafaelrubbioli/fileapi/pkg/service"
)

type app struct {
	service service.Service
}

func (a app) Query() gqlgen.QueryResolver {
	return query{service: a.service}
}

func (a app) Mutation() gqlgen.MutationResolver {
	return mutation{service: a.service}
}

func New(service service.Service) gqlgen.ResolverRoot {
	return &app{
		service: service,
	}
}
