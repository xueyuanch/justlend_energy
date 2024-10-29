package justlend

const (
	GetRentalRateABI      = "0x3193fada"                         // getRentalRate(uint256,uint256)
	RentResourceABI       = "0xfd8527a1"                         // rentResource(address,uint256,uint256)
	ReturnResourceABI     = "0xaf6f4896"                         // returnResource(address,uint256,uint256)
	LiquidateThresholdABI = "0xfdcb648c"                         //	liquidateThreshold()
	MinFeeABI             = "0x24ec7590"                         // minFee()
	FeeRatioABI           = "0x41744dd4"                         // feeRatio()
	JustLendContract      = "TU2MJ5Veik1LRAgjeSzEdvmDYx7mefJZvd" // JustLend DAO: Energy Rental
	TokenDefaultPrecision = 1000000000000000000
)

type Service interface {
	RentResourceService
	ReturnResourceService
	FeeRatioService
}
