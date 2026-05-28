package webapp

import (
	"regexp"
	"strings"
	"testing"
)

func TestAutoSlug_LatinName_HasNameAndSuffix(t *testing.T) {
	got := autoSlug("Test Program")
	if !strings.HasPrefix(got, "test-program-") {
		t.Errorf("expected prefix 'test-program-', got %q", got)
	}
	if !regexp.MustCompile(`^test-program-\d+$`).MatchString(got) {
		t.Errorf("expected '<slug>-<digits>' format, got %q", got)
	}
}

func TestAutoSlug_CyrillicName_FallsBackToTimestamp(t *testing.T) {
	got := autoSlug("Тестовая программа")
	// No ASCII letters → only the unix-millis suffix.
	if !regexp.MustCompile(`^\d+$`).MatchString(got) {
		t.Errorf("expected pure-digit suffix for Cyrillic name, got %q", got)
	}
}

func TestAutoSlug_StripsLeadingTrailingDashes(t *testing.T) {
	got := autoSlug("  ---Hello---  ")
	if !strings.HasPrefix(got, "hello-") {
		t.Errorf("expected 'hello-<digits>', got %q", got)
	}
}

func TestAutoSlug_CollapsesMultipleSpaces(t *testing.T) {
	got := autoSlug("a   b   c")
	if !strings.HasPrefix(got, "a-b-c-") {
		t.Errorf("expected 'a-b-c-<digits>', got %q", got)
	}
}

func TestAutoSlug_NoCollisionUnderRapidCalls(t *testing.T) {
	// Two calls in succession should produce different slugs because of the
	// UnixMilli suffix — this is what prevents the UNIQUE constraint bug we
	// hit before ensureSlug was added.
	seen := map[string]bool{}
	for i := 0; i < 50; i++ {
		s := autoSlug("Same Name")
		if seen[s] {
			t.Fatalf("duplicate slug produced: %q (iteration %d)", s, i)
		}
		seen[s] = true
	}
}

func TestEnsureSlug_KeepsExplicitSlug(t *testing.T) {
	if got := ensureSlug("custom-slug", "ignored"); got != "custom-slug" {
		t.Errorf("ensureSlug should keep explicit slug, got %q", got)
	}
}

func TestEnsureSlug_GeneratesWhenBlank(t *testing.T) {
	for _, in := range []string{"", "   ", "\t\n"} {
		got := ensureSlug(in, "Test Name")
		if !strings.HasPrefix(got, "test-name-") {
			t.Errorf("ensureSlug(%q, ...) expected auto-generated, got %q", in, got)
		}
	}
}
