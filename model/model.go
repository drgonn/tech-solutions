package model

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"user-track/global"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID uint64 `gorm:"primary_key" json:"id,omitempty"`
	// CreatedAt *time.Time `json:"created_at,omitempty" db:"created_at" gorm:"column:created_at"`
	// UpdatedAt time.Time  `json:"updated_at,omitempty" db:"updated_at" gorm:"column:updated_at"`
}

type BaseDAO interface {
	Create(entity interface{}) error
	CreateOmitColumns(entity interface{}, columns []string) error
	Update(entity interface{}, index string, uid interface{}) (int64, error)
	UpdatePatch(entity interface{}, uid uint64) (int64, error)
	Delete(entity interface{}, index string, uid interface{}) (int64, error)
	FindAll(entity interface{}, opts ...DAOOption) error
	FindByKeys(entity interface{}, keys map[string]interface{}) (int64, error)
	FindByKeysNum(entity interface{}, keys map[string]interface{}) (int64, error)
	FindByPages(entity interface{}, currentPage, pageSize int) error
	FindByPagesWithKeys(entity interface{}, keys map[string]interface{}, currentPage, pageSize int, opts ...DAOOption) error
	SearchByPagesWithKeys(entity interface{}, keys, keyOpts map[string]interface{}, currentPage, pageSize int, opts ...DAOOption) error
	Count(entity interface{}, count *int64) error
	CountWithKeys(entity interface{}, count *int64, keys, keyOpts map[string]interface{}, opts ...DAOOption) error
}

type DAOOption struct {
	Select string
	Order  string
	Where  map[string]interface{}
	Limit  int
}

func (u *BaseModel) Create(entity interface{}) error {
	result := global.GormDb.Model(entity).Create(entity)
	return result.Error
}

// create
func (u *BaseModel) CreateOmitColumns(entity interface{}, columns []string) error {
	columns = append(columns, "create_time", "update_time")
	result := global.GormDb.Model(entity).Omit(columns...).Create(entity)
	return result.Error
}

func (u *BaseModel) Update(entity interface{}, index string, uid interface{}) (int64, error) {
	result := global.GormDb.Model(entity).Where(index+" = ?", uid).Updates(entity)
	return result.RowsAffected, result.Error
}

func (u *BaseModel) UpdatePatch(entity interface{}, uid uint64) (int64, error) {
	result := global.GormDb.Model(entity).Where("id = ?", uid).Updates(entity)
	return result.RowsAffected, result.Error
}

func (u *BaseModel) Delete(entity interface{}, index string, uid interface{}) (int64, error) {
	result := global.GormDb.Model(entity).Where(index+" = ?", uid).Delete(entity)
	return result.RowsAffected, result.Error
}

func (u *BaseModel) FindAll(entity interface{}, opts ...DAOOption) error {
	tx := global.GormDb.Model(entity)
	if len(opts) > 0 {
		beginCustomTx(tx, opts)
	}

	tx.Find(entity)
	return tx.Error
}

func (u *BaseModel) FindByKeys(entity interface{}, keys map[string]interface{}) (int64, error) {
	result := global.GormDb.Model(entity).Where(keys).Find(entity)
	return result.RowsAffected, result.Error
}

func (u *BaseModel) FindByKeysNum(entity interface{}, keys map[string]interface{}) (int64, error) {
	result := global.GormDb.Model(entity).Where(keys).Find(entity)
	return result.RowsAffected, result.Error
}

func (u *BaseModel) FindByPages(entity interface{}, currentPage, pageSize int) error {
	result := Paginate(currentPage, pageSize).Model(entity).Find(entity)
	return result.Error
}

func (u *BaseModel) FindByPagesWithKeys(entity interface{},
	keys map[string]interface{},
	currentPage, pageSize int,
	opts ...DAOOption) error {
	tx := Paginate(currentPage, pageSize).Model(entity)
	beginCustomTx(tx, opts)
	result := tx.Where(keys).Find(entity)
	return result.Error
}

func (u *BaseModel) Count(entity interface{}, count *int64) error {
	result := global.GormDb.Model(entity).Count(count)
	return result.Error
}

func (u *BaseModel) CountWithKeys(entity interface{},
	count *int64,
	keys, keyOpts map[string]interface{},
	opts ...DAOOption) error {
	// tx begin
	tx := global.GormDb.Model(entity)
	beginCustomTx(tx, opts)
	searchCustomTx(tx, keys, keyOpts)
	tx.Count(count)
	return tx.Error
}

// search by pages
func (u *BaseModel) SearchByPagesWithKeys(entity interface{},
	keys, keyOpts map[string]interface{},
	currentPage, pageSize int,
	opts ...DAOOption) error {
	// tx begin
	tx := Paginate(currentPage, pageSize).Model(entity)
	beginCustomTx(tx, opts)
	searchCustomTx(tx, keys, keyOpts)
	tx.Find(entity)
	return tx.Error
}

// search custom tx
func searchCustomTx(tx *gorm.DB, keys, keyOpts map[string]interface{}) {
	compareOpt := map[string]string{">": ">", "<": "<", ">=": ">=", "<=": "<=", "=": "="}
	//set transaction
	transactionOpt := "or"
	if _, ok := keyOpts["searchKeyOpt"]; ok {
		transactionOpt = keyOpts["searchKeyOpt"].(string)
	}
	//search with like
	for k, v := range keys {
		usingLike := false
		// set tx
		if _, exist := keyOpts[k]; exist {
			tx.Where(k+keyOpts[k].(string)+" ?", fmt.Sprintf("%v", v))
			if _, ok := compareOpt[keyOpts[k].(string)]; ok {
				usingLike = false
			}
			if keyOpts[k].(string) == "like" {
				usingLike = true
			}
		}
		//using like
		if usingLike {
			if transactionOpt == "or" {
				tx.Where("").Or(k+" LIKE ?", "%"+v.(string)+"%")
			} else {
				tx.Where(k+" LIKE ?", "%"+v.(string)+"%")
			}
		}
	}
}

func Paginate(currentPage, pageSize int) *gorm.DB {
	offset := (currentPage - 1) * pageSize
	return global.GormDb.Offset(offset).Limit(pageSize)
}

// begin custom tx
func beginCustomTx(tx *gorm.DB, opts []DAOOption) {
	for _, opt := range opts {
		// set tx order
		if opt.Order != "" {
			tx.Order(opt.Order)
		}
		// set tx where
		if opt.Where != nil {
			tx.Where(opt.Where)
		}
		// set tx select
		if opt.Select != "" {
			tx.Select(opt.Select)
		}
		// set tx limit
		if opt.Limit > 0 {
			tx.Limit(int(opt.Limit))
		}
	}
}

func UpdateSqlString(object interface{}, vals []string, table, key string, value interface{}) (string, error) {
	sqlValString, err := UpdateSqlvals(object, vals)
	if err != nil {
		return "", err
	}
	switch value.(type) {
	case int, uint, uint32:
		return fmt.Sprintf("UPDATE %s SET  %s where %s = %d", table, sqlValString, key, value), nil
	default:
		return fmt.Sprintf("UPDATE %s SET  %s where %s = '%v'", table, sqlValString, key, value), nil
	}

}

func UpdateSqlvals(object interface{}, vals []string) (string, error) {
	// 找到要更新的字段,然后查看是否时空值
	typeOfObject := reflect.TypeOf(object)
	getValue := reflect.ValueOf(object)
	var sqlVals []string
	for _, v := range vals {
		getType, ok := typeOfObject.FieldByName(v)
		if !ok {
			return "", errors.New("提交的结构体field错误")
		}
		ovalue := getValue.FieldByName(v)

		switch getType.Type.Name() {
		case "string":
			valueS := ovalue.Interface().(string)
			if len(valueS) > 0 {
				sqlVals = append(sqlVals, fmt.Sprintf("`%s` = '%s'", getType.Tag.Get("json"), valueS))
			}
		case "int":
			valueI := ovalue.Interface().(int)
			if valueI != 0 {
				sqlVals = append(sqlVals, fmt.Sprintf("`%s` = %d", getType.Tag.Get("json"), valueI))
			}
		default:
			sqlVals = append(sqlVals, fmt.Sprintf("`%s` = %v", getType.Tag.Get("json"), ovalue.Interface()))
		}
	}
	if len(sqlVals) == 0 {
		return "", errors.New("所有提交参数为空")
	}

	return strings.Join(sqlVals, ", "), nil
}
