package parser

import (
    "strings"
    "fmt"
    //"log"
)

type UpdateStmt struct {
    Table string
    Set []RowUnit
    Where *QueryExpNode
}

func (stmt *UpdateStmt) Content() string {
    ret := fmt.Sprintf("update %s", stmt.Table)
    if len(stmt.Set) > 0 {
        ret = ret +" set ("
        for _, unit := range stmt.Set {
            ret += unit.Content() + ", "
        }
        ret = ret + ")"
    }
    if stmt.Where != nil {
        ret = ret + " where " + stmt.Where.Content()
    }

    return ret
}

func (parser *SQLParser) parseUpdateStmt() (UpdateStmt, error) {
    var table string
    set := []RowUnit{}
    var where *QueryExpNode

    loop:
    for ;!parser.emptyToken(); {
        switch strings.ToUpper(parser.peekToken()) {
        case "UPDATE":
            parser.popToken()
            table = parser.popToken()

            if len(table) == 0 {
                return UpdateStmt{}, ErrNoTableSpe
            }
        case "SET":
            parser.popToken()
            sets, _ := parser.popTokenUntil("WHERE")
            if len(sets) == 0{
                return UpdateStmt{}, ErrNoValueSpe
            }
            
            //log.Println(fmt.Sprintf("%v", sets))

            if node, err:=parseExpression(sets); err!=nil {
                return UpdateStmt{}, err
            } else {
                if node.Type == ExpList {
                    for _, n:=range node.List {
                        if unit, err1 := getRowUnitFromNode(n); err1!=nil {
                            return UpdateStmt{}, err1
                        } else {
                            set = append(set, unit)
                        }
                    }
                } else {
                    if unit, err1 := getRowUnitFromNode(node); err1!=nil {
                        return UpdateStmt{}, err1
                    } else {
                        set = append(set, unit)
                    }
                }
            }
        case "WHERE":
            if node, err:=parser.parseWhere(); err!=nil {
                return UpdateStmt{}, err
            } else {
                where = node
            }
        case "":
            break loop
        default:
            return UpdateStmt{}, ErrUnsupported
        }
    }

    stmt := UpdateStmt{}
    stmt.Table = table
    stmt.Set = set
    stmt.Where = where

    return stmt, nil
}
