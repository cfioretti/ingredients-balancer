package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/cfioretti/ingredients-balancer/pkg/domain"
	pb "github.com/cfioretti/ingredients-balancer/pkg/infrastructure/grpc/proto/generated"
)

type MockIngredientsBalancerService struct {
	mock.Mock
}

func (m *MockIngredientsBalancerService) Balance(recipe domain.Recipe, pans domain.Pans) (*domain.RecipeAggregate, error) {
	args := m.Called(recipe, pans)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RecipeAggregate), args.Error(1)
}

func TestNewServer(t *testing.T) {
	mockService := &MockIngredientsBalancerService{}
	server := NewServer(mockService)

	assert.NotNil(t, server)
	assert.Equal(t, mockService, server.ingredientsBalancerService)
}

func TestServer_Balance_Success(t *testing.T) {
	// Setup
	mockService := &MockIngredientsBalancerService{}
	server := NewServer(mockService)

	recipeUUID := uuid.New()
	protoRequest := &pb.BalanceRequest{
		Recipe: &pb.Recipe{
			Id:          1,
			Uuid:        recipeUUID.String(),
			Name:        "Pizza Margherita",
			Description: "Classica pizza italiana",
			Author:      "Chef Mario",
			Dough: &pb.Dough{
				Name:             "Impasto base",
				PercentVariation: 10.5,
				Ingredients: []*pb.Ingredient{
					{Name: "Farina", Amount: 1000},
					{Name: "Acqua", Amount: 600},
				},
			},
			Topping: &pb.Topping{
				Name:          "Condimento Margherita",
				ReferenceArea: 300,
				Ingredients: []*pb.Ingredient{
					{Name: "Pomodoro", Amount: 200},
					{Name: "Mozzarella", Amount: 150},
				},
			},
			Steps: &pb.Steps{
				RecipeId: 1,
				Steps: []*pb.Step{
					{Id: 1, StepNumber: 1, Description: "Impastare"},
					{Id: 2, StepNumber: 2, Description: "Lievitare"},
				},
			},
		},
		Pans: &pb.Pans{
			TotalArea: 600,
			Pans: []*pb.Pan{
				{
					Shape: "circular",
					Name:  "Teglia rotonda",
					Area:  300,
					Measures: &pb.Measures{
						Diameter: int32Ptr(30),
					},
				},
			},
		},
	}

	expectedDomainRecipe := domain.Recipe{
		Id:          1,
		Uuid:        recipeUUID,
		Name:        "Pizza Margherita",
		Description: "Classica pizza italiana",
		Author:      "Chef Mario",
		Dough: domain.Dough{
			Name:             "Impasto base",
			PercentVariation: 10.5,
			Ingredients: []domain.Ingredient{
				{Name: "Farina", Amount: 1000},
				{Name: "Acqua", Amount: 600},
			},
		},
		Topping: domain.Topping{
			Name:          "Condimento Margherita",
			ReferenceArea: 300,
			Ingredients: []domain.Ingredient{
				{Name: "Pomodoro", Amount: 200},
				{Name: "Mozzarella", Amount: 150},
			},
		},
		Steps: domain.Steps{
			RecipeId: 1,
			Steps: []domain.Step{
				{Id: 1, StepNumber: 1, Description: "Impastare"},
				{Id: 2, StepNumber: 2, Description: "Lievitare"},
			},
		},
	}

	expectedDomainPans := domain.Pans{
		TotalArea: 600,
		Pans: []domain.Pan{
			{
				Shape: "circular",
				Name:  "Teglia rotonda",
				Area:  300,
				Measures: domain.Measures{
					Diameter: intPtr(30),
				},
			},
		},
	}

	mockResult := &domain.RecipeAggregate{
		Recipe: expectedDomainRecipe,
		SplitIngredients: domain.SplitIngredients{
			SplitDough: []domain.Dough{
				{
					Name:             "Impasto base - Teglia 1",
					PercentVariation: 10.5,
					Ingredients: []domain.Ingredient{
						{Name: "Farina", Amount: 500},
						{Name: "Acqua", Amount: 300},
					},
				},
			},
			SplitTopping: []domain.Topping{
				{
					Name:          "Condimento Margherita - Teglia 1",
					ReferenceArea: 300,
					Ingredients: []domain.Ingredient{
						{Name: "Pomodoro", Amount: 100},
						{Name: "Mozzarella", Amount: 75},
					},
				},
			},
		},
	}

	mockService.On("Balance", expectedDomainRecipe, expectedDomainPans).Return(mockResult, nil)

	// Execute
	response, err := server.Balance(context.Background(), protoRequest)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotNil(t, response.RecipeAggregate)
	assert.Equal(t, "Pizza Margherita", response.RecipeAggregate.Recipe.Name)
	assert.Len(t, response.RecipeAggregate.SplitIngredients.SplitDough, 1)
	assert.Len(t, response.RecipeAggregate.SplitIngredients.SplitTopping, 1)

	mockService.AssertExpectations(t)
}

func TestServer_Balance_ServiceError(t *testing.T) {
	// Setup
	mockService := &MockIngredientsBalancerService{}
	server := NewServer(mockService)

	protoRequest := &pb.BalanceRequest{
		Recipe: &pb.Recipe{
			Id:   1,
			Uuid: uuid.New().String(),
			Name: "Test Recipe",
			Dough: &pb.Dough{
				Name:        "Test Dough",
				Ingredients: []*pb.Ingredient{},
			},
			Topping: &pb.Topping{
				Name:        "Test Topping",
				Ingredients: []*pb.Ingredient{},
			},
			Steps: &pb.Steps{},
		},
		Pans: &pb.Pans{
			Pans: []*pb.Pan{},
		},
	}

	expectedError := errors.New("servizio non disponibile")
	mockService.On("Balance", mock.AnythingOfType("domain.Recipe"), mock.AnythingOfType("domain.Pans")).Return(nil, expectedError)

	// Execute
	response, err := server.Balance(context.Background(), protoRequest)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, expectedError, err)

	mockService.AssertExpectations(t)
}

func TestToDomainRecipe(t *testing.T) {
	recipeUUID := uuid.New()
	protoRecipe := &pb.Recipe{
		Id:          1,
		Uuid:        recipeUUID.String(),
		Name:        "Test Recipe",
		Description: "Test Description",
		Author:      "Test Author",
		Dough: &pb.Dough{
			Name:             "Test Dough",
			PercentVariation: 5.0,
			Ingredients: []*pb.Ingredient{
				{Name: "Ingredient1", Amount: 100},
			},
		},
		Topping: &pb.Topping{
			Name:          "Test Topping",
			ReferenceArea: 200,
			Ingredients: []*pb.Ingredient{
				{Name: "Ingredient2", Amount: 50},
			},
		},
		Steps: &pb.Steps{
			RecipeId: 1,
			Steps: []*pb.Step{
				{Id: 1, StepNumber: 1, Description: "Step 1"},
			},
		},
	}

	result := toDomainRecipe(protoRecipe)

	assert.Equal(t, 1, result.Id)
	assert.Equal(t, recipeUUID, result.Uuid)
	assert.Equal(t, "Test Recipe", result.Name)
	assert.Equal(t, "Test Description", result.Description)
	assert.Equal(t, "Test Author", result.Author)
	assert.Equal(t, "Test Dough", result.Dough.Name)
	assert.Equal(t, 5.0, result.Dough.PercentVariation)
	assert.Len(t, result.Dough.Ingredients, 1)
	assert.Equal(t, "Test Topping", result.Topping.Name)
	assert.Len(t, result.Steps.Steps, 1)
}

func TestToDomainPans(t *testing.T) {
	protoPans := &pb.Pans{
		TotalArea: 500,
		Pans: []*pb.Pan{
			{
				Shape: "rectangular",
				Name:  "Teglia rettangolare",
				Area:  250,
				Measures: &pb.Measures{
					Width:  int32Ptr(20),
					Length: int32Ptr(30),
				},
			},
			{
				Shape: "circular",
				Name:  "Teglia rotonda",
				Area:  250,
				Measures: &pb.Measures{
					Diameter: int32Ptr(25),
				},
			},
		},
	}

	result := toDomainPans(protoPans)

	assert.Equal(t, 500.0, result.TotalArea)
	assert.Len(t, result.Pans, 2)

	// Test prima teglia (rettangolare)
	assert.Equal(t, "rectangular", result.Pans[0].Shape)
	assert.Equal(t, "Teglia rettangolare", result.Pans[0].Name)
	assert.Equal(t, 250.0, result.Pans[0].Area)
	assert.Equal(t, 20, *result.Pans[0].Measures.Width)
	assert.Equal(t, 30, *result.Pans[0].Measures.Length)
	assert.Nil(t, result.Pans[0].Measures.Diameter)

	// Test seconda teglia (rotonda)
	assert.Equal(t, "circular", result.Pans[1].Shape)
	assert.Equal(t, "Teglia rotonda", result.Pans[1].Name)
	assert.Equal(t, 25, *result.Pans[1].Measures.Diameter)
	assert.Nil(t, result.Pans[1].Measures.Width)
	assert.Nil(t, result.Pans[1].Measures.Length)
}

func TestToProtoRecipeAggregate(t *testing.T) {
	recipeUUID := uuid.New()
	domainAggregate := &domain.RecipeAggregate{
		Recipe: domain.Recipe{
			Id:          1,
			Uuid:        recipeUUID,
			Name:        "Test Recipe",
			Description: "Test Description",
			Author:      "Test Author",
			Dough: domain.Dough{
				Name:        "Test Dough",
				Ingredients: []domain.Ingredient{{Name: "Farina", Amount: 1000}},
			},
			Topping: domain.Topping{
				Name:        "Test Topping",
				Ingredients: []domain.Ingredient{{Name: "Pomodoro", Amount: 200}},
			},
			Steps: domain.Steps{
				RecipeId: 1,
				Steps:    []domain.Step{{Id: 1, StepNumber: 1, Description: "Test Step"}},
			},
		},
		SplitIngredients: domain.SplitIngredients{
			SplitDough: []domain.Dough{
				{Name: "Split Dough 1", Ingredients: []domain.Ingredient{{Name: "Farina", Amount: 500}}},
			},
			SplitTopping: []domain.Topping{
				{Name: "Split Topping 1", Ingredients: []domain.Ingredient{{Name: "Pomodoro", Amount: 100}}},
			},
		},
	}

	result := toProtoRecipeAggregate(domainAggregate)

	assert.NotNil(t, result)
	assert.NotNil(t, result.Recipe)
	assert.Equal(t, int32(1), result.Recipe.Id)
	assert.Equal(t, recipeUUID.String(), result.Recipe.Uuid)
	assert.Equal(t, "Test Recipe", result.Recipe.Name)
	assert.NotNil(t, result.SplitIngredients)
	assert.Len(t, result.SplitIngredients.SplitDough, 1)
	assert.Len(t, result.SplitIngredients.SplitTopping, 1)
}

func TestToPointer(t *testing.T) {
	// Test con valore non nil
	value := int32(42)
	result := toPointer(&value)
	assert.NotNil(t, result)
	assert.Equal(t, 42, *result)

	// Test con valore nil
	result = toPointer(nil)
	assert.Nil(t, result)
}

func TestToDomainIngredients(t *testing.T) {
	protoIngredients := []*pb.Ingredient{
		{Name: "Farina", Amount: 1000},
		{Name: "Acqua", Amount: 600},
		{Name: "Sale", Amount: 20},
	}

	result := toDomainIngredients(protoIngredients)

	assert.Len(t, result, 3)
	assert.Equal(t, "Farina", result[0].Name)
	assert.Equal(t, 1000.0, result[0].Amount)
	assert.Equal(t, "Acqua", result[1].Name)
	assert.Equal(t, 600.0, result[1].Amount)
	assert.Equal(t, "Sale", result[2].Name)
	assert.Equal(t, 20.0, result[2].Amount)
}

func TestToProtoIngredients(t *testing.T) {
	domainIngredients := []domain.Ingredient{
		{Name: "Farina", Amount: 1000},
		{Name: "Acqua", Amount: 600},
	}

	result := toProtoIngredients(domainIngredients)

	assert.Len(t, result, 2)
	assert.Equal(t, "Farina", result[0].Name)
	assert.Equal(t, 1000.0, result[0].Amount)
	assert.Equal(t, "Acqua", result[1].Name)
	assert.Equal(t, 600.0, result[1].Amount)
}

func int32Ptr(v int32) *int32 {
	return &v
}

func intPtr(v int) *int {
	return &v
}
