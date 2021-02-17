package main

import (
    "log"
    "fmt"
    sqlParser "github.com/myrp556/sql_parse/parser"
)

func main() {
    parser := sqlParser.SQLParser{}
    result, err := parser.Parse("SELECT * FROM users WHERE a>1 AND (b<33 OR c=1) OR (d BETWEEN 2 AND 4) OR (e IN (2, 3,));")
    //result, err := parser.Parse("INSERT INTO table0 ('name', 'age', 'score') VALUES ('name1', 2, 4);")
    //result, err := parser.Parse("UPDATE Person SET name='you', city='bj' WHERE age<30 AND age>10;")
    if err == nil {
        for _, c := range result {
            switch c.Type {
            case sqlParser.SelectQuery:
                log.Println(fmt.Sprintf("%v", c.Select.Content()))
            case sqlParser.InsertQuery:
                log.Println(fmt.Sprintf("%v", c.Insert.Content()))
            case sqlParser.UpdateQuery:
                log.Println(fmt.Sprintf("%v", c.Update.Content()))
            case sqlParser.DeleteQuery:
                log.Println(fmt.Sprintf("%v", c.Delete.Content()))
            }
        }
    } else {
        log.Println(fmt.Sprintf("%v", err))
    }
}
