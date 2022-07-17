
%{

package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
	"strings"
)

var regs = make([]int, 26)
var base int

%}

// fields inside this union end up as the fields in a structure known
// as ${PREFIX}SymType, of which a reference is passed to the lexer.
%union{
	val int
}

// any non-terminal which returns a value needs a type, which is
// really a field name in the above union struct
%type <val> expr number

// same for terminals
%token <val> DIGIT LETTER TYPE_VOID TYPE_INT TYPE_STRING

%left '|'
%left '&'
%left '+'  '-'
%left '*'  '/'  '%'
%left UMINUS      /*  supplies  precedence  for  unary  minus  */

%%

list	: /* empty */
	| list stat2 '\n'
	| list function '\n'
	;

function : type functionName '(' ')' '{'
        //instrList
        '}'
        {
            fmt.Printf( "funct\n");
        }
        ;

type : TYPE_VOID | TYPE_INT | TYPE_STRING
        ;

functionName : LETTER
        ;

instrList : /* empty */
        | instrList instr
        ;

instr :    LETTER '=' DIGIT
            {
                fmt.Printf( "affect %s = %s\n", $1, $3 );
            }
        ;

stat2	:    expr
		{
			fmt.Printf( "%d\n", $1 );
		}
	|    LETTER '=' expr
		{
			regs[$1]  =  $3
		}
	;

expr	:    '(' expr ')'
		{ $$  =  $2 }
	|    expr '+' expr
		{ $$  =  $1 + $3 }
	|    expr '-' expr
		{ $$  =  $1 - $3 }
	|    expr '*' expr
		{ $$  =  $1 * $3 }
	|    expr '/' expr
		{ $$  =  $1 / $3 }
	|    expr '%' expr
		{ $$  =  $1 % $3 }
	|    expr '&' expr
		{ $$  =  $1 & $3 }
	|    expr '|' expr
		{ $$  =  $1 | $3 }
	|    '-'  expr        %prec  UMINUS
		{ $$  = -$2  }
	|    LETTER
		{ $$  = regs[$1] }
	|    number
	;

number	:    DIGIT
		{
			$$ = $1;
			if $1==0 {
				base = 8
			} else {
				base = 10
			}
		}
	|    number DIGIT
		{ $$ = base * $1 + $2 }
	;

%%      /*  start  of  programs  */

type CalcLex struct {
	s string
	pos int
}


func (l *CalcLex) Lex(lval *HephaestusSymType) int {
	var c rune = ' '
	for c == ' ' {
		if l.pos == len(l.s) {
			return 0
		}
		c = rune(l.s[l.pos])
		l.pos += 1
	}

    if(l.pos+4 < len(l.s)&&c=='v'&&rune(l.s[l.pos+1])=='o'&&rune(l.s[l.pos+2])=='i'&&rune(l.s[l.pos+3])=='d') {
        l.pos += 4;
        return TYPE_VOID;
    } else if(l.pos+3 < len(l.s)&&c=='i'&&rune(l.s[l.pos+1])=='n'&&rune(l.s[l.pos+2])=='t') {
        l.pos += 3;
        return TYPE_INT;
    } else if(l.pos+6 < len(l.s)&&c=='s'&&rune(l.s[l.pos+1])=='t'&&rune(l.s[l.pos+2])=='r'&&
            rune(l.s[l.pos+3])=='i'&&rune(l.s[l.pos+4])=='n'&&rune(l.s[l.pos+5])=='g') {
             l.pos += 6;
             return TYPE_STRING;
         }


	if unicode.IsDigit(c) {
		lval.val = int(c) - '0'
		return DIGIT
	} else if unicode.IsLower(c) {
		lval.val = int(c) - 'a'
		return LETTER
	}
	return int(c)
}

func (l *CalcLex) Error(s string) {
	fmt.Printf("syntax error: %s\n", s)
}

func Main0(filename string) {
    var fi *bufio.Reader;
    if filename=="" {
	    fi = bufio.NewReader(os.Stdin)
	} else {
        f, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}
		defer f.Close()
		fi = bufio.NewReader(f)
	}

	for {
		var eqn string
		var ok bool

		fmt.Printf("equation: ")
		if eqn, ok = readline(fi); ok {
			HephaestusParse(&CalcLex{s: eqn})
		} else {
			break
		}
	}
}

func TrimSuffix(s, suffix string) string {
    if strings.HasSuffix(s, suffix) {
        s = s[:len(s)-len(suffix)]
    }
    return s
}

func readline(fi *bufio.Reader) (string, bool) {
    /*scanner := fi.NewScanner(file)
    for scanner.Scan() {


    }*/
	s, err := fi.ReadString('\n')
	if err != nil {
		return "", false
	}
	//fmt.Println( "line:",s,"!");
	if strings.HasSuffix(s, "\r\n") {
            s = s[:len(s)-2]+"\n"
        }
	return TrimSuffix(s, "\r"), true
}
