// Code generated from EventRule.g4 by ANTLR 4.7.2. DO NOT EDIT.

package parser // EventRule

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type BaseEventRuleVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseEventRuleVisitor) VisitStart(ctx *StartContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEventRuleVisitor) VisitNot(ctx *NotContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEventRuleVisitor) VisitParenthesis(ctx *ParenthesisContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEventRuleVisitor) VisitVariable(ctx *VariableContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEventRuleVisitor) VisitCompareNumber(ctx *CompareNumberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEventRuleVisitor) VisitStringEqualContains(ctx *StringEqualContainsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEventRuleVisitor) VisitAndOr(ctx *AndOrContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEventRuleVisitor) VisitStringIn(ctx *StringInContext) interface{} {
	return v.VisitChildren(ctx)
}
