package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestExtractTypeAliasThenRename(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	content := `
var x: /*a*/{ a?: number, b?: string }/*b*/ = { };
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title:        "Extract to type alias",
		ApplyChanges: true,
		NewFileContent: `
type NewType = {
    a?: number; b?: string;
};

var x: NewType = { };
`,
	})

	f.GoToMarker(t, "a")
	f.RenameAtCaret(t, "Options")
	f.VerifyCurrentFileContent(t, `
type Options = {
    a?: number; b?: string;
};

var x: Options = { };
`)
}

func TestExtractInterfaceThenRename(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	content := `
var x: /*a*/{ a?: number, b?: string }/*b*/ = { };
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title:        "Extract to interface",
		ApplyChanges: true,
		NewFileContent: `
interface NewType {
    a?: number; b?: string;
}

var x: NewType = { };
`,
	})

	f.GoToMarker(t, "a")
	f.RenameAtCaret(t, "Options")
	f.VerifyCurrentFileContent(t, `
interface Options {
    a?: number; b?: string;
}

var x: Options = { };
`)
}
