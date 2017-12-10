// This file holds the main logic of the BLM parser

package main

import(
    "blm/structs"
	"blm/utils"
	"fmt"
	//"os"
	//"strings"
)

// Struct definition for a node in the binary tree that
// is returned by the parse function
//type Node struct {
//    Label string
//    Form string
//    Left *Node
//    Right *Node
//    Features []string
//}

//type State struct {
//    Node   *Node
//    Stream []string
//}

var (
    lex  map[string][]string
    feat map[string][]string
    null map[string][]string
)

func parse(sentence string) *structs.Node {
    return &structs.Node{"TP", "", &structs.Node{"T", "lmao", nil, nil, nil}, &structs.Node{"vP", "", nil, nil, nil}, nil}
}

func main() {
    lex = make(map[string][]string)
    feat = make(map[string][]string)
    null = make(map[string][]string)
    utils.Init_map(lex, "../lexicon.txt")
    utils.Init_map(feat, "../features.txt")
    utils.Init_map(null, "../null.txt")
    utils.PrintMap([]map[string][]string{lex, feat, null})
    sentence := utils.Wait_user()
	if sentence == "" {
		fmt.Println("Error listening to sentence.")
		return
	}
    soln := parse(sentence)
    if soln == nil {
        fmt.Println("No tree can be formed for the sentence:")
        fmt.Println("    " + sentence)
        return
    }
    fmt.Println("\\Tree")
    utils.Latex(soln, 0)
}
