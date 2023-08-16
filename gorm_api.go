package util

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

// GORMBaseModel 基本字段
type GORMBaseModel[T comparable] struct {
	// 数据库ID
	ID T `json:"id" gorm:""`
	GORMTimeModel
}

// GORMTimeModel 创建和更新时间
type GORMTimeModel struct {
	// 创建时间戳
	CreatedAt int64 `json:"createdAt" gorm:""`
	// 更新时间戳
	UpdatedAt int64 `json:"updatedAt" gorm:""`
}

// GORMQuery 查询参数的接口
type GORMQuery interface {
	Init(*gorm.DB) *gorm.DB
}

// GORMListPage 分页查询参数
type GORMListPage struct {
	// 偏移，小于 0 不匹配
	Offset *int `form:"offset" binding:"omitempty,min=0"`
	// 条数，小于 1 不匹配
	Count *int `form:"count" binding:"omitempty,min=1"`
	// 排序，"column [desc]"
	Order string `form:"order"`
}

// Init 初始化 db
func (m *GORMListPage) Init(db *gorm.DB) *gorm.DB {
	// 分页
	if m.Offset != nil {
		db = db.Offset(*m.Offset)
	}
	if m.Count != nil {
		db = db.Limit(*m.Count)
	}
	// 排序
	if m.Order != "" {
		db = db.Order(m.Order)
	}
	//
	return db
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
	D *gorm.DB
	M M
}

// NewGORMDB 返回新的 GORMDB
func NewGORMDB[K, M any](db *gorm.DB, m M) *GORMDB[K, M] {
	return &GORMDB[K, M]{
		D: db,
		M: m,
	}
}

// Init 初始化
func (g *GORMDB[K, M]) Init(db *gorm.DB, m M) {
	g.D = db
	g.M = m
}

// Model 返回
func (g *GORMDB[K, M]) Model() *gorm.DB {
	return g.D.Model(g.M)
}

// ModelWithContext 返回
func (g *GORMDB[K, M]) ModelWithContext(ctx context.Context) *gorm.DB {
	return g.D.Model(g.M).WithContext(ctx)
}

// All 返回列表查询结果
func (g *GORMDB[K, M]) All(query GORMQuery) ([]M, error) {
	return g.AllWithContext(context.Background(), query)
}

// AllWithContext 返回列表查询结果
func (g *GORMDB[K, M]) AllWithContext(ctx context.Context, query GORMQuery) ([]M, error) {
	return GORMAll[M](g.ModelWithContext(ctx), query)
}

// List 返回分页查询结果
func (g *GORMDB[K, M]) List(page *GORMListPage, query GORMQuery, res *GORMListData[M]) error {
	return g.ListWithContext(context.Background(), page, query, res)
}

// ListWithContext 返回分页查询结果
func (g *GORMDB[K, M]) ListWithContext(ctx context.Context, page *GORMListPage, query GORMQuery, res *GORMListData[M]) error {
	return GORMList(g.ModelWithContext(ctx), page, query, res)
}

// Save 保存
func (g *GORMDB[K, M]) Save(m M) (int64, error) {
	return g.SaveWithContext(context.Background(), m)
}

// SaveWithContext 保存
func (g *GORMDB[K, M]) SaveWithContext(ctx context.Context, m M) (int64, error) {
	db := g.ModelWithContext(ctx).Save(m)
	return db.RowsAffected, db.Error
}

// Add 添加
func (g *GORMDB[K, M]) Add(m M) (int64, error) {
	return g.AddWithContext(context.Background(), m)
}

// AddWithContext 添加
func (g *GORMDB[K, M]) AddWithContext(ctx context.Context, m M) (int64, error) {
	db := g.ModelWithContext(ctx).Create(m)
	return db.RowsAffected, db.Error
}

// Update 更新
func (g *GORMDB[K, M]) Update(m M) (int64, error) {
	return g.UpdateWithContext(context.Background(), m)
}

// UpdateWithContext 更新
func (g *GORMDB[K, M]) UpdateWithContext(ctx context.Context, m M) (int64, error) {
	db := g.ModelWithContext(ctx).Updates(m)
	return db.RowsAffected, db.Error
}

// Delete 删除
func (g *GORMDB[K, M]) Delete(k K) (int64, error) {
	return g.DeleteWithContext(context.Background(), k)
}

// DeleteWithContext 删除
func (g *GORMDB[K, M]) DeleteWithContext(ctx context.Context, k K) (int64, error) {
	db := g.ModelWithContext(ctx).Delete(g.M, k)
	return db.RowsAffected, db.Error
}

// BatchDelete 批量删除
func (g *GORMDB[K, M]) BatchDelete(ks []K) (int64, error) {
	return g.BatchDeleteWithContext(context.Background(), ks)
}

// BatchDeleteWithContext 批量删除
func (g *GORMDB[K, M]) BatchDeleteWithContext(ctx context.Context, ks []K) (int64, error) {
	db := g.ModelWithContext(ctx).Delete(g.M, ks)
	return db.RowsAffected, db.Error
}

// Get 查询
func (g *GORMDB[K, M]) Get(m M) (bool, error) {
	return g.GetWithContext(context.Background(), m)
}

// GetWithContext 查询
func (g *GORMDB[K, M]) GetWithContext(ctx context.Context, m M) (bool, error) {
	err := g.ModelWithContext(ctx).First(m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Select 查询选择列
func (g *GORMDB[K, M]) Select(m M, c ...string) (bool, error) {
	return g.SelectWithContext(context.Background(), m, c...)
}

// SelectWithContext 查询选择列
func (g *GORMDB[K, M]) SelectWithContext(ctx context.Context, m M, c ...string) (bool, error) {
	db := g.ModelWithContext(ctx)
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
func (g *GORMDB[K, M]) In(ks []K) ([]M, error) {
	return g.InWithContext(context.Background(), ks)
}

// InWithContext 根据主键查询，where in ks
func (g *GORMDB[K, M]) InWithContext(ctx context.Context, ks []K) ([]M, error) {
	var ms []M
	err := g.ModelWithContext(ctx).Find(&ms, ks).Error
	if err != nil {
		return nil, err
	}
	return ms, nil
}

var (
	// GORMInitQueryTag 是 GORMInitQuery 解析 tag 的名称
	GORMInitQueryTag = "gq"
)

// GORMInitQuery 将 q 格式化到 where ，全部是 AND ，略过空值
//
//	type query struct {
//	  A *int64 `gq:"eq"` db.Where("`A` = ?", A)
//	  B string `gq:"like"` db.Where("`B` LiKE %%s%%", B)
//	  C *int64 `gq:"gt=A"` db.Where("`A` < ?", C)
//	  D *int64 `gq:"gte=A"` db.Where("`A` <= ?", D)
//	  E *int64 `gq:"lt=A"` db.Where("`A` > ?", E)
//	  F *int64 `gq:"let=A"` db.Where("`A` >= ?", F)
//	  G *int64 `gq:"neq"` db.Where("`G` != ?", G)
//	}
//
// 先这样，以后遇到再加
func GORMInitQuery(db *gorm.DB, q any) *gorm.DB {
	v := reflect.ValueOf(q)
	vk := v.Kind()
	if vk == reflect.Pointer {
		v = v.Elem()
		vk = v.Kind()
	}
	if vk != reflect.Struct {
		panic("v must be struct or struct ptr")
	}
	return gormInitQuery(db, v)
}

func gormInitQuery(db *gorm.DB, v reflect.Value) *gorm.DB {
	vt := v.Type()
	for i := 0; i < vt.NumField(); i++ {
		fv := v.Field(i)
		if !fv.IsValid() {
			continue
		}
		fvk := fv.Kind()
		if fvk == reflect.Pointer {
			// 空指针
			if fv.IsNil() {
				continue
			}
			fv = fv.Elem()
			fvk = fv.Kind()
		}
		// 结构
		if fvk == reflect.Struct {
			gormInitQuery(db, fv)
			continue
		}
		if fvk == reflect.String {
			// 空字符串
			if fv.IsZero() {
				continue
			}
		}
		ft := vt.Field(i)
		tn := ft.Tag.Get(GORMInitQueryTag)
		p := strings.TrimPrefix(tn, "eq=")
		if p != tn {
			db = db.Where(fmt.Sprintf("`%s` = ?", p), fv.Interface())
			continue
		}
		if tn == "eq" {
			db = db.Where(fmt.Sprintf("`%s` = ?", ft.Name), fv.Interface())
			continue
		}
		p = strings.TrimPrefix(tn, "neq=")
		if p != tn {
			db = db.Where(fmt.Sprintf("`%s` != ?", p), fv.Interface())
			continue
		}
		if tn == "neq" {
			db = db.Where(fmt.Sprintf("`%s` != ?", ft.Name), fv.Interface())
			continue
		}
		p = strings.TrimPrefix(tn, "like=")
		if p != tn {
			db = db.Where(fmt.Sprintf("`%s` LIKE ?", p), fmt.Sprintf("%%%v%%", fv.Interface()))
			continue
		}
		if tn == "like" {
			db = db.Where(fmt.Sprintf("`%s` LIKE ?", ft.Name), fmt.Sprintf("%%%v%%", fv.Interface()))
			continue
		}
		p = strings.TrimPrefix(tn, "gt=")
		if p != tn {
			db = db.Where(fmt.Sprintf("`%s` < ?", p), fv.Interface())
			continue
		}
		if tn == "gt" {
			db = db.Where(fmt.Sprintf("`%s` < ?", ft.Name), fv.Interface())
			continue
		}
		p = strings.TrimPrefix(tn, "gte=")
		if p != tn {
			db = db.Where(fmt.Sprintf("`%s` <= ?", p), fv.Interface())
			continue
		}
		if tn == "gte" {
			db = db.Where(fmt.Sprintf("`%s` <= ?", ft.Name), fv.Interface())
			continue
		}
		p = strings.TrimPrefix(tn, "lt=")
		if p != tn {
			db = db.Where(fmt.Sprintf("`%s` > ?", p), fv.Interface())
			continue
		}
		if tn == "lt" {
			db = db.Where(fmt.Sprintf("`%s` > ?", ft.Name), fv.Interface())
			continue
		}
		p = strings.TrimPrefix(tn, "lte=")
		if p != tn {
			db = db.Where(fmt.Sprintf("`%s` >= ?", p), fv.Interface())
			continue
		}
		if tn == "lte" {
			db = db.Where(fmt.Sprintf("`%s` >= ?", ft.Name), fv.Interface())
			continue
		}
	}
	//
	return db
}

// GORMList 分页查询
func GORMList[M any](db *gorm.DB, page *GORMListPage, query GORMQuery, res *GORMListData[M]) error {
	// 条件
	if query != nil {
		db = query.Init(db)
	}
	// 总数
	err := db.Count(&res.Total).Error
	if err != nil {
		return err
	}
	// 分页
	if page != nil {
		db = page.Init(db)
	}
	// 查询
	err = db.Find(&res.Data).Error
	if err != nil {
		return err
	}
	//
	return nil
}

// GORMAll 查询
func GORMAll[M any](db *gorm.DB, query GORMQuery) (mm []M, err error) {
	// 条件
	if query != nil {
		db = query.Init(db)
	}
	// 查询
	err = db.Find(&mm).Error
	//
	return
}
