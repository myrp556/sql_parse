package parser

import (
    "fmt"
    "strconv"
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
        return fmt.Sprintf("%s: %v", unit.Col, unit.IntValue)
    case StrRowUnit:
        return fmt.Sprintf("%s: '%v'", unit.Col, unit.StrValue)
    case UnknownRowUnit:
        return fmt.Sprintf("UNKNOWN")
    default:
        return "?"
    }
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

func getRowUnitFromNode(node *QueryExpNode) (RowUnit, error) {
    if node == nil {
        return RowUnit{}, ErrInvalidExpression
    }

    if node.Type==ExpExpression && 
        node.Left!=nil && node.Right!=nil && node.Left.Type==ExpField && node.Op==OpEqual {
        unit := RowUnit{}
        valNode := node.Right

        unit.Col = node.Left.Field
        switch valNode.Type {
        case ExpIntValue:
            unit.Type = IntRowUnit
            unit.IntValue = valNode.IntValue
        case ExpStrValue:
            unit.Type = StrRowUnit
            unit.StrValue = valNode.StrValue
        default:
            unit.Type = UnknownRowUnit
        }

        return unit, nil
    } else {
        return RowUnit{}, ErrInvalidExpression
    }
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

func reverse(list []int) []int {
    for i,j:=0, len(list)-1; i<j; i,j=i+1, j-1 {
        list[i], list[j] = list[j], list[i]
    }

    return list
}
