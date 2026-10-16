package skills

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSkillsLoader_Lifecycle(t *testing.T) {
	tempDir := t.TempDir()
	loader := NewLoader()

	// Initial load should be empty
	skills, err := loader.LoadSkills(tempDir)
	if err != nil {
		t.Fatalf("unexpected error loading empty skills: %v", err)
	}
	if len(skills) != 0 {
		t.Fatalf("expected 0 skills, got %d", len(skills))
	}

	prompt := loader.FormatSkillsPrompt(skills)
	if prompt != "" {
		t.Fatalf("expected empty prompt for 0 skills")
	}

	// Ensure default skills
	loader.EnsureDefaultSkills(tempDir)

	skillFile := filepath.Join(tempDir, ".agent", "skills", "code-quality.md")
	if _, err := os.Stat(skillFile); os.IsNotExist(err) {
		t.Fatalf("default skill file not created: %s", skillFile)
	}

	// Add custom skill
	customSkillPath := filepath.Join(tempDir, ".agent", "skills", "security-audit.md")
	customSkillContent := `# Security Audit Guide
1. 校验输入路径，防止目录遍历。
2. 阻断高危系统指令。`
	if err := os.WriteFile(customSkillPath, []byte(customSkillContent), 0644); err != nil {
		t.Fatalf("failed to write custom skill: %v", err)
	}

	// Load skills
	loadedSkills, err := loader.LoadSkills(tempDir)
	if err != nil {
		t.Fatalf("LoadSkills failed: %v", err)
	}
	if len(loadedSkills) != 2 {
		t.Fatalf("expected 2 skills, got %d", len(loadedSkills))
	}

	// Format prompt
	formattedPrompt := loader.FormatSkillsPrompt(loadedSkills)
	if !strings.Contains(formattedPrompt, "工作空间专属技能库") {
		t.Fatalf("expected header in formatted prompt, got: %s", formattedPrompt)
	}
	if !strings.Contains(formattedPrompt, "Security Audit Guide") {
		t.Fatalf("expected custom skill title in prompt, got: %s", formattedPrompt)
	}
	if !strings.Contains(formattedPrompt, "Code Quality & Testing Guide") {
		t.Fatalf("expected default skill title in prompt, got: %s", formattedPrompt)
	}
}
