// This file holds the main logic of the BLM parser.

package main

import(
    "fmt"
    "blm/utils"
    "bufio"
    "os"
    "strings"
)

// Struct definition for a node in the binary tree that
// is returned by the parse function
type Node struct {
    Label string
    Form string
    Left *Node
    Right *Node
    Features []string
}

type State struct {
    Node *Node
    Stream []string
}

var (
    lex  map[string][]string
    feat map[string][]string
    null map[string][]string
)

// Waits for user input to process and parse into a tree.
func wait_user() string {
	buf := bufio.NewReader(os.Stdin)
    fmt.Print("> ")
    sentence, err := buf.ReadBytes('\n')
    if err != nil {
        fmt.Println(err)
		return ""
    } else {
        return string(sentence)
    }
}

// Attemps to merge the required specifier given the 
// word stream. Returns all possible states where
// the requested specifier can be interpreted
func spec(stream []string, uFeat string, movedFeat string) []State {
    //strings.Split(stream)
    var (
        head string
        sel string
        states []State
    )
    for i, word := range stream {
        categories := lex[word]
        for _, cat := range categories {
            cat_sel := strings.Split(cat, "_")
            head = cat_sel[0]
            sel = ""
            if len(cat_sel) == 2 {
                sel = cat_sel[1]
            }
            candidates := feat[head]
            for _, candidate := range candidates {
                bundle := strings.Split(candidate, ",")
                if bundle[0] == sel {
                    xP := &Node{}
                    xP.Label = head + "P"
                    xP.Left = &Node{head, word, nil, nil, bundle}
                    states = append(states, State{xP, stream[i+1:]})
                }
            }
        }
    }
    return states
}

func parse(sentence string) *Node {
    var (
        t *Node = &Node{}
        tBar *Node = &Node{}
    )
    tBar.Left = t
    stream := strings.Split(sentence, " ")
    t.Label = "T"
    t.Form = "$\\varnothing$"
    t.Features = []string{"FIN,*uD,uv"}
    found := spec(stream, "*uD", "")
    if found == nil {
        return nil
    }
    fmt.Println("This is what remains:", found[0].Stream)
    ret := Node{}
    ret.Label = "TP"
    ret.Left = found[0].Node
    ret.Right = tBar
    return &ret
}

// Formats the resulting tree such that it may be compiled in
// Latex.
func latex(root *Node, depth int) {
    if root == nil {
        return
    }
    offset := strings.Repeat(" ", depth*2)
    if root.Left == root.Right {
        fmt.Println(offset + "[." + root.Label + " " + root.Form + " ]")
        return
    }
    fmt.Println(offset + "[." + root.Label + " ")
    latex(root.Left, depth + 1)
    latex(root.Right, depth + 1)
    fmt.Println(offset + "]")
}

func main() {
    lex = make(map[string][]string)
    feat = make(map[string][]string)
    null = make(map[string][]string)
    utils.Init_map(lex, "../lexicon.txt")
    utils.Init_map(feat, "../features.txt")
    utils.Init_map(null, "../null.txt")
    utils.PrintMap([]map[string][]string{lex, feat, null})
    sentence := wait_user()
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
    latex(soln, 0)
}
