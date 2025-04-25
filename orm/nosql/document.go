package nosql

import (
	"context"
	"math"
	"math/rand"
	"reflect"
	"time"

	"github.com/pkg/errors"

	"github.com/gstones/moke-kit/mq/miface"
	"github.com/gstones/moke-kit/orm/nerrors"
	"github.com/gstones/moke-kit/orm/nosql/diface"
	"github.com/gstones/moke-kit/orm/nosql/key"
	"github.com/gstones/moke-kit/orm/nosql/noptions"
)

const (
	// MaxRetries 是更新操作的最大重试次数
	MaxRetries = 5
	// DefaultWriteBackDelay 默认回写延迟时间
	DefaultWriteBackDelay = 500 * time.Millisecond
	ExpireRangeMin        = 6 * time.Hour
	ExpireRangeMax        = 12 * time.Hour
	// WriteBackTopic 延迟回写的消息队列主题
	WriteBackTopic = "nats://writeback"
)

// WriteBackPayload 定义回写消息的有效载荷
type WriteBackPayload struct {
	CollectionName string           `json:"collection"`
	Key            string           `json:"key"`
	Data           map[string]any   `json:"data"`
	Version        noptions.Version `json:"version"`
}

// WriteBackOptions 定义回写选项
type WriteBackOptions struct {
	// Enabled 是否启用回写
	Enabled bool
	// Delay 回写延迟时间
	Delay time.Duration
	// MQ 消息队列客户端实例
	MQ miface.MessageQueue
}

// DefaultWriteBackOptions 返回默认回写选项
func DefaultWriteBackOptions() WriteBackOptions {
	return WriteBackOptions{
		Enabled: false,
		Delay:   DefaultWriteBackDelay,
		MQ:      nil,
	}
}

// DocumentBase 表示 NoSQL 操作的基础文档结构
type DocumentBase struct {
	Key key.Key

	clear   func()
	data    any
	version noptions.Version

	DocumentStore diface.ICollection
	cache         diface.ICache
	ctx           context.Context

	// 回写相关字段
	writeBack WriteBackOptions
}

// Init 执行 DocumentBase 的就地初始化
func (d *DocumentBase) Init(
	ctx context.Context,
	data any,
	clear func(),
	store diface.ICollection,
	key key.Key,
) {
	defaultCache := diface.DefaultDocumentCache()
	d.InitWithCache(ctx, data, clear, store, key, defaultCache)
}

// InitWithCache 使用自定义缓存执行 DocumentBase 的就地初始化
func (d *DocumentBase) InitWithCache(
	ctx context.Context,
	data any,
	clear func(),
	store diface.ICollection,
	key key.Key,
	cache diface.ICache,
) {
	d.ctx = ctx
	d.clear = clear
	d.data = data
	d.DocumentStore = store
	d.Key = key
	d.cache = cache
	d.writeBack = DefaultWriteBackOptions()
}

// EnableWriteBackWithMQ 启用基于消息队列的延迟回写功能
func (d *DocumentBase) EnableWriteBackWithMQ(mqClient miface.MessageQueue, delay time.Duration) error {
	if mqClient == nil {
		return errors.New("MQ client cannot be nil")
	}

	d.writeBack.Enabled = true
	d.writeBack.MQ = mqClient
	if delay > 0 {
		d.writeBack.Delay = delay
	}

	return nil
}

// DisableWriteBack 禁用延迟回写功能
func (d *DocumentBase) DisableWriteBack() {
	d.writeBack.Enabled = false
	d.writeBack.MQ = nil
}

// 将修改通过MQ发送延迟回写消息
func (d *DocumentBase) scheduleWriteBackWithMQ(changes map[string]any) error {
	// 准备回写数据
	payload := &WriteBackPayload{
		CollectionName: d.DocumentStore.GetName(),
		Key:            d.Key.String(),
		Data:           changes,
		Version:        d.version,
	}

	// 发送延迟消息
	opts := []miface.PubOption{
		miface.WithJSON(payload),
	}
	return d.writeBack.MQ.Publish(WriteBackTopic, opts...)
}

// Clear 清除此 DocumentBase 上的所有数据
func (d *DocumentBase) Clear() {
	d.version = noptions.NoVersion
	d.clear()
}

// Create 在数据库中创建数据和版本
func (d *DocumentBase) Create() error {
	version, err := d.DocumentStore.Set(
		d.ctx,
		d.Key,
		noptions.WithSource(d.data),
	)
	if err != nil {
		return err
	}

	d.version = version
	// 创建后更新缓存
	return d.updateCache()
}

// 生成随机的缓存过期时间，防止缓存雪崩
func randomExpiration() time.Duration {
	return ExpireRangeMin + time.Duration(rand.Int63n(int64(ExpireRangeMax-ExpireRangeMin)))
}

// 更新缓存
func (d *DocumentBase) updateCacheChanges(changes map[string]any) error {
	data, err := marshalAnyMap(changes)
	if err != nil {
		return err
	}
	return d.cache.SetCache(d.ctx, d.Key, data, randomExpiration())
}

// 全量更新缓存

func (d *DocumentBase) updateCache() error {
	dataMap, err := struct2MapShallow(d.data)
	if err != nil {
		return errors.Wrap(err, "failed to convert data to map")
	}
	return d.updateCacheChanges(dataMap)
}

// Load 实现读穿透缓存
func (d *DocumentBase) Load() error {
	d.clear()

	// 首先尝试缓存
	if data := d.cache.GetCache(d.ctx, d.Key); len(data) > 0 {
		return map2StructShallow(data, d.data)
	}

	// 缓存未命中 - 从数据库加载
	version, err := d.DocumentStore.Get(
		d.ctx,
		d.Key,
		noptions.WithDestination(d.data),
	)
	if err != nil {
		return err
	}
	d.version = version
	// 从数据库加载后更新缓存
	return d.updateCache()
}

// SaveAsync 实现异步写入，更新缓存并安排延迟回写
func (d *DocumentBase) SaveAsync(changes map[string]any) error {
	// 更新缓存
	if err := d.updateCacheChanges(changes); err != nil {
		return err
	}

	// 如果启用了回写，安排回写
	if d.writeBack.Enabled && d.writeBack.MQ != nil {
		go func() {
			if err := d.scheduleWriteBackWithMQ(changes); err != nil {
				return
			}
		}()
	}

	return nil
}

// Save 实现同步写入并删除缓存
func (d *DocumentBase) Save() error {
	// 直接同步写入数据库
	version, err := d.DocumentStore.Set(
		d.ctx,
		d.Key,
		noptions.WithSource(d.data),
		noptions.WithVersion(d.version),
	)
	if err != nil {
		return err
	}
	d.version = version

	// 删除缓存
	d.cache.DeleteCache(d.ctx, d.Key)
	return nil
}

func (d *DocumentBase) doUpdate(f func() bool, u func() error) error {
	var lastErr error
	for r := 0; r < MaxRetries; r++ {
		if !f() {
			return nerrors.ErrUpdateLogicFailed
		}

		if err := u(); err == nil {
			return nil // 成功更新
		} else {
			lastErr = err
			// 指数退避加抖动
			backoff := time.Duration(math.Pow(2, float64(r))) * time.Millisecond
			jitter := time.Duration(rand.Float64() * float64(backoff))
			time.Sleep(backoff + jitter)

			if err := d.Load(); err != nil {
				return errors.Wrap(err, "failed to reload during update retry")
			}
		}
	}

	if lastErr != nil {
		return errors.Wrap(nerrors.ErrTooManyRetries, lastErr.Error())
	}
	return nerrors.ErrTooManyRetries
}

// Update 使用给定函数更改数据并通过CAS(比较并交换)将其保存到数据库
// 如果函数返回false，更新将被中止
// 如果更新CAS失败，函数将重试最多 MaxRetries 次，并使用随机化的退避策略
func (d *DocumentBase) Update(f func() bool) error {
	return d.doUpdate(f, func() error {
		return d.Save()
	})
}

// diffMapAny 比较两个 map[string]any，返回不同的键值对
func diffMapAny(oldData map[string]any, newData map[string]any) (map[string]any, error) {
	changes := make(map[string]any)
	// Compare and collect changes
	for k, newVal := range newData {
		if oldVal, exists := oldData[k]; !exists || !reflect.DeepEqual(oldVal, newVal) {
			changes[k] = newVal
		}
	}

	return changes, nil
}

// UpdateAsync 异步更新数据，更改后安排延迟回写
func (d *DocumentBase) UpdateAsync(f func() bool) error {
	oldData, err := struct2MapShallow(d.data)
	if err != nil {
		return err
	}
	if !f() {
		return nerrors.ErrUpdateLogicFailed
	}

	newData, err := struct2MapShallow(d.data)
	if err != nil {
		return err
	}

	changes, err := diffMapAny(oldData, newData)
	if err != nil {
		return err
	}
	return d.SaveAsync(changes)
}

// Delete 从数据库中删除数据
func (d *DocumentBase) Delete() error {
	if err := d.DocumentStore.Delete(d.ctx, d.Key); err != nil {
		return err
	}
	d.cache.DeleteCache(d.ctx, d.Key)
	return nil
}
