package store

// cloneRuleIDs 返回切片的拷贝，使其与调用方传入的底层数组互不影响。
// 用于在持久化 MaskTask/Policy 等持有 RuleIDs 切片的实体时做防御性拷贝，
// 避免外部请求复用并修改同一切片后污染已保存记录、或引发删除规则时的引用冲突。
func cloneRuleIDs(ids []string) []string {
	if ids == nil {
		return nil
	}
	cp := make([]string, len(ids))
	copy(cp, ids)
	return cp
}
