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
	UserAddress      string `gorm:"primary_key"`
	AssetName        string `gorm:"primary_key"`
	Icon             string
	SupplyBalance    string
	BorrowBalance    string
	InsuranceBalance string
	IfCollateral     bool
}

func (client Client) LoadUserBalance(userAddress string) ([]UserAssetBalance, error) {
	userBalance := make([]UserAssetBalance, 0)
	err := client.db.Where("user_address = ?", userAddress).Find(&userBalance).Error
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

type AssetApy struct {
	AssetName        string `gorm:"primary_key"`
	SupplyApy        string
	BorrowApy        string
	InsuranceApy     string
	CollateralFactor string
}

func (client Client) LoadAssetApy(assetName string) (AssetApy, error) {
	var assetApy AssetApy
	err := client.db.Where(AssetApy{AssetName: assetName}).Last(&assetApy).Error
	return assetApy, err
}

func (client Client) SaveAssetApy(assetApy *AssetApy) error {
	return client.db.Save(assetApy).Error
}
