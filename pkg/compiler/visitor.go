package compiler

import (
	"fmt"
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
	"strconv"
	"strings"
)

type (
	forOption func(f *expressions.ForExpression) error

	visitor struct {
		*fql.BaseFqlParserVisitor
		src   string
		funcs map[string]core.Function
	}
)

func newVisitor(src string, funcs map[string]core.Function) *visitor {
	return &visitor{
		&fql.BaseFqlParserVisitor{},
		src,
		funcs,
	}
}

func (v *visitor) VisitProgram(ctx *fql.ProgramContext) interface{} {
	return newResultFrom(func() (interface{}, error) {
		rootScope := newRootScope()
		block, err := v.doVisitBody(ctx.Body().(*fql.BodyContext), rootScope)

		if err != nil {
			return nil, err
		}

		return runtime.NewProgram(v.src, block)
	})
}

func (v *visitor) doVisitBody(ctx *fql.BodyContext, scope *scope) (core.Expression, error) {
	statements := ctx.AllBodyStatement()
	body := expressions.NewBlockExpression(len(statements) + 1)

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

	return nil, errors.Wrap(ErrInvalidToken, ctx.GetText())
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
	valVar := ctx.ForExpressionValueVariable()
	valVarName = valVar.GetText()
	forInScope := scope.Fork()
	forInScope.SetVariable(valVarName)

	keyVar := ctx.ForExpressionKeyVariable()

	if keyVar != nil {
		keyVarName = keyVar.GetText()
		forInScope.SetVariable(keyVarName)
	}

	srcCtx := ctx.ForExpressionSource().(*fql.ForExpressionSourceContext)
	srcExp, err := v.doVisitForExpressionSource(srcCtx, forInScope)

	if err != nil {
		return nil, err
	}

	src, err := expressions.NewDataSource(
		v.getSourceMap(srcCtx),
		valVarName,
		keyVarName,
		srcExp,
	)

	// Clauses.
	// We put clauses parsing before parsing the query body because COLLECT clause overrides scope variables
	for _, clause := range ctx.AllForExpressionClause() {
		clause := clause.(*fql.ForExpressionClauseContext)

		limitCtx := clause.LimitClause()

		if limitCtx != nil {
			limit, offset, err := v.createLimit(limitCtx.(*fql.LimitClauseContext))

			if err != nil {
				return nil, err
			}

			parsedClauses = append(parsedClauses, func(f *expressions.ForExpression) error {
				return f.AddLimit(v.getSourceMap(limitCtx), limit, offset)
			})

			continue
		}

		filterCtx := clause.FilterClause()

		if filterCtx != nil {
			filterExp, err := v.createFilter(filterCtx.(*fql.FilterClauseContext), forInScope)

			if err != nil {
				return nil, err
			}

			parsedClauses = append(parsedClauses, func(f *expressions.ForExpression) error {
				return f.AddFilter(v.getSourceMap(filterCtx), filterExp)
			})

			continue
		}

		sortCtx := clause.SortClause()

		if sortCtx != nil {
			sortCtx := sortCtx.(*fql.SortClauseContext)
			sortExps, err := v.createSort(sortCtx, forInScope)

			if err != nil {
				return nil, err
			}

			parsedClauses = append(parsedClauses, func(f *expressions.ForExpression) error {
				return f.AddSort(v.getSourceMap(sortCtx), sortExps...)
			})
		}

		collectCtx := clause.CollectClause()

		if collectCtx != nil {
			collectCtx := collectCtx.(*fql.CollectClauseContext)

			params, err := v.createCollect(collectCtx, forInScope)

			if err != nil {
				return nil, err
			}

			forInScope.RemoveVariable(valVarName)

			if keyVarName != "" {
				forInScope.RemoveVariable(keyVarName)
			}

			parsedClauses = append(parsedClauses, func(f *expressions.ForExpression) error {
				return f.AddCollect(v.getSourceMap(collectCtx), params)
			})
		}
	}

	body := ctx.AllForExpressionBody()
	predicate := expressions.NewBlockExpression(len(body) + 1)

	for _, el := range body {
		el, err := v.doVisitForExpressionBody(el.(*fql.ForExpressionBodyContext), forInScope)

		if err != nil {
			return nil, err
		}

		err = predicate.Add(el)

		if err != nil {
			return nil, err
		}
	}

	var spread bool
	var distinct bool
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

		predicate.Add(returnExp)
	} else {
		forInCtx := forRetCtx.ForExpression().(*fql.ForExpressionContext)
		forInExp, err := v.doVisitForExpression(forInCtx, forInScope)

		if err != nil {
			return nil, err
		}

		spread = true

		predicate.Add(forInExp)
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

func (v *visitor) createLimit(ctx *fql.LimitClauseContext) (int, int, error) {
	var err error
	var count int
	var offset int

	intLiterals := ctx.AllIntegerLiteral()

	if len(intLiterals) > 1 {
		offset, err = v.parseInt(intLiterals[0])

		if err != nil {
			return 0, 0, err
		}

		count, err = v.parseInt(intLiterals[1])

		if err != nil {
			return 0, 0, err
		}
	} else {
		count, err = strconv.Atoi(intLiterals[0].GetText())

		if err != nil {
			return 0, 0, err
		}
	}

	return count, offset, nil
}

func (v *visitor) parseInt(node antlr.TerminalNode) (int, error) {
	return strconv.Atoi(node.GetText())
}

func (v *visitor) createFilter(ctx *fql.FilterClauseContext, scope *scope) (core.Expression, error) {
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

		logicalOp := exp.LogicalOperator()

		if logicalOp != nil {
			return operators.NewLogicalOperator(v.getSourceMap(ctx), left, right, logicalOp.GetText())
		}
	} else {
		// should be unary operator
		return v.doVisitExpression(exp, scope)
	}

	return nil, core.Error(ErrInvalidToken, ctx.GetText())
}

func (v *visitor) createSort(ctx *fql.SortClauseContext, scope *scope) ([]*clauses.SorterExpression, error) {
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

func (v *visitor) createCollect(ctx *fql.CollectClauseContext, scope *scope) (clauses.CollectParams, error) {
	params := clauses.CollectParams{}

	groupingCtx := ctx.CollectGrouping()

	if groupingCtx != nil {
		groupingCtx := groupingCtx.(*fql.CollectGroupingContext)
		collectSelectors := groupingCtx.AllCollectSelector()

		grouping := make([]*clauses.CollectSelector, 0, len(collectSelectors))

		for _, cs := range collectSelectors {
			selector, err := v.createCollectSelector(cs.(*fql.CollectSelectorContext), scope)

			if err != nil {
				return params, err
			}

			grouping = append(grouping, selector)

			if err := scope.SetVariable(selector.Variable()); err != nil {
				return params, err
			}
		}

		params.Grouping = grouping
	}

	return params, nil
}

func (v *visitor) createCollectSelector(ctx *fql.CollectSelectorContext, scope *scope) (*clauses.CollectSelector, error) {
	variable := ctx.Identifier().GetText()
	exp, err := v.doVisitExpression(ctx.Expression().(*fql.ExpressionContext), scope)

	if err != nil {
		return nil, err
	}

	return clauses.NewCollectSelector(variable, exp)
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

func (v *visitor) doVisitForExpressionBody(ctx *fql.ForExpressionBodyContext, scope *scope) (core.Expression, error) {
	varDecCtx := ctx.VariableDeclaration()

	if varDecCtx != nil {
		return v.doVisitVariableDeclaration(varDecCtx.(*fql.VariableDeclarationContext), scope)
	}

	funcCallCtx := ctx.FunctionCallExpression()

	if funcCallCtx != nil {
		return v.doVisitFunctionCallExpression(funcCallCtx.(*fql.FunctionCallExpressionContext), scope)
	}

	return nil, v.unexpectedToken(ctx)
}

func (v *visitor) doVisitMemberExpression(ctx *fql.MemberExpressionContext, scope *scope) (core.Expression, error) {
	varName := ctx.Identifier().GetText()

	_, err := scope.GetVariable(varName)

	if err != nil {
		return nil, err
	}

	children := ctx.GetChildren()
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

	member, err := expressions.NewMemberExpression(
		v.getSourceMap(ctx),
		varName,
		path,
	)

	if err != nil {
		return nil, err
	}

	return member, nil
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

		if prop != nil {
			name, err = v.doVisitPropertyNameContext(prop.(*fql.PropertyNameContext), scope)
		} else if computedProp != nil {
			name, err = v.doVisitComputedPropertyNameContext(computedProp.(*fql.ComputedPropertyNameContext), scope)
		} else {
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

		props = append(props, literals.NewObjectPropertyAssignment(name, value))
	}

	return literals.NewObjectLiteralWith(props...), nil
}

func (v *visitor) doVisitPropertyNameContext(ctx *fql.PropertyNameContext, _ *scope) (core.Expression, error) {
	var name string

	identifier := ctx.Identifier()

	if identifier != nil {
		name = identifier.GetText()
	} else {
		stringLiteral := ctx.StringLiteral()

		if stringLiteral != nil {
			runes := []rune(stringLiteral.GetText())
			name = string(runes[1 : len(runes)-1])
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

	_, err := scope.GetVariable(name)

	if err != nil {
		return nil, err
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
	text := ctx.StringLiteral().GetText()

	// remove extra quotes
	return literals.NewStringLiteral(text[1 : len(text)-1]), nil
}

func (v *visitor) doVisitBooleanLiteral(ctx *fql.BooleanLiteralContext) (core.Expression, error) {
	return literals.NewBooleanLiteral(strings.ToUpper(ctx.GetText()) == "TRUE"), nil
}

func (v *visitor) doVisitNoneLiteral(_ *fql.NoneLiteralContext) (core.Expression, error) {
	return literals.None, nil
}

func (v *visitor) doVisitVariable(ctx *fql.VariableContext, scope *scope) (core.Expression, error) {
	name := ctx.Identifier().GetText()

	// check whether the variable is defined
	_, err := scope.GetVariable(name)

	if err != nil {
		return nil, err
	}

	return expressions.NewVariableExpression(v.getSourceMap(ctx), name)
}

func (v *visitor) doVisitVariableDeclaration(ctx *fql.VariableDeclarationContext, scope *scope) (core.Expression, error) {
	var init core.Expression
	var err error

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

	name := ctx.Identifier().GetText()

	scope.SetVariable(name)

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
		return nil, errors.Wrap(ErrInvalidToken, ctx.GetText())
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

	funcName := context.Identifier().GetText()

	fun, exists := v.funcs[funcName]

	if !exists {
		return nil, core.Error(core.ErrNotFound, fmt.Sprintf("function: '%s'", funcName))
	}

	return expressions.NewFunctionCallExpression(
		v.getSourceMap(context),
		fun,
		args...,
	)
}

func (v *visitor) doVisitParamContext(context *fql.ParamContext, _ *scope) (core.Expression, error) {
	name := context.Identifier().GetText()

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
	mathOp := ctx.MathOperator().(*fql.MathOperatorContext)
	exps, err := v.doVisitAllExpressions(ctx.AllExpression(), scope)

	if err != nil {
		return nil, err
	}

	left := exps[0]
	right := exps[1]

	return operators.NewMathOperator(
		v.getSourceMap(mathOp),
		left,
		right,
		operators.MathOperatorType(mathOp.GetText()),
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
	logicalOp := ctx.LogicalOperator().(*fql.LogicalOperatorContext)
	exps, err := v.doVisitAllExpressions(ctx.AllExpression(), scope)

	if err != nil {
		return nil, err
	}

	left := exps[0]
	right := exps[1]

	return operators.NewLogicalOperator(v.getSourceMap(logicalOp), left, right, logicalOp.GetText())
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

	if ctx.InOperator() != nil {
		comparator, err = v.doVisitInOperator(ctx, scope)
	} else if ctx.EqualityOperator() != nil {
		comparator, err = v.doVisitEqualityOperator(ctx, scope)
	} else {
		return nil, v.unexpectedToken(ctx)
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

func (v *visitor) doVisitExpression(ctx *fql.ExpressionContext, scope *scope) (core.Expression, error) {
	notOp := ctx.UnaryOperator()

	if notOp != nil {
		return v.doVisitUnaryOperator(ctx, scope)
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

	funCall := ctx.FunctionCallExpression()

	if funCall != nil {
		return v.doVisitFunctionCallExpression(funCall.(*fql.FunctionCallExpressionContext), scope)
	}

	member := ctx.MemberExpression()

	if member != nil {
		return v.doVisitMemberExpression(member.(*fql.MemberExpressionContext), scope)
	}

	none := ctx.NoneLiteral()

	if none != nil {
		return v.doVisitNoneLiteral(none.(*fql.NoneLiteralContext))
	}

	arrOp := ctx.ArrayOperator()

	if arrOp != nil {
		return v.doVisitArrayOperator(ctx, scope)
	}

	inOp := ctx.InOperator()

	if inOp != nil {
		return v.doVisitInOperator(ctx, scope)
	}

	equalityOp := ctx.EqualityOperator()

	if equalityOp != nil {
		return v.doVisitEqualityOperator(ctx, scope)
	}

	logicalOp := ctx.LogicalOperator()

	if logicalOp != nil {
		return v.doVisitLogicalOperator(ctx, scope)
	}

	mathOp := ctx.MathOperator()

	if mathOp != nil {
		return v.doVisitMathOperator(ctx, scope)
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
		return make([]core.Expression, 0, 0), nil
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

	switch node.(type) {
	case *fql.BodyContext:
		out, err = v.doVisitBody(node.(*fql.BodyContext), scope)
	case *fql.ExpressionContext:
		out, err = v.doVisitExpression(node.(*fql.ExpressionContext), scope)
	case *fql.ForExpressionContext:
		out, err = v.doVisitForExpression(node.(*fql.ForExpressionContext), scope)
	case *fql.ReturnExpressionContext:
		out, err = v.doVisitReturnExpression(node.(*fql.ReturnExpressionContext), scope)
	case *fql.ArrayLiteralContext:
		out, err = v.doVisitArrayLiteral(node.(*fql.ArrayLiteralContext), scope)
	case *fql.ObjectLiteralContext:
		out, err = v.doVisitObjectLiteral(node.(*fql.ObjectLiteralContext), scope)
	case *fql.StringLiteralContext:
		out, err = v.doVisitStringLiteral(node.(*fql.StringLiteralContext))
	case *fql.IntegerLiteralContext:
		out, err = v.doVisitIntegerLiteral(node.(*fql.IntegerLiteralContext))
	case *fql.FloatLiteralContext:
		out, err = v.doVisitFloatLiteral(node.(*fql.FloatLiteralContext))
	case *fql.BooleanLiteralContext:
		out, err = v.doVisitBooleanLiteral(node.(*fql.BooleanLiteralContext))
	case *fql.NoneLiteralContext:
		out, err = v.doVisitNoneLiteral(node.(*fql.NoneLiteralContext))
	case *fql.VariableContext:
		out, err = v.doVisitVariable(node.(*fql.VariableContext), scope)
	case *fql.VariableDeclarationContext:
		out, err = v.doVisitVariableDeclaration(node.(*fql.VariableDeclarationContext), scope)
	case *fql.FunctionCallExpressionContext:
		out, err = v.doVisitFunctionCallExpression(node.(*fql.FunctionCallExpressionContext), scope)
	case *fql.ParamContext:
		out, err = v.doVisitParamContext(node.(*fql.ParamContext), scope)
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
