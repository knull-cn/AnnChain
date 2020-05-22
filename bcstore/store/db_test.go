package store

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path"
	"testing"
	"time"
)

const (
	strChars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz" // 62 characters
)

var (
	dirs []string
)

func init() {
	rand.Seed(time.Now().UnixNano())
}
func RandStr(length int) string {
	return fmt.Sprintf("%x", RandBytes(length))
}

// Constructs an alphanumeric string of given length.
func RandBytes(length int) []byte {
	chars := []byte{}
MAIN_LOOP:
	for {
		val := rand.Int63()
		for i := 0; i < 10; i++ {
			v := int(val & 0x3f) // rightmost 6 bits
			if v >= 62 {         // only 62 characters in strChars
				val >>= 6
				continue
			} else {
				chars = append(chars, strChars[v])
				if len(chars) == length {
					break MAIN_LOOP
				}
				val >>= 6
			}
		}
	}

	return chars
}

func RandInt64(n int64) int64 {
	return rand.Int63n(n)
}

func dbname() string {
	dirname := os.TempDir()
	dirname = path.Join(os.TempDir(), "test_leveldb_"+RandStr(12))
	dirs = append(dirs, dirname)
	return dirname
}

func TestMain(m *testing.M) {
	r := m.Run()
	finish()
	os.Exit(r)
}

func finish() {
	for _, dir := range dirs {
		err := os.RemoveAll(dir)
		if err != nil {
			fmt.Println("remove(", dir, ")failed :", err, ".")
		}
	}
}

func Ctx(d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.TODO(), d)
}
