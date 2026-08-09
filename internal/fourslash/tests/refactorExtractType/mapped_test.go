package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestMapped_type(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type Deep<T> = /*a*/{ [K in keyof T]: Deep<T[K]> }/*b*/
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType<T> = {
    [K in keyof T]: Deep<T[K]>;
};

type Deep<T> = NewType<T>
`,
	})
}

func TestMapped_inConditional(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type Expand<T> = T extends any
    ? /*a*/{ [K in keyof T]: Expand<T[K]> }/*b*/
    : never;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType<T> = {
    [K in keyof T]: Expand<T[K]>;
};

type Expand<T> = T extends any
    ? NewType<T>
    : never;
`,
	})
}
