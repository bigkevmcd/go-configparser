package configparser_test

import (
	"testing"

	gc "gopkg.in/check.v1"

	"github.com/bigkevmcd/go-configparser"
)

func TestChainMap(t *testing.T) { gc.TestingT(t) }

type ChainMapSuite struct {
	c     configparser.Interpolator
	dict1 configparser.Dict
	dict2 configparser.Dict
}

var _ = gc.Suite(&ChainMapSuite{})

func (s *ChainMapSuite) SetUpTest(c *gc.C) {
	s.c = configparser.NewChainMap()
	s.dict1 = make(configparser.Dict)
	s.dict2 = make(configparser.Dict)
	s.dict1["testing"] = "2"
	s.dict1["value"] = "3"
	s.dict2["value"] = "4"
}

func (s *ChainMapSuite) TestLen(c *gc.C) {
	s.c.Add(s.dict1, s.dict2)
	c.Assert(s.c.Len(), gc.Equals, 2)
}

func (s *ChainMapSuite) TestGet1(c *gc.C) {
	s.c.Add(s.dict1, s.dict2)

	result := s.c.Get("unknown")
	c.Assert(result, gc.Equals, "")
}

func (s *ChainMapSuite) TestGet2(c *gc.C) {
	s.c.Add(s.dict1, s.dict2)

	result := s.c.Get("value")
	c.Assert(result, gc.Equals, "4")
	result = s.c.Get("testing")
	c.Assert(result, gc.Equals, "2")
}

func (s *ChainMapSuite) TestGet3(c *gc.C) {
	s.c.Add(s.dict2, s.dict1)

	result := s.c.Get("value")
	c.Assert(result, gc.Equals, "3")
}
