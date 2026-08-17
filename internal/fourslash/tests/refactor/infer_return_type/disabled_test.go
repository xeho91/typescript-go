package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestRefactorInferReturnTypeDisabledWithCapability(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
function /*a*/isString/*b*/(x: any): x is string {
    return typeof x === "string";
}`
	ptrTrue := true
	capabilities := &lsproto.ClientCapabilities{
		TextDocument: &lsproto.TextDocumentClientCapabilities{
			CodeAction: &lsproto.CodeActionClientCapabilities{
				DisabledSupport: &ptrTrue,
			},
		},
	}
	f, done := fourslash.NewFourslash(t, capabilities, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorDisabled(t, inferReturnTypeRefactorTitle)
}

func TestRefactorInferReturnTypeDisabledWithoutCapability(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
function /*a*/isString/*b*/(x: any): x is string {
    return typeof x === "string";
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeRefactorTitle)
}
