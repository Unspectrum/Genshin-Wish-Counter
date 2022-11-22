package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Unspectrum/Genshin-Wish-Counter/models"
	"github.com/Unspectrum/Genshin-Wish-Counter/utils"

	"github.com/xuri/excelize/v2"
)

func check(e error) {
	if e != nil {
		panic(e)
	}

}

const (
	logLocation      = "\\AppData\\LocalLow\\miHoYo\\Genshin Impact\\output_log.txt"
	dataFileLocation = "webCaches/Cache/Cache_Data/data_2"
	warmUpStr        = "Warmup file "
	streamAssetsStr  = "StreamingAssets"
)

var API_URL = "https://hk4e-api-os.hoyoverse.com/event/gacha_info/api/getGachaLog"

func GetInstallLocation(str string) string {
	str = strings.ReplaceAll(str, warmUpStr, "")
	str = strings.Split(str, "\\")[0]
	return strings.ReplaceAll(str, streamAssetsStr, "")
}

func main() {
	// apiHost := "hk4e-api-os.hoyoverse.com"
	dir, err := os.UserHomeDir()
	check(err)
	log := dir + logLocation
	lines, err := utils.OpenFileToStringArray(log, "\n")
	check(err)
	// var temp string
	var installLocation string
	for _, line := range lines {
		if strings.Contains(line, warmUpStr) {
			installLocation = GetInstallLocation(line)
			break
		}
	}
	installLocation += dataFileLocation
	DataOpen, err := utils.OpenReadFileToString(installLocation)
	if err != nil {
		panic(err)
	}

	dataSplit := strings.Split(string(DataOpen), "1/0")
	var exp = "(authkey=.+?game_biz=)"
	re := regexp.MustCompile(exp)
	match := re.FindStringSubmatch(dataSplit[len(dataSplit)-1])[0]
	match = strings.ReplaceAll(match, "&game_biz=", "")
	match = strings.ReplaceAll(match, "authkey=", "")
	ExampleGachaReq := models.GachaRequest{
		AuthkeyVer: "1",
		SignType:   "2",
		AuthAppId:  "webview_gacha",
		InitType:   "301",
		Lang:       "en",
		AuthKey:    match,
		Page:       "1",
		Size:       "20",
		EndId:      "0",
		GachaType:  "301",
	}
	q := url.Values{}
	q.Add("authkey_ver", ExampleGachaReq.AuthkeyVer)
	q.Add("sign_type", ExampleGachaReq.SignType)
	q.Add("auth_appid", ExampleGachaReq.AuthAppId)
	q.Add("init_type", ExampleGachaReq.InitType)
	q.Add("lang", ExampleGachaReq.Lang)
	q.Add("authkey", ExampleGachaReq.AuthKey)
	q.Add("page", ExampleGachaReq.Page)
	q.Add("size", ExampleGachaReq.Size)
	q.Add("end_id", ExampleGachaReq.EndId)
	q.Add("gacha_type", ExampleGachaReq.GachaType)

	f := excelize.NewFile()
	f.SetSheetName("Sheet1", "Event Banner 1")
	f.SetCellValue("Event Banner 1", "A1", "TimeStamp")
	f.SetCellValue("Event Banner 1", "B1", "Name")
	f.SetCellValue("Event Banner 1", "C1", "Type")
	f.SetCellValue("Event Banner 1", "D1", "Rarity")
	f.NewSheet("Weapon Banner")
	f.SetCellValue("Weapon Banner", "A1", "TimeStamp")
	f.SetCellValue("Weapon Banner", "B1", "Name")
	f.SetCellValue("Weapon Banner", "C1", "Type")
	f.SetCellValue("Weapon Banner", "D1", "Rarity")
	f.NewSheet("Standard Banner")
	f.SetCellValue("Standard Banner", "A1", "TimeStamp")
	f.SetCellValue("Standard Banner", "B1", "Name")
	f.SetCellValue("Standard Banner", "C1", "Type")
	f.SetCellValue("Standard Banner", "D1", "Rarity")
	f.SetColWidth("Event Banner 1", "A", "G", 20)
	f.SetColWidth("Weapon Banner", "A", "G", 20)
	f.SetColWidth("Standard Banner", "A", "G", 20)
	styleB5, errB5 := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#FFB13F"}, Pattern: 1},
	})
	check(errB5)
	styleB4, errB4 := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#D28FD6"}, Pattern: 1},
	})
	check(errB4)
	styleB3, errB3 := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#4E7CFF"}, Pattern: 1},
	})
	check(errB3)
	styleHeader, errHeader := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#B8E8FC"}, Pattern: 1},
	})
	check(errHeader)

	loopcounter := 2
	totalwish := 0
	currentBanner := "Event Banner 1"
	DummyData := models.GachaDetail{
		RateOn:               false,
		Last5Stars:           "None",
		Last4Stars:           "None",
		CountAfterLast5Stars: 0,
		CountAfterLast4Stars: 0,
		Last5StarsFlag:       false,
		Last4StarsFlag:       false,
	}

	for {
		API_URL_EXEC := API_URL + "?authkey=" + match + "&" + q.Encode()
		response, err := http.Get(API_URL_EXEC)
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()
		b, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
		}
		var GachaResponse models.GachaResponse
		errnew := json.Unmarshal(b, &GachaResponse)
		if errnew != nil {
			panic(errnew)
		}

		if GachaResponse.Data.List == nil {
			fmt.Println("NILL")
			break
		}
		for i, value := range GachaResponse.Data.List {
			rarity, _ := strconv.ParseInt(value.RankType, 10, 64)
			// fmt.Println(value.Time + "\t" + value.Name + "\t" + value.ItemType + "\t" + value.RankType + "\t" + currentBanner)
			f.SetCellValue(currentBanner, fmt.Sprintf("A%v", strconv.Itoa(i+loopcounter)), value.Time)
			f.SetCellValue(currentBanner, fmt.Sprintf("B%v", strconv.Itoa(i+loopcounter)), value.Name)
			f.SetCellValue(currentBanner, fmt.Sprintf("C%v", strconv.Itoa(i+loopcounter)), value.ItemType)
			if rarity == 5 { //☆☆☆☆☆
				f.SetCellValue(currentBanner, fmt.Sprintf("D%v", strconv.Itoa(i+loopcounter)), "☆☆☆☆☆")
				errB5 = f.SetCellStyle(currentBanner, fmt.Sprintf("A%v", strconv.Itoa(i+loopcounter)), fmt.Sprintf("D%v", strconv.Itoa(i+loopcounter)), styleB5)
				if DummyData.Last5StarsFlag == false {
					if value.Name == "Keqing" || value.Name == "Diluc" || value.Name == "Mona" || value.Name == "Qiqi" || value.Name == "Tighnari" {
						DummyData.RateOn = true
					}
					DummyData.Last5Stars = value.Name
					DummyData.Last5StarsFlag = true
					DummyData.CountAfterLast5Stars = totalwish
				}
			} else if rarity == 4 { //☆☆☆☆
				f.SetCellValue(currentBanner, fmt.Sprintf("D%v", strconv.Itoa(i+loopcounter)), "☆☆☆☆")
				errB4 = f.SetCellStyle(currentBanner, fmt.Sprintf("A%v", strconv.Itoa(i+loopcounter)), fmt.Sprintf("D%v", strconv.Itoa(i+loopcounter)), styleB4)
				if DummyData.Last4StarsFlag == false {
					DummyData.Last4Stars = value.Name
					DummyData.Last4StarsFlag = true
					DummyData.CountAfterLast4Stars = totalwish
				}
			} else { //☆☆☆
				f.SetCellValue(currentBanner, fmt.Sprintf("D%v", strconv.Itoa(i+loopcounter)), "☆☆☆")
				errB3 = f.SetCellStyle(currentBanner, fmt.Sprintf("A%v", strconv.Itoa(i+loopcounter)), fmt.Sprintf("D%v", strconv.Itoa(i+loopcounter)), styleB3)
			}
			totalwish++
			if i == 19 {
				q.Set("end_id", value.Id)
				loopcounter += 20
			}
		}
		if len(GachaResponse.Data.List) < 20 {
			if currentBanner == "Standard Banner" || currentBanner == "Weapon Banner" {
				DummyData.RateOn = false
			}
			if currentBanner == "Event Banner 1" {
				q.Set("init_type", "302")
				q.Set("gacha_type", "302")
				q.Set("end_id", "0")
				errHeader = f.SetCellStyle(currentBanner, "A1", "D1", styleHeader)
				f.SetCellValue(currentBanner, "F1", "Total Wish:")
				f.SetCellValue(currentBanner, "G1", totalwish)
				f.SetCellValue(currentBanner, "F2", "Rate ON")
				f.SetCellValue(currentBanner, "G2", DummyData.RateOn)
				errHeader = f.SetCellStyle(currentBanner, "F1", "G2", styleHeader)
				f.SetCellValue(currentBanner, "F3", "Last 5 Stars")
				f.SetCellValue(currentBanner, "G3", DummyData.Last5Stars)
				f.SetCellValue(currentBanner, "F4", "Wish Until Next 5 Stars Estimate")
				f.SetCellValue(currentBanner, "G4", 90-DummyData.CountAfterLast5Stars)
				errB5 = f.SetCellStyle(currentBanner, "F3", "G4", styleB5)
				f.SetCellValue(currentBanner, "F5", "Last 4 Stars")
				f.SetCellValue(currentBanner, "G5", DummyData.Last4Stars)
				f.SetCellValue(currentBanner, "F6", "Wish Until Next 4 Stars Estimate")
				f.SetCellValue(currentBanner, "G6", 10-DummyData.CountAfterLast4Stars)
				errB4 = f.SetCellStyle(currentBanner, "F5", "G6", styleB4)
				totalwish = 0
				loopcounter = 2
				currentBanner = "Weapon Banner"
				DummyData.Last4StarsFlag = false
				DummyData.Last5StarsFlag = false
			} else if currentBanner == "Weapon Banner" {
				q.Set("init_type", "200")
				q.Set("gacha_type", "200")
				q.Set("end_id", "0")
				loopcounter = 2
				errHeader = f.SetCellStyle(currentBanner, "A1", "D1", styleHeader)
				f.SetCellValue(currentBanner, "F1", "Total Wish:")
				f.SetCellValue(currentBanner, "G1", totalwish)
				f.SetCellValue(currentBanner, "F2", "Rate ON")
				f.SetCellValue(currentBanner, "G2", DummyData.RateOn)
				errHeader = f.SetCellStyle(currentBanner, "F1", "G2", styleHeader)
				f.SetCellValue(currentBanner, "F3", "Last 5 Stars")
				f.SetCellValue(currentBanner, "G3", DummyData.Last5Stars)
				f.SetCellValue(currentBanner, "F4", "Wish Until Next 5 Stars Estimate")
				f.SetCellValue(currentBanner, "G4", 90-DummyData.CountAfterLast5Stars)
				errB5 = f.SetCellStyle(currentBanner, "F3", "G4", styleB5)
				f.SetCellValue(currentBanner, "F5", "Last 4 Stars")
				f.SetCellValue(currentBanner, "G5", DummyData.Last4Stars)
				f.SetCellValue(currentBanner, "F6", "Wish Until Next 4 Stars Estimate")
				f.SetCellValue(currentBanner, "G6", 10-DummyData.CountAfterLast4Stars)
				errB4 = f.SetCellStyle(currentBanner, "F5", "G6", styleB4)
				totalwish = 0
				currentBanner = "Standard Banner"
				DummyData.Last4StarsFlag = false
				DummyData.Last5StarsFlag = false
			} else if currentBanner == "Standard Banner" {
				errHeader = f.SetCellStyle(currentBanner, "A1", "D1", styleHeader)
				f.SetCellValue(currentBanner, "F1", "Total Wish:")
				f.SetCellValue(currentBanner, "G1", totalwish)
				f.SetCellValue(currentBanner, "F2", "Rate ON")
				f.SetCellValue(currentBanner, "G2", DummyData.RateOn)
				errHeader = f.SetCellStyle(currentBanner, "F1", "G2", styleHeader)
				f.SetCellValue(currentBanner, "F3", "Last 5 Stars")
				f.SetCellValue(currentBanner, "G3", DummyData.Last5Stars)
				f.SetCellValue(currentBanner, "F4", "Wish Until Next 5 Stars Estimate")
				f.SetCellValue(currentBanner, "G4", 90-DummyData.CountAfterLast5Stars)
				errB5 = f.SetCellStyle(currentBanner, "F3", "G4", styleB5)
				f.SetCellValue(currentBanner, "F5", "Last 4 Stars")
				f.SetCellValue(currentBanner, "G5", DummyData.Last4Stars)
				f.SetCellValue(currentBanner, "F6", "Wish Until Next 4 Stars Estimate")
				f.SetCellValue(currentBanner, "G6", 10-DummyData.CountAfterLast4Stars)
				errB4 = f.SetCellStyle(currentBanner, "F5", "G6", styleB4)
				break
			}
		}
	}
	if err := f.SaveAs("Genshin-Wishing.xlsx"); err != nil {
		fmt.Println(err)
	}
}
