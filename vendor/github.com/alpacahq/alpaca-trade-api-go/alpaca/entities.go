package alpaca

import (
	"time"

	"github.com/shopspring/decimal"
)

type Account struct {
	ID                    string          `json:"id"`
	AccountNumber         string          `json:"account_number"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`
	DeletedAt             *time.Time      `json:"deleted_at"`
	Status                string          `json:"status"`
	Currency              string          `json:"currency"`
	Cash                  decimal.Decimal `json:"cash"`
	CashWithdrawable      decimal.Decimal `json:"cash_withdrawable"`
	TradingBlocked        bool            `json:"trading_blocked"`
	TransfersBlocked      bool            `json:"transfers_blocked"`
	AccountBlocked        bool            `json:"account_blocked"`
	ShortingEnabled       bool            `json:"shorting_enabled"`
	BuyingPower           decimal.Decimal `json:"buying_power"`
	PatternDayTrader      bool            `json:"pattern_day_trader"`
	DaytradeCount         int64           `json:"daytrade_count"`
	DaytradingBuyingPower decimal.Decimal `json:"daytrading_buying_power"`
	RegTBuyingPower       decimal.Decimal `json:"regt_buying_power"`
	Equity                decimal.Decimal `json:"equity"`
	LastEquity            decimal.Decimal `json:"last_equity"`
	Multiplier            string          `json:"multiplier"`
	InitialMargin         decimal.Decimal `json:"initial_margin"`
	MaintenanceMargin     decimal.Decimal `json:"maintenance_margin"`
	LastMaintenanceMargin decimal.Decimal `json:"last_maintenance_margin"`
	LongMarketValue       decimal.Decimal `json:"long_market_value"`
	ShortMarketValue      decimal.Decimal `json:"short_market_value"`
	PortfolioValue        decimal.Decimal `json:"portfolio_value"`
}

type Order struct {
	ID             string           `json:"id"`
	ClientOrderID  string           `json:"client_order_id"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	SubmittedAt    time.Time        `json:"submitted_at"`
	FilledAt       *time.Time       `json:"filled_at"`
	ExpiredAt      *time.Time       `json:"expired_at"`
	CanceledAt     *time.Time       `json:"canceled_at"`
	FailedAt       *time.Time       `json:"failed_at"`
	AssetID        string           `json:"asset_id"`
	Symbol         string           `json:"symbol"`
	Exchange       string           `json:"exchange"`
	Class          string           `json:"asset_class"`
	Qty            decimal.Decimal  `json:"qty"`
	FilledQty      decimal.Decimal  `json:"filled_qty"`
	Type           OrderType        `json:"order_type"`
	Side           Side             `json:"side"`
	TimeInForce    TimeInForce      `json:"time_in_force"`
	LimitPrice     *decimal.Decimal `json:"limit_price"`
	FilledAvgPrice *decimal.Decimal `json:"filled_avg_price"`
	StopPrice      *decimal.Decimal `json:"stop_price"`
	Status         string           `json:"status"`
}

type Position struct {
	AssetID        string          `json:"asset_id"`
	Symbol         string          `json:"symbol"`
	Exchange       string          `json:"exchange"`
	Class          string          `json:"asset_class"`
	AccountID      string          `json:"account_id"`
	EntryPrice     decimal.Decimal `json:"avg_entry_price"`
	Qty            decimal.Decimal `json:"qty"`
	Side           string          `json:"side"`
	MarketValue    decimal.Decimal `json:"market_value"`
	CostBasis      decimal.Decimal `json:"cost_basis"`
	UnrealizedPL   decimal.Decimal `json:"unrealized_pl"`
	UnrealizedPLPC decimal.Decimal `json:"unrealized_plpc"`
	CurrentPrice   decimal.Decimal `json:"current_price"`
	LastdayPrice   decimal.Decimal `json:"lastday_price"`
	ChangeToday    decimal.Decimal `json:"change_today"`
}

type Asset struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Exchange string `json:"exchange"`
	Class    string `json:"asset_class"`
	Symbol   string `json:"symbol"`
	Status   string `json:"status"`
	Tradable bool   `json:"tradable"`
}

type Fundamental struct {
	AssetID           string          `json:"asset_id"`
	Symbol            string          `json:"symbol"`
	FullName          string          `json:"full_name"`
	IndustryName      string          `json:"industry_name"`
	IndustryGroup     string          `json:"industry_group"`
	Sector            string          `json:"sector"`
	PERatio           float32         `json:"pe_ratio"`
	PEGRatio          float32         `json:"peg_ratio"`
	Beta              float32         `json:"beta"`
	EPS               float32         `json:"eps"`
	MarketCap         int64           `json:"market_cap"`
	SharesOutstanding int64           `json:"shares_outstanding"`
	AvgVol            int64           `json:"avg_vol"`
	DivRate           float32         `json:"div_rate"`
	ROE               float32         `json:"roe"`
	ROA               float32         `json:"roa"`
	PS                float32         `json:"ps"`
	PC                float32         `json:"pc"`
	GrossMargin       float32         `json:"gross_margin"`
	FiftyTwoWeekHigh  decimal.Decimal `json:"fifty_two_week_high"`
	FiftyTwoWeekLow   decimal.Decimal `json:"fifty_two_week_low"`
	ShortDescription  string          `json:"short_description"`
	LongDescription   string          `json:"long_description"`
}

type Bar struct {
	Time   int64   `json:"t"`
	Open   float32 `json:"o"`
	High   float32 `json:"h"`
	Low    float32 `json:"l"`
	Close  float32 `json:"c"`
	Volume int32   `json:"v"`
}

type ListBarParams struct {
	Timeframe string     `url:"timeframe,omitempty"`
	StartDt   *time.Time `url:"start_dt,omitempty"`
	EndDt     *time.Time `url:"end_dt,omitempty"`
	Limit     *int       `url:"limit,omitempty"`
}

type Quote struct {
	BidTimestamp  time.Time `json:"bid_timestamp"`
	Bid           float32   `json:"bid"`
	AskTimestamp  time.Time `json:"ask_timestamp"`
	Ask           float32   `json:"ask"`
	LastTimestamp time.Time `json:"last_timestamp"`
	Last          float32   `json:"last"`
	AssetID       string    `json:"asset_id"`
	Symbol        string    `json:"symbol"`
	Class         string    `json:"asset_class"`
}

type CalendarDay struct {
	Date  string `json:"date"`
	Open  string `json:"open"`
	Close string `json:"close"`
}

type Clock struct {
	Timestamp time.Time `json:"timestamp"`
	IsOpen    bool      `json:"is_open"`
	NextOpen  time.Time `json:"next_open"`
	NextClose time.Time `json:"next_close"`
}

type AccountConfigurations struct {
	DtbpCheck            DtbpCheck         `json:"dtbp_check"`
	NoShorting           bool              `json:"no_shorting"`
	TradeConfirmEmail    TradeConfirmEmail `json:"trade_confirm_email"`
	TradeSuspendedByUser bool              `json:"trade_suspended_by_user"`
}

type AccountActvity struct {
	ID              string          `json:"id"`
	ActivityType    string          `json:"activity_type"`
	TransactionTime time.Time       `json:"transaction_time"`
	Type            string          `json:"type"`
	Price           decimal.Decimal `json:"price"`
	Qty             decimal.Decimal `json:"qty"`
	Side            string          `json:"side"`
	Symbol          string          `json:"symbol"`
	LeavesQty       decimal.Decimal `json:"leaves_qty"`
	CumQty          decimal.Decimal `json:"cum_qty"`
	Date            time.Time       `json:"date"`
	NetAmount       decimal.Decimal `json:"net_amount"`
	Description     string          `json:"description"`
	PerShareAmount  decimal.Decimal `json:"per_share_amount"`
}

type PlaceOrderRequest struct {
	AccountID     string           `json:"-"`
	AssetKey      *string          `json:"symbol"`
	Qty           decimal.Decimal  `json:"qty"`
	Side          Side             `json:"side"`
	Type          OrderType        `json:"type"`
	TimeInForce   TimeInForce      `json:"time_in_force"`
	LimitPrice    *decimal.Decimal `json:"limit_price"`
	StopPrice     *decimal.Decimal `json:"stop_price"`
	ClientOrderID string           `json:"client_order_id"`
}

type ReplaceOrderRequest struct {
	Qty           *decimal.Decimal `json:"qty"`
	LimitPrice    *decimal.Decimal `json:"limit_price"`
	StopPrice     *decimal.Decimal `json:"stop_price"`
	TimeInForce   TimeInForce      `json:"time_in_force"`
	ClientOrderID string           `json:"client_order_id"`
}

type AccountConfigurationsRequest struct {
	DtbpCheck            *string `json:"dtbp_check"`
	NoShorting           *bool   `json:"no_shorting"`
	TradeConfirmEmail    *string `json:"trade_confirm_email"`
	TradeSuspendedByUser *bool   `json:"trade_suspended_by_user"`
}

type AccountActivitiesRequest struct {
	ActivityTypes *[]string  `json:"activity_types"`
	Date          *time.Time `json:"date"`
	Until         *time.Time `json:"until"`
	After         *time.Time `json:"after"`
	Direction     *string    `json:"direction"`
	PageSize      *int       `json:"page_size"`
}

type Side string

const (
	Buy  Side = "buy"
	Sell Side = "sell"
)

type OrderType string

const (
	Market        OrderType = "market"
	Limit         OrderType = "limit"
	Stop          OrderType = "stop"
	StopLimit     OrderType = "stop_limit"
	MarketOnClose OrderType = "market_on_close"
	LimitOnClose  OrderType = "limit_on_close"
)

type TimeInForce string

const (
	Day TimeInForce = "day"
	GTC TimeInForce = "gtc"
	OPG TimeInForce = "opg"
	IOC TimeInForce = "ioc"
	FOK TimeInForce = "fok"
	GTX TimeInForce = "gtx"
	GTD TimeInForce = "gtd"
	CLS TimeInForce = "cls"
)

type DtbpCheck string

const (
	Entry DtbpCheck = "entry"
	Exit  DtbpCheck = "exit"
	Both  DtbpCheck = "both"
)

type TradeConfirmEmail string

const (
	None TradeConfirmEmail = "none"
	All  TradeConfirmEmail = "all"
)

// stream

// ClientMsg is the standard message sent by clients of the stream interface
type ClientMsg struct {
	Action string      `json:"action" msgpack:"action"`
	Data   interface{} `json:"data" msgpack:"data"`
}

// ServerMsg is the standard message sent by the server to update clients
// of the stream interface
type ServerMsg struct {
	Stream string      `json:"stream" msgpack:"stream"`
	Data   interface{} `json:"data"`
}

type TradeUpdate struct {
	Event string `json:"event"`
	Order Order  `json:"order"`
}
