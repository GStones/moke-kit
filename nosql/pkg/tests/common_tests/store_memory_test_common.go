package common_tests

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/gstones/platform/services/common/nosql/memory"
	"github.com/gstones/platform/services/common/nosql/memory/keys"
	"github.com/gstones/platform/services/common/nosql/memory/options"
)

type memTest struct {
	Index string
	Value interface{}
}

type listTest struct {
	MapId int32
}

var (
	ErrSetNXFailed = errors.New("SetNX error")
)

func MemoryStoreCommonTest(logger *zap.Logger, store memory.MemoryStore) (err error) {
	type testCase struct {
		key         string
		value       string
		expectedErr error
	}

	testCases := []testCase{
		{
			key:         "test:keys",
			value:       `{"a":"b"}`,
			expectedErr: nil,
		},
		{
			key:         "keys:to:test",
			value:       `{"b":"a"}`,
			expectedErr: nil,
		},
		{
			key:         "keys:to:test",
			value:       `{"a":{"b":{"c":true}}}`,
			expectedErr: nil,
		},
		{
			key:         "keys:to:test.stuff",
			value:       `{"Test": "Name", "Data": "TestData"}`,
			expectedErr: nil,
		},
	}

	// Set/Get
	for _, v := range testCases {
		err := store.Set(keys.NewKey(v.key), options.WithSource(v.value), options.WithTTL(time.Minute))
		if err != nil {
			return errors.Wrap(err, v.key)
		}
	}

	mt := &memTest{}
	mt.Index = "test_struct"
	mt.Value = 2
	key := keys.NewKey("test:struct")
	err = store.Set(key, options.WithSource(mt), options.WithTTL(time.Minute))
	if err != nil {
		return err
	}

	isExists, err := store.Exists(key)
	if err != nil {
		return err
	}
	logger.Info("Exists:", zap.Bool("isExists", isExists))

	err = store.Del(key)
	if err != nil {
		return err
	}
	logger.Info("Del")

	iKey := keys.NewKey("test:incur")
	scoreIn, err := store.Incur(iKey)
	if err != nil {
		return err
	}
	logger.Info("Incur:", zap.Any("scoreIn", scoreIn))

	pKey := keys.NewKey("zadd:test")
	member1 := keys.NewKey("zadd:test1")
	member2 := keys.NewKey("zadd:test2")
	member3 := keys.NewKey("zadd:test3")

	rankMap := make(map[string]interface{}, 3)
	rankMap[member1.String()] = float64(10)
	rankMap[member2.String()] = float64(20)
	rankMap[member3.String()] = float64(15)
	err = store.ZAdd(pKey, options.WithMultipleSource(rankMap))
	if err != nil {
		return err
	}

	//destZScan := make([]interface{}, 2)
	//err = store.ZRevRange(pKey, options.WithDestinationList(destZScan), options.WithInterval(0, 1))
	//if err != nil {
	//	return err
	//}
	//logger.Info("ZRevRange:", zap.Any("values", destZScan))

	score, err := store.ZScore(pKey, member2)
	if err != nil {
		return err
	}
	logger.Info("ZScore:", zap.Any(key.String(), score))

	rank, err := store.ZRevRank(pKey, member2)
	if err != nil {
		return err
	}
	logger.Info("zRevRank:", zap.Int64(member2.String(), rank))

	rank, err = store.ZRank(pKey, member2)
	if err != nil {
		return err
	}
	logger.Info("ZRank:", zap.Int64(member2.String(), rank))

	err = store.Set(key, options.WithSource(mt), options.WithTTL(time.Minute))
	if err != nil {
		return err
	}

	isOk, err := store.SetNX(key, options.WithSource(mt))
	if err != nil {
		return err
	}
	if isOk {
		return ErrSetNXFailed
	}

	dest := &memTest{}
	err = store.Get(key, options.WithDestination(dest))
	if err != nil {
		return
	}
	logger.Info("store get", zap.Any(key.String(), dest))

	// HSet/Get
	path := keys.NewKey("test:hset")
	for k := range testCases {
		key := fmt.Sprintf("test:struct.%d", k)
		mt := &memTest{}
		mt.Index = fmt.Sprintf("test_%d", k)
		mt.Value = k
		err = store.HSet(path, keys.NewKey(key), options.WithSource(mt))
		if err != nil {
			return
		}
	}
	destH := &memTest{}
	err = store.HGet(path, key, options.WithDestination(destH))
	if err != nil {
		logger.Error("store hget error", zap.Error(err))
		return err
	}
	logger.Info("store HGet", zap.Any(key.String(), destH))

	// HScan
	destScan := make(map[string]*memTest)
	id, err := store.Scan(path, options.WithDestination(&destScan), options.MatchRegex("test:*"))
	if err != nil {
		logger.Error("store scan error", zap.Error(err))
		return err
	}
	logger.Info("store scan", zap.Int("index", id), zap.Any("destScan", destScan))

	err = store.Remove(path)
	if err != nil {
		logger.Error("hash delete err:", zap.Error(err))
		return err
	}

	// HMSet/Get
	path = keys.NewKey("test:hmset:demo.demo1")
	const (
		NameKey   string = "profile:name"
		GenderKey string = "profile:key"
		AvatarKey string = "profile:avatar"
	)

	valueMap := make(map[string]interface{})
	valueMap[NameKey] = "1111"
	valueMap[GenderKey] = 0
	valueMap[AvatarKey] = 12
	err = store.HMSet(path, options.WithMultipleSource(valueMap))
	if err != nil {
		logger.Error("store hmset error", zap.Error(err))
		return err
	}

	valueMapResult := make(map[string]interface{})
	err = store.HGetAll(path, options.WithDestinationMap(valueMapResult))
	if err != nil {
		logger.Error("store hget all error", zap.Error(err))
		return err
	}
	logger.Debug("store hget all result", zap.Any("result", valueMapResult))

	keyList := keys.NewKeys(NameKey, GenderKey, AvatarKey)
	//valueList := make([]interface{}, len(keyList))
	//err = store.HMGet(path, keyList, options.WithDestinationList(valueList))
	//if err != nil {
	//	logger.Error("store hmget error", zap.Error(err))
	//	return err
	//}
	//logger.Debug("store hmget ", zap.Any("values", valueList))

	err = store.HDel(path, keyList...)
	if err != nil {
		logger.Error("hash delete err:", zap.Error(err))
		return err
	}

	// LPush/LPop
	path = keys.NewKey("test:ltest:demo")
	testList := make([]interface{}, 0)
	testList = append(testList, &listTest{
		MapId: 101,
	})
	testList = append(testList, &listTest{
		MapId: 102,
	})
	testList = append(testList, &listTest{
		MapId: 103,
	})

	err = store.LPush(path, options.WithSourceList(testList))
	if err != nil {
		logger.Error("store lpush error", zap.Error(err))
		return err
	}

	length, err := store.LLen(path)
	if err != nil {
		return err
	}
	logger.Info("store llen", zap.Int64("length", length))

	listTest := &listTest{}
	err = store.LPop(path, options.WithDestination(listTest))
	if err != nil {
		logger.Error("store lpop error", zap.Error(err))
		return err
	}
	logger.Debug("store lpop ", zap.Any("value", listTest))

	//Sub
	ch, err := store.Sub(key)
	if err != nil {
		logger.Error("store sub error", zap.Error(err))
		return err
	}

	go func() {
		ding := time.After(time.Second * 2)
		select {
		case <-ding:
			mt2 := &memTest{}
			mt2.Index = "test_struct2"
			mt2.Value = 5
			store.Set(key, options.WithSource(mt2))
		}
	}()

	select {
	case msg := <-ch:
		if msg == keys.Set {
			logger.Debug("store watch changes", zap.Any("type", keys.Set))
		} else if msg == keys.Expired {
			logger.Debug("store watch changes", zap.Any("type", keys.Expired))
		} else if msg == keys.Delete {
			logger.Debug("store watch changes", zap.Any("type", keys.Delete))
		}
	}

	store.Remove(key)

	return nil
}
