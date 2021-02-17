package parser

import (
//    "strings"
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

func (parser *SQLParser) parseSelectFields() ([]string, error) {
    fields, err := parser.popTokenUntil("FROM")
    if err!=nil {
        return []string{}, err
    }

    fields, err = parser.parseFields(fields)
    if err != nil {
        return []string{}, err
    }

    return fields, nil
}

func (parser *SQLParser) parseSelectStmt() (SelectStmt, error) {
    //parser.tokensRaw = parser.tokensRaw[1:]
    parser.popToken()
    stmt := SelectStmt {}

    fields, err := parser.parseSelectFields()
    if err != nil {
        return SelectStmt{}, err
    }

    if parser.peekToken() != "FROM" {
        return SelectStmt{}, ErrInvalidQuery
    }
    parser.popToken()
    table := parser.popToken()

    loop:
    for ;!parser.emptyToken(); {
        switch parser.peekToken() {
        case "WHERE":
            node, err := parser.parseSelectWhere()
            if err != nil {
                return SelectStmt{}, err
            }
            stmt.Where = node
        case "":
            break loop
        default:
            return SelectStmt{}, ErrUnsupported
        }
    }

    stmt.Fields = fields
    stmt.Table =table 
    return stmt, nil
}
