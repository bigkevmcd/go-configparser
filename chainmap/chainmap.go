package chainmap

// Dict is a simple string->string map.
type Dict map[string]string

// ChainMap contains a slice of Dicts for interpolation values.
type ChainMap struct {
	maps []Dict
}

// New creates a new ChainMap.
func New(dicts ...Dict) *ChainMap {
	chainMap := &ChainMap{
		maps: make([]Dict, 0),
	}
	chainMap.maps = append(chainMap.maps, dicts...)

	return chainMap
}

// Add adds given dicts to the ChainMap.
func (c *ChainMap) Add(dicts ...Dict) {
	c.maps = append(c.maps, dicts...)
}

// Len returns the amount of Dicts in the ChainMap.
func (c *ChainMap) Len() int {
	return len(c.maps)
}

// Get returns the value for the given key from the ChainMap.
// The last-added dicts have priority, so this checks maps from last to
// first and returns the first match. If key does not exist, returns empty string.
func (c *ChainMap) Get(key string) string {
	for i := len(c.maps) - 1; i >= 0; i-- {
		if result, present := c.maps[i][key]; present {
			return result
		}
	}
	return ""
}
