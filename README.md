go-configparser
===============

Go implementation of the Python ConfigParser class.

This can parse Python-compatible ConfigParser config files, including support for option interpolation.

## Setup
<pre>
  import (
    "github.com/bigkevmcd/configparser"
</pre>

## Parsing configuration files
It's easy to parse a configuration file.
<pre>
  p, err = configparser.NewConfigParserFromFile("example.cfg")
  if err != nil {
    ...
  }
</pre>

## Methods
The ConfigParser implements most of the Python ConfigParser API
<pre>
  v, err := p.Get("section", "option")
  err = p.Set("section", "newoption", "value")

  s := p.Sections()
</pre>

## Interpolation
The ConfigParser implements interpolation in the same format as the Python implementation.

Given the configuration

<pre>
  [DEFAULTS]
  dir: testing

  [testing]
  something: %(dir)s/whatever
</pre>

<pre>
  v, err := p.GetInterpolated("testing, something")
</pre>

It's also possible to override the values to use when interpolating values by providing a Dict to lookup values in.
<pre>
  d := make(configparser.Dict)
  d["dir"] = "/a/non/existent/path"
  result, err := p.GetInterpolatedWithVars("testing", "something", d)
</pre>

Will get <code>testing/whatever</code> as the value