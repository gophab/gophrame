package i18n

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/gophab/gophrame/core/file"
	"github.com/gophab/gophrame/core/regex"
	"github.com/nacos-group/nacos-sdk-go/common/logger"
)

// pathType is the type for i18n file path.
type pathType string

const (
	pathTypeNone   pathType = "none"
	pathTypeNormal pathType = "normal"
)

// Manager for i18n contents, it is concurrent safe, supporting hot reload.
type Manager struct {
	mu       sync.RWMutex
	data     map[string]map[string]interface{} // Translating map.
	pattern  string                            // Pattern for regex parsing.
	pathType pathType                          // Path type for i18n files.
	options  Options                           // configuration options.
}

// Options is used for i18n object configuration.
type Options struct {
	Path       string   // I18n files storage path.
	Language   string   // Default local language.
	Delimiters []string // Delimiters for variable parsing.
}

var (
	// defaultLanguage defines the default language if user does not specify in options.
	defaultLanguage = "en"

	// defaultDelimiters defines the default key variable delimiters.
	defaultDelimiters = []string{"{#", "}"}

	// i18n files searching folders.
	searchFolders = []string{"language", "i18n"}
)

// New creates and returns a new i18n manager.
// The optional parameter `option` specifies the custom options for i18n manager.
// It uses a default one if it's not passed.
func New(options ...Options) *Manager {
	var opts Options
	var pathType = pathTypeNone
	if len(options) > 0 {
		opts = options[0]
		pathType = opts.checkPathType(opts.Path)
	} else {
		opts = Options{}
		for _, folder := range searchFolders {
			pathType = opts.checkPathType(folder)
			if pathType != pathTypeNone {
				break
			}
		}
	}
	if len(opts.Language) == 0 {
		opts.Language = defaultLanguage
	}
	if len(opts.Delimiters) == 0 {
		opts.Delimiters = defaultDelimiters
	}
	m := &Manager{
		options: opts,
		pattern: fmt.Sprintf(
			`%s(.+?)%s`,
			regexp.QuoteMeta(opts.Delimiters[0]),
			regexp.QuoteMeta(opts.Delimiters[1]),
		),
		pathType: pathType,
	}
	return m
}

// checkPathType checks and returns the path type for given directory path.
func (o *Options) checkPathType(dirPath string) pathType {
	if dirPath == "" {
		return pathTypeNone
	}

	if file.Exist(dirPath) {
		o.Path = dirPath
		return pathTypeNormal
	}

	return pathTypeNone
}

// SetPath sets the directory path storing i18n files.
func (m *Manager) SetPath(path string) error {
	pathType := m.options.checkPathType(path)
	if pathType == pathTypeNone {
		return fmt.Errorf(`%s does not exist`, path)
	}

	m.pathType = pathType
	// Reset the manager after path changed.
	m.reset()
	return nil
}

// SetLanguage sets the language for translator.
func (m *Manager) SetLanguage(language string) {
	m.options.Language = language
}

// SetDelimiters sets the delimiters for translator.
func (m *Manager) SetDelimiters(left, right string) {
	m.pattern = fmt.Sprintf(`%s(.+?)%s`, regexp.QuoteMeta(left), regexp.QuoteMeta(right))
}

// T is alias of Translate for convenience.
func (m *Manager) T(content string) string {
	return m.Translate(content)
}

// T is alias of Translate for convenience.
func (m *Manager) LT(locale, content string) string {
	return m.Translate(content)
}

// Tf is alias of TranslateFormat for convenience.
func (m *Manager) LTf(locale, format string, values ...interface{}) string {
	return m.TranslateFormat(format, values...)
}

// TranslateFormat translates, formats and returns the `format` with configured language
// and given `values`.
func (m *Manager) TranslateFormat(format string, values ...interface{}) string {
	return fmt.Sprintf(m.Translate(format), values...)
}

// TranslateFormat translates, formats and returns the `format` with configured language
// and given `values`.
func (m *Manager) LocaleTranslateFormat(locale, format string, values ...interface{}) string {
	return fmt.Sprintf(m.LocaleTranslate(locale, format), values...)
}

func getMapValue(config interface{}, path string) (interface{}, bool) {
	if config == nil {
		return nil, false
	}

	segs := strings.SplitN(path, ".", 2)
	if len(segs) == 2 {
		if node, b := getMapValue(config, segs[0]); b {
			return getMapValue(node, segs[1])
		}
	} else {
		switch reflect.TypeOf(config).Kind() {
		case reflect.Map:
			value := reflect.ValueOf(config).MapIndex(reflect.ValueOf(path))
			if value.IsValid() && !value.IsNil() && !value.IsZero() {
				return value.Interface(), true
			}
		case reflect.Array:
			if index, err := strconv.ParseInt(path, 10, 32); err == nil {
				array := reflect.ValueOf(config)
				if 0 <= index && index <= int64(array.Len()) {
					value := reflect.ValueOf(config).Index(int(index))
					if value.IsValid() && !value.IsNil() && !value.IsZero() {
						return value.Interface(), true
					}
				}
			}
		case reflect.Struct:
			value := reflect.ValueOf(config)
			field := value.FieldByName(path)
			if field.IsValid() && !field.IsNil() && !field.IsZero() {
				return field.Interface(), true
			}
		case reflect.Ptr, reflect.Interface:
			value := reflect.ValueOf(config).Elem()
			if value.IsValid() && !value.IsNil() && !value.IsZero() {
				return getMapValue(value.Interface(), path)
			}
		default:
		}
	}
	return nil, false
}

// Translate translates `content` with configured language.
func (m *Manager) Translate(content string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	transLang := m.options.Language
	if lang := GetCurrentLanguage(); lang != "" {
		transLang = lang
	}
	data := m.data[transLang]
	if data == nil {
		return content
	}
	// Parse content as name.
	if v, ok := getMapValue(data, content); ok {
		return fmt.Sprintf("%v", v)
	}
	// Parse content as variables container.
	result, _ := regex.ReplaceStringFuncMatch(
		m.pattern, content,
		func(match []string) string {
			if v, ok := data[match[1]]; ok {
				return fmt.Sprintf("%v", v)
			}
			// return match[1] will return the content between delimiters
			// return match[0] will return the original content
			return match[0]
		})
	return result
}

func (m *Manager) LocaleTranslate(lang, content string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	transLang := m.options.Language
	if lang != "" {
		transLang = lang
	}
	data := m.data[transLang]
	if data == nil {
		return content
	}
	// Parse content as name.
	if v, ok := getMapValue(data, content); ok {
		return fmt.Sprintf("%v", v)
	}
	// Parse content as variables container.
	result, _ := regex.ReplaceStringFuncMatch(
		m.pattern, content,
		func(match []string) string {
			if v, ok := data[match[1]]; ok {
				return fmt.Sprintf("%v", v)
			}
			// return match[1] will return the content between delimiters
			// return match[0] will return the original content
			return match[0]
		})
	return result
}

// GetContent retrieves and returns the configured content for given key and specified language.
// It returns an empty string if not found.
func (m *Manager) GetContent(key string) interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	transLang := m.options.Language
	if lang := GetCurrentLanguage(); lang != "" {
		transLang = lang
	}
	if data, ok := m.data[transLang]; ok {
		return data[key]
	}
	return ""
}

// reset reset data of the manager.
func (m *Manager) reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = nil
}

// init initializes the manager for lazy initialization design.
// The i18n manager is only initialized once.
func (m *Manager) init() {
	m.mu.RLock()
	// If the data is not nil, means it's already initialized.
	if m.data != nil {
		m.mu.RUnlock()
		return
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()
	switch m.pathType {
	case pathTypeNormal:
		files, _ := file.ScanDirFile(m.options.Path, "*.*", true)
		if len(files) == 0 {
			return
		}
		var (
			lang  string
			array []string
		)
		m.data = make(map[string]map[string]interface{})
		for _, fileName := range files {
			array = strings.Split(fileName, file.Separator)
			if len(array) > 1 {
				lang = file.Name(array[len(array)-1])
			} else if len(array) == 1 {
				lang = file.Name(array[0])
			}
			if m.data[lang] == nil {
				m.data[lang] = make(map[string]interface{})
			}
			var ext = m.data[lang]
			if err := json.Unmarshal(file.GetBytes(fileName), &ext); err != nil {
				logger.Errorf("load i18n file '%s' failed: %+v", fileName, err)
			} else {
				m.data[lang] = ext
			}
		}
	}
}

var i18nManager *Manager

func init() {
	i18nManager = New()
	i18nManager.init()
}

func T(content string) string {
	return i18nManager.Translate(content)
}

func Tf(format string, values ...interface{}) string {
	return i18nManager.TranslateFormat(format, values...)
}

func LT(locale, content string) string {
	return i18nManager.LocaleTranslate(locale, content)
}

func LTf(locale, format string, values ...interface{}) string {
	return i18nManager.LocaleTranslateFormat(locale, format, values...)
}
