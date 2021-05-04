//go:generate go run -v github.com/99designs/gqlgen

package graphql

import "context"

type Resolver struct {
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r Resolver) __Directive() __DirectiveResolver {
	panic("implement me")
}

type queryResolver struct{ *Resolver }

func (q queryResolver) CalculatePrice(ctx context.Context, typeArg *Type, margin *float64, exchangeRate *float64) (*float64, error) {
	panic("implement me")
}
