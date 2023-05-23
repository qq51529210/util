package util

import (
	"sync"

	"gorm.io/gorm"
)

// GORMCache 用于缓存数据
type GORMCache[K comparable, M any] struct {
	sync.Mutex
	// 数据
	d map[K]M
	// gorm.Model 要的模型
	m M
	// 标记数据 d 是否需要重新加载
	ok bool
	// 数据库
	db *gorm.DB
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
	newFunc func() M,
	keyFunc func(M) K,
	whereKeyFunc func(*gorm.DB, K) *gorm.DB,
	whereKeysFunc func(*gorm.DB, []K) *gorm.DB) *GORMCache[K, M] {
	c := new(GORMCache[K, M])
	c.d = make(map[K]M)
	c.db = db
	c.new = newFunc
	c.key = keyFunc
	c.whereKey = whereKeyFunc
	c.whereKeys = whereKeysFunc
	c.m = newFunc()
	return c
}

// loadAll 加载所有数据到内存
func (c *GORMCache[K, M]) loadAll() error {
	var models []M
	// 数据库
	err := c.db.Find(&models).Error
	if err != nil {
		c.ok = false
		return err
	}
	// 重置
	c.d = make(map[K]M)
	// 添加
	for _, model := range models {
		c.d[c.key(model)] = model
	}
	c.ok = true
	return nil
}

// Check 检查内存数据是否需要重新加载，同步
func (c *GORMCache[K, M]) Check() error {
	c.Lock()
	err := c.check()
	c.Unlock()
	return err
}

// check 检查内存数据是否需要重新加载
func (c *GORMCache[K, M]) check() error {
	// 数据是否有效
	if !c.ok {
		// 加载
		return c.loadAll()
	}
	return nil
}

// Load 尝试加载，添加和修改时候调用，同步
func (c *GORMCache[K, M]) Load(k K) {
	c.Lock()
	c.load(k)
	c.Unlock()
}

// load 尝试加载，添加和修改时候调用
func (c *GORMCache[K, M]) load(k K) error {
	// 读取
	m := c.new()
	err := c.whereKey(c.db, k).First(m).Error
	// 失败
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.ok = true
			return nil
		}
		c.ok = false
		return err
	}
	// 成功
	c.d[k] = m
	c.ok = true
	//
	return nil
}

// All 返回所有内存，不要修改返回的指针，同步
func (c *GORMCache[K, M]) All() ([]M, error) {
	// 上锁
	c.Lock()
	defer c.Unlock()
	// 确保数据
	err := c.check()
	if err != nil {
		return nil, err
	}
	// 列表
	var models []M
	for _, v := range c.d {
		models = append(models, v)
	}
	// 返回
	return models, nil
}

// Get 返回指定，不要修改返回的指针，同步
func (c *GORMCache[K, M]) Get(k K) (m M, err error) {
	// 上锁
	c.Lock()
	defer c.Unlock()
	// 确保加载
	err = c.check()
	if err != nil {
		return
	}
	// 返回
	m = c.d[k]
	// 返回
	return
}

// Add 添加，同步
func (c *GORMCache[K, M]) Add(m M) (int64, error) {
	// 数据库
	db := c.db.Create(m)
	if db.Error != nil {
		return db.RowsAffected, db.Error
	}
	// 内存
	if db.RowsAffected > 0 {
		c.Load(c.key(m))
	}
	// 返回
	return db.RowsAffected, nil
}

// Update 更新，同步
func (c *GORMCache[K, M]) Update(m M) (int64, error) {
	// 数据库
	k := c.key(m)
	db := c.whereKey(c.db, k).Updates(m)
	if db.Error != nil {
		return db.RowsAffected, db.Error
	}
	// 内存
	if db.RowsAffected > 0 {
		c.Load(c.key(m))
	}
	// 返回
	return db.RowsAffected, nil
}

// Save 保存，同步
func (c *GORMCache[K, M]) Save(m M) (int64, error) {
	// 数据库
	k := c.key(m)
	db := c.whereKey(c.db, k).Save(m)
	if db.Error != nil {
		return db.RowsAffected, db.Error
	}
	// 内存
	if db.RowsAffected > 0 {
		c.Load(c.key(m))
	}
	return db.RowsAffected, nil
}

// Delete 删除
func (c *GORMCache[K, M]) Delete(k K) (int64, error) {
	// 数据库
	db := c.whereKey(c.db, k).Delete(c.m)
	if db.Error != nil {
		return db.RowsAffected, db.Error
	}
	// 内存
	if db.RowsAffected > 0 {
		c.Lock()
		delete(c.d, k)
		c.Unlock()
	}
	// 返回
	return db.RowsAffected, nil
}

// BatchDelete 批量删除
func (c *GORMCache[K, M]) BatchDelete(ks []K) (int64, error) {
	// 数据库
	db := c.whereKeys(c.db, ks).Delete(c.m)
	if db.Error != nil {
		return db.RowsAffected, db.Error
	}
	// 内存
	if db.RowsAffected > 0 {
		// 上锁
		c.Lock()
		for _, k := range ks {
			delete(c.d, k)
		}
		c.Unlock()
	}
	// 返回
	return db.RowsAffected, nil
}
