package main

func main() {
	photo, _ := fileToBase64("example.jpeg")

	makeLesson([]string{photo})
}
