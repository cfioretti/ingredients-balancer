package application

import (
	"errors"

	"github.com/cfioretti/ingredients-balancer/internal/domain"
	"github.com/cfioretti/ingredients-balancer/internal/domain/strategies"
)

type DoughCalculatorService struct{}

func NewIngredientsCalculatorService() *DoughCalculatorService {
	return &DoughCalculatorService{}
}

type Input struct {
	Pans []PanInput `json:"pans"`
}

type PanInput struct {
	Shape    string                 `json:"shape"`
	Measures map[string]interface{} `json:"measures"`
}

func (dc DoughCalculatorService) TotalDoughWeightByPans(body domain.Pans) (*domain.Pans, error) {
	var result domain.Pans
	for _, item := range body.Pans {
		strategy, err := strategies.GetStrategy(item.Shape)
		if err != nil {
			return nil, errors.New("unsupported shape")
		}

		pan, err := strategy.Calculate(item.Measures)
		if err != nil {
			return nil, errors.New("error processing pan")
		}

		result.Pans = append(result.Pans, pan)
		result.TotalArea += pan.Area
	}
	return &result, nil
}
