package controllers

import (
	"fmt"
	"os"
)

func GitGoGet(checksum string){
	_, err := os.ReadDir(BASE_REL_PATH);
	if(err != nil){
		fmt.Println("-----GetGo Not Initialized----");
		return;
	}
	// Get the Hash and respective FileData and save it into the file
	Tree := CommitTree{}
	Tree.ReadFromJSON(BASE_REL_PATH+COMMIT_TREE_REL_PATH+COMMIT_TREE_NAME)
	Tree.FixTree()
	That_Node := Tree.SearchCommit(checksum)
	if(That_Node == nil){
		fmt.Println("Did not Find any node with that checksum")
		return
	}
	BlobData, err := getBlobFromHash(checksum);
	if(err != nil){
		fmt.Println("FATAL_ERR: Could not retreive File at the commit")
		return
	}
	ThatFileData, err := DecompressData(BlobData);
	if(err != nil){
		fmt.Println("FATAL_ERR: Could not retreive File at the commit")
		return
	}
	filename, err := ReadGOFile(BASE_REL_PATH+FILENAME_PATH_NAME)
	if(err != nil){
		fmt.Println("FATAL_ERR: Could not retreive File at the commit")
		return
	}
	op_file_name := string(filename) + "_"+checksum[:8]
	err1 := os.WriteFile(op_file_name, ThatFileData, 0644)
	if(err1 != nil){
		fmt.Println("FATAL_ERR: Could not retreive File at the commit")
		return
	}
	fmt.Printf("Commit's Data Written to %s, You can see there.\n", op_file_name);
}