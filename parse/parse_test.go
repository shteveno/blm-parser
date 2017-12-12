package main

import (
    //"testing"
    //"fmt"
    "blm/utils"
)

/*
func ExampleParseDP() {
    found := spec("the cake")
    if len(spec) == 1 {
        utils.Latex(found[0], 0)
    }
    // Output:
    // [.DP
    //   [.D the ]
    //   [.NP 
    //     [.N cake ]
    //   ]
    // ]
}
*/

func ExampleHead() {
    lex = make(map[string][]string)
    feat = make(map[string][]string)
    null = make(map[string][]string)
    utils.Init_map(lex, "../lexicon.txt")
    utils.Init_map(feat, "../features.txt")
    utils.Init_map(null, "../null.txt")
    stream := utils.Lexify("I like cake")
    heads := head(stream, "uT_FIN")
    for _, head := range heads {
        utils.Latex(head.Tree, 0)
    }
    // Output:
    // [.T_FIN $\varnothing$ ]
    // [.T_FIN $\varnothing$ ]
    // [.T_FIN $\varnothing$ ]
}
/*
func ExampleParseTPNoDP() {
    tree := parse("I like cake")
    utils.Latex(tree, 0)
    // Output:
    // [.TP
    //   [.DP 
    //     [.D $\varnothing$ ]
    //     [.NP 
    //       [.N cake ]
    //     ]
    //   ]
    //   [.NP 
    //     [.N cake ]
    //   ]
    // ]
}

func ExampleParseTP() {
    tree := parse("I like the cake")
    utils.Latex(tree, 0)
    // Output:
    // [.TP
    //   [.DP 
    //     [.D the ]
    //     [.NP 
    //       [.N cake ]
    //     ]
    //   ]
    //   [.NP 
    //     [.N cake ]
    //   ]
    // ]
}
*/
