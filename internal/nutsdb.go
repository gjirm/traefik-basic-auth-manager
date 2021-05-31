package tbam

import (
	"io/ioutil"

	"github.com/xujiajun/nutsdb"
)

var db *nutsdb.DB

// InitDB exported
func InitDB() {
	opt := nutsdb.DefaultOptions
	fileDir := "./nutsdb"
	var noDb bool = true
	files, _ := ioutil.ReadDir(fileDir)
	for _, f := range files {
		name := f.Name()
		if name != "" {
			log.Info("Using existing DB: " + fileDir + "/" + name)
			noDb = false
			//err := os.RemoveAll(fileDir + "/" + name)
			// if err != nil {
			// 	panic(err)
			// }
		} else {
			log.Info("Creating new DB: " + fileDir + "/" + name)
		}
	}
	opt.Dir = fileDir
	opt.SegmentSize = 1024 * 1024 // 1MB
	db, _ = nutsdb.Open(opt)
	if noDb {
		PutValue("test", "ipcounter", "jedna")
		PutValue("test", "ipspace", "dva")
		PutValue("test", "server", "tri")
	}
	//bucket = "bucketForString"
}

// PutValue
func PutValue(bucket, keyname, value string) error {
	if err := db.Update(
		func(tx *nutsdb.Tx) error {
			key := []byte(keyname)
			val := []byte(value)
			return tx.Put(bucket, key, val, 0)
		}); err != nil {
		log.Debugf("Error putting value: %w", err)
		return err
	}
	return nil
}

// GetValue
func GetValue(bucket, keyname string) (string, error) {
	var getvalue string
	if err := db.View(
		func(tx *nutsdb.Tx) error {
			key := []byte(keyname)
			e, err := tx.Get(bucket, key)
			if err != nil {
				if err == nutsdb.ErrKeyNotFound {
					getvalue = "Key Not Found"
					return nil
				}
				return err
			}
			getvalue = string(e.Value)
			return nil

		}); err != nil {
		log.Debugf("Error getting value: %w", err)
		return "", err
	}
	return getvalue, nil
}

// Delete key
func DeleteKey(bucket, keyname string) error {
	if err := db.Update(
		func(tx *nutsdb.Tx) error {
			key := []byte(keyname)
			return tx.Delete(bucket, key)
		}); err != nil {
		log.Debugf("Error deleting value: %w", err)
		return err
	}
	return nil
}

// List bucket
func ListBucket(bucket string) (map[string]string, error) {
	var m = map[string]string{}
	if err := db.View(
		func(tx *nutsdb.Tx) error {
			entries, err := tx.GetAll(bucket)
			if err != nil {
				if err == nutsdb.ErrBucketEmpty {
					m["status"] = "Bucket Is Empty"
					return nil
				}
				return err
			}

			for _, entry := range entries {

				m[string(entry.Key)] = string(entry.Value)
				//fmt.Println(string(entry.Key), string(entry.Value))
			}
			return nil
		}); err != nil {
		log.Debugf("Error listing bucket: "+bucket+", %w", err)
		return map[string]string{}, err
	}
	return m, nil
}

// func delete() {
// 	if err := db.Update(
// 		func(tx *nutsdb.Tx) error {
// 			key := []byte("name1")
// 			return tx.Delete(bucket, key)
// 		}); err != nil {
// 		log.Fatal(err)
// 	}
// }

// func put() {
// 	if err := db.Update(
// 		func(tx *nutsdb.Tx) error {
// 			key := []byte("name1")
// 			val := []byte("val1")
// 			return tx.Put(bucket, key, val, 0)
// 		}); err != nil {
// 		log.Fatal(err)
// 	}
// }

// func read() {
// 	if err := db.View(
// 		func(tx *nutsdb.Tx) error {
// 			key := []byte("name1")
// 			e, err := tx.Get(bucket, key)
// 			if err != nil {
// 				return err
// 			}
// 			fmt.Println("val:", string(e.Value))

// 			return nil
// 		}); err != nil {
// 		log.Println(err)
// 	}
// }
