package model

import "gorm.io/gorm"

// Liquidation 清仓
type Liquidation struct {
	gorm.Model
	StockName        string  `gorm:"column:stock_name;type:string;comment:股票名称" json:"stock_name"`
	ProfitAndLoss    string  `gorm:"column:profit_and_loss;type:string;comment:盈亏" json:"profit_and_loss"`
	AmountOfMoney    float32 `gorm:"column:amount_of_money;type:decimal(10,5);comment:盈亏金额" json:"amount_of_money"`
	EarningRate      float32 `gorm:"column:earning_rate;type:decimal(10,5);comment:收益率" json:"earning_rate"`
	HoldStockDays    uint32  `gorm:"column:hold_stock_days;type:uint;comment:持股天数" json:"hold_stock_days"`
	TransactionTaxes float32 `gorm:"column:transaction_taxes;type:decimal(10,5);comment:交易税费" json:"transaction_taxes"`
	BuyingDay        string  `gorm:"column:buying_day;type:string;comment:建仓日期" json:"buying_day"`
	SellingDay       string  `gorm:"column:selling_day;type:string;comment:清仓日期" json:"selling_day"`
}

// TradingRecords 历史成交
type TradingRecords struct {
	gorm.Model
	TradingDay          string  `gorm:"column:trading_day;type:string;comment:成交日期"`
	TradingTime         string  `gorm:"column:trading_time;type:string;comment:成交时间"`
	SecuritiesCode      uint32  `gorm:"column:securities_code;type:uint;comment:证券代码"`
	SecuritiesName      string  `gorm:"column:securities_name;type:string;comment:证券名称"`
	Operate             string  `gorm:"column:operate;type:string;comment:操作"`
	TradingVolume       uint32  `gorm:"column:trading_volume;type:uint;comment:成交数量"`
	TradingAveragePrice float32 `gorm:"column:trading_average_price;type:decimal(10,5);comment:成交均价"`
	Turnover            float32 `gorm:"column:turnover;type:decimal(10,5);comment:成交金额"`
	ContractNumber      uint32  `gorm:"column:contract_number;type:uint;comment:合同编号"`
	TradingNumber       uint64  `gorm:"column:trading_number;type:bigint;comment:成交编号"`
	Commission          float32 `gorm:"column:commission;type:decimal(10,5);comment:手续费"`
	StampDuty           float32 `gorm:"column:stamp_duty;type:decimal(10,5);comment:印花税"`
	OtherExpenses       float32 `gorm:"column:other_expenses;type:decimal(10,5);comment:其他杂费"`
	Note                string  `gorm:"column:note;type:string;comment:备注"`
	ShareholderAccounts string  `gorm:"column:shareholder_accounts;type:string;comment:股东帐户"`
}
