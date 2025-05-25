package logging

import "fmt"

type Field struct {
	Key   string
	Value interface{}
}

type ILogger interface {
	LogInfo(message string, fields ...Field)
	LogWarning(message string, fields ...Field)
	LogError(message string, fields ...Field)
}

type ConsoleLogger struct{}

func (l ConsoleLogger) LogInfo(message string, fields ...Field) {
	fmt.Printf("INFO: %s\n", message)
	for _, field := range fields {
		fmt.Printf("  %s: %v\n", field.Key, field.Value)
	}
}

func (l ConsoleLogger) LogWarning(message string, fields ...Field) {
	fmt.Printf("WARNING: %s\n", message)
	for _, field := range fields {
		fmt.Printf("  %s: %v\n", field.Key, field.Value)
	}
}

func (l ConsoleLogger) LogError(message string, fields ...Field) {
	fmt.Printf("ERROR: %s\n", message)
	for _, field := range fields {
		fmt.Printf("  %s: %v\n", field.Key, field.Value)
	}
}
