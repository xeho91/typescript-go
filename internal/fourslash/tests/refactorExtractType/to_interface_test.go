package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestToInterface_basic(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
function foo(a: /*a*/{ a: number | string, b: string }/*b*/) { }
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to interface",
		NewFileContent: `
interface NewType {
    a: number | string; b: string;
}

function foo(a: NewType) { }
`,
	})
}

func TestToInterface_intersection(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
function foo(a: /*a*/{ a: number | string, b: string } & { c: string } & { d: boolean }/*b*/) { }
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to interface",
		NewFileContent: `
interface NewType {
    a: number | string; b: string; c: string; d: boolean;
}

function foo(a: NewType) { }
`,
	})
}

func TestToInterface_notAvailWithTypeAlias(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type T = { c: string }
function foo(a: /*a*/{ a: number | string, b: string } & T/*b*/) { }
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, "Extract to interface")
}

func TestToInterface_notAvailWithComplexType(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type T = { c: string } & { d: boolean }
function foo(a: /*a*/{ a: number | string, b: string } & T/*b*/) { }
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, "Extract to interface")
}

func TestToInterface_notAvailWithRecord(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type T = { c: string } & Record<string, string>
function foo(a: /*a*/{ a: number | string, b: string } & T/*b*/) { }
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, "Extract to interface")
}

func TestToInterface_notAvailTypeRefOnly(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type T = { c: string }
function foo(a: /*a*/T/*b*/) { }
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, "Extract to interface")
}

func TestToInterface_generic(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
function foo<U>(a: /*a*/{ a: string } & { b: U }/*b*/) { }
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to interface",
		NewFileContent: `
interface NewType<U> {
    a: string; b: U;
}

function foo<U>(a: NewType<U>) { }
`,
	})
}

func TestToInterface_notAvailConflictingProps(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
function foo(a: /*a*/{ a: number | string, b: string } & { b: string } & { d: boolean }/*b*/) { }
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, "Extract to interface")
}

func TestToInterface_notAvailConflictingTypes(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
function foo(a: /*a*/{ a: number | string, b: string } & { b: number } & { d: boolean }/*b*/) { }
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, "Extract to interface")
}

func TestToInterface_intersectionSubset(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type A = { a: string } & /*1*/{ b: string } & { c: string }/*2*/;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "1", "2")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to interface",
		NewFileContent: `
interface NewType {
    b: string; c: string;
}

type A = { a: string } & NewType;
`,
	})
}
