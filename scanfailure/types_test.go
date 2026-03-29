package scanfailure

import (
	"testing"
)

func TestReasonFriendlyText(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "known code returns friendly text",
			code:     ReasonSBOMGenerationFailed,
			expected: "Failed to generate software inventory (SBOM) for this image",
		},
		{
			name:     "scanner OOM code",
			code:     ReasonScannerOOMKilled,
			expected: "SBOM scanner was killed due to memory limits — consider increasing the scanner memory limit",
		},
		{
			name:     "empty string returns unexpected text",
			code:     "",
			expected: reasonFriendlyText[ReasonUnexpected],
		},
		{
			name:     "unknown code returns code itself",
			code:     "some_future_code",
			expected: "some_future_code",
		},
		{
			name:     "all known codes have non-empty text",
			code:     ReasonImageTooLarge,
			expected: "Image exceeds the maximum size limit for vulnerability scanning",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReasonFriendlyText(tt.code)
			if got != tt.expected {
				t.Errorf("ReasonFriendlyText(%q) = %q, want %q", tt.code, got, tt.expected)
			}
		})
	}
}

func TestAllReasonCodesHaveFriendlyText(t *testing.T) {
	codes := []string{
		ReasonSBOMGenerationFailed, ReasonImageTooLarge, ReasonSBOMTooLarge,
		ReasonSBOMIncomplete, ReasonImageAuthFailed, ReasonImageNotFound,
		ReasonCVEMatchingFailed, ReasonResultUploadFailed, ReasonSBOMStorageFailed,
		ReasonScannerOOMKilled, ReasonScanTimeout, ReasonUnexpected,
	}
	for _, code := range codes {
		text := ReasonFriendlyText(code)
		if text == "" || text == code {
			t.Errorf("ReasonFriendlyText(%q) returned %q — expected human-friendly text", code, text)
		}
	}
}
