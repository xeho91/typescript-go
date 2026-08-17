package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestExtractFromParameterType(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
function foo(a: /*a*/number/*b*/, b?: number, ...c: number[]): boolean {
    return false as boolean;
}
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType = number;

function foo(a: NewType, b?: number, ...c: number[]): boolean {
    return false as boolean;
}
`,
	})
}

func TestExtractFromOptionalParameterType(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
function foo(a: number, b?: /*a*/number/*b*/, ...c: number[]): boolean {
    return false as boolean;
}
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType = number;

function foo(a: number, b?: NewType, ...c: number[]): boolean {
    return false as boolean;
}
`,
	})
}

func TestExtractFromRestParameterType(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
function foo(a: number, b?: number, ...c: /*a*/number[]/*b*/): boolean {
    return false as boolean;
}
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType = number[];

function foo(a: number, b?: number, ...c: NewType): boolean {
    return false as boolean;
}
`,
	})
}

func TestExtractFromReturnType(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
function foo(a: number, b?: number, ...c: number[]): /*a*/boolean/*b*/ {
    return false as boolean
}
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType = boolean;

function foo(a: number, b?: number, ...c: number[]): NewType {
    return false as boolean
}
`,
	})
}

func TestExtractFromTypeAssertion(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
function foo(a: number, b?: number, ...c: number[]): boolean {
    return false as /*a*/boolean/*b*/;
}
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
function foo(a: number, b?: number, ...c: number[]): boolean {
    type NewType = boolean;

    return false as NewType;
}
`,
	})
}
