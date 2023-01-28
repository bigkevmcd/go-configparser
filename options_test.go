package configparser_test

import (
	"strconv"
	"strings"

	. "gopkg.in/check.v1"

	"github.com/bigkevmcd/go-configparser"
	"github.com/bigkevmcd/go-configparser/chainmap"
)

func (s *ConfigParserSuite) TestInterpolationOpt(c *C) {
	parsed, err := configparser.ParseReaderWithOptions(
		strings.NewReader("[DEFAULT]\ndir=/home\n[paths]\npath=%(dir)s/something\n\n"),
		configparser.Interpolation(chainmap.New()),
	)
	c.Assert(err, IsNil)

	v, err := parsed.GetInterpolated("paths", "path")
	c.Assert(err, IsNil)
	c.Assert(v, Equals, "/home/something")
}

func (s *ConfigParserSuite) TestCommentPrefixesOpt(c *C) {
	parsed, err := configparser.ParseReaderWithOptions(
		strings.NewReader("[section]\n// this is a comment\noption=value\n\n"),
		configparser.CommentPrefixes(configparser.Prefixes{"//"}),
	)
	c.Assert(err, IsNil)

	opt, err := parsed.Options("section")
	c.Assert(err, IsNil)
	c.Assert(len(opt), Equals, 1)
}

func (s *ConfigParserSuite) TestInlineCommentPrefixesOpt(c *C) {
	parsed, err := configparser.ParseReaderWithOptions(
		strings.NewReader("[section] // this is section inline comment\noption=value // this is an inline comment\n\n"),
		configparser.InlineCommentPrefixes(configparser.Prefixes{"//"}),
	)
	c.Assert(err, IsNil)

	v, err := parsed.Get("section", "option")
	c.Assert(err, IsNil)
	c.Assert(v, Equals, "value")
}

func (s *ConfigParserSuite) TestDefalutSectionOpt(c *C) {
	parsed, err := configparser.ParseReaderWithOptions(
		strings.NewReader("[NEW DEFAULT]\noption=value\n\n"),
		configparser.DefaultSection("NEW DEFAULT"),
	)
	c.Assert(err, IsNil)

	keys := parsed.Defaults().Keys()
	c.Assert(len(keys), Equals, 1)

	v, err := parsed.Get("NEW DEFAULT", "option")
	c.Assert(err, IsNil)
	c.Assert(v, Equals, "value")
}

func (s *ConfigParserSuite) TestDelimetersOpt(c *C) {
	parsed, err := configparser.ParseReaderWithOptions(
		strings.NewReader("[section]\noption==test\n\n"),
		configparser.Delimiters("=="),
	)
	c.Assert(err, IsNil)

	v, err := parsed.Get("section", "option")
	c.Assert(err, IsNil)
	c.Assert(v, Equals, "test")
}

func (s *ConfigParserSuite) TestConvertersOpt(c *C) {
	intConv := func(s string) (any, error) {
		i, err := strconv.Atoi(s)
		if err != nil {
			return -1, err
		}

		return i + 1, err
	}

	floatConv := func(s string) (any, error) {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return -1, err
		}
		return f + 1, nil
	}

	stringConv := func(s string) (any, error) {
		return s + "_updated", nil
	}

	boolConv := func(s string) (any, error) {
		return s != "", nil
	}

	conv := configparser.Converter{
		"int":    intConv,
		"float":  floatConv,
		"string": stringConv,
		"bool":   boolConv,
	}

	parsed, err := configparser.ParseReaderWithOptions(
		strings.NewReader("[section]\nint=1\nfloat=1.1\nstring=test\nbool\n\n"),
		configparser.Converters(conv),
		configparser.AllowNoValue,
	)
	c.Assert(err, IsNil)

	pInt, err := parsed.GetInt64("section", "int")
	c.Assert(err, IsNil)
	c.Assert(pInt, Equals, int64(2))

	pFloat, err := parsed.GetFloat64("section", "float")
	c.Assert(err, IsNil)
	c.Assert(pFloat, Equals, 2.1)

	pString, err := parsed.Get("section", "string")
	c.Assert(err, IsNil)
	c.Assert(pString, Equals, "test_updated")

	pBool, err := parsed.GetBool("section", "bool")
	c.Assert(err, IsNil)
	c.Assert(pBool, Equals, false)
}

func (s *ConfigParserSuite) TestAllowNoValueOptParsedFromReader(c *C) {
	parsed, err := configparser.ParseReaderWithOptions(
		strings.NewReader("[empty]\noption\n\n"), configparser.AllowNoValue,
	)
	c.Assert(err, IsNil)

	ok, err := parsed.HasOption("empty", "option")
	c.Assert(err, IsNil)
	c.Assert(ok, Equals, true)
}

func (s *ConfigParserSuite) TestAllowNoValueOptParsedFromFile(c *C) {
	parsed, err := configparser.ParseWithOptions(
		"testdata/example.cfg", configparser.AllowNoValue,
	)
	c.Assert(err, IsNil)

	ok, err := parsed.HasOption("empty", "foo")
	c.Assert(err, IsNil)
	c.Assert(ok, Equals, true)
}

func (s *ConfigParserSuite) TestStrictOptDuplicateSection(c *C) {
	_, err := configparser.ParseReaderWithOptions(
		strings.NewReader("[dubl]\noption=1\n\n[dubl]\noption=2\n\n"),
		configparser.Strict,
	)

	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "section \"dubl\" error: already exist")
}

func (s *ConfigParserSuite) TestStrictOptDuplicateOption(c *C) {
	_, err := configparser.ParseReaderWithOptions(
		strings.NewReader("[section1]\ndubl=1\n\n[section2]\ndubl=2\n\n"),
		configparser.Strict,
	)

	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "option \"dubl\" error: already exist")
}
