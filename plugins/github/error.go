package github

import (
	"fmt"
)

type repoError struct {
	text string
}

func (err *repoError) Error() string {
	return fmt.Sprintf("could not extract owner/repo from %s", err.text)
}

type issueError struct {
	text string
}

func (err *issueError) Error() string {
	return fmt.Sprintf("could not extract issue arguments from %s", err.text)
}
