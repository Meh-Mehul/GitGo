package controllers

import (
	"os"
	"os/user"
	"fmt"
)

// Lists all branches
func GitGoBranch_Noarg(){
	_, err := os.ReadDir(BASE_REL_PATH);
	if(err != nil){
		fmt.Println("-----GetGo Not Initialized----");
		return;
	}
	Tree := CommitTree{}
	err1 := Tree.ReadFromJSON(BASE_REL_PATH+COMMIT_TREE_REL_PATH+COMMIT_TREE_NAME);
	if(err1 != nil){
		fmt.Println("ERR: error occured while reading Commit History. Re-Initialize GitGO?")
		return
	}
	Tree.FixTree()
	bs := Tree.ListBranches()
	currB := GetCurrentBranch()
	fmt.Println("Current Branches:")
	for _,val := range bs{
		if(val+".ref" == currB){
			fmt.Printf("*%s\n", val)
		} else{
			fmt.Printf("%s\n", val);
		}
		
	}

}



func GetCurrentBranch() string {
	prev_hash , err := ReadGOFile(BASE_REL_PATH+"HEAD.gotem");
	if(err != nil){
		fmt.Println("FATAL_ERR: Fatal Error occured in getting Branch Info")
		return ""
	}
	if _, err := os.Stat(BASE_REL_PATH+REFS_REL_PATH); os.IsNotExist(err) {
		return ""
	}
	files, err := os.ReadDir(BASE_REL_PATH+REFS_REL_PATH)
	if err != nil {
		return ""
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filePath := BASE_REL_PATH+REFS_REL_PATH + "/" + file.Name()
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("error reading file %s: %v\n", file.Name(), err)
			continue
		}
		if(string(prev_hash) == string(data)){
			return file.Name()
		}
	}
	return ""
}



// This does NOT Preserve Your On-Going work, for that use, git stash
func GitGoBranch(name string){
	_, err := os.ReadDir(BASE_REL_PATH);
	if(err != nil){
		fmt.Println("-----GetGo Not Initialized----");
		return;
	}
	curr_branch_name := GetCurrentBranch();
	if(curr_branch_name == name+".ref"){
		fmt.Println("Already On Branch "+name);
		return
	}
	Tree := CommitTree{}
	err1 := Tree.ReadFromJSON(BASE_REL_PATH+COMMIT_TREE_REL_PATH+COMMIT_TREE_NAME);
	if(err1 != nil){
		fmt.Println("ERR: error occured while reading Commit History. Re-Initialize GitGO?")
		return
	}
	Tree.FixTree()
	ok := Tree.CheckBranchExists(name);
	if(ok){
		// Branch existed!
		div_Node, err := Tree.FindFirstDivergentNode(name);
		if(err != nil){
			fmt.Println("FATAL_ERR: An error Occured while Switching to new branch")
			return
		}
		stack := []*Node{div_Node}
		last_node := Node{ID:-1}
		for len(stack)>0{
			curr := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			last_node = *curr;
			for i := len(curr.Children) - 1; i >= 0; i-- {
				if(curr.Children[i].Branch == name){
					stack = append(stack, curr.Children[i])
				}
			}
		}
		if(last_node.ID == -1){
			fmt.Println("FATAL_ERR: An error Occured while Switching to new branch")
			return
		}
		// swicthing to that node in branch
		curr_hash := last_node.Checksum;
		// Writing to branch.ref and HEAD.gotem
		branch_path := BASE_REL_PATH+REFS_REL_PATH+last_node.Branch+".ref"
		err3 := os.WriteFile(branch_path, []byte(curr_hash), 0644);
		if(err3 != nil){
			fmt.Println("FATAL_ERR: An error Occured while Switching to new branch")
			return
		}
		Head_path := BASE_REL_PATH+"HEAD.gotem"
		err4 := os.WriteFile(Head_path, []byte(curr_hash), 0644);
		if(err4 != nil){
			fmt.Println("FATAL_ERR: An error Occured while Switching to new branch")
			return
		}
		stashBlobData, err := getBlobFromHash(curr_hash);
		if(err != nil){
			fmt.Println("FATAL_ERR: An error Occured while Switching to new branch")
			return
		}
		filename, err := ReadGOFile(BASE_REL_PATH+FILENAME_PATH_NAME)
		if(err != nil){
			fmt.Println("FATAL_ERR: An error Occured while Switching to new branch")
			return
		}
		decomp_stashBlobData, err := DecompressData(stashBlobData);
		if(err != nil){
			fmt.Println("FATAL_ERR: An error Occured while Switching to new branch")
			return
		}
		file, err := os.OpenFile(string(filename), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if(err != nil){
			fmt.Println("FATAL_ERR: An error Occured while Switching to new branch")
			return
		}
		defer file.Close()
		_, err5 := file.Write(decomp_stashBlobData);
		if(err5 != nil){
			fmt.Println("FATAL_ERR: An error Occured while Switching to new branch")
			return
		}
		fmt.Printf("Successfully Switched to Latest Commit at Branch : %s\n", name);
	}
	if(!ok){
		// We need to Create a New Branch at current Node (New Hash!, new commits, new stuff!)

		// Check if there are >1 commits
		if(Tree.lastID == 1){
			fmt.Println("For Changing branch you need to have atleast one commit apart from intial commit")
			return
		}
		// If yes, then go to parent commit of current
		curr_hash, err := ReadGOFile(BASE_REL_PATH+"HEAD.gotem")
		if(err != nil){
			fmt.Println("FATAL_ERR: An error Occured while Switching to new branch")
			return
		}
		curr_Node := Tree.SearchCommit(string(curr_hash));
		if(curr_Node == nil){
			fmt.Println("FATAL_ERR: An error Occured while Switching to new branch")
			return
		}
		pareto := curr_Node.Parent;
		// Now we generate hash for the new Branch by using file data present at that node
		BlobData, err := getBlobFromHash(string(curr_hash));
		if(err != nil){
			fmt.Println("FATAL_ERR: An error Occured while Switching to new branch")
			return
		}
		decomp_BlobData, err := DecompressData(BlobData);
		if(err != nil){
			fmt.Println("FATAL_ERR: An error Occured while Switching to new branch")
			return
		}
		new_checksum := GenrateHash(name+string(decomp_BlobData));
		author, err2 := user.Current();
		if(err2 != nil){
			fmt.Println("FATAL_ERR: An error Occured while Switching to new branch")
			return
		}
		// Then make a new children node with new branch (effectively a new commit)
		NewNode, err := Tree.AddCommit(pareto, new_checksum, author.Username,"Switched To Branch "+name)
		NewNode.Branch = name;
		if(err != nil){
			fmt.Println("FATAL_ERR: An error Occured while Switching to new branch")
			return
		}
		// Dont Forget to add to branch.ref, HEAD.gotem, etc
		Head_path := BASE_REL_PATH+"HEAD.gotem"
		err4 := os.WriteFile(Head_path, []byte(new_checksum), 0644);
		if(err4 != nil){
			fmt.Println("FATAL_ERR: An error Occured while Switching to new branch")
			return
		}
		branch_path := BASE_REL_PATH+REFS_REL_PATH+NewNode.Branch+".ref"
		err3 := os.WriteFile(branch_path, []byte(new_checksum), 0644);
		if(err3 != nil){
			fmt.Println("FATAL_ERR: An error Occured while Switching to new branch")
			return
		}
		// Write this to a new Object with the new Hash
		obj_dir := new_checksum[:2];
		createDirectory(BASE_REL_PATH+OBJECT_REL_PATH+obj_dir);
		file_name := new_checksum[2:];
		obj_path := BASE_REL_PATH+OBJECT_REL_PATH+obj_dir+"/"+file_name+".obj";
		err1 := os.WriteFile(obj_path, BlobData, 0644);
		if err1 != nil{
			fmt.Println("FATAL_ERR: An error Occured while Switching to new branch")
			return
		}
		Tree.WriteToJSON(BASE_REL_PATH+COMMIT_TREE_REL_PATH+COMMIT_TREE_NAME)
		fmt.Printf("Successfully Switched to Branch %s\n", name);
	}
}