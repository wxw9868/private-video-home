package service

import (
	"github.com/spf13/cast"
	"github.com/wxw9868/util/pagination"
	"github.com/wxw9868/video/model"
	"github.com/xuri/excelize/v2"
)

type StockService struct{}

func (s *StockService) ImportTradingRecords(filepath string) error {
	f, err := excelize.OpenFile(filepath)
	if err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	// 获取工作表列表
	Sheet1 := f.GetSheetList()[0]
	// 获取 Sheet1 上所有单元格
	rows, err := f.GetRows(Sheet1)
	if err != nil {
		return err
	}

	tradingRecords := make([]model.TradingRecords, len(rows)-1)
	for i, row := range rows {
		if i > 0 {
			tradingRecords[i-1] = model.TradingRecords{
				TradingDay:          row[0],
				TradingTime:         row[1],
				SecuritiesCode:      cast.ToUint32(row[2]),
				SecuritiesName:      row[3],
				Operate:             row[4],
				TradingVolume:       cast.ToUint32(row[5]),
				TradingAveragePrice: cast.ToFloat32(row[6]),
				Turnover:            cast.ToFloat32(row[7]),
				ContractNumber:      cast.ToUint32(row[8]),
				TradingNumber:       cast.ToUint64(row[9]),
				Commission:          cast.ToFloat32(row[10]),
				StampDuty:           cast.ToFloat32(row[11]),
				OtherExpenses:       cast.ToFloat32(row[12]),
				Note:                row[13],
				ShareholderAccounts: row[14],
			}
		}
	}

	if err = db.CreateInBatches(tradingRecords, 50).Error; err != nil {
		return err
	}
	return nil
}

type Liquidation struct {
	ID               uint    `json:"id"`
	StockName        string  `gorm:"column:stock_name;type:string;comment:股票名称" json:"stockName"`
	ProfitAndLoss    string  `gorm:"column:profit_and_loss;type:string;comment:盈亏" json:"profitAndLoss"`
	AmountOfMoney    float32 `gorm:"column:amount_of_money;type:decimal(10,5);comment:盈亏金额" json:"amountOfMoney"`
	EarningRate      float32 `gorm:"column:earning_rate;type:decimal(10,5);comment:收益率" json:"earningRate"`
	HoldStockDays    uint32  `gorm:"column:hold_stock_days;type:uint;comment:持股天数" json:"holdStockDays"`
	TransactionTaxes float32 `gorm:"column:transaction_taxes;type:decimal(10,5);comment:交易税费" json:"transactionTaxes"`
	BuyingDay        string  `gorm:"column:buying_day;type:string;comment:建仓日期" json:"buyingDay"`
	SellingDay       string  `gorm:"column:selling_day;type:string;comment:清仓日期" json:"sellingDay"`
}

func (s *StockService) Liquidation(page, pageSize int) ([]Liquidation, error) {
	var count int64
	if err := db.Table("video_Liquidation").Count(&count).Error; err != nil {
		return nil, err
	}

	var list []Liquidation
	result := db.Model(&model.Liquidation{}).Scopes(Paginate(page, pageSize, int(count))).Find(&list)
	if result.Error != nil {
		return nil, result.Error
	}
	return list, nil
}

type TradingRecords struct {
	ID                  int     `json:"id"`
	TradingDay          string  `gorm:"column:trading_day;type:string;comment:成交日期" json:"trading_day"`
	TradingTime         string  `gorm:"column:trading_time;type:string;comment:成交时间" json:"trading_time"`
	SecuritiesCode      uint32  `gorm:"column:securities_code;type:uint;comment:证券代码" json:"securities_code"`
	SecuritiesName      string  `gorm:"column:securities_name;type:string;comment:证券名称" json:"securities_name"`
	Operate             string  `gorm:"column:operate;type:string;comment:操作" json:"operate"`
	TradingVolume       uint32  `gorm:"column:trading_volume;type:uint;comment:成交数量" json:"trading_volume"`
	TradingAveragePrice float32 `gorm:"column:trading_average_price;type:decimal(10,5);comment:成交均价" json:"trading_average_price"`
	Turnover            float32 `gorm:"column:turnover;type:decimal(10,5);comment:成交金额" json:"turnover"`
	ContractNumber      uint32  `gorm:"column:contract_number;type:uint;comment:合同编号" json:"contract_number"`
	TradingNumber       uint64  `gorm:"column:trading_number;type:bigint;comment:成交编号" json:"trading_number"`
	Commission          float32 `gorm:"column:commission;type:decimal(10,5);comment:手续费" json:"commission"`
	StampDuty           float32 `gorm:"column:stamp_duty;type:decimal(10,5);comment:印花税" json:"stamp_duty"`
	OtherExpenses       float32 `gorm:"column:other_expenses;type:decimal(10,5);comment:其他杂费" json:"other_expenses"`
	Note                string  `gorm:"column:note;type:string;comment:备注" json:"note"`
	ShareholderAccounts string  `gorm:"column:shareholder_accounts;type:string;comment:股东帐户" json:"shareholder_accounts"`
}

func (s *StockService) TradingRecords(page, pageSize int) ([]TradingRecords, error) {
	var count int64
	if err := db.Table("video_TradingRecords").Count(&count).Error; err != nil {
		return nil, err
	}

	var list []TradingRecords
	result := db.Model(&model.TradingRecords{}).Scopes(Paginate(page, pageSize, int(count))).Find(&list)
	if result.Error != nil {
		return nil, result.Error
	}
	return list, nil
}

func (s *StockService) Pagination(model interface{}, page, pageSize int) (*pagination.Paginator, error) {
	var total int64
	if err := db.Model(model).Count(&total).Error; err != nil {
		return nil, err
	}
	paginator := pagination.NewPaginator(int(total), pageSize)
	paginator.SetCurrentPage(page)
	return paginator, nil
}
