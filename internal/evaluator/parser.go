package evaluator

import (
	"fmt"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	LOGOR       // OR
	LOGAND      // AND
	COMPARISON  // ==, !=, >, <, >=, <=
)

var precedences = map[TokenType]int{
	TokenOR:       LOGOR,
	TokenAND:      LOGAND,
	TokenEqual:    COMPARISON,
	TokenNotEqual: COMPARISON,
	TokenGT:       COMPARISON,
	TokenLT:       COMPARISON,
	TokenGTE:      COMPARISON,
	TokenLTE:      COMPARISON,
}

type Parser struct {
	l         *Lexer
	curToken  Token
	peekToken Token
	errors    []string
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	// Read two tokens so curToken and peekToken are both populated
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) ParseRule() (Node, error) {
	expression := p.parseExpression(LOWEST)
	if expression == nil {
		return nil, fmt.Errorf("failed to parse valid rule hierarchy. Errors: %v", p.errors)
	}
	
	if p.curToken.Type != TokenEOF {
		return nil, fmt.Errorf("unexpected trailing syntax token detected: %s (%s)", p.curToken.Type, p.curToken.Literal)
	}
	return expression, nil
}

func (p *Parser) parseExpression(precedence int) Node {
	// Left-hand side prefix resolution
	leftExp := p.parsePrefix()
	if leftExp == nil {
		return nil
	}

	// Infix calculation loop while next token has higher precedence
	for p.peekToken.Type != TokenEOF && precedence < p.peekPrecedence() {
		p.nextToken()
		leftExp = p.parseInfix(leftExp)
	}

	return leftExp
}

func (p *Parser) parsePrefix() Node {
	switch p.curToken.Type {
	case TokenIdentifier:
		return &IdentifierNode{Name: p.curToken.Literal}
	case TokenString:
		return &LiteralNode{Value: p.curToken.Literal, Kind: TypeString}
	case TokenNumber:
		val, err := strconv.ParseFloat(p.curToken.Literal, 64)
		if err != nil {
			p.errors = append(p.errors, fmt.Sprintf("could not parse number literal %s", p.curToken.Literal))
			return nil
		}
		return &LiteralNode{Value: val, Kind: TypeNumber}
	default:
		p.errors = append(p.errors, fmt.Sprintf("no prefix parse function found for token type %s", p.curToken.Type))
		return nil
	}
}

func (p *Parser) parseInfix(left Node) Node {
	tokenType := p.curToken.Type
	precedence := p.curPrecedence()
	literal := p.curToken.Literal
	
	p.nextToken()
	right := p.parseExpression(precedence)

	switch tokenType {
	case TokenAND, TokenOR:
		return &LogicalNode{Left: left, Operator: literal, Right: right}
	default:
		return &ComparisonNode{Left: left, Operator: literal, Right: right}
	}
}