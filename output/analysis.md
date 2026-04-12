# API Documentation

This documentation was automatically generated from source code.

## Table of Contents

- [Project Overview](#project-overview)
- [Packages](#packages)
- [Architecture Analysis](#architecture-analysis)

## Project Overview

**Total Packages:** 1

**Total Elements:** 7

## Packages

### Package: `main`

**Imports:**
- `fmt`

#### Exported Elements

**Functions:**

- **`ExamplePrintMessage`** (function)
  ```go
  func ExamplePrintMessage(message string)
  ```
  Example function for demonstration
This function shows how to use the package

- **`NewServiceA`** (function)
  ```go
  func NewServiceA(name string) *ServiceA
  ```
  NewServiceA creates a new ServiceA instance

- **`NewServiceB`** (function)
  ```go
  func NewServiceB(serviceA *ServiceA) *ServiceB
  ```
  NewServiceB creates ServiceB with dependency

**Methods:**

- **`Process`** (method)
  ```go
  func (*ServiceA) Process(data string) string
  ```
  Process does some processing

- **`Execute`** (method)
  ```go
  func (*ServiceB) Execute(input string) string
  ```
  Execute executes the service

**Structs:**

- **`ServiceA`** (struct)
  ServiceA provides core functionality

- **`ServiceB`** (struct)
  ServiceB depends on ServiceA


---

## Architecture Analysis

### Architectural Layers

**Layer 0:**
- main

### Complex Packages (God Objects)

Packages with high complexity that might need refactoring:

- **main** (Complexity: 7, Dependencies: 0)

### Dependency Graph

```
```

