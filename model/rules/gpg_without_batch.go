package rules

import (
	"strings"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/validations"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type gpgWithoutBatch struct {
}

func (g *gpgWithoutBatch) Name() string {
	return "gpg-without-batch"
}

func (g *gpgWithoutBatch) Summary() string {
	return "GPG call without --batch (or --no-tty) may error."
}

func (g *gpgWithoutBatch) Details() string {
	return "Running GPG without --batch (or --no-tty) may cause GPG to fail opening /dev/tty, resulting in docker build failures."
}

func (g *gpgWithoutBatch) Priority() model.Priority {
	return model.MediumPriority
}

func (g *gpgWithoutBatch) Commands() []commands.DockerCommand {
	return []commands.DockerCommand{commands.Run}
}

func (g *gpgWithoutBatch) Category() *string {
	return nil
}

func (g *gpgWithoutBatch) URL() *string {
	return model.StringPtr("https://bugs.debian.org/cgi-bin/bugreport.cgi?bug=913614")
}

func (g *gpgWithoutBatch) LintID() string {
	return validations.LintID(g)
}

func (g *gpgWithoutBatch) Evaluate(node *parser.Node, validationContext validations.ValidationContext) *validations.ValidationResult {
	trimStart := len(node.Value) + 1 // command plus trailing space
	matchAgainst := node.Original[trimStart:]
	result := model.Success
	if model.NewPattern(`(?s)\bgpg\b.*?--recv-keys.*?`).Matches(matchAgainst) {
		hasBatch := strings.Contains(matchAgainst, "--batch")
		hasNoTty := strings.Contains(matchAgainst, "--no-tty")

		if !hasBatch || !hasNoTty {
			result = model.Failure
			validationContext.CausedFailure = true
		}
	}
	return &validations.ValidationResult{
		Result:   result,
		Details:  g.Summary(),
		Contexts: []validations.ValidationContext{validationContext},
	}
}

func init() {
	r := gpgWithoutBatch{}
	AddRule(&r)
}
