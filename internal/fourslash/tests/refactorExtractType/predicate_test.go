package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestPredicate_extractParam(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type A = (v: /*a*/string | number/*b*/) => v is string;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType = string | number;

type A = (v: NewType) => v is string;
`,
	})
}

func TestPredicate_extractReturn(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type A = (v: string | number) => v is /*a*/string/*b*/;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType = string;

type A = (v: string | number) => v is NewType;
`,
	})
}

func TestPredicate_notAvailFull(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type A = (v: string | number) => /*a*/v is string/*b*/
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, "Extract to type alias")
	f.VerifyRefactorNotAvailable(t, "Extract to interface")
}

func TestPredicate_extractWholeFunc(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type A = /*a*/(v: string | number) => v is string/*b*/;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType = (v: string | number) => v is string;

type A = NewType;
`,
	})
}

func TestTypeof_notAvailLocal(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type A = (v: string | number) => /*a*/typeof v/*b*/
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, "Extract to type alias")
	f.VerifyRefactorNotAvailable(t, "Extract to interface")
}

func TestTypeof_extractWholeFuncWithLocal(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type A = /*a*/(v: string | number) => typeof v/*b*/
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType = (v: string | number) => typeof v;

type A = NewType
`,
	})
}

func TestTypeof_notAvailUnionWithLocal(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type A = (v: string | number) => /*a*/number | typeof v/*b*/
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, "Extract to type alias")
	f.VerifyRefactorNotAvailable(t, "Extract to interface")
}

func TestTypeof_extractModuleConst(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
const a = 1;
type A = (v: string | number) => /*a*/typeof a/*b*/;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
const a = 1;
type NewType = typeof a;

type A = (v: string | number) => NewType;
`,
	})
}

func TestTypeof_extractNamespaceConst(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
namespace A { export const b = 1; }
function a(b: string): /*a*/typeof A.b/*b*/ { return 1; }
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
namespace A { export const b = 1; }
type NewType = typeof A.b;

function a(b: string): NewType { return 1; }
`,
	})
}

func TestTypeof_notAvailThis(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
interface I { f: (this: O, b: number) => /*a*/typeof this.a/*b*/ };
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, "Extract to type alias")
	f.VerifyRefactorNotAvailable(t, "Extract to interface")
}

func TestThis_notAvailThisRef(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
interface I { f: (this: O, b: number) => /*a*/this/*b*/ };
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, "Extract to type alias")
	f.VerifyRefactorNotAvailable(t, "Extract to interface")
}

func TestThis_notAvailIndexedThis(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
interface I { f: (this: O, b: number) => /*a*/typeof this["a"]/*b*/ };
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, "Extract to type alias")
	f.VerifyRefactorNotAvailable(t, "Extract to interface")
}

func TestThis_extractFullFuncWithThis(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
interface I { f: /*a*/(this: O, b: number) => typeof this.a/*b*/ };
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType = (this: O, b: number) => typeof this.a;

interface I { f: NewType };
`,
	})
}

func TestThis_extractFullFuncWithThisRef(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
interface I { f: /*a*/(this: O, b: number) => this/*b*/ };
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType = (this: O, b: number) => this;

interface I { f: NewType };
`,
	})
}

func TestThis_extractFullFuncWithIndexedThis(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
interface I { f: /*a*/(this: O, b: number) => typeof this["a"]/*b*/ };
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType = (this: O, b: number) => typeof this["a"];

interface I { f: NewType };
`,
	})
}

func TestTypeof_notAvailOtherParam(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
function f(a: string, b: /*a*/typeof a/*b*/): typeof b { return ''; }
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, "Extract to type alias")
	f.VerifyRefactorNotAvailable(t, "Extract to interface")
}

func TestTypeof_extractInBlock(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
function id<T>(x: T): T {
    return (() => {
        const s: /*a*/typeof x/*b*/ = x;
        return s;
    })();
}
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
function id<T>(x: T): T {
    return (() => {
        type NewType = typeof x;

        const s: NewType = x;
        return s;
    })();
}
`,
	})
}

func TestTypeof_templateLiteral(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
const key = "key";
type Foo = /*a*/` + "`" + `${typeof key}Foo` + "`" + `/*b*/;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
const key = "key";
type NewType = ` + "`" + `${typeof key}Foo` + "`" + `;

type Foo = NewType;
`,
	})
}

func TestThis_notAvailMixedWithThis(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
interface I { f: (this: O, b: number) => /*a*/ true | this | false /*b*/ };
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, "Extract to type alias")
	f.VerifyRefactorNotAvailable(t, "Extract to interface")
}

func TestThis_notAvailClassThis(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
class C {
    m<T>(): /*a*/T | this | number/*b*/ {
        return {} as any
    }
}
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, "Extract to type alias")
	f.VerifyRefactorNotAvailable(t, "Extract to interface")
}
