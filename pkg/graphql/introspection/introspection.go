package introspection

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type Introspection struct{}

func (c Introspection) ExtensionName() string {
	return "Introspection"
}

func (c Introspection) Validate(_ graphql.ExecutableSchema) error {
	return nil
}

func (c Introspection) MutateOperationContext(_ context.Context, rc *graphql.OperationContext) *gqlerror.Error {
	rc.DisableIntrospection = false
	return nil
}
