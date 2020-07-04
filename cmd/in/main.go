package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/mtharrison/github-pr-comment-resource/internal/resource"
)

func main() {
	input, err := resource.GetInput(os.Stdin, true)
	if err != nil {
		log.Fatal("Could not read input: ", err)
	}

	comment, err := resource.GithubComment(input)
	if err != nil {
		log.Fatal("Could not fetch comment: ", err)
	}

	// Write resource to file

	dirName := os.Args[1]
	outFile := path.Join(dirName, "comment.json")
	contents, err := json.Marshal(comment)
	if err != nil {
		log.Fatal("JSON marshal error: ", err)
	}
	err = ioutil.WriteFile(outFile, contents, 0644)
	if err != nil {
		log.Fatal("Error writing file: ", err)
	}

	// Write output to STDOUT

	output := resource.Output{
		Version:  *input.Version,
		Metadata: comment.Metadata(),
	}
	out, err := json.Marshal(output)
	if err != nil {
		log.Fatal("JSON marshal error: ", err)
	}
	_, err = os.Stdout.Write(out)
	if err != nil {
		log.Fatal("STDOUT write error: ", err)
	}
}
