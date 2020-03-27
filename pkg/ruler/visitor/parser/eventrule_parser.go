// Code generated from EventRule.g4 by ANTLR 4.7.2. DO NOT EDIT.

package parser // EventRule

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = reflect.Copy
var _ = strconv.Itoa

var parserATN = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 20, 46, 4,
	2, 9, 2, 4, 3, 9, 3, 3, 2, 3, 2, 3, 2, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 7, 3, 26,
	10, 3, 12, 3, 14, 3, 29, 11, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 5, 3, 36,
	10, 3, 3, 3, 3, 3, 3, 3, 7, 3, 41, 10, 3, 12, 3, 14, 3, 44, 11, 3, 3, 3,
	2, 3, 4, 4, 2, 4, 2, 5, 4, 2, 8, 8, 14, 14, 3, 2, 8, 13, 3, 2, 5, 6, 2,
	50, 2, 6, 3, 2, 2, 2, 4, 35, 3, 2, 2, 2, 6, 7, 5, 4, 3, 2, 7, 8, 7, 2,
	2, 3, 8, 3, 3, 2, 2, 2, 9, 10, 8, 3, 1, 2, 10, 11, 7, 7, 2, 2, 11, 36,
	5, 4, 3, 8, 12, 13, 7, 3, 2, 2, 13, 14, 5, 4, 3, 2, 14, 15, 7, 4, 2, 2,
	15, 36, 3, 2, 2, 2, 16, 17, 7, 18, 2, 2, 17, 18, 9, 2, 2, 2, 18, 36, 7,
	19, 2, 2, 19, 20, 7, 18, 2, 2, 20, 21, 7, 15, 2, 2, 21, 22, 7, 3, 2, 2,
	22, 27, 7, 19, 2, 2, 23, 24, 7, 16, 2, 2, 24, 26, 7, 19, 2, 2, 25, 23,
	3, 2, 2, 2, 26, 29, 3, 2, 2, 2, 27, 25, 3, 2, 2, 2, 27, 28, 3, 2, 2, 2,
	28, 30, 3, 2, 2, 2, 29, 27, 3, 2, 2, 2, 30, 36, 7, 4, 2, 2, 31, 32, 7,
	18, 2, 2, 32, 33, 9, 3, 2, 2, 33, 36, 7, 17, 2, 2, 34, 36, 7, 18, 2, 2,
	35, 9, 3, 2, 2, 2, 35, 12, 3, 2, 2, 2, 35, 16, 3, 2, 2, 2, 35, 19, 3, 2,
	2, 2, 35, 31, 3, 2, 2, 2, 35, 34, 3, 2, 2, 2, 36, 42, 3, 2, 2, 2, 37, 38,
	12, 9, 2, 2, 38, 39, 9, 4, 2, 2, 39, 41, 5, 4, 3, 10, 40, 37, 3, 2, 2,
	2, 41, 44, 3, 2, 2, 2, 42, 40, 3, 2, 2, 2, 42, 43, 3, 2, 2, 2, 43, 5, 3,
	2, 2, 2, 44, 42, 3, 2, 2, 2, 5, 27, 35, 42,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "'('", "')'", "'and'", "'or'", "'not'", "'='", "'!='", "'>'", "'<'",
	"'>='", "'<='", "'contains'", "'in'", "','",
}
var symbolicNames = []string{
	"", "", "", "AND", "OR", "NOT", "EQU", "NEQ", "GT", "LT", "GTE", "LTE",
	"CONTAINS", "IN", "COMMA", "NUMBER", "VAR", "STRING", "WHITESPACE",
}

var ruleNames = []string{
	"start", "expression",
}
var decisionToDFA = make([]*antlr.DFA, len(deserializedATN.DecisionToState))

func init() {
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

type EventRuleParser struct {
	*antlr.BaseParser
}

func NewEventRuleParser(input antlr.TokenStream) *EventRuleParser {
	this := new(EventRuleParser)

	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "EventRule.g4"

	return this
}

// EventRuleParser tokens.
const (
	EventRuleParserEOF        = antlr.TokenEOF
	EventRuleParserT__0       = 1
	EventRuleParserT__1       = 2
	EventRuleParserAND        = 3
	EventRuleParserOR         = 4
	EventRuleParserNOT        = 5
	EventRuleParserEQU        = 6
	EventRuleParserNEQ        = 7
	EventRuleParserGT         = 8
	EventRuleParserLT         = 9
	EventRuleParserGTE        = 10
	EventRuleParserLTE        = 11
	EventRuleParserCONTAINS   = 12
	EventRuleParserIN         = 13
	EventRuleParserCOMMA      = 14
	EventRuleParserNUMBER     = 15
	EventRuleParserVAR        = 16
	EventRuleParserSTRING     = 17
	EventRuleParserWHITESPACE = 18
)

// EventRuleParser rules.
const (
	EventRuleParserRULE_start      = 0
	EventRuleParserRULE_expression = 1
)

// IStartContext is an interface to support dynamic dispatch.
type IStartContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsStartContext differentiates from other interfaces.
	IsStartContext()
}

type StartContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStartContext() *StartContext {
	var p = new(StartContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EventRuleParserRULE_start
	return p
}

func (*StartContext) IsStartContext() {}

func NewStartContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StartContext {
	var p = new(StartContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EventRuleParserRULE_start

	return p
}

func (s *StartContext) GetParser() antlr.Parser { return s.parser }

func (s *StartContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *StartContext) EOF() antlr.TerminalNode {
	return s.GetToken(EventRuleParserEOF, 0)
}

func (s *StartContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StartContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StartContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EventRuleVisitor:
		return t.VisitStart(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EventRuleParser) Start() (localctx IStartContext) {
	localctx = NewStartContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, EventRuleParserRULE_start)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(4)
		p.expression(0)
	}
	{
		p.SetState(5)
		p.Match(EventRuleParserEOF)
	}

	return localctx
}

// IExpressionContext is an interface to support dynamic dispatch.
type IExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsExpressionContext differentiates from other interfaces.
	IsExpressionContext()
}

type ExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpressionContext() *ExpressionContext {
	var p = new(ExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EventRuleParserRULE_expression
	return p
}

func (*ExpressionContext) IsExpressionContext() {}

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext {
	var p = new(ExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EventRuleParserRULE_expression

	return p
}

func (s *ExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionContext) CopyFrom(ctx *ExpressionContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *ExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type NotContext struct {
	*ExpressionContext
}

func NewNotContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NotContext {
	var p = new(NotContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *NotContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NotContext) NOT() antlr.TerminalNode {
	return s.GetToken(EventRuleParserNOT, 0)
}

func (s *NotContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *NotContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EventRuleVisitor:
		return t.VisitNot(s)

	default:
		return t.VisitChildren(s)
	}
}

type ParenthesisContext struct {
	*ExpressionContext
}

func NewParenthesisContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ParenthesisContext {
	var p = new(ParenthesisContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *ParenthesisContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParenthesisContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ParenthesisContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EventRuleVisitor:
		return t.VisitParenthesis(s)

	default:
		return t.VisitChildren(s)
	}
}

type VariableContext struct {
	*ExpressionContext
}

func NewVariableContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *VariableContext {
	var p = new(VariableContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *VariableContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VariableContext) VAR() antlr.TerminalNode {
	return s.GetToken(EventRuleParserVAR, 0)
}

func (s *VariableContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EventRuleVisitor:
		return t.VisitVariable(s)

	default:
		return t.VisitChildren(s)
	}
}

type CompareNumberContext struct {
	*ExpressionContext
	op antlr.Token
}

func NewCompareNumberContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CompareNumberContext {
	var p = new(CompareNumberContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *CompareNumberContext) GetOp() antlr.Token { return s.op }

func (s *CompareNumberContext) SetOp(v antlr.Token) { s.op = v }

func (s *CompareNumberContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CompareNumberContext) VAR() antlr.TerminalNode {
	return s.GetToken(EventRuleParserVAR, 0)
}

func (s *CompareNumberContext) NUMBER() antlr.TerminalNode {
	return s.GetToken(EventRuleParserNUMBER, 0)
}

func (s *CompareNumberContext) EQU() antlr.TerminalNode {
	return s.GetToken(EventRuleParserEQU, 0)
}

func (s *CompareNumberContext) NEQ() antlr.TerminalNode {
	return s.GetToken(EventRuleParserNEQ, 0)
}

func (s *CompareNumberContext) GT() antlr.TerminalNode {
	return s.GetToken(EventRuleParserGT, 0)
}

func (s *CompareNumberContext) LT() antlr.TerminalNode {
	return s.GetToken(EventRuleParserLT, 0)
}

func (s *CompareNumberContext) GTE() antlr.TerminalNode {
	return s.GetToken(EventRuleParserGTE, 0)
}

func (s *CompareNumberContext) LTE() antlr.TerminalNode {
	return s.GetToken(EventRuleParserLTE, 0)
}

func (s *CompareNumberContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EventRuleVisitor:
		return t.VisitCompareNumber(s)

	default:
		return t.VisitChildren(s)
	}
}

type StringEqualContainsContext struct {
	*ExpressionContext
	op antlr.Token
}

func NewStringEqualContainsContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *StringEqualContainsContext {
	var p = new(StringEqualContainsContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *StringEqualContainsContext) GetOp() antlr.Token { return s.op }

func (s *StringEqualContainsContext) SetOp(v antlr.Token) { s.op = v }

func (s *StringEqualContainsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StringEqualContainsContext) VAR() antlr.TerminalNode {
	return s.GetToken(EventRuleParserVAR, 0)
}

func (s *StringEqualContainsContext) STRING() antlr.TerminalNode {
	return s.GetToken(EventRuleParserSTRING, 0)
}

func (s *StringEqualContainsContext) EQU() antlr.TerminalNode {
	return s.GetToken(EventRuleParserEQU, 0)
}

func (s *StringEqualContainsContext) CONTAINS() antlr.TerminalNode {
	return s.GetToken(EventRuleParserCONTAINS, 0)
}

func (s *StringEqualContainsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EventRuleVisitor:
		return t.VisitStringEqualContains(s)

	default:
		return t.VisitChildren(s)
	}
}

type AndOrContext struct {
	*ExpressionContext
	op antlr.Token
}

func NewAndOrContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *AndOrContext {
	var p = new(AndOrContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *AndOrContext) GetOp() antlr.Token { return s.op }

func (s *AndOrContext) SetOp(v antlr.Token) { s.op = v }

func (s *AndOrContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AndOrContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *AndOrContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *AndOrContext) AND() antlr.TerminalNode {
	return s.GetToken(EventRuleParserAND, 0)
}

func (s *AndOrContext) OR() antlr.TerminalNode {
	return s.GetToken(EventRuleParserOR, 0)
}

func (s *AndOrContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EventRuleVisitor:
		return t.VisitAndOr(s)

	default:
		return t.VisitChildren(s)
	}
}

type StringInContext struct {
	*ExpressionContext
}

func NewStringInContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *StringInContext {
	var p = new(StringInContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *StringInContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StringInContext) VAR() antlr.TerminalNode {
	return s.GetToken(EventRuleParserVAR, 0)
}

func (s *StringInContext) IN() antlr.TerminalNode {
	return s.GetToken(EventRuleParserIN, 0)
}

func (s *StringInContext) AllSTRING() []antlr.TerminalNode {
	return s.GetTokens(EventRuleParserSTRING)
}

func (s *StringInContext) STRING(i int) antlr.TerminalNode {
	return s.GetToken(EventRuleParserSTRING, i)
}

func (s *StringInContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(EventRuleParserCOMMA)
}

func (s *StringInContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(EventRuleParserCOMMA, i)
}

func (s *StringInContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EventRuleVisitor:
		return t.VisitStringIn(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EventRuleParser) Expression() (localctx IExpressionContext) {
	return p.expression(0)
}

func (p *EventRuleParser) expression(_p int) (localctx IExpressionContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewExpressionContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IExpressionContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 2
	p.EnterRecursionRule(localctx, 2, EventRuleParserRULE_expression, _p)
	var _la int

	defer func() {
		p.UnrollRecursionContexts(_parentctx)
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(33)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 1, p.GetParserRuleContext()) {
	case 1:
		localctx = NewNotContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx

		{
			p.SetState(8)
			p.Match(EventRuleParserNOT)
		}
		{
			p.SetState(9)
			p.expression(6)
		}

	case 2:
		localctx = NewParenthesisContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(10)
			p.Match(EventRuleParserT__0)
		}
		{
			p.SetState(11)
			p.expression(0)
		}
		{
			p.SetState(12)
			p.Match(EventRuleParserT__1)
		}

	case 3:
		localctx = NewStringEqualContainsContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(14)
			p.Match(EventRuleParserVAR)
		}
		{
			p.SetState(15)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*StringEqualContainsContext).op = _lt

			_la = p.GetTokenStream().LA(1)

			if !(_la == EventRuleParserEQU || _la == EventRuleParserCONTAINS) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*StringEqualContainsContext).op = _ri
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(16)
			p.Match(EventRuleParserSTRING)
		}

	case 4:
		localctx = NewStringInContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(17)
			p.Match(EventRuleParserVAR)
		}
		{
			p.SetState(18)
			p.Match(EventRuleParserIN)
		}
		{
			p.SetState(19)
			p.Match(EventRuleParserT__0)
		}
		{
			p.SetState(20)
			p.Match(EventRuleParserSTRING)
		}
		p.SetState(25)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == EventRuleParserCOMMA {
			{
				p.SetState(21)
				p.Match(EventRuleParserCOMMA)
			}
			{
				p.SetState(22)
				p.Match(EventRuleParserSTRING)
			}

			p.SetState(27)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(28)
			p.Match(EventRuleParserT__1)
		}

	case 5:
		localctx = NewCompareNumberContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(29)
			p.Match(EventRuleParserVAR)
		}
		{
			p.SetState(30)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*CompareNumberContext).op = _lt

			_la = p.GetTokenStream().LA(1)

			if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<EventRuleParserEQU)|(1<<EventRuleParserNEQ)|(1<<EventRuleParserGT)|(1<<EventRuleParserLT)|(1<<EventRuleParserGTE)|(1<<EventRuleParserLTE))) != 0) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*CompareNumberContext).op = _ri
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(31)
			p.Match(EventRuleParserNUMBER)
		}

	case 6:
		localctx = NewVariableContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(32)
			p.Match(EventRuleParserVAR)
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(40)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 2, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewAndOrContext(p, NewExpressionContext(p, _parentctx, _parentState))
			p.PushNewRecursionContext(localctx, _startState, EventRuleParserRULE_expression)
			p.SetState(35)

			if !(p.Precpred(p.GetParserRuleContext(), 7)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 7)", ""))
			}
			{
				p.SetState(36)

				var _lt = p.GetTokenStream().LT(1)

				localctx.(*AndOrContext).op = _lt

				_la = p.GetTokenStream().LA(1)

				if !(_la == EventRuleParserAND || _la == EventRuleParserOR) {
					var _ri = p.GetErrorHandler().RecoverInline(p)

					localctx.(*AndOrContext).op = _ri
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(37)
				p.expression(8)
			}

		}
		p.SetState(42)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 2, p.GetParserRuleContext())
	}

	return localctx
}

func (p *EventRuleParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 1:
		var t *ExpressionContext = nil
		if localctx != nil {
			t = localctx.(*ExpressionContext)
		}
		return p.Expression_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *EventRuleParser) Expression_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 7)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
