package configparser_test

import (
	"github.com/bigkevmcd/go-configparser"

	gc "gopkg.in/check.v1"
)

// Defaults() should return a map containing the parser defaults.
func (s *ConfigParserSuite) TestDefaults(c *gc.C) {
	d := s.p.Defaults()
	c.Assert(d["base_dir"], gc.Equals, "/srv")
}

// Defaults() should return an empty Dict if there are no parser defaults
func (s *ConfigParserSuite) TestDefaultsWithNoDefaults(c *gc.C) {
	p := configparser.New()
	d := p.Defaults()
	c.Assert(d, gc.DeepEquals, configparser.Dict{})
}

// Sections() should return a list of section names excluding [DEFAULT]
func (s *ConfigParserSuite) TestSections(c *gc.C) {
	result := s.p.Sections()
	c.Assert(result, gc.DeepEquals, []string{"follower", "whitespace"})
}

// AddSection(section) should create a new section in the configuration
func (s *ConfigParserSuite) TestAddSection(c *gc.C) {
	newParser := configparser.New()

	err := newParser.AddSection("newsection")

	c.Assert(err, gc.IsNil)
	c.Assert(newParser.Sections(), gc.DeepEquals, []string{"newsection"})
}

// AddSection(section) should return an appropriate error if the section already exists
func (s *ConfigParserSuite) TestAddSectionDuplicate(c *gc.C) {
	err := s.p.AddSection("follower")

	c.Assert(err, gc.ErrorMatches, "Section 'follower' already exists")
}

// AddSection(section) should return an appropriate error if we attempt to add a default section
func (s *ConfigParserSuite) TestAddSectionDefaultLowercase(c *gc.C) {
	newParser := configparser.New()
	err := newParser.AddSection("default")

	c.Assert(err, gc.ErrorMatches, "Invalid section name: 'default'")
}

// AddSection(section) should return an appropriate error if we attempt to add a DEFAULT section
func (s *ConfigParserSuite) TestAddSectionDefaultUppercase(c *gc.C) {
	newParser := configparser.New()
	err := newParser.AddSection("DEFAULT")

	c.Assert(err, gc.ErrorMatches, "Invalid section name: 'DEFAULT'")
}

// Options(section) should return an appropriate error if the section doesn't exist
func (s *ConfigParserSuite) TestOptionsWithNoSection(c *gc.C) {
	_, err := s.p.Options("unknown")
	c.Assert(err, gc.ErrorMatches, "No section: 'unknown'")
}

// Options(section) should return a list of option names for a given section mixed in with the defaults
func (s *ConfigParserSuite) TestOptionsWithSection(c *gc.C) {
	result, err := s.p.Options("follower")
	c.Assert(err, gc.IsNil)
	c.Assert(result, gc.DeepEquals, []string{"FrobTimeout", "TableName", "base_dir", "bin_dir", "builder_command", "log_dir", "max_build_time"})
}

// Options(section) should return an empty slice if there are no options in a section
func (s *ConfigParserSuite) TestOptionsWithEmptySection(c *gc.C) {
	newParser := configparser.New()
	newParser.AddSection("testing")
	result, err := newParser.Options("testing")
	c.Assert(err, gc.IsNil)
	c.Assert(result, gc.DeepEquals, []string{})
}

// Get(section, option) should return an appropriate error if the section does not exist
func (s *ConfigParserSuite) TestGetWithMissingSection(c *gc.C) {
	_, err := s.p.Get("missing", "value")
	c.Assert(err, gc.ErrorMatches, "No section: 'missing'")
}

// Get(section, option) should return an appropriate error if the option does not exist within the section
func (s *ConfigParserSuite) TestGetWithMissingOptionInSection(c *gc.C) {
	_, err := s.p.Get("follower", "missing")
	c.Assert(err, gc.ErrorMatches, "No option 'missing' in section: 'follower'")
}

// Get(section, option) should return the option value for the named section
func (s *ConfigParserSuite) TestGet(c *gc.C) {
	result, err := s.p.Get("follower", "max_build_time")
	c.Assert(err, gc.IsNil)
	c.Assert(result, gc.Equals, "200")
}

// Get(section, option) should return the option value for the named section
// regardless of case
func (s *ConfigParserSuite) TestGetCamelCase(c *gc.C) {
	result, err := s.p.Get("follower", "FrobTimeout")
	c.Assert(err, gc.IsNil)
	c.Assert(result, gc.Equals, "5")
}

// Get(section, option) should return the option value for the named section
// without mangling the value's case
func (s *ConfigParserSuite) TestValueCasePreservation(c *gc.C) {
	result, err := s.p.Get("follower", "TableName")
	c.Assert(err, gc.IsNil)
	c.Assert(result, gc.Equals, "MyCaseSensitiveTableName")
}

// Get(section, option) should lookup the option in the DEFAULT section if requested
func (s *ConfigParserSuite) TestGetDefaultSection(c *gc.C) {
	result, err := s.p.Get("DEFAULT", "bin_dir")
	c.Assert(err, gc.IsNil)
	c.Assert(result, gc.Equals, "%(base_dir)s/bin")
}

// Get(section, option) should lookup the option in the DEFAULT section if requested
func (s *ConfigParserSuite) TestGetDefaultSectionLowercase(c *gc.C) {
	result, err := s.p.Get("default", "bin_dir")
	c.Assert(err, gc.IsNil)
	c.Assert(result, gc.Equals, "%(base_dir)s/bin")
}

// Get(section, option) should lookup the value in the default section if it doesn't exist in the section
func (s *ConfigParserSuite) TestGetWithMissingOptionInSectionButDefaultProvided(c *gc.C) {
	result, err := s.p.Get("follower", "base_dir")
	c.Assert(err, gc.IsNil)
	c.Assert(result, gc.Equals, "/srv")
}

// Get(section, option) should be case insensitive with respect to options
func (s *ConfigParserSuite) TestGetCaseInsensitiveWithOptions(c *gc.C) {
	result, err := s.p.Get("follower", "MAX_BUILD_TIME")
	c.Assert(err, gc.IsNil)
	c.Assert(result, gc.Equals, "200")
}

// Set(section, option, value) should return an error if the section doesn't exist
func (s *ConfigParserSuite) TestSetWithNoSection(c *gc.C) {
	err := s.p.Set("unknown", "my_value", "testing")
	c.Assert(err, gc.ErrorMatches, "No section: 'unknown'")
}

// Set(section, option, value) should set a default value if the section is the DEFAULT section
func (s *ConfigParserSuite) TestSetDefaultSection(c *gc.C) {
	s.p.Set("DEFAULT", "my_value", "testing")
	defaults := s.p.Defaults()
	c.Assert(defaults["my_value"], gc.Equals, "testing")
}

// Set(section, option, value) should record the specified value in the correct section
func (s *ConfigParserSuite) TestSet(c *gc.C) {
	s.p.Set("follower", "my_value", "newvalue")
	result, err := s.p.Get("follower", "my_value")
	c.Assert(err, gc.IsNil)
	c.Assert(result, gc.Equals, "newvalue")
}

// HasSection(section) should return false if the section does not exist in the configuration
func (s *ConfigParserSuite) TestHasSectionWithoutSection(c *gc.C) {
	newParser := configparser.New()
	c.Assert(newParser.HasSection("mysection"), gc.Equals, false)
}

// HasSection(section) should return true if the section does exist in the configuration
func (s *ConfigParserSuite) TestHasSectionWithSection(c *gc.C) {
	c.Assert(s.p.HasSection("follower"), gc.Equals, true)
}

// Items(section) should return an appropriate error if the section doesn't exist
func (s *ConfigParserSuite) TestItemsWithNoSection(c *gc.C) {
	_, err := s.p.Items("unknown")
	c.Assert(err, gc.ErrorMatches, "No section: 'unknown'")
}

// Items(section) should return a copy of the dict for the section
func (s *ConfigParserSuite) TestItemsWithSection(c *gc.C) {
	result, err := s.p.Items("follower")
	c.Assert(err, gc.IsNil)
	c.Assert(result, gc.DeepEquals, configparser.Dict{
		"FrobTimeout":     "5",
		"TableName":       "MyCaseSensitiveTableName",
		"max_build_time":  "200",
		"builder_command": "%(bin_dir)s/build",
		"log_dir":         "%(base_dir)s/logs"})
}

// Items(section) should return a copy of the dict for the section
func (s *ConfigParserSuite) TestItemsWithDefaults(c *gc.C) {
	result, err := s.p.ItemsWithDefaults("follower")
	c.Assert(err, gc.IsNil)
	c.Assert(result, gc.DeepEquals, configparser.Dict{
		"FrobTimeout":     "5",
		"TableName":       "MyCaseSensitiveTableName",
		"max_build_time":  "200",
		"base_dir":        "/srv",
		"builder_command": "%(bin_dir)s/build",
		"log_dir":         "%(base_dir)s/logs",
		"bin_dir":         "%(base_dir)s/bin"})
}

// GetInt64(section, option) should return the option value for the named section as an Int64 value
func (s *ConfigParserSuite) TestGetInt64(c *gc.C) {
	newParser := configparser.New()
	newParser.AddSection("testing")
	newParser.Set("testing", "value", "200")

	result, err := newParser.GetInt64("testing", "value")
	c.Assert(err, gc.IsNil)
	c.Assert(result, gc.Equals, int64(200))
}

// GetInt64(section, option) should return an appropriate error if the option does not exist
func (s *ConfigParserSuite) TestGetInt64MissingOption(c *gc.C) {
	newParser := configparser.New()
	newParser.AddSection("testing")

	_, err := newParser.GetInt64("testing", "value")
	c.Assert(err, gc.ErrorMatches, "No option 'value' in section: 'testing'")
}

// GetInt64(section, option) should return an appropriate error if the value can't be converted
func (s *ConfigParserSuite) TestGetInt64InvalidOption(c *gc.C) {
	newParser := configparser.New()
	newParser.AddSection("testing")
	newParser.Set("testing", "value", "invalid")

	_, err := newParser.GetInt64("testing", "value")
	c.Assert(err, gc.ErrorMatches, ".*invalid syntax.*")
}

// GetFloat64(section, option) should return the option value for the named section as a Float64 value
func (s *ConfigParserSuite) TestGetFloat64(c *gc.C) {
	newParser := configparser.New()
	newParser.AddSection("testing")
	newParser.Set("testing", "value", "3.14159265")

	result, err := newParser.GetFloat64("testing", "value")
	c.Assert(err, gc.IsNil)
	c.Assert(result, gc.Equals, float64(3.14159265))
}

// GetFloat64(section, option) should return an appropriate error if the option does not exist
func (s *ConfigParserSuite) TestGetFloat64MissingOption(c *gc.C) {
	newParser := configparser.New()
	newParser.AddSection("testing")

	_, err := newParser.GetFloat64("testing", "value")
	c.Assert(err, gc.ErrorMatches, "No option 'value' in section: 'testing'")
}

// GetFloat64(section, option) should return an appropriate error if the value can't be converted
func (s *ConfigParserSuite) TestGetFloat64InvalidOption(c *gc.C) {
	newParser := configparser.New()
	newParser.AddSection("testing")
	newParser.Set("testing", "value", "invalid")

	_, err := newParser.GetFloat64("testing", "value")
	c.Assert(err, gc.ErrorMatches, ".*invalid syntax.*")
}

// GetBool(section, option) should return the option value for the named section as a Bool
func (s *ConfigParserSuite) TestGetBool(c *gc.C) {
	newParser := configparser.New()
	newParser.AddSection("testing")

	for _, value := range []string{"1", "yes", "true", "on"} {
		newParser.Set("testing", "value", value)
		result, err := newParser.GetBool("testing", "value")
		c.Assert(err, gc.IsNil)
		c.Assert(result, gc.Equals, true)
	}

	for _, value := range []string{"0", "no", "false", "off"} {
		newParser.Set("testing", "value", value)
		result, err := newParser.GetBool("testing", "value")
		c.Assert(err, gc.IsNil)
		c.Assert(result, gc.Equals, false)
	}
}

// GetBool(section, option) should return an appropriate error if the value can't be converted
func (s *ConfigParserSuite) TestGetBoolInvalidValue(c *gc.C) {
	newParser := configparser.New()
	newParser.AddSection("testing")
	newParser.Set("testing", "value", "testing")

	_, err := newParser.GetBool("testing", "value")
	c.Assert(err, gc.ErrorMatches, "Not a boolean: 'testing'")
}

// RemoveSection(section) should return an appropriate error if the section doesn't exist
func (s *ConfigParserSuite) TestRemoveSectionMissingSection(c *gc.C) {
	err := s.p.RemoveSection("unknown")
	c.Assert(err, gc.ErrorMatches, "No section: 'unknown'")
}

// RemoveSection(section) should return an appropriate error if the section doesn't exist
func (s *ConfigParserSuite) TestRemoveSection(c *gc.C) {
	newParser := configparser.New()
	newParser.AddSection("testing1")
	newParser.AddSection("testing2")
	err := newParser.RemoveSection("testing1")

	c.Assert(err, gc.IsNil)
	result := newParser.Sections()
	c.Assert(result, gc.DeepEquals, []string{"testing2"})
}

// RemoveOption(section, option) should return an appropriate error if the section doesn't exist
func (s *ConfigParserSuite) TestRemoveOptionMissingSection(c *gc.C) {
	err := s.p.RemoveOption("unknown", "web")
	c.Assert(err, gc.ErrorMatches, "No section: 'unknown'")
}

// RemoveOption(section, option) should return an appropriate error if the option doesn't exist
func (s *ConfigParserSuite) TestRemoveOptionMissingOption(c *gc.C) {
	err := s.p.RemoveOption("follower", "unknown")
	c.Assert(err, gc.ErrorMatches, "No option 'unknown' in section: 'follower'")
}

// RemoveOption(section, option) should remove an option from the specified
// section.
func (s *ConfigParserSuite) TestRemoveOption(c *gc.C) {
	err := s.p.RemoveOption("follower", "max_build_time")
	c.Assert(err, gc.IsNil)
	hasOption, err := s.p.HasOption("follower", "max_build_time")
	c.Assert(err, gc.IsNil)
	c.Assert(hasOption, gc.Equals, false)
}

// RemoveOption(section, option) does not remove options when the option doesn't
// match the specified option exactly.
func (s *ConfigParserSuite) TestRemoveOptionMatchesPrecisely(c *gc.C) {
	err := s.p.RemoveOption("follower", "max_build_TIME")
	c.Assert(err, gc.ErrorMatches, "No option 'max_build_TIME' in section: 'follower'")
}

// HasOption(section, option) should return true if section is default and the option is a default
func (s *ConfigParserSuite) TestHasOptionFromDefaults(c *gc.C) {
	result, err := s.p.HasOption("DEFAULT", "base_dir")
	c.Assert(err, gc.IsNil)
	c.Assert(result, gc.Equals, true)
}

// HasOption(section, option) should return an appropriate error if the section does not exist
func (s *ConfigParserSuite) TestHasOptionMissingSection(c *gc.C) {
	_, err := s.p.HasOption("unknown", "missing")
	c.Assert(err, gc.ErrorMatches, "No section: 'unknown'")
}

// Options(section) should strip whitespace from the keys when parsing sections.
func (s *ConfigParserSuite) TestOptionsWithSectionStripsWhitespaceFromKeys(c *gc.C) {
	result, err := s.p.Options("whitespace")
	c.Assert(err, gc.IsNil)
	c.Assert(result, gc.DeepEquals, []string{"base_dir", "bin_dir", "foo"})
}
