package configparser_test

import (
	"github.com/bigkevmcd/go-configparser"

	. "launchpad.net/gocheck"
)

// Defaults() should return a map containing the parser defaults.
func (s *ConfigParserSuite) TestDefaults(c *C) {
	d := s.p.Defaults()
	c.Assert(d["base_dir"], Equals, "/srv")
}

// Defaults() should return an empty Dict if there are no parser defaults
func (s *ConfigParserSuite) TestDefaultsWithNoDefaults(c *C) {
	p := configparser.New()
	d := p.Defaults()
	c.Assert(d, DeepEquals, configparser.Dict{})
}

// Sections() should return a list of section names excluding [DEFAULT]
func (s *ConfigParserSuite) TestSections(c *C) {
	result := s.p.Sections()
	c.Assert(result, DeepEquals, []string{"slave"})
}

// AddSection(section) should create a new section in the configuration
func (s *ConfigParserSuite) TestAddSection(c *C) {
	newParser := configparser.New()

	err := newParser.AddSection("newsection")

	c.Assert(err, IsNil)
	c.Assert(newParser.Sections(), DeepEquals, []string{"newsection"})
}

// AddSection(section) should return an appropriate error if the section already exists
func (s *ConfigParserSuite) TestAddSectionDuplicate(c *C) {
	err := s.p.AddSection("slave")

	c.Assert(err, ErrorMatches, "Section 'slave' already exists")
}

// AddSection(section) should return an appropriate error if we attempt to add a default section
func (s *ConfigParserSuite) TestAddSectionDefaultLowercase(c *C) {
	newParser := configparser.New()
	err := newParser.AddSection("default")

	c.Assert(err, ErrorMatches, "Invalid section name: 'default'")
}

// AddSection(section) should return an appropriate error if we attempt to add a DEFAULT section
func (s *ConfigParserSuite) TestAddSectionDefaultUppercase(c *C) {
	newParser := configparser.New()
	err := newParser.AddSection("DEFAULT")

	c.Assert(err, ErrorMatches, "Invalid section name: 'DEFAULT'")
}

// Options(section) should return an appropriate error if the section doesn't exist
func (s *ConfigParserSuite) TestOptionsWithNoSection(c *C) {
	_, err := s.p.Options("unknown")
	c.Assert(err, ErrorMatches, "No section: 'unknown'")
}

// Options(section) should return a list of option names for a given section mixed in with the defaults
func (s *ConfigParserSuite) TestOptionsWithSection(c *C) {
	result, err := s.p.Options("slave")
	c.Assert(err, IsNil)
	c.Assert(result, DeepEquals, []string{"base_dir", "bin_dir", "builder_command", "log_dir", "max_build_time"})
}

// Options(section) should return an empty slice if there are no options in a section
func (s *ConfigParserSuite) TestOptionsWithEmptySection(c *C) {
	newParser := configparser.New()
	newParser.AddSection("testing")
	result, err := newParser.Options("testing")
	c.Assert(err, IsNil)
	c.Assert(result, DeepEquals, []string{})
}

// Get(section, option) should return an appropriate error if the section does not exist
func (s *ConfigParserSuite) TestGetWithMissingSection(c *C) {
	_, err := s.p.Get("missing", "value")
	c.Assert(err, ErrorMatches, "No section: 'missing'")
}

// Get(section, option) should return an appropriate error if the option does not exist within the section
func (s *ConfigParserSuite) TestGetWithMissingOptionInSection(c *C) {
	_, err := s.p.Get("slave", "missing")
	c.Assert(err, ErrorMatches, "No option 'missing' in section: 'slave'")
}

// Get(section, option) should return the option value for the named section
func (s *ConfigParserSuite) TestGet(c *C) {
	result, err := s.p.Get("slave", "max_build_time")
	c.Assert(err, IsNil)
	c.Assert(result, Equals, "200")
}

// Get(section, option) should lookup the option in the DEFAULT section if requested
func (s *ConfigParserSuite) TestGetDefaultSection(c *C) {
	result, err := s.p.Get("DEFAULT", "bin_dir")
	c.Assert(err, IsNil)
	c.Assert(result, Equals, "%(base_dir)s/bin")
}

// Get(section, option) should lookup the option in the DEFAULT section if requested
func (s *ConfigParserSuite) TestGetDefaultSectionLowercase(c *C) {
	result, err := s.p.Get("default", "bin_dir")
	c.Assert(err, IsNil)
	c.Assert(result, Equals, "%(base_dir)s/bin")
}

// Get(section, option) should lookup the value in the default section if it doesn't exist in the section
func (s *ConfigParserSuite) TestGetWithMissingOptionInSectionButDefaultProvided(c *C) {
	result, err := s.p.Get("slave", "base_dir")
	c.Assert(err, IsNil)
	c.Assert(result, Equals, "/srv")
}

// Get(section, option) should be case insensitive with respect to options
func (s *ConfigParserSuite) TestGetCaseInsensitiveWithOptions(c *C) {
	result, err := s.p.Get("slave", "MAX_BUILD_TIME")
	c.Assert(err, IsNil)
	c.Assert(result, Equals, "200")
}

// Set(section, option, value) should return an error if the section doesn't exist
func (s *ConfigParserSuite) TestSetWithNoSection(c *C) {
	err := s.p.Set("unknown", "my_value", "testing")
	c.Assert(err, ErrorMatches, "No section: 'unknown'")
}

// Set(section, option, value) should set a default value if the section is the DEFAULT section
func (s *ConfigParserSuite) TestSetDefaultSection(c *C) {
	s.p.Set("DEFAULT", "my_value", "testing")
	defaults := s.p.Defaults()
	c.Assert(defaults["my_value"], Equals, "testing")
}

// Set(section, option, value) should record the specified value in the correct section
func (s *ConfigParserSuite) TestSet(c *C) {
	s.p.Set("slave", "my_value", "newvalue")
	result, err := s.p.Get("slave", "my_value")
	c.Assert(err, IsNil)
	c.Assert(result, Equals, "newvalue")
}

// HasSection(section) should return false if the section does not exist in the configuration
func (s *ConfigParserSuite) TestHasSectionWithoutSection(c *C) {
	newParser := configparser.New()
	c.Assert(newParser.HasSection("mysection"), Equals, false)
}

// HasSection(section) should return true if the section does exist in the configuration
func (s *ConfigParserSuite) TestHasSectionWithSection(c *C) {
	c.Assert(s.p.HasSection("slave"), Equals, true)
}

// Items(section) should return an appropriate error if the section doesn't exist
func (s *ConfigParserSuite) TestItemsWithNoSection(c *C) {
	_, err := s.p.Items("unknown")
	c.Assert(err, ErrorMatches, "No section: 'unknown'")
}

// Items(section) should return a copy of the dict for the section
func (s *ConfigParserSuite) TestItemsWithSection(c *C) {
	result, err := s.p.Items("slave")
	c.Assert(err, IsNil)
	c.Assert(result, DeepEquals, configparser.Dict{
		"max_build_time":  "200",
		"builder_command": "%(bin_dir)s/build",
		"log_dir":         "%(base_dir)s/logs"})
}

// Items(section) should return a copy of the dict for the section
func (s *ConfigParserSuite) TestItemsWithDefaults(c *C) {
	result, err := s.p.ItemsWithDefaults("slave")
	c.Assert(err, IsNil)
	c.Assert(result, DeepEquals, configparser.Dict{
		"max_build_time":  "200",
		"base_dir":        "/srv",
		"builder_command": "%(bin_dir)s/build",
		"log_dir":         "%(base_dir)s/logs",
		"bin_dir":         "%(base_dir)s/bin"})
}

// GetInt64(section, option) should return the option value for the named section as an Int64 value
func (s *ConfigParserSuite) TestGetInt64(c *C) {
	newParser := configparser.New()
	newParser.AddSection("testing")
	newParser.Set("testing", "value", "200")

	result, err := newParser.GetInt64("testing", "value")
	c.Assert(err, IsNil)
	c.Assert(result, Equals, int64(200))
}

// GetInt64(section, option) should return an appropriate error if the option does not exist
func (s *ConfigParserSuite) TestGetInt64MissingOption(c *C) {
	newParser := configparser.New()
	newParser.AddSection("testing")

	_, err := newParser.GetInt64("testing", "value")
	c.Assert(err, ErrorMatches, "No option 'value' in section: 'testing'")
}

// GetInt64(section, option) should return an appropriate error if the value can't be converted
func (s *ConfigParserSuite) TestGetInt64InvalidOption(c *C) {
	newParser := configparser.New()
	newParser.AddSection("testing")
	newParser.Set("testing", "value", "invalid")

	_, err := newParser.GetInt64("testing", "value")
	c.Assert(err, ErrorMatches, ".*invalid syntax.*")
}

// GetFloat64(section, option) should return the option value for the named section as a Float64 value
func (s *ConfigParserSuite) TestGetFloat64(c *C) {
	newParser := configparser.New()
	newParser.AddSection("testing")
	newParser.Set("testing", "value", "3.14159265")

	result, err := newParser.GetFloat64("testing", "value")
	c.Assert(err, IsNil)
	c.Assert(result, Equals, float64(3.14159265))
}

// GetFloat64(section, option) should return an appropriate error if the option does not exist
func (s *ConfigParserSuite) TestGetFloat64MissingOption(c *C) {
	newParser := configparser.New()
	newParser.AddSection("testing")

	_, err := newParser.GetFloat64("testing", "value")
	c.Assert(err, ErrorMatches, "No option 'value' in section: 'testing'")
}

// GetFloat64(section, option) should return an appropriate error if the value can't be converted
func (s *ConfigParserSuite) TestGetFloat64InvalidOption(c *C) {
	newParser := configparser.New()
	newParser.AddSection("testing")
	newParser.Set("testing", "value", "invalid")

	_, err := newParser.GetFloat64("testing", "value")
	c.Assert(err, ErrorMatches, ".*invalid syntax.*")
}

// GetBool(section, option) should return the option value for the named section as a Bool
func (s *ConfigParserSuite) TestGetBool(c *C) {
	newParser := configparser.New()
	newParser.AddSection("testing")

	for _, value := range []string{"1", "yes", "true", "on"} {
		newParser.Set("testing", "value", value)
		result, err := newParser.GetBool("testing", "value")
		c.Assert(err, IsNil)
		c.Assert(result, Equals, true)
	}

	for _, value := range []string{"0", "no", "false", "off"} {
		newParser.Set("testing", "value", value)
		result, err := newParser.GetBool("testing", "value")
		c.Assert(err, IsNil)
		c.Assert(result, Equals, false)
	}
}

// GetBool(section, option) should return an appropriate error if the value can't be converted
func (s *ConfigParserSuite) TestGetBoolInvalidValue(c *C) {
	newParser := configparser.New()
	newParser.AddSection("testing")
	newParser.Set("testing", "value", "testing")

	_, err := newParser.GetBool("testing", "value")
	c.Assert(err, ErrorMatches, "Not a boolean: 'testing'")
}

// RemoveSection(section) should return an appropriate error if the section doesn't exist
func (s *ConfigParserSuite) TestRemoveSectionMissingSection(c *C) {
	err := s.p.RemoveSection("unknown")
	c.Assert(err, ErrorMatches, "No section: 'unknown'")
}

// RemoveSection(section) should return an appropriate error if the section doesn't exist
func (s *ConfigParserSuite) TestRemoveSection(c *C) {
	newParser := configparser.New()
	newParser.AddSection("testing1")
	newParser.AddSection("testing2")
	err := newParser.RemoveSection("testing1")

	c.Assert(err, IsNil)
	result := newParser.Sections()
	c.Assert(result, DeepEquals, []string{"testing2"})
}

// // RemoveOption(section, option) should return an appropriate error if the section doesn't exist
// func (s *ConfigParserSuite) TestRemoveOptionMissingSection(c *C) {
// 	err := s.p.RemoveOption("unknown", "web")
// 	c.Assert(err, ErrorMatches, "No section: 'unknown'")
// }

// // RemoveOption(section, option) should return an appropriate error if the option doesn't exist
// func (s *ConfigParserSuite) TestRemoveOptionMissingOption(c *C) {
//  err := s.p.RemoveOption("web", "unknown")
//  c.Assert(err, ErrorMatches, "No option 'unknown' in section: 'web'")
// }

// HasOption(section, option) should return true if section is default and the option is a default
func (s *ConfigParserSuite) TestHasOptionFromDefaults(c *C) {
	result, err := s.p.HasOption("DEFAULT", "base_dir")
	c.Assert(err, IsNil)
	c.Assert(result, Equals, true)
}

// HasOption(section, option) should return an appropriate error if the section does not exist
func (s *ConfigParserSuite) TestHasOptionMissingSection(c *C) {
	_, err := s.p.HasOption("unknown", "missing")
	c.Assert(err, ErrorMatches, "No section: 'unknown'")
}
