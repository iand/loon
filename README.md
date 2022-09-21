# loon 

A line-oriented object notation parser

[![Check Status](https://github.com/iand/loon/actions/workflows/check.yml/badge.svg)](https://github.com/iand/loon/actions/workflows/check.yml)
[![Test Status](https://github.com/iand/loon/actions/workflows/test.yml/badge.svg)](https://github.com/iand/loon/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/iand/loon)](https://goreportcard.com/report/github.com/iand/loon)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/iand/loon)

## Overview

**loon** is a parser for an ultra simple line-oriented notation for describing objects. It was originally conceived
as a simple resource definition format for games.

An object is declared by writing its type and name on a single line, followed by any number of lines containing directives 
that apply to the object. The object declaration is terminated by the keyword `end` on its own line:

```
<type> <name>
	<directive name> <arg> <arg> <arg>
end
```

Everything between the declaration and the `end` keyword is interpreted as a directive. 
A directive has a name and zero or more arguments which are simply space separated strings.

Notes:

 - All leading and trailing whitespace is ignored
 - Lines starting with `#` are comments and ignored

An example:

```
# Define some resources
resource iron_ore
	name Iron Ore
end

resource iron_ingot
	name Iron Ingot
end

# Define some production rules
rule smelter
	in iron_ore 3
	out iron_ingot 3
end

rule blacksmith
	in iron_ingot 1
	out horseshoe 1
end
```


## Author

* [Ian Davis](http://github.com/iand) - <http://iandavis.com/>

# License

This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying [`UNLICENSE`](UNLICENSE) file.

