run:
	@wtc

prettier:
	prettier --write "pkg/**/*.graphql"

gqlgen:
	go generate pkg/graphql/handler.go
