package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// MedScan is the flagship model-as-a-file example: write a radiology report
// to /ai/models/medscan/infer and the triage streams back as structured
// tokens on the resource's watch stream. Deterministic + dependency-free,
// it proves the inference-resource flow end to end without an external LLM.

type MedScanResult struct {
	Triage    string   `json:"triage"`     // low | moderate | urgent
	RiskScore int      `json:"risk_score"` // 0-100
	Findings  []string `json:"findings"`
	FollowUp  string   `json:"follow_up"`
}

var (
	reUrgent   = regexp.MustCompile(`(?i)\b(mass|nodule\s*[>=]\s*8|malignant|metasta|spiculated|hemorrh|pulmonary\s+embol|pericardial|pleural\s+effusion)\b`)
	reModerate = regexp.MustCompile(`(?i)\b(nodule|opacit|consolidat|infiltrate|atelectas|excrescence|pleural\s+thick)\b`)
	reLow      = regexp.MustCompile(`(?i)\b(normal|clear|no\s+acute|stable|benign|unchanged|unremarkable)\b`)
)

// inferMedScan analyzes a report text and emits a streaming report through
// the given node's emit path (mirroring how real token streaming would work).
func inferMedScan(path, prompt string, n *Node) *MedScanResult {
	mod := len(reModerate.FindAllStringIndex(prompt, -1))
	urg := len(reUrgent.FindAllStringIndex(prompt, -1))

	var findings []string
	var triage string
	score := 0
	for _, m := range reUrgent.FindAllString(prompt, -1) {
		findings = append(findings, strings.ToLower(m))
		score += 28
	}
	for _, m := range reModerate.FindAllString(prompt, -1) {
		findings = append(findings, strings.ToLower(m))
		score += 14
	}
	for range reLow.FindAllString(prompt, -1) {
		score -= 8
	}
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	switch {
	case urg > 0 || score >= 70:
		triage = "urgent"
	case mod > 0 || score >= 25:
		triage = "moderate"
	default:
		triage = "low"
	}

	var followUp string
	switch triage {
	case "urgent":
		followUp = "Immediate pulmonology consult; confirmatory CT within 24h."
	case "moderate":
		followUp = "Re-image in 6-12 months; correlate clinically."
	default:
		followUp = "Routine follow-up; no acute intervention indicated."
	}

	// stream the reasoning/result through the resource a little at a time,
	// mirroring token streaming: POST -> infer -> watchers see tokens.
	n.emit(fmt.Sprintf("model medscan: analyzing %d findings", len(findings)))
	for _, f := range findings {
		n.emit(fmt.Sprintf("finding: %s", f))
		time.Sleep(30 * time.Millisecond)
	}
	n.emit(fmt.Sprintf("triage: %s (risk %d/100)", triage, score))
	n.emit(fmt.Sprintf("follow_up: %s", followUp))
	n.emit(fmt.Sprintf("done %s", path))

	return &MedScanResult{
		Triage:    triage,
		RiskScore: score,
		Findings:  findings,
		FollowUp:  followUp,
	}
}
