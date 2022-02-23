// Copyright (c) 2022, Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found s the LICENSE file.

package main

import (
	"encoding/binary"

	"github.com/ethersphere/bee/pkg/postage"
	"github.com/ethersphere/bee/pkg/shed"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type localstore struct {
	shed               *shed.DB
	retrievalDataIndex shed.Index
}

func newLocalstore(path string) (s *localstore, err error) {
	ldb, err := leveldb.OpenFile(path, &opt.Options{
		OpenFilesCacheCapacity: 200,
		BlockCacheCapacity:     32 * 1024 * 1024,
		WriteBuffer:            32 * 1024 * 1024,
		DisableSeeksCompaction: false,
		ErrorIfMissing:         true,
		ReadOnly:               true,
	})
	if err != nil {
		return nil, err
	}

	s = new(localstore)

	s.shed, err = shed.NewDBWrap(ldb)
	if err != nil {
		return nil, err
	}

	// Based on github.com/ethersphere/bee v1.4.3

	headerSize := 16 + postage.StampSize
	s.retrievalDataIndex, err = s.shed.NewIndex("Address->StoreTimestamp|BinID|BatchID|BatchIndex|Sig|Data", shed.IndexFuncs{
		EncodeKey: func(fields shed.Item) (key []byte, err error) {
			return fields.Address, nil
		},
		DecodeKey: func(key []byte) (e shed.Item, err error) {
			e.Address = key
			return e, nil
		},
		EncodeValue: func(fields shed.Item) (value []byte, err error) {
			b := make([]byte, headerSize)
			binary.BigEndian.PutUint64(b[:8], fields.BinID)
			binary.BigEndian.PutUint64(b[8:16], uint64(fields.StoreTimestamp))
			stamp, err := postage.NewStamp(fields.BatchID, fields.Index, fields.Timestamp, fields.Sig).MarshalBinary()
			if err != nil {
				return nil, err
			}
			copy(b[16:], stamp)
			value = append(b, fields.Data...)
			return value, nil
		},
		DecodeValue: func(keyItem shed.Item, value []byte) (e shed.Item, err error) {
			// do not calculate unneeded fields for this particular purpose
			//
			// e.StoreTimestamp = int64(binary.BigEndian.Uint64(value[8:16]))
			// e.BinID = binary.BigEndian.Uint64(value[:8])
			// stamp := new(postage.Stamp)
			// if err = stamp.UnmarshalBinary(value[16:headerSize]); err != nil {
			// 	return e, err
			// }
			// e.BatchID = stamp.BatchID()
			// e.Index = stamp.Index()
			// e.Timestamp = stamp.Timestamp()
			// e.Sig = stamp.Sig()
			e.Data = value[headerSize:]
			return e, nil
		},
	})
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *localstore) IterateChunkData(fn func(data []byte) (stop bool, err error)) error {
	return s.retrievalDataIndex.Iterate(func(item shed.Item) (stop bool, err error) {
		return fn(item.Data)
	}, nil)
}

func (s *localstore) Close() error {
	return s.shed.Close()
}
