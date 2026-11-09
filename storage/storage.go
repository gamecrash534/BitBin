package storage

import (
	"log"
	"math/rand/v2"
	"time"
)

type Config struct {
	CompressFiles     bool
	DbCompressOnStart bool
	DbCacheMiB        int
	DataTtlMin        int64
}

const keyCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var conf Config

var rnd *rand.Rand

func Init(p string, c Config) {
	conf = c
	src := rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixMicro()))
	rnd = rand.New(src)

	InitDB(Paths(p, "database"), c.DbCompressOnStart, c.DbCacheMiB)
	InitFS(Paths(p, "files"), c.CompressFiles)
}

func Create(d []byte, mime string, alreadyCompressed bool) (*string, *Data, error) {
	key, err := nonexistentStr()
	if err != nil {
		return nil, nil, err
	}

	data := Data{
		Deletion: 0,
		MimeType: mime,
	}

	err = StoreFile(*key, d, alreadyCompressed)
	if err != nil {
		return nil, nil, err
	}

	err = PutDB(*key, data)
	if err != nil {
		return nil, nil, err
	}

	log.Println("Created, key: '", key, "'")
	return key, &data, nil
}

func Get(key string) ([]byte, *Data, error) {
	if e, err := Exists(key); err != nil {
		return nil, nil, err
	} else if !e {
		return nil, nil, nil
	}

	data, err := GetDB(key)
	if err != nil {
		return nil, nil, err
	}

	contents, err := GetContents(key)
	if err != nil {
		return nil, nil, err
	}

	log.Println("Retrieved, key: '", key, "'")
	return contents, data, nil
}

func Delete(key string) error {
	err := RemoveDB(key)
	if err != nil {
		return err
	}

	err = DeleteFile(key)

	if err != nil {
		return err
	}

	log.Println("Deleted, key: '", key, "'")
	return nil
}

func randomStr() string {
	res := make([]byte, 8)
	for i := range res {
		res[i] = keyCharset[rnd.IntN(len(keyCharset))]
	}

	return string(res)
}

func nonexistentStr() (*string, error) {
	var k string
	for {
		k = randomStr()
		exists, err := Exists(k)

		if err != nil {
			return nil, err
		}

		if !exists {
			break
		}
	}

	return &k, nil
}
