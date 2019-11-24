package urlparser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name      string
		rawUrl    string
		format    string
		expected  string
		expectErr bool
	}{
		{
			name:     "all basic fields",
			format:   "{scheme}://{authority}{relativeUrl}",
			expected: refUrl,
		},
		{
			name:     "duplicate some fields",
			format:   "{scheme} {tld} {scheme} {path} {tld}",
			expected: "https co.uk https /lorem/ipsum/dolor/sit.html co.uk",
		},
		{
			name:     "with partial placeholders",
			format:   "{{scheme}{{tld}{err{scheme}}}{{tld}}",
			expected: "{https{co.uk{errhttps}}{co.uk}",
		},
		{
			name:      "with invalid placeholders",
			format:    "{scheme} {tld} {wrong}",
			expectErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)

			if test.rawUrl == "" {
				test.rawUrl = refUrl
			}

			result, err := Format(test.rawUrl, test.format)

			if test.expectErr {
				require.Error(err)
			} else {
				require.NoError(err)
			}

			require.Equal(test.expected, result)
		})
	}
}
