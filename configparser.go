package configparser

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

var (
	sectionHeader = regexp.MustCompile(`^\[([^]]+)\]$`)
	interpolater  = regexp.MustCompile(`%\(([^)]*)\)s`)
)

var ErrAlreadyExist = errors.New("already exist")

var boolMapping = map[string]bool{
	"1":     true,
	"true":  true,
	"on":    true,
	"yes":   true,
	"0":     false,
	"false": false,
	"off":   false,
	"no":    false,
}

// Dict is a simple string->string map.
type Dict map[string]string

// Config represents a Python style configuration file.
type Config map[string]*Section

// ConfigParser ties together a Config and default values for use in
// interpolated configuration values.
type ConfigParser struct {
	config   Config
	defaults *Section
	opt      *options
}

// Keys returns a sorted slice of keys
func (d Dict) Keys() []string {
	var keys []string

	for key := range d {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys
}

func getNoSectionError(section string) error {
	return fmt.Errorf("no section: %q", section)
}

func getNoOptionError(section, option string) error {
	return fmt.Errorf("no option %q in section: %q", option, section)
}

// New creates a new ConfigParser.
func New() *ConfigParser {
	return &ConfigParser{
		config:   make(Config),
		defaults: newSection(defaultSectionName),
		opt:      defaultOptions(),
	}
}

// NewWithOptions creates a new ConfigParser with options.
func NewWithOptions(opts ...optFunc) *ConfigParser {
	opt := defaultOptions()
	for _, fn := range opts {
		fn(opt)
	}

	return &ConfigParser{
		config:   make(Config),
		defaults: newSection(opt.defaultSection),
		opt:      opt,
	}
}

// NewWithDefaults allows creation of a new ConfigParser with a pre-existing
// Dict.
func NewWithDefaults(defaults Dict) (*ConfigParser, error) {
	p := New()
	for key, value := range defaults {
		if err := p.defaults.Add(key, value); err != nil {
			return nil, fmt.Errorf("failed to add %q to %q: %w", key, value, err)
		}
	}
	return p, nil
}

// NewConfigParserFromFile creates a new ConfigParser struct populated from the
// supplied filename.
func NewConfigParserFromFile(filename string) (*ConfigParser, error) {
	p, err := Parse(filename)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// ParseReader parses a ConfigParser from the provided input.
func ParseReader(in io.Reader) (*ConfigParser, error) {
	p := New()
	err := p.ParseReader(in)

	return p, err
}

// ParseReaderWithOptions parses a ConfigParser from the provided input with given options.
func ParseReaderWithOptions(in io.Reader, opts ...optFunc) (*ConfigParser, error) {
	p := NewWithOptions(opts...)
	err := p.ParseReader(in)

	return p, err
}

// Parse takes a filename and parses it into a ConfigParser value.
func Parse(filename string) (*ConfigParser, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	p, err := ParseReader(file)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// ParseWithOptions takes a filename and parses it into a ConfigParser value with given options.
func ParseWithOptions(filename string, opts ...optFunc) (*ConfigParser, error) {
	p := NewWithOptions(opts...)
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = p.ParseReader(bytes.NewReader(data))
	return p, err
}

func writeSection(file *os.File, delimiter string, section *Section) error {
	_, err := file.WriteString(fmt.Sprintf("[%s]\n", section.Name))
	if err != nil {
		return err
	}

	for _, option := range section.Options() {
		_, err = file.WriteString(fmt.Sprintf("%s %s %s\n", option, delimiter, section.options[option]))
		if err != nil {
			return err
		}
	}
	_, err = file.WriteString("\n")
	return err
}

// SaveWithDelimiter writes the current state of the ConfigParser to the named
// file with the specified delimiter.
func (p *ConfigParser) SaveWithDelimiter(filename, delimiter string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if len(p.defaults.Options()) > 0 {
		err = writeSection(f, delimiter, p.defaults)
		if err != nil {
			return err
		}
	}

	for _, s := range p.Sections() {
		err = writeSection(f, delimiter, p.config[s])
		if err != nil {
			return err
		}
	}

	return nil
}

// ParseReader parses data into ConfigParser from provided reader.
func (p *ConfigParser) ParseReader(in io.Reader) error {
	reader := bufio.NewReader(in)
	var lineNo int
	var curSect *Section

	keyValue := regexp.MustCompile(
		fmt.Sprintf(
			`([^%[1]s\s][^%[1]s]*)\s*(?P<vi>[%[1]s]+)\s*(.*)$`,
			p.opt.delimeters,
		),
	)
	keyWNoValue := regexp.MustCompile(
		fmt.Sprintf(
			`([^%[1]s\s][^%[1]s]*)\s*((?P<vi>[%[1]s]+)\s*(.*)$)?`,
			p.opt.delimeters,
		),
	)

	for {
		l, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		lineNo++
		if len(l) == 0 {
			continue
		}
		line := strings.TrimFunc(string(l), unicode.IsSpace) // ensures sectionHeader regex will match

		// Skip comment lines and empty lines
		if p.opt.commentPrefixes.AtFirst(line) || line == "" {
			continue
		}

		if match := sectionHeader.FindStringSubmatch(line); len(match) > 0 {
			section := match[1]
			if section == p.opt.defaultSection {
				curSect = p.defaults
			} else if _, present := p.config[section]; !present {
				curSect = newSection(section)
				p.config[section] = curSect
			} else if p.opt.strict {
				return fmt.Errorf("section %q error: %w", section, ErrAlreadyExist)
			}
		} else if match = keyValue.FindStringSubmatch(line); len(match) > 0 {
			if curSect == nil {
				return fmt.Errorf("missing section header: %d %s", lineNo, line)
			}
			key := strings.TrimSpace(match[1])
			if p.opt.strict {
				if err := p.inOptions(key); err != nil {
					return err
				}
			}

			value := p.opt.inlineCommentPrefixes.Split(match[3])
			if err := curSect.Add(key, value); err != nil {
				return fmt.Errorf("failed to add %q = %q: %w", key, value, err)
			}
		} else if match = keyWNoValue.FindStringSubmatch(line); len(match) > 0 &&
			p.opt.allowNoValue && curSect != nil {
			key := strings.TrimSpace(match[1])
			if p.opt.strict {
				if err := p.inOptions(key); err != nil {
					return err
				}
			}

			value := p.opt.inlineCommentPrefixes.Split(match[4])
			if err := curSect.Add(key, value); err != nil {
				return fmt.Errorf("failed to add %q = %q: %w", key, value, err)
			}
		}
	}

	return nil
}
