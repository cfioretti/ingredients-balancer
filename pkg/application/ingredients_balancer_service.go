package application

import (
	"errors"
	"math"

	"github.com/cfioretti/ingredients-balancer/pkg/domain"
)

const totalPercentage = 100

type IngredientsBalancerService struct{}

func NewIngredientsBalancerService() *IngredientsBalancerService {
	return &IngredientsBalancerService{}
}

func (bs IngredientsBalancerService) Balance(recipe domain.Recipe, pans domain.Pans) (*domain.RecipeAggregate, error) {
	if pans.TotalArea <= 0 || getFirstIngredientAmount(recipe.Dough.Ingredients) <= 0 {
		return nil, errors.New("invalid dough weight")
	}

	totalDoughWeight := pans.TotalArea / 2
	doughPercentVariation := totalDoughWeight * recipe.Dough.PercentVariation / 100
	doughConversionRatio := (totalDoughWeight + doughPercentVariation) / totalPercentage
	balancedDough := domain.Dough{
		PercentVariation: recipe.Dough.PercentVariation,
		Ingredients:      balanceIngredients(recipe.Dough.Ingredients, doughConversionRatio),
	}

	toppingConversionRatio := pans.TotalArea / recipe.Topping.ReferenceArea
	balancedTopping := domain.Topping{
		ReferenceArea: recipe.Topping.ReferenceArea,
		Ingredients:   balanceIngredients(recipe.Topping.Ingredients, toppingConversionRatio),
	}

	recipeAggregate := &domain.RecipeAggregate{
		Recipe: recipe,
		SplitIngredients: domain.SplitIngredients{
			SplitDough:   calculateSplitDoughs(balancedDough, pans),
			SplitTopping: []domain.Topping{},
		},
	}
	recipeAggregate.Dough = balancedDough
	recipeAggregate.Topping = balancedTopping

	return recipeAggregate, nil
}

func calculateSplitDoughs(totalDough domain.Dough, pans domain.Pans) []domain.Dough {
	var splitDoughs []domain.Dough

	for _, pan := range pans.Pans {
		ratio := pan.Area / pans.TotalArea

		splitDough := domain.Dough{
			Name:        pan.Name,
			Ingredients: make([]domain.Ingredient, len(totalDough.Ingredients)),
		}
		splitDough.Ingredients = balanceIngredients(totalDough.Ingredients, ratio)

		splitDoughs = append(splitDoughs, splitDough)
	}

	return splitDoughs
}

func balanceIngredients(ingredients []domain.Ingredient, ratio float64) []domain.Ingredient {
	balancedIngredients := make([]domain.Ingredient, len(ingredients))
	for i, ingredient := range ingredients {
		balancedIngredients[i] = domain.Ingredient{
			Name:   ingredient.Name,
			Amount: round(ingredient.Amount * ratio),
		}
	}
	return balancedIngredients
}

func getFirstIngredientAmount(ingredients []domain.Ingredient) float64 {
	if len(ingredients) == 0 {
		return 0
	}
	return ingredients[0].Amount
}

func round(num float64) float64 {
	return math.Round(num*10) / 10
}
