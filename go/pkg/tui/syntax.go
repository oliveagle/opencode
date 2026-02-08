package tui

import (
	"regexp"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type SyntaxHighlighter struct {
	language string
	keywords map[string]bool
	types    map[string]bool
	builtins map[string]bool
}

func NewSyntaxHighlighter(language string) *SyntaxHighlighter {
	sh := &SyntaxHighlighter{
		language: language,
		keywords: map[string]bool{},
		types:    map[string]bool{},
		builtins: map[string]bool{},
	}

	sh.loadLanguage(language)
	return sh
}

func (sh *SyntaxHighlighter) loadLanguage(lang string) {
	switch strings.ToLower(lang) {
	case "go":
		sh.keywords = map[string]bool{
			"package":  true, "import": true, "const": true, "var": true,
			"func": true, "if": true, "else": true, "for": true, "range": true,
			"switch": true, "case": true, "default": true, "break": true,
			"continue": true, "return": true, "defer": true, "go": true,
			"chan": true, "select": true, "interface": true, "type": true,
			"struct": true,
		}
		sh.types = map[string]bool{
			"int": true, "int8": true, "int16": true, "int32": true, "int64": true,
			"uint": true, "uint8": true, "uint16": true, "uint32": true, "uint64": true,
			"float32": true, "float64": true, "complex64": true, "complex128": true,
			"string": true, "bool": true, "byte": true, "rune": true,
			"error": true, "interface{}": true,
		}
		sh.builtins = map[string]bool{
			"len": true, "cap": true, "make": true, "new": true,
			"append": true, "copy": true, "delete": true, "panic": true,
			"recover": true, "print": true, "println": true, "complex": true,
			"real": true, "imag": true,
		}

	case "python":
		sh.keywords = map[string]bool{
			"def": true, "class": true, "import": true, "from": true, "as": true,
			"if": true, "elif": true, "else": true, "for": true, "while": true,
			"break": true, "continue": true, "return": true, "yield": true,
			"try": true, "except": true, "finally": true, "with": true, "pass": true,
			"lambda": true, "and": true, "or": true, "not": true, "in": true, "is": true,
		}
		sh.types = map[string]bool{
			"int": true, "float": true, "str": true, "bool": true,
			"list": true, "dict": true, "set": true, "tuple": true,
			"None": true, "True": true, "False": true,
		}
		sh.builtins = map[string]bool{
			"print": true, "len": true, "range": true, "enumerate": true,
			"zip": true, "map": true, "filter": true, "sorted": true,
		}

	case "javascript", "typescript", "js", "ts":
		sh.keywords = map[string]bool{
			"function": true, "const": true, "let": true, "var": true,
			"if": true, "else": true, "for": true, "while": true, "do": true,
			"switch": true, "case": true, "break": true, "continue": true,
			"return": true, "try": true, "catch": true, "finally": true,
			"throw": true, "async": true, "await": true, "class": true,
			"extends": true, "super": true, "this": true, "new": true,
			"import": true, "export": true, "default": true, "from": true,
		}
		sh.types = map[string]bool{
			"string": true, "number": true, "boolean": true, "object": true,
			"undefined": true, "null": true, "void": true,
		}
		sh.builtins = map[string]bool{
			"console": true, "log": true, "Array": true, "Object": true,
			"String": true, "Number": true, "Boolean": true, "Date": true,
			"Math": true, "JSON": true, "Promise": true, "Symbol": true,
		}

	case "rust":
		sh.keywords = map[string]bool{
			"fn": true, "let": true, "mut": true, "const": true, "static": true,
			"struct": true, "enum": true, "trait": true, "impl": true, "pub": true,
			"use": true, "mod": true, "crate": true, "self": true, "super": true,
			"if": true, "else": true, "match": true, "for": true, "while": true,
			"loop": true, "break": true, "continue": true, "return": true,
			"unsafe": true, "async": true, "await": true, "move": true,
		}
		sh.types = map[string]bool{
			"i8": true, "i16": true, "i32": true, "i64": true, "i128": true,
			"u8": true, "u16": true, "u32": true, "u64": true, "u128": true,
			"f32": true, "f64": true, "bool": true, "char": true, "str": true,
			"String": true, "Vec": true, "Option": true, "Result": true,
		}
		sh.builtins = map[string]bool{
			"println": true, "print": true, "panic": true, "unwrap": true,
			"expect": true, "Some": true, "None": true, "Ok": true, "Err": true,
		}

	case "yaml", "yml":
		// YAML is simpler, focusing on structure
		sh.keywords = map[string]bool{}
		sh.types = map[string]bool{
			"true": true, "false": true, "null": true, "~": true,
		}
		sh.builtins = map[string]bool{}

	case "json":
		sh.keywords = map[string]bool{}
		sh.types = map[string]bool{
			"true": true, "false": true, "null": true,
		}
		sh.builtins = map[string]bool{}
	}
}

func (sh *SyntaxHighlighter) ColorizeText(text string) map[int]tcell.Color {
	colors := map[int]tcell.Color{}

	// String highlighting (handles ", ', and `)
	stringPatterns := []string{
		`"[^"]*"`,
		`'[^']*'`,
		"`[^`]*`",
	}

	for _, pattern := range stringPatterns {
		stringRegex := regexp.MustCompile(pattern)
		for _, match := range stringRegex.FindAllStringIndex(text, -1) {
			for i := match[0]; i < match[1]; i++ {
				colors[i] = tcell.ColorGreen
			}
		}
	}

	// Comment highlighting
	commentRegex := regexp.MustCompile(`//.*|#.*`)
	for _, match := range commentRegex.FindAllStringIndex(text, -1) {
		for i := match[0]; i < match[1]; i++ {
			colors[i] = tcell.ColorGray
		}
	}

	// Number highlighting
	numberRegex := regexp.MustCompile(`\b\d+(\.\d+)?\b`)
	for _, match := range numberRegex.FindAllStringIndex(text, -1) {
		for i := match[0]; i < match[1]; i++ {
			colors[i] = tcell.ColorYellow
		}
	}

	// Keyword highlighting
	words := regexp.MustCompile(`\b\w+\b`)
	for _, match := range words.FindAllStringIndex(text, -1) {
		word := text[match[0]:match[1]]

		color := tcell.ColorDefault
		if sh.keywords[word] {
			color = tcell.ColorDarkMagenta
		} else if sh.types[word] {
			color = tcell.ColorDarkCyan
		} else if sh.builtins[word] {
			color = tcell.ColorDarkBlue
		}

		if color != tcell.ColorDefault {
			for i := match[0]; i < match[1]; i++ {
				if colors[i] == tcell.ColorDefault {
					colors[i] = color
				}
			}
		}
	}

	return colors
}

func (sh *SyntaxHighlighter) GetLanguage() string {
	return sh.language
}
