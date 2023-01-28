package configparser

import (
	"strings"

	"github.com/bigkevmcd/go-configparser/chainmap"
)

const defaultSectionName = "DEFAULT"

// options allows to control parser behavior.
type options struct {
	interpolation         Interpolator
	commentPrefixes       Prefixes
	inlineCommentPrefixes Prefixes
	defaultSection        string
	delimeters            string
	converters            Converter
	allowNoValue          bool
	strict                bool

	// TODO: under investigation, have no effect now.
	emptyLines bool
}

// Converter contains custom convert functions for available types.
//
// List of types:
//   - int -> int64
//   - float -> float64
//   - string
//   - bool
//
// The caller should guarantee type assertion to the requested type
// after custom processing!
type Converter map[string]ConvertFunc

// ConvertFunc is a custom datatype converter.
type ConvertFunc func(string) (any, error)

// Prefixes stores available prefixes for comments.
type Prefixes []string

// AtFirst checks if str starts with any of the prefixes.
func (pr Prefixes) AtFirst(str string) bool {
	for _, p := range pr {
		if strings.HasPrefix(str, p) {
			return true
		}
	}

	return false
}

// Split splits str with the first prefix found.
// Returns original string if no matches.
func (pr Prefixes) Split(str string) string {
	for _, p := range pr {
		if strings.Contains(str, p) {
			return strings.Split(str, p)[0]
		}
	}

	return str
}

// Interpolator defines interpolation instance.
// For more details, check [chainmap.ChainMap] realisation.
type Interpolator interface {
	Add(...chainmap.Dict)
	Len() int
	Get(string) string
}

// defaultOptions presets required options.
func defaultOptions() *options {
	return &options{
		interpolation:   chainmap.New(),
		defaultSection:  defaultSectionName,
		delimeters:      ":=",
		commentPrefixes: Prefixes{"#"},
	}
}

type optFunc func(*options)

// Interpolation sets custom interpolator.
func Interpolation(i Interpolator) optFunc {
	return func(o *options) {
		o.interpolation = i
	}
}

// CommentPrefixes sets a slice of comment prefixes.
// Lines of configuration file which starts with
// the first match in this slice will be skipped.
func CommentPrefixes(pr Prefixes) optFunc {
	return func(o *options) {
		o.commentPrefixes = pr
	}
}

// InlineCommentPrefixes sets a slice of inline comment delimeters.
// When parsing a value, it will be split with
// the first match in this slice.
func InlineCommentPrefixes(pr Prefixes) optFunc {
	return func(o *options) {
		o.inlineCommentPrefixes = pr
	}
}

// DefaultSection sets the name of the default section.
func DefaultSection(n string) optFunc {
	return func(o *options) {
		o.defaultSection = n
	}
}

// Delimeters sets a string of delimeters for option-value pairs.
func Delimeters(d string) optFunc {
	return func(o *options) {
		o.delimeters = d
	}
}

// Converters sets custom convertion functions. Will apply to return
// value of the Get* methods instead of the default convertion.
//
// NOTE: the caller should guarantee type assetion to the requested type
// after custom processing.
func Converters(conv Converter) optFunc {
	return func(o *options) {
		o.converters = conv
	}
}

// AllowNoValue allows option with no value to be saved
// as empty line.
func AllowNoValue(o *options) { o.allowNoValue = true }

// Strict prohibits the duplicates of options and values.
func Strict(o *options) { o.strict = true }
