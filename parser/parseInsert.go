package parser

import (
    "strings"
    "strconv"
    "fmt"
    "log"
)

type UnitType string

const (
    IntRowUnit     UnitType = "IntRowUnit"
    StrRowUnit     UnitType = "StrRowUnit"
    UnknownRowUnit UnitType = "KnownRowUnit"
)

type RowUnit struct {
    Type UnitType
    Col string
    IntValue  int64
    StrValue  string
}

func (unit *RowUnit) Content() string {
    switch unit.Type {
    case IntRowUnit:
        return fmt.Sprintf("'%s': %v", unit.Col, unit.IntValue)
    case StrRowUnit:
        return fmt.Sprintf("'%s': '%v'", unit.Col, unit.StrValue)
    case UnknownRowUnit:
        return fmt.Sprintf("UNKNOWN")
    default:
        return "?"
    }
}

type InsertStmt struct {
    Table string
    Lists [][]RowUnit
}

func (stmt *InsertStmt) Content() string {
    ret := "select " + stmt.Table + " ("
    for _, lis := range stmt.Lists {
        str := "("
        for _, unit := range lis {
            str = str + unit.Content() + ", "
        }
        str = str+"), "
        ret = ret + str
    }
    ret = ret +")"

    return ret
}

func getRowUnit(i int, colNode *QueryExpNode, valNode *QueryExpNode) RowUnit {
    unit := RowUnit{}

    if colNode!=nil {
        switch colNode.Type {
        case ExpStrValue:
            unit.Col = colNode.StrValue
        case ExpIntValue:
            unit.Col = strconv.Itoa(int(colNode.IntValue))
        default:
            unit.Col = strconv.Itoa(i)
        }
    } else {
        unit.Col = strconv.Itoa(i)
    }

    if valNode!=nil {
        switch valNode.Type {
        case ExpStrValue:
            unit.Type = StrRowUnit
            unit.StrValue = valNode.StrValue
        case ExpIntValue:
            unit.Type = IntRowUnit
            unit.IntValue = valNode.IntValue
        default:
            unit.Type = UnknownRowUnit
        }
    } else {
        unit.Type = UnknownRowUnit
    }

    return unit
}

func parseInsertRow(cols []string, vals []string) ([]RowUnit, error) {
    lis := []RowUnit {}
    var colNode *QueryExpNode
    var valNode *QueryExpNode
    var err error

    if len(cols)>0 {
        if colNode, err=parseExpression(cols); err!=nil {
            return nil, err
        }
    }

    if len(vals)>0 {
        if valNode, err=parseExpression(vals); err!=nil {
            return nil, err
        }
    }

    if valNode==nil || (colNode!=nil &&
        ((colNode.Type==ExpList&&valNode.Type==ExpList&&len(colNode.List)!=len(valNode.List)) ||
         (colNode.Type!=valNode.Type && (colNode.Type==ExpList || valNode.Type==ExpList)))) {

        log.Println(fmt.Sprintf("1"))
        return nil, ErrInvalidQuery
    }

    //log.Println(fmt.Sprintf("%v", valNode.Content()))

    var unit RowUnit
    if valNode.Type==ExpList {
        for i:=0; i<len(valNode.List); i++ {
            if colNode!=nil {
                unit = getRowUnit(i, colNode.List[i], valNode.List[i])
            } else {
                unit = getRowUnit(i, nil, valNode.List[i])
            }
            lis = append(lis, unit)
        }
    } else {
        if colNode!=nil {
            unit = getRowUnit(0, colNode, valNode)
        } else {
            unit = getRowUnit(0, nil, valNode)
        }
        lis = append(lis, unit)
    }

    return lis, nil
}

func (parser *SQLParser) parseInsertStmt() (InsertStmt, error) {
    insertStr := parser.popToken()
    intoStr := parser.popToken()
    if strings.ToUpper(insertStr)!="INSERT" ||
        strings.ToUpper(intoStr)!="INTO"{
        return InsertStmt{}, ErrInvalidQuery
    }


    lists := [][]RowUnit {}
    table := parser.popToken()
    if len(table)==0 {
        return InsertStmt{}, ErrNoTableSpe
    }

    cols, _ := parser.popTokenUntil("VALUES")

    valueStr := parser.popToken()
    if strings.ToUpper(valueStr)!="VALUES" {
        return InsertStmt{}, ErrInvalidQuery
    }
    vals := parser.popAllToken()

    log.Println(fmt.Sprintf("insert %v %v", cols, vals))
    if lis,err:=parseInsertRow(cols, vals); err!=nil {
        return InsertStmt{}, err
    } else {
        //log.Println(fmt.Sprintf("%v", lis))
        if len(lis)>0 {
            lists = append(lists, lis)
        }
    }

    stmt := InsertStmt {}
    stmt.Table = table
    stmt.Lists = lists

    return stmt, nil
}
