package lexer

type TokenType string

const (
	// Literals
	INT_LIT    TokenType = "INT_LIT"
	FLOAT_LIT  TokenType = "FLOAT_LIT"
	STRING_LIT TokenType = "STRING_LIT"
	BOOL_LIT   TokenType = "BOOL_LIT"
	NULL_LIT   TokenType = "NULL_LIT"
	IDENT      TokenType = "IDENT"

	// Keywords
	LET      TokenType = "let"
	CONST    TokenType = "const"
	FN       TokenType = "fn"
	RETURN   TokenType = "return"
	IF       TokenType = "if"
	ELSE     TokenType = "else"
	ELSEIF   TokenType = "elseif"
	WHILE    TokenType = "while"
	FOR      TokenType = "for"
	IN       TokenType = "in"
	BREAK    TokenType = "break"
	CONTINUE TokenType = "continue"
	LOG      TokenType = "log"
	IMPORT   TokenType = "import"
	FROM     TokenType = "from"
	TYPE     TokenType = "type"
	STRUCT   TokenType = "struct"
	INTERFACE TokenType = "interface"
	NEW      TokenType = "new"
	TRUE     TokenType = "true"
	FALSE    TokenType = "false"
	NULL     TokenType = "null"

	// Types
	T_INT    TokenType = "int"
	T_FLOAT  TokenType = "float"
	T_STRING TokenType = "string"
	T_BOOL   TokenType = "bool"
	T_ANY    TokenType = "any"
	T_VOID   TokenType = "void"

	// Operators
	ASSIGN   TokenType = "="
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	STAR     TokenType = "*"
	SLASH    TokenType = "/"
	PERCENT  TokenType = "%"
	BANG     TokenType = "!"
	EQ       TokenType = "=="
	NEQ      TokenType = "!="
	LT       TokenType = "<"
	GT       TokenType = ">"
	LTE      TokenType = "<="
	GTE      TokenType = ">="
	AND      TokenType = "&&"
	OR       TokenType = "||"
	PLUS_ASSIGN  TokenType = "+="
	MINUS_ASSIGN TokenType = "-="
	ARROW    TokenType = "->"

	// Delimiters
	LPAREN    TokenType = "("
	RPAREN    TokenType = ")"
	LBRACE    TokenType = "{"
	RBRACE    TokenType = "}"
	LBRACKET  TokenType = "["
	RBRACKET  TokenType = "]"
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"
	COLON     TokenType = ":"
	DOT       TokenType = "."
	NEWLINE   TokenType = "NEWLINE"

	// Special
	EOF     TokenType = "EOF"
	ILLEGAL TokenType = "ILLEGAL"
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Col     int
}

var keywords = map[string]TokenType{
	"let":       LET,
	"const":     CONST,
	"fn":        FN,
	"return":    RETURN,
	"if":        IF,
	"else":      ELSE,
	"elseif":    ELSEIF,
	"while":     WHILE,
	"for":       FOR,
	"in":        IN,
	"break":     BREAK,
	"continue":  CONTINUE,
	"log":       LOG,
	"import":    IMPORT,
	"from":      FROM,
	"type":      TYPE,
	"struct":    STRUCT,
	"interface": INTERFACE,
	"new":       NEW,
	"true":      TRUE,
	"false":     FALSE,
	"null":      NULL,
	"int":       T_INT,
	"float":     T_FLOAT,
	"string":    T_STRING,
	"bool":      T_BOOL,
	"any":       T_ANY,
	"void":      T_VOID,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
