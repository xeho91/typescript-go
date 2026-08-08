package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestRefactorInferReturnTypeNotInFunction(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
let x = 10;/*a*//*b*/

function func() {
    return 10;
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeRefactorTitle)
}

func TestRefactorInferReturnTypeComputedAsConst(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
function /*a*//*b*/f() {
    return {
        [-1]: 0
    } as const
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: inferReturnTypeRefactorTitle,
		NewFileContent: `
function f(): {
    readonly [-1]: 0;
} {
    return {
        [-1]: 0
    } as const
}`,
	})
}
