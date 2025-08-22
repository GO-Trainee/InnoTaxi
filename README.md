# Taxi Ordering Application

## Overview

Taxi ordering application designed to facilitate taxi bookings with two primary user roles: **User** and **Driver**. Users can be assigned multiple roles. These include the Analyst role, which grants permissions to access analytics features, and the Admin role, which provides full system access and administrative privileges. The application handles finance operations and offers three taxi types: **Economy**, **Comfort**, and **Business**. Built on a microservice architecture with clean architecture principles, it ensures modularity, scalability, and maintainability.

**Important**: All repositories created for this project must remain **private** and are subject to a Non-Disclosure Agreement (NDA). They must **never** be made public.

## Microservices

The application consists of seven microservices, each handling specific functionalities:

| Service            | Description                                                                 |
|--------------------|-----------------------------------------------------------------------------|
| **User Service**   | Centralized management of all user accounts (including drivers), storing authentication data, profile details, and user roles. |
| **Driver Service** | Manages driver-specific data, driver statuses, operational availability, and tasks related to order acceptance, cancellation, and completion — all from the driver's side.. |
| **Order Service**  | Orchestrates taxi order creation, driver assignment, trip lifecycle, and pricing logic. |
| **Analytic Service** | Provides statistical insights and analytics accessible to users with Analyst permissions. |
| **Wallet Service** | Handles digital wallets, transaction history, and payment operations for both riders and drivers. |
| **Auth Service**   | Performs authentication operations and JWT token validation across services. |
| **Gateway Service** | Routes external API requests to internal services, verifies user roles, and enforces access control rules. |


## Functional Requirements

### Auth Service

- **Authentication**: Processes login requests by delegating credential verification to User Service. Issues access and refresh JWT tokens upon successful authentication.
- **Token Management**: 
  - Issues short-lived access tokens (15-30 minutes) containing user ID and roles (e.g., ["User", "Analyst", "Driver", "Admin"])
  - Issues long-lived refresh tokens (7-30 days) for token renewal
  - Stores active tokens in Redis
- **Token Validation**: Validates access tokens by checking signature, expiration, and active status in redis. Returns user ID and roles if token is valid.
- **Token Refresh**: Provides `/refresh` endpoint to exchange valid refresh tokens for new access tokens without requiring re-authentication.
- **Session Management**: Manages token pairs, provides logout functionality that invalidates both access and refresh tokens.

### Gateway Service

- **Request Routing**: Routes incoming requests to appropriate microservices based on URL patterns and business logic.
- **Role-Based Access Control**: Validates user roles and enforces access policies before forwarding requests to target services.
- **Token Validation**: Validates JWT tokens through Auth Service and extracts user information for authorization decisions.
- **Load Balancing**: Distributes requests across multiple instances of backend services.
- **Rate Limiting**: Implements rate limiting policies to prevent abuse and ensure fair usage.
- **Request/Response Transformation**: Handles request/response transformations and API versioning if needed.

### User Service

- **Registration**: Users sign up with name, phone number, email, role and password. Password is hashed before storage. Allowed registration roles: "User", "Driver".
- **Authentication Support**: Provides credential verification endpoint for Auth Service during login.
- **Profile Management**: Users can view, update (name, phone number, email), change password, or soft-delete their profiles.
- **Role Management**: Admin functionality to assign "Analyst" role to users (manual process via admin interface).
- **Rating System**: View calculated rating based on last 20 trips. For both users and drivers.
- **Order Management**: 
  
- **Wallet Integration**: Interface with Wallet Service for balance checks and payment operations.
- **Authentication**: All operations require valid JWT token. Role validation is handled by Gateway Service.

### Driver Service

- **Registration**: The Driver service registers data that is specific only to a driver. It is invoked by the User service when a user with the 'Driver' registration role initiates the process.
- **Profile Management**: Drivers can view, update profile information.
- **Status Management**: Drivers can toggle their availability status (available/on-trip/offline).
- **Trip Management**: 
  - Accept/decline trip requests from Order Service
  - Update trip status (started, completed)
- **Authentication**: All operations require valid JWT token. Role validation is handled by Gateway Service.

### Order Service

- **Order Management**:
- - Initiate taxi orders through Order Service (User only)
  - View order history and status
  - Rate completed trips 
  - Creates and manages taxi orders
  - Handles order lifecycle: created → driver_assigned → in_progress → completed/cancelled
  - Stores complete trip history with ratings and comments
- **Driver Assignment**: 
  - Implements driver matching algorithm based on location, taxi type, and rating
  - Manages driver queue and availability
  - Handles driver acceptance/rejection logic
- **Pricing Engine**: 
  - Calculates trip cost based on distance, time, taxi type, and surge pricing
  - Manages pricing configuration
- **Trip History**: Stores all trip data including routes, timestamps, costs, and ratings.
- **Search & Analytics**: Provides search functionality for users with "Analyst" role.
- **Payment Integration**: Coordinates with Wallet Service for payment processing.
- **Authentication**: All operations require valid JWT token. Role validation is handled by Gateway Service.

### Analytic Service

- **Data Collection**: Receives events from other services via Kafka (user registrations, completed orders, ratings).
- **Statistics**: Provides order statistics, user behavior analytics, and business metrics.
- **Ratings Analytics**: Aggregated ratings for drivers and overall service quality.
- **Reports**: Generates business reports for stakeholders.
- **Real-time Dashboards**: Provides real-time metrics for operations monitoring.
- **Authentication**: Access restricted to users with "Analyst" role. Role validation is handled by Gateway Service.

### Wallet Service

- **Wallet Creation**: 
  - Automatically creates personal wallets upon user/driver registration
- **Transaction Management**: 
  - Processes payments via immediate deduction with compensatory rollback for cancellations.
  - Maintains transaction history with detailed status tracking
  - Handles refunds and cancellations
- **Balance Management**: 
  - Real-time balance checking
  - Cash-in operations with payment gateway integration
- **Security**: 
  - Transaction verification and fraud detection
  - Audit logging for all financial operations
- **Authentication**: Internal service authentication for secure inter-service communication. Role validation is handled by Gateway Service.

## Nonfunctional Requirements

### General

- **Repository Structure**: All microservices are organized in a single monorepo with clear service boundaries and shared libraries.
- **Service Architecture**: Each microservice follows Clean Architecture principles with mandatory 3-layer separation:
  - **Handler Layer**: HTTP handlers, gRPC servers, Kafka consumers, validation
  - **Service Layer**: Business logic, orchestration
  - **Repository Layer**: Data access, database operations
- **Version Control**: All services use GitHub with consistent branching strategy.
- **Branching Strategy**: GitHub Flow with **main** (production) and **develop** (development) branches. Feature branches from **develop**, named `feature/TASK-ID-description`.
- **Pull Requests**: Submit PRs to **develop** with task ID in the name. Include proof of work for frontend PRs. Code review required before merging.
- **CI/CD**: GitHub Actions for all services with pipeline stages:
  - Code quality checks (linting, formatting)
  - Unit and integration tests
  - Security vulnerability scanning
  - Docker image build and push to registry (main branch only)
- **Deployment**: 
  - Development: Docker Compose
  - Production: Kubernetes with Helm charts
- **Configuration**: Environment-based configuration using .env files and Kubernetes ConfigMaps/Secrets.
- **Documentation**: Each service includes comprehensive README, API documentation (Swagger/OpenAPI for REST, protobuf for gRPC/Kafka, GraphQL schemas), and architecture diagrams.
- **Testing**: Unit tests, integration tests, and E2E test suites. Postman collections for API testing.
- **Monitoring**: Prometheus metrics, structured logging, distributed tracing with Jaeger.

### Service-Specific Technical Stack

#### User Service
- **Database**: MongoDB for user profiles and authentication data
- **HTTP Framework**: Gin
- **Database Layer**: mongo-driver for interacting with database
- **Migrations**: golang-migrate
- **Testing**: testify, gomock, dockertest for integration tests
- **Frontend**: Vue.js 3 with Composition API and Pinia for state management

#### Driver Service
- **Database**: MongoDB for driver profiles (consistency with User Service)
- **HTTP Framework**: Gin
- **Database Layer**: mongo-driver
- **Migrations**: golang-migrate
- **Testing**: testify, gomock, dockertest
- **Frontend**: Vue.js 3 (consistent with User Service)

#### Order Service
- **Database**: Elasticsearch for orders data
- **Search Engine**: Elasticsearch for advanced search capabilities
- **HTTP Framework**: Gin with GraphQL support for complex queries
- **Message Queue**: Kafka for event publishing
- **Caching**: Redis for caching recent orders
- **Testing**: testify, gomock, dockertest
- **Frontend**: React with Redux Toolkit for order management interface

#### Analytic Service
- **Database**: ClickHouse for analytical workloads
- **Message Broker**: Kafka consumer for real-time data ingestion
- **HTTP Framework**: Fiber for high-performance analytics API
- **Testing**: testify, integration tests with ClickHouse
- **Frontend**: React with visualization libraries (Chart.js, D3.js)

#### Wallet Service
- **Database**: PostgreSQL for financial data consistency
- **HTTP Framework**: Gin
- **Caching**: Redis for storing recent transactions
- **Database Layer**: sqlx with transaction support
- **Security**: Encryption for sensitive financial data
- **Testing**: testify, financial transaction testing suites

#### Auth Service
- **Database**: Redis for session storage and refresh token management
- **HTTP Framework**: Gin with middleware for token validation
- **JWT Library**: golang-jwt with secure key management for access and refresh tokens
- **Token Storage**: Redis with different TTL strategies for access vs refresh tokens
- **Testing**: testify, Redis integration tests, JWT token validation tests

#### Gateway Service
- **Gateway Platform**: NGINX (used as the primary API gateway)
- **Routing**: Configured via NGINX location blocks and upstreams for directing traffic to appropriate services depending on user role
- **Authentication**: JWT validation via Auth service
- **Rate Limiting**: Rate limiting using NGINX modules (limit_req, e.g.)
- **Load Balancing**: Built-in NGINX load balancing strategies 
- **Testing**: Integration tests using HTTP clients

## Technical Requirements

- **Inter-Service Communication**: gRPC for synchronous calls, Kafka for asynchronous events
- **API Gateway**: Custom Gateway Service handles all external API requests and routing
- **Message Broker**: Kafka for event-driven architecture
- **Containerization**: Docker with multi-stage builds
- **Orchestration**: Kubernetes with proper resource limits and health checks
- **Security**: 
  - TLS for all communications
  - Secrets management with Kubernetes Secrets
  - Regular security audits
- **Code Quality**: 
  - Go linting with golangci-lint
  - Code coverage requirements (>80%)
  - Static analysis tools
- **Performance**: 
  - Database query optimization
  - Caching strategies (Redis)
  - Connection pooling

## Authentication Flow

### Registration
1. Gateway sends registration data to user service
2. User service validates data and stores hashed password
3. User service creates default role assignment
4. Service triggers wallet creation via Wallet Service
5. Service publishes user_registered event to Analytics

### Login
1. Client sends credentials to Auth Service
2. Auth Service delegates verification to User Service
3. User/Driver Service validates credentials (email/phone + password)
4. On success, Auth Service receives user ID and roles
5. Auth Service issues access token (short TTL) and refresh token (long TTL) with user ID and roles
6. Both tokens are cached in Redis with appropriate TTLs
7. Client receives both tokens and stores them securely

### Request Authentication
1. Client sends request with access token in Authorization header to Gateway Service
2. Gateway Service validates access token via Auth Service
3. Auth Service checks access token validity (signature, expiration, redis)
4. Auth Service returns user ID and roles if valid
5. Gateway Service enforces role-based permissions and routes request to target service
6. Target service processes request without additional authentication checks

### Token Refresh
1. Client's access token expires during session
2. Client sends refresh token to Auth Service `/refresh` endpoint
3. Auth Service validates refresh token (signature, expiration, redis)
4. If valid, Auth Service issues new access nad resfresh tokens
5. Client updates stored tokens and retries original request

### Logout
1. Client sends logout request to Auth Service with both access and refresh tokens
2. Auth Service deletes both tokens from redis
3. Both tokens become immediately invalid
4. User session is terminated

## Inter-Service Communication Patterns

### Synchronous (gRPC)
- Gateway to Auth Service for token validation
- Gateway to target services for request forwarding
- Wallet balance checks
- Driver availability queries
- Real-time data requirements

### Asynchronous (Kafka Events)
- User/Driver registration
- Order status changes
- Payment completions
- Analytics data collection

### Database Consistency
- Each service owns its data
- No direct database access between services
- Event-driven eventual consistency
- Saga pattern for distributed transactions


## Error Handling and Resilience

- **Circuit Breaker**: Implement circuit breaker pattern for service-to-service calls
- **Retry Logic**: Exponential backoff for transient failures
- **Graceful Degradation**: Service should handle partial functionality when dependencies are unavailable
- **Health Checks**: Kubernetes liveness and readiness probes
- **Timeout Management**: Appropriate timeouts for all external calls

## Security Considerations

- **Input Validation**: Strict validation on all input data
- **SQL Injection Prevention**: Use parameterized queries
- **XSS Protection**: Proper output encoding
- **Rate Limiting**: API rate limiting to prevent abuse
- **Audit Logging**: Log all security-relevant events
- **Data Encryption**: Encrypt sensitive data at rest and in transit

## Deployment Strategy

### Development Environment
- Docker Compose with all services
- Local databases and Redis
- Mock external services

### Production Environment
- Kubernetes deployment
- High-availability databases
- Load balancing and auto-scaling
- Monitoring and alerting

## Monitoring and Observability

- **Metrics**: Prometheus with Grafana dashboards
- **Logging**: Structured logging with centralized collection (ELK/EFK)
- **Tracing**: Distributed tracing with Jaeger
- **Alerting**: Prometheus integration for critical issues
- **SLA Monitoring**: Uptime and performance SLAs

## Project Implementation Steps

### Phase 1: Architecture Foundation
1. **Project Structure Setup**
    - check file link
    - better to do in drawio

### Phase 2: Database Design
2. **Database Schema Creation**
   - Design PostgreSQL schemas for each service:
     - **User Service**
     - **Driver Service**
     - **Order Service**
     - **Wallet Service**
   - Design Redis schemas for Auth Service
   - Design ClickHouse schemas for Analytic Service
   - Document table relationships and constraints
   - better to do in drawio


### Phase 3: API Contracts Definition
3. **API Documentation**
   - **REST APIs**: Create Swagger/OpenAPI specifications for all HTTP endpoints
   - **gRPC APIs**: Define .proto files for inter-service communication
   - **Kafka Events**: Define .proto files for event schemas
   - **GraphQL**: Create GraphQL schemas for complex query operations
   - Document request/response formats, error codes, and validation rules

### Phase 4: Data Layer Implementation
4. **Repository Layer Development**
   - Implement repository interfaces for each service
   - Setup database connections and connection pooling
   - Implement CRUD operations with proper error handling
   - Setup database migrations and seeders
   - Add database-level tests

### Phase 5: Communication Layer Implementation  
5. **Handler Layer Development**
   - **REST Handlers**: Implement HTTP handlers with request validation
   - **gRPC Servers**: Implement gRPC service implementations
   - **GraphQL Resolvers**: Implement GraphQL query and mutation resolvers
   - **Kafka Consumers**: Implement event consumers with proper error handling
   - **Client Generation**: Generate gRPC clients and Kafka producers

### Phase 6: Business Logic Implementation
6. **Service Layer Development**
   - Implement business logic in service layer
   - Implement inter-service communication patterns
   - Add comprehensive unit tests for business logic
   - Implement transaction management and saga patterns

### Phase 7: System Integration
7. **Microservices Integration**
   - Setup Gateway Service with routing and authentication
   - Configure Auth Service with token management
   - Implement end-to-end request flows
   - Setup monitoring and logging infrastructure
   - Conduct integration testing
   - Performance testing and optimization
   - Deploy to k8s for system testing
