package cmd

import (
	"errors"
	"strings"
	"testing"
)

func TestCommandsRunCommandTable(t *testing.T) {
	tests := []struct {
		name       string
		cmdName    string
		register   bool
		handlerErr error
		wantErr    bool
		wantInErr  string
		wantCalled bool
	}{
		{
			name:       "registered command executes handler",
			cmdName:    "login",
			register:   true,
			wantErr:    false,
			wantCalled: true,
		},
		{
			name:       "handler error is returned",
			cmdName:    "register",
			register:   true,
			handlerErr: errors.New("boom"),
			wantErr:    true,
			wantInErr:  "boom",
			wantCalled: true,
		},
		{
			name:      "unknown command returns not found",
			cmdName:   "missing",
			register:  false,
			wantErr:   true,
			wantInErr: "command 'missing' not found",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			commands := &Commands{}
			called := false

			if tc.register {
				commands.Register(tc.cmdName, func(_ *state, _ CliCommand) error {
					called = true
					return tc.handlerErr
				})
			}

			err := commands.RunCommand(nil, CliCommand{Name: tc.cmdName})

			if tc.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !tc.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if tc.wantInErr != "" && (err == nil || !strings.Contains(err.Error(), tc.wantInErr)) {
				t.Fatalf("expected error to contain %q, got %v", tc.wantInErr, err)
			}

			if called != tc.wantCalled {
				t.Fatalf("expected called=%v, got %v", tc.wantCalled, called)
			}
		})
	}
}
