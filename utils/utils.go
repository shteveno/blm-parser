// A file that contains functions for reading text files
// and basic string manipulation
package utils

import(
    "blm/structs"
    "fmt"
    "bufio"
    "log"
    "os"
    "strings"
)

// Takes as input a .txt "lexicon" file that lists of all
// of the lexical heads known to the program. Converts 
// the .txt file into a map with the phonological form
// as the key (a string) and the possible lexical categories as
// the value (a list of strings). A word may have multiple 
// lexical categories and will be stored as a list.

// TODO: Add functionality that word can have multiple
// lexical categories!

func Init_map(mapping map[string][]string, fileName string) {
	file, err := os.Open(fileName)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        //fmt.Println(scanner.Text())
        word_cat := strings.Split(scanner.Text(), ": ")
        switch len(word_cat) {
        case 0:
            fmt.Println("Parsing error with ", scanner.Text())
            continue
        case 1:
            mapping[word_cat[0]] = []string{}
        default:
            word := word_cat[0]
            cat := word_cat[1]
            if mapping[word] != nil {
                mapping[word] = append(mapping[word], cat)
            } else {
                mapping[word] = []string{cat}
            }
        }
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}

func PrintMap(maps []map[string][]string) {
    for _, m := range maps {
        for k, v := range m {
            fmt.Println("k is ", k, " and v is ", v, " with length ", len(v))
        }
        fmt.Println("")
    }
}

// Waits for user input to process and parse into a tree.
func Wait_user() string {
	buf := bufio.NewReader(os.Stdin)
    fmt.Print("> ")
    sentence, err := buf.ReadBytes('\n')
    if err != nil {
        fmt.Println(err)
		return ""
    } else {
        return strings.TrimSpace(string(sentence))
    }
}

/* Turns the sentence into a stream of lexical heads. */
func Lexify(sentence string) []string {
    return strings.Split(sentence, " ")
}

/* Makes a copy of a tree. */
func Copy(root *structs.Node) *structs.Node {
    if root == nil {
        return nil
    }
    ret := &structs.Node{}
    ret.Label = root.Label
    ret.Form = root.Form
    ret.Left = Copy(root.Left)
    ret.Right = Copy(root.Right)
    ret.Features = root.Features
    return ret
}

/* Forms a tree given the specifier, head, and complement. */
func FormTree(spec *structs.Node, head *structs.Node, comp *structs.Node) *structs.Node {
    var (
        xP *structs.Node = &structs.Node{}
        xBar *structs.Node = xP
    )
    cat := strings.Split(head.Label, "_")[0]
    xP.Label = cat + "P"
    if spec != nil {
        xBar = &structs.Node{}
        xP.Left = Copy(spec)
        xP.Right = xBar
        xBar.Label = cat + "'"
    }
    xBar.Left = Copy(head)
    xBar.Right = Copy(comp)
    return xP
}

// Formats the resulting tree such that it may be compiled in
// Latex.
func Latex(root *structs.Node, depth int) {
    if root == nil {
        return
    }
    offset := strings.Repeat(" ", depth*2)
    if root.Left == root.Right {
        cat_sel := strings.Split(root.Form, "_")
        if len(cat_sel) == 2 {
            root.Form = cat_sel[0] + "_{\\textsc{" + cat_sel[1] + "}}"
        }
        fmt.Printf("%s[.%s %s ]\n", offset, root.Label, root.Form)
        return
    }
    fmt.Printf("%s[.%s \n", offset, root.Label)
    Latex(root.Left, depth + 1)
    Latex(root.Right, depth + 1)
    fmt.Printf("%s]\n", offset)
}
