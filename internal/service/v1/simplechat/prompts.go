package simplechat

import "strings"

func mergePrompts(systemPrompt, userPrompt, vaultPrompt string, useSystemPrompt bool) string {
	sections := make([]string, 0, 3)
	if useSystemPrompt && strings.TrimSpace(systemPrompt) != "" {
		sections = append(sections, "System instructions:\n"+strings.TrimSpace(systemPrompt))
	}
	if strings.TrimSpace(userPrompt) != "" {
		sections = append(sections, "User instructions:\n"+strings.TrimSpace(userPrompt))
	}
	if strings.TrimSpace(vaultPrompt) != "" {
		sections = append(sections, "Vault instructions:\n"+strings.TrimSpace(vaultPrompt))
	}
	return strings.Join(sections, "\n\n")
}
