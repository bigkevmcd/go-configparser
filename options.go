package configparser

type options struct {
	allowNoValue bool
}

// OptFunc modifies options.
type OptFunc func(*options)

// AllowNoValue allows option with no value to be saved
// as empty line.
func AllowNoValue(o *options) { o.allowNoValue = true }
