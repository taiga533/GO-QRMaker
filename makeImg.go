package QRMaker

import (
    "image"
    "image/color"
    "os"
    "image/png"
)

type qrImg struct {
    imgSize int
    imgArray [][]int
    mskNumBin string
    img *image.NRGBA
}

func newQrImg (data []int, imgSize int) (qri *qrImg) {
    qri = new(qrImg)
    //インデックスの最大は各軸のピクセル数より1小さいので減算する
    qri.imgSize = imgSize
    qri.imgArray = make([][]int, 0, qri.imgSize)
    for len(qri.imgArray) < cap(qri.imgArray) {
        qri.imgArray = append(qri.imgArray, make([]int, qri.imgSize))
    }

    qri.imgSize--
    qri.setPosPttrns()
    qri.setTmpFmtPttrns()
    qri.setTimingPttrns()
    if qrInfo.version > 1 {
        qri.setAlignmentPttrn()
    }

    qri.setDataToArray(data)
    qri.mskNumBin = getMaskedData(qri.imgArray, qri.imgSize + 1)
    fmtInfoBin := getFmtInfoBin(qri.mskNumBin)
    qri.setFmtPttrn(fmtInfoBin)
    return
}
func (qri *qrImg) arrayWriteToImg (fileName string, outputSize int) error {
    qri.imgSize++
    qri.imgSize = (qri.imgSize * outputSize) + (8 * outputSize)
    qri.img = image.NewNRGBA(image.Rect(0, 0, qri.imgSize, qri.imgSize))
    f, err := os.Create(fileName + ".png")
    if err != nil {
        return err
    }
    //余白を生成する
    for p := 0; p < qri.imgSize; p++ {
        for lineNum := 0; lineNum <= 4 * outputSize; lineNum++ {
            qri.img.Set(p, lineNum, color.NRGBA{255, 255, 255, 255})
            qri.img.Set(lineNum, p, color.NRGBA{255, 255, 255, 255})
            qri.img.Set(qri.imgSize - p, qri.imgSize - lineNum, color.NRGBA{255, 255, 255, 255})
            qri.img.Set(qri.imgSize - lineNum, qri.imgSize - p, color.NRGBA{255, 255, 255, 255})
        }
    }
    //サイズに合わせてQRコードを出力する
    outside := 4 * outputSize
    for dataY,_ := range qri.imgArray {
        for dataX,_ := range qri.imgArray {
            for addY := dataY * outputSize; addY < dataY * outputSize + outputSize; addY++ {
                for addX := dataX * outputSize; addX < dataX * outputSize + outputSize; addX++ {
                    if qri.imgArray[dataX][dataY] == 1 {
                        qri.img.Set(outside + addX, outside + addY, color.NRGBA{0, 0, 0, 255})
                    } else {
                        qri.img.Set(outside + addX, outside + addY, color.NRGBA{255, 255, 255, 255})
                    }
                }
            }
        }
    }
    defer f.Close()
    err = png.Encode(f, qri.img)
    return err
}
func (qri *qrImg) setDataToArray (data []int) {
    //スライスにデータビット列と誤り訂正ビット列を記録します
    tmpByte := make([]byte, 0, len(data) * 8)
    for _,v := range data {
        tmpByte = append(tmpByte, intToBin(v, 8)...)
    }
    dataIdx := 0
    direction := -1
    dataStr := string(tmpByte)
    y := qri.imgSize
    for x := qri.imgSize - 1; x >= 0; x -= 2 {
        for (y >= 0) && (y <= qri.imgSize) {
            if dataIdx >= len(dataStr) - 1 {
                break
            }
            switch {
            case (qri.imgArray[x][y] > 1) && (qri.imgArray[x + 1][y] > 1):
                y += direction
            case qri.imgArray[x + 1][y] > 1:
                qri.imgArray[x][y] = binToInt(string(dataStr[dataIdx]))
                dataIdx++
                y += direction
            default:
                qri.imgArray[x + 1][y] = binToInt(string(dataStr[dataIdx]))
                qri.imgArray[x][y] = binToInt(string(dataStr[dataIdx + 1]))
                dataIdx += 2
                y += direction
            }
        }
        if x == 7 {
            x--
        }
        direction = -direction
        y += direction
    }
}

func (qri *qrImg) setPosPttrns () {
    //位置合わせパターンをスライスに記録します。
    posPttrn := [][]int {
        {3, 3, 3, 3, 3, 3, 3},
        {3, 2, 2, 2, 2, 2, 3},
        {3, 2, 3, 3, 3, 2, 3},
        {3, 2, 3, 3, 3, 2, 3},
        {3, 2, 3, 3, 3, 2, 3},
        {3, 2, 2, 2, 2, 2, 3},
        {3, 3, 3, 3, 3, 3, 3},
    }
    for y := 0; y < len(posPttrn); y++ {
        for x := 0; x < len(posPttrn); x++ {
            qri.imgArray[x][y] = posPttrn[x][y]
            qri.imgArray[qri.imgSize - x][y] = posPttrn[x][y]
            qri.imgArray[x][qri.imgSize - y] = posPttrn[x][y]
        }
    }
    //位置あわせパターンの中央向きの空白領域をスライスに記録します。
    for p := 0; p < 8; p++ {
        //左上
        qri.imgArray[7][p] = 2
        qri.imgArray[p][7] = 2
        //右上
        qri.imgArray[qri.imgSize - p][7] = 2
        qri.imgArray[qri.imgSize - 7][p] = 2
        //左下
        qri.imgArray[p][qri.imgSize - 7] = 2
        qri.imgArray[7][qri.imgSize - p] = 2
    }
}
func (qri *qrImg) setTmpFmtPttrns () {
    /*
    形式情報パターンの領域を配列に書き込みます。
    マスク処理が完了するまで形式情報パターンを書き込むことができませんが
    データを書き込まない領域として認識させなければならないので配列に値を入れます。
    */
    for p := 0; p < 8; p++ {
        //左上
        qri.imgArray[p][8] = 2
        qri.imgArray[8][p] = 2
        qri.imgArray[8][8] = 2
        //右上
        qri.imgArray[qri.imgSize - p][8] = 2
        //左下
        qri.imgArray[8][qri.imgSize - p] = 2
    }
    //左下の点を書き込みます
    qri.imgArray[8][qri.imgSize - 7] = 3
    if qrInfo.version > 6 {
        for y := qri.imgSize - 10; y <= qri.imgSize - 8; y++ {
            for x := 0; x <= 5; x++ {
                qri.imgArray[y][x] = 2
                qri.imgArray[x][y] = 2
            }
        }
    }
}
func (qri *qrImg) setTimingPttrns () {
    //タイミングパターンを描写します。
    for p := 7; qri.imgArray[p][6] != 3; p++ {
        qri.imgArray[6][p] = 3 - p % 2
        qri.imgArray[p][6] = 3 - p % 2
    }
}
func (qri *qrImg) setFmtPttrn (fmtInfoBin string) {
    //形式情報をスライスの形式情報領域に記録します。
    p := 0
    qri.imgArray[8][8] = binToInt(string(fmtInfoBin[7]))
    qri.imgArray[qri.imgSize - 7][7] = binToInt(string(fmtInfoBin[7]))
    for i := 0; i < 7; i++ {
        //左上縦
        qri.imgArray[8][7 - i - p] = binToInt(string(fmtInfoBin[8 + i]))
        //左上横
        qri.imgArray[7 - i - p][8] = binToInt(string(fmtInfoBin[6 - i]))
        //右上
        qri.imgArray[qri.imgSize - i][8] = binToInt(string(fmtInfoBin[14 - i]))
        //左下
        qri.imgArray[8][qri.imgSize - i] = binToInt(string(fmtInfoBin[i]))
        //右上の形式情報パターンだけタイミングパターンに触れるのでそれを回避する
        if i == 0 {
            p++
        }
    }
    if qrInfo.version > 6 {
        versionInfo := calcVersionInfoBin(qrInfo.version)
        versionInfoIdx := 0
        for x := 5; x >= 0; x-- {
            for y := qri.imgSize - 8; y >= qri.imgSize - 10; y-- {
                qri.imgArray[y][x] = binToInt(string(versionInfo[versionInfoIdx]))
                qri.imgArray[x][y] = binToInt(string(versionInfo[versionInfoIdx]))
                versionInfoIdx++
            }
        }
    }
}
func (qri *qrImg) setAlignmentPttrn () {
    //軸合わせパターンを描写します。
    alignPttrn := [][]int {
        {3, 3, 3, 3, 3},
        {3, 2, 2, 2, 3},
        {3, 2, 3, 2, 3},
        {3, 2, 2, 2, 3},
        {3, 3, 3, 3, 3},
    }
    if qrInfo.version > 1 {
        var alignPos [][]int
        alignPos = make([][]int, 0, len(qrInfo.alignment) * len(qrInfo.alignment) - 3)
        for x := 0; x < len(qrInfo.alignment); x++ {
            for y := 0; y < len(qrInfo.alignment); y++ {
                switch {
                case x == 0 && y == 0:
                    continue
                case (x == 0 && y == len(qrInfo.alignment) - 1) || (x == len(qrInfo.alignment) - 1 && y == 0):
                    continue
                default:
                    alignPos = append(alignPos, []int{qrInfo.alignment[x] - 2, qrInfo.alignment[y] - 2})
                }
            }
        }
        for _,ap := range alignPos {
            for y := 0; y < 5; y++ {
                for x := 0; x < 5;x++ {
                    qri.imgArray[ap[0] + x][ap[1] + y] = alignPttrn[x][y]
                }
            }
        }
    }
}
