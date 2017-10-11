package QRMaker

import (
    "log"
    //入力データのutf8エンコードをするためにインポート
    "unicode/utf8"
    //入力データが数値の場合int型にエンコードするためにインポート
    "strconv"
)

var alphaNumLookupTable = []string {
    "0", "1", "2", "3", "4", "5", "6",
    "7", "8", "9", "A", "B", "C","D",
    "E", "F", "G", "H", "I", "J", "K",
    "L", "M", "N", "O", "P", "Q", "R",
    "S", "T", "U", "V", "W", "X", "Y",
    "Z", " ", "$", "%", "*", "+", "-",
    ".", "/", ":",
}

type DataCode struct {
    modeBin string
    numLenBin string
    data []int
}

func newDataCode (src string, dataCodeCnt, modeIndicator int) (dc *DataCode) {
    dc = new(DataCode)
    data := dc.setIndicator(src, modeIndicator)
    log.Printf("モード識別子%v, ビット分割子%v\n", dc.modeBin, dc.numLenBin)
    data = setEndPattrn(data, dataCodeCnt * 8)
    dc.data = convertToCodeLang(data, dataCodeCnt)
    return
}
func (dc *DataCode) setIndicator (src string, modeIndicator int) string {
    var srcBin string
    dc.modeBin = intToBin(modeIndicator, 4)
    switch modeIndicator {
    case 1:
        dc.numLenBin = intToBin(getStrLen(src), versions[qrInfo.versionIdx].dataNumIndicator[0])
        srcBin = dc.modeBin + dc.numLenBin + numToData(src)
    case 2:
        dc.numLenBin = intToBin(getStrLen(src), versions[qrInfo.versionIdx].dataNumIndicator[1])
        srcBin = dc.modeBin + dc.numLenBin + alphaNumToData(src)
    default :
        dc.numLenBin = intToBin(len(src), versions[qrInfo.versionIdx].dataNumIndicator[2])
        srcBin = dc.modeBin + dc.numLenBin + utf8ToData(src)
    }
    return srcBin
}

func utf8ToData (src string) string {
    tmp := make([]byte, 4)
    binSrc := make([]byte, 0, len(src) * 8)
    for _,r := range src {
        idx := utf8.EncodeRune(tmp, r)
        for i := 0; i < idx; i++ {
            binSrc = append(binSrc, intToBin(int(tmp[i]), 8)...)
        }
    }
    return string(binSrc)
}
func numToData (srcStr string) string {
    srcStrBinLen := len(srcStr) / 3 * 10
    srcStrBin := make([]byte, 0, srcStrBinLen)
    for i := 0; i < len(srcStr) - (len(srcStr) % 3); i += 3 {
        srcInt,_ := strconv.Atoi(string(srcStr[i:i + 3]))
        srcStrBin = append(srcStrBin, intToBin(srcInt, 10)...)
    }
    switch len(srcStr) % 3 {
    case 1:
        srcInt,_ := strconv.Atoi(string(srcStr[len(srcStr) - 1]))
        srcStrBin = append(srcStrBin, intToBin(srcInt, 4)...)
    case 2:
        srcInt,_ := strconv.Atoi(string(srcStr[len(srcStr) - 2:]))
        srcStrBin = append(srcStrBin, intToBin(srcInt, 7)...)
    }
    return string(srcStrBin)
}
func alphaNumToData (src string) string {
    anltable := alphaNumLookupTable
    div := make([]int, 0, divideRoundUp(getStrLen(src), 2))
    divIdx := 0
    /*
    引数の2文字ごとに
    一文字目に対応する数値*45の値と2文字目に対応する数値を足した数を配列に格納していく
    引数の文字数が奇数だった場合、配列の最後尾の要素 = 引数の最後の英数字に対応する値となる
    となる
    */
    for i, r := range src {
        if i % 2 == 0 {
            div = append(div, searchCharArray(anltable, r))
        } else {
            div[divIdx] = div[divIdx] * 45 + searchCharArray(anltable, r)
            divIdx++
        }
    }
    //配列の各要素を11bitで表現し、配列の最後尾の要素に含まれるのが1文字の場合6bitで表現する
    binLen := len(div) * 11
    if len(div) % 2 == 1 {
        binLen -= 5
    }
    bin := make([]byte, 0, binLen)
    for _, v := range div {
        if v < 45 {
            bin = append(bin, intToBin(v, 6)...)
        } else {
            bin = append(bin, intToBin(v, 11)...)
        }
    }
    return string(bin)
}
func setEndPattrn (src string, dataLen int) string {
    /*
    モード指定子 + 文字数指定子 + データbit列に
    最大4bitかつシンボル容量を超えない形で終端パターンを付与する
    */
    maxLen := dataLen - getStrLen(src)
    tmpBin := make([]byte, 0, 4)
    for len(tmpBin) < maxLen && len(tmpBin) < 4 {
        tmpBin = append(tmpBin, "0"...)
    }
    src += string(tmpBin)
    return src

}
func convertToCodeLang (src string, codeLen int) (result []int) {
    //得られたデータを8bit区切りのコード語に変換する
    srcLen := len(src)
    /*

    srcを8bitに区切った際、最後のbit列が8bitに満たない場合
    8bitになるようbitを追加する
    */
    if needLen := 8 - (srcLen % 8); needLen > 0 && needLen < 8 {
        tmpBin := make([]byte, 0, needLen)
        for len(tmpBin) < needLen {
            tmpBin = append(tmpBin, "0"...)
            srcLen++
        }
        src += string(tmpBin)
    }
    /*
    また、8bitの数がシンボル容量に満たない場合
    11101100と00010001を交互に追加する
    */
    for i := srcLen / 8; i < codeLen; i++ {
        if i % 2 == 0 {
            src += "11101100"
        } else {
            src += "00010001"
        }
        srcLen += 8
    }
    //8bitごとに整数に変換して終了
    result = make([]int, codeLen)
    for i := 0; i < codeLen; i++ {
        result[i] = binToInt(src[i * 8:i * 8 + 8])
    }
    return
}
