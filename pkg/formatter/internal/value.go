package internal

import "github.com/MontFerret/ferret/v2/pkg/parser/fql"

type valueFormatter struct {
	*engine
}

func (v *valueFormatter) formatStringOrRef(
	stringCtx fql.IStringLiteralContext,
	variable fql.IVariableContext,
	param fql.IParamContext,
	member fql.IMemberExpressionContext,
	call fql.IFunctionCallContext,
) {
	if stringCtx != nil {
		v.literal.formatStringLiteralNode(stringCtx)
		return
	}
	v.formatRefValue(variable, param, member, call)
}

func (v *valueFormatter) formatRefValue(
	variable fql.IVariableContext,
	param fql.IParamContext,
	member fql.IMemberExpressionContext,
	call fql.IFunctionCallContext,
) {
	switch {
	case variable != nil:
		v.expression.formatVariable(variable.(*fql.VariableContext))
	case param != nil:
		v.expression.formatParam(param.(*fql.ParamContext))
	case member != nil:
		v.member.formatMemberExpression(member.(*fql.MemberExpressionContext))
	case call != nil:
		v.expression.formatFunctionCall(call.(*fql.FunctionCallContext))
	}
}

func (v *valueFormatter) formatRefValueWithCallExpr(
	call fql.IFunctionCallExpressionContext,
	variable fql.IVariableContext,
	param fql.IParamContext,
	member fql.IMemberExpressionContext,
) {
	switch {
	case call != nil:
		v.expression.formatFunctionCallExpression(call.(*fql.FunctionCallExpressionContext))
	case variable != nil:
		v.expression.formatVariable(variable.(*fql.VariableContext))
	case param != nil:
		v.expression.formatParam(param.(*fql.ParamContext))
	case member != nil:
		v.member.formatMemberExpression(member.(*fql.MemberExpressionContext))
	}
}

func (v *valueFormatter) formatDurationNumberOrRef(
	duration fql.IDurationLiteralContext,
	integer fql.IIntegerLiteralContext,
	float fql.IFloatLiteralContext,
	variable fql.IVariableContext,
	param fql.IParamContext,
	member fql.IMemberExpressionContext,
	call fql.IFunctionCallContext,
) {
	switch {
	case duration != nil:
		v.p.write(duration.GetText())
	case integer != nil:
		v.p.write(integer.GetText())
	case float != nil:
		v.p.write(float.GetText())
	default:
		v.formatRefValue(variable, param, member, call)
	}
}

func (v *valueFormatter) formatFloatOrIntOrRef(
	float fql.IFloatLiteralContext,
	integer fql.IIntegerLiteralContext,
	variable fql.IVariableContext,
	param fql.IParamContext,
	member fql.IMemberExpressionContext,
	call fql.IFunctionCallContext,
) {
	switch {
	case float != nil:
		v.p.write(float.GetText())
	case integer != nil:
		v.p.write(integer.GetText())
	default:
		v.formatRefValue(variable, param, member, call)
	}
}

func (v *valueFormatter) formatIntOrRef(
	integer fql.IIntegerLiteralContext,
	param fql.IParamContext,
	variable fql.IVariableContext,
	call fql.IFunctionCallExpressionContext,
	member fql.IMemberExpressionContext,
) {
	switch {
	case integer != nil:
		v.p.write(integer.GetText())
	case param != nil:
		v.expression.formatParam(param.(*fql.ParamContext))
	case variable != nil:
		v.expression.formatVariable(variable.(*fql.VariableContext))
	case call != nil:
		v.expression.formatFunctionCallExpression(call.(*fql.FunctionCallExpressionContext))
	case member != nil:
		v.member.formatMemberExpression(member.(*fql.MemberExpressionContext))
	}
}
