package module

// Flags holds the command-line configuration for the HTTP scan module.
// Populated by the framework.
//
// TODO: Custom headers?
type HTTP struct {
	Base            `group:"base"`
	TLSFlags        `group:"TLS"`
	Method          string `long:"method" default:"GET" description:"Set HTTP request method type"`
	Endpoint        string `long:"endpoint" default:"/" description:"Send an HTTP request to an endpoint"`
	FailHTTPToHTTPS bool   `long:"fail-http-to-https" description:"Trigger retry-https logic on known HTTP/400 protocol mismatch responses"`
	UserAgent       string `long:"user-agent" default:"Mozilla/5.0 zgrab/0.x" description:"Set a custom user agent"`
	RetryHTTPS      bool   `long:"retry-https" description:"If the initial request fails, reconnect and try with HTTPS."`
	MaxSize         int    `long:"max-size" default:"256" description:"Max kilobytes to read in response to an HTTP request"`
	MaxRedirects    int    `long:"max-redirects" default:"0" description:"Max number of redirects to follow"`

	// FollowLocalhostRedirects overrides the default behavior to return
	// ErrRedirLocalhost whenever a redirect points to localhost.
	FollowLocalhostRedirects bool `long:"follow-localhost-redirects" description:"Follow HTTP redirects to localhost"`

	// UseHTTPS causes the first request to be over TLS, without requiring a
	// redirect to HTTPS. It does not change the port used for the connection.
	UseHTTPS bool `long:"use-https" description:"Perform an HTTPS connection on the initial host"`

	// RedirectsSucceed causes the ErrTooManRedirects error to be suppressed
	RedirectsSucceed bool `long:"redirects-succeed" description:"Redirects are always a success, even if max-redirects is exceeded"`

	// Set arbitrary HTTP headers
	CustomHeadersNames     string `long:"custom-headers-names" description:"CSV of custom HTTP headers to send to server"`
	CustomHeadersValues    string `long:"custom-headers-values" description:"CSV of custom HTTP header values to send to server. Should match order of custom-headers-names."`
	CustomHeadersDelimiter string `long:"custom-headers-delimiter" description:"Delimiter for customer header name/value CSVs"`

	OverrideSH bool `long:"override-sig-hash" description:"Override the default SignatureAndHashes TLS option with more expansive default"`

	// ComputeDecodedBodyHashAlgorithm enables computing the body hash later than the default,
	// using the specified algorithm, allowing a user of the response to recompute a matching hash
	ComputeDecodedBodyHashAlgorithm string `long:"compute-decoded-body-hash-algorithm" choice:"sha256" choice:"sha1" description:"Choose algorithm for BodyHash field"`

	// WithBodyLength enables adding the body_size field to the Response
	WithBodyLength bool `long:"with-body-size" description:"Enable the body_size attribute, for how many bytes actually read"`
}

// Description returns an overview of this module.
func (flags *HTTP) Description() string {
	return "Send an HTTP request and read the response, optionally following redirects."
}

// Execute performs any needed validation on the arguments.
func (flags *HTTP) Execute(args []string) error {
	return nil
}

// Help returns module-specific help.
func (flags *HTTP) Help() string {
	return ""
}
