// Package store encapsulates all database interaction.
package store

import (
	"bytes"
	"database/sql/driver"
	"encoding/csv"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/siovanus/wingServer/http/common"
	"github.com/siovanus/wingServer/store/migrations"
)

const (
	sqlDialect = "postgres"
)

// SQLStringArray is a string array stored in the database as comma separated values.
type SQLStringArray []string

// Scan implements the sql Scanner interface.
func (arr *SQLStringArray) Scan(src interface{}) error {
	if src == nil {
		*arr = nil
	}
	v, err := driver.String.ConvertValue(src)
	if err != nil {
		return fmt.Errorf("failed to scan StringArray")
	}
	str, ok := v.(string)
	if !ok {
		return nil
	}

	buf := bytes.NewBufferString(str)
	r := csv.NewReader(buf)
	ret, err := r.Read()
	if err != nil {
		return fmt.Errorf("badly formatted csv string array: %s", err)
	}
	*arr = ret
	return nil
}

// Value implements the driver Valuer interface.
func (arr SQLStringArray) Value() (driver.Value, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	err := w.Write(arr)
	if err != nil {
		return nil, fmt.Errorf("csv encoding of string array: %s", err)
	}
	w.Flush()
	return buf.String(), nil
}

// Client holds a connection to the database.
type Client struct {
	db *gorm.DB
}

// ConnectToDB attempts to connect to the database URI provided,
// and returns a new Client instance if successful.
func ConnectToDb(uri string) (*Client, error) {
	db, err := gorm.Open(sqlDialect, uri)
	if err != nil {
		return nil, fmt.Errorf("unable to open %s for gorm DB: %+v", uri, err)
	}
	if err = migrations.Migrate(db); err != nil {
		return nil, fmt.Errorf("newDBStore#Migrate: %s", err)
	}
	store := &Client{
		db: db.Set("gorm:auto_preload", true),
	}
	return store, nil
}

// Close will close the connection to the database.
func (client Client) Close() error {
	return client.db.Close()
}

type FlashPoolDetail struct {
	Timestamp      uint64 `gorm:"primary_key"`
	TotalSupply    string
	TotalBorrow    string
	TotalInsurance string
}

func (client Client) LoadLatestFlashPoolDetail() (FlashPoolDetail, error) {
	var flashPoolDetail FlashPoolDetail
	err := client.db.Last(&flashPoolDetail).Error
	return flashPoolDetail, err
}

func (client Client) SaveFlashPoolDetail(flashPoolDetail *FlashPoolDetail) error {
	return client.db.Create(flashPoolDetail).Error
}

type FlashPoolMarket struct {
	ID             uint64
	Name           string
	Timestamp      uint64
	TotalSupply    string
	TotalBorrow    string
	TotalInsurance string
}

func (client Client) LoadLatestFlashPoolMarket(name string) (FlashPoolMarket, error) {
	var flashPoolMarket FlashPoolMarket
	err := client.db.Where(FlashPoolMarket{Name: name}).Last(&flashPoolMarket).Error
	return flashPoolMarket, err
}

func (client Client) SaveFlashPoolMarket(flashPoolMarket *FlashPoolMarket) error {
	return client.db.Create(flashPoolMarket).Error
}

type Price struct {
	Name  string `gorm:"primary_key"`
	Price string
}

func (client Client) LoadPrice(name string) (Price, error) {
	var price Price
	err := client.db.Where(Price{Name: name}).Last(&price).Error
	return price, err
}

func (client Client) SavePrice(Price *Price) error {
	return client.db.Save(Price).Error
}

type TrackHeight struct {
	Name   string `gorm:"primary_key"`
	Height uint32
}

func (client Client) LoadTrackHeight() (uint32, error) {
	var trackHeight TrackHeight
	err := client.db.Where(TrackHeight{Name: "TrackHeight"}).Last(&trackHeight).Error
	return trackHeight.Height, err
}

func (client Client) SaveTrackHeight(height uint32) error {
	trackHeight := &TrackHeight{
		Name:   "TrackHeight",
		Height: height,
	}
	return client.db.Save(trackHeight).Error
}

type UserAssetBalance struct {
	UserAddress  string `gorm:"primary_key"`
	AssetName    string `gorm:"primary_key"`
	AssetAddress string
	Icon         string
	FToken       string
	BorrowAmount string
	BorrowIndex  string
	Itoken       string
	IfCollateral bool
}

func (client Client) LoadUserBalance(userAddress string) ([]UserAssetBalance, error) {
	userBalance := make([]UserAssetBalance, 0)
	err := client.db.Where("user_address = ?", userAddress).Find(&userBalance).Error
	if err != nil {
		return userBalance, err
	}
	return userBalance, err
}

func (client Client) LoadBorrowUsers() ([]UserAssetBalance, error) {
	userBalance := make([]UserAssetBalance, 0)
	err := client.db.Select("user_address").Where("borrow_amount <> ?", "0").Find(&userBalance).Error
	if err != nil {
		return userBalance, err
	}
	return userBalance, err
}

func (client Client) SaveUserAssetBalance(input *UserAssetBalance) error {
	return client.db.Save(input).Error
}

func (client Client) LoadFlashMarket(name string) (common.Market, error) {
	var market common.Market
	err := client.db.Where(common.Market{Name: name}).Last(&market).Error
	return market, err
}

func (client Client) SaveFlashMarket(market *common.Market) error {
	return client.db.Save(market).Error
}

func (client Client) UpdateFlashMarketBorrowIndex(name, borrowIndex string) error {
	return client.db.Model(&common.Market{Name: name}).Update("borrow_index", borrowIndex).Error
}

func (client Client) LoadWingApys() ([]common.WingApy, error) {
	wingApys := make([]common.WingApy, 0)
	err := client.db.Find(&wingApys).Error
	return wingApys, err
}

func (client Client) LoadWingApy(assetName string) (common.WingApy, error) {
	var wingApy common.WingApy
	err := client.db.Where(common.WingApy{AssetName: assetName}).Last(&wingApy).Error
	return wingApy, err
}

func (client Client) SaveWingApy(wingApy *common.WingApy) error {
	return client.db.Save(wingApy).Error
}

type IFInfo struct {
	Name  string `gorm:"primary_key"`
	Total string
	Cap   string
}

func (client Client) LoadIFInfo() (IFInfo, error) {
	var ifInfo IFInfo
	err := client.db.Where(IFInfo{Name: "IFInfo"}).Last(&ifInfo).Error
	return ifInfo, err
}

func (client Client) SaveIFInfo(ifInfo *IFInfo) error {
	ifInfo.Name = "IFInfo"
	return client.db.Save(ifInfo).Error
}

type IFMarketInfo struct {
	Name             string `gorm:"primary_key"`
	TotalCash        string
	TotalDebt        string
	InterestIndex    string
	TotalInsurance   string
	InterestRate     uint64
	CollateralFactor uint64
}

func (client Client) LoadIFMarketInfo(name string) (IFMarketInfo, error) {
	var ifMarketInfo IFMarketInfo
	err := client.db.Where(IFMarketInfo{Name: name}).Last(&ifMarketInfo).Error
	return ifMarketInfo, err
}

func (client Client) SaveIFMarketInfo(ifMarketInfo *IFMarketInfo) error {
	return client.db.Save(ifMarketInfo).Error
}

type IfPoolHistory struct {
	ID        uint64 `gorm:"primary_key"`
	Address   string
	Token     string
	Operation string
	Amount    string
	Timestamp uint64
	TxHash    string
}

func (client Client) SaveIFHistory(history *IfPoolHistory) error {
	return client.db.Save(history).Error
}

func (client Client) LoadIFHistory(asset, operation string, start, end, pageNo, pageSize uint64) ([]IfPoolHistory, error) {
	startPage := pageSize * (pageNo - 1)
	if startPage < 0 {
		startPage = 0
	}
	IfPoolHistory := make([]IfPoolHistory, 0)
	db := client.db
	if asset != "" {
		db.Where("token = ?", asset)
	}
	if operation != "" {
		db.Where("operation = ?", operation)
	}
	if start > 0 {
		db.Where("timestamp >= ?", start)
	}
	if end > 0 {
		db.Where("timestamp <= ?", end)
	}
	db.Order("timestamp desc")
	//count := db.Count(&count)
	db.Offset(startPage)
	db.Limit(pageSize)
	err := db.Find(&IfPoolHistory).Error
	if err != nil {
		return IfPoolHistory, err
	}
	return IfPoolHistory, err
}
