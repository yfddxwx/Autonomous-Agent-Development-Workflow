package skills

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"agent/pkg/models"
)

// Loader discovers and loads local skill markdown documents from .agent/skills/ and skills/.
type Loader struct{}

// NewLoader creates a new Loader instance.
func NewLoader() *Loader {
	return &Loader{}
}

// EnsureDefaultSkills creates sample skills in the workspace if none exist.
func (l *Loader) EnsureDefaultSkills(workspaceDir string) {
	skillsDir := filepath.Join(workspaceDir, ".agent", "skills")
	_ = os.MkdirAll(skillsDir, 0755)

	skillFile := filepath.Join(skillsDir, "code-quality.md")
	if _, err := os.Stat(skillFile); os.IsNotExist(err) {
		defaultSkill := `# Code Quality & Testing Guide

## 核心行为准则:
1. 每次编写或修改代码时，严格遵守单一职责原则（SRP），杜绝超长函数。
2. 修改已有代码必须使用 file_edit 进行微创块级精准差异替换，严禁非必要全文件重写。
3. 编码完成后必须通过终端执行编译与单元测试（如 go test ./...），确保 100% 验证通过。
4. 遇到测试用例报错，主动阅读出错信息并反思根本原因，精准修正。
`
		_ = os.WriteFile(skillFile, []byte(defaultSkill), 0644)
	}

	codeQualityDir := filepath.Join(skillsDir, "code-quality")
	stdSkillFile := filepath.Join(codeQualityDir, "SKILL.md")
	if _, err := os.Stat(stdSkillFile); os.IsNotExist(err) {
		_ = os.MkdirAll(codeQualityDir, 0755)
		stdSkill := `---
name: code-quality
description: 保证工程代码质量、模块高内聚低耦合与自动化单元测试通过
---

# Code Quality & Testing Guide

## 核心行为准则:
1. 每次编写或修改代码时，严格遵守单一职责原则（SRP），杜绝超长函数。
2. 修改已有代码必须使用 file_edit 进行微创块级精准差异替换，严禁非必要全文件重写。
3. 编码完成后必须通过终端执行编译与单元测试（如 go test ./...），确保 100% 验证通过。
4. 遇到测试用例报错，主动阅读出错信息并反思根本原因，精准修正。
`
		_ = os.WriteFile(stdSkillFile, []byte(stdSkill), 0644)
	}
}

// LoadSkills scans the workspace for local skills (supporting both SKILL.md folders and .md files).
func (l *Loader) LoadSkills(workspaceDir string) ([]models.SkillDefinition, error) {
	candidateDirs := []string{
		filepath.Join(workspaceDir, ".agent", "skills"),
		filepath.Join(workspaceDir, "skills"),
	}

	seenNames := make(map[string]bool)
	var skills []models.SkillDefinition

	for _, dir := range candidateDirs {
		info, err := os.Stat(dir)
		if err != nil || !info.IsDir() {
			continue
		}

		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			fullPath := filepath.Join(dir, entry.Name())

			// Case 1: Subdirectory with SKILL.md (standard Agent Skills format)
			if entry.IsDir() {
				skillMD := filepath.Join(fullPath, "SKILL.md")
				if _, err := os.Stat(skillMD); os.IsNotExist(err) {
					skillMD = filepath.Join(fullPath, "skill.md")
				}
				if _, err := os.Stat(skillMD); err == nil {
					key := strings.ToLower(entry.Name())
					if def, ok := l.parseSkillFile(skillMD, entry.Name()); ok {
						if !seenNames[key] {
							seenNames[key] = true
							skills = append(skills, def)
						}
					}
				}
				continue
			}

			// Case 2: Flat markdown file (.md)
			if strings.HasSuffix(strings.ToLower(entry.Name()), ".md") {
				baseName := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
				key := strings.ToLower(baseName)
				if def, ok := l.parseSkillFile(fullPath, baseName); ok {
					if !seenNames[key] {
						seenNames[key] = true
						skills = append(skills, def)
					}
				}
			}
		}
	}

	return skills, nil
}

func (l *Loader) parseSkillFile(filePath, defaultName string) (models.SkillDefinition, bool) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return models.SkillDefinition{}, false
	}

	raw := string(bytes)
	name := defaultName
	desc := "工作区专属自主研发技能"
	body := raw

	// Check for YAML Frontmatter (--- ... ---)
	if strings.HasPrefix(strings.TrimSpace(raw), "---") {
		parts := strings.SplitN(raw, "---", 3)
		if len(parts) >= 3 {
			frontmatter := parts[1]
			body = strings.TrimSpace(parts[2])

			for _, line := range strings.Split(frontmatter, "\n") {
				trimmed := strings.TrimSpace(line)
				if strings.HasPrefix(trimmed, "name:") {
					name = strings.TrimSpace(strings.TrimPrefix(trimmed, "name:"))
				} else if strings.HasPrefix(trimmed, "description:") {
					desc = strings.TrimSpace(strings.TrimPrefix(trimmed, "description:"))
				}
			}
		}
	}

	// Fallback to title line
	if name == defaultName {
		for _, line := range strings.Split(raw, "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "# ") {
				name = strings.TrimPrefix(trimmed, "# ")
				break
			}
		}
	}

	return models.SkillDefinition{
		Name:        name,
		Description: desc,
		Content:     body,
		Path:        filePath,
	}, true
}

// FormatSkillsPrompt formats loaded skills into a system prompt segment.
func (l *Loader) FormatSkillsPrompt(skills []models.SkillDefinition) string {
	if len(skills) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("\n\n【本地专属技能即插即用解析器 (工作空间专属技能库 - Agent Skills 标准支持)】:\n")
	sb.WriteString("你所在的运行环境已完全挂载并激活了 Agent Skills 解析引擎，支持读取工作区 `.agent/skills/` 与 `skills/` 目录下的所有技能定义（包括带 SKILL.md 的独立技能文件夹与 Markdown 规范）。\n")
	sb.WriteString(fmt.Sprintf("当前已自动扫描并挂载了 %d 项本地工程技能，请在执行任务与代码重构时严格遵循与调用：\n", len(skills)))
	for _, s := range skills {
		sb.WriteString(fmt.Sprintf("\n--- [Skill: %s (%s)] ---\n%s\n", s.Name, s.Description, s.Content))
	}

	return sb.String()
}
