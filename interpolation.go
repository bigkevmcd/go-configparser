package configparser

import (
	"github.com/bigkevmcd/configparser/chainmap"

	"strings"
)

func (p *ConfigParser) getInterpolated(section, option string, c *chainmap.ChainMap) (string, error) {
	val, err := p.Get(section, option)
	if err != nil {
		return "", err
	}
	return p.interpolate(val, c), nil
}

// return a string value for the named option.  All % interpolations are
// expanded in the return values, based on the defaults passed into the
// constructor and the DEFAULT section.  
func (p *ConfigParser) GetInterpolated(section, option string) (string, error) {
	o, err := p.Items(section)
	if err != nil {
		return "", err
	}
	c := chainmap.New(chainmap.Dict(p.Defaults()), chainmap.Dict(o))
	return p.getInterpolated(section, option, c)
}

// return a string value for the named option.  All % interpolations are
// expanded in the return values, based on the defaults passed into the
// constructor and the DEFAULT section.  Additional substitutions may be
// provided using the 'v' argument, which must be a Dict whose contents contents
// override any pre-existing defaults.
func (p *ConfigParser) GetInterpolatedWithVars(section, option string, v Dict) (string, error) {
	o, err := p.Items(section)
	if err != nil {
		return "", err
	}
	c := chainmap.New(chainmap.Dict(p.Defaults()), chainmap.Dict(o), chainmap.Dict(v))
	return p.getInterpolated(section, option, c)

}

// Private method which does the work of interpolating a value
// interpolates the value using the values in the ChainMap
// returns the interpolated string
func (p *ConfigParser) interpolate(value string, options *chainmap.ChainMap) string {

	for i := 0; i < maxInterpolationDepth; i++ {
		if strings.Contains(value, "%(") {
			value = interpolater.ReplaceAllStringFunc(value, func(m string) string {
				// No ReplaceAllStringSubMatchFunc so apply the regexp twice
				match := interpolater.FindAllStringSubmatch(m, 1)[0][1]
				replacement := options.Get(match)
				return replacement
			})
		}
	}
	return value
}

// return a copy of the dict for the section
func (p *ConfigParser) ItemsInterpolated(section string) (Dict, error) {
	s, err := p.Items(section)
	if err != nil {
		return nil, err
	}
	// TOOD: Optimise this...instantiate the ChainMap and delegate to interpolate()
	for k, v := range s {
		v, err = p.GetInterpolated(section, k)
		if err != nil {
			return nil, err
		}
		s[k] = v
	}
	return s, nil
}
