package server

import (
	"fmt"
	"sync"
)

var ErrOffsetNotFound = fmt.Errorf("offset not found")

type Log struct {
	mu	sync.Mutex		// 同時にアクセスする可能性があるため排他的ロックを提供するライブラリを用いる
	records []Record
}

type Record struct {
	Value []byte `json:"value`		
	Offset uint64 `json:offset`		// ログが格納される位置を示す整数値
}

// 新しいLogオブジェクトを生成し，ポインタを返す
func NewLog() *Log {
	return &Log{}
}

// ログのレコードの値を追加する
func (c *Log) Append(record Record) (uint64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()    // 
	record.Offset = uint64(len(c.records))
	c.records = append(c.records, record)
	return record.Offset, nil
}

// Log構造体にあるoffsetを読み取る
func (c *Log) Read(offset uint64) (Record, error) {
	c.mu.Lock()
	defer c.mu.Unlock()    
	if offset >= uint64(len(c.records)) {
		return Record{}, ErrOffsetNotFound
	}
	return c.records[offset], nil
}


