package tbam

import (
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/location"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-co-op/gocron"
)

// Upload file to uploader bucket
func UpdateBasicAuthCredentials(c *gin.Context) {

	// cLog := log.WithFields(logrus.Fields{
	// 	"path":   "/",
	// 	"remote": c.ClientIP(),
	// })

	// cookie, err := c.Cookie(config.Cookie.Name)
	// if err != nil {
	// 	msg := "Auth cookie " + config.Cookie.Name + " not found"
	// 	cLog.Error(msg)
	// 	c.String(400, msg)
	// 	return
	// }

	// // Validate cookie
	// valid, msg := ValidateCookie(cookie)
	// if valid {
	// 	// Cookie is valid -> run SSH cmd - avtivate users WireGuard peers
	// 	cLog.WithFields(logrus.Fields{
	// 		"user": msg,
	// 	}).Info("Valid request")

	// 	user := strings.Split(msg, "@")

	// 	username, password, err := UpdateCredentials(user[0])
	// 	if err != nil {
	// 		msg := "Failed to update credentials"
	// 		cLog.Error(msg)
	// 		c.String(400, msg)
	// 		return
	// 	}

	// 	c.JSON(200, gin.H{
	// 		"username": username,
	// 		"password": password,
	// 	})

	// } else {
	// 	cLog.Error(msg)
	// 	c.String(400, msg)
	// 	return
	// }
}

// List test
func DelTest(c *gin.Context) {

	user := c.Param("user")

	cLog := log.WithFields(logrus.Fields{
		"path":   "/del",
		"remote": c.ClientIP(),
	})

	cLog.Info("Del...")

	//test, err := ListBucket("test1")
	err := DeleteKey("user", user)

	if err != nil {
		msg := "Error"
		cLog.Error(msg)
		c.String(400, msg)
		return
	}
	cLog.Info("ok")
	c.JSON(200, "ok")

}

// List valid
func ListBucketItems(c *gin.Context) {

	bucket := c.Param("bucket")
	bucket = strings.TrimPrefix(bucket, "/")

	cLog := log.WithFields(logrus.Fields{
		"path":   "/list" + bucket,
		"remote": c.ClientIP(),
	})

	cLog.Info("Listing bucket")

	bucketList, err := ListBucket(bucket)
	//test, err := GetValue("test", "kolo")
	if err != nil {
		msg := "Error"
		cLog.Error(msg)
		c.String(400, msg)
		return
	}
	//cLog.Info(test)
	c.JSON(200, bucketList)

}

// Add user test
func AddUser(c *gin.Context) {

	user := c.Param("user")

	cLog := log.WithFields(logrus.Fields{
		"path":   "/add/:user",
		"remote": c.ClientIP(),
	})

	cLog.Info("Add...")

	pwd, hash, err := GetRandomBcryptHash()
	if err != nil {
		msg := "Error"
		cLog.Error(msg)
		c.String(400, msg)
		return
	}

	err = PutValue("users", user, hash)
	//test, err := GetValue("test", "kolo")
	if err != nil {
		msg := "Error"
		cLog.Error(msg)
		c.String(400, msg)
		return
	}

	expireTime := time.Now().Add(30 + time.Minute)
	err = PutValue("valid", user, strconv.FormatInt(expireTime.Unix(), 10))
	//test, err := GetValue("test", "kolo")
	if err != nil {
		msg := "Error"
		cLog.Error(msg)
		c.String(400, msg)
		return
	}

	err = UpdateCredentials(user)
	if err != nil {
		msg := "Error"
		cLog.Error(msg)
		c.String(400, msg)
		return
	}

	cLog.Info("ok")
	c.JSON(200, gin.H{
		"username": user,
		"password": pwd,
	})

}

// Activate user
func ActivateUser(c *gin.Context) {

	user := c.Param("user")
	user = strings.TrimPrefix(user, "/")

	cLog := log.WithFields(logrus.Fields{
		"path":   "/up/:user",
		"remote": c.ClientIP(),
	})

	cLog.Info("Activate")

	test, err := GetValue("users", user)
	if err != nil {
		msg := "Error"
		cLog.Error(msg)
		c.String(400, msg)
		return
	}

	if test == "Key Not Found" {
		msg := "User not exists"
		cLog.Error(msg)
		c.String(400, msg)
		return
	}

	expireTime := time.Now().Add(1 * time.Minute)
	err = PutValue("valid", user, strconv.FormatInt(expireTime.Unix(), 10))
	//test, err := GetValue("test", "kolo")
	if err != nil {
		msg := "Error"
		cLog.Error(msg)
		c.String(400, msg)
		return
	}

	err = UpdateCredentials(user)
	if err != nil {
		msg := "Error"
		cLog.Error(msg)
		c.String(400, msg)
		return
	}

	cLog.Info("ok")

}

// List test
func Expire(c *gin.Context) {

	cLog := log.WithFields(logrus.Fields{
		"path":   "/expire",
		"remote": c.ClientIP(),
	})

	cLog.Info("Expiring...")

	err := checkExpire("valid")
	if err != nil {
		msg := "Error"
		cLog.Error(msg)
		c.String(400, msg)
		return
	}

	cLog.Info("ok")
	c.JSON(200, "ok")

}

var task = func() {
	bucket := "valid"
	users, _ := ListBucket(bucket)

	log.Info("Running task")
	for user, expire := range users {
		i, _ := strconv.ParseInt(string(expire), 10, 64)
		expireTime := time.Unix(i, 0).Unix()

		if time.Now().Unix() > expireTime {
			DeleteKey(bucket, user)
			RemoveCredentials(user)
		}

	}
}

func checkExpire(bucket string) error {

	users, err := ListBucket(bucket)
	if err != nil {
		return err
	}

	for user, expire := range users {
		i, _ := strconv.ParseInt(string(expire), 10, 64)
		expireTime := time.Unix(i, 0).Unix()

		if time.Now().Unix() > expireTime {
			err := DeleteKey(bucket, user)
			if err != nil {
				return err
			}

			err = RemoveCredentials(user)
			if err != nil {
				return err
			}
		}
		//fmt.Println(string(entry.Key), string(entry.Value))
	}
	return nil

}

// MyServer server instance
func ApiServer() {

	if !(config.Webserver.Debug) {
		gin.SetMode(gin.ReleaseMode)
	}
	// Disable Console Color
	gin.DisableConsoleColor()

	// Disable gin logging
	gin.DefaultWriter = ioutil.Discard

	// Create router
	r := gin.Default()
	r.Use(location.Default())
	r.Use(cors.Default())

	s := gocron.NewScheduler(time.UTC)

	s.Every(5).Seconds().Do(task)

	// you can start running the scheduler in two different ways:
	// starts the scheduler asynchronously
	s.StartAsync()

	// Upload file to uploader
	r.GET("/", UpdateBasicAuthCredentials)

	r.GET("/list/:bucket", ListBucketItems)

	r.GET("/del/:user", DelTest)

	r.GET("/add/:user", AddUser)

	r.GET("/up/:user", ActivateUser)

	r.GET("/expire", Expire)

	portNumber := ":" + strconv.Itoa(config.Webserver.Port)
	r.Run(portNumber)
}
