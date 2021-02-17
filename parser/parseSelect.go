package parser

import (
    "strings"
//    "log"
    "fmt"
)

type SelectStmt struct {

    Fields []string
    Table string
    Where *QueryExpNode
}

func (stmt *SelectStmt) Content() string {
    ret := fmt.Sprintf("select %s from %v", stmt.Fields, stmt.Table)
    if where:=stmt.Where.Content(); len(where)>0 {
        ret = ret + fmt.Sprintf(" where %s", where)
    }
    
    return ret
}

func (parser *SQLParser) parseSelectStmt() (SelectStmt, error) {
    fields := []string {}
    var table string
    var where *QueryExpNode

    loop:
    for ;!parser.emptyToken(); {
        switch strings.ToUpper(parser.peekToken()) {
        case "SELECT":
            parser.popToken()
            if f, err := parser.parseFields(); err!=nil {
                return SelectStmt{}, err
            } else {
                fields = f
            }
        case "FROM":
            parser.popToken()
            table := parser.popToken()
            if len(table)==0 {
                return SelectStmt{}, ErrNoTableSpe
            }
        case "WHERE":
            node, err := parser.parseWhere()
            if err != nil {
                return SelectStmt{}, err
            }
            where = node
        case "":
            break loop
        default:
            return SelectStmt{}, ErrUnsupported
        }
    }

    stmt := SelectStmt {}
    stmt.Fields = fields
    stmt.Table =table
    stmt.Where = where
    return stmt, nil
}
