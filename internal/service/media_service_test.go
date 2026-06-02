package service

import "testing"

// Guards against silent shrinkage of the upload caps. If this test fails,
// someone lowered a limit — make sure that was intentional and update the
// frontend (VideoUploader.vue, ImageUploader.vue) in lockstep.
func TestAllowedContentTypes_Caps(t *testing.T) {
	cases := []struct {
		ct      string
		wantMin int64
	}{
		{"image/jpeg", 10 * 1024 * 1024},
		{"image/png", 10 * 1024 * 1024},
		{"image/webp", 10 * 1024 * 1024},
		{"video/mp4", 500 * 1024 * 1024},
	}
	for _, c := range cases {
		got, ok := allowedContentTypes[c.ct]
		if !ok {
			t.Errorf("%s missing from allowedContentTypes", c.ct)
			continue
		}
		if got < c.wantMin {
			t.Errorf("%s cap shrunk: got %d, want >= %d", c.ct, got, c.wantMin)
		}
	}
}
