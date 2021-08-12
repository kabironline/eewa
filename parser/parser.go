package parser

import (
	"fmt"
	"strconv"

	"github.com/kabironline/monke/ast"
	"github.com/kabironline/monke/lexer"
	"github.com/kabironline/monke/tokens"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

var precedences = map[tokens.TokenType]int{
	tokens.EQ:       EQUALS,
	tokens.NOT_EQ:   EQUALS,
	tokens.LT:       LESSGREATER,
	tokens.GT:       LESSGREATER,
	tokens.PLUS:     SUM,
	tokens.MINUS:    SUM,
	tokens.SLASH:    PRODUCT,
	tokens.ASTERISK: PRODUCT,
}

type Parser struct {
	l *lexer.Lexer

	errors []string

	curToken  tokens.Token
	peekToken tokens.Token

	prefixParseFns map[tokens.TokenType]prefixParseFn
	infixParseFns  map[tokens.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[tokens.TokenType]prefixParseFn)
	p.registerPrefix(tokens.IDENT, p.parseIdentifier)
	p.registerPrefix(tokens.INT, p.parseIntegerLiteral)
	p.registerPrefix(tokens.BANG, p.parsePrefixExpression)
	p.registerPrefix(tokens.MINUS, p.parsePrefixExpression)

	p.infixParseFns = make(map[tokens.TokenType]infixParseFn)
	p.registerInfix(tokens.PLUS, p.parseInfixExpression)
	p.registerInfix(tokens.MINUS, p.parseInfixExpression)
	p.registerInfix(tokens.SLASH, p.parseInfixExpression)
	p.registerInfix(tokens.ASTERISK, p.parseInfixExpression)
	p.registerInfix(tokens.EQ, p.parseInfixExpression)
	p.registerInfix(tokens.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(tokens.LT, p.parseInfixExpression)
	p.registerInfix(tokens.GT, p.parseInfixExpression)

	p.registerPrefix(tokens.TRUE, p.parseBoolean)
	p.registerPrefix(tokens.FALSE, p.parseBoolean)
	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.curToken.Type != tokens.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case tokens.LET:
		return p.parseLetStatement()
	case tokens.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseIdentifier() ast.Expession {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

//Parsing
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeek(tokens.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(tokens.ASSIGN) {
		return nil
	}
	// TODO: We're skipping the expressions until we
	// encounter a semicolon
	for !p.curTokenIs(tokens.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()
	// TODO: We're skipping the expressions until we
	// encounter a semicolon
	for !p.curTokenIs(tokens.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

//Token Checking
func (p *Parser) curTokenIs(t tokens.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t tokens.TokenType) bool {
	return p.peekToken.Type == t
}
func (p *Parser) expectPeek(t tokens.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

//Expression Parsing
type (
	prefixParseFn func() ast.Expession
	infixParseFn  func(ast.Expession) ast.Expession
)

func (p *Parser) registerPrefix(tokenType tokens.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
func (p *Parser) registerInfix(tokenType tokens.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expession = p.parseExpression(LOWEST)
	if p.peekTokenIs(tokens.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

//Prefix Parsing
func (p *Parser) parseExpression(precedence int) ast.Expession {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()
	for !p.peekTokenIs(tokens.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parsePrefixExpression() ast.Expession {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

//Integer Literal Parsing
func (p *Parser) parseIntegerLiteral() ast.Expession {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

//Infix Parsing
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseInfixExpression(left ast.Expession) ast.Expession {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

//Boolean Parsing

func (p *Parser) parseBoolean() ast.Expession {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(tokens.TRUE)}
}

//Error Handling
func (p *Parser) Errors() []string {
	return p.errors
}
func (p *Parser) peekError(t tokens.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t tokens.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}
