package configparser_test

import (
	"github.com/bigkevmcd/go-configparser"
	. "gopkg.in/check.v1"

	"fmt"
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
	s.p, _ = configparser.NewConfigParserFromFile("example.cfg")
}

// NewWithDefaults should add defaults to the configuration
func (s *ConfigParserSuite) TestNewWithDefaults(c *C) {
	n := make(configparser.Dict)
	n["testing"] = "value"

	p := configparser.NewWithDefaults(n)

	d := p.Defaults()
	c.Assert(d["testing"], Equals, "value")
}

// NewWithDefaults should copy the items from the provided map
func (s *ConfigParserSuite) TestNewWithDefaultsCopied(c *C) {
	n := make(configparser.Dict)
	n["testing"] = "value"
	p := configparser.NewWithDefaults(n)

	n["testing2"] = "myvalue"

	d := p.Defaults()
	c.Assert(d["testing2"], Equals, "")
}

// SaveWithDelimiter(filename) should write the current state of the
// configuration to the provided filename with the specified delimiter
func (s *ConfigParserSuite) TestSaveWithDelimiter(c *C) {
	p := configparser.New()

	p.AddSection("testing")
	p.Set("testing", "myoption", "value")
	p.AddSection("othersection")
	p.Set("othersection", "newoption", "novalue")
	p.Set("othersection", "myoption", "myvalue")

	tempfile := path.Join(c.MkDir(), "config.cfg")
	err := p.SaveWithDelimiter(tempfile, "=")
	c.Assert(err, IsNil)

	f, err := os.Open(tempfile)
	c.Assert(err, IsNil)

	data, err := ioutil.ReadAll(f)
	c.Assert(err, IsNil)
	f.Close()
	c.Assert(fmt.Sprintf("%s", data), Equals, "[othersection]\nmyoption = myvalue\nnewoption = novalue\n\n[testing]\nmyoption = value\n\n")
}

// Save(filename) should correctly write out the defaults section with the
// defaults
func (s *ConfigParserSuite) TestSaveWithDelimiterAndDefaults(c *C) {
	n := make(configparser.Dict)
	n["testing"] = "value"
	p := configparser.NewWithDefaults(n)

	p.AddSection("testing")
	p.Set("testing", "myoption", "value")
	p.AddSection("othersection")
	p.Set("othersection", "newoption", "novalue")
	p.Set("othersection", "myoption", "myvalue")

	tempfile := path.Join(c.MkDir(), "config.cfg")
	err := p.SaveWithDelimiter(tempfile, "=")
	c.Assert(err, IsNil)

	f, err := os.Open(tempfile)
	c.Assert(err, IsNil)

	data, err := ioutil.ReadAll(f)
	c.Assert(err, IsNil)
	f.Close()
	c.Assert(fmt.Sprintf("%s", data), Equals, "[defaults]\ntesting = value\n\n[othersection]\nmyoption = myvalue\nnewoption = novalue\n\n[testing]\nmyoption = value\n\n")
}
