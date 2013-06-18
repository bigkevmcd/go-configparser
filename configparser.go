package configparser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"
)

const (
	defaultSectionName        = "DEFAULT"
	maxInterpolationDepth int = 10
)

var (
	sectionHeader    = regexp.MustCompile("\\[([^]]+)\\]")
	keyValue         = regexp.MustCompile("([^:=\\s][^:=]*)\\s*(?P<vi>[:=])\\s*(.*)$")
	continuationLine = regexp.MustCompile("\\w+(.*)$")
	interpolater     = regexp.MustCompile("%\\(([^)]*)\\)s")
)

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

type Dict map[string]string
type Config map[string]Dict

type ConfigParser struct {
	config   Config
	defaults Dict
}

func getNoSectionError(section string) error {
	return fmt.Errorf("No section: '%s'", section)
}

func getNoOptionError(section, option string) error {
	return fmt.Errorf("No option '%s' in section: '%s'", option, section)
}

func New() *ConfigParser {
	return &ConfigParser{
		config:   make(Config),
		defaults: make(Dict),
	}
}

func NewWithDefaults(defaults Dict) *ConfigParser {
	p := ConfigParser{
		config:   make(Config),
		defaults: make(Dict),
	}
	for key, value := range defaults {
		p.defaults[key] = value
	}
	return &p
}

// Create a new ConfigParser struct populated from the supplied filename
func NewConfigParserFromFile(filename string) (*ConfigParser, error) {
	p, err := Parse(filename)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func parseFile(file *os.File) (*ConfigParser, error) {
	p := New()
	defer file.Close()

	reader := bufio.NewReader(file)
	var lineNo int
	var err error
	var curSect Dict

	for err == nil {
		l, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		lineNo++
		if len(l) == 0 {
			continue
		}
		line := strings.TrimFunc(string(l), unicode.IsSpace)

		// Skip comment lines and empty lines
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		if match := sectionHeader.FindStringSubmatch(line); len(match) > 0 {
			section := match[1]
			if section == defaultSectionName {
				curSect = p.defaults
			} else if _, present := p.config[section]; !present {
				curSect = make(Dict)
				p.config[section] = curSect
			}
		} else if match = keyValue.FindStringSubmatch(line); len(match) > 0 {
			if curSect == nil {
				return nil, fmt.Errorf("Missing Section Header: %d %s", lineNo, line)
			} else {
				option := match[1]
				// separator := match[2]
				value := p.transformOption(strings.TrimFunc(match[3], unicode.IsSpace))
				curSect[option] = value
			}
		}
	}
	return p, nil
}

func Parse(filename string) (*ConfigParser, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	p, err := parseFile(file)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func writeSection(file *os.File, name, delimiter string, options Dict) error {
	_, err := file.WriteString(fmt.Sprintf("[%s]\n", name))
	if err != nil {
		return err
	}
	for k, v := range options {
		_, err = file.WriteString(fmt.Sprintf("%s %s %s\n", k, delimiter, v))
		if err != nil {
			return err
		}
	}
	_, err = file.WriteString("\n")
	return err
}

// Save the current state of the ConfigParser to the named file with the specified delimiter
func (p *ConfigParser) SaveWithDelimiter(filename, delimiter string) error {
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		return err
	}

	d := p.Defaults()
	if len(d) > 0 {
		err = writeSection(f, "defaults", delimiter, d)
		if err != nil {
			return err
		}
	}

	for _, s := range p.Sections() {
		d, err = p.Items(s)
		if err != nil {
			return err
		}
		err = writeSection(f, s, delimiter, d)
		if err != nil {
			return err
		}
	}

	return nil
}
