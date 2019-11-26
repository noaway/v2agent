package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// GetMD5HashString func
func GetMD5HashString(str string) string {
	return GetMD5HashBytes([]byte(str))
}

// GetMD5HashBytes func
func GetMD5HashBytes(data []byte) string {
	hasher := md5.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}

// GetRandomString func
func GetRandomString(n int, alphabets ...byte) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	for i, b := range bytes {
		if len(alphabets) == 0 {
			bytes[i] = alphanum[b%byte(len(alphanum))]
		} else {
			bytes[i] = alphabets[b%byte(len(alphabets))]
		}
	}
	return string(bytes)
}

// Errors func
func Errors(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// WaitTimeout fn
func WaitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false
	case <-time.After(timeout):
		return true
	}
}

// DeferError fn
func DeferError(errorfn func(string, interface{}), dones ...func()) {
	if err := recover(); err != nil {
		var buf [2 << 10]byte
		errorfn(string(buf[:runtime.Stack(buf[:], false)]), err)
	}
	for _, done := range dones {
		done()
	}
}

// ParseIPAndPort func
func ParseIPAndPort(addr string) (string, int, error) {
	host, strPort, err := net.SplitHostPort(addr)
	if err != nil {
		return "", 0, err
	}

	port, err := StrTo(strPort).Int()
	if err != nil {
		return "", 0, err
	}

	return host, port, nil
}

func EnsurePath(path string, dir bool) error {
	if !dir {
		path = filepath.Dir(path)
	}
	return os.MkdirAll(path, 0755)
}

// Append func
func Append(slice []string, data ...string) []string {
	l := len(slice)
	if l+len(data) > cap(slice) {
		newSlice := make([]string, (l+len(data))*2)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0 : l+len(data)]
	copy(slice[l:], data)
	return slice
}

const (
	WeightDel uint8 = 1
)

type item struct {
	value  interface{}
	weight int
	status uint8
}

type Weight struct {
	items         []*item
	rm            bool
	currItemValue interface{}
}

func (w *Weight) Add(k string, v interface{}, weight int) *Weight {
	if w.items == nil {
		w.items = []*item{}
	}
	w.items = append(w.items, &item{value: v, weight: weight})
	return w
}

func (w *Weight) isItemRM(i *item) bool {
	return w.rm && i.status == WeightDel
}

func (w *Weight) RandomValue() interface{} {
	total := 0
	for _, item := range w.items {
		if w.isItemRM(item) {
			continue
		}
		total += item.weight
	}
	if total == 0 {
		return nil
	}
	rd := rand.Intn(total)
	currsum := 0
	for _, item := range w.items {
		if w.isItemRM(item) {
			continue
		}
		currsum += item.weight
		if rd <= currsum {
			if w.rm {
				item.status = WeightDel
			}
			return item.value
		}
	}
	return nil
}

func (w *Weight) NextRandom() bool {
	w.rm = true
	w.currItemValue = w.RandomValue()
	return w.currItemValue != nil
}

func (w *Weight) Value() interface{} {
	return w.currItemValue
}

func CheckIP(addr string) error {
	conn, err := net.DialTimeout("tcp", addr, time.Second*5)
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	return err
}

const (
	formatStart = iota
	formatAllAnd
	formatAllOr
	formatGroupStart
	formatGroupEnd
	formatGroupEndAnd
)

var loc, _ = time.LoadLocation("Asia/Shanghai")

func TimeParseInShanghai(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, loc)
}

func ShanghaiNowTime() time.Time { return time.Now().In(loc) }
