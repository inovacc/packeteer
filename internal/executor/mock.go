package executor

import (
	"context"
	"fmt"
)

// MockExecutor is a test double for CommandExecutor.
// Set the response fields before calling Execute.
type MockExecutor struct {
	// Responses maps "binary args..." to a mock response.
	Responses map[string]MockResponse
	// Calls records all executed commands for assertions.
	Calls []MockCall
	// DefaultResponse is returned when no matching response is found.
	DefaultResponse *MockResponse
}

// MockResponse holds canned output for a mocked command.
type MockResponse struct {
	Stdout []byte
	Stderr []byte
	Err    error
}

// MockCall records a single command execution.
type MockCall struct {
	Binary string
	Args   []string
}

// NewMockExecutor creates a mock executor with an empty response map.
func NewMockExecutor() *MockExecutor {
	return &MockExecutor{
		Responses: make(map[string]MockResponse),
	}
}

func (m *MockExecutor) Execute(_ context.Context, binary string, args []string) ([]byte, []byte, error) {
	m.Calls = append(m.Calls, MockCall{Binary: binary, Args: args})

	key := binary
	if len(args) > 0 {
		key = binary + " " + args[0]
	}

	if resp, ok := m.Responses[key]; ok {
		return resp.Stdout, resp.Stderr, resp.Err
	}

	if resp, ok := m.Responses[binary]; ok {
		return resp.Stdout, resp.Stderr, resp.Err
	}

	if m.DefaultResponse != nil {
		return m.DefaultResponse.Stdout, m.DefaultResponse.Stderr, m.DefaultResponse.Err
	}

	return nil, nil, fmt.Errorf("no mock response for %s", key)
}

func (m *MockExecutor) BinaryPath(name string) (string, error) {
	return "/usr/bin/" + name, nil
}
