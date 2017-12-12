// This file holds the main logic of the BLM parser

package main

import(
    "blm/structs"
	"blm/utils"
	"fmt"
	"strings"
)

var (
    lex  map[string][]string
    feat map[string][]string
    null map[string][]string
)

/*  Looks for a constituent headed by the target and returns it
    along with the remaining words in the stream. */
//func search(stream []string, head string, form string, ) &structs.State {

//}

//return []*structs.State{&structs.State{&structs.Node{"DP", "I", nil, nil, nil}, []string{"like", "cake"}}}

/*  Takes in a fill uFeature and returns the lexical class,
    a selected feature, and whether or not it was strong. */
func defeat(uFeat string) (string, string, bool) {
    var (
        isStrong  bool               = uFeat[0] == '*'
        full_feat string             = strings.Split(uFeat, "u")[1]
        temp      []string           = strings.Split(full_feat, "_")
        uCat      string
        uSel      string
    )
    uCat = temp[0]
    if len(temp) == 2 {
        uSel = temp[1]
    } else {
        uSel = "."
    }
    return uCat, uSel, isStrong
}

/*  Attemps to find the required specifier headed by the 
    target, given the word stream. Returns all possible states
    where the requested specifier can be interpreted. */
func spec(stream []string, uFeat string) []*structs.State {
    if stream == nil || uFeat == "." {
        return nil
    }
    uCat, _, _ := defeat(uFeat)
    found := search(stream, "u" + uCat, nil)
    return found
}

func head(stream []string, uFeat string, moved *structs.Node) []*structs.State {
    var (
        states    []*structs.State
        bundles   []string
        wordsUsed int
    )
    uCat, uSel, isStrong := defeat(uFeat)
    if isStrong {
        bundles = feat[uCat]
        for _, bundle := range bundles {
            feats := strings.Split(bundle, ",")
            fmt.Println("*uV feat bundle:",bundle, "stream:", stream, "feats1:", feats[1], "feats2:", feats[2])
            x := &structs.Node{}
            x.Label = uCat + "_" + uSel
            x.Form = "$\\sout{" + uCat + "P}$"
            x.Features = feats
            states = append(states, &structs.State{x, stream, wordsUsed+1, feats[1], feats[2]})
        }
        return states
    }
    for i, word := range stream {
        categories := lex[word]
        for _, cat := range categories {
            cat_sel := strings.Split(cat, "_")
            h := cat_sel[0]
            sel := "."
            if len(cat_sel) == 2 {
                sel = cat_sel[1]
            }
            // You must use a null head
            if h != uCat {
                bundles = null[uCat]
                wordsUsed = i
                word = "$\\varnothing$"
                if uCat == "T" {
                    sel = "FIN"
                }
                if uCat == "v" {
                    sel = "AG"
                }
            // There is an overt head present
            } else {
                bundles = feat[uCat]
                wordsUsed = i + 1
            }
            for _, bundle := range bundles {
                feats := strings.Split(bundle, ",")
                if feats[1] != "." && wordsUsed == 0 && moved == nil {
                    continue
                }
                if feats[0] == sel {
                    x := &structs.Node{}
                    x.Label = uCat + "_" + sel
                    if feats[2][0] == '*' {
                        fmt.Println("Strong feature found: " + feats[2])
                        fmt.Println("Its features:", feats[1], feats[2])
                        fmt.Println("Remaining stream:", stream[wordsUsed + 1:])
                        word = stream[wordsUsed]
                        wordsUsed += 1
                    }
                    x.Form = word
                    x.Features = feats
                    states = append(states, &structs.State{x, stream, wordsUsed, feats[1], feats[2]})
                }
            }
        }
    }
    return states
    //return []*structs.State{&structs.State{&structs.Node{"T", "$\\varnothing_{\\textsc{pres}}$", nil, nil, nil}, []string{"like", "cake"}}}
}

func comp(stream []string, uFeat string, moved *structs.Node) []*structs.State {
    if stream == nil || uFeat == "." {
        return nil
    }
    fmt.Println("search called on:", stream, uFeat, moved)
    found := search(stream, uFeat, moved)
    return found
    //return []*structs.State{&structs.State{&structs.Node{"vP", "", &structs.Node{"v_{\\textsc{ag}}", "like", nil, nil, nil}, &structs.Node{"VP", "", &structs.Node{"DP", "", nil, nil, nil}, &structs.Node{"V", "", nil, nil, nil}, nil}, nil}, nil}}
}

func search(stream []string, uFeat string, moved *structs.Node) []*structs.State {
    if stream == nil || uFeat == "." {
        return nil
    }
    var ret []*structs.State
    heads := head(stream, uFeat, moved)
    if heads == nil {
        return nil
    }
    for _, x := range heads {
        var specifiers []*structs.State
        if moved != nil  {
            maybe_move := strings.Split(x.Spec, "u")[1] + "P"
            //fmt.Println(maybe_move)
            //fmt.Println(moved.Label)
            if maybe_move == moved.Label {
                toAdd := &structs.State{}
                toAdd.Tree =  &structs.Node{maybe_move, "", nil, nil, nil}
                toAdd.Remaining = x.Remaining
                toAdd.HeadPos = 0
                toAdd.Spec = "."
                toAdd.Comp = "."
                specifiers = append(specifiers, toAdd)
                fmt.Println("uFeat trying to be resolved by movement", x.Spec)
                utils.Latex(toAdd.Tree, 0)
            }
        }
        if x.HeadPos > len(stream) {
            continue
        }
        specifiers = append(specifiers, spec(stream[:x.HeadPos], x.Spec)...)
        if specifiers == nil {
            if x.Spec != "." {
                continue
            }
            //fmt.Println("I made a head, no specifier!")
            //utils.Latex(x.Tree, 0)
            complements := comp(stream[x.HeadPos:], x.Comp, moved)
            if complements == nil {
                if x.Comp != "." {
                    continue
                }
                toAdd := &structs.State{}
                toAdd.Tree = utils.FormTree(nil, x.Tree, nil)
                toAdd.Remaining = stream[x.HeadPos:]
                toAdd.HeadPos = 0
                toAdd.Spec = "."
                toAdd.Comp = "."
                ret = append(ret, toAdd)
                continue
            }
            for _, yP := range complements {
                //fmt.Println("OMFG")
                //utils.Latex(x.Tree, 0)
                //utils.Latex(yP.Tree, 0)
                //fmt.Println(yP.Remaining)
                if len(yP.Remaining) == 0 {
                    toAdd := &structs.State{}
                    toAdd.Tree = utils.FormTree(nil, x.Tree, yP.Tree)
                    toAdd.Remaining = stream[yP.HeadPos:]
                    toAdd.HeadPos = 0
                    toAdd.Spec = "."
                    toAdd.Comp = "."
                    ret = append(ret, toAdd)
                }
            }
        }
        for _, wP := range specifiers {
            moved = nil
            fmt.Println("I made a specifier for", uFeat)
            utils.Latex(wP.Tree, 0)
            fmt.Println("Remaining:", stream[x.HeadPos:])
            fmt.Println("Next comp feature:", x.Comp)
            if x.Spec[0] == '*' {
                moved = wP.Tree
            }
            fmt.Println("comp called on: ", stream[x.HeadPos:], x.Comp, moved)
            complements := comp(stream[x.HeadPos:], x.Comp, moved)
            if complements == nil {
                if x.Comp != "." {
                    continue
                }
                toAdd := &structs.State{}
                toAdd.Tree = utils.FormTree(wP.Tree, x.Tree, nil)
                toAdd.Remaining = []string{}
                toAdd.HeadPos = 0
                toAdd.Spec = "."
                toAdd.Comp = "."
                ret = append(ret, toAdd)
                continue
            }
            for _, yP := range complements {
                fmt.Println("Got full constituent for ", uFeat)
                fmt.Println("Remaining of string:", yP.Remaining)
                if len(yP.Remaining) == 0 {
                    toAdd := &structs.State{}
                    toAdd.Tree = utils.FormTree(wP.Tree, x.Tree, yP.Tree)
                    toAdd.Remaining = []string{}
                    toAdd.HeadPos = 0
                    toAdd.Spec = "."
                    toAdd.Comp = "."
                    ret = append(ret, toAdd)
                }
            }

        }
    }
    return ret
}

func parseTP(stream []string) []*structs.Node {
    var ret []*structs.Node
    for _, state  := range search(stream, "uT_FIN", nil) {
        ret = append(ret, state.Tree)
    }
    return ret
}

func parse(sentence string) *structs.Node {
    stream := utils.Lexify(sentence)
    trees := parseTP(stream)
    if trees == nil {
        return nil
    }
    return trees[0]
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
