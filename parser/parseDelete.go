package parser

import (
    "strings"
)

type DeleteStmt struct {
    Table string
    Fields []string
    Where *QueryExpNode
}

func (stmt *DeleteStmt) Content() string {
    ret := "delete from " + stmt.Table
    ret = ret + "["
    for i, field := range stmt.Fields {
        ret = ret + field
        if i>0 {
            ret = ret + ","
        }
    }
    ret = ret + "]"
    if stmt.Where != nil {
        ret = ret + " where " + stmt.Where.Content()
    }

    return ret
}

func (parser *SQLParser) parseDeleteStmt() (DeleteStmt, error) {
    var table string
    fields := []string{}
    var where *QueryExpNode

    loop:
    for ;!parser.emptyToken(); {
        switch strings.ToUpper(parser.peekToken()) {
        case "DELETE":
            parser.popToken()
            if f, err:=parser.parseFields(); err!=nil {
                return DeleteStmt{}, err
            } else {
                fields = f
            }
        case "FROM":
            parser.popToken()
            table = parser.popToken()
            if len(table)==0 {
                return DeleteStmt{}, ErrNoTableSpe
            }
        case "WHERE":
            if node, err:=parser.parseWhere(); err!=nil {
                return DeleteStmt{}, err
            } else {
                where = node
            }
        case "":
            break loop
        default:
            return DeleteStmt{}, ErrUnsupported
        }
    }

    stmt := DeleteStmt{}
    stmt.Table = table
    stmt.Fields = fields
    stmt.Where = where

    return stmt, nil
}
