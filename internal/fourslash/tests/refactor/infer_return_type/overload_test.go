package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestRefactorInferReturnTypeOverloadImplString(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
function f(x: number) {
    return 1;
}
function /*a*/f/*b*/(x: number) {
    return "1";
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: inferReturnTypeRefactorTitle,
		NewFileContent: `
function f(x: number) {
    return 1;
}
function f(x: number): string {
    return "1";
}`,
	})
}

func TestRefactorInferReturnTypeOverloadSignatureOnlyNotAvailable(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
function /*a*/f/*b*/(x: string): number;
function f(x: string | number) {
  return 1;
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeRefactorTitle)
}

func TestRefactorInferReturnTypeOverloadUnion(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
function f(x: number): string;
function f(x: string): number;
function /*a*/f/*b*/(x: string | number) {
    return x === 1 ? 1 : "quit";
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: inferReturnTypeRefactorTitle,
		NewFileContent: `
function f(x: number): string;
function f(x: string): number;
function f(x: string | number): string | number {
    return x === 1 ? 1 : "quit";
}`,
	})
}

func TestRefactorInferReturnTypeOverloadNamedUnion(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
interface Foo {
    x: number;
}
function f(x: number): Foo;
function f(x: string): number;
function /*a*/f/*b*/(x: string | number) {
    return x === 1 ? 1 : { x };
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: inferReturnTypeRefactorTitle,
		NewFileContent: `
interface Foo {
    x: number;
}
function f(x: number): Foo;
function f(x: string): number;
function f(x: string | number): number | Foo {
    return x === 1 ? 1 : { x };
}`,
	})
}
