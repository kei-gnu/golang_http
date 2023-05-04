package server

import (
	"reflect"
	"testing"
)

func TestLogAppend(t *testing.T) {
	log := NewLog()
	record := Record{Value: []byte("testlog")}    // recordの作成
	offset, err := log.Append(record)			// logオブジェクトに対してappend関数を実行

	if err != nil {
		t.Errorf("Append returned unexpected error: %v", err)
	}

	if offset != 0 {
		t.Errorf("Append returned unexpected offset: got %v, want 0", offset)
	}

	if !reflect.DeepEqual(log.records[0], record) {    // !reflectとは？
		t.Errorf("Append did not store corret record: got %v, want %v", log.records[0], record)
	}
}

func TestLogRead(t *testing.T) {
	log := NewLog()
	record := Record{Value: []byte("testreadlog")}
	log.records = append(log.records, record)

	got, err := log.Read(0)

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