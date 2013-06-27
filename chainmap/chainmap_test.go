package chainmap_test

import (
	"github.com/bigkevmcd/go-configparser/chainmap"
	"testing"

	. "launchpad.net/gocheck"
)

func Test(t *testing.T) { TestingT(t) }

type ChainMapSuite struct {
	dict1 chainmap.Dict
	dict2 chainmap.Dict
}

var _ = Suite(&ChainMapSuite{})

func (s *ChainMapSuite) SetUpTest(c *C) {
	s.dict1 = make(chainmap.Dict)
	s.dict2 = make(chainmap.Dict)
	s.dict1["testing"] = "2"
	s.dict1["value"] = "3"
	s.dict2["value"] = "4"
}

func (s *ChainMapSuite) TestLen(c *C) {
	chainMap := chainmap.New(s.dict1, s.dict2)
	c.Assert(chainMap.Len(), Equals, 2)
}

func (s *ChainMapSuite) TestGet1(c *C) {
	chainMap := chainmap.New(s.dict1, s.dict2)

	result := chainMap.Get("unknown")
	c.Assert(result, Equals, "")
}

func (s *ChainMapSuite) TestGet2(c *C) {
	chainMap := chainmap.New(s.dict1, s.dict2)

	result := chainMap.Get("value")
	c.Assert(result, Equals, "4")
	result = chainMap.Get("testing")
	c.Assert(result, Equals, "2")
}

func (s *ChainMapSuite) TestGet3(c *C) {
	chainMap := chainmap.New(s.dict2, s.dict1)

	result := chainMap.Get("value")
	c.Assert(result, Equals, "3")
}
