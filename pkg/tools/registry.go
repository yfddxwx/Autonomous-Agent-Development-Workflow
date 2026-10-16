package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"agent/pkg/models"
)

// Tool represents an executable function callable by LLM.
type Tool interface {
	Name() string
	Description() string
	Parameters() map[string]interface{}
	Execute(ctx context.Context, args map[string]interface{}) (string, error)
}

// Registry manages available tools.
type Registry struct {
	tools map[string]Tool
}

// NewRegistry creates a new tool registry.
func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]Tool),
	}
}

// Register registers a tool.
func (r *Registry) Register(tool Tool) {
	r.tools[tool.Name()] = tool
}

// Get finds a tool by name.
func (r *Registry) Get(name string) (Tool, bool) {
	t, ok := r.tools[name]
	return t, ok
}

// GetDefinitions returns definitions for all registered tools formatted for LLM.
func (r *Registry) GetDefinitions() []models.ToolDefinition {
	defs := make([]models.ToolDefinition, 0, len(r.tools))
	for _, t := range r.tools {
		defs = append(defs, models.ToolDefinition{
			Type: "function",
			Function: models.ToolFunctionDefinition{
				Name:        t.Name(),
				Description: t.Description(),
				Parameters:  t.Parameters(),
			},
		})
	}
	return defs
}

// Execute parses JSON arguments and executes the tool.
func (r *Registry) Execute(ctx context.Context, name string, argsJSON string) (string, error) {
	t, ok := r.Get(name)
	if !ok {
		return "", fmt.Errorf("tool '%s' not found", name)
	}

	var args map[string]interface{}
	if argsJSON != "" {
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("failed to parse arguments JSON: %w", err)
		}
	} else {
		args = make(map[string]interface{})
	}

	return t.Execute(ctx, args)
}
