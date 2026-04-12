# API Documentation

This documentation was automatically generated from source code.

## Table of Contents

- [Project Overview](#project-overview)
- [Packages](#packages)
- [Architecture Analysis](#architecture-analysis)

## Project Overview

**Total Packages:** 4

**Total Elements:** 27

## Packages

### Package: `api`

**Description:** Package api provides HTTP API layer

**Imports:**
- `awdoc/examples/complex/service`
- `fmt`

#### Exported Elements

**Functions:**

- **`NewUserHandler`** (function)
  ```go
  func NewUserHandler(svc *service.UserService) *UserHandler
  ```
  NewUserHandler creates a new handler

**Methods:**

- **`HandleGetUser`** (method)
  ```go
  func (*UserHandler) HandleGetUser(id int) (string, error)
  ```
  HandleGetUser handles GET /users/:id request

- **`HandleCreateUser`** (method)
  ```go
  func (*UserHandler) HandleCreateUser(name string) error
  ```
  HandleCreateUser handles POST /users request

- **`HandleListUsers`** (method)
  ```go
  func (*UserHandler) HandleListUsers() ([]string, error)
  ```
  HandleListUsers handles GET /users request

- **`HandleDeleteUser`** (method)
  ```go
  func (*UserHandler) HandleDeleteUser(id int) error
  ```
  HandleDeleteUser handles DELETE /users/:id request

**Structs:**

- **`UserHandler`** (struct)
  UserHandler handles HTTP requests for users


---

### Package: `database`

**Description:** Package database provides database access layer

#### Exported Elements

**Functions:**

- **`NewPostgreSQL`** (function)
  ```go
  func NewPostgreSQL() *PostgreSQL
  ```
  NewPostgreSQL creates a new PostgreSQL instance

**Methods:**

- **`Connect`** (method)
  ```go
  func (*PostgreSQL) Connect(url string) error
  ```
  Connect establishes connection to PostgreSQL

- **`Query`** (method)
  ```go
  func (*PostgreSQL) Query(sql string) ([]map[string]interface{}, error)
  ```
  Query executes a SELECT query

- **`Exec`** (method)
  ```go
  func (*PostgreSQL) Exec(sql string, args unknown) error
  ```
  Exec executes an INSERT/UPDATE/DELETE query

- **`Close`** (method)
  ```go
  func (*PostgreSQL) Close() error
  ```
  Close closes the database connection

**Structs:**

- **`PostgreSQL`** (struct)
  PostgreSQL implements Database interface for PostgreSQL

**Interfaces:**

- **`Database`** (interface)
  Database interface defines core database operations


---

### Package: `repository`

**Description:** Package repository provides data access layer

**Imports:**
- `awdoc/examples/complex/database`

#### Exported Elements

**Functions:**

- **`NewUserRepository`** (function)
  ```go
  func NewUserRepository(db database.Database) *UserRepository
  ```
  NewUserRepository creates a new repository

**Methods:**

- **`GetByID`** (method)
  ```go
  func (*UserRepository) GetByID(id int) (*User, error)
  ```
  GetByID fetches a user by ID

- **`Save`** (method)
  ```go
  func (*UserRepository) Save(user *User) error
  ```
  Save saves a user to database

- **`Delete`** (method)
  ```go
  func (*UserRepository) Delete(id int) error
  ```
  Delete deletes a user from database

- **`GetAll`** (method)
  ```go
  func (*UserRepository) GetAll() ([]*User, error)
  ```
  GetAll retrieves all users

**Structs:**

- **`UserRepository`** (struct)
  UserRepository handles user database operations

- **`User`** (struct)
  User represents a user entity


---

### Package: `service`

**Description:** Package service provides business logic layer

**Imports:**
- `awdoc/examples/complex/repository`

#### Exported Elements

**Functions:**

- **`NewUserService`** (function)
  ```go
  func NewUserService(repo *repository.UserRepository) *UserService
  ```
  NewUserService creates a new service

**Methods:**

- **`RegisterUser`** (method)
  ```go
  func (*UserService) RegisterUser(name string) error
  ```
  RegisterUser registers a new user

- **`GetUser`** (method)
  ```go
  func (*UserService) GetUser(id int) (*repository.User, error)
  ```
  GetUser fetches user information

- **`UpdateUser`** (method)
  ```go
  func (*UserService) UpdateUser(id int, name string) error
  ```
  UpdateUser updates user information

- **`DeleteUser`** (method)
  ```go
  func (*UserService) DeleteUser(id int) error
  ```
  DeleteUser removes a user

- **`ListUsers`** (method)
  ```go
  func (*UserService) ListUsers() ([]*repository.User, error)
  ```
  ListUsers returns all users

**Structs:**

- **`UserService`** (struct)
  UserService provides user business logic


---

## Architecture Analysis

### Architectural Layers

**Layer 0:**
- api
- database
- repository
- service

### Complex Packages (God Objects)

Packages with high complexity that might need refactoring:

- **service** (Complexity: 7, Dependencies: 0)

### Dependency Graph

```
```

