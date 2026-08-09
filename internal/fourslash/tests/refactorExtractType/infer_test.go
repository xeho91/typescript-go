package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestInfer_nestedGenericArrow_hoistOuterU(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type A<T, U> = () => <T>(v: T) => (v: T) => /*a*/<T>(v: T) => U/*b*/;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType<U> = <T>(v: T) => U;

type A<T, U> = () => <T>(v: T) => (v: T) => NewType<U>;
`,
	})
}

func TestInfer_nestedGenericArrow_hoistTU(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type A<T, U> = () => <T>(v: T) => /*a*/(v: T) => <T>(v: T) => U/*b*/;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType<T, U> = (v: T) => <T>(v: T) => U;

type A<T, U> = () => <T>(v: T) => NewType<T, U>;
`,
	})
}

func TestInfer_nestedGenericArrow_wholeNested(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type A<T, U> = () => /*a*/<T>(v: T) => (v: T) => <T>(v: T) => U/*b*/;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType<U> = <T>(v: T) => (v: T) => <T>(v: T) => U;

type A<T, U> = () => NewType<U>;
`,
	})
}

func TestInfer_nestedGenericArrow_wholeAlias(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type A<T, U> = /*a*/() => <T>(v: T) => (v: T) => <T>(v: T) => U/*b*/;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType<U> = () => <T>(v: T) => (v: T) => <T>(v: T) => U;

type A<T, U> = NewType<U>;
`,
	})
}

func TestInfer_extractTConditional(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type Item<T> = /*a*/T/*b*/ extends (infer P)[] ? P : never;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType<T> = T;

type Item<T> = NewType<T> extends (infer P)[] ? P : never;
`,
	})
}

func TestInfer_extractPFromConditional(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type Item<T> = T extends (infer P)[] ? /*a*/P/*b*/ : never;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType<P> = P;

type Item<T> = T extends (infer P)[] ? NewType<P> : never;
`,
	})
}

func TestInfer_extractNeverFromConditional(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type Item<T> = T extends (infer P)[] ? P : /*a*/never/*b*/;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType = never;

type Item<T> = T extends (infer P)[] ? P : NewType;
`,
	})
}

func TestInfer_extractConditionalWhole(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type Item<T> = /*a*/T extends (infer P)[] ? P : never/*b*/;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType<T> = T extends (infer P)[] ? P : never;

type Item<T> = NewType<T>;
`,
	})
}

func TestInfer_extractUnion(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type Union<T, U> = /*a*/U | T/*b*/;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType<U, T> = U | T;

type Union<T, U> = NewType<U, T>;
`,
	})
}

func TestInfer_notAvail_inferArray(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type Item<T> = T extends /*a*/(infer P)[]/*b*/ ? P : never
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, "Extract to type alias")
	f.VerifyRefactorNotAvailable(t, "Extract to interface")
}

func TestInfer_notAvail_inferKeyword(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type Item<T> = T extends (/*a*/infer P/*b*/)[] ? P : never
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, "Extract to type alias")
	f.VerifyRefactorNotAvailable(t, "Extract to interface")
}

func TestInfer_notAvail_inferVariable(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type Item<T> = T extends (infer /*a*/P/*b*/)[] ? P : never
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, "Extract to type alias")
	f.VerifyRefactorNotAvailable(t, "Extract to interface")
}

func TestInfer_notAvail_inferInConditional(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type Crazy<T> = T extends [infer P, (/*a*/infer R extends string ? string : never/*b*/)] ? P & R : string;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, "Extract to type alias")
	f.VerifyRefactorNotAvailable(t, "Extract to interface")
}

func TestInfer_extractConditionalWholeTuple(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type Crazy<T> = /*a*/T extends [infer P, ((infer R) extends string ? string : never)] ? P & R : string/*b*/;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType<T> = T extends [
    infer P, ((infer R) extends string ? string : never)] ? P & R : string;

type Crazy<T> = NewType<T>;
`,
	})
}
