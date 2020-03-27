// Code generated from EventRule.g4 by ANTLR 4.7.2. DO NOT EDIT.

package parser

import (
	"fmt"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = unicode.IsLetter

var serializedLexerAtn = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 2, 20, 126,
	8, 1, 4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7,
	9, 7, 4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12,
	4, 13, 9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4,
	18, 9, 18, 4, 19, 9, 19, 4, 20, 9, 20, 3, 2, 3, 2, 3, 3, 3, 3, 3, 4, 3,
	4, 3, 4, 3, 4, 3, 5, 3, 5, 3, 5, 3, 6, 3, 6, 3, 6, 3, 6, 3, 7, 3, 7, 3,
	8, 3, 8, 3, 8, 3, 9, 3, 9, 3, 10, 3, 10, 3, 11, 3, 11, 3, 11, 3, 12, 3,
	12, 3, 12, 3, 13, 3, 13, 3, 13, 3, 13, 3, 13, 3, 13, 3, 13, 3, 13, 3, 13,
	3, 14, 3, 14, 3, 14, 3, 15, 3, 15, 3, 16, 5, 16, 87, 10, 16, 3, 16, 6,
	16, 90, 10, 16, 13, 16, 14, 16, 91, 3, 16, 3, 16, 6, 16, 96, 10, 16, 13,
	16, 14, 16, 97, 5, 16, 100, 10, 16, 3, 17, 6, 17, 103, 10, 17, 13, 17,
	14, 17, 104, 3, 18, 3, 18, 3, 18, 7, 18, 110, 10, 18, 12, 18, 14, 18, 113,
	11, 18, 3, 18, 3, 18, 3, 19, 3, 19, 3, 19, 3, 19, 3, 20, 3, 20, 3, 20,
	3, 20, 5, 20, 125, 10, 20, 3, 111, 2, 21, 3, 3, 5, 4, 7, 5, 9, 6, 11, 7,
	13, 8, 15, 9, 17, 10, 19, 11, 21, 12, 23, 13, 25, 14, 27, 15, 29, 16, 31,
	17, 33, 18, 35, 19, 37, 20, 39, 2, 3, 2, 6, 3, 2, 47, 47, 3, 2, 50, 59,
	7, 2, 47, 48, 50, 59, 67, 92, 97, 97, 99, 124, 5, 2, 11, 12, 15, 15, 34,
	34, 2, 132, 2, 3, 3, 2, 2, 2, 2, 5, 3, 2, 2, 2, 2, 7, 3, 2, 2, 2, 2, 9,
	3, 2, 2, 2, 2, 11, 3, 2, 2, 2, 2, 13, 3, 2, 2, 2, 2, 15, 3, 2, 2, 2, 2,
	17, 3, 2, 2, 2, 2, 19, 3, 2, 2, 2, 2, 21, 3, 2, 2, 2, 2, 23, 3, 2, 2, 2,
	2, 25, 3, 2, 2, 2, 2, 27, 3, 2, 2, 2, 2, 29, 3, 2, 2, 2, 2, 31, 3, 2, 2,
	2, 2, 33, 3, 2, 2, 2, 2, 35, 3, 2, 2, 2, 2, 37, 3, 2, 2, 2, 3, 41, 3, 2,
	2, 2, 5, 43, 3, 2, 2, 2, 7, 45, 3, 2, 2, 2, 9, 49, 3, 2, 2, 2, 11, 52,
	3, 2, 2, 2, 13, 56, 3, 2, 2, 2, 15, 58, 3, 2, 2, 2, 17, 61, 3, 2, 2, 2,
	19, 63, 3, 2, 2, 2, 21, 65, 3, 2, 2, 2, 23, 68, 3, 2, 2, 2, 25, 71, 3,
	2, 2, 2, 27, 80, 3, 2, 2, 2, 29, 83, 3, 2, 2, 2, 31, 86, 3, 2, 2, 2, 33,
	102, 3, 2, 2, 2, 35, 106, 3, 2, 2, 2, 37, 116, 3, 2, 2, 2, 39, 124, 3,
	2, 2, 2, 41, 42, 7, 42, 2, 2, 42, 4, 3, 2, 2, 2, 43, 44, 7, 43, 2, 2, 44,
	6, 3, 2, 2, 2, 45, 46, 7, 99, 2, 2, 46, 47, 7, 112, 2, 2, 47, 48, 7, 102,
	2, 2, 48, 8, 3, 2, 2, 2, 49, 50, 7, 113, 2, 2, 50, 51, 7, 116, 2, 2, 51,
	10, 3, 2, 2, 2, 52, 53, 7, 112, 2, 2, 53, 54, 7, 113, 2, 2, 54, 55, 7,
	118, 2, 2, 55, 12, 3, 2, 2, 2, 56, 57, 7, 63, 2, 2, 57, 14, 3, 2, 2, 2,
	58, 59, 7, 35, 2, 2, 59, 60, 7, 63, 2, 2, 60, 16, 3, 2, 2, 2, 61, 62, 7,
	64, 2, 2, 62, 18, 3, 2, 2, 2, 63, 64, 7, 62, 2, 2, 64, 20, 3, 2, 2, 2,
	65, 66, 7, 64, 2, 2, 66, 67, 7, 63, 2, 2, 67, 22, 3, 2, 2, 2, 68, 69, 7,
	62, 2, 2, 69, 70, 7, 63, 2, 2, 70, 24, 3, 2, 2, 2, 71, 72, 7, 101, 2, 2,
	72, 73, 7, 113, 2, 2, 73, 74, 7, 112, 2, 2, 74, 75, 7, 118, 2, 2, 75, 76,
	7, 99, 2, 2, 76, 77, 7, 107, 2, 2, 77, 78, 7, 112, 2, 2, 78, 79, 7, 117,
	2, 2, 79, 26, 3, 2, 2, 2, 80, 81, 7, 107, 2, 2, 81, 82, 7, 112, 2, 2, 82,
	28, 3, 2, 2, 2, 83, 84, 7, 46, 2, 2, 84, 30, 3, 2, 2, 2, 85, 87, 9, 2,
	2, 2, 86, 85, 3, 2, 2, 2, 86, 87, 3, 2, 2, 2, 87, 89, 3, 2, 2, 2, 88, 90,
	9, 3, 2, 2, 89, 88, 3, 2, 2, 2, 90, 91, 3, 2, 2, 2, 91, 89, 3, 2, 2, 2,
	91, 92, 3, 2, 2, 2, 92, 99, 3, 2, 2, 2, 93, 95, 7, 48, 2, 2, 94, 96, 9,
	3, 2, 2, 95, 94, 3, 2, 2, 2, 96, 97, 3, 2, 2, 2, 97, 95, 3, 2, 2, 2, 97,
	98, 3, 2, 2, 2, 98, 100, 3, 2, 2, 2, 99, 93, 3, 2, 2, 2, 99, 100, 3, 2,
	2, 2, 100, 32, 3, 2, 2, 2, 101, 103, 9, 4, 2, 2, 102, 101, 3, 2, 2, 2,
	103, 104, 3, 2, 2, 2, 104, 102, 3, 2, 2, 2, 104, 105, 3, 2, 2, 2, 105,
	34, 3, 2, 2, 2, 106, 111, 7, 36, 2, 2, 107, 110, 5, 39, 20, 2, 108, 110,
	11, 2, 2, 2, 109, 107, 3, 2, 2, 2, 109, 108, 3, 2, 2, 2, 110, 113, 3, 2,
	2, 2, 111, 112, 3, 2, 2, 2, 111, 109, 3, 2, 2, 2, 112, 114, 3, 2, 2, 2,
	113, 111, 3, 2, 2, 2, 114, 115, 7, 36, 2, 2, 115, 36, 3, 2, 2, 2, 116,
	117, 9, 5, 2, 2, 117, 118, 3, 2, 2, 2, 118, 119, 8, 19, 2, 2, 119, 38,
	3, 2, 2, 2, 120, 121, 7, 94, 2, 2, 121, 125, 7, 36, 2, 2, 122, 123, 7,
	94, 2, 2, 123, 125, 7, 94, 2, 2, 124, 120, 3, 2, 2, 2, 124, 122, 3, 2,
	2, 2, 125, 40, 3, 2, 2, 2, 11, 2, 86, 91, 97, 99, 104, 109, 111, 124, 3,
	8, 2, 2,
}

var lexerDeserializer = antlr.NewATNDeserializer(nil)
var lexerAtn = lexerDeserializer.DeserializeFromUInt16(serializedLexerAtn)

var lexerChannelNames = []string{
	"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
}

var lexerModeNames = []string{
	"DEFAULT_MODE",
}

var lexerLiteralNames = []string{
	"", "'('", "')'", "'and'", "'or'", "'not'", "'='", "'!='", "'>'", "'<'",
	"'>='", "'<='", "'contains'", "'in'", "','",
}

var lexerSymbolicNames = []string{
	"", "", "", "AND", "OR", "NOT", "EQU", "NEQ", "GT", "LT", "GTE", "LTE",
	"CONTAINS", "IN", "COMMA", "NUMBER", "VAR", "STRING", "WHITESPACE",
}

var lexerRuleNames = []string{
	"T__0", "T__1", "AND", "OR", "NOT", "EQU", "NEQ", "GT", "LT", "GTE", "LTE",
	"CONTAINS", "IN", "COMMA", "NUMBER", "VAR", "STRING", "WHITESPACE", "ESC",
}

type EventRuleLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var lexerDecisionToDFA = make([]*antlr.DFA, len(lexerAtn.DecisionToState))

func init() {
	for index, ds := range lexerAtn.DecisionToState {
		lexerDecisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

func NewEventRuleLexer(input antlr.CharStream) *EventRuleLexer {

	l := new(EventRuleLexer)

	l.BaseLexer = antlr.NewBaseLexer(input)
	l.Interpreter = antlr.NewLexerATNSimulator(l, lexerAtn, lexerDecisionToDFA, antlr.NewPredictionContextCache())

	l.channelNames = lexerChannelNames
	l.modeNames = lexerModeNames
	l.RuleNames = lexerRuleNames
	l.LiteralNames = lexerLiteralNames
	l.SymbolicNames = lexerSymbolicNames
	l.GrammarFileName = "EventRule.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// EventRuleLexer tokens.
const (
	EventRuleLexerT__0       = 1
	EventRuleLexerT__1       = 2
	EventRuleLexerAND        = 3
	EventRuleLexerOR         = 4
	EventRuleLexerNOT        = 5
	EventRuleLexerEQU        = 6
	EventRuleLexerNEQ        = 7
	EventRuleLexerGT         = 8
	EventRuleLexerLT         = 9
	EventRuleLexerGTE        = 10
	EventRuleLexerLTE        = 11
	EventRuleLexerCONTAINS   = 12
	EventRuleLexerIN         = 13
	EventRuleLexerCOMMA      = 14
	EventRuleLexerNUMBER     = 15
	EventRuleLexerVAR        = 16
	EventRuleLexerSTRING     = 17
	EventRuleLexerWHITESPACE = 18
)
