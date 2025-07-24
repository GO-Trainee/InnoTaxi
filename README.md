# InnoTaxi - Taxi Ordering Application

## Overview

InnoTaxi is a taxi ordering application designed to facilitate taxi bookings with two primary user roles: **Passenger** and **Driver**. Users can also be assigned **Analyst** role, which grants additional permissions to access analytics features. Each user has only one role. The application supports two wallet types—**Personal** and **Family**—and offers three taxi types: **Economy**, **Comfort**, and **Business**. Built on a microservice architecture with clean architecture principles, it ensures modularity, scalability, and maintainability.

**Important**: All repositories created for this project must remain **private** and are subject to a Non-Disclosure Agreement (NDA). They must **never** be made public.

## Microservices

The application consists of seven microservices, each handling specific functionalities:

| Service          | Description                                                                 |
|------------------|-----------------------------------------------------------------------------|
| **Gateway Service** | Handles request routing, authentication validation, and role-based access control using nginx. |
| **User Service** | Manages all user profiles, authentication data, and role assignments for all users. |
| **Driver Service** | Manages driver-specific information (vehicle, license, status) linked to users via user_id. |
| **Order Service** | Orchestrates taxi order creation, driver assignment, trip management, and pricing. |
| **Analytic Service** | Provides statistical insights and analytics for users with Analyst role. |
| **Wallet Service** | Manages wallets, transactions, and payment processing for users and drivers. |
| **Auth Service** | Handles JWT token generation and validation. Contains minimal user information (ID and role). |

## Functional Requirements

### Gateway Service

- **Request Routing**: Routes incoming requests to appropriate microservices based on URL paths using nginx reverse proxy configuration.
- **Authentication Validation**: Validates JWT tokens by calling Auth Service before forwarding requests to target services.
- **Role-Based Access Control**: Enforces role-based permissions by checking user roles against endpoint requirements before routing requests.
- **Load Balancing**: Distributes requests across multiple service instances for high availability.
- **Rate Limiting**: Implements rate limiting to prevent API abuse and ensure fair usage.
- **Security Headers**: Adds security headers (CORS, CSP, etc.) to all responses.
- **Request/Response Logging**: Logs all incoming requests and responses for monitoring and debugging.

### Auth Service

- **Token Generation**: Issues JWT tokens containing user ID and single role upon successful authentication.
- **Token Validation**: Validates JWT tokens by checking signature, expiration, and blacklist status.
- **Token Structure**: Minimal JWT payload with user_id, role ("Passenger", "Driver", or "Analyst"), and standard claims (exp, iat).
- **Token Refresh**: Provides token refresh mechanism for extending session duration.
- **Logout**: Adds tokens to blacklist in Redis, making them immediately invalid.
- **Authentication Delegation**: Receives authentication requests from Gateway and delegates credential verification to User Service.

### User Service

- **Registration**: All users (passengers, drivers, analysts) register through this service with name, phone number, email, password, and role ("Passenger", "Driver", or "Analyst"). Password is hashed before storage.
- **Authentication Support**: Provides credential verification endpoint for Auth Service during login process.
- **Profile Management**: Users can view and update basic profile information (name, phone number, email), change password, or soft-delete their profiles.
- **User Data API**: Provides user information to other services (Driver Service, Order Service, etc.) via internal APIs.
- **Order Integration**: 
  - Passengers can initiate taxi orders through Order Service
  - View order history and status
  - Rate completed trips
- **Wallet Integration**: Interface with Wallet Service for balance checks and payment operations.
- **Internal API**: All operations are called internally by Gateway Service after authentication and authorization.

### Driver Service

- **Driver Profile Management**: Manages driver-specific information linked to User Service via user_id:
  - Taxi type (Economy, Comfort, Business)
  - Vehicle information (make, model, year, license plate, color)
  - Driver license information (number, expiration date, category)
  - Insurance and registration documents
- **Driver Status Management**: 
  - Toggle availability status (free/busy/offline)
  - Location tracking and updates
  - Working hours and schedule management
- **Trip Management**: 
  - Accept/decline trip requests from Order Service
  - Update trip status (started, en route, completed)
  - Rate passengers after completed trips
- **Rating System**: 
  - Maintains and calculates driver ratings based on trip feedback
  - Stores rating history and statistics
- **Driver Analytics**: 
  - Trip statistics (completed trips, earnings, ratings)
  - Performance metrics and reports
- **Integration with User Service**: 
  - Validates user existence and "Driver" role before creating driver profile
  - Synchronizes basic user data changes
- **Integration with Order Service**: 
  - Provides driver availability and location data
  - Handles trip assignment and status updates
- **Internal API**: All operations are called internally by Gateway Service after authentication and authorization.

### Order Service

- **Order Management**: 
  - Creates and manages taxi orders
  - Handles order lifecycle: created → driver_assigned → in_progress → completed/cancelled
  - Stores complete trip history with ratings and comments
- **Driver Assignment**: 
  - Implements driver matching algorithm based on location, taxi type, and rating
  - Manages driver queue and availability by communicating with Driver Service
  - Handles driver acceptance/rejection logic through Driver Service
- **Pricing Engine**: 
  - Calculates trip cost based on distance, time, taxi type, and surge pricing
  - Manages pricing configuration
- **Trip History**: Stores all trip data including routes, timestamps, costs, and ratings.
- **Search & Analytics**: Provides search functionality for users with "Analyst" role.
- **Payment Integration**: Coordinates with Wallet Service for payment processing.
- **Internal API**: All operations are called internally by Gateway Service after authentication and authorization.

### Analytic Service

- **Data Collection**: Receives events from other services via Kafka (user registrations, completed orders, ratings).
- **Statistics**: Provides order statistics, user behavior analytics, and business metrics.
- **Ratings Analytics**: Aggregated ratings for drivers and overall service quality.
- **Reports**: Generates business reports for stakeholders.
- **Real-time Dashboards**: Provides real-time metrics for operations monitoring.
- **Internal API**: All operations are called internally by Gateway Service after role validation (restricted to "Analyst" role).

### Wallet Service

- **Wallet Creation**: 
  - Automatically creates personal wallets upon user registration (Passenger and Driver Role)
  - Creates family wallets upon user request
- **Family Wallet Management**: 
  - Adds/removes members by phone number
  - Manages ownership and permissions
- **Transaction Management**: 
  - Maintains transaction history with detailed status tracking
  - Handles refunds and cancellations
- **Balance Management**: 
  - Real-time balance checking
  - Cash-in operations
- **Security**: 
  - Audit logging for all financial operations
- **Internal API**: All operations are called internally by Gateway Service or other services after proper authentication.

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

#### Gateway Service
- **Reverse Proxy**: nginx with Lua scripting for custom logic
- **Configuration**: nginx.conf with upstream definitions and location blocks
- **Authentication**: JWT validation via Auth Service HTTP calls
- **Load Balancing**: nginx upstream with health checks
- **Monitoring**: nginx access logs and metrics
- **Security**: Rate limiting, CORS, security headers

#### User Service
- **Database**: MongoDB for all user profiles and authentication data
- **HTTP Framework**: Gin
- **Database Layer**: MongoDB Go driver (go.mongodb.org/mongo-driver) for database operations
- **Migrations**: MongoDB migrations with custom migration scripts
- **Testing**: testify, gomock, dockertest for integration tests
- **Frontend**: Vue.js 3 with Composition API and Pinia for state management

#### Driver Service
- **Database**: MongoDB for driver-specific information (linked to User Service via user_id)
- **HTTP Framework**: Gin
- **Database Layer**: MongoDB Go driver (go.mongodb.org/mongo-driver) for database operations
- **Migrations**: MongoDB migrations with custom migration scripts
- **Testing**: testify, gomock, dockertest for integration tests
- **Location Services**: Integration with mapping services for location tracking
- **Frontend**: Vue.js 3 for driver management interface

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

- **Inter-Service Communication**: HTTP/gRPC for synchronous calls, Kafka for asynchronous events
- **API Gateway**: Mandatory nginx-based Gateway Service for all external API access
- **Message Broker**: Kafka for event-driven architecture
- **Containerization**: Docker with multi-stage builds
- **Orchestration**: Kubernetes with proper resource limits and health checks
- **Security**: 
  - TLS for all communications
  - JWT-based authentication handled by Gateway Service
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
  - nginx-based load balancing and caching

## Authentication Flow

### Registration
1. Client sends registration data (including role: "Passenger" or "Driver") to Gateway Service
2. Gateway Service routes request to User Service
3. User Service validates data and stores hashed password with role
4. If role is "Driver", User Service triggers Driver Service to create driver profile with user_id
5. Driver Service stores driver-specific information (vehicle, license, etc.)
6. User Service triggers wallet creation via Wallet Service
7. User Service publishes user_registered event to Analytics

### Login
1. Client sends credentials to Gateway Service
2. Gateway Service routes login request to Auth Service
3. Auth Service delegates credential verification to User Service
4. User Service validates credentials (email/phone + password)
5. On success, Auth Service receives user ID and single role
6. Auth Service issues JWT token with user ID and role
7. Token is cached in Redis with TTL

### Request Authentication & Authorization
1. Client includes JWT token in Authorization header for all requests
2. Gateway Service validates token via Auth Service
3. Auth Service checks token validity (signature, expiration, blacklist)
4. Auth Service returns user ID and role if valid
5. Gateway Service enforces role-based permissions for the requested endpoint
6. If authorized, Gateway Service routes request to target service
7. Target service processes request without additional authentication checks

### Logout
1. Client sends logout request to Gateway Service
2. Gateway Service routes request to Auth Service
3. Auth Service adds token to blacklist in Redis
4. Token becomes immediately invalid across all services

## Inter-Service Communication Patterns

### External Client Communication
- All client requests go through Gateway Service
- Gateway handles authentication via Auth Service
- Gateway enforces role-based access control
- Gateway routes requests to appropriate services

### Synchronous (gRPC/HTTP)
- Gateway to Auth Service for token validation
- Gateway to target services for request routing
- Service-to-service calls for business logic:
  - Order → Driver Service (driver availability, location, trip assignment)
  - Order → User Service (user data, passenger information)
  - Order → Wallet Service (payment processing)
  - User → Driver Service (driver profile creation for new drivers)
  - Driver → User Service (user validation and basic data)
- Real-time data requirements

### Asynchronous (Kafka Events)
- User registration events
- Driver profile creation/updates
- Driver status changes (availability, location)
- Order status changes
- Trip completion events
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

