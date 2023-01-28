package configparser

type chainMap struct {
	maps []Dict
}

// NewChainMap default interpolator.
func NewChainMap() *chainMap {
	return &chainMap{
		maps: make([]Dict, 0),
	}
}

func (c *chainMap) Add(dicts ...Dict) {
	c.maps = append(c.maps, dicts...)
}

func (c *chainMap) Len() int {
	return len(c.maps)
}

func (c *chainMap) Get(key string) string {
	var value string

	for _, dict := range c.maps {
		if result, present := dict[key]; present {
			value = result
		}
	}
	return value
}
