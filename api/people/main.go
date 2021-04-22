package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	graphql "github.com/graph-gophers/graphql-go"
)

// Schema : GraphQL schema definition. This is an example schema
var Schema = `
	schema {
		query: Query
	}

	type Person{
		id: ID!
		firstName: String!
		lastName: String
	}

	type Query{
		person(id: ID!): Person
	}
`

type person struct {
	ID        graphql.ID
	FirstName string
	LastName  string
}

var people = []*person{
	{
		ID:        "1000",
		FirstName: "Pedro",
		LastName:  "Marquez",
	},

	{
		ID:        "1001",
		FirstName: "John",
		LastName:  "Doe",
	},
}

type personResolver struct {
	p *person
}

func (r *personResolver) ID() graphql.ID {
	return r.p.ID
}

func (r *personResolver) FirstName() string {
	return r.p.FirstName
}

func (r *personResolver) LastName() *string {
	return &r.p.LastName
}

// Resolver : Struct with all the resolver functions
type resolver struct{}

// Person : Resolver function for the "Person" query
func (r *resolver) Person(args struct{ ID graphql.ID }) *personResolver {
	if p := peopleData[args.ID]; p != nil {
		return &personResolver{p}
	}
	return nil
}

var peopleData = make(map[graphql.ID]*person)

var mainSchema *graphql.Schema

var (
	// ErrNameNotProvided is thrown when a name is not provided
	QueryNameNotProvided = errors.New("no query was provided in the HTTP body")
)

func Handler(context context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

	// If no query is provided in the HTTP request body, throw an error
	if len(request.Body) < 1 {
		return events.APIGatewayProxyResponse{}, QueryNameNotProvided
	}

	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}

	if err := json.Unmarshal([]byte(request.Body), &params); err != nil {
		log.Print("Could not decode body", err)
	}

	response := mainSchema.Exec(context, params.Query, params.OperationName, params.Variables)
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Print("Could not decode body")
	}

	return events.APIGatewayProxyResponse{
		Body:       string(responseJSON),
		StatusCode: 200,
	}, nil

}

func init() {
	for _, p := range people {
		peopleData[p.ID] = p
	}
	mainSchema = graphql.MustParseSchema(Schema, &resolver{})
}

func main() {
	lambda.Start(Handler)
}
