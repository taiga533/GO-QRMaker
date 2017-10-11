package QRMaker

import (
)

type errCheckCode struct {
    eccWordCnt int
    fx []int
    gx []int
    ecc []int
}

func newECC (dataCodes []int, eccWordCnt int) (ecc *errCheckCode) {
    ecc = new(errCheckCode)
    ecc.eccWordCnt = eccWordCnt
    if eccWordCnt > len(dataCodes) {
        ecc.fx = make([]int, eccWordCnt)
    } else {
        ecc.fx = make([]int, len(dataCodes))
    }
    //１誤り訂正コード = gx2(x)なので引数を１減算して渡す
    ecc.eccWordCnt--
    ecc.calcGx()
    divideCnt := len(dataCodes)
    copy(ecc.fx, dataCodes)
    ecc.divideFxByGx(divideCnt)
    ecc.ecc = ecc.fx
    return
}
func (ecc *errCheckCode) calcGx () {
    addX := []int{0, 0}
    addY := []int{0, 1}
    for i := 0; i < ecc.eccWordCnt; i++ {
        addResult := addGx(addX, addY)
        ecc.gx = make([]int, len(addX) + 1)
        ecc.gx[0] = addResult[0][0]
        copy(ecc.gx[1:], xorGx(addResult))
        ecc.gx[len(addX)] = addResult[1][len(addResult[1]) - 1]
        addX = ecc.gx
        addY = []int{0, addY[1] + 1}
    }
}
func addGx (addX []int, addY []int) (addResult [][]int) {
    addResult = make([][]int, 2)

    for yIndex := 0; yIndex < 2; yIndex++ {
        tmpRow := make([]int, len(addX))
        for xIndex := 0; xIndex < len(addX); xIndex++ {
            tmpRow[xIndex] = addY[yIndex] + addX[xIndex]
            if tmpRow[xIndex] > 255 {
                tmpRow[xIndex] -= 255
            }
        }
        addResult[yIndex] = tmpRow
    }

    return
}
func xorGx (addResult [][]int) (xorResult []int) {
    altable := alphaLookupTable
    xorResult = make([]int, 0, len(addResult[0]) - 1)
    for xIndex := 1; xIndex < len(addResult[0]); xIndex++ {
        xorTemp := altable[addResult[0][xIndex]] ^ altable[addResult[1][xIndex - 1]]
        if xorTemp > 255 {
            xorTemp -= 255
        }
        xorResult = append(xorResult, searchIntArray(altable, xorTemp))
    }
    return
}
func (ecc *errCheckCode) divideFxByGx (divideCnt int) {
    altable := alphaLookupTable
    tmpFx := make([]int, len(ecc.fx))
    copy(tmpFx, ecc.fx)
    for i := 0; i < divideCnt; i++ {
        fxAlpha := searchIntArray(altable, tmpFx[0])
        tmpGx := make([]int, len(ecc.gx) - 1)
        copy(tmpGx, ecc.gx[1:])
        for i,v := range tmpGx {
            tmpValue := v + fxAlpha
            if tmpValue > 255 {
                tmpValue -= 255
            }
            tmpGx[i] = altable[tmpValue]
        }
        copy(tmpFx, tmpFx[1:])
        tmpFx[len(tmpFx) - 1] = 0
        for i,v := range tmpGx {
            tmpFx[i] = tmpFx[i] ^ v
        }
    }
    ecc.fx = make([]int, 0, len(ecc.gx))
    ecc.fx = append(ecc.fx, tmpFx[:len(ecc.gx) - 1]...)
}

/*
errCheckCode計算時のg(x)を求めるのに用いる
aの乗数と整数の対応表
*/
var alphaLookupTable = []int{
    1,   2,   4,   8,  16,  32,  64, 128,
    29,  58, 116, 232, 205, 135,  19,  38,
    76, 152,  45,  90, 180, 117, 234, 201,
    143,   3,   6,  12,  24,  48,  96, 192,
    157,  39,  78, 156,  37,  74, 148,  53,
    106, 212, 181, 119, 238, 193, 159,  35,
    70, 140,   5,  10,  20,  40,  80, 160,
    93, 186, 105, 210, 185, 111, 222, 161,
    95, 190,  97, 194, 153,  47,  94, 188,
    101, 202, 137,  15,  30,  60, 120, 240,
    253, 231, 211, 187, 107, 214, 177, 127,
    254, 225, 223, 163,  91, 182, 113, 226,
    217, 175,  67, 134,  17,  34,  68, 136,
    13,  26,  52, 104, 208, 189, 103, 206,
    129,  31,  62, 124, 248, 237, 199, 147,
    59, 118, 236, 197, 151,  51, 102, 204,
    133,  23,  46,  92, 184, 109, 218, 169,
    79, 158,  33,  66, 132,  21,  42,  84,
    168,  77, 154,  41,  82, 164,  85, 170,
    73, 146,  57, 114, 228, 213, 183, 115,
    230, 209, 191,  99, 198, 145,  63, 126,
    252, 229, 215, 179, 123, 246, 241, 255,
    227, 219, 171,  75, 150,  49,  98, 196,
    149,  55, 110, 220, 165,  87, 174,  65,
    130,  25,  50, 100, 200, 141,   7,  14,
    28,  56, 112, 224, 221, 167,  83, 166,
    81, 162,  89, 178, 121, 242, 249, 239,
    195, 155,  43,  86, 172,  69, 138,   9,
    18,  36,  72, 144,  61, 122, 244, 245,
    247, 243, 251, 235, 203, 139,  11,  22,
    44,  88, 176, 125, 250, 233, 207, 131,
    27,  54, 108, 216, 173,  71, 142,   1,
}
