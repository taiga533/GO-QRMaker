package QRMaker

import (
)
type mask struct {
    mPttrnBin string
    bestMask [][]int
    bestPscore int
    mskSize int
}
func getMaskedData (srcArray [][]int, mskSize int) (pttrnBin string) {
    var mskedArray [][]int
    m := new(mask)
    m.bestMask = make([][]int, 0, mskSize)
    m.mskSize = mskSize
    mskedArray = make([][]int, 0, mskSize)
    for len(mskedArray) < cap(mskedArray) {
        mskedArray = append(mskedArray, make([]int, mskSize))
        m.bestMask = append(m.bestMask, make([]int, mskSize))
    }
    for mskNum := 0; mskNum < 8; mskNum++ {
        mskedArray, pttrnBin = calcMaskPttrn(srcArray, mskSize, mskNum)
        m.evaluateMsk(mskedArray, pttrnBin)
    }
    pttrnBin = m.mPttrnBin
    copy(srcArray, m.bestMask)
    return
}
func makeCloneMap (srcMap, newMap [][]int) [][]int {
    //2次元スライスはcopy(srcMap, newMap)だけではコピーできないので作った関数
    for i,_ := range newMap {
        copy(newMap[i], srcMap[i])
    }
    return newMap
}
func (m *mask)evaluateMsk (mskedArray [][]int, pttrnBin string) (pScore int) {
    pScore = penalty1(mskedArray)
    pScore += penalty2(mskedArray, m.mskSize)
    pScore += penalty4(mskedArray, m.mskSize)
    if pScore < m.bestPscore || m.bestPscore == 0 {
        m.bestMask = mskedArray
        m.bestPscore = pScore
        m.mPttrnBin = pttrnBin
    }
    return
}
func penalty1 (mskedArray [][]int) (penaltyScore int) {
    // 縦軸と横軸で黒or白が5つ以上続いた場合に減点
    // 黒が7つ続いたら7-5で2になりこれを3倍(重み)した数分減点
    var xPenaltyCnt,yPenaltyCnt int
    var xBefore = - 1
    var yBefore = - 1
    for y,_ := range mskedArray {
        for x,_ := range mskedArray {
            if mskedArray[y][x] == xBefore {
                xPenaltyCnt++
            } else {
                if xPenaltyCnt > 5 {
                    penaltyScore += 3 + (xPenaltyCnt - 5)
                } else if xPenaltyCnt > 4 {
                    penaltyScore += 3
                }
                xPenaltyCnt = 0
            }
            xBefore = mskedArray[y][x]
            if mskedArray[x][y] == yBefore {
                yPenaltyCnt++
            } else {
                if yPenaltyCnt > 5 {
                    penaltyScore += 3 + (yPenaltyCnt - 5)
                } else if yPenaltyCnt > 4 {
                    penaltyScore += 3
                }
                yPenaltyCnt = 0
            }
            yBefore = mskedArray[x][y]

        }
    }
    return
}
func penalty2 (mskedArray [][]int, mskSize int) (penaltyScore int) {
    // 2x2の黒または白のブロックが存在した場合は減点
    // ここで3x3の黒または白のブロックには4つの2x2のブロックが存在するとみなす
    for y := 0; y < mskSize - 1; y++ {
        for x := 0; x < mskSize - 1; x++ {
            xy := mskedArray[y][x]
            if xy == mskedArray[y][x + 1] && xy == mskedArray[y + 1][x] && xy == mskedArray[y + 1][x + 1] {
                penaltyScore++
            }
        }
    }
    penaltyScore = penaltyScore*3
    return
}
func penalty4 (mskedArray [][]int, mskSize int) (penaltyScore int) {
    // 黒の割合が50%から5%ごとに増減することにより10点ずつ減点
    var penaltyCnt int
    for y := 0; y < mskSize - 1; y++ {
        for x := 0; x < mskSize - 1; x++ {
            xy := mskedArray[y][x]
            if xy == 1 {
                penaltyCnt++
            }
        }
    }
    penaltyScore = 50 - (penaltyCnt * 100 / mskSize) / 5
    if penaltyScore < 0 {
        penaltyScore = -penaltyScore
    }
    penaltyScore = penaltyScore*2
    return
}
func calcMaskPttrn (srcArray [][]int, mskSize int, mskNum int) (mskPttrnArray [][]int, pttrnBin string) {
    mskPttrnArray = make([][]int, 0, mskSize)
    pttrnBin = intToBin(mskNum, 3)
    for len(mskPttrnArray) < cap(mskPttrnArray) {
        mskPttrnArray = append(mskPttrnArray, make([]int, mskSize))
    }
    makeCloneMap(srcArray, mskPttrnArray)
    var mskPttrnBool bool
    for y,_ := range mskPttrnArray {
        for x,_ := range mskPttrnArray {
            switch mskNum {
            case 0:
                mskPttrnBool = (mskPttrnArray[x][y] < 2) && (((x + y) % 2) == 0)
            case 1:
                mskPttrnBool = (mskPttrnArray[x][y] < 2) && ((y % 2) == 0)
            case 2:
                mskPttrnBool = (mskPttrnArray[x][y] < 2) && ((x % 3) == 0)
            case 3:
                mskPttrnBool = (mskPttrnArray[x][y] < 2) && (((x + y) % 3) == 0)
            case 4:
                mskPttrnBool = (mskPttrnArray[x][y] < 2) && (((y / 2) + (x / 3)) % 2 == 0)
            case 5:
                mskPttrnBool = (mskPttrnArray[x][y] < 2) && ((x*y) % 2 + (x*y) % 3 == 0)
            case 6:
                mskPttrnBool = (mskPttrnArray[x][y] < 2) && (((x*y) % 2 + (x*y) % 3) % 2 == 0)
            case 7:
                mskPttrnBool = (mskPttrnArray[x][y] < 2) && (((x*y) % 3 + (x+y) % 2) % 2 == 0)
            }
            if mskPttrnBool {
                //ビットを反転させる
                mskPttrnArray[x][y] = 1 - mskPttrnArray[x][y]
            } else if mskPttrnArray[x][y] % 2 == 1 {
                // データ書き込み時に特定のパターンに対するビットの書き込みを拒否するのに
                // 使っていた判別用の3と2を1と0に直す
                mskPttrnArray[x][y] = 1
            } else {
                mskPttrnArray[x][y] = 0
            }
        }
    }
    return
}
