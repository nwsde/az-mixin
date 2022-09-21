package az

import (
	"context"
	"strings"

	"get.porter.sh/porter/pkg/exec/builder"
)

var (
	_ TypedCommand             = &GroupCommand{}
	_ builder.HasErrorHandling = &GroupCommand{}
)

// GroupCommand ensures that a group exist or not
type GroupCommand struct {
	action      string
	Description string `yaml:"description"`
	Name        string `yaml:"name"`
	Location    string `yaml:"location"`
}

func (c *GroupCommand) HandleError(ctx context.Context, err builder.ExitError, stdout string, stderr string) error {
	switch c.action {
	case "uninstall":
		// It's okay if we try to delete the resource group and it's already gone
		if strings.Contains(stderr, "ResourceGroupNotFound") {
			return nil
		}
	default:
		// NOTE: We are treating "already exists" as an okay state and reusing the existing group,
		// but the user could have requested a different location.
		// We may want a feature that recreates the group to get it into the correct location.
		if strings.Contains(stderr, "Resource group already exists") {
			return nil
		}
	}

	return err
}

func (c *GroupCommand) GetWorkingDir() string {
	return ""
}

func (c *GroupCommand) SetAction(action string) {
	c.action = action
}

func (c *GroupCommand) GetCommand() string {
	return "az"
}

func (c *GroupCommand) GetArguments() []string {
	switch c.action {
	case "uninstall":
		return []string{"group", "delete", "--yes"}
	default:
		return []string{"group", "create"}
	}
}

func (c *GroupCommand) GetFlags() builder.Flags {
	var flags builder.Flags
	flags = append(flags, builder.NewFlag("name", c.Name))

	if c.action != "uninstall" {
		flags = append(flags, builder.NewFlag("location", c.Location))
	}

	return flags
}

func (c *GroupCommand) SuppressesOutput() bool {
	return false
}
