# Product Requirements Document (PRD)

## 1. Project Overview
The Task Management API is a robust, highly performant RESTful backend service designed to handle user authentication, role management, and task tracking. Built with Go (Golang) and adhering to Clean Architecture principles, the system ensures high maintainability, secure data transactions, and rapid response times capable of handling production-level traffic.

---

## 2. Objectives & Goals

- Performance : Achieve sub-10ms response times for high-frequency read operations using in-memory caching.
- Security : Implement robust authentication (JWT) and authorization (RBAC) to protect sensitive user data and administrative endpoints.
- Reliability : Guarantee data integrity during complex operations using strict database transactions (ACID properties).
- Developer Experience : Provide comprehensive, interactive API documentation and fully automated CI pipelines for seamless integration and deployment.

---

## 3. System Architecture
The application is structured using a Modular Monolithic Clean Architecture, dividing responsibilities into distinct layers to decouple business logic from framework-specific routing and database drivers:

- Router/Middleware Layer: Intercepts HTTP requests, handles CORS, and verifies JWT/Roles.
- Controller Layer: Parses incoming JSON payloads and formats HTTP responses.
- Service Layer: Executes core business logic and manages database transactions.
- Repository Layer: Abstracts direct database interactions and raw SQL/ORM queries.
- Model Layer: Defines database schemas and entity relationships.

---

## 4. Key Features & Requirements
#### 4.1. Authentication & Authorization (RBAC)

- User Registration & Login: Secure password hashing (Bcrypt) and JWT generation.
- Role-Based Access Control: Users are assigned distinct roles (user or admin).
- user: Can read/update their own profile and manage their personal tasks.
- admin: Possesses elevated privileges to retrieve all users and execute account deletions.
- Identity Extraction: Endpoints strictly extract user_id and role from the verified JWT payload rather than URL parameters to prevent IDOR (Insecure Direct Object Reference) vulnerabilities.

#### 4.2. User Management

- Get My Profile: Retrieves the authenticated user's personal details.
- Update Profile: Allows users to modify their personal information.
- Admin User Management: Restricted endpoints for monitoring and managing the entire user base.

#### 4.3. Task Management
- Task CRUD Operations: Users can create, read, update, and delete their tasks.
- Get My Tasks: Fetches a filtered list of tasks exclusively belonging to the authenticated user.
- Automated Onboarding: A default welcome task is automatically generated within the same database transaction when a new user registers.

#### 4.4. Performance & Reliability Optimization
- Redis Cache-Aside Pattern: Frequently accessed data (e.g., retrieving all users) is temporarily stored in Redis, drastically reducing PostgreSQL load and improving response times.
- Database Transactions: Registration logic utilizes gorm.Transaction to ensure that if the automated task creation fails, the user record is seamlessly rolled back.
- Database Indexing: B-Tree indexing is applied to the email column to accelerate login queries.

---

## 5. Technology Stack
| Component | Technology | Description |
| -------- | -------- | -------- |
| Language |    Go    | Version 1.25+ |
| Web Framework | Gin | GonicHigh-performance HTTP router |
| Database | PostgreSQL | Primary relational database |
| ORM      |   GORM   | Object-Relational Mapping library |
| Caching  |  Redis   | In-memory data store for high-speed retrieval |
| API Documentation | Swagger (swaggo) | Auto-generated interactive API docs | 
| Containerization | Docker & Compose | Unified environment for App, DB, and Cache | 
| CI/CD    |  GitHub Actions | Automated unit testing and build verification |

---

## 6. Infrastructure & Deployment
- Dockerized Environment: The entire application stack (Go app, PostgreSQL, Redis) is containerized and orchestrated via docker-compose.yml for isolated and reproducible local development.
- Continuous Integration (CI): A GitHub Actions workflow (ci.yml) is triggered on every push to the main branch to automatically resolve dependencies, execute unit tests, and verify build integrity to prevent regressions.

---
