# loon 

A line-oriented object notation parser

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/iand/loon)
[![Check Status](https://github.com/iand/loon/actions/workflows/check.yml/badge.svg)](https://github.com/iand/loon/actions/workflows/check.yml)
[![Test Status](https://github.com/iand/loon/actions/workflows/test.yml/badge.svg)](https://github.com/iand/loon/actions/workflows/test.yml)

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
 - Lines starting with `#` are comments. 
 - The first line of a loon document may optionally include a syntax version tag written as `@version <num>`. Only `@version 1` is supported currently. Documents without a version tag are assumed to be version 1. 
 - Comments are preserved when parsing and will round trip when printed:
    - Comments before an object or directive are assumed to belong to the object or directive.
    - Comments before the first object in a document followed by an empty line are assumed to belong to the document. 
    - Comments after the last directive are associated with the object as trailing comments.
    - Comments after the last object are associated with the document as trailing comments.

An example:

```
# This is a comment about the document
# Since there is an empty line between it and
# the object comment below

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

# a trailing document comment
```

## Usage

```Go
import (
	"os"
	"log"

	"github.com/iand/loon"
)

func main() {
	f, err := os.Open("rules.loon")
	if err != nil {
		log.Fatalf("open: %v", err)
	}
	defer f.Close()

	p := NewParser(f)
	doc, err := p.Parse()
	if err != nil {
		log.Fatalf("parse: %v", err)
	}

	for _, o := range doc.Objects {
		log.Printf("object type=%s, name=%s\n", o.Type, o.Name)
	}
}
```

## Author

* [Ian Davis](http://github.com/iand) - <http://iandavis.com/>

# License

This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying [`UNLICENSE`](UNLICENSE) file.

