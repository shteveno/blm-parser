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
func spec(stream []string, uFeat string) []*structs.Node {
    if stream == nil || uFeat == "." {
        return nil
    }
    uCat, _, _ := defeat(uFeat)
    fmt.Println(stream, uFeat)
    found := search(stream, "u" + uCat)
    //for _, tree := range found {
    //    utils.Latex(tree, 0)
    //}
    return found
}

func head(stream []string, uFeat string) []*structs.State {
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
            x := &structs.Node{}
            x.Label = uCat + "_" + uSel
            x.Form = "$\\sout{" + uCat + "P}$"
            x.Features = feats
            states = append(states, &structs.State{x, stream, wordsUsed, feats[1], feats[2]})
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
            // There is an overt head present
            } else {
                bundles = feat[uCat]
                wordsUsed = i + 1
            }
            for _, bundle := range bundles {
                feats := strings.Split(bundle, ",")
                if feats[1] != "." && wordsUsed == 0 {
                    continue
                }
                if feats[0] == sel {
                    x := &structs.Node{}
                    x.Label = uCat + "_" + sel
                    if feats[2][0] == '*' {
                        fmt.Println("Strong feature found: " + feats[2])
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

func comp(stream []string, uFeat string) []*structs.State {
    return nil
    //return []*structs.State{&structs.State{&structs.Node{"vP", "", &structs.Node{"v_{\\textsc{ag}}", "like", nil, nil, nil}, &structs.Node{"VP", "", &structs.Node{"DP", "", nil, nil, nil}, &structs.Node{"V", "", nil, nil, nil}, nil}, nil}, nil}}
}

func search(stream []string, uFeat string) []*structs.Node {
    if stream == nil || uFeat == "." {
        return nil
    }
    var ret []*structs.Node
    heads := head(stream, uFeat)
    if heads == nil {
        return nil
    }
    for _, x := range heads {
        specifiers := spec(stream[:x.HeadPos], x.Spec)
        if specifiers == nil {
            if x.Spec != "." {
                continue
            }
            //fmt.Println("I made a head, no specifier!")
            //utils.Latex(x.Tree, 0)
            complements := comp(x.Remaining, x.Comp)
            if complements == nil {
                if x.Comp != "." {
                    continue
                }
                ret = append(ret, utils.FormTree(nil, x.Tree, nil))
                continue
            }
            for _, yP := range complements {
                if yP.Remaining == nil {
                    ret = append(ret, utils.FormTree(nil, x.Tree, yP.Tree))
                }
            }
        }
        for _, wP := range specifiers {
            fmt.Println("I made a specifier!")
            utils.Latex(wP, 0)
            complements := comp(x.Remaining, x.Comp)
            if complements == nil {
                if x.Comp != "." {
                    continue
                }
                ret = append(ret, utils.FormTree(wP, x.Tree, nil))
                continue
            }
            for _, yP := range complements {
                if yP.Remaining == nil {
                    ret = append(ret, utils.FormTree(wP, x.Tree, yP.Tree))
                }
            }

        }
    }
    return ret
}

func parseTP(stream []string) []*structs.Node {
    return search(stream, "uT_FIN")
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
