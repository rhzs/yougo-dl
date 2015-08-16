package command

import (
	"strings"
)

type GetCommand struct {
	Meta
}

func (c *GetCommand) Run(args []string) int {
	video, err := Get(args[0])
	if err != nil {
		return -1
	}

	video.Download(0, args[1])

	return 0
}

func (c *GetCommand) Synopsis() string {
	return ""
}

func (c *GetCommand) Help() string {
	helpText := `
Get video
`
	return strings.TrimSpace(helpText)
}
