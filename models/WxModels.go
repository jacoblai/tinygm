package models

type AuthLoginBody struct {
	Code     string      `json:"code"`
	UserInfo ResUserInfo `json:"userInfo"`
}

type WXLoginAccessToken struct {
	AccessToken string `json:"access_token"`
	expiresIn   int    `json:"expires_in"`
}
type WXLoginGetPaidUnionId struct {
	UnionID string `json:"unionid"`
}
type Watermark struct {
	AppID     string `json:"appid"`
	TimeStamp int64  `json:"timestamp"`
}
type WXUserInfo struct {
	OpenID     string    `json:"openId,omitempty"`
	NickName   string    `json:"nickName"`
	AvatarUrl  string    `json:"avatarUrl"`
	Gender     int       `json:"gender"`
	Country    string    `json:"country"`
	Province   string    `json:"province"`
	City       string    `json:"city"`
	UnionID    string    `json:"unionId,omitempty"`
	Language   string    `json:"language"`
	Watermark  Watermark `json:"watermark,omitempty"`
	SessionKey string    `json:"session_key"`
}
type ResUserInfo struct {
	UserInfo      WXUserInfo `json:"userInfo"`
	RawData       string     `json:"rawData"`
	Signature     string     `json:"signature"`
	EncryptedData string     `json:"encryptedData"`
	IV            string     `json:"iv"`
}

// WXPayNotifyReq 当用户支付成功的时候微信返回的回调参数
// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_1
type WXPayNotifyReq struct {
	Id           string        `json:"id"`            //通知的唯一ID
	CreateTime   string        `json:"create_time"`   //通知创建的时间
	EventType    string        `json:"event_type"`    //通知的类型，支付成功通知的类型为TRANSACTION.SUCCESS
	ResourceType string        `json:"resource_type"` //通知的资源数据类型，支付成功通知为encrypt-resource
	Resource     InformTheData `json:"resource"`      //通知资源数据
	Summary      string        `json:"summary"`       //回调摘要
	//NonceStr      string `xml:"nonce_str"`
	//Openid        string `xml:"openid"`
	//OutTradeNo    string `xml:"out_trade_no"`
	//ResultCode    string `xml:"result_code"`
	//ReturnCode    string `xml:"return_code"`
	//Sign          string `xml:"sign"`
	//TimeEnd       string `xml:"time_end"`
	//TotalFee      int    `xml:"total_fee"`
	//TradeType     string `xml:"trade_type"`
	//TransactionID string `xml:"transaction_id"`
}

//通知资源数据
type InformTheData struct {
	Agorithm       string `json:"algorithm"`       //对开启结果数据进行加密的加密算法，目前只支持AEAD_AES_256_GCM
	Ciphertext     string `json:"ciphertext"`      //Base64编码后的开启/停用结果数据密文
	AssociatedData string `json:"associated_data"` //附加数据
	OriginalType   string `json:"original_type"`   //原始类型 原始回调类型，为transaction
	Nonce          string `json:"nonce"`           //随机串 加密使用的随机串

}

//type WXPayNotifyReq struct {
//	Appid         string `xml:"appid"`
//	BankType      string `xml:"bank_type"`
//	CashFee       int    `xml:"cash_fee"`
//	FeeType       string `xml:"fee_type"`
//	IsSubscribe   string `xml:"is_subscribe"`
//	MchID         string `xml:"mch_id"`
//	NonceStr      string `xml:"nonce_str"`
//	Openid        string `xml:"openid"`
//	OutTradeNo    string `xml:"out_trade_no"`
//	ResultCode    string `xml:"result_code"`
//	ReturnCode    string `xml:"return_code"`
//	Sign          string `xml:"sign"`
//	TimeEnd       string `xml:"time_end"`
//	TotalFee      int    `xml:"total_fee"`
//	TradeType     string `xml:"trade_type"`
//	TransactionID string `xml:"transaction_id"`
//}

// WXRefundNotifyReq 当用户退款成功的时候微信返回的回调参数
type WXRefundNotifyReq struct {
	Appid    string `xml:"appid"`
	MchID    string `xml:"mch_id"`
	NonceStr string `xml:"nonce_str"`
	ReqInfo  string `xml:"req_info"`
}

// WXReqInfo 当用户退款成功的时候微信Req_info 里面的加密字段
type WXReqInfo struct {
	//加密字段
	TransactionID       string `xml:"transaction_id"`        //微信订单号
	OutTradeNo          string `xml:"out_trade_no"`          //商户订单号
	RefundID            string `xml:"refund_id"`             //微信退款单号
	OutRefundNo         string `xml:"out_refund_no"`         //商户退款单号
	TotalFee            int    `xml:"total_fee"`             //订单金额
	RefundFee           int    `xml:"refund_fee"`            //申请退款金额
	SettlementRefundFee int    `xml:"settlement_refund_fee"` //退款金额 =申请退款金额-非充值代金券退款金额，退款金额<=申请退款金额
	RefundStatus        string `xml:"refund_status"`         //退款状态：SUCCESS-退款成功 CHANGE-退款异常 REFUNDCLOSE—退款关闭
	RefundRecvAccout    string `xml:"refund_recv_accout"`    //退款入账账户 1:退回银行卡： {银行名称}{卡类型}{卡尾号} 2:退回支付用户零钱: 支付用户零钱 3:退还商户: 商户基本账户 商户结算银行账户 4:退回支付用户零钱通: 支付用户零钱通
	RefundAccount       string `xml:"refund_account"`        //退款资金来源，REFUND_SOURCE_RECHARGE_FUNDS 可用余额退款/基本账户 REFUND_SOURCE_UNSETTLED_FUNDS 未结算资金退款
	RefundRequestSource string `xml:"refund_request_source"` //退款发起来源,API接口 VENDOR_PLATFORM商户平台
}

type WXRefundResV3 struct {
	RefundId            string         `json:"refund_id"`             //微信支付退款号
	OutRefundNo         string         `json:"out_refund_no"`         //商户退款订单号
	TransactionId       string         `json:"transaction_id"`        //微信支付交易订单号。
	OutTradeNo          string         `json:"out_trade_no"`          //原支付交易对应的商户订单号
	Channel             string         `json:"channel"`               //商户退款订单号
	UserReceivedAccount string         `json:"user_received_account"` //退款入账账户
	SuccessTime         string         `json:"success_time"`          //退款成功时间
	CreateTime          string         `json:"create_time"`           //退款创建时间
	Status              string         `json:"status"`                //退款状态
	FundsAccount        string         `json:"funds_account"`         //资金账户
	Amount              WxRefundAmount `json:"amount"`                //商户退款订单号
}
type WxRefundAmount struct {
	Total            int64  `json:"total"`             //订单金额
	Refund           int64  `json:"refund"`            //退款入账账户
	PayerTotal       int64  `json:"payer_total"`       //用户支付金额
	PayerRefund      int64  `json:"payer_refund"`      //用户退款金额
	SettlementRefund int64  `json:"settlement_refund"` //应结退款金额
	SettlementTotal  int64  `json:"settlement_total"`  //应结订单金额
	DiscountRefund   int64  `json:"discount_refund"`   //优惠退款金额
	Currency         string `json:"currency"`          //退款币种
}
