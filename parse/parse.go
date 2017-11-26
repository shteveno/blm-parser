// This file holds the main logic of the BLM parser.

package main

import(
    "fmt"
    "blm/utils"
)

// Struct definition for a node in the binary tree that
// is returned by the parse function
type Node struct {
    Form string
    Spec *Node
    Comp *Node
    Features []string 
}

// Waits for user input to process and parse into a tree.
//func wait_input(str string) {
//    fmt.Println("Hey!")
//}

var (
    lex  map[string][]string
    feat map[string][]string
    null map[string][]string
)

func main() {
    lex = make(map[string][]string)
    feat = make(map[string][]string)
    null = make(map[string][]string)
    utils.Init_map(lex, "../lexicon.txt")
    utils.Init_map(feat, "../features.txt")
    utils.Init_map(null, "../null.txt")
    fmt.Println("lex:")
    for k, v := range lex {
        fmt.Println("k is ", k, " and v is ", v, " with length ", len(v))
    }
    fmt.Println("feat:")
    for k, v := range feat {
        fmt.Println("k is ", k, " and v is ", v, " with length ", len(v))
    }
    fmt.Println("null:")
    for k, v := range null {
        fmt.Println("k is ", k, " and v is ", v, " with length ", len(v))
    }
}
