# Ingredients-Balancer Service

**Ingredients-Balancer Service** - Part of PizzaMaker Microservices Architecture - A microservice for pizza ingredient balancing and optimization based on Domain-Driven Design (DDD) architecture with complete observability and monitoring.

## Main Features

- **Ingredient Balancing**: Optimize ingredient distribution across multiple pizza pans
- **Pan Optimization**: Distribute ingredients optimally based on pan sizes and quantities
- **Business Metrics**: Collects domain-specific metrics (balancing accuracy, waste percentage, utilization)

## Technologies

- **Go** - Primary language
- **gRPC** - Service communication
- **Prometheus** - Metrics and monitoring
- **OpenTelemetry + Jaeger** - Distributed tracing
- **Logrus** - Structured logging
- **Docker** - Containerization

## Endpoints

### gRPC Services
- **Port**: 50052 (configurable)
- **Service**: `IngredientsBalancerServer`
- **Methods**: 
  - `Balance(BalanceRequest) -> BalanceResponse`
  - `ValidateRecipe(ValidateRequest) -> ValidateResponse`

### HTTP Endpoints
- **Port**: 8081 (configurable)
- `GET /metrics` - Prometheus metrics
- `GET /health` - Health check

## Observability

### Structured Logging
- **Correlation ID** for cross-service request tracking
- **Structured JSON** for easy parsing
- **Configurable levels** (Debug, Info, Warn, Error)

### Distributed Tracing
- **OpenTelemetry** for instrumentation
- **Jaeger** for trace visualization
- **Automatic spans** for gRPC operations

### Prometheus Metrics
The service exposes both **business** and **technical** metrics:

#### Business Metrics
- `ingredients_balancer_balance_operations_total` - Total balance operations by recipe type
- `ingredients_balancer_balancing_accuracy` - Balancing accuracy percentage
- `ingredients_balancer_ingredient_wastage_percentage` - Ingredient wastage percentage
- `ingredients_balancer_pan_utilization_percentage` - Pan utilization efficiency
- `ingredients_balancer_recipe_portions` - Number of recipe portions processed
- `ingredients_balancer_optimization_strategies_total` - Optimization strategies used
- `ingredients_balancer_quality_checks_total` - Quality validation checks

#### Technical Metrics
- `ingredients_balancer_grpc_requests_total` - Total gRPC requests
- `ingredients_balancer_grpc_request_duration_seconds` - gRPC request duration
- `ingredients_balancer_active_grpc_connections` - Active gRPC connections
- `ingredients_balancer_active_balance_operations` - Active balance operations

#### Processing Metrics
- `ingredients_balancer_ingredient_processing_total` - Ingredient processing operations
- `ingredients_balancer_recipe_analysis_total` - Recipe analysis operations
- `ingredients_balancer_pan_distributions_total` - Pan distribution operations
- `ingredients_balancer_ingredient_optimizations_total` - Ingredient optimizations
