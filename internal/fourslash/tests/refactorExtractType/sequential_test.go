package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func invokedTrigger() *lsproto.CodeActionTriggerKind {
	kind := lsproto.CodeActionTriggerKindInvoked
	return &kind
}

func TestSequential(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `
interface Yadda<T> { x: T }

export let blah: Yadda/*a*/<string>/*b*/;

interface YaddaWithDefault<T = boolean/*c*/> { x: T/*d*/ }
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToMarker(t, "a")
	f.VerifyRefactorAvailableForTriggerReason(t, "invoked", "Extract to type alias")

	f.GoToMarker(t, "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title:        "Extract to type alias",
		ApplyChanges: true,
		TriggerKind:  invokedTrigger(),
		NewFileContent: `
interface Yadda<T> { x: T }

type NewType = Yadda<string>;

export let blah: NewType;

interface YaddaWithDefault<T = boolean> { x: T }
`,
	})

	f.GoToMarker(t, "c")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title:        "Extract to type alias",
		ApplyChanges: true,
		TriggerKind:  invokedTrigger(),
		NewFileContent: `
interface Yadda<T> { x: T }

type NewType = Yadda<string>;

export let blah: NewType;

type NewType_1 = boolean;

interface YaddaWithDefault<T = NewType_1> { x: T }
`,
	})

	f.GoToMarker(t, "d")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title:        "Extract to type alias",
		ApplyChanges: true,
		TriggerKind:  invokedTrigger(),
		NewFileContent: `
interface Yadda<T> { x: T }

type NewType = Yadda<string>;

export let blah: NewType;

type NewType_1 = boolean;

type NewType_2<T> = T;

interface YaddaWithDefault<T = NewType_1> { x: NewType_2<T> }
`,
	})
}
