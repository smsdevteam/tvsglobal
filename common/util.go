package common

import (
	"fmt"
	"strconv"
	"time"
)

// StrToInt : Convert string to int (base 10)
// error : return 0
func StrToInt(iStr string) int {
	i, err := strconv.Atoi(iStr)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return i
}

// StrToInt64 : Convert string to int (base 64)
// error : return 0
func StrToInt64(iStr string) int64 {
	i, err := strconv.ParseInt(iStr, 10, 64)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return i
}

// StrTofloat64 : Convert string to float   (base 64)
// error : return 0
func StrTofloat64(iStr string) float64 {
	i, err := strconv.ParseFloat(iStr, 64)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return i
}

// IntToStr : Convert int (base 10) to string
func IntToStr(iNum int) string {
	s := strconv.Itoa(iNum)
	return s
}

// Int64ToStr : Convert int (base 64) to string
func Int64ToStr(iNum int64) string {
	s := strconv.FormatInt(iNum, 10)
	return s
}

// Float64ToStr : Convert int (base 64) to string
func Float64ToStr(iNum float64) string {
	s := strconv.FormatFloat(iNum, 'E', -1, 64)
	return s
}

// TimeToStr : Convert time to format
func TimeToStr(itimein time.Time) string {
	//s := itimein.String()
	s := itimein.Format("2006-01-02T15:04:05.99")
	return s
}

// StrToBool : Convert string (1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False) to boolean (true or false)
func StrToBool(iStr string) bool {
	b, err := strconv.ParseBool(iStr)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return b
}

/* func main() {
	var str string
	var num int

	fmt.Print("Input : ")
	fmt.Scanln(&str)

	var fsi = StrToInt
	var fis = IntToStr

	num = fsi(str)
	fmt.Println("Number : ", num)

	word := fis(num)
	fmt.Println("String : ", word)
} */
