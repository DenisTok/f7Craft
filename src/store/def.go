package store

import (
	"context"
	"github.com/DenisTok/f7Craft/src/config"
	"sync"
	"time"

	"github.com/dgraph-io/badger/v2"
	"github.com/dgraph-io/badger/v2/options"
	"github.com/rs/zerolog/log"
)

type DefStore struct {
	db *badger.DB
	sync.Mutex
}

func NewDefStore(ctx context.Context, name string) (*DefStore, error) {
	opt := badger.DefaultOptions(config.DBDir + name).
		//WithValueLogLoadingMode(options.FileIO).
		WithValueLogFileSize((1<<20 - 1) * 10).
		//WithTableLoadingMode(options.FileIO).
		WithCompression(options.ZSTD).
		WithZSTDCompressionLevel(1)

	db, err := badger.Open(opt)
	if err != nil {
		return nil, err
	}

	go func(name string) {
		select {
		case <-ctx.Done():
			log.Info().Msg("Close db " + name)
			err := db.Close()
			if err != nil {
				log.Warn().Err(err).Send()
			}
		}
		log.Info().Msg("Close done " + name)
	}(name)

	def := DefStore{db: db}

	def.runGC(ctx, time.Minute*5)

	return &def, err
}

func (kv *DefStore) IsClose() bool {
	kv.Lock()
	defer kv.Unlock()
	return kv.db.IsClosed()
}

func (kv *DefStore) Update(f func(txn *badger.Txn) error) error {
	kv.Lock()
	defer kv.Unlock()
	err := kv.db.Update(f)
	if err != nil {
		return err
	}
	return err
}

func (kv *DefStore) Write(key, value []byte) error {
	kv.Lock()
	defer kv.Unlock()
	err := kv.db.Update(func(txn *badger.Txn) error {
		return txn.SetEntry(badger.NewEntry(key, value))
	})
	if err != nil {
		return err
	}
	return err
}

func (kv *DefStore) Get(key []byte, dsc func(v []byte) error) error {
	kv.Lock()
	defer kv.Unlock()
	err := kv.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		err = item.Value(dsc)
		if err != nil {
			return err
		}

		return err
	})
	if err != nil {
		return err
	}
	return err
}

func (kv *DefStore) Exist(key []byte) error {
	kv.Lock()
	defer kv.Unlock()
	err := kv.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(key)
		if err != nil {
			return err
		}

		return err
	})
	if err != nil {
		return err
	}

	return err
}

func (kv *DefStore) lastKey() (key []byte, err error) {
	kv.Lock()
	defer kv.Unlock()
	err = kv.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Reverse = true

		it := txn.NewIterator(opts)

		defer it.Close()

		for it.Rewind(); it.Valid(); {
			item := it.Item()
			key = item.KeyCopy(nil)
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return key, err
}

func (kv *DefStore) Delete(keys [][]byte) error {
	kv.Lock()
	defer kv.Unlock()
	err := kv.db.Update(func(txn *badger.Txn) error {
		for _, key := range keys {
			err := txn.Delete(key)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return err
}

func (kv *DefStore) runGC(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			log.Info().Msg("start gc")
			select {
			case <-ticker.C:
				start := time.Now()

				for {
					err := kv.db.RunValueLogGC(config.DiscardRatio)
					if err != nil {
						if err != badger.ErrNoRewrite {
							log.Warn().Err(err).Send()
							break
						}
						break
					}
				}

				log.Info().Float64("garbage collection finished", float64(time.Since(start))/float64(time.Second))
			case <-ctx.Done():
				log.Info().Msg("stop gc")
				return
			}
		}
	}()
}
