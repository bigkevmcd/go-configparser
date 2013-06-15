package chainmap

type Dict map[string]string

type ChainMap struct {
	maps []Dict
}

func New(dicts ...Dict) *ChainMap {
	chainMap := &ChainMap{
		maps: make([]Dict, 0),
	}

	for _, dict := range dicts {
		chainMap.maps = append(chainMap.maps, dict)
	}
	return chainMap
}

func (c *ChainMap) Len() int {
	return len(c.maps)
}

func (c *ChainMap) Get(key string) string {
	var value string

	for _, dict := range c.maps {
		if result, present := dict[key]; present {
			value = result
		}
	}
	return value
}
