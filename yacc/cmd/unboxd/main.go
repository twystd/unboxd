package main

import (
    "bufio"
    "bytes"
    "io"
    "log"
    "os"
)

type Parser struct {
    result     float64

    // the following fields are used by yyParse to reduce allocation.
    cache  []yySymType
    yylval yySymType
    yyVAL  *yySymType
}

func main() {


    input := bufio.NewReader(os.Stdin)
    for {
        line, _, err := input.ReadLine()
        if err == io.EOF {
            break
        }

        lexer := NewLexer(bytes.NewReader(line))

        parser := &Parser{
            cache: make([]yySymType, 200),
        }

        // yyParse(lexer, parser)
        yyParse(lexer)

        log.Printf("result=%.2f", parser.result)
    }

}