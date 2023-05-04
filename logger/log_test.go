package logger

import (
	"reflect"
	"testing"
	"fmt"
)

func TestLogAppend(t *testing.T) {
	log := NewLog()
	record := Record{Value: []byte("abcde")}    // recordの作成
	offset, err := log.Append(record)			// logオブジェクトに対してappend関数を実行
	fmt.Printf("log: %v, record: %v, offset: %v, err: %v\n", log, record, offset, err)
	if err != nil {
		t.Errorf("Append returned unexpected error: %v", err)
	}

	if offset != 0 {
		t.Errorf("Append returned unexpected offset: got %v, want 0", offset)
	}

	if !reflect.DeepEqual(log.records[0], record) {    // !reflectとは？ log.records[0]とrecordの内容が等しい華道家を判定している(同じ型，変数の内容まで比較する)
		t.Errorf("Append did not store corret record: got %v, want %v", log.records[0], record)
	}
}

func TestLogRead(t *testing.T) {
	log := NewLog()
	record := Record{Value: []byte("testreadlog")}
	log.records = append(log.records, record)

	got, err := log.Read(0)
	fmt.Printf("log: %v, record: %v, got: %v, err: %v\n", log, record, got, err)
	if err != nil {
		t.Errorf("Read returned unexpected error: %v", err)
	}

	if !reflect.DeepEqual(got, record) {
		t.Errorf("Read returned unexpected record: got %v, want %v", got, record)
	}

	_, err = log.Read(1)
	if err == nil {
		t.Errorf("Read did not return expected error")
	}
}