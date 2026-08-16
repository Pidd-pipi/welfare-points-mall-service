package util

// points_calculator.go 积分扣减/返还与兑换金额计算工具（前后端 pointsCalculator.ts 对应）。
// ExchangeCost 兑换所需积分 = 商品积分单价 * 数量。
func ExchangeCost(pointsCost, quantity int) int {
	if quantity <= 0 {
		return 0
	}
	return pointsCost * quantity
}

// RefundPoints 取消订单返还积分。
func RefundPoints(pointsCost, quantity int) int {
	return ExchangeCost(pointsCost, quantity)
}

// CanAfford 判断积分余额是否足够。
func CanAfford(balance, cost int) bool {
	return balance >= cost
}

// GrantPoints 发放积分后的余额。
func GrantPoints(balance, amount int) int {
	return balance + amount
}

// DeductPoints 扣减积分后的余额（不允许扣成负数，由上层校验）。
func DeductPoints(balance, amount int) int {
	if balance < amount {
		return balance
	}
	return balance - amount
}
