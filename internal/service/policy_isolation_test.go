package service

import (
	"testing"

	"datamasking/internal/config"
	"datamasking/internal/model"
	"datamasking/internal/store"
	"datamasking/pkg/logger"
)

func TestPolicyEvaluationDoesNotExposeStoredRules(t *testing.T) {
	svc := New(store.NewMemoryStore(), logger.NewLevel(logger.LevelError), &config.Config{MaxPageSize: 100})
	rule, err := svc.CreateMaskRule(model.MaskRule{Name: "phone", MaskType: "mask", PrefixKeep: 3, SuffixKeep: 4, Enabled: true})
	if err != nil { t.Fatalf("create rule: %v", err) }
	policy, err := svc.CreatePolicy(model.Policy{Name: "default", Scope: "all", RuleIDs: []string{rule.ID}, Enabled: true})
	if err != nil { t.Fatalf("create policy: %v", err) }
	selected, err := svc.EvaluatePolicy(policy.ID)
	if err != nil { t.Fatalf("evaluate policy: %v", err) }
	selected[0].MaskType = model.MaskTypeReplace
	selected[0].Replacement = "LEAK"
	masked, err := svc.ApplyMaskRule(rule.ID, "13800138000")
	if err != nil { t.Fatalf("apply rule: %v", err) }
	if masked != "138****8000" { t.Fatalf("policy result mutated stored rule: %q", masked) }
}
