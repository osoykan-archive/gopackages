package filerenamer

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func RenameFiles(folderPath string, textToRemove string) error {
	files, err := getFiles(folderPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		_ ,err := renameFileName(file, textToRemove)
		if err != nil {
			log.Printf("Rename got an error : %s", err.Error())
		}
	}

	return nil
}

func getFiles(folderPath string) ([]string, error) {
	empty := "."
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return []string{}, nil
	}

	if filepath.Dir(folderPath) == empty {
		return []string{}, nil
	}

	var files []string
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		return []string{}, err
	}

	return files, nil
}

func renameFileName(filePath string, textToRemove string) (string, error) {
	contains := strings.Contains(filePath, textToRemove)
	if contains == true {
		log.Printf("File is beign renamed: %s", filePath)
		replace := strings.Replace(filePath, textToRemove, "", 1)
		extension:=filepath.Ext(replace)
		withoutExtension:= strings.TrimSpace(strings.Replace(replace, extension, "", 1))
		newPath := withoutExtension + extension
		err := os.Rename(filePath, newPath)
		return newPath, err
	}
	return filePath, nil
}
