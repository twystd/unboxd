package main

import (
    "fmt"
    "io"
    "log"
    "strconv"
    "text/scanner"
)

type Lexer struct {
    s         scanner.Scanner
}

func NewLexer(r io.Reader) *Lexer {
    l := new(Lexer)
    l.s.Init(r)
    l.s.Mode = scanner.GoTokens
    return l
}

func (l *Lexer) Lex(lval *yySymType) int {
    tok := l.s.Scan()
    switch tok {
    case scanner.EOF:
        log.Printf("eof\n")
        return 0
    case scanner.Float:
        lval.number, _ = strconv.ParseFloat(l.s.TokenText(), 64)
        log.Printf("read float %.2f\n", lval.number)
        return number
    case scanner.Int:
        lval.number, _ = strconv.ParseFloat(l.s.TokenText(), 64)
        log.Printf("read float %.2f\n", lval.number)
        return number
    default:
        symbol := l.s.TokenText()
        log.Printf("read symbol %s\n", symbol)
        return int(tok)
    }
}

func (l *Lexer) Errorf(format string, a ...interface{}) error {
    err := fmt.Sprintf(format, a...)
    return fmt.Errorf("lexer err: %s", err)
}

func (l *Lexer) Error(err string)  {
    fmt.Printf("lexer err: %s", err)
}


func (l *Lexer) AppendError(err error) {
    log.Printf("append err: %s", err.Error())
}

func (l *Lexer) Errors() (warns []error, errs []error) {
    return warns, errs
}

