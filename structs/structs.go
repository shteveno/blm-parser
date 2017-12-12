package structs


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
    Tree          *Node
    Remaining     []string
    HeadPos       int
    Spec          string
    Comp          string
}
