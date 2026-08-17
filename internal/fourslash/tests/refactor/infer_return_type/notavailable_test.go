package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestRefactorInferReturnTypeNotAvailableWithReturnType(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
function /*a*/simple/*b*/(): number {
    return 42;
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeRefactorTitle)
}

func TestRefactorInferReturnTypeNotAvailableConstructor(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
class MyClass {
    /*a*/constructor/*b*/(x: number) {
    }
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeRefactorTitle)
}

func TestRefactorInferReturnTypeNotAvailableGetAccessor(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
class MyClass {
    get /*a*/myProp/*b*/() {
        return 42;
    }
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeRefactorTitle)
}

func TestRefactorInferReturnTypeNotAvailableSetAccessor(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
class MyClass {
    set /*a*/myProp/*b*/(value: number) {
    }
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeRefactorTitle)
}

func TestRefactorInferReturnTypeNotAvailableArrowFunctionExpression(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
const obj = {
    /*a*/myMethod/*b*/: (x: number) => x * 2
};`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeRefactorTitle)
}

func TestRefactorInferReturnTypeNotAvailableCursorOnLeadingTrivia(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
/*a*//*b*/
function hello() {
    return 42;
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeRefactorTitle)
}

func TestRefactorInferReturnTypeNotAvailableCursorOnLeadingTriviaBeforeArrow(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
/*a*//*b*/
const f = (x: number) => x * 2;`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeRefactorTitle)
}

func TestRefactorInferReturnTypeNotAvailableCaretOnVariableName(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
const /*a*//*b*/f = () => 1;`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeRefactorTitle)
}

func TestRefactorInferReturnTypeNotAvailableCaretOnEquals(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
const f =/*a*//*b*/ () => 1;`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeRefactorTitle)
}

func TestRefactorInferReturnTypeNotAvailableConstVariable(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
const /*a*/myConst/*b*/ = 42;`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeRefactorTitle)
}

func TestRefactorInferReturnTypeNotAvailableLetVariable(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
let /*a*/myLet/*b*/ = 42;`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeRefactorTitle)
}

func TestRefactorInferReturnTypeNotAvailableTypePredicate(t *testing.T) {
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

func TestRefactorInferReturnTypeNotAvailableAssertionSignature(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
function /*a*/assertString/*b*/(x: any): asserts x is string {
    if (typeof x !== "string") {
        throw new Error("Not a string");
    }
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeRefactorTitle)
}
