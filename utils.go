package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

func saveBase64(photo string, path string) {
	// Decode the string into bytes
	unbased, err := base64.StdEncoding.DecodeString(photo)
	if err != nil {
		fmt.Println("Cannot decode string:", err)
		return
	}

	// Create the output file
	err = os.WriteFile(path, unbased, 0644)
	if err != nil {
		fmt.Println("Cannot write to file:", err)
		return
	}
}

func fileToBase64(filePath string) (string, error) {
	// 1. Read the entire file into a byte slice
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// 2. Encode the bytes to a Base64 string
	base64Str := base64.StdEncoding.EncodeToString(bytes)

	return base64Str, nil
}

func getSceneName(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Regex to find: class [Name](Scene): or class [Name](MovingCameraScene): etc.
	re := regexp.MustCompile(`class\s+(\w+)\s*\((?:.*)Scene\)`)
	match := re.FindStringSubmatch(string(content))

	if len(match) > 1 {
		return match[1], nil // Returns the first captured group (the class name)
	}

	return "", fmt.Errorf("no Scene class found in file")
}

func renderManim(pythonFile string, outputName string) (string, error) {
	sceneName, _ := getSceneName(pythonFile)
	home, _ := os.UserHomeDir()
	mediaDir := filepath.Join(home, "output_media")

	// Run Manim
	cmd := exec.Command("python3", "-m", "manim", pythonFile, sceneName, "-o", outputName, "--media_dir", mediaDir)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("render failed: %s\n%s", err, string(out))
	}

	// Get script name without extension
	scriptBase := filepath.Base(pythonFile)
	scriptName := scriptBase[:len(scriptBase)-len(filepath.Ext(scriptBase))]

	// Search for the file instead of hardcoding '1080p60'
	// Pattern: {mediaDir}/videos/{scriptName}/*/{outputName}.mp4
	searchPattern := filepath.Join(mediaDir, "videos", scriptName, "*", outputName+".mp4")
	matches, err := filepath.Glob(searchPattern)

	if err != nil || len(matches) == 0 {
		return "", fmt.Errorf("could not find rendered file at %s", searchPattern)
	}

	// Return the first match found
	return matches[0], nil
}

func copyToDir(srcPath, dstDir string) error {
	// 1. Open the source file
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	// 2. Extract the filename and create the destination path
	fileName := filepath.Base(srcPath)
	dstPath := filepath.Join(dstDir, fileName)

	// 3. Create the destination file
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// 4. Copy the contents
	_, err = io.Copy(dst, src)
	return err
}
