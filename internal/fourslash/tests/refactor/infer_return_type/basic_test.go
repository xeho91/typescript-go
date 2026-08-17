package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

const inferReturnTypeRefactorTitle = "Infer function return type"

func TestRefactorInferReturnTypeObjectLiteral(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
function /*a*/foo/*b*/() {
    return { x: 1, y: 1 };
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: inferReturnTypeRefactorTitle,
		NewFileContent: `
function foo(): {
    x: number;
    y: number;
} {
    return { x: 1, y: 1 };
}`,
	})
}

func TestRefactorInferReturnTypeMethod(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
class Foo {
    /*a*/method/*b*/() {
        return { x: 1, y: 1 };
    }
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: inferReturnTypeRefactorTitle,
		NewFileContent: `
class Foo {
    method(): {
        x: number;
        y: number;
    } {
        return { x: 1, y: 1 };
    }
}`,
	})
}

func TestRefactorInferReturnTypeFunctionExpression(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
const foo = /*a*/function/*b*/() {
    return { x: 1, y: 1 };
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: inferReturnTypeRefactorTitle,
		NewFileContent: `
const foo = function(): {
    x: number;
    y: number;
} {
    return { x: 1, y: 1 };
}`,
	})
}

func TestRefactorInferReturnTypeArrowFunction(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
const foo = /*a*/()/*b*/ => {
    return { x: 1, y: 1 };
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: inferReturnTypeRefactorTitle,
		NewFileContent: `
const foo = (): {
    x: number;
    y: number;
} => {
    return { x: 1, y: 1 };
}`,
	})
}

func TestRefactorInferReturnTypeNumber(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
function /*a*/foo/*b*/() {
    return 1;
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: inferReturnTypeRefactorTitle,
		NewFileContent: `
function foo(): number {
    return 1;
}`,
	})
}

func TestRefactorInferReturnTypeString(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
function /*a*/foo/*b*/() {
    return "";
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: inferReturnTypeRefactorTitle,
		NewFileContent: `
function foo(): string {
    return "";
}`,
	})
}

func TestRefactorInferReturnTypeVoid(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
function /*a*/foo/*b*/() {
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: inferReturnTypeRefactorTitle,
		NewFileContent: `
function foo(): void {
}`,
	})
}

func TestRefactorInferReturnTypeMethodExpression(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
const obj = {
    /*a*/myMethod/*b*/() {
        return 42;
    }
};`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: inferReturnTypeRefactorTitle,
		NewFileContent: `
const obj = {
    myMethod(): number {
        return 42;
    }
};`,
	})
}

func TestRefactorInferReturnTypeExportedFunction(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
export /*a*/function/*b*/() {
    return 42;
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: inferReturnTypeRefactorTitle,
		NewFileContent: `
export function(): number {
    return 42;
}`,
	})
}

func TestRefactorInferReturnTypeDefaultExportedFunction(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
export default /*a*/function/*b*/() {
    return 42;
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: inferReturnTypeRefactorTitle,
		NewFileContent: `
export default function(): number {
    return 42;
}`,
	})
}
