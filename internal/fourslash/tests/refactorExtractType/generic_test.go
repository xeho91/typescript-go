package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestGeneric_extractTypeParam(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type A<T = boolean> = string | number | /*a*/T/*b*/;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType<T> = T;

type A<T = boolean> = string | number | NewType<T>;
`,
	})
}

func TestGeneric_partial(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type A<B, C, D = B> = /*a*/Partial<C | string>/*b*/ & D | C;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType<C> = Partial<C | string>;

type A<B, C, D = B> = NewType<C> & D | C;
`,
	})
}

func TestGeneric_partialMultiParam(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type A<B, C, D = B> = /*a*/Partial<C | string | D>/*b*/ & D | C;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType<C, D> = Partial<C | string | D>;

type A<B, C, D = B> = NewType<C, D> & D | C;
`,
	})
}

func TestGeneric_nestedShadowed_firstNestedT(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type A<T, U> = () => <T>(v: /*a*/T/*b*/) => (v: T) => <T>(v: T) => U;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType<T> = T;

type A<T, U> = () => <T>(v: NewType<T>) => (v: T) => <T>(v: T) => U;
`,
	})
}

func TestGeneric_nestedShadowed_secondNestedT(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type A<T, U> = () => <T>(v: T) => (v: /*a*/T/*b*/) => <T>(v: T) => U;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType<T> = T;

type A<T, U> = () => <T>(v: T) => (v: NewType<T>) => <T>(v: T) => U;
`,
	})
}

func TestGeneric_nestedShadowed_thirdNestedT(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type A<T, U> = () => <T>(v: T) => (v: T) => <T>(v: /*a*/T/*b*/) => U;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType<T> = T;

type A<T, U> = () => <T>(v: T) => (v: T) => <T>(v: NewType<T>) => U;
`,
	})
}

func TestGeneric_nestedShadowed_outerU(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type A<T, U> = () => <T>(v: T) => (v: T) => <T>(v: T) => /*a*/U/*b*/;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType<U> = U;

type A<T, U> = () => <T>(v: T) => (v: T) => <T>(v: T) => NewType<U>;
`,
	})
}

func TestGeneric_unionWithGeneric(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type A<T,S> = /*1*/{ a: string } | { b: T } | { c: string }/*2*/ | { d: string } | S;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "1", "2")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType<T> = {
    a: string;
} | {
    b: T;
} | {
    c: string;
};

type A<T,S> = NewType<T> | { d: string } | S;
`,
	})
}

func TestGeneric_constrained(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type A<T extends string> = T;
type B<T extends string> = /*a*/A<T>/*b*/;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type A<T extends string> = T;
type NewType<T extends string> = A<T>;

type B<T extends string> = NewType<T>;
`,
	})
}

func TestGeneric_templateLiteral(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type Foo<T> = {
    fn:
        keyof T extends string
            ? /*a*/(x: ` + "`" + `${keyof T}Foo` + "`" + `, callback: (y: keyof T) => void) => void/*b*/
            : never;
}
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType<T> = (x: ` + "`" + `${keyof T}Foo` + "`" + `, callback: (y: keyof T) => void) => void;

type Foo<T> = {
    fn:
        keyof T extends string
            ? NewType<T>
            : never;
}
`,
	})
}

func TestGeneric_multiParam(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type Foo<T1, T2, T3> = {
    fn: /*a*/(a: T1, b: T2, c: T3, a1: T1, a2: T2, a3: T3) => void;/*b*/
}
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType<T1, T2, T3> = (a: T1, b: T2, c: T3, a1: T1, a2: T2, a3: T3) => void;

type Foo<T1, T2, T3> = {
    fn: NewType<T1, T2, T3>;
}
`,
	})
}
