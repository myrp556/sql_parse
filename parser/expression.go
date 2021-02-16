package parser

import (
    "strconv"
    "strings"
    "log"
    "fmt"
    "errors"
)

type ExpNodeType string

const (
    ExpExpression  ExpNodeType = "Expression"
    ExpField       ExpNodeType = "Field"
    ExpIntValue    ExpNodeType = "IntValue"
    ExpStrValue    ExpNodeType = "StrValue"
    ExpList        ExpNodeType = "List"
)

type ExpOperation string

const (
    OpEqual         ExpOperation = "="
    OpNotEqual      ExpOperation = "<>"
    OpGreater       ExpOperation = ">"
    OpLess          ExpOperation = "<"
    OpGreaterEqual  ExpOperation = ">="
    OpLessEqual     ExpOperation = "<="
    OpBetween       ExpOperation = "BETWEEN"
    OpLike          ExpOperation = "LIKE"
    OpIn            ExpOperation = "IN"
    OpAnd           ExpOperation = "AND"
    OpOr            ExpOperation = "OR"
    OpUnknown       ExpOperation = "?"
)

var opPatterns = [...]ExpOperation {
    OpEqual, OpNotEqual,
    OpGreater, OpLess,
    OpGreaterEqual,OpLessEqual,
    OpBetween,OpLike,OpIn,
    OpAnd, OpOr,
}

var splitPatterns = [...]ExpOperation {
    ")", ",", "(",
    OpEqual, OpNotEqual,
    OpGreater, OpLess,
    OpGreaterEqual,OpLessEqual,
    OpBetween,OpLike,OpIn,
    OpAnd, OpOr,
}

type QueryExpNode struct {
    Type ExpNodeType
    Left *QueryExpNode
    Op ExpOperation
    Right *QueryExpNode
    Field string
    IntValue int64
    StrValue string
    List []*QueryExpNode
}

func (node *QueryExpNode) check() (bool, error) {
    switch node.Type {
    case ExpField:
        return true, nil
    case ExpIntValue, ExpStrValue:
        return false, nil
    case ExpList:
        if len(node.List)>0 {
            t := node.List[0].Type
            if t==ExpIntValue || t==ExpStrValue {
                for _, n := range node.List {
                    if n.Type != t {
                        return false, ErrInvalidList
                    }
                }
            } else {
                return false, ErrInvalidList
            }
        }

        return false, nil
    case ExpExpression:
    }

    return true, nil
}

func (node *QueryExpNode) Content() string {
    switch node.Type {
    case ExpField:
        return fmt.Sprintf("%v", node.Field)
    case ExpIntValue:
        return fmt.Sprintf("%v", node.IntValue)
    case ExpStrValue:
        return fmt.Sprintf("'%v'", node.StrValue)
    case ExpList:
        ret := "["
        for _, n := range node.List {
            ret = ret + n.Content()
        }
        return ret + "]"
    case ExpExpression:
        ret := "("+node.Left.Content()+") "+string(node.Op)+" ("+node.Right.Content()+")"
        return ret
    }

    return "?"
}

func maybeAppend(lis []string, str string) []string {
    if len(str)>0 {
        return append(lis, str)
    }
    return lis
}

func canSplit(token string) []string {
    ret := []string {}
    for {
        i := -1
        for _, pattern := range splitPatterns {
            i = strings.LastIndex(token, string(pattern))
            if i>=0 {
                ret = maybeAppend(ret, token[i+len(pattern):])
                ret = maybeAppend(ret, string(pattern))
                token = token[:i]
                break
            }
        }
        if i<0 {
            break
        }
    }

    ret = maybeAppend(ret, token)
    return ret
}

func patternToOp(token string) (ExpOperation, int) {
    switch strings.ToUpper(token) {
    case "BETWEEN":
        return OpBetween, 4
    case "LIKE":
        return OpLike, 2
    case "IN":
        return OpIn, 2
    case "AND":
        return OpAnd, 3
    case "OR":
        return OpOr, 3
    }

    switch token {
    case "=":
        return OpEqual, 1
    case "<>":
        return OpNotEqual, 1
    case ">":
        return  OpGreater, 1
    case "<":
        return OpLess, 1
    case ">=":
        return OpGreaterEqual, 1
    case "<=":
        return OpLessEqual, 1
    }

    return OpUnknown, 0
}

func isOp(token string) bool {
    for _, op := range opPatterns {
        if token==string(op) || strings.ToUpper(token)==string(op) {
            return true
        }
    }
    return false
}

func parseNode(tokens []string, leftNode bool) (*QueryExpNode, error) {
    if len(tokens)==0 {
        return nil, nil
    }

    quote := 0
    var operation ExpOperation
    opIndex := -1
    opRank := 0
    for i:=len(tokens)-1; i>=0; i-- {
        token := tokens[i]
        switch token {
        case "(":
            quote ++
        case ")":
            quote --
        case ",":
            // nothing
        default:
            if quote==0 && isOp(token) {
                op, rank := patternToOp(token)
                if rank > opRank {
                    operation = op
                    opIndex = i
                    opRank = rank
                }
            }
        }

    }

    if opIndex > 0 {
        node := &QueryExpNode{}
        node.Type = ExpExpression
        node.Op = operation
        left, err1 := parseNode(tokens[:opIndex], true)
        right, err2 := parseNode(tokens[opIndex+1:], false)
        if err1 != nil {
            return nil, err1
        }
        if err2 != nil {
            return nil, err2
        }

        node.Left = left
        node.Right = right
        return node ,nil
    } else {
        //log.Println(fmt.Sprintf("-%v", tokens))
        if len(tokens)>=2 && tokens[0]=="(" && tokens[len(tokens)-1]==")" {
            return parseNode(tokens[1:len(tokens)-1], leftNode)
        }

        node := &QueryExpNode{}
        if len(tokens) == 1 {
            token := tokens[0]
            if len(token)==0 || (token[0:1]=="'"&&token[len(token)-1:len(token)]=="'") {
                node.Type = ExpStrValue
                node.StrValue = token[1:len(token)]
            } else {
                if v,err:=strconv.ParseInt(tokens[0], 10, 64); err==nil {
                    node.Type = ExpIntValue
                    node.IntValue = v
                } else {
                    if leftNode {
                        node.Type = ExpField
                        node.Field = token
                    } else {
                        return nil, ErrInvalidExpression
                    }
                }
            }
        } else {
            node.Type = ExpList
            for _, token := range tokens {
                if token != "," {
                    n, err := parseNode([]string{token}, false)
                    if err!=nil {
                        return nil, err
                    } else {
                        node.List = append(node.List, n)
                    }
                }
            }
        }

        return node, nil
    }

}


func parseExpression(tokens []string) (*QueryExpNode, error) {
    newTokens := []string{}
    for _, token := range tokens {
        s := canSplit(token)
        //log.Println(fmt.Sprintf("%s %v", token, s))
        for i:=len(s)-1; i>=0; i-- {
            newTokens = append(newTokens, s[i])
        }
    }

    log.Println(fmt.Sprintf("parse exp %v", newTokens))
    if node, err := parseNode(newTokens, false); err!=nil {
        return nil, err
    } else {
        return node, nil
    }
}

var (
    ErrInvalidExpression = errors.New("invalid expression")
    ErrInvalidList = errors.New("invalid list")
)
