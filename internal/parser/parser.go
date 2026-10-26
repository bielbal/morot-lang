package parser

import (
	"fmt"
	"morot-lang/internal/ast"
	"morot-lang/internal/lexer"
	"strconv"
)

type Parser struct {
	tokens []lexer.Token
	pos    int
}

func New(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) peek() lexer.Token {
	for p.pos < len(p.tokens) && p.tokens[p.pos].Type == lexer.NEWLINE {
		// don't skip here, caller decides
		break
	}
	if p.pos >= len(p.tokens) {
		return lexer.Token{Type: lexer.EOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) peekSkipNL() lexer.Token {
	i := p.pos
	for i < len(p.tokens) && p.tokens[i].Type == lexer.NEWLINE {
		i++
	}
	if i >= len(p.tokens) {
		return lexer.Token{Type: lexer.EOF}
	}
	return p.tokens[i]
}

func (p *Parser) advance() lexer.Token {
	tok := p.tokens[p.pos]
	p.pos++
	return tok
}

func (p *Parser) skipNL() {
	for p.pos < len(p.tokens) && p.tokens[p.pos].Type == lexer.NEWLINE {
		p.pos++
	}
}

func (p *Parser) skipTerminators() {
	for p.pos < len(p.tokens) && (p.tokens[p.pos].Type == lexer.NEWLINE || p.tokens[p.pos].Type == lexer.SEMICOLON) {
		p.pos++
	}
}

func (p *Parser) expect(t lexer.TokenType) (lexer.Token, error) {
	p.skipNL()
	tok := p.peek()
	if tok.Type != t {
		return tok, fmt.Errorf("line %d: expected %s, got %s (%q)", tok.Line, t, tok.Type, tok.Literal)
	}
	return p.advance(), nil
}

func (p *Parser) Parse() (*ast.Program, error) {
	prog := &ast.Program{}
	for {
		p.skipTerminators()
		if p.peek().Type == lexer.EOF {
			break
		}
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			prog.Statements = append(prog.Statements, stmt)
		}
	}
	return prog, nil
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	tok := p.peek()
	switch tok.Type {
	case lexer.LET, lexer.CONST:
		return p.parseVarDecl()
	case lexer.FN:
		return p.parseFnDecl()
	case lexer.RETURN:
		return p.parseReturn()
	case lexer.IF:
		return p.parseIf()
	case lexer.WHILE:
		return p.parseWhile()
	case lexer.FOR:
		return p.parseFor()
	case lexer.BREAK:
		p.advance()
		return &ast.BreakStmt{Line: tok.Line}, nil
	case lexer.CONTINUE:
		p.advance()
		return &ast.ContinueStmt{Line: tok.Line}, nil
	case lexer.LOG:
		return p.parseLog()
	case lexer.IMPORT:
		return p.parseImport()
	case lexer.TYPE:
		return p.parseTypeDecl()
	default:
		return p.parseExprOrAssign()
	}
}

func (p *Parser) parseVarDecl() (*ast.VarDecl, error) {
	tok := p.advance()
	isConst := tok.Type == lexer.CONST
	line := tok.Line

	nameTok, err := p.expect(lexer.IDENT)
	if err != nil {
		return nil, err
	}

	var typeHint string
	if p.peekSkipNL().Type == lexer.COLON {
		p.skipNL()
		p.advance() // ':'
		p.skipNL()
		typeTok := p.advance()
		typeHint = typeTok.Literal
	}

	p.skipNL()
	if _, err := p.expect(lexer.ASSIGN); err != nil {
		return nil, err
	}
	p.skipNL()

	val, err := p.parseExpr(0)
	if err != nil {
		return nil, err
	}

	return &ast.VarDecl{
		Name:     nameTok.Literal,
		TypeHint: typeHint,
		IsConst:  isConst,
		Value:    val,
		Line:     line,
	}, nil
}

func (p *Parser) parseFnDecl() (*ast.FnDecl, error) {
	tok := p.advance() // 'fn'
	line := tok.Line

	nameTok, err := p.expect(lexer.IDENT)
	if err != nil {
		return nil, err
	}

	params, err := p.parseParams()
	if err != nil {
		return nil, err
	}

	var retType string
	p.skipNL()
	if p.peek().Type == lexer.ARROW {
		p.advance()
		p.skipNL()
		retTok := p.advance()
		retType = retTok.Literal
	}

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &ast.FnDecl{
		Name:       nameTok.Literal,
		Params:     params,
		ReturnType: retType,
		Body:       body,
		Line:       line,
	}, nil
}

func (p *Parser) parseParams() ([]ast.Param, error) {
	if _, err := p.expect(lexer.LPAREN); err != nil {
		return nil, err
	}
	var params []ast.Param
	p.skipNL()
	for p.peek().Type != lexer.RPAREN && p.peek().Type != lexer.EOF {
		nameTok, err := p.expect(lexer.IDENT)
		if err != nil {
			return nil, err
		}
		var typeHint string
		p.skipNL()
		if p.peek().Type == lexer.COLON {
			p.advance()
			p.skipNL()
			typeTok := p.advance()
			typeHint = typeTok.Literal
		}
		params = append(params, ast.Param{Name: nameTok.Literal, TypeHint: typeHint})
		p.skipNL()
		if p.peek().Type == lexer.COMMA {
			p.advance()
			p.skipNL()
		}
	}
	if _, err := p.expect(lexer.RPAREN); err != nil {
		return nil, err
	}
	return params, nil
}

func (p *Parser) parseBlock() (*ast.BlockStmt, error) {
	p.skipNL()
	if _, err := p.expect(lexer.LBRACE); err != nil {
		return nil, err
	}
	block := &ast.BlockStmt{}
	for {
		p.skipTerminators()
		if p.peek().Type == lexer.RBRACE || p.peek().Type == lexer.EOF {
			break
		}
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
	}
	if _, err := p.expect(lexer.RBRACE); err != nil {
		return nil, err
	}
	return block, nil
}

func (p *Parser) parseReturn() (*ast.ReturnStmt, error) {
	tok := p.advance()
	line := tok.Line

	// if next meaningful token is on same "line" area, parse value
	next := p.peek()
	if next.Type == lexer.NEWLINE || next.Type == lexer.SEMICOLON || next.Type == lexer.RBRACE || next.Type == lexer.EOF {
		return &ast.ReturnStmt{Line: line}, nil
	}

	val, err := p.parseExpr(0)
	if err != nil {
		return nil, err
	}
	return &ast.ReturnStmt{Value: val, Line: line}, nil
}

func (p *Parser) parseIf() (*ast.IfStmt, error) {
	tok := p.advance() // 'if'
	line := tok.Line

	cond, err := p.parseExpr(0)
	if err != nil {
		return nil, err
	}

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	stmt := &ast.IfStmt{Condition: cond, Consequence: body, Line: line}

	for {
		p.skipNL()
		if p.peek().Type == lexer.ELSEIF {
			p.advance()
			eic, err := p.parseExpr(0)
			if err != nil {
				return nil, err
			}
			eib, err := p.parseBlock()
			if err != nil {
				return nil, err
			}
			stmt.ElseIfs = append(stmt.ElseIfs, ast.ElseIfClause{Condition: eic, Body: eib})
		} else if p.peek().Type == lexer.ELSE {
			p.advance()
			alt, err := p.parseBlock()
			if err != nil {
				return nil, err
			}
			stmt.Alternative = alt
			break
		} else {
			break
		}
	}

	return stmt, nil
}

func (p *Parser) parseWhile() (*ast.WhileStmt, error) {
	tok := p.advance()
	line := tok.Line
	cond, err := p.parseExpr(0)
	if err != nil {
		return nil, err
	}
	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	return &ast.WhileStmt{Condition: cond, Body: body, Line: line}, nil
}

func (p *Parser) parseFor() (*ast.ForInStmt, error) {
	tok := p.advance()
	line := tok.Line
	p.skipNL()
	varTok, err := p.expect(lexer.IDENT)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.IN); err != nil {
		return nil, err
	}
	p.skipNL()
	iter, err := p.parseExpr(0)
	if err != nil {
		return nil, err
	}
	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	return &ast.ForInStmt{Var: varTok.Literal, Iterable: iter, Body: body, Line: line}, nil
}

func (p *Parser) parseLog() (*ast.LogStmt, error) {
	tok := p.advance()
	line := tok.Line
	if _, err := p.expect(lexer.LPAREN); err != nil {
		return nil, err
	}
	var args []ast.Expression
	p.skipNL()
	for p.peek().Type != lexer.RPAREN && p.peek().Type != lexer.EOF {
		arg, err := p.parseExpr(0)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		p.skipNL()
		if p.peek().Type == lexer.COMMA {
			p.advance()
			p.skipNL()
		}
	}
	if _, err := p.expect(lexer.RPAREN); err != nil {
		return nil, err
	}
	return &ast.LogStmt{Args: args, Line: line}, nil
}

func (p *Parser) parseImport() (*ast.ImportStmt, error) {
	tok := p.advance()
	line := tok.Line
	p.skipNL()

	var names []string

	if p.peek().Type == lexer.LBRACE {
		p.advance()
		p.skipNL()
		for p.peek().Type != lexer.RBRACE && p.peek().Type != lexer.EOF {
			n, err := p.expect(lexer.IDENT)
			if err != nil {
				return nil, err
			}
			names = append(names, n.Literal)
			p.skipNL()
			if p.peek().Type == lexer.COMMA {
				p.advance()
				p.skipNL()
			}
		}
		if _, err := p.expect(lexer.RBRACE); err != nil {
			return nil, err
		}
		p.skipNL()
	}

	if _, err := p.expect(lexer.FROM); err != nil {
		return nil, err
	}
	p.skipNL()
	modTok, err := p.expect(lexer.STRING_LIT)
	if err != nil {
		return nil, err
	}

	return &ast.ImportStmt{Names: names, Module: modTok.Literal, Line: line}, nil
}

func (p *Parser) parseTypeDecl() (*ast.TypeDecl, error) {
	tok := p.advance()
	line := tok.Line
	p.skipNL()
	nameTok, err := p.expect(lexer.IDENT)
	if err != nil {
		return nil, err
	}
	p.skipNL()
	if _, err := p.expect(lexer.LBRACE); err != nil {
		return nil, err
	}
	var fields []ast.TypeField
	for {
		p.skipTerminators()
		if p.peek().Type == lexer.RBRACE || p.peek().Type == lexer.EOF {
			break
		}
		fieldName, err := p.expect(lexer.IDENT)
		if err != nil {
			return nil, err
		}
		optional := false
		p.skipNL()
		if p.peek().Type == lexer.BANG {
			// not optional
		}
		// check for '?'  - we use optional marker
		// simplified: always require colon
		var typeHint string
		if p.peek().Type == lexer.COLON {
			p.advance()
			p.skipNL()
			typeTok := p.advance()
			typeHint = typeTok.Literal
		}
		fields = append(fields, ast.TypeField{Name: fieldName.Literal, TypeHint: typeHint, Optional: optional})
	}
	if _, err := p.expect(lexer.RBRACE); err != nil {
		return nil, err
	}
	return &ast.TypeDecl{Name: nameTok.Literal, Fields: fields, Line: line}, nil
}

// ─── Expression parsing (Pratt) ────────────────────────────────────────────

func (p *Parser) parseExprOrAssign() (ast.Statement, error) {
	expr, err := p.parseExpr(0)
	if err != nil {
		return nil, err
	}

	line := p.peek().Line
	switch p.peek().Type {
	case lexer.ASSIGN, lexer.PLUS_ASSIGN, lexer.MINUS_ASSIGN:
		op := p.advance().Literal
		p.skipNL()
		val, err := p.parseExpr(0)
		if err != nil {
			return nil, err
		}
		return &ast.AssignStmt{Target: expr, Op: op, Value: val, Line: line}, nil
	}

	return &ast.ExprStmt{Expr: expr, Line: line}, nil
}

type prefixParseFn func() (ast.Expression, error)
type infixParseFn func(ast.Expression) (ast.Expression, error)

func (p *Parser) prefixPrecedence(t lexer.TokenType) int {
	switch t {
	case lexer.BANG, lexer.MINUS:
		return 7
	}
	return 0
}

func (p *Parser) infixPrecedence(t lexer.TokenType) int {
	switch t {
	case lexer.OR:
		return 1
	case lexer.AND:
		return 2
	case lexer.EQ, lexer.NEQ:
		return 3
	case lexer.LT, lexer.GT, lexer.LTE, lexer.GTE:
		return 4
	case lexer.PLUS, lexer.MINUS:
		return 5
	case lexer.STAR, lexer.SLASH, lexer.PERCENT:
		return 6
	case lexer.LPAREN:
		return 8
	case lexer.DOT:
		return 9
	case lexer.LBRACKET:
		return 9
	}
	return 0
}

func (p *Parser) parseExpr(minPrec int) (ast.Expression, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	for {
		tok := p.peek()
		prec := p.infixPrecedence(tok.Type)
		if prec <= minPrec {
			break
		}

		switch tok.Type {
		case lexer.LPAREN:
			left, err = p.parseCall(left)
		case lexer.DOT:
			left, err = p.parseMember(left)
		case lexer.LBRACKET:
			left, err = p.parseIndex(left)
		default:
			p.advance()
			p.skipNL()
			right, rerr := p.parseExpr(prec)
			if rerr != nil {
				return nil, rerr
			}
			left = &ast.BinaryExpr{Left: left, Op: tok.Literal, Right: right, Line: tok.Line}
		}
		if err != nil {
			return nil, err
		}
	}
	return left, nil
}

func (p *Parser) parseUnary() (ast.Expression, error) {
	tok := p.peek()
	switch tok.Type {
	case lexer.BANG, lexer.MINUS:
		p.advance()
		operand, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{Op: tok.Literal, Operand: operand, Line: tok.Line}, nil
	}
	return p.parsePrimary()
}

func (p *Parser) parsePrimary() (ast.Expression, error) {
	tok := p.peek()
	switch tok.Type {
	case lexer.INT_LIT:
		p.advance()
		v, _ := strconv.ParseInt(tok.Literal, 10, 64)
		return &ast.IntLit{Value: v, Line: tok.Line}, nil
	case lexer.FLOAT_LIT:
		p.advance()
		v, _ := strconv.ParseFloat(tok.Literal, 64)
		return &ast.FloatLit{Value: v, Line: tok.Line}, nil
	case lexer.STRING_LIT:
		p.advance()
		return &ast.StringLit{Value: tok.Literal, Line: tok.Line}, nil
	case lexer.TRUE:
		p.advance()
		return &ast.BoolLit{Value: true, Line: tok.Line}, nil
	case lexer.FALSE:
		p.advance()
		return &ast.BoolLit{Value: false, Line: tok.Line}, nil
	case lexer.NULL:
		p.advance()
		return &ast.NullLit{Line: tok.Line}, nil
	case lexer.IDENT:
		p.advance()
		return &ast.Identifier{Name: tok.Literal, Line: tok.Line}, nil
	case lexer.LPAREN:
		p.advance()
		p.skipNL()
		expr, err := p.parseExpr(0)
		if err != nil {
			return nil, err
		}
		p.skipNL()
		if _, err := p.expect(lexer.RPAREN); err != nil {
			return nil, err
		}
		return expr, nil
	case lexer.LBRACKET:
		return p.parseArray()
	case lexer.LBRACE:
		return p.parseObject()
	case lexer.FN:
		return p.parseFnLit()
	case lexer.NEW:
		return p.parseNew()
	default:
		return nil, fmt.Errorf("line %d: unexpected token %s (%q)", tok.Line, tok.Type, tok.Literal)
	}
}

func (p *Parser) parseCall(callee ast.Expression) (*ast.CallExpr, error) {
	tok := p.advance() // '('
	var args []ast.Expression
	p.skipNL()
	for p.peek().Type != lexer.RPAREN && p.peek().Type != lexer.EOF {
		arg, err := p.parseExpr(0)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		p.skipNL()
		if p.peek().Type == lexer.COMMA {
			p.advance()
			p.skipNL()
		}
	}
	if _, err := p.expect(lexer.RPAREN); err != nil {
		return nil, err
	}
	return &ast.CallExpr{Callee: callee, Args: args, Line: tok.Line}, nil
}

func (p *Parser) parseMember(obj ast.Expression) (*ast.MemberExpr, error) {
	tok := p.advance() // '.'
	prop, err := p.expect(lexer.IDENT)
	if err != nil {
		return nil, err
	}
	return &ast.MemberExpr{Object: obj, Property: prop.Literal, Line: tok.Line}, nil
}

func (p *Parser) parseIndex(obj ast.Expression) (*ast.IndexExpr, error) {
	tok := p.advance() // '['
	idx, err := p.parseExpr(0)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.RBRACKET); err != nil {
		return nil, err
	}
	return &ast.IndexExpr{Object: obj, Index: idx, Line: tok.Line}, nil
}

func (p *Parser) parseArray() (*ast.ArrayLit, error) {
	tok := p.advance() // '['
	var elems []ast.Expression
	p.skipNL()
	for p.peek().Type != lexer.RBRACKET && p.peek().Type != lexer.EOF {
		el, err := p.parseExpr(0)
		if err != nil {
			return nil, err
		}
		elems = append(elems, el)
		p.skipNL()
		if p.peek().Type == lexer.COMMA {
			p.advance()
			p.skipNL()
		}
	}
	if _, err := p.expect(lexer.RBRACKET); err != nil {
		return nil, err
	}
	return &ast.ArrayLit{Elements: elems, Line: tok.Line}, nil
}

func (p *Parser) parseObject() (*ast.ObjectLit, error) {
	tok := p.advance() // '{'
	var fields []ast.ObjectField
	p.skipNL()
	for p.peek().Type != lexer.RBRACE && p.peek().Type != lexer.EOF {
		// key: value
		var key string
		if p.peek().Type == lexer.STRING_LIT {
			key = p.advance().Literal
		} else {
			keyTok, err := p.expect(lexer.IDENT)
			if err != nil {
				return nil, err
			}
			key = keyTok.Literal
		}
		p.skipNL()
		if _, err := p.expect(lexer.COLON); err != nil {
			return nil, err
		}
		p.skipNL()
		val, err := p.parseExpr(0)
		if err != nil {
			return nil, err
		}
		fields = append(fields, ast.ObjectField{Key: key, Value: val})
		p.skipNL()
		if p.peek().Type == lexer.COMMA {
			p.advance()
			p.skipNL()
		}
	}
	if _, err := p.expect(lexer.RBRACE); err != nil {
		return nil, err
	}
	return &ast.ObjectLit{Fields: fields, Line: tok.Line}, nil
}

func (p *Parser) parseFnLit() (*ast.FnLit, error) {
	tok := p.advance() // 'fn'
	params, err := p.parseParams()
	if err != nil {
		return nil, err
	}
	var retType string
	p.skipNL()
	if p.peek().Type == lexer.ARROW {
		p.advance()
		p.skipNL()
		retTok := p.advance()
		retType = retTok.Literal
	}
	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	return &ast.FnLit{Params: params, ReturnType: retType, Body: body, Line: tok.Line}, nil
}

func (p *Parser) parseNew() (*ast.NewExpr, error) {
	tok := p.advance() // 'new'
	p.skipNL()
	nameTok, err := p.expect(lexer.IDENT)
	if err != nil {
		return nil, err
	}
	p.skipNL()
	if _, err := p.expect(lexer.LBRACE); err != nil {
		return nil, err
	}
	var fields []ast.ObjectField
	p.skipNL()
	for p.peek().Type != lexer.RBRACE && p.peek().Type != lexer.EOF {
		keyTok, err := p.expect(lexer.IDENT)
		if err != nil {
			return nil, err
		}
		p.skipNL()
		if _, err := p.expect(lexer.COLON); err != nil {
			return nil, err
		}
		p.skipNL()
		val, err := p.parseExpr(0)
		if err != nil {
			return nil, err
		}
		fields = append(fields, ast.ObjectField{Key: keyTok.Literal, Value: val})
		p.skipNL()
		if p.peek().Type == lexer.COMMA {
			p.advance()
			p.skipNL()
		}
	}
	if _, err := p.expect(lexer.RBRACE); err != nil {
		return nil, err
	}
	return &ast.NewExpr{TypeName: nameTok.Literal, Fields: fields, Line: tok.Line}, nil
}
