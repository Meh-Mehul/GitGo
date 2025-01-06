package controllers

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)
type Node struct {
	ID       int
	Time     time.Time
	Checksum string
	Author   string
	Branch   string
	Parent   *Node `json:"-"`
	Children []*Node
	CommitMessage string
}
type CommitTree struct {
	Root   *Node
	lastID int
}
func NewCommitTree(rootChecksum string, author string, branch string) *CommitTree {
	rootNode := &Node{
		ID:       1,
		Time:     time.Now(),
		Checksum: rootChecksum,
		Author:   author,
		Branch:   branch,
		Parent:   nil,
		Children: []*Node{},
		CommitMessage: "Intial Commit",
	}
	return &CommitTree{
		Root:   rootNode,
		lastID: 1,
	}
}
func (ct *CommitTree) CreateNewNode(checksum string, author string, branch string, commitmsg string) *Node {
	ct.lastID++
	return &Node{
		ID:       ct.lastID,
		Time:     time.Now(),
		Checksum: checksum,
		Author:   author,
		Branch:   branch,
		Parent:   nil,
		Children: []*Node{},
		CommitMessage: commitmsg,
	}
}
func (ct *CommitTree) AddCommit(parentNode *Node, checksum string, author string,commitmsg string) (*Node, error) {
	if parentNode == nil {
		return nil, errors.New("parent node cannot be nil")
	}
	newNode := ct.CreateNewNode(checksum, author, parentNode.Branch, commitmsg)
	newNode.Parent = parentNode
	parentNode.Children = append(parentNode.Children, newNode)
	return newNode, nil
}

func (ct *CommitTree) ChangeBranch(currCommit *Node, newBranch string, checksum string, author string) (*Node, error) {
	if currCommit == nil {
		return nil, errors.New("current commit node cannot be nil")
	}
	newBranchNode := ct.CreateNewNode(checksum, author, newBranch, currCommit.CommitMessage)
	newBranchNode.Parent = currCommit
	currCommit.Children = append(currCommit.Children, newBranchNode)
	return newBranchNode, nil
}
// i used dfs to search in this DAG
// TODO: Return an error to easen my debugging
func (ct *CommitTree) SearchCommit(checksum string) *Node {
	stack := []*Node{ct.Root}
	for len(stack) > 0 {
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if curr.Checksum == checksum {
			return curr
		}
		for i := len(curr.Children) - 1; i >= 0; i-- {
			stack = append(stack, curr.Children[i])
		}
	}
	return nil
}

func (ct *CommitTree) getMaxID(node *Node) int {
	if node == nil {
		return 0
	}
	maxID := node.ID
	for _, child := range node.Children {
		childMax := ct.getMaxID(child)
		if childMax > maxID {
			maxID = childMax
		}
	}
	return maxID
}
func (ct *CommitTree) WriteToJSON(filePath string) error {
	data, err := json.MarshalIndent(ct, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

func (ct *CommitTree) ReadFromJSON(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, ct)
	if err != nil {
		return err
	}
	ct.lastID = ct.getMaxID(ct.Root)
	return nil
}
// Since i've not written tree in Marshall(due to cyclic dep. reasons), i will have to separatel fix connections each time i read the tree from json
// i used simple DFS but to fix connections back-side too
func (ct *CommitTree) FixTree() {
	stack := []*Node{ct.Root}
	for len(stack) > 0 {
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for i := len(curr.Children) - 1; i >= 0; i-- {
			curr.Children[i].Parent = curr;
			stack = append(stack, curr.Children[i])
		}
	}
}

func (ct *CommitTree) ListBranches() []string{
	mp := map[string]int{};
	stack := []*Node{ct.Root}
	for len(stack) > 0 {
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		mp[curr.Branch]+=1;
		for i := len(curr.Children) - 1; i >= 0; i-- {
			curr.Children[i].Parent = curr;
			stack = append(stack, curr.Children[i])
		}
	}
	op := []string{}
	for val  := range mp {
		op = append(op, val);
	}
	return op;
}

func (ct *CommitTree) CheckBranchExists(name string) bool {
	branches := ct.ListBranches();
	for _,val  := range branches {
		if(val == name){
			return true;
		}
	}
	return false;
}

func (ct *CommitTree) FindFirstDivergentNode(name string) (*Node, error){
	if(!ct.CheckBranchExists(name)){
		return nil, errors.New("invalid smth")
	}
	stack := []*Node{ct.Root};
	for len(stack)>0 {
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for i := len(curr.Children) - 1; i >= 0; i-- {
			if(curr.Children[i].Branch == name){
				return curr, nil;
			}
			stack = append(stack, curr.Children[i])
		}
	}
	return nil, errors.New("did not find")

}