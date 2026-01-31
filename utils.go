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
	// Extract scene name
	sceneName, _ := getSceneName(pythonFile)

	// Define where you want the media to live
	home, _ := os.UserHomeDir()
	cwd, _ := filepath.Abs(home)
	mediaDir := filepath.Join(cwd, "output_media")

	// Prepare the command
	// -p: preview (optional)
	// -o: specific filename
	// --media_dir: where to save the folders
	cmd := exec.Command("python3", "-m", "manim", pythonFile, sceneName, "-o", outputName, "--media_dir", mediaDir)

	// Capture output for debugging
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("render failed: %s\n%s", err, string(out))
	}

	// Manim follows a specific folder structure:
	// {media_dir}/videos/{python_file_name}/{quality}/{outputName}.mp4
	outputPath := filepath.Join(mediaDir, "videos",
		filepath.Base(pythonFile[:len(pythonFile)-len(filepath.Ext(pythonFile))]),
		"1080p60", // Default quality folder
		outputName+".mp4")

	return outputPath, nil
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
