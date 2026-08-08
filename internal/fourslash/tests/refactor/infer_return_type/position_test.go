package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestRefactorInferReturnTypeFunctionPositions(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
function /*a*//*b*/f1() {
    return { x: 1, y: 1 };
}
func/*c*//*d*/tion f2() {
    return { x: 1, y: 1 };
}
function f3(/*e*//*f*/) {
    return { x: 1, y: 1 };
}
function f4() {/*g*//*h*/
    return { x: 1, y: 1 };
}
function f5() {
    return { x: 1, y: 1 /*i*//*j*/};
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorAvailable(t, inferReturnTypeRefactorTitle)

	f.GoToSelect(t, "c", "d")
	f.VerifyRefactorAvailable(t, inferReturnTypeRefactorTitle)

	f.GoToSelect(t, "e", "f")
	f.VerifyRefactorAvailable(t, inferReturnTypeRefactorTitle)

	f.GoToSelect(t, "g", "h")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeRefactorTitle)

	f.GoToSelect(t, "i", "j")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeRefactorTitle)
}

func TestRefactorInferReturnTypeMethodPositions(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
class F1 {
    /*a*//*b*/method() {
        return { x: 1, y: 1 };
    }
}
class F2 {
    met/*c*//*d*/hod(){
        return { x: 1, y: 1 };
    }
}
class F3 {
    method(/*e*//*f*/) {
        return { x: 1, y: 1 };
    }
}
class F4 {
    method() {
        return { x: 1, y: 1 /*g*//*h*/};
    }
}
class F5 {
    method() {
        return { x: 1, y: 1 /*i*//*j*/};
    }
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorAvailable(t, inferReturnTypeRefactorTitle)

	f.GoToSelect(t, "c", "d")
	f.VerifyRefactorAvailable(t, inferReturnTypeRefactorTitle)

	f.GoToSelect(t, "e", "f")
	f.VerifyRefactorAvailable(t, inferReturnTypeRefactorTitle)

	f.GoToSelect(t, "g", "h")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeRefactorTitle)

	f.GoToSelect(t, "i", "j")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeRefactorTitle)
}

func TestRefactorInferReturnTypeFunctionExpressionPositions(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
const f1 = /*a*//*b*/function () {
    return { x: 1, y: 1 };
}
const f2 = func/*c*//*d*/tion() {
    return { x: 1, y: 1 };
}
const f3 = function (/*e*//*f*/) {
    return { x: 1, y: 1 };
}
const f4 = function () {/*g*//*h*/
    return { x: 1, y: 1 };
}
const f5 = function () {
    return { x: 1, y: 1 /*i*//*j*/};
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorAvailable(t, inferReturnTypeRefactorTitle)

	f.GoToSelect(t, "c", "d")
	f.VerifyRefactorAvailable(t, inferReturnTypeRefactorTitle)

	f.GoToSelect(t, "e", "f")
	f.VerifyRefactorAvailable(t, inferReturnTypeRefactorTitle)

	f.GoToSelect(t, "g", "h")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeRefactorTitle)

	f.GoToSelect(t, "i", "j")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeRefactorTitle)
}

func TestRefactorInferReturnTypeArrowPositions(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
const f1 = /*a*//*b*/() =>{
    return { x: 1, y: 1 };
}
const f2 = (/*c*//*d*/) => {
    return { x: 1, y: 1 };
}
const f3 = () => /*e*//*f*/{
    return { x: 1, y: 1 };
}
const f4 = () => {/*g*//*h*/
    return { x: 1, y: 1 };
}
const f5 = () => {
    return { x: 1, y: 1/*i*//*j*/ };
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorAvailable(t, inferReturnTypeRefactorTitle)

	f.GoToSelect(t, "c", "d")
	f.VerifyRefactorAvailable(t, inferReturnTypeRefactorTitle)

	f.GoToSelect(t, "e", "f")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeRefactorTitle)

	f.GoToSelect(t, "g", "h")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeRefactorTitle)

	f.GoToSelect(t, "i", "j")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeRefactorTitle)
}

func TestRefactorInferReturnTypeParameterPositions(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
function f1(/*a*//*b*/x, { y }) {
    return { x, y: 1 };
}
function f2(x, { /*c*//*d*/y }) {
    return { x, y: 1 };
}
function f3(x, /*e*//*f*/{ y }) {
    return { x, y: 1 };
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorAvailable(t, inferReturnTypeRefactorTitle)

	f.GoToSelect(t, "c", "d")
	f.VerifyRefactorAvailable(t, inferReturnTypeRefactorTitle)

	f.GoToSelect(t, "e", "f")
	f.VerifyRefactorAvailable(t, inferReturnTypeRefactorTitle)
}

func TestRefactorInferReturnTypeArrowBodyPositions(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
const f1 = /*a*//*b*/(x: number) => x;
const f2 = (x: number) /*c*//*d*/=> x;
const f3 = (x: number) => /*e*//*f*/x
const f4= (x: number) => x/*g*//*h*/`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorAvailable(t, inferReturnTypeRefactorTitle)

	f.GoToSelect(t, "c", "d")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeRefactorTitle)

	f.GoToSelect(t, "e", "f")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeRefactorTitle)

	f.GoToSelect(t, "g", "h")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeRefactorTitle)
}
