package configparser_test

import (
	"github.com/bigkevmcd/configparser"
	"testing"

	. "launchpad.net/gocheck"
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
