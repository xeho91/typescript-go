package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestRefactorInferReturnTypeOnlyExactKind(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
function /*a*/simple/*b*/() {
    return 42;
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorWithOnlyAvailable(t, inferReturnTypeRefactorTitle, []lsproto.CodeActionKind{"refactor.rewrite.function.returnType"})
}

func TestRefactorInferReturnTypeOnlyRewriteKind(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
function /*a*/simple/*b*/() {
    return 42;
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorWithOnlyAvailable(t, inferReturnTypeRefactorTitle, []lsproto.CodeActionKind{"refactor.rewrite"})
}

func TestRefactorInferReturnTypeOnlyEmptyKind(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
function /*a*/simple/*b*/() {
    return 42;
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorWithOnlyAvailable(t, inferReturnTypeRefactorTitle, []lsproto.CodeActionKind{""})
}

func TestRefactorInferReturnTypeOnlyExtractKind(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
function /*a*/simple/*b*/() {
    return 42;
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorWithOnlyNotAvailable(t, inferReturnTypeRefactorTitle, []lsproto.CodeActionKind{"refactor.extract"})
}
