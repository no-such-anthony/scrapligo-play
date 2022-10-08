package configparser

import (
    "bufio"
    "fmt"
    //"io"
    "os"
    //"strings"
	"regexp"
)

// based around the python ciscoconfparse
// https://github.com/mpenning/ciscoconfparse
// and node trees
// https://stackoverflow.com/questions/22957638/make-a-tree-from-a-table-using-golang

type Node struct {
    Line     string
    Children []*Node
	Parents []*Node
	Indent	string
}

type CP map[int]*Node


func (cp CP) FindNodes(r string) []*Node {
	var a []*Node

	re := regexp.MustCompile(r)

	for c := 1; c < len(cp)+1; c++ {
		if re.MatchString(cp[c].Line) {
			a = append(a, cp[c])
		}
	}

	return a
}

// TODO Insert and Delete config?
// TODO More find functions?
// TODO Banner handling?
// TODO Is there a better data structure?
// TODO Different models, vendors, etc?


func Scan(f string) CP {

	nodeTable2 := CP{}

	file, err := os.Open(f)
	if err != nil {
		fmt.Printf("Could not open the file due to this %s error \n", err)
	}
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	c := 1
	re := regexp.MustCompile(`^( +)`)
	var prevIndent string

	for fileScanner.Scan(){
		line := fileScanner.Text()
		indent := re.FindString(line)
		node := &Node{Line: line, Children: []*Node{}, Parents: []*Node{}, Indent: indent}
		
		if c == 1 {
			//fmt.Println(c, line, node.Parents)
			nodeTable2[c] = node
			c += 1
			continue
		}

		prevIndent = nodeTable2[c-1].Indent
		if len(indent) > 0 {
			node.Parents = append(node.Parents, nodeTable2[c-1].Parents...)
		}
		if len(indent) > len(prevIndent) {
			parent := nodeTable2[c-1]
			node.Parents = append(node.Parents, parent)
			parent.Children = append(parent.Children, node)
		} else if len(indent) > 0 && len(indent) == len(prevIndent) {
			parent := node.Parents[len(node.Parents)-1]
			parent.Children = append(parent.Children, node)
		} else if len(indent) > 0 && len(indent) < len(prevIndent)  {
			node.Parents = node.Parents[:len(node.Parents)-1]
			parent := node.Parents[len(node.Parents)-1]
			parent.Children = append(parent.Children, node)
		}
		
		//fmt.Println(c, line, node.Parents)
		nodeTable2[c] = node
		c += 1

	}

	//TODO Banner handling

	return nodeTable2
}

