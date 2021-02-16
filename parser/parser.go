package parser

import (
    "errors"
    "strings"
)

type QueryType string

const (
    SelectQuery     QueryType   = "SELECT"
    InsertQuery     QueryType   = "INSERT"
    DeleteQuery     QueryType   = "DELETE"
    UpdateQuery     QueryType   = "UPDATE"
    UnknownQuery    QueryType   = "UNKNOWN"
)

type QueryResult struct {
    Type QueryType

    Select SelectStmt
}

type SQLParser struct {
    queryRaw string
    tokensRaw []string
}

func NewSqlParser() *SQLParser {
    parser := &SQLParser{}
    return parser
}

func (parser *SQLParser) emptyToken() bool {
    return len(parser.tokensRaw)==0
}

func (parser *SQLParser) popToken() string {
    if parser.emptyToken() {
        return ""
    }
    token := parser.tokensRaw[0]
    parser.tokensRaw = parser.tokensRaw[1:]

    return token
}

func (parser *SQLParser) peekToken() string {
    if parser.emptyToken() {
        return ""
    }

    return parser.tokensRaw[0]
}

func (parser *SQLParser) popTokenUntil(pattern string) ([]string, error) {
    ret := []string{}
    i := 0
    for ; i<len(parser.tokensRaw); i++ {
        if parser.tokensRaw[i]==pattern {
            break
        }
    }

    if i>0 {
        ret = parser.tokensRaw[:i]
        if i<len(parser.tokensRaw) {
            parser.tokensRaw = parser.tokensRaw[i:]
        } else {
            parser.tokensRaw = []string {}
        }
    }

    return ret, nil
}

func (parser *SQLParser) parseToken(query string) error {
    tokens := strings.Fields(strings.Trim(query, " "))
    parser.tokensRaw = tokens
    return nil
}

func (parser *SQLParser) parseFields(fields []string) ([]string, error) {
    return strings.Split(strings.ReplaceAll(strings.Join(fields, ""),  " ", ""), ","), nil
}

func (parser *SQLParser) parse(query string) (QueryResult, error) {
    result := QueryResult {}

    err := parser.parseToken(query)
    if err !=nil {
        return QueryResult{}, err
    }

    switch strings.ToUpper(parser.peekToken()) {
    case "SELECT":
        result.Type = SelectQuery
        if stmt, err:=parser.parseSelectStmt(); err!=nil {
            return QueryResult{}, err
        } else {
            result.Select = stmt
        }
    case "INSERT":
    case "DELETE":
    case "UPDATE":
    default:
        //return QueryResult{}, ErrUnknownQuery
    }

    return result, nil
}

func (parser *SQLParser) Parse(queryRaw string) ([]QueryResult, error) {
    ret := []QueryResult {}

    querys := strings.Split(queryRaw, ";")
    for _, query := range querys {
        result, err := parser.parse(query)
        if err != nil {
            return []QueryResult{}, err
        }

        ret = append(ret, result)
    }

    return ret, nil
}

var (
    ErrUnknownQuery = errors.New("unknown query type")
    ErrInvalidQuery = errors.New("invalid query string")
    ErrUnsupported = errors.New("unsupported expression")
)
