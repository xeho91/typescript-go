package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestQualified_generic(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type A<T> = /*a*/B.C.D<T>/*b*/;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType<T> = B.C.D<T>;

type A<T> = NewType<T>;
`,
	})
}

func TestQualified_simple(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type A = /*a*/B.C.D/*b*/;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType = B.C.D;

type A = NewType;
`,
	})
}
