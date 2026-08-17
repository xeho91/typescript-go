package ls

import (
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/parser"
	"github.com/microsoft/typescript-go/internal/stringutil"
)

func parseInferReturnTypeSource(t *testing.T, text string) *ast.SourceFile {
	t.Helper()
	return parser.ParseSourceFile(ast.SourceFileParseOptions{
		FileName: "/index.ts",
		Path:     "/index.ts",
	}, text, core.ScriptKindTS)
}

func inferReturnTypeTokenAt(text string, sourceFile *ast.SourceFile, marker string) (*ast.Node, int) {
	index := strings.Index(text, marker)
	if index < 0 {
		return nil, -1
	}

	pos := index + len(marker)
	for pos < len(text) && stringutil.IsWhiteSpaceLike(rune(text[pos])) {
		pos++
	}

	return astnav.GetTouchingPropertyName(sourceFile, pos), pos
}

func TestFindConvertibleAncestor(t *testing.T) {
	t.Parallel()

	text := `
// @functionDeclaration
function f(a: number) {
  // @insideFunctionBlock
  return a;
}

const arrow = (
  // @arrowParameter
  a: number,
): number => a;

const arrowBody = (a: number) =>
  // @insideArrowBody
  a;

class C {
  // @constructor
  constructor(public x: number) {}

  // @getter
  get y() { return 1; }

  // @setter
  set z(v: number) {}

  // @method
  method(a: number) { return a; }
}

// @overloadSignature
function g(a: number);

// @overloadImplementation
function h(a: number) { return a; }
`
	sourceFile := parseInferReturnTypeSource(t, text)

	tests := []struct {
		name   string
		marker string
		want   ast.Kind // ast.Kind(0) means no convertible ancestor
	}{
		{"function declaration", "// @functionDeclaration", ast.KindFunctionDeclaration},
		{"inside function block", "// @insideFunctionBlock", 0},
		{"arrow parameter", "// @arrowParameter", ast.KindArrowFunction},
		{"arrow body expression", "// @insideArrowBody", 0},
		{"constructor", "// @constructor", 0},
		{"getter", "// @getter", 0},
		{"setter", "// @setter", 0},
		{"method", "// @method", ast.KindMethodDeclaration},
		{"overload signature", "// @overloadSignature", ast.KindFunctionDeclaration},
		{"overload implementation", "// @overloadImplementation", ast.KindFunctionDeclaration},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			token, pos := inferReturnTypeTokenAt(text, sourceFile, tt.marker)
			if pos < 0 {
				t.Fatalf("marker %q not found in test source", tt.marker)
			}

			got := findConvertibleAncestor(token)
			if tt.want == 0 {
				if got != nil {
					t.Errorf("expected no convertible ancestor, got %v", got.Kind)
				}
				return
			}
			if got == nil || got.Kind != tt.want {
				t.Errorf("expected convertible ancestor %v, got %v", tt.want, nodeKindOrNil(got))
			}
		})
	}
}

func TestHasBody(t *testing.T) {
	t.Parallel()

	text := `
// @functionDeclaration
function f(a: number) {
  return a;
}

const arrow = (
  // @arrowParameter
  a: number,
): number => a;

class C {
  // @method
  method(a: number) { return a; }
}

// @signatureOnly
function g(a: number);

// @overloadSignature
function h(a: number);

// @overloadImplementation
function h(a: number) { return a; }
`
	sourceFile := parseInferReturnTypeSource(t, text)

	tests := []struct {
		name   string
		marker string
		want   bool
	}{
		{"function declaration", "// @functionDeclaration", true},
		{"arrow function", "// @arrowParameter", true},
		{"method", "// @method", true},
		{"signature-only declaration", "// @signatureOnly", false},
		{"overload signature", "// @overloadSignature", false},
		{"overload implementation", "// @overloadImplementation", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			token, pos := inferReturnTypeTokenAt(text, sourceFile, tt.marker)
			if pos < 0 {
				t.Fatalf("marker %q not found in test source", tt.marker)
			}

			declaration := findFunctionLikeAncestor(token)
			if declaration == nil {
				t.Fatal("expected a function-like ancestor")
			}
			if got := hasBody(declaration); got != tt.want {
				t.Errorf("expected hasBody %v, got %v", tt.want, got)
			}
		})
	}
}

func findFunctionLikeAncestor(node *ast.Node) *ast.Node {
	for node != nil {
		if ast.IsFunctionLikeDeclaration(node) {
			return node
		}
		node = node.Parent
	}
	return nil
}

func nodeKindOrNil(node *ast.Node) ast.Kind {
	if node == nil {
		return 0
	}
	return node.Kind
}
