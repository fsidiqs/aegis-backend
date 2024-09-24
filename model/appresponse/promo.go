package appresponse

const (
	PROMO_DEVICE_ILLEGIBLE_MSG           ErrMsgResp = "illigible"
	PROMO_INVALID_REFERRAL_CODE_MSG      ErrMsgResp = "invalid referral code"
	PROMO_REFERRAL_CODE_ALREADY_USED_MSG ErrMsgResp = "referral code already used"

	PROMO_DEVICE_ELIGIBLE_MSG SuccessMsgResp = "eligible"
)

const (
	PROMO_DEVICE_ILLEGIBLE_TYPE           ErrResponseType = "ILLEGIBLE"
	PROMO_INVALID_REFERRAL_CODE_TYPE      ErrResponseType = "INVALID_REFERRAL_CODE"
	PROMO_REFERRAL_CODE_ALREADY_USED_TYPE ErrResponseType = "REFERRAL_CODE_ALREADY_USED"

	PROMO_DEVICE_ELIGIBLE_TYPE SuccessResponseType = "ELIGIBLE"
)
