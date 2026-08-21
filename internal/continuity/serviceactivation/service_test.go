package serviceactivation_test

import (
	"carecontinuity/internal/continuity/serviceactivation"
	"errors"
	"testing"
)

func TestCommunityServiceActivationPublicBehavior(t *testing.T) {
	c := serviceactivation.NewCoordinator()
	failure := errors.New("activation rollback")
	if err := c.Activate("east", "provider-7", func() error { return failure }); !errors.Is(err, failure) {
		t.Fatalf("expected failure, got %v", err)
	}
	if owner, ok := c.Owner("east"); ok {
		t.Fatalf("failed activation retained owner %q", owner)
	}
	if err := c.Activate("east", "provider-7", func() error { return nil }); err != nil {
		t.Fatal(err)
	}
	if owner, ok := c.Owner("east"); !ok || owner != "provider-7" {
		t.Fatalf("successful route %q %v", owner, ok)
	}
}
