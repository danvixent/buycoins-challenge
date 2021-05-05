//go:generate go run -v github.com/99designs/gqlgen

package graphql

import (
	"context"
	"fmt"

	"github.com/danvixent/buycoins_challenge/handlers/margin"
)

type Resolver struct {
	marginHandler *margin.Handler
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r Resolver) __Directive() __DirectiveResolver {
	return nil
}

type queryResolver struct{ *Resolver }

func (q *queryResolver) CalculatePrice(ctx context.Context, typeArg *Type, margin *float64, exchangeRate *float64) (*float64, error) {
	switch *typeArg {
	case TypeBuy:
		NGNPrice, err := q.marginHandler.AddPriceMargin(ctx, *margin, *exchangeRate)
		if err != nil {
			return nil, fmt.Errorf("add price margin error: %v", err)
		}
		return NGNPrice, nil
	case TypeSell:
		NGNPrice, err := q.marginHandler.SubtractPriceMargin(ctx, *margin, *exchangeRate)
		if err != nil {
			return nil, fmt.Errorf("subtract price margin error: %v", err)
		}
		return NGNPrice, nil
	}
	return nil, nil
}
