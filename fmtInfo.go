package QRMaker

import (
    "log"
)

func getFmtInfoBin (mskPttrnBin string) (fmtInfoBin string) {
    ecLvBin := intToBin(qrInfo.ecLv, 2)
    mskBin := mskPttrnBin
    log.Printf("エラー訂正子%v、マスク識別子%v", ecLvBin, mskBin)
    gx := "10100110111"
    fmtInfoBin = ecLvBin + mskBin + "0000000000"
    bitPos := 0
    for bitPos < 5 {
        if rune(fmtInfoBin[bitPos]) == rune('1') {
            tmpGx := BinShiftToLeft(gx, 4 - bitPos)
            fmtInfoBin = binXor(fmtInfoBin, tmpGx)
        }
        bitPos++
    }
    fmtInfoBin = fmtInfoBin[5:]
    fmtInfoBin = ecLvBin + mskBin + fmtInfoBin
    log.Printf("prev形式情報コード%v\n", fmtInfoBin)
    fmtInfoBin = binXor(fmtInfoBin, "101010000010010")
    log.Printf("形式情報コード%v\n", fmtInfoBin)
    return
}
func calcVersionInfoBin (version int) (fmtInfoBin string) {
    // バージョン7以降で必要になるバージョン情報のモジュールの内容を計算するものです。
    versionBin := intToBin(version, 6)
    log.Printf("バージョン子%v", versionBin)
    gx := "1111100100101"
    fmtInfoBin = versionBin + "000000000000"
    bitPos := 0
    for bitPos < 6 {
        if rune(fmtInfoBin[bitPos]) == rune('1') {
            tmpGx := BinShiftToLeft(gx, 5 - bitPos)
            fmtInfoBin = binXor(fmtInfoBin, tmpGx)
        }
        bitPos++
    }
    fmtInfoBin = fmtInfoBin[6:]
    fmtInfoBin =versionBin + fmtInfoBin
    log.Printf("バージョン情報コード%v\n", fmtInfoBin)
    return
}


func binXor (left string, right string) (result string) {
    leftInt := binToInt(left)
    rightInt := binToInt(right)
    result = intToBin(leftInt ^ rightInt, len(left))
    return
}
func BinShiftToLeft (bin string, shiftCnt int) (result string) {
    //2進数を左にshiftCnt回シフトする関数です。
    tmpBin := make([]byte, 0, len(bin) + shiftCnt)
    tmpBin = append(tmpBin, bin...)
    for len(tmpBin) < cap(tmpBin) {
        tmpBin = append(tmpBin, "0"...)
    }
    result = string(tmpBin)
    return
}
