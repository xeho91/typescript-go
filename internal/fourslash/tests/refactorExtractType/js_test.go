package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestJs_extractPrimitive(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// BLOCKED BY: no JSDoc emission in the Go printer (emitJSDocNode panics, internal/printer/printer.go:4572)
	// and changeTracker.getOptionsForInsertNodeBefore panics on KindJSDoc (internal/ls/change/tracker.go:644).
	t.Skip("Extract to typedef is unimplemented in Go")
	const content = `
// @allowJs: true
// @Filename: a.js
/** @type { /*a*/string/*b*/ } */
var x;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to typedef",
		NewFileContent: `
/**
 * @typedef {string} NewType
 */

/** @type { NewType } */
var x;
`,
	})
}

func TestJs_extractLeftUnionMember(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// BLOCKED BY: no JSDoc emission in the Go printer (emitJSDocNode panics, internal/printer/printer.go:4572)
	// and changeTracker.getOptionsForInsertNodeBefore panics on KindJSDoc (internal/ls/change/tracker.go:644).
	t.Skip("Extract to typedef is unimplemented in Go")
	const content = `
// @allowJs: true
// @Filename: a.js
/** @type { /*a*/string/*b*/ | number } */
var x;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to typedef",
		NewFileContent: `
/**
 * @typedef {string} NewType
 */

/** @type { NewType | number } */
var x;
`,
	})
}

func TestJs_extractRightUnionMember(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// BLOCKED BY: no JSDoc emission in the Go printer (emitJSDocNode panics, internal/printer/printer.go:4572)
	// and changeTracker.getOptionsForInsertNodeBefore panics on KindJSDoc (internal/ls/change/tracker.go:644).
	t.Skip("Extract to typedef is unimplemented in Go")
	const content = `
// @allowJs: true
// @Filename: a.js
/** @type { /*a*/string | number/*b*/ } */
var x;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to typedef",
		NewFileContent: `
/**
 * @typedef {string | number} NewType
 */

/** @type { NewType } */
var x;
`,
	})
}

func TestJs_extractTParamsNoConstraint(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// BLOCKED BY: no JSDoc emission in the Go printer (emitJSDocNode panics, internal/printer/printer.go:4572)
	// and changeTracker.getOptionsForInsertNodeBefore panics on KindJSDoc (internal/ls/change/tracker.go:644).
	t.Skip("Extract to typedef is unimplemented in Go")
	const content = `
// @allowJs: true
// @Filename: a.js
/**
 * @template T
 * @template U
 * @param {T} b
 * @param {U} c
 * @returns {/*a*/T | U/*b*/}
 */
function a(b, c) {}
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to typedef",
		NewFileContent: `
/**
 * @template T
 * @template U
 * @typedef {T | U} NewType
 */

/**
 * @template T
 * @template U
 * @param {T} b
 * @param {U} c
 * @returns {NewType<T, U>}
 */
function a(b, c) {}
`,
	})
}

func TestJs_extractTParamsWithConstraint(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// BLOCKED BY: no JSDoc emission in the Go printer (emitJSDocNode panics, internal/printer/printer.go:4572)
	// and changeTracker.getOptionsForInsertNodeBefore panics on KindJSDoc (internal/ls/change/tracker.go:644).
	t.Skip("Extract to typedef is unimplemented in Go")
	const content = `
// @allowJs: true
// @Filename: a.js
/**
 * @template {number} T
 * @template {string} U
 * @param {T} b
 * @param {U} c
 * @returns {/*a*/T | U/*b*/}
 */
function a(b, c) {}
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to typedef",
		NewFileContent: `
/**
 * @template {number} T
 * @template {string} U
 * @typedef {T | U} NewType
 */

/**
 * @template {number} T
 * @template {string} U
 * @param {T} b
 * @param {U} c
 * @returns {NewType<T, U>}
 */
function a(b, c) {}
`,
	})
}

func TestJs_extractTParamsMixedConstraint(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// BLOCKED BY: no JSDoc emission in the Go printer (emitJSDocNode panics, internal/printer/printer.go:4572)
	// and changeTracker.getOptionsForInsertNodeBefore panics on KindJSDoc (internal/ls/change/tracker.go:644).
	t.Skip("Extract to typedef is unimplemented in Go")
	const content = `
// @allowJs: true
// @Filename: a.js
/**
 * @template {number} T, U
 * @param {T} b
 * @param {U} c
 * @returns {/*a*/T | U/*b*/}
 */
function a(b, c) {}
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to typedef",
		NewFileContent: `
/**
 * @template {number} T
 * @template U
 * @typedef {T | U} NewType
 */

/**
 * @template {number} T, U
 * @param {T} b
 * @param {U} c
 * @returns {NewType<T, U>}
 */
function a(b, c) {}
`,
	})
}

func TestJs_extractParamType(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// BLOCKED BY: no JSDoc emission in the Go printer (emitJSDocNode panics, internal/printer/printer.go:4572)
	// and changeTracker.getOptionsForInsertNodeBefore panics on KindJSDoc (internal/ls/change/tracker.go:644).
	t.Skip("Extract to typedef is unimplemented in Go")
	const content = `
// @allowJs: true
// @Filename: a.js
function a(/** @type {/*a*/string/*b*/} */ b) {}
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to typedef",
		NewFileContent: `
/**
 * @typedef {string} NewType
 */

function a(/** @type {NewType} */ b) {}
`,
	})
}

func TestJs_extractTypeTag(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// BLOCKED BY: no JSDoc emission in the Go printer (emitJSDocNode panics, internal/printer/printer.go:4572)
	// and changeTracker.getOptionsForInsertNodeBefore panics on KindJSDoc (internal/ls/change/tracker.go:644).
	t.Skip("Extract to typedef is unimplemented in Go")
	const content = `
// @allowJs: true
// @Filename: a.js
/**
 * @type {/*a*/Foo/*b*/}
 */
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to typedef",
		NewFileContent: `
/**
 * @typedef {Foo} NewType
 */

/**
 * @type {NewType}
 */
`,
	})
}

func TestJs_extractMiddleUnionMember(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// BLOCKED BY: no JSDoc emission in the Go printer (emitJSDocNode panics, internal/printer/printer.go:4572)
	// and changeTracker.getOptionsForInsertNodeBefore panics on KindJSDoc (internal/ls/change/tracker.go:644).
	t.Skip("Extract to typedef is unimplemented in Go")
	const content = `
// @allowJs: true
// @Filename: a.js
/** @type { /*a*/string | number/*b*/ | boolean } */
var x;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to typedef",
		NewFileContent: `
/**
 * @typedef {string | number} NewType
 */

/** @type { NewType | boolean } */
var x;
`,
	})
}

func TestJs_unionCommentBeforeMember(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// BLOCKED BY: no JSDoc emission in the Go printer (emitJSDocNode panics, internal/printer/printer.go:4572)
	// and changeTracker.getOptionsForInsertNodeBefore panics on KindJSDoc (internal/ls/change/tracker.go:644).
	t.Skip("Extract to typedef is unimplemented in Go")
	const content = `
// @allowJs: true
// @Filename: a.js
type Bar = /*a*/string | /* oops */ boolean/*b*/;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to typedef",
		NewFileContent: `
/**
 * @typedef {string | boolean} NewType
 */

type Bar = NewType;
`,
	})
}

func TestJs_unionCommentAfterMember(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// BLOCKED BY: no JSDoc emission in the Go printer (emitJSDocNode panics, internal/printer/printer.go:4572)
	// and changeTracker.getOptionsForInsertNodeBefore panics on KindJSDoc (internal/ls/change/tracker.go:644).
	t.Skip("Extract to typedef is unimplemented in Go")
	const content = `
// @allowJs: true
// @Filename: a.js
type Bar = /*a*/string | boolean /* oops */ |/*b*/;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to typedef",
		NewFileContent: `
/**
 * @typedef {string | boolean} NewType
 */

type Bar = NewType /* oops */ ;
`,
	})
}

func TestJs_typeLiteralCommentBeforeProperty(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// BLOCKED BY: no JSDoc emission in the Go printer (emitJSDocNode panics, internal/printer/printer.go:4572)
	// and changeTracker.getOptionsForInsertNodeBefore panics on KindJSDoc (internal/ls/change/tracker.go:644).
	t.Skip("Extract to typedef is unimplemented in Go")
	const content = `
// @allowJs: true
// @Filename: a.js
type Foo = /*a*/{
    /**
     *
     */
    oops: string;
}/*b*/;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to typedef",
		NewFileContent: `
/**
 * @typedef {{
    oops: string;
}} NewType
 */

type Foo = NewType;
`,
	})
}

func TestJs_typeLiteralCommentAfterProperty(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// BLOCKED BY: no JSDoc emission in the Go printer (emitJSDocNode panics, internal/printer/printer.go:4572)
	// and changeTracker.getOptionsForInsertNodeBefore panics on KindJSDoc (internal/ls/change/tracker.go:644).
	t.Skip("Extract to typedef is unimplemented in Go")
	const content = `
// @allowJs: true
// @Filename: a.js
type Foo = /*a*/{
    oops: string;
    /**
     *
     */
}/*b*/;
`

	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToSelect(t, "a", "b")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title: "Extract to typedef",
		NewFileContent: `
/**
 * @typedef {{
    oops: string;
}} NewType
 */

type Foo = NewType;
`,
	})
}
