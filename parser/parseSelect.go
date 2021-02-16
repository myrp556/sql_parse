package parser

import (
//    "strings"
//    "log"
    "fmt"
)

type SelectStmt struct {

    Fields []string
    From string
    Where *QueryExpNode
}

func (stmt *SelectStmt) Content() string {
    ret := fmt.Sprintf("select %s from %v", stmt.Fields, stmt.From)
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

func (parser *SQLParser) parseSelectWhere() (*QueryExpNode, error) {
    tokens := []string {}
    loop:
    for ;!parser.emptyToken(); {
        switch parser.peekToken() {
        case "", "LIMIT", "GROUP", "BY", "ORDER":
            break loop
        }
        tokens = append(tokens, parser.popToken())
    }

    return parseExpression(tokens[1:])
}

func (parser *SQLParser) parseSelectStmt() (SelectStmt, error) {
    parser.tokensRaw = parser.tokensRaw[1:]
    stmt := SelectStmt {}

    fields, err := parser.parseSelectFields()
    if err != nil {
        return SelectStmt{}, err
    }

    if parser.peekToken() != "FROM" {
        return SelectStmt{}, ErrInvalidQuery
    }
    parser.popToken()
    from := parser.popToken()

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
    stmt.From = from
    return stmt, nil
}
