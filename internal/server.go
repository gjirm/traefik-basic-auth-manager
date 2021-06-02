package tbam

import (
	"io/ioutil"
	"net/http"
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

// List valid
func LoginManager(c *gin.Context) {

	cLog := log.WithFields(logrus.Fields{
		"path":   "/",
		"remote": c.ClientIP(),
	})

	cLog.Info("Login manager")

	c.Header("Cache-Control", "no-cache")

	// Get auth cookie
	cookie, err := c.Cookie(config.Cookie.Name)
	if err != nil {
		msg := "Auth cookie " + config.Cookie.Name + " not found"
		cLog.Error(msg)
		c.String(400, msg)
		return
	}

	// Validate cookie
	valid, msg := ValidateCookie(cookie)
	if !valid {
		// Cookie is not valid
		cLog.Error(msg)
		c.String(400, msg)
		return

	}

	cLog.WithFields(logrus.Fields{
		"user": msg,
	}).Info("Valid request")

	user := strings.Split(msg, "@")

	c.HTML(http.StatusOK, "login.html", gin.H{
		"user": user[0],
	})

}

// Add user test
func GetUser(c *gin.Context) {

	cLog := log.WithFields(logrus.Fields{
		"path":   "/getuser",
		"remote": c.ClientIP(),
	})

	c.Header("Cache-Control", "no-cache")

	// Get auth cookie
	cookie, err := c.Cookie(config.Cookie.Name)
	if err != nil {
		msg := "Auth cookie " + config.Cookie.Name + " not found"
		cLog.Error(msg)
		c.String(400, msg)
		return
	}

	// Validate cookie
	valid, msg := ValidateCookie(cookie)
	if !valid {
		// Cookie is not valid
		cLog.Error(msg)
		c.String(400, msg)
		return

	}

	cLog.WithFields(logrus.Fields{
		"user": msg,
	}).Info("Valid request")

	user := strings.Split(msg, "@")

	expire, err := GetValue("uservalid", user[0])
	if err != nil {
		msg := "Error"
		cLog.Error(msg)
		c.String(400, msg)
		return
	}

	if expire == "Key Not Found" {
		msg := "User not exists"
		c.String(200, msg)
		return
	}

	i, _ := strconv.ParseInt(string(expire), 10, 64)
	expireTime := time.Unix(i, 0)

	c.JSON(200, gin.H{
		"username": user[0],
		"validity": expireTime.UTC().Format("2006-01-02T15:04:05-0700"),
	})

}

// Add user test
func GenerateCredentials(c *gin.Context) {

	cLog := log.WithFields(logrus.Fields{
		"path":   "/generate",
		"remote": c.ClientIP(),
	})

	c.Header("Cache-Control", "no-cache")

	// Get auth cookie
	cookie, err := c.Cookie(config.Cookie.Name)
	if err != nil {
		msg := "Auth cookie " + config.Cookie.Name + " not found"
		cLog.Error(msg)
		c.String(400, msg)
		return
	}

	// Validate cookie
	valid, msg := ValidateCookie(cookie)
	if !valid {
		// Cookie is not valid
		cLog.Error(msg)
		c.String(400, msg)
		return

	}

	cLog.WithFields(logrus.Fields{
		"user": msg,
	}).Info("Valid request")

	user := strings.Split(msg, "@")

	pwd, hash, err := GetRandomBcryptHash()
	if err != nil {
		msg := "Error"
		cLog.Error(msg)
		c.String(400, msg)
		return
	}

	err = PutValue("users", user[0], hash)
	//test, err := GetValue("test", "kolo")
	if err != nil {
		msg := "Error"
		cLog.Error(msg)
		c.String(400, msg)
		return
	}

	sessionExpire := time.Now().Add(time.Duration(config.Validity.Session) * time.Second)
	err = PutValue("sessions", user[0], strconv.FormatInt(sessionExpire.Unix(), 10))
	//test, err := GetValue("test", "kolo")
	if err != nil {
		msg := "Error"
		cLog.Error(msg)
		c.String(400, msg)
		return
	}

	credentialExpire := time.Now().Add(time.Duration(config.Validity.Credential) * time.Second)
	err = PutValue("uservalid", user[0], strconv.FormatInt(credentialExpire.Unix(), 10))
	//test, err := GetValue("test", "kolo")
	if err != nil {
		msg := "Error"
		cLog.Error(msg)
		c.String(400, msg)
		return
	}

	err = UpdateCredentials(user[0])
	if err != nil {
		msg := "Error"
		cLog.Error(msg)
		c.String(400, msg)
		return
	}

	cLog.Info("ok")
	c.JSON(200, gin.H{
		"username": user[0],
		"password": pwd,
	})

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
		"path":   "/list/" + bucket,
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

	expireTime := time.Now().Add(30 * time.Minute)
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

var cron = func() {
	bucket := "sessions"
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

	bucket = "uservalid"
	users, _ = ListBucket(bucket)

	for user, expire := range users {
		i, _ := strconv.ParseInt(string(expire), 10, 64)
		expireTime := time.Unix(i, 0).Unix()

		if time.Now().Unix() > expireTime {
			DeleteKey(bucket, user)
			DeleteKey("users", user)
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

	// Load html templates
	r.LoadHTMLGlob("templates/*")

	s := gocron.NewScheduler(time.UTC)

	s.Every(5).Seconds().Do(cron)

	// you can start running the scheduler in two different ways:
	// starts the scheduler asynchronously
	s.StartAsync()

	// Upload file to uploader
	r.GET("/", LoginManager)

	r.GET("/generate", GenerateCredentials)

	r.GET("/getuser", GetUser)

	r.GET("/list/:bucket", ListBucketItems)

	r.GET("/del/:user", DelTest)

	r.GET("/add/:user", AddUser)

	r.GET("/up/:user", ActivateUser)

	r.GET("/expire", Expire)

	log.Infof("Listening on :%d", config.Webserver.Port)
	portNumber := ":" + strconv.Itoa(config.Webserver.Port)
	r.Run(portNumber)
}
