package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type snippet struct {
	// Content type can be either 'text', for simple markdown text, or 'animation' for a manim animation
	Content_type string `json:"content-type"`

	// Represents either the raw markdown text or a prompt which will be used by another AI to create the manim animation
	// The prompt must be detailed, concise and clearly state that a manim animation is the expected output and nopthing else
	// The animation shouldn't be longer than a couple of seconds
	Content string `json:"content"`
}

func makeAnimation(prompt string, index int) {
	log.Println("Making animation...")
	log.Println("This will definetly take a long time!")

	// Make work folder
	workingDir := fmt.Sprintf("jobs/worker_animation%d", index)
	os.Mkdir(workingDir, 0755)

	// Ask AI
	cmd := exec.Command("gemini", "-p", prompt, "-m", "gemini-3-pro-preview")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}

	// Extract code
	code := strings.Split(stdout.String(), "```")[1]
	code = "#" + code
	os.WriteFile(workingDir+"/code.py", []byte(code), 0755)

	// Run code
	log.Println("Rendering animation...")
	out, err := renderManim(workingDir+"/code.py", fmt.Sprintf("animation%d", index))
	if err != nil {
		log.Println(err)
	}

	// Copy output to the lesson folder
	wd, _ := os.Getwd()
	err = copyToDir(out, filepath.Join(wd, "lesson"))
	if err != nil {
		log.Println(err)
	}
}

func makeLesson(photos64 []string) {
	// Create directory to store photos in
	os.Mkdir("TEMP_PHOTOS", 0755)
	defer os.RemoveAll("TEMP_PHOTOS")

	log.Println("Downloading pictures...")
	// Store photos
	photoListString := ""
	for i, p := range photos64 {
		here := fmt.Sprintf("TEMP_PHOTOS/lesson_%d.jpg", i)
		saveBase64(p, here)
		photoListString += "@" + here + ", "
	}

	// Build prompt
	log.Println("Building prompt...")
	p, err := os.ReadFile("summarize-prompt.txt")
	if err != nil {
		log.Panic(err)
	}
	prompt := strings.ReplaceAll(string(p), "%PHOTOS_LOCATION%", photoListString)
	cmd := exec.Command("gemini", "-p", prompt, "-m", "gemini-3-flash-preview")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run AI analysis
	log.Println("Running AI...")
	log.Println("This could take a while...")
	err = cmd.Run()
	if err != nil {
		log.Println(stderr.String())
		log.Panic(err)
	}

	// Extract JSON
	raw_json := strings.Split(stdout.String(), "```")[1]
	raw_json = strings.TrimPrefix(raw_json, "json")
	var DATA []snippet
	json.Unmarshal([]byte(raw_json), &DATA)

	// Make new folder
	os.RemoveAll("lesson")
	os.Mkdir("lesson", 0755)

	// Create markdown file
	markdown := ""
	for i, o := range DATA {
		if o.Content_type == "text" {
			markdown += o.Content + "\n\n"
		} else if o.Content_type == "animation" {
			makeAnimation(o.Content, i)

			// Include animation
			markdown += fmt.Sprintf("<video width=\"640\" height=\"360\" controls><source src=\"animation%d.mp4\" type=\"video/mp4\">Your browser does not support the video tag.</video>\n\n", i)
		}
	}

	// Save file
	err = os.WriteFile("lesson/main.md", []byte(markdown), 0755)
	if err != nil {
		log.Panic(err)
	}

	// Debug print
	fmt.Print(stdout.String())
}
