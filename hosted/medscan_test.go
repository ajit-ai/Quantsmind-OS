package main

import "testing"

func TestMedScanUrgent(t *testing.T) {
	node := newNode("infer")
	prompt := "Chest CT: 9mm spiculated nodule right upper lobe, malignant features, pleural effusion."
	res := inferMedScan("/ai/models/medscan/infer", prompt, node)
	if res.Triage != "urgent" {
		t.Fatalf("expected urgent triage, got %q (score %d)", res.Triage, res.RiskScore)
	}
	if res.RiskScore < 70 {
		t.Fatalf("expected high risk score, got %d", res.RiskScore)
	}
}

func TestMedScanModerate(t *testing.T) {
	node := newNode("infer")
	prompt := "Small 4mm nodule in the left base, mild opacity, otherwise stable."
	res := inferMedScan("/ai/models/medscan/infer", prompt, node)
	if res.Triage != "moderate" {
		t.Fatalf("expected moderate triage, got %q (score %d)", res.Triage, res.RiskScore)
	}
}

func TestMedScanLow(t *testing.T) {
	node := newNode("infer")
	prompt := "Lungs clear bilaterally. No acute findings. Unremarkable exam."
	res := inferMedScan("/ai/models/medscan/infer", prompt, node)
	if res.Triage != "low" {
		t.Fatalf("expected low triage, got %q (score %d)", res.Triage, res.RiskScore)
	}
}

func TestMedScanClampsScore(t *testing.T) {
	node := newNode("infer")
	// many urgent keywords -> would exceed 100 without clamping
	prompt := "mass mass mass nodule mass malignant metastatic spiculated hemorrh"
	res := inferMedScan("/ai/models/medscan/infer", prompt, node)
	if res.RiskScore > 100 {
		t.Fatalf("risk score must not exceed 100, got %d", res.RiskScore)
	}
	if res.RiskScore < 0 {
		t.Fatalf("risk score must not be negative, got %d", res.RiskScore)
	}
}

func TestMedScanEmitsThroughResource(t *testing.T) {
	node := newNode("infer")
	ch := node.Subscribe()
	defer node.Unsubscribe(ch)

	// inferMedScan emits synchronously through the node's subscriber channel
	res := inferMedScan("/ai/models/medscan/infer",
		"9mm spiculated nodule with malignant features", node)

	if res.Triage != "urgent" {
		t.Fatalf("expected urgent, got %q", res.Triage)
	}

	gotTokens := 0
	gotDone := false
	// drain until we see the done event
	for i := 0; i < 50; i++ {
		select {
		case ev := <-ch:
			if ev == "" {
				t.Fatal("expected a streamed event")
			}
			if ev == "done /ai/models/medscan/infer" {
				gotDone = true
			}
			gotTokens++
		default:
			break
		}
	}
	if gotTokens == 0 {
		t.Fatal("expected streamed tokens")
	}
	if !gotDone {
		t.Fatalf("expected stream to end with done, got %d tokens", gotTokens)
	}
}
