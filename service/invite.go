package service

type CommonService struct {
}

func NewCommonService() *CommonService {
	return &CommonService{}
}

// GadgetPoints
// When a user purchases a product, a GadgetPoints entry is generated based on the quantity of the
// purchased items, and the Customer's GadgetPoints value is increased accordingly.
// When a user redeems points, a PointRedemption entry is generated, and when the company replenishes
// the inventory, the Replenished field of the PointRedemption is set to true.

// Subscription
// When a user purchases a monthly subscription, a MonthlySubscription entry is created, and the
// Customer's SubscriptionExpiry field is set based on the type of monthly subscription purchased.
// When a user cancels a subscription, the Customer's SubscriptionExpiry field is set to the current
// time or the end of the subscription period, depending on whether the cancellation includes a refund.

// Discount
// The Discount table records various discount types (national discount, member discount, store discount)
// and specifies the associated products. Therefore, when a user makes a purchase, the discount amount is
// calculated based on the product's discount type, which may or may not stack with other discounts.
