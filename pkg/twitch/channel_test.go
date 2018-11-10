package twitch

import "testing"

func TestNewChannel(t *testing.T) {
	name := "channel"
	c := newChannel(name)

	if c == nil {
		t.Fatal("newChannel unexpectedly returned nil")
	}
	if c.Name != name {
		t.Errorf("expected channel Name to be %s, got %s", name, c.Name)
	}
	if c.enabledModules == nil {
		t.Error("expected enabledModules to not be nil")
	}
	if c.modulesMutex == nil {
		t.Error("expected modulesMutex to not be nil")
	}
	if c.modules == nil {
		t.Error("expected modules to not be nil")
	}
}
