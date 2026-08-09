package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestRecursive_selfRefMultiFile(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
// @Filename: a.ts
interface Foo<T extends { prop: T }> {}

// @Filename: b.ts
interface Foo<T extends { prop: /*a*/T/*b*/ }> {}
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToFile(t, "b.ts")
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, "Extract to type alias")
}

func TestRecursive_selfRefRanges(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
// @Filename: a.ts
 interface Foo<T extends { prop: [|T|] }> {}

// @Filename: b.ts
 // Some initial comments.
 // We need to ensure these files have different contents,
 // so their ranges differ, so we'll start a few lines below in this file.
 interface Foo<T extends { prop: [|T|] }> {}
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	for _, rangeMarker := range f.Ranges() {
		f.GoToSelectRange(t, rangeMarker)
		f.VerifyRefactorNotAvailable(t, "Extract to type alias")
	}
}

func TestRecursive_selfRefTwoDecls(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
// @Filename: a.ts
interface Foo<T extends { prop: T }> {}
interface Foo<T extends { prop: /*a*/T/*b*/ }> {}
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToFile(t, "a.ts")
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactorNotAvailable(t, "Extract to type alias")
}
