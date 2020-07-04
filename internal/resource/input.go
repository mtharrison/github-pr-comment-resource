package resource

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

// Input is provided to this resource via STDIN
type Input struct {
	Source  *Source  `json:"source"`
	Version *Version `json:"version"`
}

// Validate checks if the input is sound
func (i *Input) Validate(requireVersion bool) error {

	if i.Source == nil {
		return errors.New("source cannot be empty")
	}

	if i.Source.RepositoryString == "" {
		return errors.New("source.repository cannot be empty")
	}

	if i.Source.AccessToken == "" {
		return errors.New("source.access_token cannot be empty")
	}

	if i.Version == nil && requireVersion {
		return errors.New("version cannot be empty")
	}

	if i.Source.Regex != "" {
		_, err := regexp.Compile(i.Source.Regex)
		if err != nil {
			return errors.New("source.regex not a valid regex")
		}
	}

	if i.Version != nil {
		if i.Version.CommentID == "" || i.Version.PrNumber == "" {
			return errors.New("if set, version must include comment and pr")
		}

		_, err := strconv.ParseInt(i.Version.CommentID, 10, 64)
		if err != nil {
			return errors.New("version.comment not parsable as i64")
		}

		_, err = strconv.Atoi(i.Version.PrNumber)
		if err != nil {
			return errors.New("version.pr not parsable as int")
		}
	}

	return nil
}

// Source is the configurable settings a user provides to this resource
type Source struct {
	RepositoryString string `json:"repository"`
	AccessToken      string `json:"access_token"`
	V3Endpoint       string `json:"v3_endpoint"`
	Regex            string `json:"regex"`
}

// Owner returns the repo owner
func (s *Source) Owner() string {
	return strings.Split(s.RepositoryString, "/")[0]
}

// Repo returns the repo name
func (s *Source) Repo() string {
	return strings.Split(s.RepositoryString, "/")[1]
}

// Version represents a single version of the resource
type Version struct {
	PrNumber  string `json:"pr"`
	CommentID string `json:"comment"`
}

// GetInput takes a reader and constructs an Input
func GetInput(reader io.Reader, requireVersion bool) (Input, error) {
	input, err := readInput(reader)
	if err != nil {
		return input, err
	}

	err = input.Validate(requireVersion)
	if err != nil {
		return input, err
	}

	return input, nil
}

func readInput(reader io.Reader) (Input, error) {
	var input Input

	rawInput, err := ioutil.ReadAll(reader)
	if err != nil {
		return input, err
	}

	err = json.Unmarshal(rawInput, &input)
	if err != nil {
		return input, err
	}

	return input, nil
}
