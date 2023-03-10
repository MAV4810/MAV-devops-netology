package renderers

import (
	"fmt"

	"github.com/hashicorp/terraform/internal/command/format"
	"github.com/hashicorp/terraform/internal/command/jsonformat/computed"
	"github.com/hashicorp/terraform/internal/plans"
)

func SensitiveBlock(diff computed.Diff, beforeSensitive, afterSensitive bool) computed.DiffRenderer {
	return &sensitiveBlockRenderer{
		inner:           diff,
		beforeSensitive: beforeSensitive,
		afterSensitive:  afterSensitive,
	}
}

type sensitiveBlockRenderer struct {
	inner computed.Diff

	afterSensitive  bool
	beforeSensitive bool
}

func (renderer sensitiveBlockRenderer) RenderHuman(diff computed.Diff, indent int, opts computed.RenderHumanOpts) string {
	cachedLinePrefix := fmt.Sprintf("%s%s ", formatIndent(indent), format.DiffActionSymbol(plans.NoOp))
	return fmt.Sprintf("{%s\n%s  # At least one attribute in this block is (or was) sensitive,\n%s  # so its contents will not be displayed.\n%s}",
		forcesReplacement(diff.Replace, opts.OverrideForcesReplacement), cachedLinePrefix, cachedLinePrefix, cachedLinePrefix)
}

func (renderer sensitiveBlockRenderer) WarningsHuman(diff computed.Diff, indent int) []string {
	if (renderer.beforeSensitive == renderer.afterSensitive) || renderer.inner.Action == plans.Create || renderer.inner.Action == plans.Delete {
		// Only display warnings for sensitive values if they are changing from
		// being sensitive or to being sensitive and if they are not being
		// destroyed or created.
		return []string{}
	}

	var warning string
	if renderer.beforeSensitive {
		warning = fmt.Sprintf("  # [yellow]Warning[reset]: this block will no longer be marked as sensitive\n%s  # after applying this change.", formatIndent(indent))
	} else {
		warning = fmt.Sprintf("  # [yellow]Warning[reset]: this block will be marked as sensitive and will not\n%s  # display in UI output after applying this change.", formatIndent(indent))
	}

	if renderer.inner.Action == plans.NoOp {
		return []string{fmt.Sprintf("%s The value is unchanged.", warning)}
	}
	return []string{warning}
}
