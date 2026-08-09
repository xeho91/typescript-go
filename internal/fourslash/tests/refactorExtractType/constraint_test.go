package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestConstraint_typeRefs(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type B = string;
type C = number;

export function foo<T extends boolean | /*1*/B | C/*2*/>(x: T): T {
    return x;
}
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "1", "2")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type B = string;
type C = number;

type NewType = B | C;

export function foo<T extends boolean | NewType>(x: T): T {
    return x;
}
`,
	})
}

func TestConstraint_selfReference(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
export declare function foo<T extends { a?: /*a*/T/*b*/ }>(): void;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, "Extract to type alias")
}

func TestConstraint_multiRefExtraction(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type B = { b: string };
type C = { c: number };

interface X<T extends { prop: T | /*3*/B | C/*4*/ }> {}
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "3", "4")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type B = { b: string };
type C = { c: number };

type NewType = B | C;

interface X<T extends { prop: T | NewType }> {}
`,
	})
}

func TestConstraint_selfReferenceInner(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type B = { b: string };
type C = { c: number };

interface X<T extends { prop: /*1*/T | B | C/*2*/ }> {}
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "1", "2")
	f.VerifyRefactorNotAvailable(t, "Extract to type alias")
}
