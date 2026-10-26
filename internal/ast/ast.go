package ast

import "fmt"

type Node interface {
	nodeType() string
	String() string
}

type Statement interface {
	Node
	stmtNode()
}

type Expression interface {
	Node
	exprNode()
}

// ─── Program ───────────────────────────────────────────────────────────────

type Program struct {
	Statements []Statement
}

func (p *Program) nodeType() string { return "Program" }
func (p *Program) String() string   { return "<program>" }

// ─── Statements ────────────────────────────────────────────────────────────

type VarDecl struct {
	Name     string
	TypeHint string // "int", "float", "string", "bool", "any", or ""
	IsConst  bool
	Value    Expression
	Line     int
}

func (v *VarDecl) stmtNode()      {}
func (v *VarDecl) nodeType() string { return "VarDecl" }
func (v *VarDecl) String() string {
	kw := "let"
	if v.IsConst {
		kw = "const"
	}
	return fmt.Sprintf("%s %s = %s", kw, v.Name, v.Value)
}

type AssignStmt struct {
	Target Expression // Ident or IndexExpr or MemberExpr
	Op     string     // "=" "+=" "-="
	Value  Expression
	Line   int
}

func (a *AssignStmt) stmtNode()      {}
func (a *AssignStmt) nodeType() string { return "AssignStmt" }
func (a *AssignStmt) String() string {
	return fmt.Sprintf("%s %s %s", a.Target, a.Op, a.Value)
}

type ReturnStmt struct {
	Value Expression
	Line  int
}

func (r *ReturnStmt) stmtNode()      {}
func (r *ReturnStmt) nodeType() string { return "ReturnStmt" }
func (r *ReturnStmt) String() string { return fmt.Sprintf("return %s", r.Value) }

type BreakStmt struct{ Line int }

func (b *BreakStmt) stmtNode()      {}
func (b *BreakStmt) nodeType() string { return "BreakStmt" }
func (b *BreakStmt) String() string { return "break" }

type ContinueStmt struct{ Line int }

func (c *ContinueStmt) stmtNode()      {}
func (c *ContinueStmt) nodeType() string { return "ContinueStmt" }
func (c *ContinueStmt) String() string { return "continue" }

type LogStmt struct {
	Args []Expression
	Line int
}

func (l *LogStmt) stmtNode()      {}
func (l *LogStmt) nodeType() string { return "LogStmt" }
func (l *LogStmt) String() string { return "log(...)" }

type ExprStmt struct {
	Expr Expression
	Line int
}

func (e *ExprStmt) stmtNode()      {}
func (e *ExprStmt) nodeType() string { return "ExprStmt" }
func (e *ExprStmt) String() string { return e.Expr.String() }

type BlockStmt struct {
	Statements []Statement
}

func (b *BlockStmt) stmtNode()      {}
func (b *BlockStmt) nodeType() string { return "BlockStmt" }
func (b *BlockStmt) String() string { return "{...}" }

type IfStmt struct {
	Condition   Expression
	Consequence *BlockStmt
	ElseIfs     []ElseIfClause
	Alternative *BlockStmt
	Line        int
}

type ElseIfClause struct {
	Condition Expression
	Body      *BlockStmt
}

func (i *IfStmt) stmtNode()      {}
func (i *IfStmt) nodeType() string { return "IfStmt" }
func (i *IfStmt) String() string { return "if ..." }

type WhileStmt struct {
	Condition Expression
	Body      *BlockStmt
	Line      int
}

func (w *WhileStmt) stmtNode()      {}
func (w *WhileStmt) nodeType() string { return "WhileStmt" }
func (w *WhileStmt) String() string { return "while ..." }

type ForInStmt struct {
	Var      string
	Iterable Expression
	Body     *BlockStmt
	Line     int
}

func (f *ForInStmt) stmtNode()      {}
func (f *ForInStmt) nodeType() string { return "ForInStmt" }
func (f *ForInStmt) String() string { return "for ... in ..." }

type FnDecl struct {
	Name       string
	Params     []Param
	ReturnType string
	Body       *BlockStmt
	Line       int
}

type Param struct {
	Name     string
	TypeHint string
}

func (f *FnDecl) stmtNode()      {}
func (f *FnDecl) nodeType() string { return "FnDecl" }
func (f *FnDecl) String() string { return fmt.Sprintf("fn %s(...)", f.Name) }

type TypeDecl struct {
	Name   string
	Fields []TypeField
	Line   int
}

type TypeField struct {
	Name     string
	TypeHint string
	Optional bool
}

func (t *TypeDecl) stmtNode()      {}
func (t *TypeDecl) nodeType() string { return "TypeDecl" }
func (t *TypeDecl) String() string { return fmt.Sprintf("type %s {...}", t.Name) }

type ImportStmt struct {
	Names  []string // imported symbols, empty = import all
	Module string   // module path or file
	Line   int
}

func (i *ImportStmt) stmtNode()      {}
func (i *ImportStmt) nodeType() string { return "ImportStmt" }
func (i *ImportStmt) String() string { return fmt.Sprintf("import from \"%s\"", i.Module) }

// ─── Expressions ───────────────────────────────────────────────────────────

type IntLit struct {
	Value int64
	Line  int
}

func (i *IntLit) exprNode()       {}
func (i *IntLit) nodeType() string { return "IntLit" }
func (i *IntLit) String() string { return fmt.Sprintf("%d", i.Value) }

type FloatLit struct {
	Value float64
	Line  int
}

func (f *FloatLit) exprNode()       {}
func (f *FloatLit) nodeType() string { return "FloatLit" }
func (f *FloatLit) String() string { return fmt.Sprintf("%g", f.Value) }

type StringLit struct {
	Value string
	Line  int
}

func (s *StringLit) exprNode()       {}
func (s *StringLit) nodeType() string { return "StringLit" }
func (s *StringLit) String() string { return fmt.Sprintf("%q", s.Value) }

type BoolLit struct {
	Value bool
	Line  int
}

func (b *BoolLit) exprNode()       {}
func (b *BoolLit) nodeType() string { return "BoolLit" }
func (b *BoolLit) String() string { return fmt.Sprintf("%t", b.Value) }

type NullLit struct{ Line int }

func (n *NullLit) exprNode()       {}
func (n *NullLit) nodeType() string { return "NullLit" }
func (n *NullLit) String() string { return "null" }

type Identifier struct {
	Name string
	Line int
}

func (i *Identifier) exprNode()       {}
func (i *Identifier) nodeType() string { return "Identifier" }
func (i *Identifier) String() string { return i.Name }

type BinaryExpr struct {
	Left  Expression
	Op    string
	Right Expression
	Line  int
}

func (b *BinaryExpr) exprNode()       {}
func (b *BinaryExpr) nodeType() string { return "BinaryExpr" }
func (b *BinaryExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Left, b.Op, b.Right)
}

type UnaryExpr struct {
	Op      string
	Operand Expression
	Line    int
}

func (u *UnaryExpr) exprNode()       {}
func (u *UnaryExpr) nodeType() string { return "UnaryExpr" }
func (u *UnaryExpr) String() string { return fmt.Sprintf("(%s%s)", u.Op, u.Operand) }

type CallExpr struct {
	Callee Expression
	Args   []Expression
	Line   int
}

func (c *CallExpr) exprNode()       {}
func (c *CallExpr) nodeType() string { return "CallExpr" }
func (c *CallExpr) String() string { return fmt.Sprintf("%s(...)", c.Callee) }

type MemberExpr struct {
	Object   Expression
	Property string
	Line     int
}

func (m *MemberExpr) exprNode()       {}
func (m *MemberExpr) nodeType() string { return "MemberExpr" }
func (m *MemberExpr) String() string { return fmt.Sprintf("%s.%s", m.Object, m.Property) }

type IndexExpr struct {
	Object Expression
	Index  Expression
	Line   int
}

func (i *IndexExpr) exprNode()       {}
func (i *IndexExpr) nodeType() string { return "IndexExpr" }
func (i *IndexExpr) String() string { return fmt.Sprintf("%s[%s]", i.Object, i.Index) }

type ArrayLit struct {
	Elements []Expression
	Line     int
}

func (a *ArrayLit) exprNode()       {}
func (a *ArrayLit) nodeType() string { return "ArrayLit" }
func (a *ArrayLit) String() string { return "[...]" }

type ObjectLit struct {
	Fields []ObjectField
	Line   int
}

type ObjectField struct {
	Key   string
	Value Expression
}

func (o *ObjectLit) exprNode()       {}
func (o *ObjectLit) nodeType() string { return "ObjectLit" }
func (o *ObjectLit) String() string { return "{...}" }

type FnLit struct {
	Params     []Param
	ReturnType string
	Body       *BlockStmt
	Line       int
}

func (f *FnLit) exprNode()       {}
func (f *FnLit) nodeType() string { return "FnLit" }
func (f *FnLit) String() string { return "fn(...) {...}" }

type NewExpr struct {
	TypeName string
	Fields   []ObjectField
	Line     int
}

func (n *NewExpr) exprNode()       {}
func (n *NewExpr) nodeType() string { return "NewExpr" }
func (n *NewExpr) String() string { return fmt.Sprintf("new %s {...}", n.TypeName) }
