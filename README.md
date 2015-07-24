# go-configparser [![Build Status](https://travis-ci.org/bigkevmcd/go-configparser.png)](https://travis-ci.org/bigkevmcd/go-configparser)
Go implementation of the Python ConfigParser class.

This can parse Python-compatible ConfigParser config files, including support for option interpolation.

## Setup
```Go
  import (
    "github.com/bigkevmcd/configparser"
```

## Parsing configuration files
It's easy to parse a configuration file.
```Go
  p, err = configparser.NewConfigParserFromFile("example.cfg")
  if err != nil {
    ...
  }
```

## Methods
The ConfigParser implements most of the Python ConfigParser API
```Go
  v, err := p.Get("section", "option")
  err = p.Set("section", "newoption", "value")

  s := p.Sections()
```

## Interpolation
The ConfigParser implements interpolation in the same format as the Python implementation.

Given the configuration

```
  [DEFAULTS]
  dir: testing

  [testing]
  something: %(dir)s/whatever
```

```Go
  v, err := p.GetInterpolated("testing, something")
```

It's also possible to override the values to use when interpolating values by providing a Dict to lookup values in.
```
  d := make(configparser.Dict)
  d["dir"] = "/a/non/existent/path"
  result, err := p.GetInterpolatedWithVars("testing", "something", d)
```

Will get ```testing/whatever``` as the value
