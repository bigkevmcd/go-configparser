package chainmap_test

import (
	"github.com/bigkevmcd/go-configparser/chainmap"
	"testing"

	gc "gopkg.in/check.v1"
)

func Test(t *testing.T) { gc.TestingT(t) }

type ChainMapSuite struct {
	dict1 chainmap.Dict
	dict2 chainmap.Dict
}

var _ = gc.Suite(&ChainMapSuite{})

func (s *ChainMapSuite) SetUpTest(c *gc.C) {
	s.dict1 = make(chainmap.Dict)
	s.dict2 = make(chainmap.Dict)
	s.dict1["testing"] = "2"
	s.dict1["value"] = "3"
	s.dict2["value"] = "4"
}

func (s *ChainMapSuite) TestLen(c *gc.C) {
	chainMap := chainmap.New(s.dict1, s.dict2)
	c.Assert(chainMap.Len(), gc.Equals, 2)
}

func (s *ChainMapSuite) TestGet1(c *gc.C) {
	chainMap := chainmap.New(s.dict1, s.dict2)

	result := chainMap.Get("unknown")
	c.Assert(result, gc.Equals, "")
}

func (s *ChainMapSuite) TestGet2(c *gc.C) {
	chainMap := chainmap.New(s.dict1, s.dict2)

	result := chainMap.Get("value")
	c.Assert(result, gc.Equals, "4")
	result = chainMap.Get("testing")
	c.Assert(result, gc.Equals, "2")
}

func (s *ChainMapSuite) TestGet3(c *gc.C) {
	chainMap := chainmap.New(s.dict2, s.dict1)

	result := chainMap.Get("value")
	c.Assert(result, gc.Equals, "3")
}
