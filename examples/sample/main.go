package main

import (
	"fmt"
)

// Example function for demonstration
// This function shows how to use the package
func ExamplePrintMessage(message string) {
	fmt.Println(message)
}

// ServiceA provides core functionality
type ServiceA struct {
	name string
}

// NewServiceA creates a new ServiceA instance
func NewServiceA(name string) *ServiceA {
	return &ServiceA{name: name}
}

// Process does some processing
func (s *ServiceA) Process(data string) string {
	return fmt.Sprintf("ServiceA processed: %s", data)
}

// ServiceB depends on ServiceA
type ServiceB struct {
	serviceA *ServiceA
}

// NewServiceB creates ServiceB with dependency
func NewServiceB(serviceA *ServiceA) *ServiceB {
	return &ServiceB{serviceA: serviceA}
}

// Execute executes the service
func (s *ServiceB) Execute(input string) string {
	result := s.serviceA.Process(input)
	return fmt.Sprintf("ServiceB result: %s", result)
}
