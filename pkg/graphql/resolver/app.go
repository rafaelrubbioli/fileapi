package resolver

import (
	"github.com/rafaelrubbioli/fileapi/pkg/graphql/gqlgen"
	"github.com/rafaelrubbioli/fileapi/pkg/service"
)

type app struct {
	service service.Service
}

func (a app) Query() gqlgen.QueryResolver {
	return query{app: &a}
}

func (a app) Mutation() gqlgen.MutationResolver {
	return mutation{app: &a}
}

func New(service service.Service) gqlgen.ResolverRoot {
	return &app{
		service: service,
	}
}
