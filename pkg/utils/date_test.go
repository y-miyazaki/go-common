package utils

import (
	"reflect"
	"testing"
	"time"
)

func TestGetDateFormatString(t *testing.T) {
	type args struct {
		t      time.Time
		format string
	}
	tt := time.Date(2019, time.February, 9, 8, 7, 4, 0, time.UTC)
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				t:      tt,
				format: DateFormatHyphenYearMonthDay,
			},
			want: "2019-02-09",
		},
		{
			name: "test2",
			args: args{
				t:      tt,
				format: DateFormatHyphenYearMonthDayHourMinuteSecond,
			},
			want: "2019-02-09 08:07:04",
		},
		{
			name: "test3",
			args: args{
				t:      tt,
				format: DateFormatHyphenMonthDayYear,
			},
			want: "02-09-2019",
		},
		{
			name: "test4",
			args: args{
				t:      tt,
				format: DateFormatHyphenMonthDayYearHourMinuteSecond,
			},
			want: "02-09-2019 08:07:04",
		},
		{
			name: "test5",
			args: args{
				t:      tt,
				format: DateFormatHyphenDayMonthYear,
			},
			want: "09-02-2019",
		},
		{
			name: "test6",
			args: args{
				t:      tt,
				format: DateFormatHyphenDayMonthYearHourMinuteSecond,
			},
			want: "09-02-2019 08:07:04",
		},
		{
			name: "test7",
			args: args{
				t:      tt,
				format: DateFormatHyphenMonthDay,
			},
			want: "02-09",
		},
		{
			name: "test8",
			args: args{
				t:      tt,
				format: DateFormatHyphenDayMonth,
			},
			want: "09-02",
		},
		{
			name: "test9",
			args: args{
				t:      tt,
				format: DateFormatYearMonthDay,
			},
			want: "2019/02/09",
		},
		{
			name: "test10",
			args: args{
				t:      tt,
				format: DateFormatYearMonthDayHourMinuteSecond,
			},
			want: "2019/02/09 08:07:04",
		},
		{
			name: "test11",
			args: args{
				t:      tt,
				format: DateFormatMonthDayYear,
			},
			want: "02/09/2019",
		},
		{
			name: "test12",
			args: args{
				t:      tt,
				format: DateFormatMonthDayYearHourMinuteSecond,
			},
			want: "02/09/2019 08:07:04",
		},
		{
			name: "test13",
			args: args{
				t:      tt,
				format: DateFormatDayMonthYear,
			},
			want: "09/02/2019",
		},
		{
			name: "test14",
			args: args{
				t:      tt,
				format: DateFormatDayMonthYearHourMinuteSecond,
			},
			want: "09/02/2019 08:07:04",
		},
		{
			name: "test15",
			args: args{
				t:      tt,
				format: DateFormatMonthDay,
			},
			want: "02/09",
		},
		{
			name: "test16",
			args: args{
				t:      tt,
				format: DateFormatDayMonth,
			},
			want: "09/02",
		},
		{
			name: "test17",
			args: args{
				t:      tt,
				format: DateFormatYear,
			},
			want: "2019",
		},
		{
			name: "test18",
			args: args{
				t:      tt,
				format: DateFormatHourMinuteSecond,
			},
			want: "08:07:04",
		},
	}
	count := 1
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDateFormatString(tt.args.t, tt.args.format); got != tt.want {
				t.Errorf("GetDateFormatString() = %v, want %v", got, tt.want)
			}
		})
		count++
	}
}

func TestGetDateTime(t *testing.T) {
	tt := time.Date(2019, time.February, 9, 8, 7, 4, 0, time.UTC)
	tt2 := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	type args struct {
		str    string
		format string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				str:    "2019-02-09 08:07:04",
				format: DateFormatHyphenYearMonthDayHourMinuteSecond,
			},
			want:    tt,
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				str:    "2019-02-09 08:0:04",
				format: DateFormatHyphenYearMonthDayHourMinuteSecond,
			},
			want:    tt2,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDateTime(tt.args.str, tt.args.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDateTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDateTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertUTCToJST(t *testing.T) {
	tt := time.Date(2019, time.February, 9, 8, 7, 4, 0, time.UTC)
	tt2 := tt.In(time.FixedZone("Asia/Tokyo", DateConvertUTCToJSTOffset))
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "test1",
			args: args{
				t: tt,
			},
			want: tt2,
		},
	}
	tests2 := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				t: tt,
			},
			want: "17:07:04",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertUTCToJST(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertUTCToJST() = %v, want %v", got, tt.want)
			}
		})
	}
	for _, tt := range tests2 {
		t.Run(tt.name, func(t *testing.T) {
			got := GetDateFormatString(ConvertUTCToJST(tt.args.t), DateFormatHourMinuteSecond)
			if got != tt.want {
				t.Errorf("ConvertUTCToJST() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertJSTToUTC(t *testing.T) {
	tt := time.Date(2019, time.February, 9, 8, 7, 4, 0, time.UTC)
	tt2 := tt.In(time.FixedZone("Asia/Tokyo", DateConvertUTCToJSTOffset))
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				t: tt2,
			},
			want: "08:07:04",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetDateFormatString(ConvertJSTToUTC(tt.args.t), DateFormatHourMinuteSecond)
			if got != tt.want {
				t.Errorf("ConvertJSTToUTC() = %v, want %v", got, tt.want)
			}
		})
	}
}
