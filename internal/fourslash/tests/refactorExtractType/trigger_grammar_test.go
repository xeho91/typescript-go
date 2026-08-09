package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestTrigger_emptySpanInvokedOnly(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
var x: str/*a*//*b*/ing;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailableForTriggerReason(t, "implicit", "Extract to type alias")
	f.VerifyRefactorAvailableForTriggerReason(t, "invoked", "Extract to type alias")
}

func TestTrigger_nonEmptySpanAvailableForBoth(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
var x: s/*a*/tr/*b*/ing;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorAvailableForTriggerReason(t, "implicit", "Extract to type alias")
	f.VerifyRefactorAvailableForTriggerReason(t, "invoked", "Extract to type alias")
}

func TestGrammarError_extractType(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type Foo = /*a*/{ x: string = a }/*b*/
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType = {
    x: string;
};

type Foo = NewType
`,
	})
}

func TestGrammarError_extractConstraint(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
type Foo<T extends /*a*/{ x: string = a }/*b*/> = T
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType = {
    x: string;
};

type Foo<T extends NewType> = T
`,
	})
}
