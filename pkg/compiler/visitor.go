package compiler

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/clauses"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/literals"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/operators"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/pkg/errors"
)

type (
	forOption func(f *expressions.ForExpression) error

	visitor struct {
		*fql.BaseFqlParserVisitor
		src   string
		funcs *core.Functions
	}
)

func newVisitor(src string, funcs *core.Functions) *visitor {
	return &visitor{
		&fql.BaseFqlParserVisitor{},
		src,
		funcs,
	}
}

func (v *visitor) VisitProgram(ctx *fql.ProgramContext) interface{} {
	return newResultFrom(func() (interface{}, error) {
		err := v.doVisitHeads(ctx.AllHead())
		if err != nil {
			return nil, err
		}

		gs := newGlobalScope()
		rs := newRootScope(gs)
		block, err := v.doVisitBody(ctx.Body().(*fql.BodyContext), rs)

		if err != nil {
			return nil, err
		}

		return runtime.NewProgram(v.src, block, gs.params)
	})
}

func (v *visitor) doVisitHeads(heads []fql.IHeadContext) error {
	namespaces := map[string]struct{}{}

	for _, head := range heads {
		err := v.doVisitHead(head.(*fql.HeadContext), namespaces)
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *visitor) doVisitHead(head *fql.HeadContext, namespaces map[string]struct{}) error {
	useexpr := head.UseExpression().(*fql.UseExpressionContext)

	// TODO: Think about improving collision analysis to display more detailed errors.
	// For example, "namespaces X and Y both contain function F"
	if iuse := useexpr.Use(); iuse != nil {
		ns := iuse.(*fql.UseContext).
			NamespaceIdentifier().
			GetText()

		if _, exists := namespaces[ns]; exists {
			return errors.Errorf(`namespace "%s" already used`, ns)
		}

		namespaces[ns] = struct{}{}

		err := copyFromNamespace(v.funcs, ns)
		if err != nil {
			return errors.Wrapf(err, `copy from namespace "%s"`, ns)
		}
	}

	return nil
}

func copyFromNamespace(fns *core.Functions, namespace string) error {
	// In the name of the function "A::B::C", the namespace is "A::B",
	// not "A::B::".
	//
	// So add "::" at the end.
	namespace += "::"

	// core.Functions cast every function name to upper case. Thus
	// namespace also shoud be in upper case.
	namespace = strings.ToUpper(namespace)

	for _, name := range fns.Names() {
		if !strings.HasPrefix(name, namespace) {
			continue
		}

		noprefix := strings.Replace(name, namespace, "", 1)

		if _, exists := fns.Get(noprefix); exists {
			return errors.Errorf(
				`collision occured: "%s" already registered`,
				noprefix,
			)
		}

		fn, _ := fns.Get(name)
		fns.Set(noprefix, fn)
	}

	return nil
}

func (v *visitor) doVisitBody(ctx *fql.BodyContext, scope *scope) (core.Expression, error) {
	statements := ctx.AllBodyStatement()
	body := expressions.NewBodyExpression(len(statements) + 1)

	for _, stmt := range statements {
		e, err := v.doVisitBodyStatement(stmt.(*fql.BodyStatementContext), scope)

		if err != nil {
			return nil, err
		}

		body.Add(e)
	}

	exp := ctx.BodyExpression()

	if exp != nil {
		e, err := v.doVisitBodyExpression(exp.(*fql.BodyExpressionContext), scope)

		if err != nil {
			return nil, err
		}

		if e != nil {
			body.Add(e)
		}
	}

	return body, nil
}

func (v *visitor) doVisitBodyStatement(ctx *fql.BodyStatementContext, scope *scope) (core.Expression, error) {
	variable := ctx.VariableDeclaration()

	if variable != nil {
		return v.doVisitVariableDeclaration(variable.(*fql.VariableDeclarationContext), scope)
	}

	funcCall := ctx.FunctionCallExpression()

	if funcCall != nil {
		return v.doVisitFunctionCallExpression(funcCall.(*fql.FunctionCallExpressionContext), scope)
	}

	return nil, core.Error(ErrInvalidToken, ctx.GetText())
}

func (v *visitor) doVisitBodyExpression(ctx *fql.BodyExpressionContext, scope *scope) (core.Expression, error) {
	forIn := ctx.ForExpression()

	if forIn != nil {
		return v.doVisitForExpression(forIn.(*fql.ForExpressionContext), scope)
	}

	ret := ctx.ReturnExpression()

	if ret != nil {
		return v.doVisitReturnExpression(ret.(*fql.ReturnExpressionContext), scope)
	}

	return nil, nil
}

func (v *visitor) doVisitReturnExpression(ctx *fql.ReturnExpressionContext, scope *scope) (core.Expression, error) {
	var exp core.Expression
	expCtx := ctx.Expression()

	if expCtx != nil {
		out, err := v.doVisitExpression(expCtx.(*fql.ExpressionContext), scope)

		if err != nil {
			return nil, err
		}

		exp = out

		return expressions.NewReturnExpression(v.getSourceMap(ctx), exp)
	}

	forIn := ctx.ForExpression()

	if forIn != nil {
		out, err := v.doVisitForExpression(ctx.ForExpression().(*fql.ForExpressionContext), scope.Fork())

		if err != nil {
			return nil, err
		}

		exp = out

		return expressions.NewReturnExpression(v.getSourceMap(ctx), exp)
	}

	forInTernary := ctx.ForTernaryExpression()

	if forInTernary != nil {
		out, err := v.doVisitForTernaryExpression(forInTernary.(*fql.ForTernaryExpressionContext), scope)

		if err != nil {
			return nil, err
		}

		return expressions.NewReturnExpression(v.getSourceMap(ctx), out)
	}

	return nil, ErrNotImplemented
}

func (v *visitor) doVisitForExpression(ctx *fql.ForExpressionContext, scope *scope) (core.Expression, error) {
	var valVarName string
	var keyVarName string

	parsedClauses := make([]forOption, 0, 10)
	forInScope := scope.Fork()

	srcCtx := ctx.ForExpressionSource().(*fql.ForExpressionSourceContext)
	srcExp, err := v.doVisitForExpressionSource(srcCtx, forInScope)

	if err != nil {
		return nil, err
	}

	valVar := ctx.ForExpressionValueVariable()
	valVarName = valVar.GetText()

	if err := forInScope.SetVariable(valVarName); err != nil {
		return nil, err
	}

	keyVar := ctx.ForExpressionKeyVariable()

	if keyVar != nil {
		keyVarName = keyVar.GetText()

		if err := forInScope.SetVariable(keyVarName); err != nil {
			return nil, err
		}
	}

	src, err := expressions.NewDataSource(
		v.getSourceMap(srcCtx),
		valVarName,
		keyVarName,
		srcExp,
	)

	if err != nil {
		return nil, err
	}

	// Clauses.
	// We put clauses parsing before parsing the query body because COLLECT clause overrides scope variables
	for _, e := range ctx.AllForExpressionBody() {
		e := e.(*fql.ForExpressionBodyContext)
		clauseCtx := e.ForExpressionClause()
		statementCtx := e.ForExpressionStatement()

		if clauseCtx != nil {
			setter, err := v.doVisitForExpressionClause(
				clauseCtx.(*fql.ForExpressionClauseContext),
				forInScope,
				valVarName,
				keyVarName,
			)

			if err != nil {
				return nil, err
			}

			parsedClauses = append(parsedClauses, setter)
		} else if statementCtx != nil {
			exp, err := v.doVisitForExpressionStatement(
				statementCtx.(*fql.ForExpressionStatementContext),
				forInScope,
			)

			if err != nil {
				return nil, err
			}

			parsedClauses = append(parsedClauses, exp)
		}
	}

	var spread bool
	var distinct bool
	var predicate core.Expression
	forRetCtx := ctx.ForExpressionReturn().(*fql.ForExpressionReturnContext)
	returnCtx := forRetCtx.ReturnExpression()

	if returnCtx != nil {
		returnCtx := returnCtx.(*fql.ReturnExpressionContext)
		returnExp, err := v.doVisitReturnExpression(returnCtx, forInScope)

		if err != nil {
			return nil, err
		}

		distinctCtx := returnCtx.Distinct()

		if distinctCtx != nil {
			distinct = true
		}

		predicate = returnExp
	} else {
		forInCtx := forRetCtx.ForExpression().(*fql.ForExpressionContext)
		forInExp, err := v.doVisitForExpression(forInCtx, forInScope)

		if err != nil {
			return nil, err
		}

		spread = true

		predicate = forInExp
	}

	forExp, err := expressions.NewForExpression(
		v.getSourceMap(ctx),
		src,
		predicate,
		distinct,
		spread,
	)

	if err != nil {
		return nil, err
	}

	// add all available clauses
	for _, clause := range parsedClauses {
		if err := clause(forExp); err != nil {
			return nil, err
		}
	}

	return forExp, nil
}

func (v *visitor) doVisitLimitClause(ctx *fql.LimitClauseContext, scope *scope) (core.Expression, core.Expression, error) {
	var err error
	var count core.Expression
	var offset core.Expression

	clauseValues := ctx.AllLimitClauseValue()

	if len(clauseValues) > 1 {
		offset, err = v.doVisitLimitClauseValue(clauseValues[0].(*fql.LimitClauseValueContext), scope)

		if err != nil {
			return nil, nil, err
		}

		count, err = v.doVisitLimitClauseValue(clauseValues[1].(*fql.LimitClauseValueContext), scope)

		if err != nil {
			return nil, nil, err
		}
	} else {
		count, err = v.doVisitLimitClauseValue(clauseValues[0].(*fql.LimitClauseValueContext), scope)

		if err != nil {
			return nil, nil, err
		}

		offset = literals.NewIntLiteral(0)
	}

	return count, offset, nil
}

func (v *visitor) doVisitLimitClauseValue(ctx *fql.LimitClauseValueContext, scope *scope) (core.Expression, error) {
	literalCtx := ctx.IntegerLiteral()

	if literalCtx != nil {
		i, err := strconv.Atoi(literalCtx.GetText())

		if err != nil {
			return nil, err
		}

		return literals.NewIntLiteral(i), nil
	}

	paramCtx := ctx.Param()

	return v.doVisitParamContext(paramCtx.(*fql.ParamContext), scope)
}

func (v *visitor) doVisitFilterClause(ctx *fql.FilterClauseContext, scope *scope) (core.Expression, error) {
	exp := ctx.Expression().(*fql.ExpressionContext)

	exps, err := v.doVisitAllExpressions(exp.AllExpression(), scope)

	if err != nil {
		return nil, err
	}

	if len(exps) == 2 {
		left := exps[0]
		right := exps[1]

		equalityOp := exp.EqualityOperator()

		if equalityOp != nil {
			return operators.NewEqualityOperator(v.getSourceMap(ctx), left, right, equalityOp.GetText())
		}

		logicalAndOp := exp.LogicalAndOperator()

		if logicalAndOp != nil {
			return operators.NewLogicalOperator(v.getSourceMap(ctx), left, right, logicalAndOp.GetText())
		}

		logicalOrOp := exp.LogicalOrOperator()

		if logicalOrOp != nil {
			return operators.NewLogicalOperator(v.getSourceMap(ctx), left, right, logicalOrOp.GetText())
		}
	} else {
		// should be unary operator
		return v.doVisitExpression(exp, scope)
	}

	return nil, core.Error(ErrInvalidToken, ctx.GetText())
}

func (v *visitor) doVisitSortClause(ctx *fql.SortClauseContext, scope *scope) ([]*clauses.SorterExpression, error) {
	sortExpCtxs := ctx.AllSortClauseExpression()

	res := make([]*clauses.SorterExpression, len(sortExpCtxs))

	for idx, sortExpCtx := range sortExpCtxs {
		sortExpCtx := sortExpCtx.(*fql.SortClauseExpressionContext)
		exp, err := v.doVisitExpression(sortExpCtx.Expression().(*fql.ExpressionContext), scope)

		if err != nil {
			return nil, err
		}

		direction := collections.SortDirectionAsc
		dir := sortExpCtx.SortDirection()

		if dir != nil {
			direction = collections.SortDirectionFromString(dir.GetText())
		}

		sorterExp, err := clauses.NewSorterExpression(
			exp,
			direction,
		)

		if err != nil {
			return nil, err
		}

		res[idx] = sorterExp
	}

	return res, nil
}

func (v *visitor) doVisitCollectClause(ctx *fql.CollectClauseContext, scope *scope, valVarName string) (*clauses.Collect, error) {
	var err error
	var selectors []*clauses.CollectSelector
	var projection *clauses.CollectProjection
	var count *clauses.CollectCount
	var aggregate *clauses.CollectAggregate
	variables := make([]string, 0, 10)

	groupingCtx := ctx.CollectGrouping()

	if groupingCtx != nil {
		groupingCtx := groupingCtx.(*fql.CollectGroupingContext)
		collectSelectors := groupingCtx.AllCollectSelector()

		// group selectors
		if len(collectSelectors) > 0 {
			selectors = make([]*clauses.CollectSelector, 0, len(collectSelectors))

			for _, cs := range collectSelectors {
				selector, err := v.doVisitCollectSelector(cs.(*fql.CollectSelectorContext), scope)

				if err != nil {
					return nil, err
				}

				selectors = append(selectors, selector)
				variables = append(variables, selector.Variable())
			}
		}

		projectionCtx := ctx.CollectGroupVariable()

		if projectionCtx != nil {
			projectionCtx := projectionCtx.(*fql.CollectGroupVariableContext)
			projectionSelectorCtx := projectionCtx.CollectSelector()
			var projectionSelector *clauses.CollectSelector

			// if projection expression is defined like WITH group = { foo: i.bar }
			if projectionSelectorCtx != nil {
				selector, err := v.doVisitCollectSelector(projectionSelectorCtx.(*fql.CollectSelectorContext), scope)

				if err != nil {
					return nil, err
				}

				projectionSelector = selector
			} else {
				// otherwise, use default expression WITH group = { i }
				projectionIdentifier := projectionCtx.Identifier(0)

				if projectionIdentifier != nil {
					varExp, err := expressions.NewVariableExpression(v.getSourceMap(projectionCtx), valVarName)

					if err != nil {
						return nil, err
					}

					strLitExp := literals.NewStringLiteral(valVarName)

					propExp, err := literals.NewObjectPropertyAssignment(
						strLitExp,
						varExp,
					)

					if err != nil {
						return nil, err
					}

					projectionSelectorExp := literals.NewObjectLiteralWith(propExp)

					selector, err := clauses.NewCollectSelector(projectionIdentifier.GetText(), projectionSelectorExp)

					if err != nil {
						return nil, err
					}

					projectionSelector = selector
				}
			}

			if projectionSelector != nil {
				variables = append(variables, projectionSelector.Variable())
				projection, err = clauses.NewCollectProjection(projectionSelector)

				if err != nil {
					return nil, err
				}
			}
		}
	}

	countCtx := ctx.CollectCounter()

	if countCtx != nil {
		countCtx := countCtx.(*fql.CollectCounterContext)
		variable := countCtx.Identifier().GetText()
		variables = append(variables, variable)

		count, err = clauses.NewCollectCount(variable)

		if err != nil {
			return nil, err
		}
	}

	aggrCtx := ctx.CollectAggregator()

	if aggrCtx != nil {
		aggrCtx := aggrCtx.(*fql.CollectAggregatorContext)

		selectorCtxs := aggrCtx.AllCollectAggregateSelector()
		selectors := make([]*clauses.CollectAggregateSelector, 0, len(selectorCtxs))

		for _, sc := range selectorCtxs {
			selector, err := v.doVisitCollectAggregateSelector(sc.(*fql.CollectAggregateSelectorContext), scope)

			if err != nil {
				return nil, err
			}

			selectors = append(selectors, selector)
			variables = append(variables, selector.Variable())
		}

		aggregate, err = clauses.NewCollectAggregate(selectors)

		if err != nil {
			return nil, err
		}
	}

	// clear all variables defined before
	scope.ClearVariables()

	for _, variable := range variables {
		if err := scope.SetVariable(variable); err != nil {
			return nil, err
		}
	}

	return clauses.NewCollect(selectors, projection, count, aggregate)
}

func (v *visitor) doVisitCollectSelector(ctx *fql.CollectSelectorContext, scope *scope) (*clauses.CollectSelector, error) {
	variable := ctx.Identifier().GetText()
	exp, err := v.doVisitExpression(ctx.Expression().(*fql.ExpressionContext), scope)

	if err != nil {
		return nil, err
	}

	return clauses.NewCollectSelector(variable, exp)
}

func (v *visitor) doVisitCollectAggregateSelector(ctx *fql.CollectAggregateSelectorContext, scope *scope) (*clauses.CollectAggregateSelector, error) {
	variable := ctx.Identifier().GetText()
	fnCtx := ctx.FunctionCallExpression()

	if fnCtx != nil {
		exp, err := v.doVisitFunctionCallExpression(fnCtx.(*fql.FunctionCallExpressionContext), scope)

		if err != nil {
			return nil, err
		}

		fnExp, ok := exp.(*expressions.FunctionCallExpression)

		if !ok {
			return nil, core.Error(core.ErrInvalidType, "expected function expression")
		}

		return clauses.NewCollectAggregateSelector(variable, fnExp.Arguments(), fnExp.Function())
	}

	return nil, core.Error(core.ErrNotFound, "function expression")
}

func (v *visitor) doVisitForExpressionSource(ctx *fql.ForExpressionSourceContext, scope *scope) (core.Expression, error) {
	arr := ctx.ArrayLiteral()

	if arr != nil {
		return v.doVisitArrayLiteral(arr.(*fql.ArrayLiteralContext), scope)
	}

	obj := ctx.ObjectLiteral()

	if obj != nil {
		return v.doVisitObjectLiteral(obj.(*fql.ObjectLiteralContext), scope)
	}

	variable := ctx.Variable()

	if variable != nil {
		return v.doVisitVariable(variable.(*fql.VariableContext), scope)
	}

	funcCall := ctx.FunctionCallExpression()

	if funcCall != nil {
		return v.doVisitFunctionCallExpression(funcCall.(*fql.FunctionCallExpressionContext), scope)
	}

	memberExp := ctx.MemberExpression()

	if memberExp != nil {
		return v.doVisitMemberExpression(memberExp.(*fql.MemberExpressionContext), scope)
	}

	rangeOp := ctx.RangeOperator()

	if rangeOp != nil {
		return v.doVisitRangeOperator(rangeOp.(*fql.RangeOperatorContext), scope)
	}

	param := ctx.Param()

	if param != nil {
		return v.doVisitParamContext(param.(*fql.ParamContext), scope)
	}

	return nil, core.Error(ErrInvalidDataSource, ctx.GetText())
}

func (v *visitor) doVisitForExpressionClause(ctx *fql.ForExpressionClauseContext, scope *scope, valVarName, _ string) (func(f *expressions.ForExpression) error, error) {
	limitCtx := ctx.LimitClause()

	if limitCtx != nil {
		limit, offset, err := v.doVisitLimitClause(limitCtx.(*fql.LimitClauseContext), scope)

		if err != nil {
			return nil, err
		}

		return func(f *expressions.ForExpression) error {
			return f.AddLimit(v.getSourceMap(limitCtx), limit, offset)
		}, nil
	}

	filterCtx := ctx.FilterClause()

	if filterCtx != nil {
		filterExp, err := v.doVisitFilterClause(filterCtx.(*fql.FilterClauseContext), scope)

		if err != nil {
			return nil, err
		}

		return func(f *expressions.ForExpression) error {
			return f.AddFilter(v.getSourceMap(filterCtx), filterExp)
		}, nil
	}

	sortCtx := ctx.SortClause()

	if sortCtx != nil {
		sortCtx := sortCtx.(*fql.SortClauseContext)
		sortExps, err := v.doVisitSortClause(sortCtx, scope)

		if err != nil {
			return nil, err
		}

		return func(f *expressions.ForExpression) error {
			return f.AddSort(v.getSourceMap(sortCtx), sortExps...)
		}, nil
	}

	collectCtx := ctx.CollectClause()

	if collectCtx != nil {
		collectCtx := collectCtx.(*fql.CollectClauseContext)
		params, err := v.doVisitCollectClause(collectCtx, scope, valVarName)

		if err != nil {
			return nil, err
		}

		return func(f *expressions.ForExpression) error {
			return f.AddCollect(v.getSourceMap(collectCtx), params)
		}, nil
	}

	return nil, v.unexpectedToken(ctx)
}

func (v *visitor) doVisitForExpressionStatement(ctx *fql.ForExpressionStatementContext, scope *scope) (func(f *expressions.ForExpression) error, error) {
	variableCtx := ctx.VariableDeclaration()

	if variableCtx != nil {
		variableExp, err := v.doVisitVariableDeclaration(
			variableCtx.(*fql.VariableDeclarationContext),
			scope,
		)

		if err != nil {
			return nil, err
		}

		return func(f *expressions.ForExpression) error {
			return f.AddStatement(variableExp)
		}, nil
	}

	fnCallCtx := ctx.FunctionCallExpression()

	if fnCallCtx != nil {
		fnCallExp, err := v.doVisitFunctionCallExpression(
			fnCallCtx.(*fql.FunctionCallExpressionContext),
			scope,
		)

		if err != nil {
			return nil, err
		}

		return func(f *expressions.ForExpression) error {
			return f.AddStatement(fnCallExp)
		}, nil
	}

	return nil, v.unexpectedToken(ctx)
}

func (v *visitor) doVisitMemberExpression(ctx *fql.MemberExpressionContext, scope *scope) (core.Expression, error) {
	member, err := v.doVisitMember(ctx.Member().(*fql.MemberContext), scope)

	if err != nil {
		return nil, err
	}

	children := ctx.MemberPath().GetChildren()
	path := make([]core.Expression, 0, len(children))

	for _, child := range children {
		_, ok := child.(antlr.TerminalNode)

		if ok {
			continue
		}

		var exp core.Expression
		var err error
		var parsed bool

		prop, ok := child.(*fql.PropertyNameContext)

		if ok {
			exp, err = v.doVisitPropertyNameContext(prop, scope)
			parsed = true
		} else {
			computedProp, ok := child.(*fql.ComputedPropertyNameContext)

			if ok {
				exp, err = v.doVisitComputedPropertyNameContext(computedProp, scope)
				parsed = true
			}
		}

		if err != nil {
			return nil, err
		}

		if !parsed {
			// TODO: add more contextual information
			return nil, ErrInvalidToken
		}

		path = append(path, exp)
	}

	exp, err := expressions.NewMemberExpression(
		v.getSourceMap(ctx),
		member,
		path,
	)

	if err != nil {
		return nil, err
	}

	return exp, nil
}

func (v *visitor) doVisitMember(ctx *fql.MemberContext, scope *scope) (core.Expression, error) {
	identifier := ctx.Identifier()

	if identifier != nil {
		varName := ctx.Identifier().GetText()

		if !scope.HasVariable(varName) {
			return nil, core.Error(ErrVariableNotFound, varName)
		}

		exp, err := expressions.NewVariableExpression(v.getSourceMap(ctx), varName)

		if err != nil {
			return nil, err
		}

		return exp, nil
	}

	fnCall := ctx.FunctionCallExpression()

	if fnCall != nil {
		exp, err := v.doVisitFunctionCallExpression(fnCall.(*fql.FunctionCallExpressionContext), scope)

		if err != nil {
			return nil, err
		}

		return exp, nil
	}

	param := ctx.Param()

	exp, err := v.doVisitParamContext(param.(*fql.ParamContext), scope)

	if err != nil {
		return nil, err
	}

	return exp, nil
}

func (v *visitor) doVisitObjectLiteral(ctx *fql.ObjectLiteralContext, scope *scope) (core.Expression, error) {
	assignments := ctx.AllPropertyAssignment()
	props := make([]*literals.ObjectPropertyAssignment, 0, len(assignments))

	for _, assignment := range assignments {
		var name core.Expression
		var value core.Expression
		var err error

		assignment := assignment.(*fql.PropertyAssignmentContext)

		prop := assignment.PropertyName()
		computedProp := assignment.ComputedPropertyName()
		shortHand := assignment.ShorthandPropertyName()

		switch {
		case prop != nil:
			name, err = v.doVisitPropertyNameContext(prop.(*fql.PropertyNameContext), scope)
		case computedProp != nil:
			name, err = v.doVisitComputedPropertyNameContext(computedProp.(*fql.ComputedPropertyNameContext), scope)
		default:
			name, err = v.doVisitShorthandPropertyNameContext(shortHand.(*fql.ShorthandPropertyNameContext), scope)
		}

		if err != nil {
			return nil, err
		}

		if shortHand == nil {
			value, err = v.visit(assignment.Expression(), scope)
		} else {
			value, err = v.doVisitVariable(shortHand.(*fql.ShorthandPropertyNameContext).Variable().(*fql.VariableContext), scope)
		}

		if err != nil {
			return nil, err
		}

		pa, err := literals.NewObjectPropertyAssignment(name, value)

		if err != nil {
			return nil, err
		}

		props = append(props, pa)
	}

	return literals.NewObjectLiteralWith(props...), nil
}

func (v *visitor) doVisitPropertyNameContext(ctx *fql.PropertyNameContext, scope *scope) (core.Expression, error) {
	var name string

	identifier := ctx.Identifier()

	if identifier != nil {
		name = identifier.GetText()
	} else {
		stringLiteral := ctx.StringLiteral()

		if stringLiteral != nil {
			runes := []rune(stringLiteral.GetText())
			name = string(runes[1 : len(runes)-1])
		} else {
			param, err := v.doVisitParamContext(ctx.Param().(*fql.ParamContext), scope)

			if err != nil {
				return nil, err
			}

			return param, nil
		}
	}

	if name == "" {
		return nil, core.Error(core.ErrNotFound, "property name")
	}

	return literals.NewStringLiteral(name), nil
}

func (v *visitor) doVisitComputedPropertyNameContext(ctx *fql.ComputedPropertyNameContext, scope *scope) (core.Expression, error) {
	return v.doVisitExpression(ctx.Expression().(*fql.ExpressionContext), scope)
}

func (v *visitor) doVisitShorthandPropertyNameContext(ctx *fql.ShorthandPropertyNameContext, scope *scope) (core.Expression, error) {
	name := ctx.Variable().GetText()

	if !scope.HasVariable(name) {
		return nil, core.Error(ErrVariableNotFound, name)
	}

	return literals.NewStringLiteral(ctx.Variable().GetText()), nil
}

func (v *visitor) doVisitArrayLiteral(ctx *fql.ArrayLiteralContext, scope *scope) (core.Expression, error) {
	listCtx := ctx.ArrayElementList()

	if listCtx == nil {
		return literals.NewArrayLiteral(0), nil
	}

	list := listCtx.(*fql.ArrayElementListContext)
	exp := list.AllExpression()
	elements := make([]core.Expression, 0, len(exp))

	for _, e := range exp {
		element, err := v.visit(e, scope)

		if err != nil {
			return nil, err
		}

		elements = append(elements, element)
	}

	return literals.NewArrayLiteralWith(elements...), nil
}

func (v *visitor) doVisitFloatLiteral(ctx *fql.FloatLiteralContext) (core.Expression, error) {
	val, err := strconv.ParseFloat(ctx.GetText(), 64)

	if err != nil {
		return nil, err
	}

	return literals.NewFloatLiteral(val), nil
}

func (v *visitor) doVisitIntegerLiteral(ctx *fql.IntegerLiteralContext) (core.Expression, error) {
	val, err := strconv.Atoi(ctx.GetText())

	if err != nil {
		return nil, err
	}

	return literals.NewIntLiteral(val), nil
}

func (v *visitor) doVisitStringLiteral(ctx *fql.StringLiteralContext) (core.Expression, error) {
	var text string

	strLiteral := ctx.StringLiteral()

	if strLiteral != nil {
		text = strLiteral.GetText()
	}

	// remove extra quotes
	return literals.NewStringLiteral(text[1 : len(text)-1]), nil
}

func (v *visitor) doVisitBooleanLiteral(ctx *fql.BooleanLiteralContext) (core.Expression, error) {
	return literals.NewBooleanLiteral(strings.EqualFold(ctx.GetText(), "TRUE")), nil
}

func (v *visitor) doVisitNoneLiteral(_ *fql.NoneLiteralContext) (core.Expression, error) {
	return literals.None, nil
}

func (v *visitor) doVisitVariable(ctx *fql.VariableContext, scope *scope) (core.Expression, error) {
	name := ctx.Identifier().GetText()

	// check whether the variable is defined
	if !scope.HasVariable(name) {
		return nil, core.Error(ErrVariableNotFound, name)
	}

	return expressions.NewVariableExpression(v.getSourceMap(ctx), name)
}

func (v *visitor) doVisitVariableDeclaration(ctx *fql.VariableDeclarationContext, scope *scope) (core.Expression, error) {
	var init core.Expression
	var err error

	name := ctx.Identifier().GetText()
	err = scope.SetVariable(name)

	if err != nil {
		return nil, err
	}

	exp := ctx.Expression()

	if exp != nil {
		init, err = v.doVisitExpression(ctx.Expression().(*fql.ExpressionContext), scope)
	}

	if init == nil && err == nil {
		forIn := ctx.ForExpression()

		if forIn != nil {
			init, err = v.doVisitForExpression(forIn.(*fql.ForExpressionContext), scope)
		}
	}

	if init == nil && err == nil {
		forTer := ctx.ForTernaryExpression()

		if forTer != nil {
			init, err = v.doVisitForTernaryExpression(forTer.(*fql.ForTernaryExpressionContext), scope)
		}
	}

	if err != nil {
		return nil, err
	}

	return expressions.NewVariableDeclarationExpression(
		v.getSourceMap(ctx),
		name,
		init,
	)
}

func (v *visitor) doVisitRangeOperator(ctx *fql.RangeOperatorContext, scope *scope) (core.Expression, error) {
	exp, err := v.doVisitChildren(ctx, scope)

	if err != nil {
		return nil, err
	}

	if len(exp) < 2 {
		return nil, core.Error(ErrInvalidToken, ctx.GetText())
	}

	left := exp[0]
	right := exp[1]

	return operators.NewRangeOperator(
		v.getSourceMap(ctx),
		left,
		right,
	)
}

func (v *visitor) doVisitFunctionCallExpression(context *fql.FunctionCallExpressionContext, scope *scope) (core.Expression, error) {
	args := make([]core.Expression, 0, 5)
	argsCtx := context.Arguments()

	if argsCtx != nil {
		argsCtx := argsCtx.(*fql.ArgumentsContext)

		for _, arg := range argsCtx.AllExpression() {
			exp, err := v.doVisitExpression(arg.(*fql.ExpressionContext), scope)

			if err != nil {
				return nil, err
			}

			args = append(args, exp)
		}
	}

	var name string

	funcNS := context.Namespace()

	if funcNS != nil {
		name += funcNS.GetText()
	}

	name += context.FunctionIdentifier().GetText()

	fun, exists := v.funcs.Get(name)

	if !exists {
		return nil, core.Error(core.ErrNotFound, fmt.Sprintf("function: '%s'", name))
	}

	return expressions.NewFunctionCallExpression(
		v.getSourceMap(context),
		fun,
		args...,
	)
}

func (v *visitor) doVisitParamContext(context *fql.ParamContext, s *scope) (core.Expression, error) {
	name := context.Identifier().GetText()

	s.AddParam(name)

	return expressions.NewParameterExpression(
		v.getSourceMap(context),
		name,
	)
}

func (v *visitor) doVisitAllExpressions(contexts []fql.IExpressionContext, scope *scope) ([]core.Expression, error) {
	ret := make([]core.Expression, 0, len(contexts))

	for _, ctx := range contexts {
		exp, err := v.doVisitExpression(ctx.(*fql.ExpressionContext), scope)

		if err != nil {
			return nil, err
		}

		ret = append(ret, exp)
	}

	return ret, nil
}

func (v *visitor) doVisitMathOperator(ctx *fql.ExpressionContext, scope *scope) (core.OperatorExpression, error) {
	var operator operators.MathOperatorType
	multiCtx := ctx.MultiplicativeOperator()

	if multiCtx != nil {
		operator = operators.MathOperatorType(multiCtx.GetText())
	} else {
		additiveCtx := ctx.AdditiveOperator()

		if additiveCtx == nil {
			return nil, ErrInvalidToken
		}

		operator = operators.MathOperatorType(additiveCtx.GetText())
	}

	exps, err := v.doVisitAllExpressions(ctx.AllExpression(), scope)

	if err != nil {
		return nil, err
	}

	left := exps[0]
	right := exps[1]

	return operators.NewMathOperator(
		v.getSourceMap(ctx),
		left,
		right,
		operator,
	)
}

func (v *visitor) doVisitUnaryOperator(ctx *fql.ExpressionContext, scope *scope) (core.OperatorExpression, error) {
	op := ctx.UnaryOperator().(*fql.UnaryOperatorContext)

	exps, err := v.doVisitAllExpressions(ctx.AllExpression(), scope)

	if err != nil {
		return nil, err
	}

	exp := exps[0]

	return operators.NewUnaryOperator(
		v.getSourceMap(ctx),
		exp,
		operators.UnaryOperatorType(op.GetText()),
	)
}

func (v *visitor) doVisitLogicalOperator(ctx *fql.ExpressionContext, scope *scope) (core.OperatorExpression, error) {
	var operator string

	logicalAndOp := ctx.LogicalAndOperator()

	if logicalAndOp != nil {
		operator = logicalAndOp.GetText()
	} else {
		logicalOrOp := ctx.LogicalOrOperator()

		if logicalOrOp == nil {
			return nil, ErrInvalidToken
		}

		operator = logicalOrOp.GetText()
	}

	exps, err := v.doVisitAllExpressions(ctx.AllExpression(), scope)

	if err != nil {
		return nil, err
	}

	left := exps[0]
	right := exps[1]

	return operators.NewLogicalOperator(v.getSourceMap(ctx), left, right, operator)
}

func (v *visitor) doVisitEqualityOperator(ctx *fql.ExpressionContext, scope *scope) (core.OperatorExpression, error) {
	equalityOp := ctx.EqualityOperator().(*fql.EqualityOperatorContext)
	exps, err := v.doVisitAllExpressions(ctx.AllExpression(), scope)

	if err != nil {
		return nil, err
	}

	left := exps[0]
	right := exps[1]

	return operators.NewEqualityOperator(v.getSourceMap(equalityOp), left, right, equalityOp.GetText())
}

func (v *visitor) doVisitRegexpOperator(ctx *fql.ExpressionContext, scope *scope) (core.Expression, error) {
	regexpOp := ctx.RegexpOperator().(*fql.RegexpOperatorContext)
	rawExps := ctx.AllExpression()
	exps, err := v.doVisitAllExpressions(rawExps, scope)

	if err != nil {
		return nil, err
	}

	left := exps[0]
	right := exps[1]

	switch lit := right.(type) {
	case literals.StringLiteral:
		_, err := regexp.Compile(string(lit))

		if err != nil {
			src := v.getSourceMap(rawExps[1])

			return nil, errors.Wrap(err, src.String())
		}
	case *literals.ArrayLiteral, *literals.ObjectLiteral, literals.BooleanLiteral, literals.FloatLiteral, literals.IntLiteral:
		src := v.getSourceMap(rawExps[1])

		return nil, errors.Wrap(errors.New("expected a string literal or a function call"), src.String())
	}

	return operators.NewRegexpOperator(v.getSourceMap(regexpOp), left, right, regexpOp.GetText())
}

func (v *visitor) doVisitInOperator(ctx *fql.ExpressionContext, scope *scope) (core.OperatorExpression, error) {
	exps, err := v.doVisitAllExpressions(ctx.AllExpression(), scope)

	if err != nil {
		return nil, err
	}

	op := ctx.InOperator().(*fql.InOperatorContext)

	left := exps[0]
	right := exps[1]

	if len(exps) != 2 {
		return nil, v.unexpectedToken(ctx)
	}

	return operators.NewInOperator(
		v.getSourceMap(ctx),
		left,
		right,
		op.Not() != nil,
	)
}

func (v *visitor) doVisitArrayOperator(ctx *fql.ExpressionContext, scope *scope) (core.OperatorExpression, error) {
	var comparator core.OperatorExpression
	var err error

	switch {
	case ctx.InOperator() != nil:
		comparator, err = v.doVisitInOperator(ctx, scope)
	case ctx.EqualityOperator() != nil:
		comparator, err = v.doVisitEqualityOperator(ctx, scope)
	default:
		return nil, v.unexpectedToken(ctx)
	}

	if err != nil {
		return nil, err
	}

	exps, err := v.doVisitAllExpressions(ctx.AllExpression(), scope)

	if err != nil {
		return nil, err
	}

	if len(exps) != 2 {
		return nil, v.unexpectedToken(ctx)
	}

	left := exps[0]
	right := exps[1]

	aotype, err := operators.ToIsValidArrayOperatorType(ctx.ArrayOperator().GetText())

	if err != nil {
		return nil, err
	}

	return operators.NewArrayOperator(
		v.getSourceMap(ctx),
		left,
		right,
		aotype,
		comparator,
	)
}

func (v *visitor) doVisitExpressionGroup(ctx *fql.ExpressionGroupContext, scope *scope) (core.Expression, error) {
	exp := ctx.Expression()

	if exp == nil {
		return nil, ErrInvalidToken
	}

	return v.doVisitExpression(exp.(*fql.ExpressionContext), scope)
}

func (v *visitor) doVisitExpression(ctx *fql.ExpressionContext, scope *scope) (core.Expression, error) {
	seq := ctx.ExpressionGroup()

	if seq != nil {
		return v.doVisitExpressionGroup(seq.(*fql.ExpressionGroupContext), scope)
	}

	member := ctx.MemberExpression()

	if member != nil {
		return v.doVisitMemberExpression(member.(*fql.MemberExpressionContext), scope)
	}

	funCall := ctx.FunctionCallExpression()

	if funCall != nil {
		return v.doVisitFunctionCallExpression(funCall.(*fql.FunctionCallExpressionContext), scope)
	}

	notOp := ctx.UnaryOperator()

	if notOp != nil {
		return v.doVisitUnaryOperator(ctx, scope)
	}

	multiOp := ctx.MultiplicativeOperator()

	if multiOp != nil {
		return v.doVisitMathOperator(ctx, scope)
	}

	addOp := ctx.AdditiveOperator()

	if addOp != nil {
		return v.doVisitMathOperator(ctx, scope)
	}

	arrOp := ctx.ArrayOperator()

	if arrOp != nil {
		return v.doVisitArrayOperator(ctx, scope)
	}

	equalityOp := ctx.EqualityOperator()

	if equalityOp != nil {
		return v.doVisitEqualityOperator(ctx, scope)
	}

	inOp := ctx.InOperator()

	if inOp != nil {
		return v.doVisitInOperator(ctx, scope)
	}

	logicalAndOp := ctx.LogicalAndOperator()

	if logicalAndOp != nil {
		return v.doVisitLogicalOperator(ctx, scope)
	}

	logicalOrOp := ctx.LogicalOrOperator()

	if logicalOrOp != nil {
		return v.doVisitLogicalOperator(ctx, scope)
	}

	regexpOp := ctx.RegexpOperator()

	if regexpOp != nil {
		return v.doVisitRegexpOperator(ctx, scope)
	}

	variable := ctx.Variable()

	if variable != nil {
		return v.doVisitVariable(variable.(*fql.VariableContext), scope)
	}

	str := ctx.StringLiteral()

	if str != nil {
		return v.doVisitStringLiteral(str.(*fql.StringLiteralContext))
	}

	integ := ctx.IntegerLiteral()

	if integ != nil {
		return v.doVisitIntegerLiteral(integ.(*fql.IntegerLiteralContext))
	}

	float := ctx.FloatLiteral()

	if float != nil {
		return v.doVisitFloatLiteral(float.(*fql.FloatLiteralContext))
	}

	boolean := ctx.BooleanLiteral()

	if boolean != nil {
		return v.doVisitBooleanLiteral(boolean.(*fql.BooleanLiteralContext))
	}

	arr := ctx.ArrayLiteral()

	if arr != nil {
		return v.doVisitArrayLiteral(arr.(*fql.ArrayLiteralContext), scope)
	}

	obj := ctx.ObjectLiteral()

	if obj != nil {
		return v.doVisitObjectLiteral(obj.(*fql.ObjectLiteralContext), scope)
	}

	none := ctx.NoneLiteral()

	if none != nil {
		return v.doVisitNoneLiteral(none.(*fql.NoneLiteralContext))
	}

	questionCtx := ctx.QuestionMark()

	if questionCtx != nil {
		exps, err := v.doVisitAllExpressions(ctx.AllExpression(), scope)

		if err != nil {
			return nil, err
		}

		return v.createTernaryOperator(
			v.getSourceMap(ctx),
			exps,
			scope,
		)
	}

	rangeOp := ctx.RangeOperator()

	if rangeOp != nil {
		return v.doVisitRangeOperator(rangeOp.(*fql.RangeOperatorContext), scope)
	}

	param := ctx.Param()

	if param != nil {
		return v.doVisitParamContext(param.(*fql.ParamContext), scope)
	}

	// TODO: Complete it
	return nil, ErrNotImplemented
}

func (v *visitor) doVisitChildren(node antlr.RuleNode, scope *scope) ([]core.Expression, error) {
	children := node.GetChildren()

	if children == nil {
		return make([]core.Expression, 0), nil
	}

	result := make([]core.Expression, 0, len(children))

	for _, child := range children {
		_, ok := child.(antlr.TerminalNode)

		if ok {
			continue
		}

		out, err := v.visit(child, scope)

		if err != nil {
			return nil, err
		}

		result = append(result, out)
	}

	return result, nil
}

func (v *visitor) visit(node antlr.Tree, scope *scope) (core.Expression, error) {
	var out core.Expression
	var err error

	switch ctx := node.(type) {
	case *fql.BodyContext:
		out, err = v.doVisitBody(ctx, scope)
	case *fql.ExpressionContext:
		out, err = v.doVisitExpression(ctx, scope)
	case *fql.ForExpressionContext:
		out, err = v.doVisitForExpression(ctx, scope)
	case *fql.ReturnExpressionContext:
		out, err = v.doVisitReturnExpression(ctx, scope)
	case *fql.ArrayLiteralContext:
		out, err = v.doVisitArrayLiteral(ctx, scope)
	case *fql.ObjectLiteralContext:
		out, err = v.doVisitObjectLiteral(ctx, scope)
	case *fql.StringLiteralContext:
		out, err = v.doVisitStringLiteral(ctx)
	case *fql.IntegerLiteralContext:
		out, err = v.doVisitIntegerLiteral(ctx)
	case *fql.FloatLiteralContext:
		out, err = v.doVisitFloatLiteral(ctx)
	case *fql.BooleanLiteralContext:
		out, err = v.doVisitBooleanLiteral(ctx)
	case *fql.NoneLiteralContext:
		out, err = v.doVisitNoneLiteral(ctx)
	case *fql.VariableContext:
		out, err = v.doVisitVariable(ctx, scope)
	case *fql.VariableDeclarationContext:
		out, err = v.doVisitVariableDeclaration(ctx, scope)
	case *fql.FunctionCallExpressionContext:
		out, err = v.doVisitFunctionCallExpression(ctx, scope)
	case *fql.ParamContext:
		out, err = v.doVisitParamContext(ctx, scope)
	default:
		err = v.unexpectedToken(node)
	}

	return out, err
}

func (v *visitor) doVisitForTernaryExpression(ctx *fql.ForTernaryExpressionContext, scope *scope) (*expressions.ConditionExpression, error) {
	exps, err := v.doVisitChildren(ctx, scope)

	if err != nil {
		return nil, err
	}

	return v.createTernaryOperator(
		v.getSourceMap(ctx),
		exps,
		scope,
	)
}

func (v *visitor) createTernaryOperator(src core.SourceMap, exps []core.Expression, _ *scope) (*expressions.ConditionExpression, error) {
	var test core.Expression
	var consequent core.Expression
	var alternate core.Expression

	if len(exps) == 3 {
		test = exps[0]
		consequent = exps[1]
		alternate = exps[2]
	} else {
		test = exps[0]
		alternate = exps[1]
	}

	return expressions.NewConditionExpression(
		src,
		test,
		consequent,
		alternate,
	)
}

func (v *visitor) unexpectedToken(node antlr.Tree) error {
	name := "undefined"
	ctx, ok := node.(antlr.RuleContext)

	if ok {
		name = ctx.GetText()
	}

	return errors.Errorf("unexpected token: %s", name)
}

func (v *visitor) getSourceMap(rule antlr.ParserRuleContext) core.SourceMap {
	start := rule.GetStart()

	return core.NewSourceMap(
		rule.GetText(),
		start.GetLine(),
		start.GetColumn(),
	)
}
