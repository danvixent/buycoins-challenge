package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/99designs/gqlgen/graphql"
)

func TestCalculatePrice(t *testing.T) {
	type args struct {
		typeArg      string
		margin       float64
		exchangeRate float64
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
		wantErr  bool
	}{
		{
			name: "calculatePrice1",
			args: args{
				typeArg:      "buy",
				margin:       0.2,
				exchangeRate: 505,
			},
			wantCode: http.StatusOK,
			wantErr:  false,
		},
		{
			name: "calculatePrice2",
			args: args{
				typeArg:      "sell",
				margin:       0.4,
				exchangeRate: 484,
			},
			wantCode: http.StatusOK,
			wantErr:  false,
		},
		{
			name: "calculatePrice3",
			args: args{
				typeArg:      "sold",
				margin:       0.4,
				exchangeRate: 484,
			},
			wantCode: http.StatusUnprocessableEntity,
			wantErr:  false,
		},
	}

	const baseQuery = `query{
		calculatePrice(type: %s,margin: %f,exchangeRate: %f)
	}`

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			query := fmt.Sprintf(baseQuery, tt.args.typeArg, tt.args.margin, tt.args.exchangeRate)
			gql := graphql.RawParams{Query: query}

			resp, err := calculatePrice(gql)
			if (err != nil) != tt.wantErr {
				t.Errorf("calculatePrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if resp.StatusCode != tt.wantCode {
				t.Errorf("calculatePrice() status code = %v, want %v", resp.StatusCode, tt.wantCode)
			}
		})
	}
}

func calculatePrice(body interface{}) (*http.Response, error) {
	return POST(baseURL+"/graphql", serialize(body))
}

// serialize obj into json bytes
func serialize(obj interface{}) *bytes.Buffer {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(obj); err != nil {
		// if json encoding fails, stop the test immediately
		log.Fatalf("unable to serialize obj: %v", err)
	}
	return buf
}

func POST(url string, body *bytes.Buffer) (*http.Response, error) {
	return http.Post(url, "application/json", body)
}
