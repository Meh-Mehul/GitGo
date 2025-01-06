package controllers

import (
	"fmt"
	"os"
)
// Removes the current Commit and its successors from the Commit-Tree
func GitGoRevert_Noarg(){
	_, err := os.ReadDir(BASE_REL_PATH);
	if(err != nil){
		fmt.Println("-----GetGo Not Initialized----");
		return;
	}
    filename, err := ReadGOFile(BASE_REL_PATH+FILENAME_PATH_NAME)
	if(err != nil){
		fmt.Println("FATAL_ERR: Fatal Error occured in Reverting")
		return
	}
	RN_hash , err := ReadGOFile(BASE_REL_PATH+"HEAD.gotem");
	if(err != nil){
		fmt.Println("FATAL_ERR: Fatal Error occured in Reverting")
		return
	}
    Tree := CommitTree{}
    Tree.ReadFromJSON(BASE_REL_PATH+COMMIT_TREE_REL_PATH+COMMIT_TREE_NAME)
    Tree.FixTree()
    if(string(RN_hash) == Tree.Root.Checksum){
        fmt.Println("USER_ERR: No Commit found, except intial commit");
        return
    }
    curr_Node := Tree.SearchCommit(string(RN_hash));
    pareto := curr_Node.Parent;
    tobeBlobData, err := getBlobFromHash(pareto.Checksum);
    if(err != nil){
        fmt.Println("FATAL_ERR: Fatal Error occured in Reverting")
		return
    }
    child_array := curr_Node.Parent.Children;
    currid := -1;
    for id, val  := range child_array {
        if(val.Checksum == curr_Node.Checksum){
            currid = id;
            break;
        }
    }
    if(currid == -1){
        fmt.Println("FATAL_ERR: Fatal Error occured in Reverting")
        return
    }
    fmt.Println("Deleting Current Commits data....")
    // We remove the current Node as a children of previous, thereby erasing its existence
    child_array = append(child_array[:currid],child_array[currid+1:]...);
    curr_Node.Parent.Children = child_array;
    Tree.WriteToJSON(BASE_REL_PATH+COMMIT_TREE_REL_PATH+COMMIT_TREE_NAME);
    // Now we delete the blobs too! (using dfs)
	stack := []*Node{curr_Node}
	for len(stack) > 0 {
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
        DeleteDirectory(curr.Checksum);
        for i := len(curr.Children) - 1; i >= 0; i-- {
			stack = append(stack, curr.Children[i])
		}
	}
    fmt.Println("Writing reverted data into the File....")
    // now we write that (pareto's) hash to current file
    tobeFileData, err := DecompressData(tobeBlobData);
    if(err != nil){
        fmt.Println("FATAL_ERR: Fatal Error occured in Reverting")
		return
    }
	file, err := os.OpenFile(string(filename), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
    if(err != nil){
        fmt.Println("FATAL_ERR: Fatal Error occured in Reverting")
		return
    }
	defer file.Close()
	_, err1 := file.Write(tobeFileData);
    if(err1 != nil){
        fmt.Println("FATAL_ERR: Fatal Error occured in Reverting")
		return
    }
	branch_path := BASE_REL_PATH+REFS_REL_PATH+pareto.Branch+".ref"
	err3 := os.WriteFile(branch_path, []byte(pareto.Checksum), 0644);
	if(err3 != nil){
		fmt.Println("COMMIT_ERR: Fatal Error occured in Reverting")
		return
	}
	Head_path := BASE_REL_PATH+"HEAD.gotem"
	err4 := os.WriteFile(Head_path, []byte(pareto.Checksum), 0644);
	if(err4 != nil){
		fmt.Println("COMMIT_ERR: Fatal Error occured in Reverting")
		return
	}
    fmt.Printf("Successfully Reverted to the Commit: %s\n Author: %s\n Message: %s\n Date: %v\n", pareto.Checksum, pareto.Author, pareto.CommitMessage, pareto.Time)

}

// Removes everything after a commit and reverts to that point in time (Even ALL branches)
func GitGoRevert(checksum string){
	_, err := os.ReadDir(BASE_REL_PATH);
	if(err != nil){
		fmt.Println("-----GetGo Not Initialized----");
		return;
	}
    filename, err := ReadGOFile(BASE_REL_PATH+FILENAME_PATH_NAME)
	if(err != nil){
		fmt.Println("FATAL_ERR: Fatal Error occured in Reverting")
		return
	}
	currHash , err := ReadGOFile(BASE_REL_PATH+"HEAD.gotem");
	if(err != nil){
		fmt.Println("FATAL_ERR: Fatal Error occured in Reverting")
		return
	}
    RN_hash := checksum;
    if(RN_hash == string(currHash)){
        fmt.Println("Already at that Commit.")
        return
    }
    Tree := CommitTree{}
    Tree.ReadFromJSON(BASE_REL_PATH+COMMIT_TREE_REL_PATH+COMMIT_TREE_NAME)
    Tree.FixTree()
    // if(string(RN_hash) == Tree.Root.Checksum){
    //     fmt.Println("USER_ERR: No Commit found, except intial commit");
    //     return
    // }
    curr_Node := Tree.SearchCommit(string(RN_hash));
    if(curr_Node == nil){
        fmt.Println("No Such commit Found :(")
        return
    }
    pareto := curr_Node; // since i copied above code, its better to keep it that way.
    tobeBlobData, err := getBlobFromHash(pareto.Checksum);
    if(err != nil){
        fmt.Println("FATAL_ERR: Fatal Error occured in Reverting")
		return
    }
    
    fmt.Println("Deleting Future Commits data....")
   
    // Now we delete the blobs  (using dfs)
	stack := curr_Node.Children
	for len(stack) > 0 {
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
        DeleteDirectory(curr.Checksum);
        for i := len(curr.Children) - 1; i >= 0; i-- {
			stack = append(stack, curr.Children[i])
		}
	}
    curr_Node.Children = []*Node{}
    Tree.WriteToJSON(BASE_REL_PATH+COMMIT_TREE_REL_PATH+COMMIT_TREE_NAME);
    fmt.Println("Writing reverted data into the File....")
    // now we write that (pareto's) hash to current file
    tobeFileData, err := DecompressData(tobeBlobData);
    if(err != nil){
        fmt.Println("FATAL_ERR: Fatal Error occured in Reverting")
		return
    }
	file, err := os.OpenFile(string(filename), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
    if(err != nil){
        fmt.Println("FATAL_ERR: Fatal Error occured in Reverting")
		return
    }
	defer file.Close()
	_, err1 := file.Write(tobeFileData);
    if(err1 != nil){
        fmt.Println("FATAL_ERR: Fatal Error occured in Reverting")
		return
    }
	branch_path := BASE_REL_PATH+REFS_REL_PATH+pareto.Branch+".ref"
	err3 := os.WriteFile(branch_path, []byte(pareto.Checksum), 0644);
	if(err3 != nil){
		fmt.Println("FATAL_ERR: Fatal Error occured in Reverting")
		return
	}
	Head_path := BASE_REL_PATH+"HEAD.gotem"
	err4 := os.WriteFile(Head_path, []byte(pareto.Checksum), 0644);
	if(err4 != nil){
		fmt.Println("FATAL_ERR: Fatal Error occured in Reverting")
		return
	}
    fmt.Printf("Successfully Reverted to the Commit: %s\n Author: %s\n Message: %s\n Date: %v\n", pareto.Checksum, pareto.Author, pareto.CommitMessage, pareto.Time)
}