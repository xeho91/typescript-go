package ls

import (
	"context"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/locale"
	"github.com/microsoft/typescript-go/internal/ls/change"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/nodebuilder"
)

const inferReturnTypeRefactorKind = lsproto.CodeActionKind("refactor.rewrite.function.returnType")

// InferReturnTypeProvider is a RefactorProvider that adds return type annotations
// by inferring the type from the function body.
var InferReturnTypeProvider = &RefactorProvider{
	RefactorActions: []RefactorAction{
		{
			Title:   "Infer function return type",
			ID:      "inferReturnType",
			Kinds:   []lsproto.CodeActionKind{inferReturnTypeRefactorKind},
			Factory: getInferReturnTypeCodeActions,
		},
	},
}

func getInferReturnTypeCodeActions(ctx context.Context, refactorContext *RefactorContext, refactorID string) ([]*CodeAction, error) {
	if ast.IsInJSFile(refactorContext.SourceFile.AsNode()) {
		return nil, nil
	}

	title := diagnostics.Infer_function_return_type.Localize(locale.FromContext(ctx))

	token := astnav.GetTouchingPropertyName(refactorContext.SourceFile, refactorContext.Range.Pos())

	declaration := findConvertibleAncestor(token)
	if declaration == nil || !hasBody(declaration) || declaration.Type() != nil {
		reason := diagnostics.Return_type_must_be_inferred_from_a_function.Localize(locale.FromContext(ctx))
		return []*CodeAction{{Description: title, Kind: inferReturnTypeRefactorKind, DisabledReason: reason}}, nil
	}

	ch, done := refactorContext.Program.GetTypeCheckerForFile(ctx, refactorContext.SourceFile)
	defer done()

	typeNode := getInferredReturnTypeNode(ch, declaration, refactorContext.SourceFile)
	if typeNode == nil {
		reason := diagnostics.Could_not_determine_function_return_type.Localize(locale.FromContext(ctx))
		return []*CodeAction{{Description: title, Kind: inferReturnTypeRefactorKind, DisabledReason: reason}}, nil
	}

	changeTracker := change.NewTracker(
		ctx,
		refactorContext.Program.Options(),
		refactorContext.LS.FormatOptions(),
		refactorContext.LS.converters,
	)
	if ast.IsArrowFunction(declaration) {
		changeTracker.ParenthesizeArrowParameters(refactorContext.SourceFile, declaration)
	}

	changeTracker.TryInsertTypeAnnotation(refactorContext.SourceFile, declaration, typeNode)

	changes := changeTracker.GetChanges()
	if len(changes) == 0 {
		return nil, nil
	}

	actions := []*CodeAction{{
		Description: title,
		Changes:     changes[refactorContext.SourceFile.FileName()],
		FixID:       refactorID,
		Kind:        inferReturnTypeRefactorKind,
	}}
	return actions, nil
}

func getUnionReturnTypeNode(ch *checker.Checker, declaration *ast.Node, signatures []*checker.Signature, idToSymbol map[*ast.IdentifierNode]*ast.Symbol) *ast.Node {
	if len(signatures) <= 1 {
		return nil
	}

	returnTypes := make([]*checker.Type, 0, len(signatures))

	for _, sig := range signatures {
		rt := ch.GetReturnTypeOfSignature(sig)
		if rt != nil {
			returnTypes = append(returnTypes, rt)
		}
	}

	if len(returnTypes) == 0 {
		return nil
	}

	return ch.TypeToTypeNodeEx(
		ch.GetUnionType(returnTypes),
		declaration,
		nodebuilder.FlagsNoTruncation,
		nodebuilder.InternalFlagsAllowUnresolvedNames,
		idToSymbol,
	)
}

func getOverloadReturnTypeNode(ch *checker.Checker, declaration *ast.Node, idToSymbol map[*ast.IdentifierNode]*ast.Symbol) *ast.Node {
	if !ch.GetEmitResolver().IsImplementationOfOverload(declaration) {
		return nil
	}

	fnType := ch.GetTypeAtLocation(declaration)
	if fnType == nil {
		return nil
	}

	return getUnionReturnTypeNode(ch, declaration, ch.GetCallSignatures(fnType), idToSymbol)
}

func getTypePredicateReturnTypeNode(ch *checker.Checker, declaration *ast.Node, signature *checker.Signature, sourceFile *ast.SourceFile, idToSymbol map[*ast.IdentifierNode]*ast.Symbol) *ast.Node {
	typePredicate := ch.GetTypePredicateOfSignature(signature)
	if typePredicate == nil || typePredicate.Type() == nil {
		return nil
	}

	enclosingDecl := ast.FindAncestor(declaration, ast.IsDeclaration)
	if enclosingDecl == nil {
		enclosingDecl = sourceFile.AsNode()
	}

	return ch.TypePredicateToTypePredicateNodeEx(
		typePredicate,
		enclosingDecl,
		nodebuilder.FlagsNoTruncation,
		nodebuilder.InternalFlagsAllowUnresolvedNames,
		idToSymbol,
	)
}

func getInferredReturnTypeNode(ch *checker.Checker, declaration *ast.Node, sourceFile *ast.SourceFile) *ast.Node {
	idToSymbol := make(map[*ast.IdentifierNode]*ast.Symbol)

	if typeNode := getOverloadReturnTypeNode(ch, declaration, idToSymbol); typeNode != nil {
		return typeNode
	}

	signature := ch.GetSignatureFromDeclaration(declaration)
	if signature == nil {
		return nil
	}

	if typeNode := getTypePredicateReturnTypeNode(ch, declaration, signature, sourceFile, idToSymbol); typeNode != nil {
		return typeNode
	}

	return ch.TypeToTypeNodeEx(
		ch.GetReturnTypeOfSignature(signature),
		declaration,
		nodebuilder.FlagsNoTruncation,
		nodebuilder.InternalFlagsAllowUnresolvedNames,
		idToSymbol,
	)
}

func findConvertibleAncestor(node *ast.Node) *ast.Node {
	for node != nil {
		if ast.IsBlock(node) {
			return nil
		}

		if isArrowFunctionBodyOrToken(node) {
			return nil
		}

		if isConvertibleDeclaration(node) {
			return node
		}

		node = node.Parent
	}

	return nil
}

func isArrowFunctionBodyOrToken(node *ast.Node) bool {
	return node.Parent != nil && ast.IsArrowFunction(node.Parent) &&
		(node.Kind == ast.KindEqualsGreaterThanToken || node.Parent.AsArrowFunction().Body == node)
}

func isConvertibleDeclaration(node *ast.Node) bool {
	return ast.IsFunctionLikeDeclaration(node) &&
		!ast.IsConstructorDeclaration(node) &&
		!ast.IsGetAccessorDeclaration(node) &&
		!ast.IsSetAccessorDeclaration(node)
}

func hasBody(node *ast.Node) bool {
	switch node.Kind {
	case ast.KindFunctionDeclaration:
		return node.AsFunctionDeclaration().Body != nil
	case ast.KindFunctionExpression:
		return node.AsFunctionExpression().Body != nil
	case ast.KindArrowFunction:
		return node.AsArrowFunction().Body != nil
	case ast.KindMethodDeclaration:
		return node.AsMethodDeclaration().Body != nil
	}

	return false
}
