package QRMaker

import (
    "errors"
)
//
const low = 1
const medium = 0
const high = 3
const highest = 2
type dataNumIndicator struct {
    num int
    alphaNum int
    bytemode int
}
var dataNumIndicator1to9 = [3]int {
    10,
    9,
    8,
}
var dataNumIndicator10to26 = [3]int {
    12,
    11,
    16,
}
var dataNumIndicator27to40 = [3]int {
    14,
    13,
    16,
}
var alignment = [][]int {
    //バージョン１は軸あわせパターンが無い
    {},
    {6, 18},
    {6, 22},
    {6, 26},
    {6, 30},
    {6, 34},
    {6, 22, 38},
}
type qrCodeInfo struct {
    version int
    imgSize int
    ecLv int
    dataCodeCnt int
    alignment []int
    versionIdx int
    modeIndicator int
}
type version struct {
    version int
    dataNumIndicator [3]int
    ecLevel int
    blocks [][]int
}
var versions = []version{
    {
        1,
        dataNumIndicator1to9,
        low,
        [][]int{
            {
                1,
                26,
                19,
            },
        },
    },
    {
        1,
        dataNumIndicator1to9,
        medium,
        [][]int{
            {
                1,
                26,
                16,
            },
        },
    },
    {
        1,
        dataNumIndicator1to9,
        high,
        [][]int{
            {
                1,
                26,
                13,
            },
        },
    },
    {
        1,
        dataNumIndicator1to9,
        highest,
        [][]int{
            {
                1,
                26,
                9,
            },
        },
    },
    {
        2,
        dataNumIndicator1to9,
        low,
        [][]int{
            {
                1,
                44,
                34,
            },
        },
    },
    {
        2,
        dataNumIndicator1to9,
        medium,
        [][]int{
            {
                1,
                44,
                28,
            },
        },
    },
    {
        2,
        dataNumIndicator1to9,
        high,
        [][]int{
            {
                1,
                44,
                22,
            },
        },
    },
    {
        2,
        dataNumIndicator1to9,
        highest,
        [][]int{
            {
                1,
                44,
                16,
            },
        },
    },
    {
        3,
        dataNumIndicator1to9,
        low,
        [][]int{
            {
                1,
                70,
                55,
            },
        },
    },
    {
        3,
        dataNumIndicator1to9,
        medium,
        [][]int{
            {
                1,
                70,
                44,
            },
        },
    },
    {
        3,
        dataNumIndicator1to9,
        high,
        [][]int{
            {
                2,
                35,
                17,
            },
        },
    },
    {
        3,
        dataNumIndicator1to9,
        highest,
        [][]int{
            {
                2,
                35,
                13,
            },
        },
    },
    {
        4,
        dataNumIndicator1to9,
        low,
        [][]int{
            {
                1,
                100,
                80,
            },
        },
    },
    {
        4,
        dataNumIndicator1to9,
        medium,
        [][]int{
            {
                2,
                50,
                32,
            },
        },
    },
    {
        4,
        dataNumIndicator1to9,
        high,
        [][]int{
            {
                2,
                50,
                24,
            },
        },
    },
    {
        4,
        dataNumIndicator1to9,
        highest,
        [][]int{
            {
                4,
                25,
                9,
            },
        },
    },
    {
        5,
        dataNumIndicator1to9,
        low,
        [][]int{
            {
                1,
                134,
                108,
            },
        },
    },
    {
        5,
        dataNumIndicator1to9,
        medium,
        [][]int{
            {
                2,
                67,
                43,
            },
        },
    },
    {
        5,
        dataNumIndicator1to9,
        high,
        [][]int{
            {
                2,
                33,
                15,
            },
            {
                2,
                34,
                16,
            },
        },
    },
    {
        5,
        dataNumIndicator1to9,
        highest,
        [][]int{
            {
                2,
                33,
                11,
            },
            {
                2,
                34,
                12,
            },
        },
    },
    {
        6,
        dataNumIndicator1to9,
        low,
        [][]int{
            {
                2,
                86,
                68,
            },
        },
    },
    {
        6,
        dataNumIndicator1to9,
        medium,
        [][]int{
            {
                4,
                43,
                27,
            },
        },
    },
    {
        6,
        dataNumIndicator1to9,
        high,
        [][]int{
            {
                4,
                43,
                19,
            },
        },
    },
    {
        6,
        dataNumIndicator1to9,
        highest,
        [][]int{
            {
                4,
                43,
                15,
            },
        },
    },
    {
        7,
        dataNumIndicator1to9,
        low,
        [][]int{
            {
                2,
                98,
                78,
            },
        },
    },
    {
        7,
        dataNumIndicator1to9,
        medium,
        [][]int{
            {
                4,
                49,
                31,
            },
        },
    },
    {
        7,
        dataNumIndicator1to9,
        high,
        [][]int{
            {
                2,
                32,
                14,
            },
            {
                4,
                33,
                15,
            },
        },
    },
    {
        7,
        dataNumIndicator1to9,
        highest,
        [][]int{
            {
                4,
                39,
                13,
            },
            {
                1,
                40,
                14,
            },
        },
    },
}
func newQrInfo (srcStr string, ecLv int) (qrInfo *qrCodeInfo, err error) {
    qrInfo = new(qrCodeInfo)
    if ecLv > 4 || ecLv < 1 {
        err = errors.New("ErrorCheckLevel value is invalid")
        return
    }
    switch {
    case checkRegexp(`^[0-9]+$`, srcStr):
        qrInfo.modeIndicator = 1
        //数値は１バイト文字なのでlenで数える
        strBitCnt := (len(srcStr) / 3) * 10
        switch len(srcStr) % 3 {
        case 1:
            strBitCnt += 4
        case 2:
            strBitCnt += 7
        }
        qrInfo.version, qrInfo.versionIdx, qrInfo.dataCodeCnt,err = setQrVersion(strBitCnt, ecLv, 0)
    case checkRegexp(`^[A-Z0-9+-.:*\s%$/]+$`, srcStr):
        qrInfo.modeIndicator = 2
        strCnt := getStrLen(srcStr)
        strBitCnt := (strCnt / 2) * 11 + (strCnt % 2) * 6
        qrInfo.version, qrInfo.versionIdx, qrInfo.dataCodeCnt,err = setQrVersion(strBitCnt, ecLv, 1)
    default :
        qrInfo.modeIndicator = 4
        //１コード語1バイトなのでマルチバイト対応の文字数計算ではなくlen()で文字数を数える
        strBitCnt := len(srcStr) * 8
        qrInfo.version, qrInfo.versionIdx, qrInfo.dataCodeCnt,err = setQrVersion(strBitCnt, ecLv, 2)
    }
    qrInfo.imgSize = 21 + 4 * (qrInfo.version - 1)
    qrInfo.alignment = alignment[qrInfo.version - 1]
    switch ecLv {
    case 1:
        qrInfo.ecLv = 1
    case 2:
        qrInfo.ecLv = 0
    case 3:
        qrInfo.ecLv = 3
    case 4:
        qrInfo.ecLv = 2
    }
    return
}
func setQrVersion (srcBitCnt, ecLv, dataMode int) (version, versionIdx, maxDataCnt int, err error) {
    var dataCodeCnt int
    for i := ecLv - 1; i < len(versions); i += 4 {
        tmpSrcBitCnt := (srcBitCnt + 4 + versions[i].dataNumIndicator[dataMode])
        dataCodeCnt = divideRoundUp(tmpSrcBitCnt, 8)
        maxDataCnt = 0
        for _,block := range versions[i].blocks {
            maxDataCnt += block[2] * block[0]
        }
        version = versions[i].version
        if dataCodeCnt <= maxDataCnt {
            versionIdx = i
            break
        }
    }
    if dataCodeCnt > maxDataCnt {
        err = errors.New("Input data is too long, please shorten it")
    }
    return
}
func tomato () (err error) {
    setQrVersion(143, 1, 1)
    return
}
