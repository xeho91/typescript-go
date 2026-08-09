package ls

import (
	"context"
	"fmt"
	"slices"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/locale"
	"github.com/microsoft/typescript-go/internal/ls/change"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/scanner"
)

const (
	extractTypeRefactorKind      = lsproto.CodeActionKind("refactor.extract")
	extractToTypeAliasActionKind = lsproto.CodeActionKind("refactor.extract.type")
	extractToInterfaceActionKind = lsproto.CodeActionKind("refactor.extract.interface")
	extractToTypeDefActionKind   = lsproto.CodeActionKind("refactor.extract.typedef")
)

// ExtractTypeProvider provides the "Extract type" refactoring.
var ExtractTypeProvider = &RefactorProvider{
	RefactorActions: []RefactorAction{
		{
			Title: "Extract type",
			ID:    "extractType",
			Kinds: []lsproto.CodeActionKind{
				extractTypeRefactorKind,
				extractToTypeAliasActionKind,
				extractToInterfaceActionKind,
				extractToTypeDefActionKind,
			},
			Factory: getExtractTypeCodeActions,
		},
	},
}

type extractTypeSelection struct {
	first *ast.Node
	last  *ast.Node
	nodes []*ast.Node
}

type extractTypeInfo struct {
	selection      extractTypeSelection
	enclosingNode  *ast.Node
	typeParameters []*ast.Node
	typeElements   []*ast.Node
}

// Empty (cursor) selections only offer the refactor when invoked.
func isRefactorTriggeredInvoked(refactorContext *RefactorContext) bool {
	triggerKind := refactorContext.CodeActionTriggerKind()
	return triggerKind != nil && *triggerKind == lsproto.CodeActionTriggerKindInvoked
}

func getExtractTypeCodeActions(ctx context.Context, refactorContext *RefactorContext, refactorID string) ([]*CodeAction, error) {
	file := refactorContext.SourceFile
	isJS := ast.IsInJSFile(file.AsNode())

	showNotApplicable := showNotApplicableReasons(ctx, refactorContext.LS)

	ch, done := refactorContext.Program.GetTypeCheckerForFile(ctx, refactorContext.SourceFile)
	defer done()

	range_ := refactorContext.Range
	considerEmptySpans := range_.Pos() == range_.End() && isRefactorTriggeredInvoked(refactorContext)

	result := getRangeToExtract(refactorContext, ch, considerEmptySpans)
	if result == nil {
		return buildNotApplicableAction(ctx, showNotApplicable), nil
	}

	return buildExtractTypeActions(ctx, refactorContext, result, isJS, refactorContext.LS.FormatOptions()), nil
}

func buildExtractTypeActions(ctx context.Context, refactorContext *RefactorContext, info *extractTypeInfo, isJS bool, formatOptions lsutil.FormatCodeSettings) []*CodeAction {
	if isJS {
		return nil
	}

	var actions []*CodeAction

	act := buildTypeAliasAction(ctx, refactorContext, info, formatOptions)
	if act != nil {
		actions = append(actions, act)
	}

	if info.typeElements != nil {
		act := buildInterfaceAction(ctx, refactorContext, info, formatOptions)
		if act != nil {
			actions = append(actions, act)
		}
	}

	return actions
}

func buildNotApplicableAction(ctx context.Context, showNotApplicable bool) []*CodeAction {
	if !showNotApplicable {
		return nil
	}

	return []*CodeAction{{
		Description:    diagnostics.Extract_type.Localize(locale.FromContext(ctx)),
		Kind:           extractTypeRefactorKind,
		DisabledReason: diagnostics.Selection_is_not_a_valid_type_node.Localize(locale.FromContext(ctx)),
	}}
}

func buildTypeAliasAction(ctx context.Context, refactorContext *RefactorContext, info *extractTypeInfo, formatOptions lsutil.FormatCodeSettings) *CodeAction {
	title := diagnostics.Extract_to_type_alias.Localize(locale.FromContext(ctx))
	file := refactorContext.SourceFile

	changeTracker := change.NewTracker(ctx, refactorContext.Program.Options(), formatOptions, refactorContext.LS.converters)

	name := newExtractedTypeName(changeTracker, file)

	newTypeNode := getNewTypeNode(info, changeTracker.NodeFactory)
	if newTypeNode == nil {
		newTypeNode = info.selection.first
	}

	typeAlias := changeTracker.NewTypeAliasDeclaration(
		nil,
		name,
		buildTypeParameterList(info.typeParameters, changeTracker.NodeFactory),
		newTypeNode,
	)
	changeTracker.InsertNodeBefore(file, info.enclosingNode, typeAlias, true, change.LeadingTriviaOptionNone)
	replaceSelectionWithTypeReference(changeTracker, file, info, name)

	return finalizeCodeAction(changeTracker, file, title, extractToTypeAliasActionKind)
}

func buildInterfaceAction(ctx context.Context, refactorContext *RefactorContext, info *extractTypeInfo, formatOptions lsutil.FormatCodeSettings) *CodeAction {
	title := diagnostics.Extract_to_interface.Localize(locale.FromContext(ctx))
	file := refactorContext.SourceFile

	changeTracker := change.NewTracker(ctx, refactorContext.Program.Options(), formatOptions, refactorContext.LS.converters)

	name := newExtractedTypeName(changeTracker, file)

	interface_ := changeTracker.NewInterfaceDeclaration(
		nil,
		name,
		buildTypeParameterList(info.typeParameters, changeTracker.NodeFactory),
		nil,
		changeTracker.NewNodeList(info.typeElements),
	)
	changeTracker.InsertNodeBefore(file, info.enclosingNode, interface_, true, change.LeadingTriviaOptionNone)
	replaceSelectionWithTypeReference(changeTracker, file, info, name)

	return finalizeCodeAction(changeTracker, file, title, extractToInterfaceActionKind)
}

func newExtractedTypeName(changeTracker *change.Tracker, file *ast.SourceFile) *ast.Node {
	return changeTracker.NewIdentifier(getUniqueName("NewType", file))
}

func replaceSelectionWithTypeReference(changeTracker *change.Tracker, file *ast.SourceFile, info *extractTypeInfo, name *ast.Node) {
	refTypeParams := buildTypeParamRefList(info.typeParameters, changeTracker.NodeFactory)
	typeRef := changeTracker.NewTypeReferenceNode(name, refTypeParams)
	rng := changeTracker.GetAdjustedRange(file, info.selection.first, info.selection.last, change.LeadingTriviaOptionExclude, change.TrailingTriviaOptionExclude)
	changeTracker.ReplaceRange(file, rng, typeRef, change.NodeOptions{LeadingTriviaOption: change.LeadingTriviaOptionExclude, TrailingTriviaOption: change.TrailingTriviaOptionExclude})
}

func finalizeCodeAction(changeTracker *change.Tracker, file *ast.SourceFile, title string, kind lsproto.CodeActionKind) *CodeAction {
	changes := changeTracker.GetChanges()
	if len(changes) == 0 {
		return nil
	}

	edits := changes[file.FileName()]
	if len(edits) == 0 {
		return nil
	}

	return &CodeAction{
		Description: title,
		Changes:     edits,
		FixID:       "extractType",
		Kind:        kind,
	}
}

func getUniqueName(baseName string, file *ast.SourceFile) string {
	name := baseName
	for i := 1; !isFileLevelUniqueName(file, name); i++ {
		name = fmt.Sprintf("%s_%d", baseName, i)
	}
	return name
}

func isFileLevelUniqueName(file *ast.SourceFile, name string) bool {
	return !file.HasIdentifier(name)
}

func getRangeToExtract(refactorContext *RefactorContext, ch *checker.Checker, considerEmptySpans bool) *extractTypeInfo {
	file := refactorContext.SourceFile
	range_ := refactorContext.Range
	isCursorRequest := range_.Pos() == range_.End() && considerEmptySpans

	firstType := getFirstTypeAt(file, int(range_.Pos()), range_, isCursorRequest)
	if firstType == nil || !ast.IsTypeNode(firstType) {
		return nil
	}

	isJS := ast.IsInJSFile(file.AsNode())

	enclosingNode := getEnclosingNode(firstType, isJS)
	if enclosingNode == nil {
		return nil
	}

	expandedFirstType := getExpandedSelectionNode(firstType, enclosingNode)
	if expandedFirstType == nil || !ast.IsTypeNode(expandedFirstType) {
		return nil
	}

	typeList := getOverlappingTypes(expandedFirstType, firstType, range_)

	var selection extractTypeSelection
	if len(typeList) > 1 {
		selection = extractTypeSelection{
			first: typeList[0],
			last:  typeList[len(typeList)-1],
			nodes: typeList,
		}
	} else {
		selection = extractTypeSelection{
			first: expandedFirstType,
			last:  expandedFirstType,
		}
	}

	typeParameters, ok := collectTypeParameters(ch, selection, enclosingNode, file)
	if !ok {
		return nil
	}

	typeElements := flattenTypeLiteralNodeReference(selection)

	return &extractTypeInfo{
		selection:      selection,
		enclosingNode:  enclosingNode,
		typeParameters: typeParameters,
		typeElements:   typeElements,
	}
}

func getOverlappingTypes(expandedFirstType, firstType *ast.Node, range_ core.TextRange) []*ast.Node {
	parent := expandedFirstType.Parent
	if parent == nil {
		return nil
	}

	isUnionOrIntersection := isUnionOrIntersectionType(parent)
	selectionExtendsBeyondFirstType := range_.End() > firstType.End()
	if !isUnionOrIntersection || !selectionExtendsBeyondFirstType {
		return nil
	}

	var types []*ast.Node
	if ast.IsUnionTypeNode(parent) {
		types = parent.AsUnionTypeNode().Types.Nodes
	} else {
		types = parent.AsIntersectionTypeNode().Types.Nodes
	}

	var typeList []*ast.Node
	for _, t := range types {
		if nodeOverlapsRange(t, range_) {
			typeList = append(typeList, t)
		}
	}

	return typeList
}

func isUnionOrIntersectionType(node *ast.Node) bool {
	return ast.IsUnionTypeNode(node) || ast.IsIntersectionTypeNode(node)
}

func getFirstTypeAt(file *ast.SourceFile, startPosition int, range_ core.TextRange, isCursorRequest bool) *ast.Node {
	if firstType := findFirstTypeInRange(astnav.GetTokenAtPosition(file, startPosition), file, range_, isCursorRequest); firstType != nil {
		return firstType
	}

	if firstType := findFirstTypeInRange(astnav.GetTouchingTokenIncludingPreceding(file, startPosition), file, range_, isCursorRequest); firstType != nil {
		return firstType
	}

	return nil
}

func findFirstTypeInRange(current *ast.Node, file *ast.SourceFile, range_ core.TextRange, isCursorRequest bool) *ast.Node {
	if current == nil {
		return nil
	}

	return findTypeAncestorInRange(current, file, range_, isCursorRequest, nodeOverlapsRange(current, range_))
}

func findTypeAncestorInRange(node *ast.Node, file *ast.SourceFile, range_ core.TextRange, isCursorRequest bool, overlappingRange bool) *ast.Node {
	return ast.FindAncestor(node, func(node *ast.Node) bool {
		if node.Parent == nil || !ast.IsTypeNode(node) {
			return false
		}

		if rangeContainsSkipTrivia(range_, node.Parent, file) {
			return false
		}

		if !isCursorRequest && !overlappingRange {
			return false
		}

		return true
	})
}

func getNewTypeNode(info *extractTypeInfo, factory *ast.NodeFactory) *ast.Node {
	if len(info.selection.nodes) > 0 {
		nodes := info.selection.nodes
		parent := nodes[0].Parent
		if parent != nil && ast.IsUnionTypeNode(parent) {
			return factory.NewUnionTypeNode(factory.NewNodeList(nodes))
		}

		return factory.NewIntersectionTypeNode(factory.NewNodeList(nodes))
	}

	return info.selection.first
}

func buildTypeParameterList(params []*ast.Node, factory *ast.NodeFactory) *ast.NodeList {
	if len(params) == 0 {
		return nil
	}

	nodes := make([]*ast.Node, len(params))
	for i, p := range params {
		tpd := p.AsTypeParameterDeclaration()
		nodes[i] = factory.NewTypeParameterDeclaration(tpd.Modifiers(), tpd.Name(), tpd.Constraint, nil, nil)
	}

	return factory.NewNodeList(nodes)
}

func buildTypeParamRefList(params []*ast.Node, factory *ast.NodeFactory) *ast.NodeList {
	if len(params) == 0 {
		return nil
	}

	nodes := make([]*ast.Node, len(params))
	for i, p := range params {
		nodes[i] = factory.NewTypeReferenceNode(p.AsTypeParameterDeclaration().Name(), nil)
	}

	return factory.NewNodeList(nodes)
}

func getEnclosingNode(node *ast.Node, isJS bool) *ast.Node {
	enclosing := ast.FindAncestor(node, func(n *ast.Node) bool {
		return ast.IsStatement(n)
	})
	if enclosing == nil && isJS {
		enclosing = ast.FindAncestor(node, func(n *ast.Node) bool {
			return ast.IsJSDoc(n)
		})
	}

	return enclosing
}

func getExpandedSelectionNode(firstType *ast.Node, enclosingNode *ast.Node) *ast.Node {
	expanded := ast.FindAncestorOrQuit(firstType, func(node *ast.Node) ast.FindAncestorResult {
		if node == enclosingNode {
			return ast.FindAncestorQuit
		}

		if node.Parent != nil && isUnionOrIntersectionType(node.Parent) {
			return ast.FindAncestorTrue
		}

		return ast.FindAncestorFalse
	})

	if expanded != nil {
		return expanded
	}

	return firstType
}

func nodeOverlapsRange(node *ast.Node, range_ core.TextRange) bool {
	nodeRange := core.NewTextRange(node.Pos(), node.End())
	return nodeRange.Overlaps(range_)
}

func rangeContainsSkipTrivia(r1 core.TextRange, node *ast.Node, file *ast.SourceFile) bool {
	startPos := scanner.SkipTrivia(file.Text(), node.Pos())
	return core.NewTextRange(startPos, node.End()).ContainedBy(r1)
}

func collectTypeParameters(ch *checker.Checker, selection extractTypeSelection, enclosingNode *ast.Node, file *ast.SourceFile) ([]*ast.Node, bool) {
	selectionArray := selection.nodes
	if len(selectionArray) == 0 && selection.first != nil {
		selectionArray = []*ast.Node{selection.first}
	}
	if len(selectionArray) == 0 {
		return nil, false
	}

	selectionStart := scanner.SkipTrivia(file.Text(), selectionArray[0].Pos())
	selectionRange := core.NewTextRange(selectionStart, selectionArray[len(selectionArray)-1].End())

	enclosingRange := core.NewTextRange(enclosingNode.Pos(), enclosingNode.End())

	var typeParameters []*ast.Node
	for _, node := range selectionArray {
		var abort bool
		typeParameters, abort = collectTypeParametersInNode(ch, node, file, selectionRange, enclosingRange, typeParameters)
		if abort {
			return nil, false
		}
	}

	return typeParameters, true
}

func collectTypeParametersInNode(ch *checker.Checker, node *ast.Node, file *ast.SourceFile, selectionRange, enclosingRange core.TextRange, typeParameters []*ast.Node) ([]*ast.Node, bool) {
	if node == nil {
		return typeParameters, false
	}

	switch {
	case ast.IsTypeReferenceNode(node):
		var abort bool
		typeParameters, abort = collectOrAbortTypeReference(ch, node, file, selectionRange, enclosingRange, typeParameters)
		if abort {
			return typeParameters, true
		}
	case ast.IsInferTypeNode(node):
		if abortForInferType(node, file, selectionRange) {
			return typeParameters, true
		}
	case ast.IsTypePredicateNode(node) || ast.IsThisTypeNode(node):
		if abortForPredicateOrThisType(node, file, selectionRange) {
			return typeParameters, true
		}
	case ast.IsTypeQueryNode(node):
		if abortForTypeQuery(ch, node, file, selectionRange, enclosingRange) {
			return typeParameters, true
		}
	}

	var childAbort bool
	node.ForEachChild(func(child *ast.Node) bool {
		if childAbort {
			return true
		}

		var abort bool
		typeParameters, abort = collectTypeParametersInNode(ch, child, file, selectionRange, enclosingRange, typeParameters)
		childAbort = abort
		return abort
	})

	return typeParameters, childAbort
}

func collectOrAbortTypeReference(ch *checker.Checker, node *ast.Node, file *ast.SourceFile, selectionRange, enclosingRange core.TextRange, typeParameters []*ast.Node) ([]*ast.Node, bool) {
	typeName := node.AsTypeReferenceNode().TypeName
	if !ast.IsIdentifier(typeName) {
		return typeParameters, false
	}

	symbol := ch.ResolveName(typeName.Text(), typeName, ast.SymbolFlagsTypeParameter, true)
	if symbol == nil {
		return typeParameters, false
	}

	for _, decl := range symbol.Declarations {
		if !ast.IsTypeParameterDeclaration(decl) || ast.GetSourceFileOfNode(decl) != file {
			continue
		}

		if decl.Name().Text() == typeName.Text() && selectionRange.ContainedBy(core.NewTextRange(decl.Pos(), decl.End())) {
			return typeParameters, true
		}

		if shouldHoistTypeParameter(typeParameters, decl, enclosingRange, selectionRange, file) {
			typeParameters = append(typeParameters, decl)
			break
		}
	}

	return typeParameters, false
}

func shouldHoistTypeParameter(typeParameters []*ast.Node, decl *ast.Node, enclosingRange, selectionRange core.TextRange, file *ast.SourceFile) bool {
	return rangeContainsSkipTrivia(enclosingRange, decl, file) &&
		!rangeContainsSkipTrivia(selectionRange, decl, file) &&
		!slices.Contains(typeParameters, decl)
}

func abortForInferType(node *ast.Node, file *ast.SourceFile, selectionRange core.TextRange) bool {
	conditionalTypeNode := ast.FindAncestor(node, func(n *ast.Node) bool {
		return ast.IsConditionalTypeNode(n) &&
			inferTypeInConditionalExtends(node, n, file)
	})

	return conditionalTypeNode == nil || !rangeContainsSkipTrivia(selectionRange, conditionalTypeNode, file)
}

func inferTypeInConditionalExtends(node, n *ast.Node, file *ast.SourceFile) bool {
	extendsType := n.AsConditionalTypeNode().ExtendsType
	return rangeContainsSkipTrivia(core.NewTextRange(extendsType.Pos(), extendsType.End()), node, file)
}

func abortForPredicateOrThisType(node *ast.Node, file *ast.SourceFile, selectionRange core.TextRange) bool {
	functionLikeNode := ast.FindAncestor(node.Parent, ast.IsFunctionLike)
	if functionLikeNode == nil || functionLikeNode.Type() == nil {
		return false
	}

	isInReturnType := rangeContainsSkipTrivia(core.NewTextRange(functionLikeNode.Type().Pos(), functionLikeNode.Type().End()), node, file)
	functionOutsideSelection := !rangeContainsSkipTrivia(selectionRange, functionLikeNode, file)

	return isInReturnType && functionOutsideSelection
}

func abortForTypeQuery(ch *checker.Checker, node *ast.Node, file *ast.SourceFile, selectionRange, enclosingRange core.TextRange) bool {
	if ast.IsIdentifier(node.AsTypeQueryNode().ExprName) {
		symbol := ch.ResolveName(node.AsTypeQueryNode().ExprName.Text(), node.AsTypeQueryNode().ExprName, ast.SymbolFlagsValue, false)
		if symbol == nil || symbol.ValueDeclaration == nil {
			return false
		}

		declaredInEnclosingNode := rangeContainsSkipTrivia(enclosingRange, symbol.ValueDeclaration, file)
		declaredOutsideSelection := !rangeContainsSkipTrivia(selectionRange, symbol.ValueDeclaration, file)

		return declaredInEnclosingNode && declaredOutsideSelection
	}

	return ast.IsThisIdentifier(node.AsTypeQueryNode().ExprName.AsQualifiedName().Left) && !rangeContainsSkipTrivia(selectionRange, node.Parent, file)
}

func flattenTypeLiteralNodeReference(selection extractTypeSelection) []*ast.Node {
	if len(selection.nodes) > 0 {
		var allMembers []*ast.Node

		for _, s := range selection.nodes {
			members := flattenSingleTypeLiteral(s)
			if members == nil {
				return nil
			}

			allMembers = append(allMembers, members...)
		}

		return allMembers
	}

	if selection.first != nil {
		return flattenSingleTypeLiteral(selection.first)
	}

	return nil
}

func flattenSingleTypeLiteral(node *ast.Node) []*ast.Node {
	switch {
	case ast.IsIntersectionTypeNode(node):
		return flattenIntersectionTypeLiteral(node)
	case ast.IsParenthesizedTypeNode(node):
		return flattenSingleTypeLiteral(node.AsParenthesizedTypeNode().Type)
	case ast.IsTypeLiteralNode(node):
		members := node.AsTypeLiteralNode().Members
		if members == nil {
			return nil
		}
		return members.Nodes
	}

	return nil
}

func flattenIntersectionTypeLiteral(node *ast.Node) []*ast.Node {
	var allMembers []*ast.Node
	seen := make(map[string]bool)

	for _, t := range node.AsIntersectionTypeNode().Types.Nodes {
		members := flattenSingleTypeLiteral(t)
		if members == nil {
			return nil
		}

		var ok bool
		allMembers, ok = addUniqueMembers(allMembers, members, seen)
		if !ok {
			return nil
		}
	}

	return allMembers
}

func addUniqueMembers(allMembers, members []*ast.Node, seen map[string]bool) ([]*ast.Node, bool) {
	for _, member := range members {
		name := getNameOfTypeElement(member)
		if name != "" {
			if seen[name] {
				return allMembers, false
			}

			seen[name] = true
		}

		allMembers = append(allMembers, member)
	}

	return allMembers, true
}

func getNameOfTypeElement(node *ast.Node) string {
	name := node.Name()
	if name != nil {
		return name.Text()
	}

	return ""
}
