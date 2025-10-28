package configparser

import (
	"strings"

	"github.com/bigkevmcd/go-configparser/chainmap"
)

const maxInterpolationDepth int = 10

func (p *ConfigParser) getInterpolated(section, option string, c Interpolator) (string, error) {
	val, err := p.get(section, option)
	if err != nil {
		return "", err
	}
	return p.interpolate(val, c), nil
}

// GetInterpolated returns a string value for the named option.
//
// All % interpolations are expanded in the return values, based on
// the defaults passed into the constructor and the DEFAULT section.
func (p *ConfigParser) GetInterpolated(section, option string) (string, error) {
	o, err := p.Items(section)
	if err != nil {
		return "", err
	}
	// Create a fresh interpolator for each call to avoid accumulating state.
	// If using the default ChainMap interpolator, create a new one.
	// If using a custom interpolator, use it directly (and add dicts to it).
	var interp Interpolator
	if _, isChainMap := p.opt.interpolation.(*chainmap.ChainMap); isChainMap {
		// Default case: create fresh ChainMap to avoid mutation
		interp = chainmap.New(chainmap.Dict(p.Defaults()), chainmap.Dict(o))
	} else {
		// Custom interpolator: use it directly (legacy behavior for compatibility)
		p.opt.interpolation.Add(chainmap.Dict(p.Defaults()), chainmap.Dict(o))
		interp = p.opt.interpolation
	}
	return p.getInterpolated(section, option, interp)
}

// GetInterpolatedWithVars returns a string value for the named option.
//
// All % interpolations are expanded in the return values, based on the defaults passed
// into the constructor and the DEFAULT section.  Additional substitutions may be
// provided using the 'v' argument, which must be a Dict whose contents contents
// override any pre-existing defaults.
func (p *ConfigParser) GetInterpolatedWithVars(section, option string, v Dict) (string, error) {
	o, err := p.Items(section)
	if err != nil {
		return "", err
	}
	// Create a fresh interpolator for each call to avoid accumulating state.
	// If using the default ChainMap interpolator, create a new one.
	// If using a custom interpolator, use it directly (and add dicts to it).
	var interp Interpolator
	if _, isChainMap := p.opt.interpolation.(*chainmap.ChainMap); isChainMap {
		// Default case: create fresh ChainMap to avoid mutation
		interp = chainmap.New(chainmap.Dict(p.Defaults()), chainmap.Dict(o), chainmap.Dict(v))
	} else {
		// Custom interpolator: use it directly (legacy behavior for compatibility)
		p.opt.interpolation.Add(chainmap.Dict(p.Defaults()), chainmap.Dict(o), chainmap.Dict(v))
		interp = p.opt.interpolation
	}

	return p.getInterpolated(section, option, interp)
}

// Private method which does the work of interpolating a value
// interpolates the value using the values in the ChainMap
// returns the interpolated string.
func (p *ConfigParser) interpolate(value string, options Interpolator) string {
	for i := 0; i < maxInterpolationDepth; i++ {
		if !strings.Contains(value, "%(") {
			break
		}

		var changed bool
		value = interpolater.ReplaceAllStringFunc(value, func(m string) string {
			ms := interpolater.FindStringSubmatch(m)
			if len(ms) < 2 {
				return m
			}

			match := ms[1]
			replacement := options.Get(match)
			if replacement != m {
				changed = true
			}

			return replacement
		})

		if !changed {
			break
		}
	}

	return value
}

// ItemsWithDefaultsInterpolated returns a copy of the dict for the section.
func (p *ConfigParser) ItemsWithDefaultsInterpolated(section string) (Dict, error) {
	s, err := p.ItemsWithDefaults(section)
	if err != nil {
		return nil, err
	}
	// TODO: Optimise this... instantiate the ChainMap and delegate to interpolate()
	for k := range s {
		v, err := p.GetInterpolated(section, k)
		if err != nil {
			return nil, err
		}
		s[k] = v
	}
	return s, nil
}
