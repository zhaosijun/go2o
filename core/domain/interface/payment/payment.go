/**
 * Copyright 2015 @ to2.net.
 * name : payment
 * author : jarryliu
 * date : 2016-07-02 23:06
 * description : 支付单据
 * history :
 */

// 支付单,不限于订单,可以生成支付单,即一个支付请求
package payment

import (
	"go2o/core/domain/interface/promotion"
	"go2o/core/infrastructure/domain"
)

// 支付通道
const (
	// 余额抵扣通道
	MBalance = 1 << 0
	// 钱包支付通道
	MWallet = 1 << 1
	// 积分兑换通道
	MIntegral = 1 << 2
	// 用户卡通道
	MUserCard = 1 << 3
	// 用户券通道
	MUserCoupon = 1 << 4
	// 现金支付通道
	MCash = 1 << 5
	// 银行卡支付通道
	MBankCard = 1 << 6
	// 第三方支付
	MPaySP = 1 << 7
	// 卖家支付通道
	MSellerPay = 1 << 8
	// 系统支付通道
	MSystemPay = 1 << 9
)

// 所有支付方式
const PAllFlag = MBalance | MWallet | MIntegral | MUserCard |
	MUserCoupon | MCash | MBankCard | MPaySP | MSellerPay | MSystemPay

// 支付单状态
const (
	// 待支付
	StateAwaitingPayment = 1
	// 已支付
	StateFinished = 2
	// 已取消
	StateCancelled = 3
	// 已终止（超时关闭）
	StateAborted = 4
)

var (
	ErrNoSuchPaymentOrder = domain.NewError(
		"err_no_such_payment_order", "支付单不存在")

	ErrExistsTradeNo = domain.NewError(
		"err_payment_exists_trade_no", "支付单号重复")

	ErrPaymentNotSave = domain.NewError(
		"err_payment_not_save", "支付单需存后才能执行操作")

	ErrFinalFee = domain.NewError(
		"err_final_fee", "支付单金额有误")

	ErrNotSupportPaymentChannel = domain.NewError(
		"err_payment_not_support_channel", "不支持此支付方式,无法完成付款")
	ErrItemAmount    = domain.NewError("err_payment_item_amount", "支付单金额不能为零")
	ErrOutOfFinalFee = domain.NewError("err_out_of_final_fee",
		"超出支付单金额")
	ErrNotMatchFinalFee = domain.NewError("err_not_match_final_fee",
		"金额与实际金额不符，无法完成付款")
	ErrTradeNoPrefix = domain.NewError(
		"err_payment_trade_no_prefix", "支付单号前缀不正确")
	ErrTradeNoExistsPrefix = domain.NewError(
		"err_payment_trade_no_exists_prefix", "支付单号已存在前缀")

	ErrOrderCommitted = domain.NewError(
		"err_payment_order_committed", "支付单已提交")

	ErrOrderPayed = domain.NewError(
		"err_payment_order_payed", "订单已支付")

	ErrOrderCancelled = domain.NewError("err_payment_order_has_cancel", "订单已经取消")

	ErrOrderNotPayed = domain.NewError("err_payment_order_not_payed", "订单未支付")

	ErrCanNotUseBalance = domain.NewError("err_can_not_use_balance", "不能使用余额支付")

	ErrNotEnoughAmount = domain.NewError("err_payment_not_enough_amount", "余额不足,无法完成支付")

	ErrCanNotUseIntegral = domain.NewError("err_can_not_use_integral", "不能使用积分抵扣")

	ErrCanNotUseCoupon = domain.NewError("err_can_not_use_coupon", "不能使用优惠券")

	ErrCanNotSystemDiscount = domain.NewError("err_can_not_system_discount", "不允许系统支付")

	ErrOuterNo = domain.NewError("err_outer_no", "第三方交易号错误")
)

type (
	// 支付单接口
	IPaymentOrder interface {
		// 获取聚合根编号
		GetAggregateRootId() int
		// 获取支付单的值
		Get() Order
		// 获取交易号
		TradeNo() string
		// 支付单状态
		State() int
		// 支付方式
		Flag() int
		// 支付途径支付信息
		TradeMethods() []*TradeMethodData
		// 在支付之前检查订单状态
		CheckPaymentState() error
		// 提交支付单
		Submit() error
		// 合并支付
		MergePay(orders []IPaymentOrder) (mergeTradeNo string, finalFee int, err error)
		// 取消支付
		Cancel() error
		// 线下现金/刷卡支付,cash:现金,bank:刷卡金额,finalZero:是否金额必须为零
		OfflineDiscount(cash int, bank int, finalZero bool) error
		// 交易完成
		TradeFinish() error
		// 支付完成并保存,传入第三名支付名称,以及外部的交易号
		PaymentFinish(spName string, outTradeNo string) error
		// 优惠券抵扣
		CouponDiscount(coupon promotion.ICouponPromotion) (int, error)
		// 使用会员的余额抵扣
		BalanceDiscount(remark string) error
		// 使用会员积分抵扣,返回抵扣的金额及错误,ignoreOut:是否忽略超出订单金额的积分
		IntegralDiscount(integral int, ignoreOut bool) (amount int, err error)
		// 系统支付金额
		SystemPayment(amount int) error
		// 钱包账户支付
		PaymentByWallet(remark string) error
		// 使用会员卡支付,cardCode:会员卡编码,amount:支付金额
		PaymentWithCard(cardCode string, amount int) error
		// 余额钱包混合支付，优先扣除余额。
		HybridPayment(remark string) error
		// 设置支付方式
		SetTradeSP(spName string) error

		// 调整金额,如调整金额与实付金额相加小于等于零,则支付成功。
		Adjust(amount int) error
		// 退款
		Refund(amount int) error
		// 获取支付通道字符串
		ChanName(method int) string
	}

	// 支付仓储
	IPaymentRepo interface {
		// 根据编号获取支付单
		GetPaymentOrderById(id int) IPaymentOrder
		// 根据支付单号获取支付单
		GetPaymentOrder(tradeNo string) IPaymentOrder
		// 根据订单号获取支付单
		GetPaymentBySalesOrderId(orderId int64) IPaymentOrder
		// 根据支付单号获取支付单
		GetPaymentOrderByOrderNo(orderType int, orderNo string) IPaymentOrder
		// 创建支付单
		CreatePaymentOrder(p *Order) IPaymentOrder
		// 保存支付单
		SavePaymentOrder(v *Order) (int, error)
		// 检查支付单号是否匹配
		CheckTradeNoMatch(tradeNo string, id int) bool
		// 获取交易途径支付信息
		GetTradeChannelItems(tradeNo string) []*TradeMethodData
		// 保存支付途径支付信息
		SavePaymentTradeChan(tradeNo string, tradeChan *TradeMethodData) (int, error)
		// 获取合并支付的订单
		GetMergePayOrders(mergeTradeNo string) []IPaymentOrder
		// 清除欲合并的支付单
		ResetMergePaymentOrders(tradeNos []string) error
		//  保存合并的支付单
		SaveMergePaymentOrders(s string, tradeNos []string) error
	}

	// 请求支付数据
	RequestPayData struct {
		// 支付方式
		method int
		// 支付方式代码
		code string
		// 支付金额
		amount int
	}

	// 支付单
	Order struct {
		// 编号
		ID int `db:"id" pk:"yes" auto:"yes"`
		// 卖家编号
		SellerId int `db:"seller_id"`
		// 交易类型
		TradeType string `db:"trade_type"`
		// 交易号
		TradeNo string `db:"trade_no"`
		// 支付单的类型，如购物或其他
		OrderType int `db:"order_type"`
		// 是否为子订单
		SubOrder int `db:"sub_order"`
		// 外部订单号
		OutOrderNo string `db:"out_order_no"`
		// 支付单详情
		Subject string `db:"subject"`
		// 买家编号
		BuyerId int64 `db:"buyer_id"`
		// 支付用户编号
		PayUid int64 `db:"pay_uid"`
		// 商品金额
		ItemAmount int `db:"item_amount"`
		// 优惠金额
		DiscountAmount int `db:"discount_amount"`
		// 调整金额
		AdjustAmount int `db:"adjust_amount"`
		// 共计金额，包含抵扣金额
		TotalAmount int `db:"total_amount"`
		// 抵扣金额
		DeductAmount int `db:"deduct_amount"`
		// 手续费
		ProcedureFee int `db:"procedure_fee"`
		// 最终支付金额，包含手续费，不包含抵扣金额
		FinalFee int `db:"final_fee"`
		// 实付金额
		PaidFee int `db:"paid_fee"`
		// 可⽤支付方式
		PayFlag int `db:"pay_flag"`
		// 实际支付方式
		FinalFlag int `db:"final_flag"`
		// 其他支付信息
		ExtraData string `db:"extra_data"`
		// 交易支付渠道
		TradeChannel int `db:"trade_channel"`
		// 外部交易提供商
		OutTradeSp string `db:"out_trade_sp"`
		// 外部交易订单号
		OutTradeNo string `db:"out_trade_no"`
		// 订单状态
		State int `db:"state"`
		// 提交时间
		SubmitTime int64 `db:"submit_time"`
		// 过期时间
		ExpiresTime int64 `db:"expires_time"`
		// 支付时间
		PaidTime int64 `db:"paid_time"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
		// 交易途径支付信息
		TradeMethods []*TradeMethodData `db:"-"`
	}

	// 支付单项
	TradeMethodData struct {
		// 编号
		ID int `db:"id" pk:"yes" auto:"yes"`
		// 交易单号
		TradeNo string `db:"trade_no"`
		// 支付途径
		Method int `db:"pay_method"`
		// 支付代码
		Code string `db:"pay_code"`
		// 是否为内置支付途径
		Internal int `db:"internal"`
		// 支付金额
		Amount int `db:"pay_amount"`
		// 外部交易单号
		OutTradeNo string `db:"out_trade_no"`
		// 支付时间
		PayTime int64 `db:"pay_time"`
	}

	// 合并的支付单
	MergeOrder struct {
		// 编号
		ID int `db:"id"`
		// 合并交易单号
		MergeTradeNo string `db:"merge_trade_no"`
		// 交易号
		OrderTradeNo string `db:"order_trade_no"`
		// 提交时间
		SubmitTime int64 `db:"submit_time"`
	}

	// SP支付交易
	PaySpTrade struct {
		// 编号
		ID int `db:"id"`
		// 交易SP
		TradeSp string `db:"trade_sp"`
		// 交易号
		TradeNo string `db:"trade_no"`
		// 合并的订单号,交易号用"|"分割
		TradeOrders string `db:"trade_orders"`
		// 交易状态
		TradeState int `db:"trade_state"`
		// 交易结果
		TradeResult int `db:"trade_result"`
		// 交易备注
		TradeRemark string `db:"trade_remark"`
		// 交易时间
		TradeTime int `db:"trade_time"`
	}
)
