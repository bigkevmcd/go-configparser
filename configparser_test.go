package configparser_test

import (
	"io"
	"os"
	"path"
	"strings"
	"testing"

	. "gopkg.in/check.v1"

	"github.com/bigkevmcd/go-configparser"
)

func Test(t *testing.T) { TestingT(t) }

type ConfigParserSuite struct {
	p *configparser.ConfigParser
}

var _ = Suite(&ConfigParserSuite{})

func (s *ConfigParserSuite) SetUpTest(c *C) {
	s.p, _ = configparser.NewConfigParserFromFile("testdata/example.cfg")
}

// NewWithDefaults should add defaults to the configuration
func (s *ConfigParserSuite) TestNewWithDefaults(c *C) {
	n := make(configparser.Dict)
	n["testing"] = "value"

	p, err := configparser.NewWithDefaults(n)
	c.Assert(err, IsNil)

	d := p.Defaults()
	c.Assert(d["testing"], Equals, "value")
}

// NewWithDefaults should copy the items from the provided map
func (s *ConfigParserSuite) TestNewWithDefaultsCopied(c *C) {
	n := make(configparser.Dict)
	n["testing"] = "value"
	p, err := configparser.NewWithDefaults(n)
	c.Assert(err, IsNil)

	n["testing2"] = "myvalue"

	d := p.Defaults()
	c.Assert(d["testing2"], Equals, "")
}

// SaveWithDelimiter(filename) should write the current state of the
// configuration to the provided filename with the specified delimiter
func (s *ConfigParserSuite) TestSaveWithDelimiter(c *C) {
	p := configparser.New()

	err := p.AddSection("testing")
	c.Assert(err, IsNil)
	err = p.Set("testing", "myoption", "value")
	c.Assert(err, IsNil)
	err = p.AddSection("othersection")
	c.Assert(err, IsNil)
	err = p.Set("othersection", "newoption", "novalue")
	c.Assert(err, IsNil)
	err = p.Set("othersection", "myoption", "myvalue")
	c.Assert(err, IsNil)

	tempfile := path.Join(c.MkDir(), "config.cfg")
	err = p.SaveWithDelimiter(tempfile, "=")
	c.Assert(err, IsNil)

	f, err := os.Open(tempfile)
	c.Assert(err, IsNil)

	data, err := io.ReadAll(f)
	c.Assert(err, IsNil)
	f.Close()
	c.Assert(string(data), Equals, "[othersection]\nmyoption = myvalue\nnewoption = novalue\n\n[testing]\nmyoption = value\n\n")
}

// Save(filename) should correctly write out the defaults section with the
// defaults
func (s *ConfigParserSuite) TestSaveWithDelimiterAndDefaults(c *C) {
	n := make(configparser.Dict)
	n["testing"] = "value"
	p, err := configparser.NewWithDefaults(n)
	c.Assert(err, IsNil)

	err = p.AddSection("testing")
	c.Assert(err, IsNil)
	err = p.Set("testing", "myoption", "value")
	c.Assert(err, IsNil)
	err = p.AddSection("othersection")
	c.Assert(err, IsNil)
	err = p.Set("othersection", "newoption", "novalue")
	c.Assert(err, IsNil)
	err = p.Set("othersection", "myoption", "myvalue")
	c.Assert(err, IsNil)

	tempfile := path.Join(c.MkDir(), "config.cfg")
	err = p.SaveWithDelimiter(tempfile, "=")
	c.Assert(err, IsNil)

	f, err := os.Open(tempfile)
	c.Assert(err, IsNil)

	data, err := io.ReadAll(f)
	c.Assert(err, IsNil)
	f.Close()
	c.Assert(string(data), Equals, "[DEFAULT]\ntesting = value\n\n[othersection]\nmyoption = myvalue\nnewoption = novalue\n\n[testing]\nmyoption = value\n\n")
}

// ParseFromReader() parses the Config data from an io.Reader.
func (s *ConfigParserSuite) TestParseFromReader(c *C) {
	parsed, err := configparser.ParseReader(strings.NewReader("[DEFAULT]\ntesting = value\n\n[othersection]\nmyoption = myvalue\nnewoption = novalue\nfinal = foo[bar]\n\n[testing]\nmyoption = value\nemptyoption\n\n"))
	c.Assert(err, IsNil)

	result, err := parsed.Items("othersection")
	c.Assert(err, IsNil)
	c.Assert(result, DeepEquals, configparser.Dict{
		"myoption":  "myvalue",
		"newoption": "novalue",
		"final":     "foo[bar]",
	})
}

// TestMultilineValue tests multiline value parsing.
func (s *ConfigParserSuite) TestMultilineValue(c *C) {
	parsed, err := configparser.ParseReader(
		strings.NewReader(`[DEFAULT]
testing = multiline
 value

myoption = another
 multiline
		value

broken_option = this value will miss

 its multiline
`),
	)
	c.Assert(err, IsNil)
	result, err := parsed.Items("DEFAULT")
	c.Assert(err, IsNil)
	c.Assert(result, DeepEquals, configparser.Dict{
		"testing":       "multiline\nvalue",
		"myoption":      "another\nmultiline\nvalue",
		"broken_option": "this value will miss",
	})
}

func assertSuccessful(c *C, err error) {
	c.Assert(err, IsNil)
}
