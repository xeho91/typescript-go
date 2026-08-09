package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestObjectLiteral_typeMember(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
var x: { a?: /*a*/number/*b*/, b?: string } = { };
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to type alias",
		NewFileContent: `
type NewType = number;

var x: { a?: NewType, b?: string } = { };
`,
	})
}
