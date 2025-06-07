package test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/cfioretti/ingredients-balancer/pkg/application"
	grpcServer "github.com/cfioretti/ingredients-balancer/pkg/infrastructure/grpc"
	pb "github.com/cfioretti/ingredients-balancer/pkg/infrastructure/grpc/proto/generated"
)

func TestIngredientsBalancerIntegration(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)
	port := lis.Addr().(*net.TCPAddr).Port
	lis.Close()

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
			t.Errorf("Failed to listen: %v", err)
			return
		}

		ingredientsBalancerService := application.NewIngredientsBalancerService()
		server := grpcServer.NewServer(ingredientsBalancerService)
		grpcNewServer := grpc.NewServer()
		pb.RegisterIngredientsBalancerServer(grpcNewServer, server)

		if err := grpcNewServer.Serve(lis); err != nil {
			t.Errorf("Failed to serve: %v", err)
		}
	}()
	time.Sleep(100 * time.Millisecond)

	conn, err := grpc.Dial(
		fmt.Sprintf("localhost:%d", port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewIngredientsBalancerClient(conn)

	// Create a test recipe
	recipeUUID := uuid.New()
	recipe := &pb.Recipe{
		Id:          1,
		Uuid:        recipeUUID.String(),
		Name:        "Test Recipe",
		Description: "Test Description",
		Author:      "Test Author",
		Dough: &pb.Dough{
			Name:             "Test Dough",
			PercentVariation: 10.0,
			Ingredients: []*pb.Ingredient{
				{
					Name:   "Flour",
					Amount: 100.0,
				},
				{
					Name:   "Water",
					Amount: 65.0,
				},
			},
		},
		Topping: &pb.Topping{
			Name:          "Test Topping",
			ReferenceArea: 1000.0,
			Ingredients: []*pb.Ingredient{
				{
					Name:   "Tomato Sauce",
					Amount: 200.0,
				},
				{
					Name:   "Mozzarella",
					Amount: 150.0,
				},
			},
		},
		Steps: &pb.Steps{
			RecipeId: 1,
			Steps: []*pb.Step{
				{
					Id:          1,
					StepNumber:  1,
					Description: "Mix flour and water",
				},
				{
					Id:          2,
					StepNumber:  2,
					Description: "Knead the dough",
				},
			},
		},
	}

	// Create test pans
	diameter := int32(28)
	pans := &pb.Pans{
		Pans: []*pb.Pan{
			{
				Shape: "round",
				Measures: &pb.Measures{
					Diameter: &diameter,
				},
				Name: "round 28 cm",
				Area: 615.75,
			},
		},
		TotalArea: 615.75,
	}

	// Create the request
	request := &pb.BalanceRequest{
		Recipe: recipe,
		Pans:   pans,
	}

	// Call the Balance method
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.Balance(ctx, request)
	require.NoError(t, err)

	// Assert the response
	assert.NotNil(t, response)
	assert.NotNil(t, response.RecipeAggregate)
	assert.NotNil(t, response.RecipeAggregate.Recipe)
	assert.Equal(t, recipe.Id, response.RecipeAggregate.Recipe.Id)
	assert.Equal(t, recipe.Name, response.RecipeAggregate.Recipe.Name)
	assert.NotNil(t, response.RecipeAggregate.SplitIngredients)
	assert.NotEmpty(t, response.RecipeAggregate.SplitIngredients.SplitDough)
}
