# InnoTaxi - Taxi Ordering Application

## Overview

InnoTaxi is a taxi ordering application designed to facilitate taxi bookings with two primary user roles: **User** and **Driver**. Users can be assigned multiple roles, including **Analyst**, which grants additional permissions to access analytics features. The application supports two wallet types—**Personal** and **Family**—and offers three taxi types: **Economy**, **Comfort**, and **Business**. Built on a microservice architecture with clean architecture principles, it ensures modularity, scalability, and maintainability.

**Important**: All repositories created for this project must remain **private** and are subject to a Non-Disclosure Agreement (NDA). They must **never** be made public.

## Microservices

The application consists of six microservices, each handling specific functionalities:

| Service          | Description                                                                 |
|------------------|-----------------------------------------------------------------------------|
| **User Service** | Manages user profiles, taxi orders, and wallet interactions for users.       |
| **Driver Service** | Handles driver profiles, trip statuses, and wallet operations for drivers. |
| **Order Service** | Orchestrates taxi order creation, driver assignment, and transaction processing. |
| **Analytic Service** | Provides statistical insights and ratings for users with Analyst permissions. |
| **Wallet Service** | Manages wallets and transactions for users and drivers.                     |
| **Auth Service** | Centralizes authentication and authorization for all roles.                 |

## Functional Requirements

### Auth Service

- **Role Management**: Issues JWT tokens containing a list of user roles (e.g., ["User", "Analyst", "Driver"]). Users can have multiple roles, enabling flexible permission management.
- **Authentication**: Processes login and logout requests for users and drivers. Verifies credentials by interacting with User or Driver Services.
- **Token Management**: Stores blacklisted tokens in Redis for logout functionality. Validates tokens by checking signature, expiration, and blacklist status.
- **Access Control**: Provides role information upon token validation. Services enforce role-based access control, ensuring users cannot access driver or analyst functionalities without appropriate roles (e.g., "Permission Denied" for unauthorized actions).

### User Service

- **Registration**: Users sign up with name, phone number, email, and password. A personal wallet is created via the Wallet Service. Users can be assigned roles such as "User" and optionally "Analyst" during registration or through administrative actions.
- **Profile Management**: Users can view, update (name, phone number, email), or soft-delete their profiles.
- **Wallet Operations**:
  - View available wallets (personal and family) through the Wallet Service.
  - Create family wallets and add members by phone number.
  - Cash in personal or family wallets (only the owner for family wallets).
  - View transaction history (restricted to owners for family wallets).
- **Order Taxi**:
  - Specify taxi type, start and end locations, and select a wallet.
  - Verifies wallet balance via the Wallet Service.
  - Requests a free driver through the Order Service. If no drivers are available, users join a queue with a configurable wait time. If no driver is found, a rejection response is sent.
- **Rate Trip**: Rate the last trip (1–5) with an optional comment, if within a configurable time since completion.
- **View Trips**: View past trips, including taxi type, driver, and route details.
- **Restriction**: Orders are blocked if the selected wallet has insufficient funds.
- **Authentication**: All operations require a valid JWT token, validated through the Auth Service to confirm the "User" role. Users with the "Analyst" role can access additional analytics features.

### Driver Service

- **Registration**: Drivers sign up with name, phone number, email, password, and taxi type. A single driver wallet is created via the Wallet Service. Drivers are assigned the "Driver" role.
- **Status Management**: Drivers toggle their status (free/busy) after a trip, updating the order status via the Order Service.
- **Rate Trip**: Rate the last trip (1–5), if within a configurable time since completion.
- **View Rating**: View rating calculated from the last 20 trips.
- **View Trips**: View past trips, including taxi type, user, and route details.
- **Wallet Operations**: View single wallet balance and transaction history via the Wallet Service.
- **Authentication**: All operations require a valid JWT token, validated through the Auth Service to confirm the "Driver" role.

### Order Service

- **Order Orchestration**: Manages taxi order creation, driver assignment, and status updates. Interacts with the Wallet Service to block and complete funds.
- **Order List**: Provides a filtered and paginated list of orders for authorized roles.
- **Order Fields**: Includes user, driver, start and end locations, taxi type, date, status (in progress, finished), and comment.
- **Pricing**: Maintains pricing for each taxi type.
- **Search**: Users with the "Analyst" role can search orders with partial field matching.
- **Driver Selection**: Supports selecting drivers based on user ratings.
- **Authentication**: Requests are validated through the Auth Service to ensure role-based access (e.g., "User" for ordering, "Driver" for status updates, "Analyst" for searches).

### Analytic Service

- **Statistics**: Users with the "Analyst" role can view order statistics (e.g., counts by day or month).
- **Ratings**: View ratings for all drivers and users.
- **Account**: Analyst permissions are granted to specific users via the "Analyst" role, managed through the User Service.
- **Data Recording**: Records all registrations and completed orders.
- **Authentication**: Access requires a valid JWT token with the "Analyst" role, validated through the Auth Service.

### Wallet Service

- **Wallet Creation**:
  - Creates personal wallets for users and a single wallet for drivers during registration.
  - Creates family wallets upon user request.
- **Wallet Management**:
  - Adds members to family wallets by phone number.
  - Processes cash-in operations (only owners for family wallets).
- **Transaction Management**:
  - Manages order transactions with statuses: create, blocked, success, canceled.
  - Verifies wallet balance during order creation, setting transactions to **blocked** or **canceled**.
  - Completes transactions upon order completion, deducting from user’s wallet and crediting driver’s wallet.
- **Transaction History**: Provides history, restricted to owners for family wallets.
- **Authentication**: Internal calls include user or driver information after token validation by the calling service.

## Nonfunctional Requirements

### General

- **GitHub Flow**: Maintain **main** (stable releases) and **dev** (development) branches. Create feature branches from **dev**, named with the Jira task ID.
- **Pull Requests**: Submit PRs to **dev** with the Jira task ID in the name. Include proof of work (e.g., video) for frontend PRs. Squash commits before merging.
- **CI/CD**: Configure for each service with steps: tests, linter, protofile linter, vulnerability check, and image build/upload to Docker Hub (master branch).
- **Deployment**: Deploy in Docker (docker-compose) and Kubernetes (Helm).
- **Configuration**: Use environment variables for settings (e.g., database connections, wait times).
- **Documentation**: Each service includes a README (startup instructions, environment variables) and Swagger (endpoint details).
- **Testing**: Implement unit and integration tests, with Postman collections and tests.
- **Authentication**: JWT managed by the Auth Service. Tokens are stored in Redis, and services validate tokens via the Auth Service, ensuring role-based access.

### Service-Specific

#### User Service
- **Database**: PostgreSQL for user and trip data.
- **Caching**: No token storage; managed by Auth Service.
- **Metrics**: Prometheus and Grafana.
- **Frontend**: Vue.js 3.0 with Composition API, using components and Pinia.
- **VCS**: GitHub; **CI/CD**: GitHub Actions.
- **Go Tools**: Gin for HTTP, golangci-lint, sqlc/sqlx/squirrel for PostgreSQL (no ORMs). Goose or go-migrate for migrations.
- **Driver Queue**: Use goroutines, channels, and sync packages.
- **Testing**: Table tests, gomock, testify/suite, ginkgo, gomega, dockertest.

#### Driver Service
- **Database**: MongoDB for driver info, trips, ratings, and balance.
- **Handlers**: Generate from Swagger.
- **Profiling**: Pprof, PGO, and an additional profiler.
- **Frontend**: Angular (latest) for forms and profile management.
- **VCS**: GitLab; **CI/CD**: GitLab CI/CD.
- **Rating**: Based on last 20 trips.

#### Order Service
- **Transport**: GraphQL for field-based searches and pagination.
- **Search**: Elasticsearch for prefix, full-text, transliteration, and lexical error searches.
- **Frontend**: React with Redux/Redux Toolkit for filterable main page.
- **VCS**: BitBucket; **CI/CD**: Bitbucket Pipelines.

#### Analytic Service
- **Database**: ClickHouse, consuming Kafka messages.
- **VCS**: GitHub; **CI/CD**: Circle CI.
- **HTTP Library**: Fiber.

#### Wallet Service
- **Database**: PostgreSQL for wallets and transactions.
- **HTTP Library**: Gin (assumed).
- **VCS**: GitHub (assumed); **CI/CD**: GitHub Actions (assumed).

#### Auth Service
- **Database**: Redis for token storage and blacklisting.
- **HTTP Library**: Gin (assumed).
- **VCS**: GitHub (assumed); **CI/CD**: GitHub Actions (assumed).

## Technical Requirements

- **Message Broker**: Kafka for event-driven communication.
- **RPC**: gRPC for inter-service communication.
- **Containerization**: Dockerfile and Makefile for each service to test, build, and deploy.
- **Go Guidelines**: Follow [Rakyll's Style Guide](https://rakyll.org/style-packages/).
- **Diagram Update**: Include updated schema diagram in .png and .drawio formats in the repository.

**Note**: Database schemas and API request/response structures must be designed independently and documented in each service’s Swagger files.

## Authentication Flow

- **Registration**:
  - Users and drivers register via User or Driver Service, creating wallets via the Wallet Service. Users may be assigned "User" and optionally "Analyst" roles.
- **Login**:
  - Clients send login requests to the Auth Service with credentials.
  - Auth Service verifies credentials with User or Driver Service, retrieves roles, and issues a JWT token containing user ID and roles (e.g., ["User", "Analyst"]).
- **Subsequent Requests**:
  - Clients include JWT token in headers.
  - Services validate tokens via Auth Service, which checks signature, expiration, and blacklist, returning roles.
  - Services enforce role-based access, rejecting requests if required roles (e.g., "User", "Driver", "Analyst") are missing.
- **Logout**:
  - Clients send logout requests to Auth Service, which blacklists the token in Redis.

## Inter-Service Interactions

- **User and Driver Services**:
  - Manage profiles and wallets via Wallet Service, with validated user/driver information.
  - Initiate orders via Order Service after token validation.
- **Order Service**:
  - Coordinates with Wallet Service for transaction lifecycles.
  - Assigns drivers and updates statuses, ensuring role-based access.
- **Wallet Service**:
  - Processes requests from other services, relying on validated tokens.
- **Analytic Service**:
  - Accessible only to users with "Analyst" role, validated via Auth Service.
- **Auth Service**:
  - Centralizes authentication and authorization, ensuring secure role-based access.

## Development Guidelines

- **Repository Structure**: Each service has a private repository on the specified VCS platform.
- **Branching Strategy**: GitHub Flow with **main** and **dev** branches. Feature branches include Jira task ID.
- **Pull Request Process**:
  - Create PRs to **dev** with Jira task ID.
  - Include proof of work for frontend changes.
  - Address mentor feedback and squash commits before merging.
- **CI/CD Pipelines**:
  - Tests, linters, protofile checks, vulnerability scans.
  - Build and push Docker images to Docker Hub (master branch).
- **Testing**:
  - Unit and integration tests.
  - Postman collections and tests for API validation.
- **Documentation**:
  - Detailed README per service.
  - Swagger documentation for API endpoints.
- **Deployment**:
  - Docker (docker-compose) and Kubernetes (Helm).
  - Configure settings via environment variables.
