package compiler

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/pkg/errors"

	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/clauses"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/literals"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/operators"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	forOption func(f *expressions.ForExpression) error

	visitor struct {
		*fql.BaseFqlParserVisitor
		src   string
		funcs *core.Functions
	}
)

const (
	waitPseudoVariable = "CURRENT"
)

const (
	waitScope = "waitfor"
	forScope  = "for"
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
		err := v.visitHeads(ctx.AllHead())
		if err != nil {
			return nil, err
		}

		gs := newGlobalScope()
		rs := newRootScope(gs)
		block, err := v.visitBody(ctx.Body().(fql.IBodyContext), rs)
		if err != nil {
			return nil, err
		}

		return runtime.NewProgram(v.src, block, gs.params)
	})
}

func (v *visitor) visitHeads(heads []fql.IHeadContext) error {
	namespaces := map[string]struct{}{}

	for _, head := range heads {
		err := v.visitHead(head.(fql.IHeadContext), namespaces)
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *visitor) visitHead(c fql.IHeadContext, namespaces map[string]struct{}) error {
	ctx := c.(*fql.HeadContext)
	useexpr := ctx.UseExpression().(*fql.UseExpressionContext)

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
	// namespace also should be in upper case.
	namespace = strings.ToUpper(namespace)

	for _, name := range fns.Names() {
		if !strings.HasPrefix(name, namespace) {
			continue
		}

		noprefix := strings.Replace(name, namespace, "", 1)

		if _, exists := fns.Get(noprefix); exists {
			return errors.Errorf(
				`collision occurred: "%s" already registered`,
				noprefix,
			)
		}

		fn, _ := fns.Get(name)
		fns.Set(noprefix, fn)
	}

	return nil
}

func (v *visitor) visitBody(c fql.IBodyContext, scope *scope) (core.Expression, error) {
	ctx := c.(*fql.BodyContext)
	statements := ctx.AllBodyStatement()
	body := expressions.NewBodyExpression(len(statements) + 1)

	for _, stmt := range statements {
		e, err := v.visitBodyStatement(stmt.(fql.IBodyStatementContext), scope)
		if err != nil {
			return nil, err
		}

		body.Add(e)
	}

	exp := ctx.BodyExpression()

	if exp != nil {
		e, err := v.visitBodyExpression(exp.(fql.IBodyExpressionContext), scope)
		if err != nil {
			return nil, err
		}

		if e != nil {
			body.Add(e)
		}
	}

	return body, nil
}

func (v *visitor) visitBodyStatement(c fql.IBodyStatementContext, scope *scope) (core.Expression, error) {
	ctx := c.(*fql.BodyStatementContext)

	if variable := ctx.VariableDeclaration(); variable != nil {
		return v.visitVariableDeclaration(variable.(fql.IVariableDeclarationContext), scope)
	}

	if funcCall := ctx.FunctionCallExpression(); funcCall != nil {
		return v.visitFunctionCallExpression(funcCall.(fql.IFunctionCallExpressionContext), scope)
	}

	if waitfor := ctx.WaitForExpression(); waitfor != nil {
		return v.visitWaitForExpression(waitfor.(fql.IWaitForExpressionContext), scope)
	}

	return nil, v.unexpectedToken(ctx)
}

func (v *visitor) visitBodyExpression(c fql.IBodyExpressionContext, scope *scope) (core.Expression, error) {
	ctx := c.(*fql.BodyExpressionContext)

	if exp := ctx.ForExpression(); exp != nil {
		return v.visitForExpression(exp.(fql.IForExpressionContext), scope)
	}

	if exp := ctx.ReturnExpression(); exp != nil {
		return v.visitReturnExpression(exp.(fql.IReturnExpressionContext), scope)
	}

	return nil, v.unexpectedToken(ctx)
}

func (v *visitor) visitReturnExpression(c fql.IReturnExpressionContext, scope *scope) (core.Expression, error) {
	var out core.Expression
	var err error
	ctx := c.(*fql.ReturnExpressionContext)

	if exp := ctx.Expression(); exp != nil {
		out, err = v.visitExpression(exp.(fql.IExpressionContext), scope)
	} else {
		return nil, v.unexpectedToken(ctx)
	}

	if err != nil {
		return nil, err
	}

	return expressions.NewReturnExpression(v.getSourceMap(ctx), out)
}

func (v *visitor) visitForExpression(c fql.IForExpressionContext, scope *scope) (core.Expression, error) {
	var err error
	var valVarName string
	var keyVarName string
	var ds collections.Iterable

	ctx := c.(*fql.ForExpressionContext)
	expVars := ctx.AllIdentifier()

	if len(expVars) > 0 {
		valVarName = expVars[0].GetText()
	}

	if len(expVars) > 1 {
		keyVarName = expVars[1].GetText()
	}

	isWhileLoop := ctx.In() == nil

	if !isWhileLoop {
		srcCtx := ctx.ForExpressionSource().(*fql.ForExpressionSourceContext)
		srcExp, err := v.visitForExpressionSource(srcCtx, scope)

		if err != nil {
			return nil, err
		}

		ds, err = expressions.NewForInIterableExpression(
			v.getSourceMap(srcCtx),
			valVarName,
			keyVarName,
			srcExp,
		)

		if err != nil {
			return nil, err
		}
	} else {
		whileExpCtx := ctx.Expression().(fql.IExpressionContext)
		conditionExp, err := v.visitExpression(whileExpCtx, scope)

		if err != nil {
			return nil, err
		}

		var mode collections.WhileMode

		if ctx.Do() != nil {
			mode = collections.WhileModePre
		}

		ds, err = expressions.NewForWhileIterableExpression(
			v.getSourceMap(whileExpCtx),
			mode,
			valVarName,
			conditionExp,
		)

		if err != nil {
			return nil, err
		}
	}

	forInScope := scope.Fork(forScope)
	if err := forInScope.SetVariable(valVarName); err != nil {
		return nil, err
	}

	if keyVarName != "" {
		if err := forInScope.SetVariable(keyVarName); err != nil {
			return nil, err
		}
	}

	parsedClauses := make([]forOption, 0, 10)

	// Clauses.
	// We put clauses parsing before parsing the query body because COLLECT clause overrides scope variables
	for _, e := range ctx.AllForExpressionBody() {
		e := e.(*fql.ForExpressionBodyContext)
		clauseCtx := e.ForExpressionClause()
		statementCtx := e.ForExpressionStatement()

		if clauseCtx != nil {
			setter, err := v.visitForExpressionClause(
				clauseCtx.(fql.IForExpressionClauseContext),
				forInScope,
				valVarName,
				keyVarName,
			)

			if err != nil {
				return nil, err
			}

			parsedClauses = append(parsedClauses, setter)
		} else if statementCtx != nil {
			exp, err := v.visitForExpressionStatement(
				statementCtx.(fql.IForExpressionStatementContext),
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
	var passThrough bool
	forRetCtx := ctx.ForExpressionReturn().(*fql.ForExpressionReturnContext)
	returnCtx := forRetCtx.ReturnExpression()

	if returnCtx != nil {
		returnCtx := returnCtx.(*fql.ReturnExpressionContext)
		returnExp, err := v.visitReturnExpression(returnCtx, forInScope)

		if err != nil {
			return nil, err
		}

		distinctCtx := returnCtx.Distinct()

		if distinctCtx != nil {
			distinct = true
		}

		predicate = returnExp

		ret := returnExp.(*expressions.ReturnExpression)
		passThrough = literals.IsNone(ret.Predicate())
	} else {
		forInCtx := forRetCtx.ForExpression().(*fql.ForExpressionContext)
		forInExp, err := v.visitForExpression(forInCtx, forInScope)

		if err != nil {
			return nil, err
		}

		spread = true

		predicate = forInExp
	}

	forExp, err := expressions.NewForExpression(
		v.getSourceMap(ctx),
		ds,
		predicate,
		distinct,
		spread,
		passThrough,
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

func (v *visitor) visitLimitClause(c fql.ILimitClauseContext, scope *scope) (core.Expression, core.Expression, error) {
	var err error
	var count core.Expression
	var offset core.Expression

	ctx := c.(*fql.LimitClauseContext)
	clauseValues := ctx.AllLimitClauseValue()

	if len(clauseValues) > 1 {
		offset, err = v.visitLimitClauseValue(clauseValues[0].(fql.ILimitClauseValueContext), scope)

		if err != nil {
			return nil, nil, err
		}

		count, err = v.visitLimitClauseValue(clauseValues[1].(fql.ILimitClauseValueContext), scope)

		if err != nil {
			return nil, nil, err
		}
	} else {
		count, err = v.visitLimitClauseValue(clauseValues[0].(fql.ILimitClauseValueContext), scope)

		if err != nil {
			return nil, nil, err
		}

		offset = literals.NewIntLiteral(0)
	}

	return count, offset, nil
}

func (v *visitor) visitLimitClauseValue(c fql.ILimitClauseValueContext, scope *scope) (core.Expression, error) {
	ctx := c.(*fql.LimitClauseValueContext)

	if literalCtx := ctx.IntegerLiteral(); literalCtx != nil {
		i, err := strconv.Atoi(literalCtx.GetText())
		if err != nil {
			return nil, err
		}

		return literals.NewIntLiteral(i), nil
	} else if paramCtx := ctx.Param(); paramCtx != nil {
		return v.visitParam(paramCtx.(fql.IParamContext), scope)
	} else if variableCtx := ctx.Variable(); variableCtx != nil {
		return v.visitVariable(variableCtx, scope)
	} else if funcCtx := ctx.FunctionCallExpression(); funcCtx != nil {
		return v.visitFunctionCallExpression(funcCtx, scope)
	} else if memCtx := ctx.MemberExpression(); memCtx != nil {
		return v.visitMemberExpression(memCtx, scope)
	}

	return nil, v.unexpectedToken(ctx)
}

func (v *visitor) visitFilterClause(c fql.IFilterClauseContext, scope *scope) (core.Expression, error) {
	return v.visitExpression(c.(*fql.FilterClauseContext).Expression(), scope)
}

func (v *visitor) visitSortClause(c fql.ISortClauseContext, scope *scope) ([]*clauses.SorterExpression, error) {
	ctx := c.(*fql.SortClauseContext)
	sortExpCtxs := ctx.AllSortClauseExpression()

	res := make([]*clauses.SorterExpression, len(sortExpCtxs))

	for idx, sortExpCtx := range sortExpCtxs {
		sortExpCtx := sortExpCtx.(*fql.SortClauseExpressionContext)
		exp, err := v.visitExpression(sortExpCtx.Expression().(fql.IExpressionContext), scope)
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

func (v *visitor) visitCollectClause(c fql.ICollectClauseContext, scope *scope, valVarName string) (*clauses.Collect, error) {
	ctx := c.(*fql.CollectClauseContext)
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
				selector, err := v.visitCollectSelector(cs.(fql.ICollectSelectorContext), scope)
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
				selector, err := v.visitCollectSelector(projectionSelectorCtx.(fql.ICollectSelectorContext), scope)
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
			selector, err := v.visitCollectAggregateSelector(sc.(fql.ICollectAggregateSelectorContext), scope)
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

func (v *visitor) visitCollectSelector(c fql.ICollectSelectorContext, scope *scope) (*clauses.CollectSelector, error) {
	ctx := c.(*fql.CollectSelectorContext)
	variable := ctx.Identifier().GetText()
	exp, err := v.visitExpression(ctx.Expression().(fql.IExpressionContext), scope)

	if err != nil {
		return nil, err
	}

	return clauses.NewCollectSelector(variable, exp)
}

func (v *visitor) visitCollectAggregateSelector(c fql.ICollectAggregateSelectorContext, scope *scope) (*clauses.CollectAggregateSelector, error) {
	ctx := c.(*fql.CollectAggregateSelectorContext)
	variable := ctx.Identifier().GetText()
	fnCtx := ctx.FunctionCallExpression()

	if fnCtx != nil {
		exp, err := v.visitFunctionCallExpression(fnCtx.(fql.IFunctionCallExpressionContext), scope)

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

func (v *visitor) visitForExpressionSource(c fql.IForExpressionSourceContext, scope *scope) (core.Expression, error) {
	ctx := c.(*fql.ForExpressionSourceContext)

	arr := ctx.ArrayLiteral()

	if arr != nil {
		return v.visitArrayLiteral(arr.(fql.IArrayLiteralContext), scope)
	}

	obj := ctx.ObjectLiteral()

	if obj != nil {
		return v.visitObjectLiteral(obj.(fql.IObjectLiteralContext), scope)
	}

	variable := ctx.Variable()

	if variable != nil {
		return v.visitVariable(variable.(fql.IVariableContext), scope)
	}

	funcCall := ctx.FunctionCallExpression()

	if funcCall != nil {
		return v.visitFunctionCallExpression(funcCall.(fql.IFunctionCallExpressionContext), scope)
	}

	memberExp := ctx.MemberExpression()

	if memberExp != nil {
		return v.visitMemberExpression(memberExp.(fql.IMemberExpressionContext), scope)
	}

	rangeOp := ctx.RangeOperator()

	if rangeOp != nil {
		return v.visitRangeOperator(rangeOp.(fql.IRangeOperatorContext), scope)
	}

	param := ctx.Param()

	if param != nil {
		return v.visitParam(param.(fql.IParamContext), scope)
	}

	return nil, core.Error(ErrInvalidDataSource, ctx.GetText())
}

func (v *visitor) visitForExpressionClause(c fql.IForExpressionClauseContext, scope *scope, valVarName, _ string) (func(f *expressions.ForExpression) error, error) {
	ctx := c.(*fql.ForExpressionClauseContext)

	limitCtx := ctx.LimitClause()

	if limitCtx != nil {
		limit, offset, err := v.visitLimitClause(limitCtx.(fql.ILimitClauseContext), scope)
		if err != nil {
			return nil, err
		}

		return func(f *expressions.ForExpression) error {
			return f.AddLimit(v.getSourceMap(limitCtx), limit, offset)
		}, nil
	}

	filterCtx := ctx.FilterClause()

	if filterCtx != nil {
		filterExp, err := v.visitFilterClause(filterCtx.(fql.IFilterClauseContext), scope)
		if err != nil {
			return nil, err
		}

		return func(f *expressions.ForExpression) error {
			return f.AddFilter(v.getSourceMap(filterCtx), filterExp)
		}, nil
	}

	sortCtx := ctx.SortClause()

	if sortCtx != nil {
		sortCtx := sortCtx.(fql.ISortClauseContext)
		sortExps, err := v.visitSortClause(sortCtx, scope)
		if err != nil {
			return nil, err
		}

		return func(f *expressions.ForExpression) error {
			return f.AddSort(v.getSourceMap(sortCtx), sortExps...)
		}, nil
	}

	collectCtx := ctx.CollectClause()

	if collectCtx != nil {
		collectCtx := collectCtx.(fql.ICollectClauseContext)
		params, err := v.visitCollectClause(collectCtx, scope, valVarName)
		if err != nil {
			return nil, err
		}

		return func(f *expressions.ForExpression) error {
			return f.AddCollect(v.getSourceMap(collectCtx), params)
		}, nil
	}

	return nil, v.unexpectedToken(ctx)
}

func (v *visitor) visitForExpressionStatement(c fql.IForExpressionStatementContext, scope *scope) (func(f *expressions.ForExpression) error, error) {
	ctx := c.(*fql.ForExpressionStatementContext)

	variableCtx := ctx.VariableDeclaration()

	if variableCtx != nil {
		variableExp, err := v.visitVariableDeclaration(
			variableCtx.(fql.IVariableDeclarationContext),
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
		fnCallExp, err := v.visitFunctionCallExpression(
			fnCallCtx.(fql.IFunctionCallExpressionContext),
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

func (v *visitor) visitOptionsClause(c fql.IOptionsClauseContext, s *scope) (core.Expression, error) {
	ctx := c.(*fql.OptionsClauseContext)

	return v.visitObjectLiteral(ctx.ObjectLiteral(), s)
}

func (v *visitor) visitWaitForEventNameContext(c fql.IWaitForEventNameContext, s *scope) (core.Expression, error) {
	ctx := c.(*fql.WaitForEventNameContext)

	if str := ctx.StringLiteral(); str != nil {
		return v.visitStringLiteral(str)
	}

	if variable := ctx.Variable(); variable != nil {
		return v.visitVariable(variable, s)
	}

	if param := ctx.Param(); param != nil {
		return v.visitParam(param, s)
	}

	if member := ctx.MemberExpression(); member != nil {
		return v.visitMemberExpression(member, s)
	}

	if fnCall := ctx.FunctionCallExpression(); fnCall != nil {
		return v.visitFunctionCallExpression(fnCall, s)
	}

	return nil, ErrNotImplemented
}

func (v *visitor) visitWaitForEventSourceContext(c fql.IWaitForEventSourceContext, s *scope) (core.Expression, error) {
	ctx := c.(*fql.WaitForEventSourceContext)

	if variable := ctx.Variable(); variable != nil {
		return v.visitVariable(variable, s)
	}

	if member := ctx.MemberExpression(); member != nil {
		return v.visitMemberExpression(member, s)
	}

	if fnCall := ctx.FunctionCallExpression(); fnCall != nil {
		return v.visitFunctionCallExpression(fnCall, s)
	}

	return nil, ErrNotImplemented
}

func (v *visitor) visitTimeoutClause(c fql.ITimeoutClauseContext, s *scope) (core.Expression, error) {
	ctx := c.(*fql.TimeoutClauseContext)

	if integer := ctx.IntegerLiteral(); integer != nil {
		return v.visitIntegerLiteral(integer)
	}

	if variable := ctx.Variable(); variable != nil {
		return v.visitVariable(variable, s)
	}

	if param := ctx.Param(); param != nil {
		return v.visitParam(param, s)
	}

	if member := ctx.MemberExpression(); member != nil {
		return v.visitMemberExpression(member, s)
	}

	if fnCall := ctx.FunctionCall(); fnCall != nil {
		return v.visitFunctionCall(fnCall, s)
	}

	return nil, ErrNotImplemented
}

func (v *visitor) visitWaitForExpression(c fql.IWaitForExpressionContext, s *scope) (core.Expression, error) {
	ctx := c.(*fql.WaitForExpressionContext)

	eventName, err := v.visitWaitForEventNameContext(ctx.WaitForEventName(), s)

	if err != nil {
		return nil, errors.Wrap(err, "invalid event name")
	}

	eventSource, err := v.visitWaitForEventSourceContext(ctx.WaitForEventSource(), s)

	if err != nil {
		return nil, errors.Wrap(err, "invalid event source")
	}

	waitForExp, err := expressions.NewWaitForEventExpression(
		v.getSourceMap(ctx),
		eventName,
		eventSource,
	)

	if err != nil {
		return nil, err
	}

	if optionsCtx := ctx.OptionsClause(); optionsCtx != nil {
		optionsExp, err := v.visitOptionsClause(optionsCtx, s)

		if err != nil {
			return nil, errors.Wrap(err, "invalid options")
		}

		if err := waitForExp.SetOptions(optionsExp); err != nil {
			return nil, err
		}
	}

	if filterCtx := ctx.FilterClause(); filterCtx != nil {
		nextScope := s.Fork(waitScope)

		if err := nextScope.SetVariable(waitPseudoVariable); err != nil {
			return nil, err
		}

		filterExp, err := v.visitFilterClause(filterCtx, nextScope)

		if err != nil {
			return nil, err
		}

		if err := waitForExp.SetFilter(v.getSourceMap(filterCtx), waitPseudoVariable, filterExp); err != nil {
			return nil, err
		}
	}

	if timeoutCtx := ctx.TimeoutClause(); timeoutCtx != nil {
		timeoutExp, err := v.visitTimeoutClause(timeoutCtx, s)

		if err != nil {
			return nil, errors.Wrap(err, "invalid timeout")
		}

		if err := waitForExp.SetTimeout(timeoutExp); err != nil {
			return nil, err
		}
	}

	return waitForExp, nil
}

func (v *visitor) visitMemberExpression(c fql.IMemberExpressionContext, scope *scope) (core.Expression, error) {
	ctx := c.(*fql.MemberExpressionContext)

	source, err := v.visitMemberExpressionSource(ctx.MemberExpressionSource(), scope)

	if err != nil {
		return nil, err
	}

	children := ctx.AllMemberExpressionPath()
	path := make([]*expressions.MemberPathSegment, 0, len(children))
	preCompiledPath := make([]core.Value, 0, len(children))
	skipOptimization := false

	for _, memberPath := range children {
		var exp core.Expression
		var err error

		memberPath := memberPath.(*fql.MemberExpressionPathContext)
		optional := memberPath.ErrorOperator() != nil

		if prop := memberPath.PropertyName(); prop != nil {
			exp, err = v.visitPropertyName(prop, scope)
		} else if prop := memberPath.ComputedPropertyName(); prop != nil {
			exp, err = v.visitComputedPropertyName(prop, scope)
		} else {
			return nil, v.unexpectedToken(memberPath)
		}

		if err != nil {
			return nil, err
		}

		segment, err := expressions.NewMemberPathSegment(exp, optional)

		if err != nil {
			return nil, err
		}

		if !skipOptimization {
			switch t := exp.(type) {
			case literals.StringLiteral:
				preCompiledPath = append(preCompiledPath, values.NewString(string(t)))
			case literals.IntLiteral:
				preCompiledPath = append(preCompiledPath, values.NewInt(int(t)))
			default:
				skipOptimization = true
				preCompiledPath = nil
			}
		}

		path = append(path, segment)
	}

	return expressions.NewMemberExpression(
		v.getSourceMap(ctx),
		source,
		path,
		preCompiledPath,
	)
}

func (v *visitor) visitMemberExpressionSource(c fql.IMemberExpressionSourceContext, scope *scope) (core.Expression, error) {
	ctx := c.(*fql.MemberExpressionSourceContext)

	if variable := ctx.Variable(); variable != nil {
		varName := variable.GetText()

		if strings.ToUpper(varName) == waitPseudoVariable {
			if scope.Name() == waitScope {
				varName = waitPseudoVariable
			}
		}

		if !scope.HasVariable(varName) {
			return nil, core.Error(ErrVariableNotFound, varName)
		}

		return expressions.NewVariableExpression(v.getSourceMap(ctx), varName)
	}

	if param := ctx.Param(); param != nil {
		return v.visitParam(param, scope)
	}

	if fnCall := ctx.FunctionCall(); fnCall != nil {
		return v.visitFunctionCall(fnCall, scope)
	}

	if objectLiteral := ctx.ObjectLiteral(); objectLiteral != nil {
		return v.visitObjectLiteral(objectLiteral, scope)
	}

	if arrayLiteral := ctx.ArrayLiteral(); arrayLiteral != nil {
		return v.visitArrayLiteral(arrayLiteral, scope)
	}

	return nil, v.unexpectedToken(ctx)
}

func (v *visitor) visitObjectLiteral(c fql.IObjectLiteralContext, scope *scope) (core.Expression, error) {
	ctx := c.(*fql.ObjectLiteralContext)
	assignments := ctx.AllPropertyAssignment()
	props := make([]*literals.ObjectPropertyAssignment, 0, len(assignments))

	for _, assignment := range assignments {
		var name core.Expression
		var value core.Expression
		var shortHand bool
		var err error

		pac := assignment.(*fql.PropertyAssignmentContext)

		if prop := pac.PropertyName(); prop != nil {
			name, err = v.visitPropertyName(prop, scope)
		} else if comProp := pac.ComputedPropertyName(); comProp != nil {
			name, err = v.visitComputedPropertyName(comProp, scope)
		} else if variable := pac.Variable(); variable != nil {
			shortHand = true
			name = literals.NewStringLiteral(variable.GetText())
			value, err = v.visitVariable(variable, scope)
		} else {
			return nil, v.unexpectedToken(pac)
		}

		if err != nil {
			return nil, err
		}

		if !shortHand {
			value, err = v.visitExpression(pac.Expression(), scope)

			if err != nil {
				return nil, err
			}
		}

		pa, err := literals.NewObjectPropertyAssignment(name, value)

		if err != nil {
			return nil, err
		}

		props = append(props, pa)
	}

	return literals.NewObjectLiteral(props), nil
}

func (v *visitor) visitPropertyName(c fql.IPropertyNameContext, scope *scope) (core.Expression, error) {
	ctx := c.(*fql.PropertyNameContext)

	if id := ctx.Identifier(); id != nil {
		return literals.NewStringLiteral(id.GetText()), nil
	}

	if rw := ctx.SafeReservedWord(); rw != nil {
		return literals.NewStringLiteral(rw.GetText()), nil
	}

	if rw := ctx.UnsafeReservedWord(); rw != nil {
		return literals.NewStringLiteral(rw.GetText()), nil
	}

	if stringLiteral := ctx.StringLiteral(); stringLiteral != nil {
		runes := []rune(stringLiteral.GetText())

		return literals.NewStringLiteral(string(runes[1 : len(runes)-1])), nil
	}

	if param := ctx.Param(); param != nil {
		return v.visitParam(param, scope)
	}

	return nil, v.unexpectedToken(ctx)
}

func (v *visitor) visitComputedPropertyName(c fql.IComputedPropertyNameContext, scope *scope) (core.Expression, error) {
	ctx := c.(*fql.ComputedPropertyNameContext)

	return v.visitExpression(ctx.Expression(), scope)
}

func (v *visitor) visitArrayLiteral(c fql.IArrayLiteralContext, scope *scope) (core.Expression, error) {
	ctx := c.(*fql.ArrayLiteralContext)
	list := ctx.ArgumentList()

	if list == nil {
		return literals.NewArrayLiteralWith(make([]core.Expression, 0, 0)), nil
	}

	elements, err := v.visitArgumentList(list, scope)

	if err != nil {
		return nil, err
	}

	return literals.NewArrayLiteralWith(elements), nil
}

func (v *visitor) visitFloatLiteral(ctx fql.IFloatLiteralContext) (core.Expression, error) {
	val, err := strconv.ParseFloat(ctx.GetText(), 64)
	if err != nil {
		return nil, err
	}

	return literals.NewFloatLiteral(val), nil
}

func (v *visitor) visitIntegerLiteral(ctx fql.IIntegerLiteralContext) (core.Expression, error) {
	val, err := strconv.Atoi(ctx.GetText())

	if err != nil {
		return nil, err
	}

	return literals.NewIntLiteral(val), nil
}

func (v *visitor) visitStringLiteral(ctx fql.IStringLiteralContext) (core.Expression, error) {
	var b strings.Builder

	for _, child := range ctx.GetChildren() {
		tree := child.(antlr.TerminalNode)
		sym := tree.GetSymbol()
		input := sym.GetInputStream()

		if input == nil {
			continue
		}

		size := input.Size()
		// skip quotes
		start := sym.GetStart() + 1
		stop := sym.GetStop() - 1

		if stop >= size {
			stop = size - 1
		}

		if start < size && stop < size {
			for i := start; i <= stop; i++ {
				c := input.GetText(i, i)

				switch c {
				case "\\":
					c2 := input.GetText(i, i+1)

					switch c2 {
					case "\\n":
						b.WriteString("\n")
					case "\\t":
						b.WriteString("\t")
					default:
						b.WriteString(c2)
					}

					i++
				default:
					b.WriteString(c)
				}
			}
		}
	}

	return literals.NewStringLiteral(b.String()), nil
}

func (v *visitor) visitBooleanLiteral(ctx fql.IBooleanLiteralContext) (core.Expression, error) {
	return literals.NewBooleanLiteral(strings.EqualFold(ctx.GetText(), "TRUE")), nil
}

func (v *visitor) visitNoneLiteral() (core.Expression, error) {
	return literals.None, nil
}

func (v *visitor) visitVariable(ctx fql.IVariableContext, scope *scope) (core.Expression, error) {
	name := ctx.GetText()

	// check whether the variable is defined
	if !scope.HasVariable(name) {
		if scope.Name() == waitScope && strings.ToUpper(name) == waitPseudoVariable {
			return expressions.NewVariableExpression(v.getSourceMap(ctx), waitPseudoVariable)
		}

		return nil, core.Error(ErrVariableNotFound, name)
	}

	return expressions.NewVariableExpression(v.getSourceMap(ctx), name)
}

func (v *visitor) visitVariableDeclaration(c fql.IVariableDeclarationContext, scope *scope) (core.Expression, error) {
	ctx := c.(*fql.VariableDeclarationContext)
	var init core.Expression
	var err error
	name := core.IgnorableVariable

	if id := ctx.Identifier(); id != nil {
		name = id.GetText()
	} else if reserved := ctx.SafeReservedWord(); reserved != nil {
		name = reserved.GetText()
	}

	err = scope.SetVariable(name)

	if err != nil {
		return nil, err
	}

	if exp := ctx.Expression(); exp != nil {
		init, err = v.visitExpression(ctx.Expression().(fql.IExpressionContext), scope)
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

func (v *visitor) visitRangeOperator(ctx fql.IRangeOperatorContext, scope *scope) (core.Expression, error) {
	left, err := v.visitRangeOperand(ctx.GetLeft(), scope)

	if err != nil {
		return nil, err
	}

	right, err := v.visitRangeOperand(ctx.GetRight(), scope)

	if err != nil {
		return nil, err
	}

	return operators.NewRangeOperator(
		v.getSourceMap(ctx),
		left,
		right,
	)
}

func (v *visitor) visitRangeOperand(c fql.IRangeOperandContext, scope *scope) (core.Expression, error) {
	ctx := c.(*fql.RangeOperandContext)

	if lit := ctx.IntegerLiteral(); lit != nil {
		return v.visitIntegerLiteral(lit)
	}

	if variable := ctx.Variable(); variable != nil {
		return v.visitVariable(variable, scope)
	}

	if param := ctx.Param(); param != nil {
		return v.visitParam(param, scope)
	}

	return nil, v.unexpectedToken(c)
}

func (v *visitor) visitFunctionCall(c fql.IFunctionCallContext, scope *scope) (core.Expression, error) {
	ctx := c.(*fql.FunctionCallContext)

	var args []core.Expression

	if arguments := ctx.ArgumentList(); arguments != nil {
		argList, err := v.visitArgumentList(arguments, scope)

		if err != nil {
			return nil, err
		}

		args = argList
	}

	var name string

	funcNS := ctx.Namespace()

	if funcNS != nil {
		name += funcNS.GetText()
	}

	name += ctx.FunctionName().GetText()

	fun, exists := v.funcs.Get(name)

	if !exists {
		return nil, core.Error(core.ErrNotFound, fmt.Sprintf("function: '%s'", name))
	}

	return expressions.NewFunctionCallExpression(
		v.getSourceMap(ctx),
		fun,
		args,
	)
}

func (v *visitor) visitFunctionCallExpression(c fql.IFunctionCallExpressionContext, scope *scope) (core.Expression, error) {
	ctx := c.(*fql.FunctionCallExpressionContext)

	exp, err := v.visitFunctionCall(
		ctx.FunctionCall(),
		scope,
	)

	if err != nil {
		return nil, err
	}

	if ctx.ErrorOperator() == nil {
		return exp, nil
	}

	return expressions.SuppressErrors(exp)
}

func (v *visitor) visitArgumentList(c fql.IArgumentListContext, scope *scope) ([]core.Expression, error) {
	ctx := c.(*fql.ArgumentListContext)
	exps := ctx.AllExpression()
	args := make([]core.Expression, 0, len(exps))

	for _, arg := range exps {
		exp, err := v.visitExpression(arg.(fql.IExpressionContext), scope)

		if err != nil {
			return nil, err
		}

		args = append(args, exp)
	}

	return args, nil
}

func (v *visitor) visitParam(c fql.IParamContext, scope *scope) (core.Expression, error) {
	ctx := c.(*fql.ParamContext)

	name := ctx.Identifier().GetText()

	scope.AddParam(name)

	return expressions.NewParameterExpression(
		v.getSourceMap(ctx),
		name,
	)
}

func (v *visitor) visitLiteral(c fql.ILiteralContext, scope *scope) (core.Expression, error) {
	ctx := c.(*fql.LiteralContext)

	if str := ctx.StringLiteral(); str != nil {
		return v.visitStringLiteral(str)
	}

	if i := ctx.IntegerLiteral(); i != nil {
		return v.visitIntegerLiteral(i)
	}

	if fl := ctx.FloatLiteral(); fl != nil {
		return v.visitFloatLiteral(fl)
	}

	if b := ctx.BooleanLiteral(); b != nil {
		return v.visitBooleanLiteral(b)
	}

	if arr := ctx.ArrayLiteral(); arr != nil {
		return v.visitArrayLiteral(arr, scope)
	}

	if obj := ctx.ObjectLiteral(); obj != nil {
		return v.visitObjectLiteral(obj, scope)
	}

	if none := ctx.NoneLiteral(); none != nil {
		return v.visitNoneLiteral()
	}

	return nil, v.unexpectedToken(ctx)
}

func (v *visitor) visitExpression(c fql.IExpressionContext, scope *scope) (core.Expression, error) {
	ctx, ok := c.(*fql.ExpressionContext)

	if !ok {
		return nil, v.invalidToken(ctx)
	}

	if op := ctx.UnaryOperator(); op != nil {
		exp, err := v.visitExpression(ctx.GetRight(), scope)

		if err != nil {
			return nil, err
		}

		return operators.NewUnaryOperator(
			v.getSourceMap(ctx),
			exp,
			op.GetText(),
		)
	}

	if ctx.GetTernaryOperator() != nil {
		cond, err := v.visitExpression(ctx.GetCondition().(fql.IExpressionContext), scope)

		if err != nil {
			return nil, err
		}

		var consequent core.Expression

		if onTrue := ctx.GetOnTrue(); onTrue != nil {
			exp, err := v.visitExpression(onTrue.(fql.IExpressionContext), scope)

			if err != nil {
				return nil, err
			}

			consequent = exp
		}

		alternate, err := v.visitExpression(ctx.GetOnFalse().(fql.IExpressionContext), scope)

		if err != nil {
			return nil, err
		}

		return expressions.NewConditionExpression(
			v.getSourceMap(ctx),
			cond,
			consequent,
			alternate,
		)
	}

	var logical fql.IExpressionContext
	var logicalOp string

	if andOp := ctx.LogicalAndOperator(); andOp != nil {
		logical = ctx
		logicalOp = andOp.GetText()
	} else if orOp := ctx.LogicalOrOperator(); orOp != nil {
		logical = ctx
		logicalOp = orOp.GetText()
	}

	if logical != nil {
		left, err := v.visitExpression(ctx.GetLeft(), scope)

		if err != nil {
			return nil, err
		}

		right, err := v.visitExpression(ctx.GetRight(), scope)

		if err != nil {
			return nil, err
		}

		return operators.NewLogicalOperator(v.getSourceMap(ctx), left, right, logicalOp)
	}

	if pred := ctx.Predicate(); pred != nil {
		return v.visitPredicate(pred, scope)
	}

	return nil, ErrNotImplemented
}

func (v *visitor) visitPredicate(c fql.IPredicateContext, scope *scope) (core.Expression, error) {
	ctx, ok := c.(*fql.PredicateContext)

	if !ok {
		return nil, v.invalidToken(ctx)
	}

	if ea := ctx.ExpressionAtom(); ea != nil {
		return v.visitExpressionAtom(ea, scope)
	}

	src := v.getSourceMap(c)

	left, err := v.visitPredicate(c.GetLeft(), scope)

	if err != nil {
		return nil, err
	}

	right, err := v.visitPredicate(c.GetRight(), scope)

	if err != nil {
		return nil, err
	}

	if op := ctx.EqualityOperator(); op != nil {
		return operators.NewEqualityOperator(src, left, right, op.GetText())
	}

	if op := ctx.ArrayOperator(); op != nil {
		op := op.(*fql.ArrayOperatorContext)

		var comparator core.Predicate

		if eqOp := op.EqualityOperator(); eqOp != nil {
			op, err := operators.NewEqualityOperator(
				v.getSourceMap(eqOp),
				left,
				right,
				eqOp.GetText(),
			)
			if err != nil {
				return nil, err
			}

			comparator = op
		} else if inOp := op.InOperator(); inOp != nil {
			inOp := inOp.(*fql.InOperatorContext)

			op, err := operators.NewInOperator(
				v.getSourceMap(inOp),
				left,
				right,
				inOp.Not() != nil,
			)

			if err != nil {
				return nil, err
			}

			comparator = op
		}

		if comparator == nil {
			return nil, v.invalidToken(c)
		}

		return operators.NewArrayOperator(
			v.getSourceMap(c),
			left,
			right,
			op.GetOperator().GetText(),
			comparator,
		)
	}

	if op := ctx.InOperator(); op != nil {
		op := op.(*fql.InOperatorContext)

		return operators.NewInOperator(
			v.getSourceMap(op),
			left,
			right,
			op.Not() != nil,
		)
	}

	if op := ctx.LikeOperator(); op != nil {
		op := op.(*fql.LikeOperatorContext)

		return operators.NewLikeOperator(v.getSourceMap(op), left, right, op.Not() != nil)
	}

	return nil, v.invalidToken(c)
}

func (v *visitor) visitExpressionAtom(c fql.IExpressionAtomContext, scope *scope) (core.Expression, error) {
	ctx, ok := c.(*fql.ExpressionAtomContext)

	if !ok {
		return nil, core.Error(ErrInvalidToken, c.GetText())
	}

	src := v.getSourceMap(c)

	if literal := ctx.Literal(); literal != nil {
		return v.visitLiteral(literal, scope)
	}

	if variable := ctx.Variable(); variable != nil {
		return v.visitVariable(variable, scope)
	}

	if param := ctx.Param(); param != nil {
		return v.visitParam(param, scope)
	}

	if member := ctx.MemberExpression(); member != nil {
		return v.visitMemberExpression(member, scope)
	}

	if fc := ctx.FunctionCallExpression(); fc != nil {
		fc := fc.(*fql.FunctionCallExpressionContext)

		exp, err := v.visitFunctionCall(fc.FunctionCall(), scope)

		if err != nil {
			return nil, err
		}

		if fc.ErrorOperator() == nil {
			return exp, nil
		}

		return expressions.SuppressErrors(exp)
	}

	if rop := ctx.RangeOperator(); rop != nil {
		return v.visitRangeOperator(rop, scope)
	}

	var exp core.Expression
	var err error

	if forIn := ctx.ForExpression(); forIn != nil {
		exp, err = v.visitForExpression(forIn, scope)
	} else if waitFor := ctx.WaitForExpression(); waitFor != nil {
		exp, err = v.visitWaitForExpression(waitFor, scope)
	} else if e := ctx.Expression(); e != nil {
		exp, err = v.visitExpression(e, scope)
	}

	if err != nil {
		return nil, err
	}

	if exp != nil {
		if ctx.ErrorOperator() == nil {
			return exp, nil
		}

		return expressions.SuppressErrors(exp)
	}

	leftExp := ctx.GetLeft()
	left, err := v.visitExpressionAtom(leftExp, scope)

	if err != nil {
		return nil, err
	}

	rightExp := ctx.GetRight()
	right, err := v.visitExpressionAtom(rightExp, scope)

	if err != nil {
		return nil, err
	}

	var math fql.IExpressionAtomContext
	var mathOp string

	if op := ctx.MultiplicativeOperator(); op != nil {
		math = ctx
		mathOp = op.GetText()
	} else if op := ctx.AdditiveOperator(); op != nil {
		math = ctx
		mathOp = op.GetText()
	}

	if math != nil {
		return operators.NewMathOperator(src, left, right, mathOp)
	}

	if op := ctx.RegexpOperator(); op != nil {
		switch lit := right.(type) {
		case literals.StringLiteral:
			_, err := regexp.Compile(string(lit))
			if err != nil {
				src := v.getSourceMap(rightExp)

				return nil, errors.Wrap(err, src.String())
			}
		case *literals.ArrayLiteral, *literals.ObjectLiteral, literals.BooleanLiteral, literals.FloatLiteral, literals.IntLiteral:
			src := v.getSourceMap(rightExp)

			return nil, errors.Wrap(errors.New("expected a string literal or a function call"), src.String())
		}

		return operators.NewRegexpOperator(src, left, right, op.GetText())
	}

	return nil, v.unexpectedToken(ctx)
}

func (v *visitor) unexpectedToken(tree antlr.ParseTree) error {
	return core.Error(ErrUnexpectedToken, tree.GetText())
}

func (v *visitor) invalidToken(tree antlr.ParseTree) error {
	return core.Error(ErrInvalidToken, tree.GetText())
}

func (v *visitor) getSourceMap(rule antlr.ParserRuleContext) core.SourceMap {
	start := rule.GetStart()

	return core.NewSourceMap(
		rule.GetText(),
		start.GetLine(),
		start.GetColumn(),
	)
}
