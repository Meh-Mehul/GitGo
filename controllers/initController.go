package controllers
// File to Manage creation of .gitGo Folder

// Structure of .gitGo folder (in my implementation)
// .gitGo
// 	|-->commits
//	|---->tree.gotem
//	|-->objects
//	|---->(value-value store)
//	|-->refs
//	|---->Branch Pointers
//	|-->HEAD.gotem (Current Pointer)

import(
	"os"
	"os/user"
	"fmt"
)

func createDirectory(path string){
	err := os.Mkdir(path,0755) // read/write/execute perms
	if(err != nil){
		fmt.Println("Some error Occured");
	}
}

func GitGoInit(path string, filename string){
	// check for previous .gitGo
	_, err := os.ReadDir(BASE_REL_PATH);
	if(err == nil){
		fmt.Println("-----GetGo Already Initialized----");
		return;
	}
	// Create .gitGo and read the file into Filedata
	fileData, err := os.ReadFile(filename)
	if err != nil{
		fmt.Println("FATAL_ERR: An error Occured while initializing GitGo. Maybe file Does not exist?");
		return
	}
	createDirectory(BASE_REL_PATH);
	// Get Hash
	checksum:=GenrateHash("main"+string(fileData));
	// create refs, refs/main, objects/main, 
	createDirectory(BASE_REL_PATH+REFS_REL_PATH);
	createDirectory(BASE_REL_PATH+OBJECT_REL_PATH)
	createDirectory(BASE_REL_PATH+COMMIT_TREE_REL_PATH)
	// assign HEAD_PTR
	Head_path := BASE_REL_PATH+"HEAD.gotem";
	file, err:= os.Create(Head_path);
	if(err != nil){
		fmt.Println("FATAL_ERR: An error Occured while initializing GitGo.");
		return
	}
	defer file.Close()
	_, err1 := file.WriteString(checksum);
	if(err1 != nil){
		fmt.Println("FATAL_ERR: An error Occured while initializing GitGo.");
		return
	}
	// Intialize tree and Node
	author, err2 := user.Current();
	if(err2 != nil){
		fmt.Println("FATAL_ERR: An error Occured while initializing GitGo.");
		return
	}
	Tree := NewCommitTree(checksum, author.Username, "main");
	// Store in TRee and store the tree itself in commit/
	Tree.WriteToJSON(BASE_REL_PATH+COMMIT_TREE_REL_PATH+COMMIT_TREE_NAME);
	// read the file, create Blob, and then Store into object/root_data
	blobData, err:= CompressData(fileData);
	if(err != nil){
		fmt.Println("FATAL_ERR: An error Occured while initializing GitGo.");
		return
	}
	obj_dir := checksum[:2]; // this is how its done in github too!, so we have O(256) different directories
	createDirectory(BASE_REL_PATH+OBJECT_REL_PATH+obj_dir);
	file_name := checksum[2:];
	obj_path := BASE_REL_PATH+OBJECT_REL_PATH+obj_dir+"/"+file_name+".obj";

	err5 := os.WriteFile(obj_path, blobData, 0644);
	if err5 != nil{
		fmt.Println("FATAL_ERR: An error Occured while initializing GitGo.");
		return
	}
	// link the Node to The refs/main
	refPath := BASE_REL_PATH + REFS_REL_PATH + "main.ref"
	err4 := os.WriteFile(refPath, []byte(checksum), 0644)
	if(err4 != nil){
		fmt.Println("FATAL_ERR: An error Occured while initializing GitGo.");
		return
	}
	// add the filename to filename.gotem
	filename_path := BASE_REL_PATH+FILENAME_PATH_NAME
	file1, err6 := os.Create(filename_path)
	if err6 != nil{
		fmt.Println("FATAL_ERR: An error Occured while initializing GitGo.");
		return
	}
	defer file1.Close();
	file1.Write([]byte(filename));
	// return file hash
	fmt.Println("gitGo Tracker Initialized for: "+filename+" with the checksum: "+checksum);
}
