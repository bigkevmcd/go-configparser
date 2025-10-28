package configparser_test

import (
	"strings"

	"github.com/bigkevmcd/go-configparser"
	"github.com/bigkevmcd/go-configparser/chainmap"

	. "gopkg.in/check.v1"
)

// GetInterpolated(section, option) should return an appropriate error if the section does not exist
func (s *ConfigParserSuite) TestGetInterpolatedWithMissingSection(c *C) {
	_, err := s.p.GetInterpolated("unknown", "missing")
	c.Assert(err, ErrorMatches, "no section: \"unknown\"")
}

// GetInterpolated(section, option) should interpolate the result
func (s *ConfigParserSuite) TestGetInterpolated(c *C) {
	result, err := s.p.GetInterpolated("follower", "builder_command")

	c.Assert(err, IsNil)
	c.Assert(result, Equals, "/srv/bin/build")
}

// GetInterpolatedWithVars(section, option, vars) should interpolate the result
// with the additional variables provided
func (s *ConfigParserSuite) TestGetInterpolatedWithVars(c *C) {
	d := make(configparser.Dict)
	d["bin_dir"] = "/a/non/existent/path"

	result, err := s.p.GetInterpolatedWithVars("follower", "builder_command", d)

	c.Assert(err, IsNil)
	c.Assert(result, Equals, "/a/non/existent/path/build")
}

// ItemsInterpolated(section) should return a copy of the section Dict
// but with the values interpolated
func (s *ConfigParserSuite) TestItemsWithDefaultsInterpolated(c *C) {
	result, err := s.p.ItemsWithDefaultsInterpolated("follower")

	c.Assert(err, IsNil)
	c.Assert(result, DeepEquals, configparser.Dict{
		"builder_command": "/srv/bin/build",
		"bin_dir":         "/srv/bin",
		"FrobTimeout":     "5",
		"TableName":       "MyCaseSensitiveTableName",
		"max_build_time":  "200",
		"log_dir":         "/srv/logs",
		"base_dir":        "/srv",
	})
}

// GetInterpolated should not accumulate state across calls
// This test verifies that multiple calls to GetInterpolated don't
// mutate the shared interpolator in a way that affects subsequent calls
func (s *ConfigParserSuite) TestGetInterpolatedDoesNotMutateGlobalState(c *C) {
	// Create a custom interpolator that we can inspect
	customInterp := chainmap.New()
	p, err := configparser.ParseReaderWithOptions(
		strings.NewReader(`[DEFAULT]
base_dir = /srv
bin_dir = %(base_dir)s/bin

[follower]
builder_command = %(bin_dir)s/build
`),
		configparser.Interpolation(customInterp),
	)
	c.Assert(err, IsNil)

	// Check initial state
	initialLen := customInterp.Len()

	// First call
	result1, err := p.GetInterpolated("follower", "builder_command")
	c.Assert(err, IsNil)
	c.Assert(result1, Equals, "/srv/bin/build")
	lenAfterFirst := customInterp.Len()

	// Second call - the interpolator should NOT grow
	result2, err := p.GetInterpolated("follower", "builder_command")
	c.Assert(err, IsNil)
	c.Assert(result2, Equals, "/srv/bin/build")
	lenAfterSecond := customInterp.Len()

	c.Assert(lenAfterSecond, Equals, lenAfterFirst,
		Commentf("Interpolator grew from %d to %d dicts after second call (started at %d)",
			lenAfterFirst, lenAfterSecond, initialLen))
}
