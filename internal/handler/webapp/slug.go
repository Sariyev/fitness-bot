package webapp

import (
	"fmt"
	"strings"
	"time"
)

// autoSlug returns a URL-safe slug derived from `name`, suffixed with a
// timestamp so two items with identical names (or Russian-only names that
// reduce to empty) don't collide on the slug UNIQUE constraint.
//
// Used as a fallback for admin forms where the slug field is optional —
// the form placeholder is "auto-from-name" but until now the server was
// inserting literal "" on submit, which exploded the second time around
// (first empty slug claimed the column; subsequent ones violated UNIQUE).
//
// Examples:
//   autoSlug("Test Program")     -> "test-program-1748520000000"
//   autoSlug("Тестовая программа") -> "1748520000000"   (no ASCII letters → just the suffix)
//   autoSlug("Lower-back rehab") -> "lower-back-rehab-1748520000000"
func autoSlug(name string) string {
	var sb strings.Builder
	prevDash := true
	for _, r := range strings.ToLower(strings.TrimSpace(name)) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			sb.WriteRune(r)
			prevDash = false
		case r == ' ' || r == '_' || r == '-':
			if !prevDash {
				sb.WriteByte('-')
				prevDash = true
			}
		}
	}
	s := strings.Trim(sb.String(), "-")
	suffix := fmt.Sprintf("%d", time.Now().UnixMilli())
	if s == "" {
		return suffix
	}
	return s + "-" + suffix
}

// ensureSlug returns `slug` if non-empty, otherwise an auto-generated slug
// from `name`. Convenience wrapper for the admin create/update handlers.
func ensureSlug(slug, name string) string {
	if strings.TrimSpace(slug) != "" {
		return slug
	}
	return autoSlug(name)
}
