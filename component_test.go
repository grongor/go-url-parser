package urlparser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const refUrl = "https://user:password@neque.erat.example.co.uk:1234" +
	"/lorem/ipsum/dolor/sit.html?amet=elit&metus=lectus#at=nostra&unde=omnis"

func TestGetComponent(t *testing.T) {
	tests := []struct {
		name      string
		rawUrl    string
		component string
		expected  string
		expectErr bool
	}{
		// ----------- scheme
		{
			name:      "scheme ref url",
			component: "scheme",
			expected:  "https",
		},
		{
			name:      "scheme only path",
			rawUrl:    "/lorem/ipsum.html",
			component: "scheme",
		},
		// ----------- authority
		{
			name:      "authority ref url",
			component: "authority",
			expected:  "user:password@neque.erat.example.co.uk:1234",
		},
		{
			name:      "authority user without password",
			rawUrl:    "http://user@example.co.uk",
			component: "authority",
			expected:  "user@example.co.uk",
		},
		{
			name:      "authority without auth",
			rawUrl:    "http://example.co.uk",
			component: "authority",
			expected:  "example.co.uk",
		},
		{
			name:      "authority only path",
			rawUrl:    "/lorem/ipsum.html",
			component: "authority",
		},
		// ----------- auth
		{
			name:      "auth ref url",
			component: "auth",
			expected:  "user:password",
		},
		{
			name:      "auth only user",
			rawUrl:    "http://user@example.co.uk",
			component: "auth",
			expected:  "user",
		},
		{
			name:      "auth no auth",
			rawUrl:    "http://example.co.uk",
			component: "auth",
		},
		// ----------- user
		{
			name:      "user ref url",
			component: "user",
			expected:  "user",
		},
		{
			name:      "user user only",
			rawUrl:    "http://user@example.co.uk",
			component: "user",
			expected:  "user",
		},
		{
			name:      "user no auth",
			rawUrl:    "http://example.co.uk",
			component: "user",
		},
		// ----------- password
		{
			name:      "password ref url",
			component: "password",
			expected:  "password",
		},
		{
			name:      "password only user",
			rawUrl:    "http://user@example.co.uk",
			component: "password",
		},
		{
			name:      "password no auth",
			rawUrl:    "http://example.co.uk",
			component: "password",
		},
		// ----------- hostPort
		{
			name:      "hostPort ref url",
			component: "hostPort",
			expected:  "neque.erat.example.co.uk:1234",
		},
		{
			name:      "hostPort no port",
			rawUrl:    "http://example.co.uk",
			component: "hostPort",
			expected:  "example.co.uk",
		},
		{
			name:      "hostPort only path",
			rawUrl:    "/lorem/ipsum.html",
			component: "hostPort",
		},
		// ----------- host
		{
			name:      "host ref url",
			component: "host",
			expected:  "neque.erat.example.co.uk",
		},
		{
			name:      "host only path",
			rawUrl:    "/lorem/ipsum.html",
			component: "host",
		},
		// ----------- tld
		{
			name:      "tld ref url",
			component: "tld",
			expected:  "co.uk",
		},
		{
			name:      "tld made up domain",
			rawUrl:    "http://what.is.this.garbage",
			component: "tld",
			expected:  "garbage",
		},
		{
			name:      "tld ipv4",
			rawUrl:    "http://192.168.1.1",
			component: "tld",
		},
		{
			name:      "tld ipv6",
			rawUrl:    "http://[2001:db8::1428:57ab]",
			component: "tld",
		},
		{
			name:      "tld only path",
			rawUrl:    "/lorem/ipsum.html",
			component: "tld",
		},
		// ----------- port
		{
			name:      "port ref url",
			component: "port",
			expected:  "1234",
		},
		{
			name:      "port only path",
			rawUrl:    "/lorem/ipsum.html",
			component: "port",
		},
		// ----------- path
		{
			name:      "path ref url",
			component: "path",
			expected:  "/lorem/ipsum/dolor/sit.html",
		},
		{
			name:      "path file",
			component: "path",
			rawUrl:    "http://example.co.uk/sit.html",
			expected:  "/sit.html",
		},
		{
			name:      "path slash",
			component: "path",
			rawUrl:    "http://example.co.uk/",
			expected:  "/",
		},
		{
			name:      "path only host",
			rawUrl:    "http://example.co.uk",
			component: "path",
		},
		// ----------- query
		{
			name:      "query ref url",
			component: "query",
			expected:  "amet=elit&metus=lectus",
		},
		{
			name:      "query only path",
			rawUrl:    "/lorem/ipsum.html",
			component: "query",
		},
		// ----------- fragment
		{
			name:      "fragment ref url",
			component: "fragment",
			expected:  "at=nostra&unde=omnis",
		},
		{
			name:      "fragment only path",
			rawUrl:    "/lorem/ipsum.html",
			component: "fragment",
		},
		// ----------- basePath
		{
			name:      "basePath ref url",
			component: "basePath",
			expected:  "/lorem/ipsum/dolor/",
		},
		{
			name:      "basePath only file",
			rawUrl:    "ipsum.html",
			component: "basePath",
			expected:  "/",
		},
		{
			name:      "basePath only host",
			rawUrl:    "http://example.co.uk",
			component: "basePath",
			expected:  "/",
		},
		// ----------- file
		{
			name:      "file ref url",
			component: "file",
			expected:  "sit.html",
		},
		{
			name:      "file no file",
			rawUrl:    "/lorem/ipsum",
			component: "file",
		},
		{
			name:      "file only host",
			rawUrl:    "http://example.co.uk",
			component: "file",
		},
		// ----------- ext
		{
			name:      "ext ref url",
			component: "ext",
			expected:  "html",
		},
		{
			name:      "ext no ext",
			rawUrl:    "/lorem/ipsum",
			component: "ext",
		},
		{
			name:      "ext only host",
			rawUrl:    "http://example.co.uk",
			component: "ext",
		},
		// ----------- relativeUrl
		{
			name:      "relativeUrl ref url",
			component: "relativeUrl",
			expected:  "/lorem/ipsum/dolor/sit.html?amet=elit&metus=lectus#at=nostra&unde=omnis",
		},
		{
			name:      "relativeUrl only path",
			rawUrl:    "/lorem/sit.html",
			component: "relativeUrl",
			expected:  "/lorem/sit.html",
		},
		{
			name:      "relativeUrl only query",
			rawUrl:    "?amet=elit",
			component: "relativeUrl",
			expected:  "?amet=elit",
		},
		{
			name:      "relativeUrl only fragment",
			rawUrl:    "#at=nostra",
			component: "relativeUrl",
			expected:  "#at=nostra",
		},
		{
			name:      "relativeUrl only path and query",
			rawUrl:    "/lorem/sit.html?amet=elit",
			component: "relativeUrl",
			expected:  "/lorem/sit.html?amet=elit",
		},
		{
			name:      "relativeUrl only path and fragment",
			rawUrl:    "/lorem/sit.html#at=nostra",
			component: "relativeUrl",
			expected:  "/lorem/sit.html#at=nostra",
		},
		{
			name:      "relativeUrl only query and fragment",
			rawUrl:    "?amet=elit#at=nostra",
			component: "relativeUrl",
			expected:  "?amet=elit#at=nostra",
		},
		{
			name:      "relativeUrl only host",
			rawUrl:    "http://example.co.uk",
			component: "relativeUrl",
		},
		// ----------- host:x
		{
			name:      "host:x ref url whole host",
			component: "host:0",
			expected:  "neque.erat.example.co.uk",
		},
		{
			name:      "host:x ref url skip first",
			component: "host:1",
			expected:  "erat.example.co.uk",
		},
		{
			name:      "host:x ref url skip last",
			component: "host:-1",
			expected:  "neque.erat.example",
		},
		{
			name:      "host:x ref url out of bounds",
			component: "host:10",
		},
		{
			name:      "host:x only path",
			rawUrl:    "/lorem/ipsum.html",
			component: "host:2",
		},
		// ----------- host:x:y
		{
			name:      "host:x:y ref url get second",
			component: "host:1:1",
			expected:  "erat",
		},
		{
			name:      "host:x:y ref url get last 2",
			component: "host:-0:2",
			expected:  "example.co.uk",
		},
		{
			name:      "host:x:y ref url out of bounds",
			component: "host:10:5",
		},
		{
			name:      "host:x:y only path",
			rawUrl:    "/lorem/ipsum.html",
			component: "host:2:5",
		},
		// ----------- hostPort:x
		{
			name:      "hostPort:x ref url whole host",
			component: "hostPort:0",
			expected:  "neque.erat.example.co.uk:1234",
		},
		{
			name:      "hostPort:x ref url skip first",
			component: "hostPort:1",
			expected:  "erat.example.co.uk:1234",
		},
		{
			name:      "hostPort:x ref url skip last",
			component: "hostPort:-1",
			expected:  "neque.erat.example:1234",
		},
		{
			name:      "hostPort:x ref url out of bounds",
			component: "hostPort:10",
			expected:  ":1234",
		},
		{
			name:      "hostPort:x only path",
			rawUrl:    "/lorem/ipsum.html",
			component: "hostPort:2",
		},
		// ----------- hostPort:x:y
		{
			name:      "hostPort:x:y ref url get second",
			component: "hostPort:1:1",
			expected:  "erat:1234",
		},
		{
			name:      "hostPort:x:y ref url get last 2",
			component: "hostPort:-0:2",
			expected:  "example.co.uk:1234",
		},
		{
			name:      "hostPort:x:y ref url out of bounds",
			component: "hostPort:10:5",
			expected:  ":1234",
		},
		{
			name:      "hostPort:x:y only path",
			rawUrl:    "/lorem/ipsum.html",
			component: "host:2:5",
		},
		// ----------- path:x
		{
			name:      "path:x ref url whole path",
			component: "path:0",
			expected:  "/lorem/ipsum/dolor/sit.html",
		},
		{
			name:      "path:x ref url skip first",
			component: "path:1",
			expected:  "/ipsum/dolor/sit.html",
		},
		{
			name:      "path:x ref url skip last",
			component: "path:-1",
			expected:  "/lorem/ipsum/dolor",
		},
		{
			name:      "path:x ref url out of bounds",
			component: "path:10",
		},
		{
			name:      "path:x only host",
			rawUrl:    "http://example.com",
			component: "path:2",
		},
		// ----------- path:x:y
		{
			name:      "path:x:y ref url get second",
			component: "path:1:1",
			expected:  "/ipsum",
		},
		{
			name:      "path:x:y ref url get last 2",
			component: "path:-0:2",
			expected:  "/dolor/sit.html",
		},
		{
			name:      "path:x:y ref url out of bounds",
			component: "path:10:5",
		},
		{
			name:      "path:x:y only host",
			rawUrl:    "http://example.com",
			component: "host:2:5",
		},
		// ----------- query:NAME
		{
			name:      "query:NAME ref url name amet",
			component: "query:amet",
			expected:  "elit",
		},
		{
			name:      "query:NAME ref url name not present",
			component: "query:not-here",
		},
		{
			name:      "query:NAME only path",
			rawUrl:    "/lorem/ipsum.html",
			component: "query:amet",
		},
		// ----------- fragment:NAME
		{
			name:      "fragment:NAME ref url name unde",
			component: "fragment:unde",
			expected:  "omnis",
		},
		{
			name:      "fragment:NAME ref url name not present",
			component: "fragment:not-here",
		},
		{
			name:      "fragment:NAME non-formatted fragment",
			rawUrl:    "just-some-paragraph-fragment",
			component: "fragment:unde",
		},
		{
			name:      "fragment:NAME only path",
			rawUrl:    "/lorem/ipsum.html",
			component: "fragment:unde",
		},
		// ----------- invalid components
		{
			name:      "invalid components not existing",
			component: "something",
			expectErr: true,
		},
		{
			name:      "invalid components invalid bounds",
			component: "path:--2",
			expectErr: true,
		},
		{
			name:      "invalid components invalid bounds 2",
			component: "path:1:5:3",
			expectErr: true,
		},
		{
			name:      "invalid components invalid bounds 3",
			component: "path:1:-2",
			expectErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)

			if test.rawUrl == "" {
				test.rawUrl = refUrl
			}

			result, err := Component(test.rawUrl, test.component)

			if test.expectErr {
				require.Error(err)
			} else {
				require.NoError(err)
			}

			require.Equal(test.expected, result)
		})
	}
}

func TestFactories(t *testing.T) {
	tests := []struct {
		name         string
		factoryValue string
		expected     string
	}{
		{
			name:         "factory PartialHost",
			factoryValue: PartialHost(-5, 3),
			expected:     "host:-5:3",
		},
		{
			name:         "factory PartialHostFrom",
			factoryValue: PartialHostFrom(2),
			expected:     "host:2",
		},
		{
			name:         "factory PartialHostPort",
			factoryValue: PartialHostPort(-5, 3),
			expected:     "hostPort:-5:3",
		},
		{
			name:         "factory PartialHostPortFrom",
			factoryValue: PartialHostPortFrom(2),
			expected:     "hostPort:2",
		},
		{
			name:         "factory PartialPath",
			factoryValue: PartialPath(-5, 3),
			expected:     "path:-5:3",
		},
		{
			name:         "factory PartialPathFrom",
			factoryValue: PartialPathFrom(2),
			expected:     "path:2",
		},
		{
			name:         "factory SingleQuery",
			factoryValue: SingleQuery("lorem"),
			expected:     "query:lorem",
		},
		{
			name:         "factory SingleFragment",
			factoryValue: SingleFragment("lorem"),
			expected:     "fragment:lorem",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expected, test.factoryValue)
		})
	}
}
