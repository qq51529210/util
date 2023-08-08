package util

import (
	"context"
	"sync"

	"gorm.io/gorm"
)

// GORMCache 用于缓存数据
type GORMCache[K comparable, M any] struct {
	sync.Mutex
	// 数据库
	db *gorm.DB
	// 是否开启缓存
	cache bool
	// 数据
	D map[K]M
	// gorm.Model 要的模型
	M M
	// 标记数据是否无效，需要重新加载
	OK bool
	// 创建函数
	new func() M
	// 返回 M 的主键
	key func(M) K
	// 返回 M 的主键，用于单个，查询，修改，删除
	whereKey func(*gorm.DB, K) *gorm.DB
	// 返回 M 的主键列表，用于批量删除
	whereKeys func(*gorm.DB, []K) *gorm.DB
}

// NewGORMCache 返回新的缓存，enable 为 false 则不开启缓存
func NewGORMCache[K comparable, M any](
	db *gorm.DB,
	cache bool,
	newFunc func() M,
	keyFunc func(M) K,
	whereKeyFunc func(*gorm.DB, K) *gorm.DB,
	whereKeysFunc func(*gorm.DB, []K) *gorm.DB,
) *GORMCache[K, M] {
	c := new(GORMCache[K, M])
	c.Init(db, cache, newFunc, keyFunc, whereKeyFunc, whereKeysFunc)
	return c
}

// Init 初始化字段
func (c *GORMCache[K, M]) Init(
	db *gorm.DB,
	cache bool,
	newFunc func() M,
	keyFunc func(M) K,
	whereKeyFunc func(*gorm.DB, K) *gorm.DB,
	whereKeysFunc func(*gorm.DB, []K) *gorm.DB,
) {
	c.db = db
	c.cache = cache
	c.D = make(map[K]M)
	c.new = newFunc
	c.key = keyFunc
	c.whereKey = whereKeyFunc
	c.whereKeys = whereKeysFunc
	c.M = newFunc()
}

// IsCache 返回是否启用
func (c *GORMCache[K, M]) IsCache() bool {
	return c.cache
}

// DB 返回 db
func (c *GORMCache[K, M]) DB() *gorm.DB {
	return c.db
}

// Model 返回加载模型的 db
func (c *GORMCache[K, M]) Model() *gorm.DB {
	return c.db.Model(c.M)
}

// ModelWithContext 返回加载模型的 db
func (c *GORMCache[K, M]) ModelWithContext(ctx context.Context) *gorm.DB {
	return c.db.Model(c.M).WithContext(ctx)
}

// loadOne 加载单个，添加和修改时候调用，db 在外面初始化好
// 注意返回错误，要设置 ok 为 false
func (c *GORMCache[K, M]) loadOne(db *gorm.DB, k K) error {
	// 读取
	m := c.new()
	err := c.whereKey(db, k).First(m).Error
	// 失败
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	// 成功
	c.D[k] = m
	//
	return nil
}

// loadMultiple 加载多个，db 在外面初始化好
func (c *GORMCache[K, M]) loadMultiple(db *gorm.DB) error {
	// 查询
	var ms []M
	err := db.Find(&ms).Error
	if err != nil {
		return err
	}
	// 加载或替换
	for _, m := range ms {
		c.D[c.key(m)] = m
	}
	//
	return nil
}

// LoadWhere 加载并替换指定条件的数据，同步
func (c *GORMCache[K, M]) LoadWhere(whereFunc func(db *gorm.DB) *gorm.DB) error {
	return c.LoadWhereWithContext(context.Background(), whereFunc)
}

// LoadWhereWithContext 加载并替换指定条件的数据，同步
func (c *GORMCache[K, M]) LoadWhereWithContext(ctx context.Context, whereFunc func(db *gorm.DB) *gorm.DB) (err error) {
	// 启用
	if c.cache {
		// 上锁
		c.Lock()
		//
		if c.OK {
			// 原数据有效，根据条件加载
			err = c.loadMultiple(whereFunc(c.ModelWithContext(ctx)))
			if err != nil {
				// 标记
				c.OK = false
			}
		} else {
			// 原数据无效直接全部加载
			err = c.loadMultiple(c.ModelWithContext(ctx))
			if err == nil {
				// 标记
				c.OK = true
			}
		}
		// 解锁
		c.Unlock()
	}
	//
	return
}

// LoadAll 重新加载，同步
func (c *GORMCache[K, M]) LoadAll() error {
	return c.LoadAllWithContext(context.Background())
}

// LoadAllWithContext 重新加载，同步
func (c *GORMCache[K, M]) LoadAllWithContext(ctx context.Context) (err error) {
	// 启用
	if c.cache {
		// 上锁
		c.Lock()
		// 加载
		err = c.loadMultiple(c.ModelWithContext(ctx))
		// 标记
		c.OK = err == nil
		// 解锁
		c.Unlock()
	}
	//
	return
}

// Check 检查内存数据是否需要重新加载，同步
func (c *GORMCache[K, M]) Check() error {
	return c.CheckWithContext(context.Background())
}

// CheckWithContext 检查内存数据是否需要重新加载，同步
func (c *GORMCache[K, M]) CheckWithContext(ctx context.Context) (err error) {
	// 启用
	if c.cache {
		// 上锁
		c.Lock()
		// 加载
		err = c.check(c.ModelWithContext(ctx))
		// 解锁
		c.Unlock()
	}
	//
	return
}

// Check 检查内存数据是否需要重新加载，同步
func (c *GORMCache[K, M]) check(db *gorm.DB) (err error) {
	// 确保数据
	if !c.OK {
		err = c.loadMultiple(db)
		// 标记
		c.OK = err == nil
	}
	return
}

// Load 加载单个，添加和修改时候调用，同步
func (c *GORMCache[K, M]) Load(k K) error {
	return c.LoadWithContext(context.Background(), k)
}

// LoadWithContext 加载单个，添加和修改时候调用，同步
func (c *GORMCache[K, M]) LoadWithContext(ctx context.Context, k K) (err error) {
	// 启用
	if c.cache {
		db := c.ModelWithContext(ctx)
		// 上锁
		c.Lock()
		if c.OK {
			// 原数据有效，加载单个
			err = c.loadOne(db, k)
			if err != nil {
				// 标记
				c.OK = false
			}
		} else {
			// 原数据无效直接全部加载
			err = c.loadMultiple(db)
			if err == nil {
				// 标记
				c.OK = true
			}
		}
		// 解锁
		c.Unlock()
	}
	//
	return
}

// All 返回所有内存，不要修改返回的指针，同步
func (c *GORMCache[K, M]) All() ([]M, error) {
	return c.AllWithContext(context.Background())
}

// AllWithContext 返回所有内存，不要修改返回的指针，同步
func (c *GORMCache[K, M]) AllWithContext(ctx context.Context) (ms []M, err error) {
	db := c.ModelWithContext(ctx)
	// 不启用
	if !c.cache {
		err = db.Find(&ms).Error
		if err != nil {
			return nil, err
		}
		return
	}
	// 上锁
	c.Lock()
	// 确保数据
	err = c.check(c.ModelWithContext(ctx))
	if err == nil {
		for _, v := range c.D {
			ms = append(ms, v)
		}
	}
	// 解锁
	c.Unlock()
	//
	return
}

// Get 返回指定，不要修改返回的指针，同步
func (c *GORMCache[K, M]) Get(k K) (m M, err error) {
	return c.GetWithContext(context.Background(), k)
}

// GetWithContext 返回指定，不要修改返回的指针，同步
func (c *GORMCache[K, M]) GetWithContext(ctx context.Context, k K) (m M, err error) {
	// 不启用
	if !c.cache {
		mm := c.new()
		err = c.whereKey(c.ModelWithContext(ctx), k).First(mm).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				err = nil
			}
			return
		}
		m = mm
		//
		return
	}
	// 上锁
	c.Lock()
	// 确保数据
	err = c.check(c.ModelWithContext(ctx))
	if err == nil {
		m = c.D[k]
	}
	// 解锁
	c.Unlock()
	//
	return
}

// Add 添加，同步
func (c *GORMCache[K, M]) Add(m M) (int64, error) {
	return c.AddWithContext(context.Background(), m)
}

// AddWithContext 添加，同步
func (c *GORMCache[K, M]) AddWithContext(ctx context.Context, m M) (int64, error) {
	// 数据库
	db := c.ModelWithContext(ctx).Create(m)
	if db.Error == nil {
		// 内存
		if db.RowsAffected > 0 {
			c.LoadWithContext(ctx, c.key(m))
		}
	}
	//
	return db.RowsAffected, db.Error
}

// Update 更新，同步
func (c *GORMCache[K, M]) Update(m M) (int64, error) {
	return c.UpdateWithContext(context.Background(), m)
}

// UpdateWithContext 更新，同步
func (c *GORMCache[K, M]) UpdateWithContext(ctx context.Context, m M) (int64, error) {
	// 数据库
	k := c.key(m)
	db := c.whereKey(c.ModelWithContext(ctx), k).Updates(m)
	if db.Error == nil {
		// 内存
		if db.RowsAffected > 0 {
			c.LoadWithContext(ctx, k)
		}
	}
	//
	return db.RowsAffected, db.Error
}

// BatchUpdate 事务更新，同步
func (c *GORMCache[K, M]) BatchUpdate(ms []M) (int64, error) {
	return c.BatchUpdateWithContext(context.Background(), ms)
}

// BatchUpdateWithContext 事务更新，同步
func (c *GORMCache[K, M]) BatchUpdateWithContext(ctx context.Context, ms []M) (int64, error) {
	var ks []K
	// 数据库
	err := c.ModelWithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, m := range ms {
			k := c.key(m)
			db := c.whereKey(tx, k).Updates(m)
			if db.Error != nil {
				return db.Error
			}
			ks = append(ks, k)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	// 内存
	c.LoadWhereWithContext(ctx, func(db *gorm.DB) *gorm.DB {
		return c.whereKeys(db, ks)
	})
	//
	return int64(len(ms)), nil
}

// UpdateCache 更新内存，同步。回调有可能为 nil
func (c *GORMCache[K, M]) UpdateCache(k K, fn func(M)) {
	// 上锁
	c.Lock()
	// 更新
	fn(c.D[k])
	// 解锁
	c.Unlock()
}

// Save 保存，同步
func (c *GORMCache[K, M]) Save(m M) (int64, error) {
	return c.SaveWithContext(context.Background(), m)
}

// SaveWithContext 保存，同步
func (c *GORMCache[K, M]) SaveWithContext(ctx context.Context, m M) (int64, error) {
	// 数据库
	k := c.key(m)
	db := c.whereKey(c.ModelWithContext(ctx), k).Save(m)
	if db.Error == nil {
		// 内存
		if db.RowsAffected > 0 {
			c.LoadWithContext(ctx, k)
		}
	}
	//
	return db.RowsAffected, db.Error
}

// BatchSave 事务保存，同步
func (c *GORMCache[K, M]) BatchSave(ms []M) (int64, error) {
	return c.BatchSaveWithContext(context.Background(), ms)
}

// BatchSaveWithContext 事务保存，同步
func (c *GORMCache[K, M]) BatchSaveWithContext(ctx context.Context, ms []M) (int64, error) {
	var ks []K
	// 数据库
	err := c.ModelWithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, m := range ms {
			k := c.key(m)
			db := c.whereKey(tx, k).Save(m)
			if db.Error != nil {
				return db.Error
			}
			ks = append(ks, k)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	// 内存
	c.LoadWhereWithContext(ctx, func(db *gorm.DB) *gorm.DB {
		return c.whereKeys(db, ks)
	})
	//
	return int64(len(ms)), nil
}

// Delete 删除，同步
func (c *GORMCache[K, M]) Delete(k K) (int64, error) {
	return c.DeleteWithContext(context.Background(), k)
}

// DeleteWithContext 删除，同步
func (c *GORMCache[K, M]) DeleteWithContext(ctx context.Context, k K) (int64, error) {
	// 数据库
	db := c.whereKey(c.ModelWithContext(ctx), k).Delete(c.M)
	if db.Error == nil {
		// 内存
		if db.RowsAffected > 0 && c.cache {
			c.DeleteCache(k)
		}
	}
	//
	return db.RowsAffected, db.Error
}

// DeleteCache 删除内存
func (c *GORMCache[K, M]) DeleteCache(k K) {
	// 上锁
	c.Lock()
	// 删除
	delete(c.D, k)
	// 解锁
	c.Unlock()
}

// BatchDelete 批量删除，同步
func (c *GORMCache[K, M]) BatchDelete(ks []K) (int64, error) {
	return c.BatchDeleteWithContext(context.Background(), ks)
}

// BatchDeleteWithContext 批量删除，同步
func (c *GORMCache[K, M]) BatchDeleteWithContext(ctx context.Context, ks []K) (int64, error) {
	// 数据库
	db := c.whereKeys(c.ModelWithContext(ctx), ks).Delete(c.M)
	if db.Error == nil {
		// 内存
		if db.RowsAffected > 0 && c.cache {
			c.BatchDeleteCache(ks)
		}
	}
	//
	return db.RowsAffected, db.Error
}

// BatchDeleteCache 删除内存，同步
func (c *GORMCache[K, M]) BatchDeleteCache(ks []K) {
	// 上锁
	c.Lock()
	// 删除
	for _, k := range ks {
		delete(c.D, k)
	}
	// 解锁
	c.Unlock()
}

// ForeachCache 遍历缓存，同步
func (c *GORMCache[K, M]) ForeachCache(fc func(M)) (err error) {
	return c.ForeachCacheWithContext(context.Background(), fc)
}

// ForeachCacheWithContext 遍历缓存，同步
func (c *GORMCache[K, M]) ForeachCacheWithContext(ctx context.Context, fc func(M)) (err error) {
	// 上锁
	c.Lock()
	// 确保数据
	err = c.check(c.ModelWithContext(ctx))
	if err == nil {
		// 循环
		for _, m := range c.D {
			fc(m)
		}
	}
	// 解锁
	c.Unlock()
	//
	return
}

// SearchCache 在内存中查找所有，同步
func (c *GORMCache[K, M]) SearchCache(match func(M) bool) ([]M, error) {
	return c.SearchCacheWithContext(context.Background(), match)
}

// SearchCacheWithContext 在内存中查找所有，同步
func (c *GORMCache[K, M]) SearchCacheWithContext(ctx context.Context, match func(M) bool) (mm []M, err error) {
	// 上锁
	c.Lock()
	// 确保数据
	err = c.check(c.ModelWithContext(ctx))
	if err == nil {
		// 查找
		for _, m := range c.D {
			if match(m) {
				mm = append(mm, m)
			}
		}
	}
	// 解锁
	c.Unlock()
	//
	return
}

// SearchCacheIn 在内存中查找 key ，同步
func (c *GORMCache[K, M]) SearchCacheIn(ks []K) ([]M, error) {
	return c.SearchCacheInWithContext(context.Background(), ks)
}

// SearchCacheInWithContext 在内存中查找 key ，同步
func (c *GORMCache[K, M]) SearchCacheInWithContext(ctx context.Context, ks []K) (mm []M, err error) {
	// 上锁
	c.Lock()
	// 确保数据
	err = c.check(c.ModelWithContext(ctx))
	if err == nil {
		// 查找
		for i := 0; i < len(ks); i++ {
			m, ok := c.D[ks[i]]
			if ok {
				mm = append(mm, m)
			}
		}
	}
	// 解锁
	c.Unlock()
	//
	return
}

// SearchCacheOne 在内存中查找第一个，同步
func (c *GORMCache[K, M]) SearchCacheOne(match func(M) bool) (m M, err error) {
	return c.SearchCacheOneWithContext(context.Background(), match)
}

// SearchCacheOneWithContext 在内存中查找第一个，同步
func (c *GORMCache[K, M]) SearchCacheOneWithContext(ctx context.Context, match func(M) bool) (m M, err error) {
	// 上锁
	c.Lock()
	// 确保数据
	err = c.check(c.ModelWithContext(ctx))
	if err == nil {
		// 查找
		for _, v := range c.D {
			if match(v) {
				m = v
				break
			}
		}
	}
	// 解锁
	c.Unlock()
	//
	return
}

// CacheCount 返回内存匹配数量，同步
func (c *GORMCache[K, M]) CacheCount(match func(M) bool) (int64, error) {
	return c.CacheCountWithContext(context.Background(), match)
}

// CacheCountWithContext 返回内存匹配数量，同步
func (c *GORMCache[K, M]) CacheCountWithContext(ctx context.Context, match func(M) bool) (n int64, err error) {
	// 上锁
	c.Lock()
	// 确保数据
	err = c.check(c.ModelWithContext(ctx))
	if err == nil {
		// 查找
		for _, v := range c.D {
			if match(v) {
				n++
			}
		}
	}
	// 解锁
	c.Unlock()
	//
	return
}

// CacheTotal 返回内存总量，同步
func (c *GORMCache[K, M]) CacheTotal() (int64, error) {
	return c.CacheTotalWithContext(context.Background())
}

// CacheTotalWithContext 返回内存总量，同步
func (c *GORMCache[K, M]) CacheTotalWithContext(ctx context.Context) (n int64, err error) {
	// 上锁
	c.Lock()
	// 确保数据
	err = c.check(c.ModelWithContext(ctx))
	if err == nil {
		n = int64(len(c.D))
	}
	// 解锁
	c.Unlock()
	//
	return
}

// List 查询数据库
func (c *GORMCache[K, M]) List(page *GORMPage, query GORMQuery, res *GORMList[M]) error {
	return c.ListWithContext(context.Background(), page, query, res)
}

// ListWithContext 查询数据库
func (c *GORMCache[K, M]) ListWithContext(ctx context.Context, page *GORMPage, query GORMQuery, res *GORMList[M]) error {
	return gormList(c.ModelWithContext(ctx), page, query, res)
}

// GORMSearchCache 模板化的 Search
func GORMSearchCache[T any, K comparable, M any](c *GORMCache[K, M], match func(M) (bool, T)) ([]T, error) {
	return GORMSearchCacheWithContext(context.Background(), c, match)
}

// GORMSearchCacheWithContext 模板化的 Search
func GORMSearchCacheWithContext[T any, K comparable, M any](ctx context.Context, c *GORMCache[K, M], match func(M) (bool, T)) ([]T, error) {
	var vv []T
	// 上锁
	c.Lock()
	// 确保数据
	err := c.check(c.ModelWithContext(ctx))
	if err == nil {
		// 查找
		for _, m := range c.D {
			o, v := match(m)
			if o {
				vv = append(vv, v)
			}
		}
	}
	// 解锁
	c.Unlock()
	//
	return vv, err
}
