// Package store encapsulates all database interaction.
package store

import (
	"bytes"
	"database/sql/driver"
	"encoding/csv"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	sqlDialect         = "postgres"
	FlashPoolDetailKey = "flashpooldetail"
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
	TotalSupply    uint64
	TotalBorrow    uint64
	TotalInsurance uint64
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
	TotalSupply    uint64
	TotalBorrow    uint64
	TotalInsurance uint64
}

func (client Client) LoadLatestFlashPoolMarket(name string) (FlashPoolMarket, error) {
	var flashPoolMarket FlashPoolMarket
	err := client.db.Where(FlashPoolMarket{Name: name}).Last(&flashPoolMarket).Error
	return flashPoolMarket, err
}

func (client Client) SaveFlashPoolMarket(flashPoolMarket *FlashPoolMarket) error {
	return client.db.Create(flashPoolMarket).Error
}
