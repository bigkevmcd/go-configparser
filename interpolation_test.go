package configparser_test

import (
	"github.com/bigkevmcd/configparser"

	. "launchpad.net/gocheck"
)

// GetInterpolated(section, option) should return an appropriate error if the section does not exist
func (s *ConfigParserSuite) TestGetInterpolatedWithMissingSection(c *C) {
	_, err := s.p.GetInterpolated("unknown", "missing")
	c.Assert(err, ErrorMatches, "No section: 'unknown'")
}

// GetInterpolated(section, option) should interpolate the result
func (s *ConfigParserSuite) TestGetInterpolated(c *C) {
	result, err := s.p.GetInterpolated("slave", "builder_command")

	c.Assert(err, IsNil)
	c.Assert(result, Equals, "/srv/bin/build")
}

// GetInterpolatedWithVars(section, option, vars) should interpolate the result
// with the additional variables provided
func (s *ConfigParserSuite) TestGetInterpolatedWithVars(c *C) {
	d := make(configparser.Dict)
	d["bin_dir"] = "/a/non/existent/path"

	result, err := s.p.GetInterpolatedWithVars("slave", "builder_command", d)

	c.Assert(err, IsNil)
	c.Assert(result, Equals, "/a/non/existent/path/build")
}

// ItemsInterpolated(section) should return a copy of the section Dict
// but with the values interpolated
func (s *ConfigParserSuite) TestItemsInterpolated(c *C) {
	result, err := s.p.ItemsInterpolated("slave")

	c.Assert(err, IsNil)
	c.Assert(result, DeepEquals, configparser.Dict{
		"builder_command": "/srv/bin/build",
		"bin_dir":         "/srv/bin",
		"max_build_time":  "200",
		"log_dir":         "/srv/logs",
		"base_dir":        "/srv"})
}
