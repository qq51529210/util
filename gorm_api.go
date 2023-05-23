package util

import (
	"fmt"

	"gorm.io/gorm"
)

// GORMPage 分页查询参数
type GORMPage struct {
	// 偏移，小于 0 不匹配
	Offset *int `form:"offset" binding:"omitempty,min=0"`
	// 条数，小于 1 不匹配
	Count *int `form:"count" binding:"omitempty,min=1"`
	// 排序，"column [desc]"
	Order string `form:"order"`
}

// GORMQuery 是 All 函数格式化查询参数的接口
type GORMQuery interface {
	Init(*gorm.DB) *gorm.DB
}

// GORMListData 是 GORMList 的返回值
type GORMListData[M any] struct {
	// 总数
	Total int64 `json:"total"`
	// 列表
	Data []M `json:"data"`
}

// GORMDB 模板 api
type GORMDB[K, M any] struct {
	db *gorm.DB
	m  M
}

// NewGORMDB 返回新的 GORMDB
func NewGORMDB[K, M any](db *gorm.DB, m M) *GORMDB[K, M] {
	return &GORMDB[K, M]{
		db: db,
		m:  m,
	}
}

// All 返回列表查询结果
func (g *GORMDB[K, M]) All(query GORMQuery) ([]M, error) {
	db := g.db
	// 条件
	if query != nil {
		db = query.Init(db)
	}
	// 列表
	var models []M
	err := db.Find(&models).Error
	if err != nil {
		return nil, err
	}
	// 返回
	return models, nil
}

// List 返回列表查询结果
func (g *GORMDB[K, M]) List(query GORMQuery, page *GORMPage, res *GORMListData[M]) error {
	db := g.db
	// 条件
	if query != nil {
		db = query.Init(db)
	}
	// 总数
	err := db.Count(&res.Total).Error
	if err != nil {
		return err
	}
	if page != nil {
		// 分页
		if page.Offset != nil {
			db = db.Offset(*page.Offset)
		}
		if page.Count != nil {
			db = db.Limit(*page.Count)
		}
		// 排序
		if page.Order != "" {
			db = db.Order(page.Order)
		}
	}
	err = db.Find(&res.Data).Error
	if err != nil {
		return err
	}
	//
	return nil
}

// Like 生成 column like %name%
func (g *GORMDB[K, M]) Like(column, value string) *gorm.DB {
	if value == "" {
		return g.db
	}
	return g.db.Where(fmt.Sprintf("`%s` LIKE '%%%s%%'", column, value))
}

// Save 添加
func (g *GORMDB[K, M]) Save(m M) (int64, error) {
	db := g.db.Save(m)
	return db.RowsAffected, db.Error
}

// Add 添加
func (g *GORMDB[K, M]) Add(m *M) (int64, error) {
	db := g.db.Create(m)
	return db.RowsAffected, db.Error
}

// Update 根据主键更新
func (g *GORMDB[K, M]) Update(m *M) (int64, error) {
	db := g.db.Updates(m)
	return db.RowsAffected, db.Error
}

// Delete 根据主键删除
func (g *GORMDB[K, M]) Delete(m *M) (int64, error) {
	db := g.db.Delete(m)
	return db.RowsAffected, db.Error
}

// BatchDelete 根据主键批量删除
func (g *GORMDB[K, M]) BatchDelete(ks []K) (int64, error) {
	db := g.db.Delete(g.m, ks)
	return db.RowsAffected, db.Error
}

// Get 根据主键查询
func (g *GORMDB[K, M]) Get(m *M) (bool, error) {
	err := g.db.First(m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Select 根据主键查询，可选择列
func (g *GORMDB[K, M]) Select(m *M, c ...string) (bool, error) {
	db := g.db
	if len(c) > 0 {
		db = db.Select(c)
	}
	err := db.First(m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// In 根据主键查询，where in ks
func (g *GORMDB[K, M]) In(ks []K) ([]*M, error) {
	var ms []*M
	err := g.db.Find(&ms, ks).Error
	if err != nil {
		return nil, err
	}
	return ms, nil
}
