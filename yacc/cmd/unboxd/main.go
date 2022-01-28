package main

import (
	"fmt"
)

func main() {
	fmt.Println("YACC'ING")
	l := &lex{
		// {key1 = value1 | key2 = {key3 = value3} | key4 = {key5 = { key6 = value6 }}}
		[]token{
			{'{', ""},
			{KEY, "key1"},
			{'=', ""},
			{VAL, "value1"},
			{'|', ""},
			{KEY, "key2"},
			{'=', ""},
			{'{', ""},
			{KEY, "key3"},
			{'=', ""},
			{VAL, "value3"},
			{'}', ""},
			{'|', ""},
			{KEY, "key4"},
			{'=', ""},
			{'{', ""},
			{KEY, "key5"},
			{'=', ""},
			{'{', ""},
			{KEY, "key6"},
			{'=', ""},
			{VAL, "value6"},
			{'}', ""},
			{'}', ""},
			{'}', ""},
		},
		map[interface{}]interface{}{},
	}

	yyParse(l)

	fmt.Println(l.m)
}
