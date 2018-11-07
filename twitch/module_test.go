package twitch

import "testing"

func TestNewModule(t *testing.T) {
	name := "test"
	m := newModule(name)

	if m == nil {
		t.Fatal("newModule unexpectedly returned nil")
	}
	if m.Name != name {
		t.Errorf("expected module name to be %s, got %s", name, m.Name)
	}
	if m.commands == nil {
		t.Error("module's commands map was not initialised")
	}
	if m.enabledCommands == nil {
		t.Error("module's enabledCommands map was not initialised")
	}
}

func TestAddCommand(t *testing.T) {
	commandName := "test"
	c := &Command{Name: commandName}
	m := newModule("module")

	if err := m.AddCommand(c); err != nil {
		t.Fatalf("AddCommand unexpectedly returned an error: %v", err)
	}
	c, ok := m.commands[commandName]
	if !ok {
		t.Fatalf("Command with name %s does not exist in the module's commands map", commandName)
	}

}

func TestAddCommand_NilCommand(t *testing.T) {
	m := newModule("module")
	if err := m.AddCommand(nil); err == nil {
		t.Error("AddCommand did not throw error when adding nil Command")
	}
}

func TestAddCommand_ExistingCommand(t *testing.T) {
	c := &Command{Name: "test"}
	m := newModule("module")

	if err := m.AddCommand(c); err != nil {
		t.Errorf("AddCommand returned unexpected error: %v", err)
	}
	if err := m.AddCommand(c); err == nil {
		t.Error("expected AddCommand to return error")
	}
}

func TestCommands(t *testing.T) {
	m := newModule("module")

	tests := []struct {
		name     string
		commands map[string]*Command
	}{
		{
			name:     "no commands",
			commands: make(map[string]*Command),
		},
		{
			name: "single command",
			commands: map[string]*Command{
				"command": &Command{Name: "command"},
			},
		},
		{
			name: "multiple commands",
			commands: map[string]*Command{
				"command1": &Command{Name: "command1"},
				"command2": &Command{Name: "command2"},
			},
		},
	}

	for _, test := range tests {
		m.commands = test.commands
		commands := m.Commands()
		if len(commands) != len(test.commands) {
			t.Errorf("expected %d commands, got %d", len(test.commands), len(commands))
		}
		for name, command := range test.commands {
			found := false
			for _, c := range commands {
				if c.Name == name {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("comamnd with name %s was not found in the returned commands", name)
			}
			if name != command.Name {
				t.Errorf("name in map differs from actual command name %s", command.Name)
			}
		}
	}
}
