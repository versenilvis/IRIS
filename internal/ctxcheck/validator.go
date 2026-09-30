package ctxcheck

func Validate(cmd, cwd string, d Dialect) Verdict {
	parsed := Parse(cmd, d)
	if parsed.IsUnknown {
		return Unknown
	}
	if len(parsed.Segments) == 0 {
		return Free
	}

	hasUnknown := false
	hasFree := false

	for _, seg := range parsed.Segments {
		if seg.IsUnknown {
			hasUnknown = true
			continue
		}
		if len(seg.Tokens) == 0 {
			continue
		}

		v := validateSegment(seg.Tokens, cwd)
		switch v {
		case Invalid:
			return Invalid
		case Unknown:
			hasUnknown = true
		case Free:
			hasFree = true
		case Valid:
			// continue checking remaining segments
		}
	}

	if hasUnknown {
		return Unknown
	}
	if hasFree {
		return Free
	}
	return Valid
}

func validateSegment(tokens []string, cwd string) Verdict {
	if len(tokens) == 0 {
		return Free
	}

	cmdWord := tokens[0]

	if cmdWord == "just" {
		return ValidateJust(tokens, cwd)
	}

	if nodePackageManagers[cmdWord] || cmdWord == "npx" || cmdWord == "bunx" {
		return ValidateNode(tokens, cwd)
	}

	if cmdWord == "make" {
		return ValidateMake(tokens, cwd)
	}

	if cmdWord == "cd" {
		return ValidateCd(tokens, cwd)
	}

	return ValidatePathTokens(tokens, cwd)
}
