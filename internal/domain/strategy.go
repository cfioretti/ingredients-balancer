package domain

import (
	"github.com/cfioretti/ingredients-balancer/pkg/domain"
)

type Strategy func(data map[string]string) domain.Pan
