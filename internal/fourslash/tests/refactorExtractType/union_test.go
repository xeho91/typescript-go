package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestUnion_middleSubset(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type A = { a: string } | /*1*/{ b: string } | { c: string }/*2*/ | { d: string };
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "1", "2")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType = {
    b: string;
} | {
    c: string;
};

type A = { a: string } | NewType | { d: string };
`,
	})
}

func TestUnion_typeRefs(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type B = string;
type C = number;
type A = { a: string } | /*1*/B | C/*2*/;
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

type A = { a: string } | NewType;
`,
	})
}

func TestUnion_selectionInKeyword(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type A = { a: str/*1*/ing } | { b: string } | { c: string }/*2*/;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "1", "2")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType = {
    a: string;
} | {
    b: string;
} | {
    c: string;
};

type A = NewType;
`,
	})
}

func TestUnion_selectionAcrossBraces(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type A = { a: string /*1*/} | { b: string } | { c: string }/*2*/;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "1", "2")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType = {
    a: string;
} | {
    b: string;
} | {
    c: string;
};

type A = NewType;
`,
	})
}

func TestUnion_selectionBeforeOperator(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type A = { a: string } /*1*/| { b: string } | { c: string }/*2*/;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "1", "2")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType = {
    a: string;
} | {
    b: string;
} | {
    c: string;
};

type A = NewType;
`,
	})
}

func TestUnion_selectionInMemberBraces(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type A = {/*1*/ a: string } | { b: string } | { /*2*/c: string };
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "1", "2")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType = {
    a: string;
} | {
    b: string;
} | {
    c: string;
};

type A = NewType;
`,
	})
}

func TestUnion_selectionAtOperator(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type A = /*1*/{ a: string } | { b: string } |/*2*/ { c: string };
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "1", "2")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType = {
    a: string;
} | {
    b: string;
};

type A = NewType | { c: string };
`,
	})
}
