package grpc

import (
	"context"

	"github.com/google/uuid"

	"github.com/cfioretti/ingredients-balancer/pkg/domain"
	pb "github.com/cfioretti/ingredients-balancer/pkg/infrastructure/grpc/proto/generated"
)

type BalancerService interface {
	Balance(context.Context, domain.Recipe, domain.Pans) (*domain.RecipeAggregate, error)
}

type Server struct {
	pb.UnimplementedIngredientsBalancerServer
	ingredientsBalancerService BalancerService
}

func NewServer(ingredientsBalancerService BalancerService) *Server {
	return &Server{
		ingredientsBalancerService: ingredientsBalancerService,
	}
}

func (s *Server) Balance(ctx context.Context, req *pb.BalanceRequest) (*pb.BalanceResponse, error) {
	recipe := toDomainRecipe(req.GetRecipe())
	pans := toDomainPans(req.GetPans())

	result, err := s.ingredientsBalancerService.Balance(ctx, recipe, pans)
	if err != nil {
		return nil, err
	}

	responseProto := toProtoRecipeAggregate(result)

	return &pb.BalanceResponse{
		RecipeAggregate: responseProto,
	}, nil
}

func toDomainRecipe(protoRecipe *pb.Recipe) domain.Recipe {
	recipeUUID, _ := uuid.Parse(protoRecipe.Uuid)

	return domain.Recipe{
		Id:          int(protoRecipe.Id),
		Uuid:        recipeUUID,
		Name:        protoRecipe.Name,
		Description: protoRecipe.Description,
		Author:      protoRecipe.Author,
		Dough:       toDomainDough(protoRecipe.Dough),
		Topping:     toDomainTopping(protoRecipe.Topping),
		Steps:       toDomainSteps(protoRecipe.Steps),
	}
}

func toDomainDough(protoDough *pb.Dough) domain.Dough {
	return domain.Dough{
		Name:             protoDough.Name,
		PercentVariation: protoDough.PercentVariation,
		Ingredients:      toDomainIngredients(protoDough.Ingredients),
	}
}

func toDomainTopping(protoTopping *pb.Topping) domain.Topping {
	return domain.Topping{
		Name:          protoTopping.Name,
		ReferenceArea: protoTopping.ReferenceArea,
		Ingredients:   toDomainIngredients(protoTopping.Ingredients),
	}
}

func toDomainIngredients(protoIngredients []*pb.Ingredient) []domain.Ingredient {
	ingredients := make([]domain.Ingredient, 0, len(protoIngredients))
	for _, protoIngredient := range protoIngredients {
		ingredients = append(ingredients, domain.Ingredient{
			Name:   protoIngredient.Name,
			Amount: protoIngredient.Amount,
		})
	}
	return ingredients
}

func toDomainSteps(protoSteps *pb.Steps) domain.Steps {
	steps := make([]domain.Step, 0, len(protoSteps.Steps))
	for _, protoStep := range protoSteps.Steps {
		steps = append(steps, domain.Step{
			Id:          int(protoStep.Id),
			StepNumber:  int(protoStep.StepNumber),
			Description: protoStep.Description,
		})
	}
	return domain.Steps{
		RecipeId: int(protoSteps.RecipeId),
		Steps:    steps,
	}
}

func toDomainPans(protoPans *pb.Pans) domain.Pans {
	pans := make([]domain.Pan, 0, len(protoPans.Pans))
	for _, protoPan := range protoPans.Pans {
		pans = append(pans, domain.Pan{
			Shape: protoPan.Shape,
			Measures: domain.Measures{
				Diameter: toPointer(protoPan.Measures.Diameter),
				Edge:     toPointer(protoPan.Measures.Edge),
				Width:    toPointer(protoPan.Measures.Width),
				Length:   toPointer(protoPan.Measures.Length),
			},
			Name: protoPan.Name,
			Area: protoPan.Area,
		})
	}
	return domain.Pans{
		Pans:      pans,
		TotalArea: protoPans.TotalArea,
	}
}

func toProtoRecipeAggregate(domainRecipeAggregate *domain.RecipeAggregate) *pb.RecipeAggregate {
	return &pb.RecipeAggregate{
		Recipe:           toProtoRecipe(domainRecipeAggregate.Recipe),
		SplitIngredients: toProtoSplitIngredients(domainRecipeAggregate.SplitIngredients),
	}
}

func toProtoRecipe(domainRecipe domain.Recipe) *pb.Recipe {
	return &pb.Recipe{
		Id:          int32(domainRecipe.Id),
		Uuid:        domainRecipe.Uuid.String(),
		Name:        domainRecipe.Name,
		Description: domainRecipe.Description,
		Author:      domainRecipe.Author,
		Dough:       toProtoDough(domainRecipe.Dough),
		Topping:     toProtoTopping(domainRecipe.Topping),
		Steps:       toProtoSteps(domainRecipe.Steps),
	}
}

func toProtoDough(domainDough domain.Dough) *pb.Dough {
	return &pb.Dough{
		Name:             domainDough.Name,
		PercentVariation: domainDough.PercentVariation,
		Ingredients:      toProtoIngredients(domainDough.Ingredients),
	}
}

func toProtoTopping(domainTopping domain.Topping) *pb.Topping {
	return &pb.Topping{
		Name:          domainTopping.Name,
		ReferenceArea: domainTopping.ReferenceArea,
		Ingredients:   toProtoIngredients(domainTopping.Ingredients),
	}
}

func toProtoIngredients(domainIngredients []domain.Ingredient) []*pb.Ingredient {
	protoIngredients := make([]*pb.Ingredient, 0, len(domainIngredients))
	for _, domainIngredient := range domainIngredients {
		protoIngredients = append(protoIngredients, &pb.Ingredient{
			Name:   domainIngredient.Name,
			Amount: domainIngredient.Amount,
		})
	}
	return protoIngredients
}

func toProtoSteps(domainSteps domain.Steps) *pb.Steps {
	protoSteps := make([]*pb.Step, 0, len(domainSteps.Steps))
	for _, domainStep := range domainSteps.Steps {
		protoSteps = append(protoSteps, &pb.Step{
			Id:          int32(domainStep.Id),
			StepNumber:  int32(domainStep.StepNumber),
			Description: domainStep.Description,
		})
	}
	return &pb.Steps{
		RecipeId: int32(domainSteps.RecipeId),
		Steps:    protoSteps,
	}
}

func toProtoSplitIngredients(domainSplitIngredients domain.SplitIngredients) *pb.SplitIngredients {
	protoSplitDoughs := make([]*pb.Dough, 0, len(domainSplitIngredients.SplitDough))
	for _, domainDough := range domainSplitIngredients.SplitDough {
		protoSplitDoughs = append(protoSplitDoughs, toProtoDough(domainDough))
	}

	protoSplitToppings := make([]*pb.Topping, 0, len(domainSplitIngredients.SplitTopping))
	for _, domainTopping := range domainSplitIngredients.SplitTopping {
		protoSplitToppings = append(protoSplitToppings, toProtoTopping(domainTopping))
	}

	return &pb.SplitIngredients{
		SplitDough:   protoSplitDoughs,
		SplitTopping: protoSplitToppings,
	}
}

func toPointer(value *int32) *int {
	if value == nil {
		return nil
	}
	val := int(*value)
	return &val
}
