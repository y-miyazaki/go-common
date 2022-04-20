package utils

import "time"

const (
	// DateFormatHyphenYearMonthDay : yyyy-mm-dd
	DateFormatHyphenYearMonthDay = "2006-01-02"
	// DateFormatHyphenYearMonthDayHourMinuteSecond : yyyy-mm-dd hh:mm:ss
	DateFormatHyphenYearMonthDayHourMinuteSecond = "2006-01-02 15:04:05"
	// DateFormatHyphenMonthDayYear : mm-dd-yyyy
	DateFormatHyphenMonthDayYear = "01-02-2006"
	// DateFormatHyphenMonthDayYearHourMinuteSecond : mm-dd-yyyy hh:mm:ss
	DateFormatHyphenMonthDayYearHourMinuteSecond = "01-02-2006 15:04:05"
	// DateFormatHyphenDayMonthYear : dd-mm-yyyy
	DateFormatHyphenDayMonthYear = "02-01-2006"
	// DateFormatHyphenDayMonthYearHourMinuteSecond : dd-mm-yyyy hh:mm:ss
	DateFormatHyphenDayMonthYearHourMinuteSecond = "02-01-2006 15:04:05"
	// DateFormatHyphenMonthDay : mm-dd
	DateFormatHyphenMonthDay = "01-02"
	// DateFormatHyphenDayMonth : dd-mm
	DateFormatHyphenDayMonth = "02-01"

	// DateFormatSlashYearMonthDay : yyyy/mm/dd
	DateFormatSlashYearMonthDay = "2006/01/02"
	// DateFormatSlashYearMonthDayHourMinuteSecond : yyyy/mm/dd hh:mm:ss
	DateFormatSlashYearMonthDayHourMinuteSecond = "2006/01/02 15:04:05"
	// DateFormatSlashMonthDayYear : mm/dd/yyyy
	DateFormatSlashMonthDayYear = "01/02/2006"
	// DateFormatSlashMonthDayYearHourMinuteSecond : mm/dd/yyyy hh:mm:ss
	DateFormatSlashMonthDayYearHourMinuteSecond = "01/02/2006 15:04:05"
	// DateFormatSlashDayMonthYear : dd/mm/yyyy
	DateFormatSlashDayMonthYear = "02/01/2006"
	// DateFormatSlashDayMonthYearHourMinuteSecond : dd/mm/yyyy hh:mm:ss
	DateFormatSlashDayMonthYearHourMinuteSecond = "02/01/2006 15:04:05"
	// DateFormatSlashMonthDay : mm/dd
	DateFormatSlashMonthDay = "01/02"
	// DateFormatSlashDayMonth : dd/mm
	DateFormatSlashDayMonth = "02/01"

	// DateFormatYear : yyyy
	DateFormatYear = "2006"
	// DateFormatHourMinuteSecond : hh:mm:ss
	DateFormatHourMinuteSecond = "15:04:05"

	dateConvertUTCToJSTOffset = 9 * 60 * 60
)

// GetDateFormatString gets format date string.
func GetDateFormatString(t time.Time, format string) string {
	return t.Format(format)
}

// GetDateTime gets time.
func GetDateTime(str, format string) (time.Time, error) {
	return time.Parse(format, str)
}

// ConvertUTCToJST converts UTC to JST
func ConvertUTCToJST(t time.Time) time.Time {
	return t.In(time.FixedZone("Asia/Tokyo", dateConvertUTCToJSTOffset))
}

// ConvertJSTToUTC converts JST to UTC
func ConvertJSTToUTC(t time.Time) time.Time {
	return t.In(time.FixedZone("UTC", 0))
}
