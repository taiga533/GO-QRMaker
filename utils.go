package QRMaker

import (
    //入力データに対してモードを識別するためにインポート
    "regexp"
    //数値から２進数文字列に変換するためにインポート
    "strconv"
    "unicode/utf8"
)

func getStrLen (src string) int {
    return utf8.RuneCountInString(src)
}

func searchIntArray (array []int, value int) (index int) {
    for i, r := range array {
        if r == value {
            index = i
            return
        }
    }
    index = 0
    return
}
func searchCharArray (array []string, value rune) (index int) {
    for i, r := range array {
        if r == string(value) {
            index = i
            return
        }
    }
    index = 0
    return
}

func checkRegexp (reg, src string) (result bool) {
    result = regexp.MustCompile(reg).MatchString(src)
    return
}
func divideInt (src, interval int) (result []int) {
    str := strconv.FormatUint(uint64(src), 10)
    strLen := getStrLen(str)
    resultLen := strLen / interval
    if strLen % interval != 0 {
        resultLen += 1
    }
    result = make([]int, resultLen)
    for i := resultLen - 1; i >= 0; i-- {
        tmpValue, _ := strconv.Atoi(str[i * interval: strLen])
        result[i] = tmpValue
        strLen -= strLen - i * interval
    }
    return
}
//引数のビット長を可変長として扱うことで汎用性を適度持たせる
func intToBin (src int, bitLen ...int) (binary string) {
    binary = strconv.FormatUint(uint64(src), 2)
    /*
    もしビット長が指定されなかった場合にBitLen[0]を参照するとOut of Indexとなる為
    最初にビット長指定の有無があるか確認する
    */
    if len(bitLen) != 0 {
        if binaryLen := getStrLen(binary); binaryLen < bitLen[0] {
            tempBinary := make([]byte, 0, bitLen[0])
            for i := 0; i < cap(tempBinary) - binaryLen; i++ {
                tempBinary = append(tempBinary, "0"...)
            }
            binary = string(tempBinary) + binary
        }
    }
    return
}
func binToInt (src string) (result int) {
    tmpUint, _ := strconv.ParseUint(src, 2, 64)
    result = int(tmpUint)
    return
}

func divideRoundUp (left, right int) (result int) {
    result = left / right
    if left % right != 0 {
        result++
    }
    return
}
