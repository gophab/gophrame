package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/wjshen/gophrame/core/global"
	"github.com/wjshen/gophrame/core/logger"
	"github.com/wjshen/gophrame/core/snowflake"

	"reflect"
	"strings"
	"time"
)

const (
	ErrorsGormDBCreateParamsNotPtr string = "gorm Create 函数的参数必须是一个指针, 为了完美支持 gorm 的所有回调函数,请在参数前面添加 & "
	ErrorsGormDBUpdateParamsNotPtr string = "gorm 的 Update、Save 函数的参数必须是一个指针(GinSkeleton ≥ v1.5.29 版本新增验证，为了完美支持 gorm 的所有回调函数,请在参数前面添加 & )"
)

// 这里的函数都是gorm的hook函数，拦截一些官方我们认为不合格的操作行为，提升项目整体的完美性

// MaskNotDataError 解决gorm v2 包在查询无数据时，报错问题（record not found），但是官方认为报错是应该是，我们认为查询无数据，代码一切ok，不应该报错
func MaskNotDataError(gormDB *gorm.DB) {
	gormDB.Statement.RaiseErrorOnNotFound = false
}

func UpdateCreatedTimeHook(db *gorm.DB) {
	ctx := db.Statement.Context

	timeFieldsToInit := []string{"CreatedTime"}
	for _, field := range timeFieldsToInit {
		if timeField := db.Statement.Schema.LookUpField(field); timeField != nil {
			switch db.Statement.ReflectValue.Kind() {
			case reflect.Slice, reflect.Array:
				for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
					if _, isZero := timeField.ValueOf(ctx, db.Statement.ReflectValue.Index(i)); isZero {
						timeField.Set(ctx, db.Statement.ReflectValue.Index(i), time.Now())
					}
				}
			case reflect.Struct:
				if _, isZero := timeField.ValueOf(ctx, db.Statement.ReflectValue); isZero {
					timeField.Set(ctx, db.Statement.ReflectValue, time.Now())
				}
			}
		}
	}
}

func UpdateLastModifiedTimeHook(db *gorm.DB) {
	ctx := db.Statement.Context

	timeFieldsToInit := []string{"LastModifiedTime"}
	for _, field := range timeFieldsToInit {
		if timeField := db.Statement.Schema.LookUpField(field); timeField != nil {
			switch db.Statement.ReflectValue.Kind() {
			case reflect.Slice, reflect.Array:
				for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
					if _, isZero := timeField.ValueOf(ctx, db.Statement.ReflectValue.Index(i)); isZero {
						timeField.Set(ctx, db.Statement.ReflectValue.Index(i), time.Now())
					}
				}
			case reflect.Struct:
				if _, isZero := timeField.ValueOf(ctx, db.Statement.ReflectValue); isZero {
					timeField.Set(ctx, db.Statement.ReflectValue, time.Now())
				}
			}
		}
	}
}

func UpdateDeletedTimeHook(db *gorm.DB) {
	ctx := db.Statement.Context

	timeFieldsToInit := []string{"DeletedTime"}
	for _, field := range timeFieldsToInit {
		if timeField := db.Statement.Schema.LookUpField(field); timeField != nil {
			switch db.Statement.ReflectValue.Kind() {
			case reflect.Slice, reflect.Array:
				for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
					if _, isZero := timeField.ValueOf(ctx, db.Statement.ReflectValue.Index(i)); isZero {
						timeField.Set(ctx, db.Statement.ReflectValue.Index(i), time.Now())
					}
				}
			case reflect.Struct:
				if _, isZero := timeField.ValueOf(ctx, db.Statement.ReflectValue); isZero {
					timeField.Set(ctx, db.Statement.ReflectValue, time.Now())
				}
			}
		}
	}
}

func UpdateIdHook(db *gorm.DB) {
	ctx := db.Statement.Context

	idFieldsToInit := []string{"Id"}
	for _, field := range idFieldsToInit {
		if timeField := db.Statement.Schema.LookUpField(field); timeField != nil {
			idGenerator := func() (result interface{}) { return }
			tags := timeField.Tag.Get("gorm")
			if strings.Contains(tags, "default:uuid") {
				idGenerator = func() interface{} { return uuid.NewString() }
			}
			if strings.Contains(tags, "default:snowflake") {
				idGenerator = func() interface{} { return snowflake.SnowflakeIdGenerator().GetId() }
			}

			switch db.Statement.ReflectValue.Kind() {
			case reflect.Slice, reflect.Array:
				for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
					if _, isZero := timeField.ValueOf(ctx, db.Statement.ReflectValue.Index(i)); isZero {
						timeField.Set(ctx, db.Statement.ReflectValue.Index(i), idGenerator())
					}
				}
			case reflect.Struct:
				if _, isZero := timeField.ValueOf(ctx, db.Statement.ReflectValue); isZero {
					timeField.Set(ctx, db.Statement.ReflectValue, idGenerator())
				}
			}
		}
	}
}

// InterceptCreatePramsNotPtrError 拦截 create 函数参数如果是非指针类型的错误,新用户最容犯此错误
func CreateBeforeHook(gormDB *gorm.DB) {
	if reflect.TypeOf(gormDB.Statement.Dest).Kind() != reflect.Ptr {
		logger.Warn(ErrorsGormDBCreateParamsNotPtr)
	} else {
		destValueOf := reflect.ValueOf(gormDB.Statement.Dest).Elem()
		if destValueOf.Type().Kind() == reflect.Slice || destValueOf.Type().Kind() == reflect.Array {
			inLen := destValueOf.Len()
			for i := 0; i < inLen; i++ {
				row := destValueOf.Index(i)
				if row.Type().Kind() == reflect.Struct {
					if b, column := structHasSpecialField("Id", row); b {
						field := row.FieldByName(column)
						switch field.Kind() {
						case reflect.String:
							if field.IsZero() {
								field.Set(reflect.ValueOf(uuid.NewString()))
							}
						case reflect.Int:
						case reflect.Int8:
						case reflect.Int16:
						case reflect.Int32:
						case reflect.Int64:
							if field.IsZero() {
								field.Set(reflect.ValueOf(snowflake.SnowflakeIdGenerator().GetId()))
							}
						}
					}
					if b, column := structHasSpecialField("CreatedTime", row); b {
						row.FieldByName(column).Set(reflect.ValueOf(time.Now().Format(global.DateFormat)))
					}
					if b, column := structHasSpecialField("ModifiedTime", row); b {
						row.FieldByName(column).Set(reflect.ValueOf(time.Now().Format(global.DateFormat)))
					}

				} else if row.Type().Kind() == reflect.Map {
					if b, column := structHasSpecialField("id", row); !b || row.MapIndex(reflect.ValueOf(column)).IsZero() {
						row.SetMapIndex(reflect.ValueOf(column), reflect.ValueOf(snowflake.SnowflakeIdGenerator().GetId()))
					}
					if b, column := structHasSpecialField("created_time", row); b {
						row.SetMapIndex(reflect.ValueOf(column), reflect.ValueOf(time.Now().Format(global.DateFormat)))
					}
					if b, column := structHasSpecialField("modified_time", row); b {
						row.SetMapIndex(reflect.ValueOf(column), reflect.ValueOf(time.Now().Format(global.DateFormat)))
					}
				}
			}
		} else if destValueOf.Type().Kind() == reflect.Struct {
			//  if destValueOf.Type().Kind() == reflect.Struct
			if b, column := structHasSpecialField("ID", gormDB.Statement.Dest); b {
				if id, _ := gormDB.Statement.Get("ID"); id == nil {
					gormDB.Statement.SetColumn(column, time.Now().Format(global.DateFormat))
				}
			}
			// 参数校验无错误自动设置 CreatedAt、 UpdatedAt
			if b, column := structHasSpecialField("CreatedTime", gormDB.Statement.Dest); b {
				gormDB.Statement.SetColumn(column, time.Now().Format(global.DateFormat))
			}
			if b, column := structHasSpecialField("ModifiedTime", gormDB.Statement.Dest); b {
				gormDB.Statement.SetColumn(column, time.Now().Format(global.DateFormat))
			}
		} else if destValueOf.Type().Kind() == reflect.Map {
			if b, column := structHasSpecialField("id", gormDB.Statement.Dest); !b || destValueOf.MapIndex(reflect.ValueOf(column)).IsZero() {
				destValueOf.SetMapIndex(reflect.ValueOf(column), reflect.ValueOf(snowflake.SnowflakeIdGenerator().GetId()))
			}
			if b, column := structHasSpecialField("created_time", gormDB.Statement.Dest); b {
				destValueOf.SetMapIndex(reflect.ValueOf(column), reflect.ValueOf(time.Now().Format(global.DateFormat)))
			}
			if b, column := structHasSpecialField("modified_time", gormDB.Statement.Dest); b {
				destValueOf.SetMapIndex(reflect.ValueOf(column), reflect.ValueOf(time.Now().Format(global.DateFormat)))
			}
		}
	}
}

// UpdateBeforeHook
// InterceptUpdatePramsNotPtrError 拦截 save、update 函数参数如果是非指针类型的错误
// 对于开发者来说，以结构体形式更新数，只需要在 update 、save 函数的参数前面添加 & 即可
// 最终就可以完美兼支持、兼容 gorm 的所有回调函数
// 但是如果是指定字段更新，例如： UpdateColumn 函数则只传递值即可，不需要做校验
func UpdateBeforeHook(gormDB *gorm.DB) {
	if reflect.TypeOf(gormDB.Statement.Dest).Kind() == reflect.Struct {
		logger.Warn(ErrorsGormDBUpdateParamsNotPtr)
	} else if reflect.TypeOf(gormDB.Statement.Dest).Kind() == reflect.Map {
		// 如果是调用了 gorm.Update 、updates 函数 , 在参数没有传递指针的情况下，无法触发回调函数

	} else if reflect.TypeOf(gormDB.Statement.Dest).Kind() == reflect.Ptr && reflect.ValueOf(gormDB.Statement.Dest).Elem().Kind() == reflect.Struct {
		// 参数校验无错误自动设置 UpdatedAt
		if b, column := structHasSpecialField("ModifiedTime", gormDB.Statement.Dest); b {
			gormDB.Statement.SetColumn(column, time.Now().Format(global.DateFormat))
		}
	} else if reflect.TypeOf(gormDB.Statement.Dest).Kind() == reflect.Ptr && reflect.ValueOf(gormDB.Statement.Dest).Elem().Kind() == reflect.Map {
		if b, column := structHasSpecialField("modified_time", gormDB.Statement.Dest); b {
			destValueOf := reflect.ValueOf(gormDB.Statement.Dest).Elem()
			destValueOf.SetMapIndex(reflect.ValueOf(column), reflect.ValueOf(time.Now().Format(global.DateFormat)))
		}
	}
}

// structHasSpecialField  检查结构体是否有特定字段
func structHasSpecialField(fieldName string, anyStructPtr interface{}) (bool, string) {
	var tmp reflect.Type
	if reflect.TypeOf(anyStructPtr).Kind() == reflect.Ptr && reflect.ValueOf(anyStructPtr).Elem().Kind() == reflect.Map {
		destValueOf := reflect.ValueOf(anyStructPtr).Elem()
		for _, item := range destValueOf.MapKeys() {
			if item.String() == fieldName {
				return true, fieldName
			}
		}
	} else if reflect.TypeOf(anyStructPtr).Kind() == reflect.Ptr && reflect.ValueOf(anyStructPtr).Elem().Kind() == reflect.Struct {
		destValueOf := reflect.ValueOf(anyStructPtr).Elem()
		tf := destValueOf.Type()
		for i := 0; i < tf.NumField(); i++ {
			if !tf.Field(i).Anonymous && tf.Field(i).Type.Kind() != reflect.Struct {
				if tf.Field(i).Name == fieldName {
					return true, getColumnNameFromGormTag(fieldName, tf.Field(i).Tag.Get("gorm"))
				}
			} else if tf.Field(i).Type.Kind() == reflect.Struct {
				tmp = tf.Field(i).Type
				for j := 0; j < tmp.NumField(); j++ {
					if tmp.Field(j).Name == fieldName {
						return true, getColumnNameFromGormTag(fieldName, tmp.Field(j).Tag.Get("gorm"))
					}
				}
			}
		}
	} else if reflect.Indirect(anyStructPtr.(reflect.Value)).Type().Kind() == reflect.Struct {
		// 处理结构体
		destValueOf := anyStructPtr.(reflect.Value)
		tf := destValueOf.Type()
		for i := 0; i < tf.NumField(); i++ {
			if !tf.Field(i).Anonymous && tf.Field(i).Type.Kind() != reflect.Struct {
				if tf.Field(i).Name == fieldName {
					return true, getColumnNameFromGormTag(fieldName, tf.Field(i).Tag.Get("gorm"))
				}
			} else if tf.Field(i).Type.Kind() == reflect.Struct {
				tmp = tf.Field(i).Type
				for j := 0; j < tmp.NumField(); j++ {
					if tmp.Field(j).Name == fieldName {
						return true, getColumnNameFromGormTag(fieldName, tmp.Field(j).Tag.Get("gorm"))
					}
				}
			}
		}
	} else if reflect.Indirect(anyStructPtr.(reflect.Value)).Type().Kind() == reflect.Map {
		destValueOf := anyStructPtr.(reflect.Value)
		for _, item := range destValueOf.MapKeys() {
			if item.String() == fieldName {
				return true, fieldName
			}
		}
	}
	return false, ""
}

// getColumnNameFromGormTag 从 gorm 标签中获取字段名
// @defaultColumn 如果没有 gorm：column 标签为字段重命名，则使用默认字段名
// @TagValue 字段中含有的gorm："column:created_at" 标签值，可能的格式：1. column:created_at    、2. default:null;  column:created_at  、3.  column:created_at; default:null
func getColumnNameFromGormTag(defaultColumn, TagValue string) (str string) {
	pos1 := strings.Index(TagValue, "column:")
	if pos1 == -1 {
		str = defaultColumn
		return
	} else {
		TagValue = TagValue[pos1+7:]
	}
	pos2 := strings.Index(TagValue, ";")
	if pos2 == -1 {
		str = TagValue
	} else {
		str = TagValue[:pos2]
	}
	return strings.ReplaceAll(str, " ", "")
}
