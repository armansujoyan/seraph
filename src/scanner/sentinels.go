package scanner

var (
	ProgramToken    = Token{Category: "term", Value: "program"}
	SemicolonToken  = Token{Category: "term", Value: ";"}
	CommaToken      = Token{Category: "term", Value: ","}
	VarToken        = Token{Category: "term", Value: "var"}
	BeginToken      = Token{Category: "term", Value: "begin"}
	EndToken        = Token{Category: "term", Value: "end"}
	DotToken        = Token{Category: "term", Value: "."}
	ColonToken      = Token{Category: "term", Value: ":"}
	AssignmentToken = Token{Category: "term", Value: ":="}
	IntegerToken    = Token{Category: "term", Value: "integer"}
	StringToken     = Token{Category: "term", Value: "string"}
	PlusToken       = Token{Category: "term", Value: "+"}
	MinusToken      = Token{Category: "term", Value: "-"}
)
