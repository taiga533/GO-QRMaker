package QRMaker

import (
    "log"
)

type codeBlock struct {
    dataCode []int
    ecc []int
    eccCnt int
}

var qrInfo *qrCodeInfo
func SimplyEncode (data, fileName string, ecLv,outputSize int) (err error) {
    qrInfo,err = newQrInfo(data, ecLv)
    if err != nil {
        return
    }
    log.Printf("QrInfo%+v\n", qrInfo)
    dc := newDataCode(data, qrInfo.dataCodeCnt, qrInfo.modeIndicator)
    code := make([]int, 0)
    code = append(code, makeCodeBlock(dc.data)...)
    qrImg := newQrImg(code, qrInfo.imgSize)
    err = qrImg.arrayWriteToImg(fileName, outputSize)
    if err != nil {
        return
    }
    return
}
func makeCodeBlock (dataCode []int) []int {
    codeBlocks := make([]codeBlock, 0)
    var dcIdx, blockIdx int
    for _,block := range versions[qrInfo.versionIdx].blocks {
        for dcIdx < block[2] * block[0] + blockIdx {
            tmpCb := new(codeBlock)
            tmpCb.dataCode = make([]int, len(dataCode[dcIdx:dcIdx + block[2]]))
            copy(tmpCb.dataCode, dataCode[dcIdx:dcIdx + block[2]])
            tmpCb.eccCnt = block[1] - block[2]
            codeBlocks = append(codeBlocks, *tmpCb)
            dcIdx += block[2]
        }
        blockIdx += dcIdx
    }
    for i,_ := range codeBlocks {
        ecc := newECC(codeBlocks[i].dataCode, codeBlocks[i].eccCnt)
        codeBlocks[i].ecc = make([]int, 0, codeBlocks[i].eccCnt)
        codeBlocks[i].ecc = append(codeBlocks[i].ecc, ecc.ecc...)
        log.Printf("dc[%v].dc: %v\n", i, codeBlocks[i].dataCode)
        log.Printf("dc[%v].ecc: %v\n\n", i, codeBlocks[i].ecc)
    }
    //log.Printf("dcblocks%#v\n", codeBlocks)
    var shapedCode []int
    shapedCode = replaceCode(codeBlocks)
    log.Printf("インタリーブ配置後のデータコード列%v\n", shapedCode)
    return shapedCode
}
func replaceCode (codeBlocks []codeBlock) []int {
    var maxEccCnt, maxDcCnt int
    for _,cb := range codeBlocks {
        if len(cb.dataCode) > maxDcCnt {
            maxDcCnt = len(cb.dataCode)
        }
        if cb.eccCnt > maxEccCnt {
            maxEccCnt = cb.eccCnt
        }
    }
    shapedCode := make([]int, 0)

    for dcIdx := 0; dcIdx < maxDcCnt; dcIdx++ {
        for _,cb := range codeBlocks {
            if dcIdx < len(cb.dataCode) {
                shapedCode = append(shapedCode, cb.dataCode[dcIdx])
            }
        }
    }

    for eccIdx := 0; eccIdx < maxEccCnt; eccIdx++ {
        for _,cb := range codeBlocks {
            if eccIdx < cb.eccCnt {
                shapedCode = append(shapedCode, cb.ecc[eccIdx])
            }
        }
    }
    return shapedCode
}
