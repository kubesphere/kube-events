package visitor

import (
	"strconv"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/kubesphere/kube-events/pkg/ruler/visitor/parser"
	"k8s.io/klog"
)

type Visitor struct {
	parser.BaseEventRuleVisitor
	valueStack []bool
	m          map[string]interface{}
}

func NewVisitor(m map[string]interface{}) *Visitor {
	return &Visitor{
		m: m,
	}
}

func (l *Visitor) pushValue(i bool) {
	l.valueStack = append(l.valueStack, i)
}

func (l *Visitor) popValue() bool {
	if len(l.valueStack) < 1 {
		panic("valueStack is empty unable to pop")
	}

	// Get the last value from the stack.
	result := l.valueStack[len(l.valueStack)-1]

	// Remove the last element from the stack.
	l.valueStack = l.valueStack[:len(l.valueStack)-1]

	return result
}

func (v *Visitor) visitRule(node antlr.RuleNode) interface{} {
	node.Accept(v)
	return nil
}

func (v *Visitor) VisitStart(ctx *parser.StartContext) interface{} {
	return v.visitRule(ctx.Expression())
}

func (v *Visitor) VisitAndOr(ctx *parser.AndOrContext) interface{} {
	klog.V(3).Info("VisitAndOr")

	//push expression result to stack
	v.visitRule(ctx.Expression(0))
	v.visitRule(ctx.Expression(1))

	//push result to stack
	var t = ctx.GetOp()
	right := v.popValue()
	left := v.popValue()
	switch t.GetTokenType() {
	case parser.EventRuleParserAND:
		v.pushValue(left && right)
	case parser.EventRuleParserOR:
		v.pushValue(left || right)
	default:
		panic("should not happen")
	}

	return nil
}

func (v *Visitor) VisitNot(ctx *parser.NotContext) interface{} {
	klog.V(3).Info("VisitNot")

	v.visitRule(ctx.Expression())

	value := v.popValue()
	v.pushValue(!value)

	return nil
}

func (v *Visitor) VisitStringEqualContains(ctx *parser.StringEqualContainsContext) interface{} {
	varName := ctx.VAR().GetText()
	strValue := ctx.STRING().GetText()
	var t = ctx.GetOp()

	strValue = strings.TrimLeft(strValue, `"`)
	strValue = strings.TrimRight(strValue, `"`)

	klog.V(3).Infof("VisitStringEqualContains %s %d %s\n", varName, t.GetTokenType(), strValue)

	varValue := ""
	if vv, y := v.m[varName]; y {
		varValue = vv.(string)
	}

	switch t.GetTokenType() {
	case parser.EventRuleParserEQU:
		v.pushValue(varValue == strValue)
	case parser.EventRuleParserCONTAINS:
		v.pushValue(strings.Contains(varValue, strValue))
	}

	return nil
}

func (v *Visitor) VisitStringIn(ctx *parser.StringInContext) interface{} {
	varName := ctx.VAR().GetText()
	length := len(ctx.AllSTRING())

	varValue := ""
	if vv, y := v.m[varName]; y {
		varValue = vv.(string)
	}

	strValues := []string{}
	for i := 0; i < length; i++ {
		strValue := ctx.STRING(i).GetText()
		strValue = strings.TrimLeft(strValue, `"`)
		strValue = strings.TrimRight(strValue, `"`)
		strValues = append(strValues, strValue)
	}

	klog.V(3).Infof("VisitStringIn %s in %v\n", varName, strValues)

	result := false
	for _, strValue := range strValues {
		if varValue == strValue {
			result = true
			break
		}
	}

	v.pushValue(result)

	return nil
}

func (v *Visitor) VisitCompareNumber(ctx *parser.CompareNumberContext) interface{} {
	varName := ctx.VAR().GetText()
	numValue, err := strconv.ParseFloat(ctx.NUMBER().GetText(), 64)
	if err != nil {
		klog.Error(err)
	}

	result := false
	if vv, y := v.m[varName]; y {
		varValue := vv.(float64)
		var t = ctx.GetOp()
		klog.V(3).Infof("VisitCompareNumber %s %d %v\n", varName, t.GetTokenType(), numValue)
		switch t.GetTokenType() {
		case parser.EventRuleParserEQU:
			result = varValue == numValue
		case parser.EventRuleParserNEQ:
			result = varValue != numValue
		case parser.EventRuleParserGT:
			result = varValue > numValue
		case parser.EventRuleParserLT:
			result = varValue < numValue
		case parser.EventRuleParserGTE:
			result = varValue >= numValue
		case parser.EventRuleParserLTE:
			result = varValue <= numValue
		}
	}
	v.pushValue(result)

	return nil
}

func (v *Visitor) VisitVariable(ctx *parser.VariableContext) interface{} {
	varName := ctx.VAR().GetText()
	klog.V(3).Infof("VisitVariable %v\n", varName)

	result := false
	if varValue, y := v.m[varName]; y {
		result = varValue.(bool)
	}
	v.pushValue(result)

	return nil
}

func (v *Visitor) VisitParenthesis(ctx *parser.ParenthesisContext) interface{} {
	v.visitRule(ctx.Expression())
	return nil
}

func EventRuleEvaluate(m map[string]interface{}, expression string) bool {
	is := antlr.NewInputStream(expression)

	// Create the Lexer
	lexer := parser.NewEventRuleLexer(is)
	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create the Parser
	p := parser.NewEventRuleParser(tokens)

	v := NewVisitor(m)
	//Start is rules name of EventRule.g4
	p.Start().Accept(v)
	return v.popValue()
}
