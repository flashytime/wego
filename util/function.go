package util

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/godcong/wego/log"
	"github.com/satori/go.uuid"
)

/*CustomHeader xml header*/
const CustomHeader = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>`

/*CDATA xml cdata defines */
type CDATA struct {
	XMLName xml.Name
	Value   string `xml:",cdata"`
}

/* error types */
var (
	ErrorSignType  = errors.New("sign type error")
	ErrorParameter = errors.New("JsonApiParameters() check error")
	ErrorToken     = errors.New("EditAddressParameters() token is nil")
)

/*RandomKind RandomKind */
type RandomKind int

/*random kinds */
const (
	RandomNum      RandomKind = iota // 纯数字
	RandomLower                      // 小写字母
	RandomUpper                      // 大写字母
	RandomLowerNum                   // 数字、小写字母
	RandomUpperNum                   // 数字、大写字母
	RandomAll                        // 数字、大小写字母
)

/*RandomString defines */
var (
	RandomString = map[RandomKind]string{
		RandomNum:      "0123456789",
		RandomLower:    "abcdefghijklmnopqrstuvwxyz",
		RandomUpper:    "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		RandomLowerNum: "0123456789abcdefghijklmnopqrstuvwxyz",
		RandomUpperNum: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		RandomAll:      "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
	}
)

/*ParseNumber parse interface to number */
func ParseNumber(v interface{}) (float64, bool) {
	switch v0 := v.(type) {
	case float64:
		return v0, true
	case float32:
		return float64(v0), true
	}
	return 0, false
}

/*ParseInt parse interface to int64 */
func ParseInt(v interface{}) (int64, bool) {
	switch v0 := v.(type) {
	case int:
		return int64(v0), true
	case int32:
		return int64(v0), true
	case int64:
		return int64(v0), true
	case uint:
		return int64(v0), true
	case uint32:
		return int64(v0), true
	case uint64:
		return int64(v0), true
	case float64:
		return int64(v0), true
	case float32:
		return int64(v0), true
	default:
	}
	return 0, false
}

/*ParseString parse interface to string */
func ParseString(v interface{}) (string, bool) {
	switch v0 := v.(type) {
	case string:
		return v0, true
	case []byte:
		return string(v0), true
	case bytes.Buffer:
		return v0.String(), true
	default:
	}
	return "", false
}

/*MapToXML Convert MAP to XML */
func MapToXML(m Map) ([]byte, error) {
	return mapToXML(m, false)
}

func convertXML(k string, v interface{}, e *xml.Encoder, start xml.StartElement) error {
	//err := e.EncodeToken(start)
	//if err != nil {
	//	return err
	//}
	var err error
	switch v1 := v.(type) {
	case Map:
		return marshalXML(v1, e, xml.StartElement{Name: xml.Name{Local: k}})
	case map[string]interface{}:
		return marshalXML(v1, e, xml.StartElement{Name: xml.Name{Local: k}})
	case string:
		if _, err := strconv.ParseInt(v1, 10, 0); err != nil {
			err = e.EncodeElement(
				CDATA{Value: v1}, xml.StartElement{Name: xml.Name{Local: k}})
			return err
		}
		err = e.EncodeElement(v1, xml.StartElement{Name: xml.Name{Local: k}})
		return err
	case float64:
		if v1 == float64(int64(v1)) {
			err = e.EncodeElement(int64(v1), xml.StartElement{Name: xml.Name{Local: k}})
			return err
		}
		err = e.EncodeElement(v1, xml.StartElement{Name: xml.Name{Local: k}})
		return err
	case bool:
		err = e.EncodeElement(v1, xml.StartElement{Name: xml.Name{Local: k}})
		return err
	case []interface{}:
		size := len(v1)
		for i := 0; i < size; i++ {
			err := convertXML(k, v1[i], e, xml.StartElement{Name: xml.Name{Local: k}})
			if err != nil {
				return err
			}
		}
		if len(v1) == 1 {
			return convertXML(k, "", e, xml.StartElement{Name: xml.Name{Local: k}})
		}

	default:
		//convertXML(k, v1, e, xml.StartElement{Name: xml.Name{Local: k}})
		log.Error(v1)
	}
	//if len(v) == 1 {
	//	convertXML(k, "dummy", e, xml.StartElement{Name: xml.Name{Local: k}})
	//}
	//return e.EncodeToken(start.End())
	return nil
}

func marshalXML(maps Map, e *xml.Encoder, start xml.StartElement) error {
	if maps == nil {
		return errors.New("nil map")
	}
	err := e.EncodeToken(start)
	if err != nil {
		return err
	}
	for k, v := range maps {
		err := convertXML(k, v, e, xml.StartElement{Name: xml.Name{Local: k}})
		if err != nil {
			return err
		}
	}
	return e.EncodeToken(start.End())
}

func mapToXML(maps Map, needHeader bool) ([]byte, error) {

	buff := bytes.NewBuffer([]byte(CustomHeader))
	if needHeader {
		buff.Write([]byte(xml.Header))
	}

	enc := xml.NewEncoder(buff)
	err := marshalXML(maps, enc, xml.StartElement{Name: xml.Name{Local: "xml"}})
	if err != nil {
		return nil, err
	}
	err = enc.Flush()
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

/*XMLToMap Convert XML to MAP */
func XMLToMap(contentXML []byte) Map {
	m, err := xmlToMap(contentXML, false)
	if err != nil {
		return nil
	}
	return m
}

/*JSONToMap Convert JSON to MAP */
func JSONToMap(content []byte) Map {
	m := Map{}
	err := json.Unmarshal(content, &m)
	if err != nil {
		log.Error(err)
	}
	return m
}

func unmarshalXML(maps Map, d *xml.Decoder, start xml.StartElement, needCast bool) error {
	//m := make(Map)
	current := ""
	var data interface{}
	//var err error
	last := ""
	arrayTmp := make(Map)
	arrayTag := ""
	var ele []string

	for t, err := d.Token(); err == nil; t, err = d.Token() {
		switch token := t.(type) {
		// 处理元素开始（标签）
		case xml.StartElement:
			if strings.ToLower(token.Name.Local) == "xml" ||
				strings.ToLower(token.Name.Local) == "root" {
				continue
			}
			ele = append(ele, token.Name.Local)
			current = strings.Join(ele, ".")
			log.Debug("EndElement", current)
			log.Debug("EndElement", last)
			log.Debug("EndElement", arrayTag)
			if current == last {
				arrayTag = current
				tmp := maps.Get(arrayTag)
				switch tmp.(type) {
				case []interface{}:
					arrayTmp.Set(arrayTag, tmp)
				default:
					arrayTmp.Set(arrayTag, []interface{}{tmp})
				}
				maps.Delete(arrayTag)
			}
			log.Debug("StartElement", ele)
			// 处理元素结束（标签）
		case xml.EndElement:
			name := token.Name.Local
			// fmt.Printf("This is the end: %s\n", name)
			if strings.ToLower(name) == "xml" ||
				strings.ToLower(name) == "root" {
				break
			}
			last = strings.Join(ele, ".")
			log.Debug("EndElement", current)
			log.Debug("EndElement", last)
			log.Debug("EndElement", arrayTag)

			if current == last {
				if data != nil {
					maps.Set(current, data)
				} else {
					//m.Set(current, nil)
				}
				data = nil
			}
			if last == arrayTag {
				arr := arrayTmp.GetArray(arrayTag)
				if arr != nil {
					if v := maps.Get(arrayTag); v != nil {
						maps.Set(arrayTag, append(arr, v))
					} else {
						maps.Set(arrayTag, arr)
					}
				} else {
					//exception doing
					maps.Set(arrayTag, []interface{}{maps.Get(arrayTag)})
				}
				arrayTmp.Delete(arrayTag)
				arrayTag = ""
			}

			ele = ele[:len(ele)-1]
			log.Debug("EndElement", ele)
			// 处理字符数据（这里就是元素的文本）
		case xml.CharData:
			if needCast {
				data, err = strconv.Atoi(string(token))
				if err == nil {
					continue
				}

				data, err = strconv.ParseFloat(string(token), 64)
				if err == nil {
					continue
				}

				data, err = strconv.ParseBool(string(token))
				if err == nil {
					continue
				}
			}

			data = string(token)
			log.Debug("CharData", data)
			// 异常处理(Log输出）
		default:
			log.Debug(token)
		}

	}

	return nil
}

func xmlToMap(contentXML []byte, hasHeader bool) (Map, error) {
	m := make(Map)
	dec := xml.NewDecoder(bytes.NewReader(contentXML))
	err := unmarshalXML(m, dec, xml.StartElement{Name: xml.Name{Local: "xml"}}, true)
	if err != nil {
		return nil, err
	}

	return m, nil
}

/*Time get time string */
func Time(t ...time.Time) string {
	if t == nil {
		return strconv.Itoa(time.Now().Nanosecond())
	}
	return strconv.Itoa(t[0].Nanosecond())
}

/*GenerateNonceStr GenerateNonceStr */
func GenerateNonceStr() string {
	return GenerateUUID()
}

/*GenerateUUID GenerateUUID */
func GenerateUUID() string {
	u1, _ := uuid.NewV1()
	s := u1.String()
	s = strings.Replace(s, "-", "", -1)
	run := ([]rune)(s)[:32]
	return string(run)
}

/*In check v is in source */
func In(source []string, v string) bool {
	size := len(source)
	for i := 0; i < size; i++ {
		if source[i] == v {
			return true
		}
	}

	return false
}

/*MapToString MapToString */
func MapToString(data Map, ignore []string) string {
	var sign []string
	m := data.Expect(ignore)
	keys := m.SortKeys()
	size := len(keys)

	for i := 0; i < size; i++ {
		v := strings.TrimSpace(m.GetString(keys[i]))
		if len(v) > 0 {
			sign = append(sign, strings.Join([]string{keys[i], v}, "="))
		}
	}

	log.Debug(strings.Join(sign, "&"))
	return strings.Join(sign, "&")
}

/*ToURLParams map to url params */
func ToURLParams(data Map, ignore []string) string {
	var sign []string
	m := data.Expect(ignore)
	keys := m.SortKeys()
	size := len(keys)
	for i := 0; i < size; i++ {
		v := strings.TrimSpace(m.GetString(keys[i]))
		if len(v) > 0 {
			sign = append(sign, strings.Join([]string{keys[i], v}, "="))
		}
	}

	return strings.Join(sign, "&")
}

// CurrentTimeStampMS get current time with millisecond
func CurrentTimeStampMS() int64 {
	return time.Now().UnixNano() / time.Millisecond.Nanoseconds()
}

// CurrentTimeStampNS get current time with nanoseconds
func CurrentTimeStampNS() int64 {
	return time.Now().UnixNano()
}

// CurrentTimeStamp get current time with unix
func CurrentTimeStamp() int64 {
	return time.Now().Unix()
}

// CurrentTimeStampString get current time to string
func CurrentTimeStampString() string {
	return strconv.FormatInt(CurrentTimeStamp(), 10)
}

// SHA1 transfer string to sha1
func SHA1(s string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(s)))
}

// HmacSha256 ...
func HmacSha256(data []byte, key string) string {
	m := hmac.New(sha256.New, []byte(key))
	m.Write(data)
	return strings.ToUpper(fmt.Sprintf("%x", m.Sum(nil)))
}

func signatureSHA1(m Map) string {
	keys := m.SortKeys()
	var sign []string
	size := len(keys)
	for i := 0; i < size; i++ {
		if v := strings.TrimSpace(m.GetString(keys[i])); v != "" {
			log.Debug(keys[i], v)
			sign = append(sign, strings.Join([]string{keys[i], v}, "="))
		} else if v, b := m.GetInt64(keys[i]); b {
			log.Debug(keys[i], v)
			sign = append(sign, strings.Join([]string{keys[i], strconv.FormatInt(v, 10)}, "="))
		}
	}

	sb := strings.Join(sign, "&")
	return SHA1(sb)
}

//GenerateRandomString2 随机字符串
func GenerateRandomString2(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{{10, 48}, {26, 97}, {26, 65}}, make([]byte, size)
	isAll := kind > 2 || kind < 0

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if isAll { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}

//GenerateRandomString 随机字符串
func GenerateRandomString(size int, kind ...RandomKind) string {
	bytes := RandomString[RandomAll]
	if kind != nil {
		if k, b := RandomString[kind[0]]; b == true {
			bytes = k
		}
	}
	var result []byte
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}
