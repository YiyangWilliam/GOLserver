
package main

import (
	"flag"
	"math/rand"
	"net"
	"net/rpc"
	"time"
	"uk.ac.bris.cs/solutions/distributed2/secretstrings/stubs"
)

type request struct {
	board [][]uint8
	i int
}
var threadUsed=0
var maxThread=1
var workDone =make(chan bool)
var rowRequest =make(chan request)
var change =make(chan []int)
var changeList [][]int

type GOLOperations struct {}

func (s *GOLOperations) NextState(req stubs.Request, res *stubs.Response) (err error) {
	go calculatePartState(req.Board,req.Start,req.End)
	waitAllDown()
	//fmt.Println(changeList)
	res.ChangeStateList=changeList
	changeList=nil
	return
}


func main(){
	pAddr := flag.String("port","8030","Port to listen on")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	rpc.Register(&GOLOperations{})
	listener, _ := net.Listen("tcp", ":"+*pAddr)
	defer listener.Close()
	rpc.Accept(listener)
}

func calculateRowState(board [][]uint8,i int,a bool){
	col:=len(board[0])
	row:=len(board)
	way := [3]int{0, 1, -1}
	for j := 0; j < col; j++ {
		aliveNum := 0
		// select each direction
		for x := 0; x < 3; x++ {
			for y := 0; y < 3; y++ {
				if way[x] != 0 || way[y] != 0 {
					r := way[x] + i
					c := way[y] + j
					if r<0 && c<0{if board[row-1][col-1]==255{aliveNum++}}
					if r<0 && (c >= 0 && c < col){if board[row-1][c]==255{aliveNum++}}
					if r<0 && c >= col{if board[row-1][0]==255{aliveNum++}}
					if (r >= 0 && r < row) && c<0{if board[r][col-1]==255{aliveNum++}}
					if (r >= 0 && r < row) && (c >= 0 && c < col){if board[r][c]==255{aliveNum++}}
					if (r >= 0 && r < row) && c >= col{if board[r][0]==255{aliveNum++}}
					if r>=row && c<0{if board[0][col-1]==255{aliveNum++}}
					if r>=row && (c >= 0 && c < col){if board[0][c]==255{aliveNum++}}
					if r>=row && c >= col{if board[0][0]==255{aliveNum++}}
				}
			}
		}
		// check status
		// alive ->dead
		if (board[i][j] == 255) && (aliveNum < 2 || aliveNum > 3) {
			change<-[]int{i,j,1}
		}
		// dead ->alive
		if (board[i][j] == 0) && (aliveNum == 3) {
			change<-[]int{i,j,0}
		}
	}
	if a{workDone<-true}
}

func calculatePartState(board [][]uint8,start int,end int) {
	threadUsed=1
	for i := start; i <=end ; i++ {
		if threadUsed < maxThread {
			rowRequest <- request{board: board, i: i}
		} else {
			calculateRowState(board, i, false)
		}
	}
	workDone<-true
}

func waitAllDown()  {
	for {
		select {
		case b:=<-change:
			changeList=append(changeList,b)
		case golParam := <- rowRequest:
			threadUsed++
			go calculateRowState(golParam.board,golParam.i,true)
		case <-workDone:
			threadUsed--
			if threadUsed==0{return}
		}
	}
}