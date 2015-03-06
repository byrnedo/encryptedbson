package encryptedbson

import (
	"crypto/rand"
	. "github.com/smartystreets/goconvey/convey"
	"labix.org/v2/mgo/bson"
	"testing"
	"time"
)

type EncryptedStruct struct {
	Bool   EncryptedBool
	String EncryptedString
	Float  EncryptedFloat
	Int    EncryptedInt
	Date   EncryptedDate
	Map    EncryptedMap
}

func TestEncryption(t *testing.T) {
	Convey("Encryption", t, func() {

		EnableEncryption = true
		rand.Reader.Read(EncryptionKey[:])

		Convey("Raw encryption/decryption", func() {
			val := "my string"

			encrypted, err := Encrypt(EncryptionKey, []byte(val))

			So(err, ShouldEqual, nil)

			decrypted, err := Decrypt(EncryptionKey, encrypted)

			So(err, ShouldEqual, nil)
			So(string(decrypted), ShouldEqual, "my string")
		})

		Convey("Encrypted types", func() {
			EnableEncryption = true
			date := time.Now().Format(iso8601Format)

			mp := map[string]interface{}{
				"foo":  "bar",
				"baz":  5,
				"boop": 10.5,
				"childMap": map[string]interface{}{
					"bing": "boop",
					"boof": "foop",
				},
			}
			myStruct := &EncryptedStruct{EncryptedBool(true), EncryptedString("foo"), EncryptedFloat(5.555), EncryptedInt(6), EncryptedDate(date), EncryptedMap(mp)}

			// connection.Session.EncryptionKey = key
			marshaled, err := bson.Marshal(myStruct)
			So(err, ShouldEqual, nil)

			newStruct := &EncryptedStruct{}
			err = bson.Unmarshal(marshaled, newStruct)
			So(err, ShouldEqual, nil)
			So(bool(newStruct.Bool), ShouldEqual, true)
			So(string(newStruct.String), ShouldEqual, "foo")
			So(float64(newStruct.Float), ShouldEqual, 5.555)
			So(int(newStruct.Int), ShouldEqual, 6)
			So(newStruct.Date, ShouldEqual, EncryptedDate(date))

			newMp := map[string]interface{}(newStruct.Map)
			So(newMp["foo"], ShouldEqual, "bar")
		})

		Convey("Encrypted types with encryption disabled", func() {
			EnableEncryption = false
			date := time.Now().Format(iso8601Format)

			mp := map[string]interface{}{
				"foo":  "bar",
				"baz":  5,
				"boop": 10.5,
				"childMap": map[string]interface{}{
					"bing": "boop",
					"boof": "foop",
				},
			}
			myStruct := &EncryptedStruct{EncryptedBool(true), EncryptedString("foo"), EncryptedFloat(5.555), EncryptedInt(6), EncryptedDate(date), EncryptedMap(mp)}

			// connection.Session.EncryptionKey = key
			marshaled, err := bson.Marshal(myStruct)
			So(err, ShouldEqual, nil)

			newStruct := &EncryptedStruct{}
			err = bson.Unmarshal(marshaled, newStruct)
			So(err, ShouldEqual, nil)
			So(bool(newStruct.Bool), ShouldEqual, true)
			So(string(newStruct.String), ShouldEqual, "foo")
			So(float64(newStruct.Float), ShouldEqual, 5.555)
			So(int(newStruct.Int), ShouldEqual, 6)
			So(newStruct.Date, ShouldEqual, EncryptedDate(date))

			newMp := map[string]interface{}(newStruct.Map)
			So(newMp["foo"], ShouldEqual, "bar")
		})
	})
}

func BenchmarkEncryptDecrypt(b *testing.B) {

	for i := 0; i < b.N; i++ {
		encrypted, _ := Encrypt(EncryptionKey, []byte("this is a test"))
		decrypted, _ := Decrypt(EncryptionKey, encrypted)

		if string(decrypted) != "this is a test" {
			panic("Failed!")
		}
	}
}
