package main

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/sirupsen/logrus"
	"github.com/y-miyazaki/go-common/example/initial/entity"
	"github.com/y-miyazaki/go-common/pkg/infrastructure"
	"github.com/y-miyazaki/go-common/pkg/repository"
	"github.com/y-miyazaki/go-common/pkg/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// --------------------------------------------------------------
	// logrus
	// --------------------------------------------------------------
	logrusLogger := &logrus.Logger{}
	logrusLogger.Formatter = &logrus.JSONFormatter{}
	logrusLogger.Out = os.Stdout
	logrusLogger.Level, _ = logrus.ParseLevel("Info")
	logger := infrastructure.NewLogger(logrusLogger)

	// --------------------------------------------------------------
	// logger for gorm
	// --------------------------------------------------------------
	loggerGorm := infrastructure.NewLoggerGorm(&infrastructure.LoggerGormConfig{
		Logger: logger.Entry.Logger,
		GormConfig: &infrastructure.GormConfig{
			// slow query time: 3 sec
			SlowThreshold:             time.Second * 3,
			IgnoreRecordNotFoundError: false,
			LogLevel:                  infrastructure.Info,
		},
	})
	gc := &gorm.Config{
		Logger: loggerGorm,
	}

	// --------------------------------------------------------------
	// MySQL
	// --------------------------------------------------------------
	mysqlDBname := os.Getenv("MYSQL_DBNAME")
	mysqlUsername := os.Getenv("MYSQL_USERNAME")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlServer := os.Getenv("MYSQL_SERVER")
	mysqlPort := os.Getenv("MYSQL_PORT")

	mysqlConfig := &infrastructure.MySQLConfig{
		Config: &mysql.Config{
			DSN:                       utils.GetMySQLDsn(mysqlUsername, mysqlPassword, mysqlServer, mysqlPort, mysqlDBname, "charset=utf8mb4&parseTime=True&loc=Local"),
			DefaultStringSize:         256,   // default size for string fields
			DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
			DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
			DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
			SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
		},
		DBConfig: infrastructure.DBConfig{
			// ConnMaxLifetime sets max life time(sec)
			ConnMaxLifetime: time.Minute * 5,
			// ConnMaxIdletime sets max idle time(sec)
			ConnMaxIdletime: time.Minute * 5,
			// MaxIdleConns sets idle connection
			MaxIdleConns: 20,
			// MaxOpenConns sets max connection
			MaxOpenConns: 100,
		},
	}
	db := infrastructure.NewMySQL(mysqlConfig, gc)

	// --------------------------------------------------------------
	// S3(minio)
	// --------------------------------------------------------------
	s3Region := os.Getenv("S3_REGION")
	s3Endpoint := os.Getenv("S3_ENDPOINT")
	s3ID := os.Getenv("S3_ID")
	s3Secret := os.Getenv("S3_SECRET")
	s3Token := os.Getenv("S3_TOKEN")

	s3SessionOptions := infrastructure.GetDefaultOptions()
	s3Config := infrastructure.GetS3Config(logger.Entry, s3ID, s3Secret, s3Token, s3Region, s3Endpoint, true)
	s3 := infrastructure.NewS3(s3SessionOptions, s3Config)

	// --------------------------------------------------------------
	// example: MySQL
	// --------------------------------------------------------------
	user := entity.User{}
	db.Take(&user)
	logrusLogger.Infof("name = %s, email = %s", user.Name, user.Email)

	// --------------------------------------------------------------
	// example: S3
	// --------------------------------------------------------------
	awsS3Repository := repository.NewAWSS3Repository(logger.Entry, s3)
	text := "aaaaaaaab"
	bucket := "test"

	// Create Bucket
	_, err := awsS3Repository.CreateBucket(bucket)
	if err != nil {
		logger.WithError(err).Errorf("can't create s3 bucket")
	}

	// ListBuckets
	listBuckets, err := awsS3Repository.ListBuckets()
	if err == nil {
		for _, b := range listBuckets.Buckets {
			logger.Infof("bucket = %s(%s)", aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
		}
	} else {
		logger.WithError(err).Errorf("can't list of s3 bucket")
	}

	// Put Object
	_, err = awsS3Repository.PutObjectText(bucket, "test.txt", &text)
	if err != nil {
		logger.WithError(err).Errorf("can't put s3 object")
	}

	// Get Object
	object, err := awsS3Repository.GetObject(bucket, "test.txt")
	if err != nil {
		logger.WithError(err).Errorf("can't get s3 object")
	}
	rc := object.Body
	defer rc.Close()

	text, err = utils.GetStringFromReadCloser(rc)
	if err != nil {
		logger.WithError(err).Errorf("can't get text")
	}
	logger.Infof("text.txt = %s", text)

	// ListObjectV2
	listObjects, err := awsS3Repository.ListObjectsV2(bucket, "")
	if err == nil {
		for _, o := range listObjects.Contents {
			logger.Infof("Object key = %s", aws.StringValue(o.Key))
		}
	} else {
		logger.WithError(err).Errorf("can't list of s3 object")
	}

	// Delete Object
	_, err = awsS3Repository.DeleteObject(bucket, "test.txt")
	if err != nil {
		logger.WithError(err).Errorf("can't delete s3 object")
	}

	// Delete Bucket
	_, err = awsS3Repository.DeleteBucket(bucket)
	if err != nil {
		logger.WithError(err).Errorf("can't delete s3 bucket")
	}
}
