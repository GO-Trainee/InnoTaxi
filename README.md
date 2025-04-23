# InnoTaxi - Taxi Ordering Application

## Overview

InnoTaxi is a taxi ordering application designed to facilitate taxi bookings with three distinct user roles: **User**, **Driver**, and **Analyst**. It supports two wallet types—**Personal** and **Family**—and offers three taxi types: **Economy**, **Comfort**, and **Business**. Built on a microservice architecture with clean architecture principles, the application ensures modularity, scalability, and maintainability.

**Important**: All repositories created for this project must remain **private** and are subject to a Non-Disclosure Agreement (NDA). They must **never** be made public.

## Microservices

The application consists of six microservices, each handling specific functionalities:

1. **User Service**: Manages user profiles, taxi orders, and wallet interactions for users.
2. **Driver Service**: Handles driver profiles, trip statuses, and wallet operations for drivers.
3. **Order Service**: Orchestrates taxi order creation, driver assignment, and transaction processing.
4. **Analytic Service**: Provides statistical insights and ratings for analysts.
5. **Wallet Service**: Manages wallets and transactions for users and drivers.
6. **Auth Service**: Centralizes authentication and authorization for all roles.

## Functional Requirements

### Auth Service

- **Authentication**: Processes login and logout requests for users, drivers, and analysts.
- **Token Management**: Issues JSON Web Tokens (JWT) upon successful login and invalidates them upon logout by adding to a Redis blacklist.
- **Credential Verification**: Communicates with User, Driver, and Analytic Services to verify login credentials.
- **Token Validation**: Validates JWT tokens for incoming requests, providing role information (User, Driver, Analyst) to enforce access control. Ensures users cannot access driver or analyst functionalities, and vice versa, resulting in a "Permission Denied" error for unauthorized actions.

### User Service

- **Registration**: Users sign up directly with name, phone number, email, and password. A personal wallet is created via the Wallet Service.
- **Profile Management**: Users can view, update (name, phone number, email), or soft-delete their profiles.
- **Wallet Operations**:
  - View available wallets (personal and family) through the Wallet Service.
  - Create family wallets and add members by phone number.
  - Cash in personal or family wallets (only the owner for family wallets).
  - View transaction history (restricted to owners for family wallets).
- **Order Taxi**:
  - Users specify taxi type, start and end locations, and select a wallet.
  - The system verifies wallet balance via the Wallet Service.
  - If sufficient, it requests a free driver through the Order Service. If no drivers are available, users join a queue with a configurable wait time. If no driver is found, a rejection response is sent.
- **Rate Trip**: Users can rate their last trip (1–5) with an optional comment, if within a configurable time since trip completion.
- **View Trips**: Users can view past trips, including taxi type, driver, and route details.
- **Restriction**: Orders are blocked if the selected wallet has insufficient funds.
- **Authentication**: All operations require a valid JWT token, validated through the Auth Service to confirm the User role.

### Driver Service

- **Registration**: Drivers sign up directly with name, phone number, email, password, and taxi type. A single driver wallet is created via the Wallet Service.
- **Status Management**: Drivers can toggle their status (free/busy) after a trip, updating the order status via the Order Service.
- **Rate Trip**: Drivers can rate their last trip (1–5), if within a configurable time since trip completion.
- **View Rating**: Drivers can view their rating, calculated from the last 20 trips.
- **View Trips**: Drivers can view past trips, including taxi type, user, and route details.
- **Wallet Operations**: Drivers can view their single wallet’s balance and transaction history via the Wallet Service.
- **Authentication**: All operations require a valid JWT token, validated through the Auth Service to confirm the Driver role.

### Order Service

- **Order Orchestration**: Manages taxi order creation, driver assignment, and status updates. Interacts with the Wallet Service to block funds during order creation and complete transactions upon trip completion.
- **Order List**: Provides a filtered and paginated list of orders for authorized roles.
- **Order Fields**: Includes user, driver, start and end locations, taxi type, date, status (in progress, finished), and comment.
- **Pricing**: Maintains pricing information for each taxi type.
- **Search**: Analysts can search orders with partial field matching.
- **Driver Selection**: Supports selecting drivers based on user ratings.
- **Authentication**: Requests are validated through the Auth Service to ensure role-based access (e.g., users for ordering, drivers for status updates, analysts for searches).

### Analytic Service

- **Statistics**: Analysts can view order statistics, such as counts by day or month.
- **Ratings**: Analysts can view ratings for all drivers and users.
- **Account**: Analysts use pre-created accounts, logging in via the Auth Service.
- **Data Recording**: Records all registrations and completed orders in its database.
- **Authentication**: Access requires a valid JWT token, validated through the Auth Service to confirm the Analyst role.

### Wallet Service

- **Wallet Creation**:
  - Creates personal wallets for users during registration.
  - Creates a single driver wallet for drivers during registration.
  - Creates family wallets upon user request.
- **Wallet Management**:
  - Adds members to family wallets by phone number.
  - Processes cash-in operations for personal or family wallets (only owners for family wallets).
- **Transaction Management**:
  - Manages order transactions with statuses: create, blocked, success, canceled.
  - Verifies sufficient wallet balance during order creation, setting transactions to **blocked** or **canceled**.
  - Completes transactions upon order completion, deducting from the user’s wallet and crediting the driver’s wallet.
- **Transaction History**: Provides transaction history, with access restricted to owners for family wallets.
- **Authentication**: Internal calls from other services include user or driver information after token validation by the calling service.

## Nonfunctional Requirements

### General

- **GitHub Flow**: Maintain two main branches: **main** (stable releases) and **dev** (development). Create feature branches from **dev**, named with the Jira task ID.
- **Pull Requests**: Submit PRs to **dev** with the Jira task ID in the name. Include proof of work (e.g., video) for frontend PRs. Squash commits before merging.
- **CI/CD**: Configure for each service with steps: tests, linter, protofile linter, vulnerability check, and image build/upload to Docker Hub (master branch).
- **Deployment**: Deploy in two environments: Docker via docker-compose and Kubernetes via Helm.
- **Configuration**: Use environment variables for configurable settings, such as database connections and wait times.
- **Documentation**: Each service must include a README (startup instructions, environment variables) and Swagger (endpoint details).
- **Testing**: Implement unit and integration tests, along with Postman collections and tests.
- **Authentication**: Use JWT managed by the Auth Service. Tokens are stored in Redis, and all services validate tokens by calling the Auth Service, ensuring role-based access control.

### Service-Specific

#### User Service

- **Database**: PostgreSQL for storing user and trip data.
- **Caching**: No token storage; tokens are managed by the Auth Service.
- **Metrics**: Use Prometheus and Grafana for monitoring.
- **Frontend**: Develop with Vue.js 3.0 using Composition API, incorporating components and Pinia for state management.
- **VCS**: GitHub; **CI/CD**: GitHub Actions.
- **Go Tools**: Use Gin for HTTP, golangci-lint for linting, and sqlc, sqlx, or squirrel for PostgreSQL queries (no ORMs). Use Goose or go-migrate for migrations.
- **Driver Queue**: Implement using goroutines, channels, and sync packages.
- **Testing**: Employ table tests, gomock, testify/suite, ginkgo, gomega, and dockertest for integration tests.

#### Driver Service

- **Database**: MongoDB for storing driver information, trips, ratings, and balance.
- **Handlers**: Generate handlers from Swagger files to avoid manual HTTP interactions.
- **Profiling**: Implement Pprof and Profile-Guided Optimization (PGO), plus an additional profiler of choice.
- **Frontend**: Use Angular (latest version) for registration, authorization, and profile management.
- **VCS**: GitLab; **CI/CD**: GitLab CI/CD.
- **Rating**: Calculate driver ratings based on the last 20 trips.

#### Order Service

- **Transport**: Use GraphQL for API interactions, supporting field-based searches and pagination.
- **Search**: Implement Elasticsearch for prefix, full-text, transliteration, and lexical error searches.
- **Frontend**: Develop with React and Redux/Redux Toolkit for a filterable main page.
- **VCS**: BitBucket; **CI/CD**: Bitbucket Pipelines.

#### Analytic Service

- **Database**: ClickHouse for storing and querying analytical data.
- **Message Consumption**: Consume messages from Kafka into ClickHouse.
- **VCS**: GitHub; **CI/CD**: Circle CI.
- **HTTP Library**: Use Fiber for HTTP handling.

#### Wallet Service

- **Database**: PostgreSQL for managing wallets and transactions.
- **HTTP Library**: Use Gin (assumed for consistency).
- **VCS**: GitHub (assumed); **CI/CD**: GitHub Actions (assumed).

#### Auth Service

- **Database**: Redis for storing and blacklisting JWT tokens.
- **HTTP Library**: Use Gin (assumed for consistency).
- **VCS**: GitHub (assumed); **CI/CD**: GitHub Actions (assumed).

## Technical Requirements

- **Message Broker**: Use Kafka for event-driven communication between services.
- **RPC**: Implement gRPC for efficient inter-service communication.
- **Containerization**: Provide a Dockerfile and Makefile for each service to handle testing, building, and deployment.
- **Go Guidelines**: Adhere to style guidelines outlined in Rakyll's Style Guide.
- **Diagram Update**: When adding new services, update the application schema diagram and include it in the repository in .png and .drawio formats.

**Note**: Database schemas and API request/response structures must be designed independently based on each service’s requirements. These are not detailed in this README and should be documented in each service’s Swagger files.

## Authentication Flow

- **Registration**:
  - Users and drivers register directly with the User or Driver Service, respectively, providing necessary details and creating wallets via the Wallet Service.
- **Login**:
  - Clients send login requests to the Auth Service, specifying credentials and role (User, Driver, Analyst).
  - The Auth Service verifies credentials by calling the corresponding service (User, Driver, or Analytic).
  - Upon successful verification, the Auth Service issues a JWT token containing the user’s role and stores it in Redis.
- **Subsequent Requests**:
  - Clients include the JWT token in request headers when accessing any service.
  - The receiving service calls the Auth Service to validate the token, checking its signature, expiration, and blacklist status in Redis, and retrieves the user’s role.
  - The service enforces role-based access control, rejecting requests with a “Permission Denied” error if the role does not match the required permissions (e.g., a User attempting Driver actions).
- **Logout**:
  - Clients send logout requests to the Auth Service, which blacklists the token in Redis, rendering it invalid for future requests.

## Inter-Service Interactions

- **User and Driver Services**:
  - Handle profile and wallet-related operations by calling the Wallet Service with validated user or driver information.
  - Initiate taxi orders by interacting with the Order Service after token validation.
- **Order Service**:
  - Coordinates with the Wallet Service to manage transaction lifecycles (blocking and completing funds).
  - Assigns drivers and updates order statuses, ensuring role-based access via Auth Service validation.
- **Wallet Service**:
  - Processes internal requests from other services, relying on the calling service to validate tokens and provide user or driver IDs.
- **Analytic Service**:
  - Provides data access to analysts, with all requests validated through the Auth Service to confirm the Analyst role.
- **Auth Service**:
  - Serves as the central authority for authentication and authorization, ensuring secure and role-appropriate access across the system.

## Development Guidelines

- **Repository Structure**: Each service must have its own private repository on the specified VCS platform (GitHub, GitLab, or BitBucket).
- **Branching Strategy**: Follow GitHub Flow with **main** and **dev** branches. Feature branches should include the Jira task ID.
- **Pull Request Process**:
  - Create PRs from feature branches to **dev**.
  - Include the Jira task ID in the PR title.
  - Provide proof of work for frontend changes (e.g., a video).
  - Address mentor feedback in the same branch and squash commits before merging.
- **CI/CD Pipelines**:
  - Run tests, linters, protofile checks, and vulnerability scans.
  - Build and push Docker images to Docker Hub for the master branch.
- **Testing**:
  - Write unit and integration tests to ensure robust functionality.
  - Create Postman collections and tests for API validation.
- **Documentation**:
  - Maintain a detailed README per service, covering setup and environment variables.
  - Generate Swagger documentation from code to describe API endpoints.
- **Deployment**:
  - Support deployment in Docker (via docker-compose) and Kubernetes (via Helm).
  - Configure all variable settings through environment variables.

## Key Citations

- Rakyll's Go Style Guidelines for Package Design
