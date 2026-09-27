package main

import (
	"strings"
	"testing"
)

func testState() *State {
	return &State{
		Version: 1,
		Sources: map[string]*SourceState{
			"otnc/skills_fix-unnatural-line-breaks": {
				Tag:     "v0.5.0",
				Dirs:    []string{"fix-unnatural-line-breaks"},
				License: &LicenseInfo{SpdxID: "MIT", URL: "https://github.com/otnc/skills_fix-unnatural-line-breaks/blob/main/LICENSE"},
			},
			"otnc/skills_ai-agent-saving-tech": {
				Tag:     "v0.1.0",
				Dirs:    []string{"token-saver-codex", "token-saver"},
				License: &LicenseInfo{SpdxID: "MIT", URL: "https://github.com/otnc/skills_ai-agent-saving-tech/blob/main/LICENSE"},
			},
		},
	}
}

func TestSkillRowsSorted(t *testing.T) {
	rows := skillRows(testState())
	var got []string
	for _, r := range rows {
		got = append(got, r.repo+"/"+r.skill)
	}
	want := []string{
		"otnc/skills_ai-agent-saving-tech/token-saver",
		"otnc/skills_ai-agent-saving-tech/token-saver-codex",
		"otnc/skills_fix-unnatural-line-breaks/fix-unnatural-line-breaks",
	}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("skillRows() = %v, want %v", got, want)
	}
}

func TestRenderTable(t *testing.T) {
	rows := skillRows(testState())
	en := renderTable("en", rows)
	if !strings.HasPrefix(en, "| Skill | Source repository | License |\n| --- | --- | --- |") {
		t.Errorf("renderTable(en) has wrong header: %q", en)
	}
	if !strings.Contains(en, "| `token-saver` | [otnc/skills_ai-agent-saving-tech](https://github.com/otnc/skills_ai-agent-saving-tech) | [MIT](https://github.com/otnc/skills_ai-agent-saving-tech/blob/main/LICENSE) |") {
		t.Errorf("renderTable(en) missing row: %q", en)
	}
	ja := renderTable("ja", rows)
	if !strings.HasPrefix(ja, "| スキル | 元のリポジトリ | ライセンス |") {
		t.Errorf("renderTable(ja) has wrong header: %q", ja)
	}
}

func TestLicenseLabel(t *testing.T) {
	cases := []struct {
		license *LicenseInfo
		want    string
	}{
		{&LicenseInfo{SpdxID: "MIT", URL: "https://example.com/LICENSE"}, "[MIT](https://example.com/LICENSE)"},
		{&LicenseInfo{SpdxID: "NOASSERTION", URL: "https://example.com/LICENSE"}, "[Other](https://example.com/LICENSE)"},
		{nil, "-"},
	}
	for _, c := range cases {
		if got := licenseLabel(c.license); got != c.want {
			t.Errorf("licenseLabel(%v) = %q, want %q", c.license, got, c.want)
		}
	}
	if !licenseEqual(nil, nil) {
		t.Error("licenseEqual(nil, nil) = false, want true")
	}
	if licenseEqual(nil, &LicenseInfo{}) {
		t.Error("licenseEqual(nil, ...) = true, want false")
	}
	if !licenseEqual(&LicenseInfo{SpdxID: "MIT", URL: "u"}, &LicenseInfo{SpdxID: "MIT", URL: "u"}) {
		t.Error("licenseEqual() = false for identical licenses")
	}
}

func TestRenderList(t *testing.T) {
	list := renderList(sourceRows(testState()))
	want := "- [otnc/skills_ai-agent-saving-tech](https://github.com/otnc/skills_ai-agent-saving-tech) ([MIT](https://github.com/otnc/skills_ai-agent-saving-tech/blob/main/LICENSE))\n" +
		"- [otnc/skills_fix-unnatural-line-breaks](https://github.com/otnc/skills_fix-unnatural-line-breaks) ([MIT](https://github.com/otnc/skills_fix-unnatural-line-breaks/blob/main/LICENSE))"
	if list != want {
		t.Errorf("renderList() = %q, want %q", list, want)
	}

	// A repository without a license file gets no parenthetical.
	st := testState()
	st.Sources["otnc/skills_ai-agent-saving-tech"].License = nil
	list = renderList(sourceRows(st))
	if !strings.Contains(list, "- [otnc/skills_ai-agent-saving-tech](https://github.com/otnc/skills_ai-agent-saving-tech)\n") {
		t.Errorf("renderList() missing unlicensed row: %q", list)
	}
}

func TestReplaceSection(t *testing.T) {
	content := "a\n<!-- m:start -->\nold\n<!-- m:end -->\nb\n<!-- m:start -->\nolder\n<!-- m:end -->\nc\n"
	got, changed, err := replaceSection(content, "m", "new")
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Error("replaceSection() reported no change")
	}
	want := "a\n<!-- m:start -->\nnew\n<!-- m:end -->\nb\n<!-- m:start -->\nnew\n<!-- m:end -->\nc\n"
	if got != want {
		t.Errorf("replaceSection() = %q, want %q", got, want)
	}

	// Identical body: no change reported.
	if got, changed, err = replaceSection(want, "m", "new"); err != nil {
		t.Fatal(err)
	}
	if changed || got != want {
		t.Errorf("replaceSection() identical body: changed=%v", changed)
	}

	// Missing markers are an error.
	if _, _, err := replaceSection("no markers here", "m", "new"); err == nil {
		t.Error("replaceSection() missing start marker: want error")
	}
	if _, _, err := replaceSection("<!-- m:start -->\nno end\n", "m", "new"); err == nil {
		t.Error("replaceSection() missing end marker: want error")
	}
}
