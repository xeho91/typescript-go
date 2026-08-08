package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestRefactorInferReturnTypeAsyncArrowParenLess(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
const foo = async /*a*//*b*/a => {
    return 1;
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: inferReturnTypeRefactorTitle,
		NewFileContent: `
const foo = async (a): Promise<number> => {
    return 1;
}`,
	})
}

func TestRefactorInferReturnTypeAsyncArrow(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
const foo = async /*a*//*b*/(a) => {
    return 1;
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: inferReturnTypeRefactorTitle,
		NewFileContent: `
const foo = async (a): Promise<number> => {
    return 1;
}`,
	})
}

func TestRefactorInferReturnTypeAsyncFunction(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
async function /*a*/asyncFunc/*b*/() {
    return 42;
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: inferReturnTypeRefactorTitle,
		NewFileContent: `
async function asyncFunc(): Promise<number> {
    return 42;
}`,
	})
}
