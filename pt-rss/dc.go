package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

type downloadCounter struct {
	data     map[string]int
	lock     sync.Mutex
	filename string
}

func (dc *downloadCounter) GetCount(key string) int {
	if count, ok := dc.data[key]; ok {
		return count
	}
	return 0
}

func (dc *downloadCounter) SetCount(key string, value int) {
	dc.lock.Lock()
	defer dc.lock.Unlock()
	dc.data[key] = value
}

func (dc *downloadCounter) IncrCount(key string) {
	dc.lock.Lock()
	defer dc.lock.Unlock()
	count := dc.GetCount(key)
	count++
	dc.SetCount(key, count)
}

func (dc *downloadCounter) Init(filename string) error {
	dc.lock.Lock()
	defer dc.lock.Unlock()
	dc.filename = filename
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "db file %s not exist, will be created\n", filename)
		dc.data = make(map[string]int, 1)
		return nil
	}
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err = decoder.Decode(&dc.data)
	if err != nil {
		return err
	}
	return nil
}

func (dc *downloadCounter) Dump() error {
	dc.lock.Lock()
	defer dc.lock.Unlock()
	var data bytes.Buffer
	encoder := gob.NewEncoder(&data)
	encoder.Encode(&dc.data)
	err := ioutil.WriteFile(dc.filename, data.Bytes(), 0644)
	if err != nil {
		return err
	}
	return nil
}
