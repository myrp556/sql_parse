package main

import (
    "log"
    "fmt"
    sqlParser "github.com/myrp556/sql_parse/parser"
)

func main() {
    parser := sqlParser.SQLParser{}
    result, err := parser.Parse("SELECT * FROM users WHERE a>1 AND (b<33 OR c=1) OR (d BETWEEN 2 AND 4) OR (e IN (2, 3,));")
    //result, err := parser.Parse("SELECT * FROM users WHERE ( a>1 );")
    if err == nil {
        for _, c := range result {
            switch c.Type {
            case sqlParser.SelectQuery:
                log.Println(fmt.Sprintf("%v", c.Select.Content()))
            }
        }
    } else {
        log.Println(fmt.Sprintf("%v", err))
    }
}
