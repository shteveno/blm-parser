package main

import (
    //"testing"
    //"fmt"
    "blm/utils"
)

func ExampleParseDP() {
    tree := parse("the cake")
    utils.Latex(tree, 0)
    // Output:
    // [.DP
    //   [.D the ]
    //   [.NP 
    //     [.N cake ]
    //   ]
    // ]
}
