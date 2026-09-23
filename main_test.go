package main

import "testing"

func TestResolveToken(t *testing.T) {
	t.Setenv("GITLAB_TOKEN", "environment-token")

	if got := resolveToken(""); got != "environment-token" {
		t.Errorf("resolveToken with no flag = %q, want environment token", got)
	}
	if got := resolveToken("flag-token"); got != "flag-token" {
		t.Errorf("resolveToken with flag = %q, want flag token", got)
	}
}
