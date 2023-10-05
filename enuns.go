package ddbmoose

type TypeFilter int

const (
	TfrEqual TypeFilter = iota
	TfrNotEqual
	TfrContains
	TfrBeginsWith
	TfrGreaterThan
	TfrGreaterThanEqual
	TfrLessThan
	TfrLessThanEqual
	TfrBetween
)

type TypeLogicalOperator int

const (
	TloNone TypeLogicalOperator = iota
	TloAnd
	TloOr
)
