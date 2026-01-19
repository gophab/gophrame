package sensitive

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/gophab/gophrame/core/global"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/starter"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func init() {
	starter.RegisterStarter(Start)
}

func Start() {
	if global.DB != nil {
		global.DB.Callback().Create().After("gorm:create").Register("SensitiveUpdateHook", SensitiveUpdateHook)
		global.DB.Callback().Update().After("gorm:update").Register("SensitiveUpdateHook", SensitiveUpdateHook)
		global.DB.Callback().Query().After("gorm:query").Register("SensitiveLoadHook", SensitiveLoadHook)
	}
}

func getTagSection(tag, key string) string {
	segs := strings.Split(tag, ";")
	for i := range segs {
		if after, ok := strings.CutPrefix(segs[i], key+":"); ok {
			return after
		}
	}
	return ""
}

func sensitiveTag(field *schema.Field) (exist bool, target string, mode string) {
	if tag, b := field.Tag.Lookup("sensitive"); b {
		target = getTagSection(tag, "target")
		mode = getTagSection(tag, "mode")
		if mode == "" {
			mode = "encrypt"
		}
		if target == "" {
			target = field.Name
		}
		exist = true
	}
	return
}

func sensitiveTag2(field reflect.StructField) (exist bool, target string, mode string) {
	if tag, b := field.Tag.Lookup("sensitive"); b {
		target = getTagSection(tag, "target")
		mode = getTagSection(tag, "mode")
		if mode == "" {
			mode = "encrypt"
		}
		if target == "" {
			target = field.Name
		}

		exist = true
	}
	return
}

type Encrypter struct {
	key string
}

func (e *Encrypter) Encrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(e.key)) // key的长度必须是16, 24或32字节以匹配AES-128, AES-192或AES-256
	if err != nil {
		return "", err
	}

	plainText := []byte(text)
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return hex.EncodeToString(cipherText), nil
}

func (e *Encrypter) Decrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(e.key))
	if err != nil {
		return "", err
	}

	encryptedBytes, err := hex.DecodeString(text)
	if err != nil {
		return "", err
	}

	if len(encryptedBytes) < aes.BlockSize {
		return "", err
	}
	iv := encryptedBytes[:aes.BlockSize]
	cipherText := encryptedBytes[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}

func base64Decode(v string) string {
	if results, err := base64.StdEncoding.DecodeString(v); err == nil {
		return string(results)
	} else {
		logger.Error("Base64 decode error: ", err.Error())
	}
	return ""
}

var encrypter = &Encrypter{
	key: "x3J9sLkqVzTbYwF2PmHnQhWcZ7UoKpN4RtGvB6yEdCfA"[0:32],
}

func EncodeSensitiveValue(v any, mode string) any {
	enc := ""
	switch v := v.(type) {
	case string:
		enc = v
	case *string:
		enc = *v
	default:
		enc = fmt.Sprint(v)
	}

	switch mode {
	case "encrypt": /* 加密 */
		if result, err := encrypter.Encrypt(enc); err == nil {
			return result
		} else {
			logger.Error("Encrypt error: ", err.Error())
		}
	case "mask": /* 掩码 */
	}
	return nil
}

func DecodeSensitiveValue(v any, mode string) any {
	enc := ""
	switch v := v.(type) {
	case string:
		enc = v
	case *string:
		enc = *v
	default:
		enc = fmt.Sprintf("%v", v)
	}

	var result any

	switch mode {
	case "encrypt": /* 解密 */
		if res, err := encrypter.Decrypt(enc); err == nil {
			result = res
		}
	case "mask": /* 掩码 */
		if len(enc) > 10 {
			result = enc[:3] + "****" + enc[len(enc)-4:]
		} else if len(enc) > 5 {
			result = "****" + enc[len(enc)-4:]
		} else if len(enc) > 2 {
			result = "**" + enc[len(enc)-1:]
		}
	}
	return result
}

func getFieldValue(db *gorm.DB, item reflect.Value, field *schema.Field) any {
	ctx := db.Statement.Context
	if v, isZero := field.ValueOf(ctx, item); !isZero {
		vt := reflect.ValueOf(v).Kind()
		switch vt {
		case reflect.Array, reflect.Slice, reflect.Map, reflect.Struct:
			if bs, err := json.Marshal(v); err == nil {
				return string(bs)
			}
		case reflect.Pointer:
			return reflect.ValueOf(v).Elem()
		default:
			return v
		}
	}

	return nil
}

// field = value
func setFieldValue(data reflect.Value, fieldName string, value any) {
	if value != nil {
		if dataField := data.FieldByName(fieldName); dataField.IsValid() {
			if dataField.Kind() == reflect.Ptr {
				z := reflect.New(dataField.Type().Elem())
				z.Elem().Set(reflect.ValueOf(value))
				dataField.Set(z)
			} else if !dataField.IsZero() {
				dataField.Set(reflect.ValueOf(value))
			}
		}
	}
}

// data[field] = value
func setSchemaFieldValue(ctx context.Context, data reflect.Value, field *schema.Field, value any) {
	if value != nil {
		field.Set(ctx, data, value)
	}
}

func buildSensitiveField(db *gorm.DB, item reflect.Value, source string, targetField *schema.Field, mode string) {
	var sourceField *schema.Field = targetField
	if source != "" {
		sourceField = db.Statement.Schema.LookUpField(source)
	}

	if sourceField == nil {
		return
	}

	if v := getFieldValue(db, item, sourceField); v != nil {
		setFieldValue(item, targetField.Name, EncodeSensitiveValue(v, mode))
	}
}

func loadSensitiveField(db *gorm.DB, item reflect.Value, sourceField *schema.Field, target string, mode string) {
	var targetField *schema.Field = sourceField
	if target != "" {
		targetField = db.Statement.Schema.LookUpField(target)
	}

	if targetField == nil {
		return
	}

	if v := getFieldValue(db, item, sourceField); v != nil {
		setFieldValue(item, target, DecodeSensitiveValue(v, mode))
	}
}

func loadSchemaSensitiveField(db *gorm.DB, item reflect.Value, sourceField *schema.Field, target string, mode string) {
	ctx := db.Statement.Context

	var targetField *schema.Field = sourceField
	if target != "" {
		targetField = db.Statement.Schema.LookUpField(target)
	}

	if targetField == nil {
		return
	}

	if v := getFieldValue(db, item, sourceField); v != nil {
		setSchemaFieldValue(ctx, item, targetField, DecodeSensitiveValue(v, mode))
	}
}

// 1. store sensitive fields
func SensitiveUpdateHook(db *gorm.DB) {
	if db.Statement.Schema == nil {
		return
	}

	if len(db.Statement.Schema.Fields) == 0 {
		return
	}

	for _, field := range db.Statement.Schema.Fields {
		if b, target, mode := sensitiveTag(field); b {
			// 1. field 是基本数据类型
			switch db.Statement.ReflectValue.Kind() {
			case reflect.Slice, reflect.Array:
				for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
					item := db.Statement.ReflectValue.Index(i)
					buildSensitiveField(db, item, target, field, mode)
				}
			case reflect.Struct:
				buildSensitiveField(db, db.Statement.ReflectValue, target, field, mode)
			}
		}
	}
}

// 2. get locale fields
func SensitiveLoadHook(db *gorm.DB) {
	if db.Statement.Schema == nil {
		return
	}

	if len(db.Statement.Schema.Fields) == 0 {
		return
	}

	var sensitiveFields = make([]*schema.Field, 0)
	for _, field := range db.Statement.Schema.Fields {
		if _, b := field.Tag.Lookup("sensitive"); b {
			sensitiveFields = append(sensitiveFields, field)
		}
	}

	if len(sensitiveFields) > 0 {
		switch db.Statement.ReflectValue.Kind() {
		case reflect.Slice, reflect.Array:
			for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
				item := db.Statement.ReflectValue.Index(i)
				var v = reflect.Indirect(item)
				if v.Kind() == reflect.Struct {
					if v.Type() == db.Statement.Schema.ModelType {
						for _, field := range sensitiveFields {
							if b, target, mode := sensitiveTag(field); b {
								loadSchemaSensitiveField(db, item, field, target, mode)
							}
						}
					} else {
						for _, field := range sensitiveFields {
							if b, target, mode := sensitiveTag(field); b {
								loadSensitiveField(db, v, field, target, mode)
							}
						}
					}
				}
			}
		case reflect.Struct:
			var v = reflect.Indirect(db.Statement.ReflectValue)
			if v.Type() == db.Statement.Schema.ModelType {
				for _, field := range sensitiveFields {
					if b, target, mode := sensitiveTag(field); b {
						loadSchemaSensitiveField(db, db.Statement.ReflectValue, field, target, mode)
					}
				}
			} else {
				for _, field := range sensitiveFields {
					if b, target, mode := sensitiveTag(field); b {
						loadSensitiveField(db, v, field, target, mode)
					}
				}
			}
		}
	}
}

func Translate(v any, enc bool) any {
	var vt = reflect.TypeOf(v)
	var vv = reflect.ValueOf(v)
	if vt.Kind() == reflect.Ptr {
		vt = vt.Elem()
		vv = vv.Elem()
	}

	switch vt.Kind() {
	case reflect.Ptr:
		value := Translate(vv.Elem().Interface(), enc)
		return &value
	case reflect.Array, reflect.Slice:
		results := make([]any, 0)
		for i := range vv.Len() {
			value := Translate(vv.Index(i).Interface(), enc)
			results = append(results, value)
		}
		return results
	case reflect.Map:
		result := make(map[string]any)
		for _, k := range vv.MapKeys() {
			value := Translate(vv.MapIndex(k).Interface(), enc)
			result[k.String()] = value
		}
		return result
	case reflect.Struct:
		for k := range vt.NumField() {
			field := vt.Field(k)
			if b, target, mode := sensitiveTag2(field); b {
				if enc {
					if target != "" && target != "-" {
						targetField := vv.FieldByName(target)
						if targetField.IsValid() && targetField.Interface() != nil {
							setFieldValue(vv, field.Name, EncodeSensitiveValue(vv.FieldByName(target).Interface(), mode))
						}
					}
				} else {
					sourceField := vv.FieldByName(field.Name)
					if sourceField.IsValid() && sourceField.Interface() != nil {
						setFieldValue(vv, target, DecodeSensitiveValue(vv.FieldByName(field.Name).Interface(), mode))
					}
				}
			}
		}
	default:
	}
	return v
}
