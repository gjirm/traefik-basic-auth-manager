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

// Validate user cookie
func validateUser(log *logrus.Entry, c *gin.Context) (string, error) {
	// Get auth cookie
	cookie, err := c.Cookie(config.Cookie.Name)
	if err != nil {
		msg := "Auth cookie " + config.Cookie.Name + " not found"
		log.Error(msg)
		return msg, err
	}

	// Validate cookie
	valid, msg := ValidateCookie(cookie)
	if !valid {
		// Cookie is not valid
		log.Error(msg)
		return msg, err
	}

	log.WithFields(logrus.Fields{
		"user": msg,
	}).Info("Valid request")

	return strings.Split(msg, "@")[0], nil
}

// Validate if user is admin
func validateAdmin(user string) bool {

	for _, item := range config.Admin {
		if item == user {
			return true
		}
	}
	return false
}

// Main index
func MainIndex(c *gin.Context) {

	cLog := log.WithFields(logrus.Fields{
		"path":   "/",
		"remote": c.ClientIP(),
	})

	c.Header("Cache-Control", "no-cache")

	user, err := validateUser(cLog, c)
	if err != nil {
		c.JSON(200, gin.H{
			"status": user,
		})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"user":  user,
		"title": config.Titleurl,
	})

}

// Check status of users's credentials
//   /status
func CheckUserStatus(c *gin.Context) {

	cLog := log.WithFields(logrus.Fields{
		"path":   "/status",
		"remote": c.ClientIP(),
	})

	c.Header("Cache-Control", "no-cache")

	user, err := validateUser(cLog, c)
	if err != nil {
		c.JSON(400, gin.H{
			"status": user,
		})
		return
	}

	expire, err := GetValue("credentials", user)
	if err != nil {
		msg := "Error"
		cLog.Error(msg)
		c.JSON(400, gin.H{
			"status": msg,
		})
		return
	}

	expireCredentialStr := ""
	if expire == "Key Not Found" {
		expireCredentialStr = "User not exists"
		c.String(200, expireCredentialStr)
		return
	} else {
		i, _ := strconv.ParseInt(string(expire), 10, 64)
		expireCredential := time.Unix(i, 0)
		expireCredentialStr = expireCredential.UTC().Format("2006-01-02T15:04:05-0700")
	}

	expire, err = GetValue("sessions", user)
	if err != nil {
		msg := "Error"
		cLog.Error(msg)
		c.JSON(400, gin.H{
			"status": msg,
		})
		return
	}

	expireSessionStr := ""
	if expire == "Key Not Found" {
		expireSessionStr = "Session not active"
	} else {
		i, _ := strconv.ParseInt(string(expire), 10, 64)
		expireSession := time.Unix(i, 0)
		expireSessionStr = expireSession.UTC().Format("2006-01-02T15:04:05-0700")
	}

	c.JSON(200, gin.H{
		"session":    expireSessionStr,
		"credential": expireCredentialStr,
	})

}

// Activate session for user
//   /activate
func ActivateSession(c *gin.Context) {

	cLog := log.WithFields(logrus.Fields{
		"path":   "/activate",
		"remote": c.ClientIP(),
	})

	c.Header("Cache-Control", "no-cache")

	user, err := validateUser(cLog, c)
	if err != nil {
		c.JSON(400, gin.H{
			"status": user,
		})
		return
	}

	test, err := GetValue("users", user)
	if err != nil {
		msg := "Error"
		cLog.Error(msg)
		c.JSON(400, gin.H{
			"status": msg,
		})
		return
	}

	if test == "Key Not Found" {
		msg := "User not exists"
		c.String(200, msg)
		return
	}

	expireTime := time.Now().Add(time.Duration(config.Validity.Session) * time.Second)
	err = PutValue("sessions", user, strconv.FormatInt(expireTime.Unix(), 10))
	if err != nil {
		msg := "Error"
		cLog.Error(msg)
		c.JSON(400, gin.H{
			"status": msg,
		})
		return
	}

	err = UpdateCredentials(user)
	if err != nil {
		msg := "Error"
		cLog.Error(msg)
		c.JSON(400, gin.H{
			"status": msg,
		})
		return
	}

	c.JSON(200, gin.H{
		"validity": expireTime.UTC().Format("2006-01-02T15:04:05-0700"),
	})

}

// Generate new basic authentication credentials for user
//   /generate
func GenerateCredentials(c *gin.Context) {

	cLog := log.WithFields(logrus.Fields{
		"path":   "/generate",
		"remote": c.ClientIP(),
	})

	c.Header("Cache-Control", "no-cache")

	user, err := validateUser(cLog, c)
	if err != nil {
		c.JSON(400, gin.H{
			"status": user,
		})
		return
	}

	pwd, hash, err := GetRandomBcryptHash()
	if err != nil {
		msg := "Error getting random hash"
		cLog.Error(msg)
		c.JSON(200, gin.H{
			"status": msg,
		})
		return
	}

	err = PutValue("users", user, hash)
	if err != nil {
		msg := "Error put users: " + user
		cLog.Error(msg)
		c.JSON(400, gin.H{
			"status": msg,
		})
		return
	}

	sessionExpire := time.Now().Add(time.Duration(config.Validity.Session) * time.Second)
	err = PutValue("sessions", user, strconv.FormatInt(sessionExpire.Unix(), 10))
	if err != nil {
		msg := "Error put sessions: " + user
		cLog.Error(msg)
		c.JSON(400, gin.H{
			"status": msg,
		})
		return
	}

	credentialExpire := time.Now().Add(time.Duration(config.Validity.Credential) * time.Second)
	err = PutValue("credentials", user, strconv.FormatInt(credentialExpire.Unix(), 10))
	if err != nil {
		msg := "Error put credentials: " + user
		cLog.Error(msg)
		c.JSON(400, gin.H{
			"status": msg,
		})
		return
	}

	err = UpdateCredentials(user)
	if err != nil {
		msg := "Error update credentials: " + user
		cLog.Error(msg)
		c.JSON(400, gin.H{
			"status": msg,
		})
		return
	}

	c.JSON(200, gin.H{
		"username": user,
		"password": pwd,
	})
}

// Delete user
//   /del/:username
func DelUser(c *gin.Context) {

	userName := c.Param("user")

	cLog := log.WithFields(logrus.Fields{
		"path":   "/del/" + userName,
		"remote": c.ClientIP(),
	})

	c.Header("Cache-Control", "no-cache")

	// Validate user authentication
	user, err := validateUser(cLog, c)
	if err != nil {
		c.JSON(400, gin.H{
			"status": user,
		})
		return
	}

	// Validate if user is admin
	if !(validateAdmin(user)) {
		msg := "User is not admin: " + user
		cLog.Error(msg)
		c.JSON(400, gin.H{
			"status": msg,
		})
		return
	}

	err = DeleteKey("users", userName)
	if err != nil {
		msg := "Error deleting: " + userName
		cLog.Error(msg)
		c.JSON(400, gin.H{
			"status": msg,
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
	})

}

// List DB bucket
//  /list/:bucketname
// Buckets:
//   users - list of usernames and their basic authentication bcrypt hashes
//   sessions - list of usernames and their login validity - until when users are allowed to login
//   credentials - list of usernames and their crednetial validity - until when the user's credentials are valid
func ListBucketObjects(c *gin.Context) {

	bucket := c.Param("bucket")
	bucket = strings.TrimPrefix(bucket, "/")

	cLog := log.WithFields(logrus.Fields{
		"path":   "/list/" + bucket,
		"remote": c.ClientIP(),
	})

	c.Header("Cache-Control", "no-cache")

	// Validate user authentication
	user, err := validateUser(cLog, c)
	if err != nil {
		c.JSON(400, gin.H{
			"status": user,
		})
		return
	}

	// Validate if user is admin
	if !(validateAdmin(user)) {
		msg := "User is not admin: " + user
		cLog.Error(msg)
		c.JSON(400, gin.H{
			"status": msg,
		})
		return
	}

	bucketList, err := ListBucket(bucket)
	if err != nil {
		msg := "Error listing: " + bucket
		cLog.Error(msg)
		c.JSON(400, gin.H{
			"status": msg,
		})
		return
	}

	c.JSON(200, bucketList)

}

// Add new user
//   /add/:username
func AddUser(c *gin.Context) {

	userName := c.Param("user")

	cLog := log.WithFields(logrus.Fields{
		"path":   "/add/" + userName,
		"remote": c.ClientIP(),
	})

	c.Header("Cache-Control", "no-cache")

	user, err := validateUser(cLog, c)
	if err != nil {
		c.JSON(400, gin.H{
			"status": user,
		})
		return
	}

	pwd, hash, err := GetRandomBcryptHash()
	if err != nil {
		msg := "Error getting random hash"
		cLog.Error(msg)
		c.JSON(400, gin.H{
			"status": msg,
		})
		return
	}

	err = PutValue("users", userName, hash)
	if err != nil {
		msg := "Error put users: " + userName
		cLog.Error(msg)
		c.JSON(400, gin.H{
			"status": msg,
		})
		return
	}

	sessionExpire := time.Now().Add(time.Duration(config.Validity.Session) * time.Second)
	err = PutValue("sessions", userName, strconv.FormatInt(sessionExpire.Unix(), 10))
	if err != nil {
		msg := "Error put sessions: " + userName
		cLog.Error(msg)
		c.JSON(400, gin.H{
			"status": msg,
		})
		return
	}

	credentialExpire := time.Now().Add(time.Duration(config.Validity.Credential) * time.Second)
	err = PutValue("credentials", userName, strconv.FormatInt(credentialExpire.Unix(), 10))
	if err != nil {
		msg := "Error put credentials: " + userName
		cLog.Error(msg)
		c.JSON(400, gin.H{
			"status": msg,
		})
		return
	}

	err = UpdateCredentials(userName)
	if err != nil {
		msg := "Error update credentials: " + userName

		cLog.Error(msg)
		c.JSON(400, gin.H{
			"status": msg,
		})
		return
	}

	c.JSON(200, gin.H{
		"username": userName,
		"password": pwd,
	})
}

// Cron task
var cron = func() {
	bucket := "sessions"
	users, _ := ListBucket(bucket)

	log.Debug("Running cron task")
	if !(users["status"] == "Bucket Is Empty") {
		for user, expire := range users {
			i, _ := strconv.ParseInt(string(expire), 10, 64)
			expireTime := time.Unix(i, 0).Unix()

			if time.Now().Unix() > expireTime {
				log.Debug("Cron - " + bucket + " - deleting: " + user)
				DeleteKey(bucket, user)
				RemoveCredentials(user)
			}

		}
	}

	bucket = "credentials"
	users, _ = ListBucket(bucket)
	if !(users["status"] == "Bucket Is Empty") {
		for user, expire := range users {
			i, _ := strconv.ParseInt(string(expire), 10, 64)
			expireTime := time.Unix(i, 0).Unix()

			if time.Now().Unix() > expireTime {
				log.Debug("Cron - " + bucket + " - deleting: " + user)
				DeleteKey(bucket, user)
				DeleteKey("users", user)
			}
		}
	}
}

// Webserver instance
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

	// Cron job for checking expiration of sessions and credentials
	s := gocron.NewScheduler(time.UTC)
	s.Every(5).Seconds().Do(cron)
	s.StartAsync()

	// Main index
	r.GET("/", MainIndex)

	// Public api for authenticated users

	r.GET("/status", CheckUserStatus)

	r.GET("/activate", ActivateSession)

	r.GET("/generate", GenerateCredentials)

	// Admin api for authenticated users and users listed as admin in the config

	// Buckets:
	//   users - list of usernames and their basic authentication bcrypt hashes
	//   sessions - list of usernames and their login validity - until when users are allowed to login
	//   credentials - list of usernames and their crednetial validity - until when the user's credentials are valid
	r.GET("/list/:bucket", ListBucketObjects)

	r.GET("/del/:user", DelUser)

	r.GET("/add/:user", AddUser)

	log.Infof("Listening on :%d", config.Webserver.Port)
	portNumber := ":" + strconv.Itoa(config.Webserver.Port)
	r.Run(portNumber)
}
