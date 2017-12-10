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
        return string(sentence)
    }
}

// Formats the resulting tree such that it may be compiled in
// Latex.
func Latex(root *structs.Node, depth int) {
    if root == nil {
        return
    }
    offset := strings.Repeat(" ", depth*2)
    if root.Left == root.Right {
        fmt.Println(offset + "[." + root.Label + " " + root.Form + " ]")
        return
    }
    fmt.Println(offset + "[." + root.Label + " ")
    Latex(root.Left, depth + 1)
    Latex(root.Right, depth + 1)
    fmt.Println(offset + "]")
}
