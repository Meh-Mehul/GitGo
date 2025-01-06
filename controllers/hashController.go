package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

func DeleteDirectory(hash string) error {
	dirPath := BASE_REL_PATH+OBJECT_REL_PATH+hash[:2] + "/"+hash[2:]+".obj"
	err := os.RemoveAll(dirPath)
	if err != nil {
		return fmt.Errorf("failed to delete directory %s: %w", dirPath, err)
	}
	return nil
}

func GenrateHash(input string) string {
	uniqueInput := input
	hasher := sha256.New()
	hasher.Write([]byte(uniqueInput))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash
}

// Function returns the compressed object of that commit with the given checksum
func getBlobFromHash(hash string ) ([]byte, error) {
	base_path := BASE_REL_PATH+OBJECT_REL_PATH+hash[:2] + "/"
	filename := base_path + hash[2:] + ".obj"
	blobData, err:= os.ReadFile(filename);
	if(err != nil){
		fmt.Println("FETCH_ERR: An error Occured while fetching commit details");
		fmt.Println(err);
		return nil, err;
	}
	return blobData, nil;
}
