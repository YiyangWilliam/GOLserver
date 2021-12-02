package stubs

var GOLNextState = "GOLOperations.NextState"



type Response struct {
	ChangeStateList [][]int
}

type Request struct {
	Width int
	Height int
	Board [][]uint8
	Start int
	End int
}

