package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestRefactorInferReturnTypeAny(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
function /*a*/foo/*b*/() {
    const bar = 1 as any;
    return bar;
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: inferReturnTypeRefactorTitle,
		NewFileContent: `
function foo(): any {
    const bar = 1 as any;
    return bar;
}`,
	})
}

func TestRefactorInferReturnTypeGeneric(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
function /*a*/foo<T>/*b*/() {
    return 1 as T;
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: inferReturnTypeRefactorTitle,
		NewFileContent: `
function foo<T>(): T {
    return 1 as T;
}`,
	})
}

func TestRefactorInferReturnTypeTernary(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
function /*a*/foo/*b*/(x: number) {
    return x ? x : x > 1;
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: inferReturnTypeRefactorTitle,
		NewFileContent: `
function foo(x: number): number | boolean {
    return x ? x : x > 1;
}`,
	})
}

func TestRefactorInferReturnTypeSwitch(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
interface F1 { x: number; y: number; }
type T1 = [number, number];

function /*a*/foo/*b*/(num: number) {
   switch (num) {
      case 1:
         return { x: num, y: num } as F1;
      case 2:
         return [num, num] as T1;
      default:
         return num;
   }
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: inferReturnTypeRefactorTitle,
		NewFileContent: `
interface F1 { x: number; y: number; }
type T1 = [number, number];

function foo(num: number): number | F1 | T1 {
   switch (num) {
      case 1:
         return { x: num, y: num } as F1;
      case 2:
         return [num, num] as T1;
      default:
         return num;
   }
}`,
	})
}

func TestRefactorInferReturnTypeGeneratorFunction(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
function* /*a*/generator/*b*/() {
    yield 1;
    yield 2;
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: inferReturnTypeRefactorTitle,
		NewFileContent: `
function* generator(): Generator<1 | 2, void, unknown> {
    yield 1;
    yield 2;
}`,
	})
}

func TestRefactorInferReturnTypeComputedProperty(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
const key = "foo";
function /*a*/withComputed/*b*/() {
    return { [key]: 42 };
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: inferReturnTypeRefactorTitle,
		NewFileContent: `
const key = "foo";
function withComputed(): {
    foo: number;
} {
    return { [key]: 42 };
}`,
	})
}
