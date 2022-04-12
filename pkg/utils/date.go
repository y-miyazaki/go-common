package utils

import "time"

const (
	DateFormatHyphenYearMonthDay                 = "2006-01-02"
	DateFormatHyphenYearMonthDayHourMinuteSecond = "2006-01-02 15:04:05"
	DateFormatHyphenMonthDayYear                 = "01-02-2006"
	DateFormatHyphenMonthDayYearHourMinuteSecond = "01-02-2006 15:04:05"
	DateFormatHyphenDayMonthYear                 = "02-01-2006"
	DateFormatHyphenDayMonthYearHourMinuteSecond = "02-01-2006 15:04:05"
	DateFormatHyphenMonthDay                     = "01-02"
	DateFormatHyphenDayMonth                     = "02-01"

	DateFormatSlashYearMonthDay                 = "2006/01/02"
	DateFormatSlashYearMonthDayHourMinuteSecond = "2006/01/02 15:04:05"
	DateFormatSlashMonthDayYear                 = "01/02/2006"
	DateFormatSlashMonthDayYearHourMinuteSecond = "01/02/2006 15:04:05"
	DateFormatSlashDayMonthYear                 = "02/01/2006"
	DateFormatSlashDayMonthYearHourMinuteSecond = "02/01/2006 15:04:05"
	DateFormatSlashMonthDay                     = "01/02"
	DateFormatSlashDayMonth                     = "02/01"

	DateFormatYear             = "2006"
	DateFormatHourMinuteSecond = "15:04:05"

	DateConvertUTCToJSTOffset = 9 * 60 * 60
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
	return t.In(time.FixedZone("Asia/Tokyo", DateConvertUTCToJSTOffset))
}

// ConvertJSTToUTC converts JST to UTC
func ConvertJSTToUTC(t time.Time) time.Time {
	return t.In(time.FixedZone("UTC", 0))
}
