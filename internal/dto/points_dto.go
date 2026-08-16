package dto

// PointsRuleCreateRequest 积分规则创建请求。
type PointsRuleCreateRequest struct {
	Name          string `json:"name" binding:"required,max=128"`
	RuleType      string `json:"rule_type" binding:"required,max=32"`
	Points        int    `json:"points" binding:"required,min=1"`
	EffectiveDate string `json:"effective_date" binding:"max=16"`
	Enabled       bool   `json:"enabled"`
	Description   string `json:"description" binding:"max=255"`
}

// PointsRuleUpdateRequest 积分规则更新请求。
type PointsRuleUpdateRequest struct {
	Name          string `json:"name" binding:"max=128"`
	RuleType      string `json:"rule_type" binding:"max=32"`
	Points        int    `json:"points" binding:"min=1"`
	EffectiveDate string `json:"effective_date" binding:"max=16"`
	Enabled       *bool  `json:"enabled"`
	Description   string `json:"description" binding:"max=255"`
}

// GrantPointsRequest 手动发放积分请求（HR）。
type GrantPointsRequest struct {
	UserID      uint   `json:"user_id" binding:"required"`
	Amount      int    `json:"amount" binding:"required,min=1"`
	Description string `json:"description" binding:"max=255"`
}
