package pathfinder_test

import (
	"testing"

	"github.com/artyomturkin/pathfinder"
)

func TestSimpleLookup(t *testing.T) {
	n := &pathfinder.Node{}
	err := n.Add("hello", "world")
	if err != nil {
		t.Fatalf("failed to add route: %v", err)
	}
	i, _, err := n.Lookup("hello")
	if err != nil {
		t.Errorf("Failed to lookup: %v", err)
	}
	if s, ok := i.(string); !ok || s != "world" {
		t.Errorf("Found unexpected value: %v", i)
	}
}

func TestMultipleSimpleLookup(t *testing.T) {
	n := &pathfinder.Node{}
	// Add simple route
	err := n.Add("hello", "world")
	if err != nil {
		t.Fatalf("failed to add route 'hello': %v", err)
	}
	// Add longer route
	err = n.Add("hello/moon", "moon")
	if err != nil {
		t.Fatalf("failed to add route 'hello/moon': %v", err)
	}

	// Lookup simple route
	i, _, err := n.Lookup("hello")
	if err != nil {
		t.Errorf("Failed to lookup 'hello': %v", err)
	}
	if s, ok := i.(string); !ok || s != "world" {
		t.Errorf("Found unexpected value for 'hello': %v", i)
	}
	// Lookup longer route
	i, _, err = n.Lookup("hello/moon")
	if err != nil {
		t.Errorf("Failed to lookup 'hello/moon': %v", err)
	}
	if s, ok := i.(string); !ok || s != "moon" {
		t.Errorf("Found unexpected value for 'hello/moon': %v", i)
	}
	// Lookup nonexistent route
	i, _, err = n.Lookup("hello/sun")
	if err == nil {
		t.Errorf("Found something on nonexistent route 'hello/sun': %v", i)
	}
}

func TestParamLookup(t *testing.T) {
	n := &pathfinder.Node{}
	err := n.Add(":goodnight/:moon/tonight", "world")
	if err != nil {
		t.Fatalf("failed to add route: %v", err)
	}

	i, params, err := n.Lookup("goodnight/moon/tonight")
	if err != nil {
		t.Errorf("Failed to lookup: %v", err)
	}
	if s, ok := i.(string); !ok || s != "world" {
		t.Errorf("Found unexpected value: %v", i)
	}
	if len(params) != 2 {
		t.Errorf("Wrong amount of params found: %v", params)
	}
	if p, ok := params[":goodnight"]; !ok || p != "goodnight" {
		t.Errorf("Param ':goodnight' is not set to expected. Got: %s", p)
	}
	if p, ok := params[":moon"]; !ok || p != "moon" {
		t.Errorf("Param ':moon' is not set to expected. Got: %s", p)
	}
}

func TestMixedLookup(t *testing.T) {
	n := &pathfinder.Node{}
	err := n.Add("hello", "world")
	if err != nil {
		t.Fatalf("failed to add route: %v", err)
	}
	err = n.Add(":goodnight/:moon", "param")
	if err != nil {
		t.Fatalf("failed to add route: %v", err)
	}

	i, _, err := n.Lookup("hello")
	if err != nil {
		t.Errorf("Failed to lookup 'hello': %v", err)
	}
	if s, ok := i.(string); !ok || s != "world" {
		t.Errorf("Found unexpected value in 'hello': %v", i)
	}

	i, params, err := n.Lookup("goodnight/moon")
	if err != nil {
		t.Errorf("Failed to lookup 'goodnight/moon': %v", err)
	}
	if s, ok := i.(string); !ok || s != "param" {
		t.Errorf("Found unexpected value in 'goodnight/moon': %v", i)
	}
	if len(params) != 2 {
		t.Errorf("Wrong amount of params found: %v", params)
	}
	if p, ok := params[":goodnight"]; !ok || p != "goodnight" {
		t.Errorf("Param ':goodnight' is not set to expected. Got: %s", p)
	}
	if p, ok := params[":moon"]; !ok || p != "moon" {
		t.Errorf("Param ':moon' is not set to expected. Got: %s", p)
	}
}

func BenchmarkLookup(b *testing.B) {
	n := &pathfinder.Node{}
	err := n.Add("1/2/3/4/5/6/7/8/9/0", "world")
	if err != nil {
		b.Fatalf("failed to add route 'hello': %v", err)
	}
	err = n.Add("1/2/:3/:4/:5/:6/7/8/9", "world")
	if err != nil {
		b.Fatalf("failed to add route 'hello': %v", err)
	}

	b.Run("simple", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, err := n.Lookup("1/2/3/4/5/6/7/8/9/0")
			if err != nil {
				b.Errorf("failed to lookup: %v", err)
			}
		}
	})
	b.Run("param", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, err := n.Lookup("1/2/p/p/p/p/7/8/9")
			if err != nil {
				b.Errorf("failed to lookup: %v", err)
			}
		}
	})
}
