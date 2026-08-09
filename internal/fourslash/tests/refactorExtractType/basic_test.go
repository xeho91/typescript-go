package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestExtractObjectTypeLiteral(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	content := `
var x: /*a*/{ a?: number, b?: string }/*b*/ = { };
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType = {
    a?: number; b?: string;
};

var x: NewType = { };
`,
	})
}

func TestExtractKeywordType(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	content := `
var x: /*a*/string/*b*/ = '';
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType = string;

var x: NewType = '';
`,
	})
}

func TestExtractUnionType(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	content := `
var x: /*a*/string | number | boolean/*b*/ = '';
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType = string | number | boolean;

var x: NewType = '';
`,
	})
}

func TestExtractLiteralType(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	content := `
var x: /*a*/1/*b*/ = 1;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType = 1;

var x: NewType = 1;
`,
	})
}

func TestExtractLiteralUnionType(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	content := `
var x: /*a*/1 | 2/*b*/ = 1;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType = 1 | 2;

var x: NewType = 1;
`,
	})
}

func TestExtractSingleTypeFromUnion(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	content := `
var x: 1 | /*a*/2/*b*/ = 1;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType = 2;

var x: 1 | NewType = 1;
`,
	})
}
