package cfg

import (
	"github.com/simplesurance/baur/v1/cfg/resolver"
)

// GolangSources specifies inputs for Golang Applications
type GolangSources struct {
	Queries     []string `toml:"queries" comment:"Queries specify the source files or packages of which the dependencies are resolved.\n Format:\n \tfile=<RELATIVE-PATH>\n \tfileglob=<GLOB-PATTERN>\t -> Supports double-star\n \tEverything else is passed directly to underlying build tool (go list by default).\n \tSee also the patterns described at:\n \t<https://github.com/golang/tools/blob/bc8aaaa29e0665201b38fa5cb5d47826788fa249/go/packages/doc.go#L17>.\n Files from Golang's stdlib are ignored.\n Valid variables: $ROOT, $APPNAME."`
	Environment []string `toml:"environment" comment:"Environment to use when discovering Golang source files\n This are environment variables understood by the Golang tools, like GOPATH, GOFLAGS, etc.\n If empty the default Go environment is used.\n Valid variables: $ROOT, $APPNAME"`
	Tests       bool     `toml:"tests" comment:"If true queries are resolved to test files, otherwise testfiles are ignored."`
}

func (g *GolangSources) IsEmpty() bool {
	return len(g.Environment) == 0 && len(g.Queries) == 0 && !g.Tests
}

// Merge merges the two GolangSources structs
func (g *GolangSources) Merge(other *GolangSources) {
	// TODO: merging this section is currently buggy,
	// https://github.com/simplesurance/baur/issues/169 must be fixed

	g.Queries = append(g.Queries, other.Queries...)
	g.Environment = append(g.Environment, other.Environment...)

	if other.Tests {
		g.Tests = other.Tests
	}
}

func (g *GolangSources) Resolve(resolvers resolver.Resolver) error {
	for i, env := range g.Environment {
		var err error

		if g.Environment[i], err = resolvers.Resolve(env); err != nil {
			return FieldErrorWrap(err, "Environment", env)
		}
	}

	for i, q := range g.Queries {
		var err error

		if g.Queries[i], err = resolvers.Resolve(q); err != nil {
			return FieldErrorWrap(err, "Paths", q)
		}
	}

	return nil
}

// Validate checks that the stored information is valid.
func (g *GolangSources) Validate() error {
	if (len(g.Environment) != 0 || g.Tests) && len(g.Queries) == 0 {
		return NewFieldError("must be set if environment or tests is set", "query")
	}

	for _, q := range g.Queries {
		if len(q) == 0 {
			return NewFieldError("empty string is an invalid query", "query")
		}
	}

	return nil
}
