package QRMaker

import (
    "testing"

)

func TestByteMode(t *testing.T) {
    err := SimplyEncode("123testByteMode%$", "./testByteMode", 4, 2)
    if err != nil {
        t.Fatalf("バイトモードのテストに失敗。 %v", err)
    }
}
func TestNumMode(t *testing.T) {
    err := SimplyEncode("0123456", "./testNumMode", 4, 2)
    if err != nil {
        t.Fatalf("数値モードのテストに失敗。 %v", err)
    }
}
func TestAlphaNumMode(t *testing.T) {
    err := SimplyEncode("012346asdfa%:54", "./testAlphaNumMode", 4, 2)
    if err != nil {
        t.Fatalf("英数字モードのテストに失敗。 %v", err)
    }
}

func TestECCLevelLow(t *testing.T) {
    err := SimplyEncode("123456", "./eccLow", 1, 2)
    if err != nil {
        t.Errorf("ECCレベル低のテストに失敗。 %v", err)
    }
}
func TestECCLevelMedium(t *testing.T) {
    err := SimplyEncode("123456", "./eccMedium", 2, 2)
    if err != nil {
        t.Errorf("ECCレベル中のテストに失敗。 %v", err)
    }
}
func TestECCLevelHigh(t *testing.T) {
    err := SimplyEncode("123456", "./eccHigh", 3, 2)
    if err != nil {
        t.Errorf("ECCレベル高のテストに失敗。 %v", err)
    }
}
func TestECCLevelHighest(t *testing.T) {
    err := SimplyEncode("123456", "./eccHighest", 4, 2)
    if err != nil {
        t.Errorf("ECCレベル最高のテストに失敗。 %v", err)
    }
}
func TestWorngECCLevel(t *testing.T) {
    err := SimplyEncode("123456", "./WorngECCLevelTest", 5, 1)
    err = SimplyEncode("123456", "./WorngECCLevelTest", 0, 1)
    if err == nil {
        t.Errorf("間違ったECCレベルの値を渡した際のテストに失敗しました %v", err)
    }
}
func TestWorngFilePath(t *testing.T) {
    err := SimplyEncode("123456", "./ldskflkaslfknjsojlnlnlkjsldfkj/test", 4, 1)
    if err == nil {
        t.Errorf("間違ったファイルパスを渡した際のテストに失敗しました %v", err)
    }
}
