package main

import storage "github.com/parsaeisa/key_value_store/internal/storage"

const (
	Capacity = 10
)

func main() {
	k := storage.NewKVStore(Capacity)
	k.Set(storage.Record{
		Key:   "key1",
		Value: "value1",
	})

	res, err := k.Get("key2")
	if err != nil {
		println(err.Error())
	} else {
		println(res.Value)
	}

}
