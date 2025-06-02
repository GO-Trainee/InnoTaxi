# InnoTaxi - Taxi Ordering Application

## Overview

InnoTaxi is a taxi ordering application designed to facilitate taxi bookings with two primary user roles: **User** and **Driver**. Users can be assigned multiple roles, including **Analyst**, which grants additional permissions to access analytics features. The application supports two wallet types—**Personal** and **Family**—and offers three taxi types: **Economy**, **Comfort**, and **Business**. Built on a microservice architecture with clean architecture principles, it ensures modularity, scalability, and maintainability.

**Important**: All repositories created for this project must remain **private** and are subject to a Non-Disclosure Agreement (NDA). They must **never** be made public.

## Microservices

The application consists of six microservices, each handling specific functionalities:

| Service          | Description                                                                 |
|------------------|-----------------------------------------------------------------------------|
| **User Service** | Manages user profiles, authentication data, and user-related operations.    |
| **Driver Service** | Handles driver profiles, authentication data, and driver status management. |
| **Order Service** | Orchestrates taxi order creation, driver assignment, trip management, and pricing. |
| **Analytic Service** | Provides statistical insights and analytics for users with Analyst permissions. |
| **Wallet Service** | Manages wallets, transactions, and payment processing for users and drivers. |
| **Auth Service** | Centralizes authentication, authorization, and session management for all roles. |

## Functional Requirements

### Auth Service

- **Authentication**: Processes login requests by delegating credential verification to User/Driver Services. Issues JWT tokens upon successful authentication.
- **Role Management**: Issues JWT tokens containing user ID and roles (e.g., ["User", "Analyst", "Driver"]). Users can have multiple roles.
- **Token Management**: Stores active sessions and blacklisted tokens in Redis. Validates tokens by checking signature, expiration, and blacklist status.
- **Session Management**: Provides token refresh mechanism and logout functionality.
- **Access Control**: Services validate tokens through Auth Service to enforce role-based access control.

### User Service

- **Registration**: Users sign up with name, phone number, email, and password. Password is hashed before storage. Default role "User" is assigned.
- **Authentication Support**: Provides credential verification endpoint for Auth Service during login.
- **Profile Management**: Users can view, update (name, phone number, email), change password, or soft-delete their profiles.
- **Role Management**: Admin functionality to assign "Analyst" role to users (manual process via admin interface).
- **Order Management**: 
  - Initiate taxi orders through Order Service
  - View order history and status
  - Rate completed trips
- **Wallet Integration**: Interface with Wallet Service for balance checks and payment operations.
- **Authentication**: All operations require valid JWT token validated through Auth Service.

### Driver Service

- **Registration**: Drivers sign up with name, phone number, email, password, and taxi type. Password is hashed before storage. Role "Driver" is assigned.
- **Authentication Support**: Provides credential verification endpoint for Auth Service during login.
- **Profile Management**: Drivers can view, update profile information, and change password.
- **Status Management**: Drivers can toggle their availability status (free/busy/offline).
- **Trip Management**: 
  - Accept/decline trip requests from Order Service
  - Update trip status (started, completed)
  - Rate completed trips
- **Rating System**: View calculated rating based on last 20 trips.
- **Wallet Integration**: Interface with Wallet Service for balance and earnings.
- **Authentication**: All operations require valid JWT token validated through Auth Service.

### Order Service

- **Order Management**: 
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
- **Authentication**: Role-based access control for different operations.

### Analytic Service

- **Data Collection**: Receives events from other services via Kafka (user registrations, completed orders, ratings).
- **Statistics**: Provides order statistics, user behavior analytics, and business metrics.
- **Ratings Analytics**: Aggregated ratings for drivers and overall service quality.
- **Reports**: Generates business reports for stakeholders.
- **Real-time Dashboards**: Provides real-time metrics for operations monitoring.
- **Authentication**: Access restricted to users with "Analyst" role.

### Wallet Service

- **Wallet Creation**: 
  - Automatically creates personal wallets upon user/driver registration
  - Creates family wallets upon user request
- **Family Wallet Management**: 
  - Adds/removes members by phone number
  - Manages ownership and permissions
- **Transaction Management**: 
  - Processes payments with two-phase commit (reserve → confirm/cancel)
  - Maintains transaction history with detailed status tracking
  - Handles refunds and cancellations
- **Balance Management**: 
  - Real-time balance checking
  - Cash-in operations with payment gateway integration
- **Security**: 
  - Transaction verification and fraud detection
  - Audit logging for all financial operations
- **Authentication**: Internal service authentication for secure inter-service communication.

## Nonfunctional Requirements

### General

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
- **Documentation**: Each service includes comprehensive README, API documentation (Swagger/OpenAPI), and architecture diagrams.
- **Testing**: Unit tests, integration tests, and E2E test suites. Postman collections for API testing.
- **Monitoring**: Prometheus metrics, structured logging, distributed tracing with Jaeger.

### Service-Specific Technical Stack

#### User Service
- **Database**: PostgreSQL for user profiles and authentication data
- **HTTP Framework**: Gin
- **Database Layer**: sqlc for type-safe SQL, squirrel for query building
- **Migrations**: go-migrate
- **Testing**: testify, gomock, dockertest for integration tests
- **Frontend**: Vue.js 3 with Composition API and Pinia for state management

#### Driver Service
- **Database**: PostgreSQL for driver profiles (consistency with User Service)
- **HTTP Framework**: Gin
- **Database Layer**: sqlc and squirrel
- **Migrations**: go-migrate
- **Testing**: testify, gomock, dockertest
- **Frontend**: Vue.js 3 (consistent with User Service)

#### Order Service
- **Database**: PostgreSQL for transactional data
- **Search Engine**: Elasticsearch for advanced search capabilities
- **HTTP Framework**: Gin with GraphQL support for complex queries
- **Message Queue**: Kafka for event publishing
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
- **Database Layer**: sqlc with transaction support
- **Security**: Encryption for sensitive financial data
- **Testing**: testify, financial transaction testing suites

#### Auth Service
- **Database**: Redis for session storage and token blacklisting
- **HTTP Framework**: Gin
- **JWT Library**: golang-jwt with secure key management
- **Testing**: testify, Redis integration tests

## Technical Requirements

- **Inter-Service Communication**: gRPC for synchronous calls, Kafka for asynchronous events
- **API Gateway**: Optional Kong or similar for external API management
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
1. User/Driver sends registration data to respective service
2. Service validates data and stores hashed password
3. Service creates default role assignment
4. Service triggers wallet creation via Wallet Service
5. Service publishes user_registered event to Analytics

### Login
1. Client sends credentials to Auth Service
2. Auth Service delegates verification to User/Driver Service
3. User/Driver Service validates credentials (email/phone + password)
4. On success, Auth Service receives user ID and roles
5. Auth Service issues JWT token with user ID and roles
6. Token is cached in Redis with TTL

### Request Authentication
1. Client includes JWT token in Authorization header
2. Target service validates token via Auth Service
3. Auth Service checks token validity (signature, expiration, blacklist)
4. Auth Service returns user ID and roles if valid
5. Target service enforces role-based permissions

### Logout
1. Client sends logout request to Auth Service
2. Auth Service adds token to blacklist in Redis
3. Token becomes immediately invalid

## Inter-Service Communication Patterns

### Synchronous (gRPC)
- Auth validation requests
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
- Mock external services (if needed)

### Staging Environment
- Kubernetes cluster + helm
- Shared databases
- Full integration testing

### Production Environment
- Kubernetes deployment
- High-availability databases
- Load balancing and auto-scaling
- Monitoring and alerting

## Monitoring and Observability

- **Metrics**: Prometheus with Grafana dashboards
- **Logging**: Structured logging with centralized collection (ELK/EFK)
- **Tracing**: Distributed tracing with Jaeger
- **Alerting**: PagerDuty/Slack integration for critical issues
- **SLA Monitoring**: Uptime and performance SLAs

