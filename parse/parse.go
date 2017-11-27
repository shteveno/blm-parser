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
    Node   *Node
    Stream []string
}

var (
    lex  map[string][]string
    feat map[string][]string
    null map[string][]string
    size int
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

// Attempts to merge the required complement given the
// word stream. Returns all possible states where the
// requested complement can be interpreted
func comp(stream []string, uFeat string, strongFeat string) []State {
    var (
        head    string
        sel     string
        states  []State
        target  string   =  strings.Split(uFeat, "u")[1]
        bundles   []string

    )
    for i, word := range stream {
        categories := lex[word]
        for _, cat := range categories {
            cat_sel := strings.Split(cat, "_")
            head = cat_sel[0]
            sel = "."
            if len(cat_sel) == 2 {
                sel = cat_sel[1]
            }
            if head != target {
                bundles = null[target]
            } else {
                bundles = feat[target]
            }
            for _, bundle := range bundles {
                //fmt.Println("comp:", head, bundle)
                feats := strings.Split(bundle, ",")
                if feats[0] == sel {
                    xP := &Node{}
                    xP.Label = head + "P"
                    xP.Left = &Node{head, word, nil, nil, feats}
                    wordsUsed := i+1
                    if head != target {
                        wordsUsed = i
                    }
                    if feats[2] != "." {
                        comp_XPs := comp(stream[wordsUsed:], feats[2], "")
                        if comp_XPs != nil {
                            for _, comp_XP := range comp_XPs {
                                xP.Right = comp_XP.Node
                                states = append(states, State{xP, comp_XP.Stream})
                            }
                        }
                    }
                    if feats[1] != "." {
                        spec_XPs := spec(stream[:i], feats[1])
                        if spec_XPs != nil {
                            for _, spec_XP := range spec_XPs {
                                X_bar := xP
                                X_bar.Label = head + "'"
                                xP = &Node{}
                                xP.Label = head + "P"
                                xP.Left = spec_XP.Node
                                xP.Right = X_bar
                                states = append(states, State{xP, stream[wordsUsed:]})
                                xP = X_bar
                            }
                        } else {
                            continue
                        }
                    } else {
                        states = append(states, State{xP, stream[wordsUsed:]})
                    }
                }
            }
        }
    }
    return states
}

// Attemps to merge the required specifier given the 
// word stream. Returns all possible states where
// the requested specifier can be interpreted
func spec(stream []string, uFeat string) []State {
    //strings.Split(stream)
    var (
        head      string
        sel       string
        states    []State
        target    string   =  strings.Split(uFeat, "u")[1]
        bundles   []string
        wordsUsed int

    )
    for i, word := range stream {
        categories := lex[word]
        for _, cat := range categories {
            cat_sel := strings.Split(cat, "_")
            head = cat_sel[0]
            sel = "."
            if len(cat_sel) == 2 {
                sel = cat_sel[1]
            }
            if head != target {
                bundles = null[target]
                wordsUsed = i
            } else {
                bundles = feat[target]
                wordsUsed = i + 1
            }
            for _, bundle := range bundles {
                offset := wordsUsed
                feats := strings.Split(bundle, ",")
                if feats[0] == sel {
                    xP := &Node{}
                    xP.Label = head + "P"
                    xP.Left = &Node{head, word, nil, nil, feats}
                    //wordsUsed := i+1
                    //if head != target {
                    //    wordsUsed = i
                    //}
                    if feats[2] != "." {
                        //fmt.Println("run comp on:",feats[2], "and", stream[i:])
                        comp_XPs := comp(stream[offset:], feats[2], "")
                        if comp_XPs != nil {
                            for _, comp_XP := range comp_XPs {
                                xP.Right = comp_XP.Node
                                offset += len(stream) - len(comp_XP.Stream)
                                //states = append(states, State{xP, comp_XP.Stream})
                                break
                            }
                        }
                    }
                    if feats[1] != "." {
                        spec_XPs := spec(stream[:i], feats[1])
                        if spec_XPs != nil {
                            for _, spec_XP := range spec_XPs {
                                X_bar := xP
                                X_bar.Label = head + "'"
                                xP = &Node{}
                                xP.Label = head + "P"
                                xP.Left = spec_XP.Node
                                xP.Right = X_bar
                                states = append(states, State{xP, stream[offset:]})
                                xP = X_bar
                            }
                        } else {
                            continue
                        }
                    } else {
                        if (feats[1] != "." && xP.Left == nil) || (feats[1] == "." && feats[2] != "." && xP.Right == nil) || (feats[1] != "." && feats[2] != "." && xP.Right.Right == nil) {
                            continue
                        }
                        states = append(states, State{xP, stream[offset:]})
                    }
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
    found := spec(stream, "*uD")
    if found == nil {
        return nil
    }
    fmt.Println(len(found))
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
