%{
package main

import (
    "fmt"
)
%}

%union {
        number float64
}

%token <number> number

%token NUMBER TOKHEAT STATE TOKTARGET TOKTEMPERATURE

%%

commands: /* empty */
        | commands command
        ;

command:
        heat_switch
        |
        target_set
        ;

heat_switch:
        TOKHEAT STATE
        {
                fmt.Printf("\tHeat turned on or off\n")
        }
        ;

target_set:
        TOKTARGET TOKTEMPERATURE NUMBER
        {
                fmt.Printf("\tTemperature set\n")
        }
        ;

