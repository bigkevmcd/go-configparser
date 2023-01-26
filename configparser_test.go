package configparser_test

import (
	"strings"

	"github.com/bigkevmcd/go-configparser"
	. "gopkg.in/check.v1"
	gc "gopkg.in/check.v1"

	"io/ioutil"
	"os"

	"path"
	"testing"
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

	data, err := ioutil.ReadAll(f)
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

	data, err := ioutil.ReadAll(f)
	c.Assert(err, IsNil)
	f.Close()
	c.Assert(string(data), Equals, "[DEFAULT]\ntesting = value\n\n[othersection]\nmyoption = myvalue\nnewoption = novalue\n\n[testing]\nmyoption = value\n\n")
}

// ParseFromReader() parses the Config data from an io.Reader.
func (s *ConfigParserSuite) TestParseFromReader(c *gc.C) {
	parsed, err := configparser.ParseReader(strings.NewReader("[DEFAULT]\ntesting = value\n\n[othersection]\nmyoption = myvalue\nnewoption = novalue\nfinal = foo[bar]\n\n[testing]\nmyoption = value\n\n"))
	c.Assert(err, gc.IsNil)

	result, err := parsed.Items("othersection")
	c.Assert(err, gc.IsNil)
	c.Assert(result, gc.DeepEquals, configparser.Dict{"myoption": "myvalue", "newoption": "novalue", "final": "foo[bar]"})
}

// If AllowNoValue is set to true, parser should recognize options without values.
func (s *ConfigParserSuite) TestParseFromReaderWNoValue(c *gc.C) {
	configparser.AllowNoValue = true
	defer func() { configparser.AllowNoValue = false }()

	parsed, err := configparser.ParseReader(strings.NewReader("[empty]\noption\n\n"))
	c.Assert(err, gc.IsNil)

	ok, err := parsed.HasOption("empty", "option")
	c.Assert(err, gc.IsNil)
	c.Assert(ok, Equals, true)
}

func assertSuccessful(c *gc.C, err error) {
	c.Assert(err, gc.IsNil)
}
