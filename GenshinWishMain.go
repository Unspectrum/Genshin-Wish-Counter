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

type GachaRequest struct {
	AuthkeyVer string `json:"authkey_ver"`
	SignType   string `json:"sign_type"`
	AuthAppId  string `json:"auth_appid"`
	InitType   string `json:"init_type"`
	Lang       string `json:"lang"`
	AuthKey    string `json:"authkey"`
	Page       string `json:"page"`
	Size       string `json:"size"`
	EndId      string `json:"end_id"`
	GachaType  string `json:"gacha_type"`
}

type ExampleGachaList struct {
	UID       string `json:"uid"`
	GachaType string `json:"gacha_type"`
	ItemId    string `json:"item_id"`
	Count     string `json:"count"`
	Time      string `json:"time"`
	Name      string `json:"name"`
	Lang      string `json:"lang"`
	ItemType  string `json:"item_type"`
	RankType  string `json:"rank_type"`
	Id        string `json:"id"`
}

type ExampleGachaData struct {
	Page   string             `json:"page"`
	Size   string             `json:"size"`
	Total  string             `json:"total"`
	List   []ExampleGachaList `json:"list"`
	Region string             `json:"region"`
}

type ExampleGachaResponse struct {
	RetCode int              `json:"retcode"`
	Message string           `json:"message"`
	Data    ExampleGachaData `json:"data"`
}

// Internal function to open file
func openFile(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// Internal function to read from file with certain delimiter
// This function return empty string and error object if any error happen
func readFileToStringArray(file *os.File, delim string) ([]string, error) {
	dataBytes, err := io.ReadAll(file)
	if err != nil {
		return []string{}, err
	}
	return strings.Split(string(dataBytes), delim), nil
}

// Function to open file from filepath and split by delim.
// This function return empty string and error object if any error happen
func OpenFileToStringArray(filepath string, delim string) ([]string, error) {
	file, err := openFile(filepath)
	if err != nil {
		return []string{}, err
	}
	lines, err := readFileToStringArray(file, delim)
	return lines, nil
}

func GetInstallLocation(str string) string {
	str = strings.ReplaceAll(str, warmUpStr, "")
	str = strings.Split(str, "\\")[0]
	return strings.ReplaceAll(str, streamAssetsStr, "")
}

func OpenReadFileToString(str string) (string, error) {
	b, err := os.ReadFile(str)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func main() {
	// apiHost := "hk4e-api-os.hoyoverse.com"
	dir, err := os.UserHomeDir()
	check(err)
	log := dir + logLocation
	lines, err := OpenFileToStringArray(log, "\n")
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
	DataOpen, err := OpenReadFileToString(installLocation)
	if err != nil {
		panic(err)
	}

	dataSplit := strings.Split(string(DataOpen), "1/0")
	var exp = "(authkey=.+?game_biz=)"
	re := regexp.MustCompile(exp)
	match := re.FindStringSubmatch(dataSplit[len(dataSplit)-1])[0]
	match = strings.ReplaceAll(match, "&game_biz=", "")
	match = strings.ReplaceAll(match, "authkey=", "")
	ExampleGachaReq := GachaRequest{
		AuthkeyVer: "1",
		SignType:   "2",
		AuthAppId:  "webview_gacha",
		InitType:   "301",
		Lang:       "en",
		AuthKey:    match,
		Page:       "1",
		Size:       "5",
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
	f.SetCellValue("Event Banner 1", "B1", "Quantity")
	f.SetCellValue("Event Banner 1", "C1", "Name")
	f.SetCellValue("Event Banner 1", "D1", "Type")
	f.SetCellValue("Event Banner 1", "E1", "Rarity")
	f.NewSheet("Event Banner 2")
	f.SetCellValue("Event Banner 2", "A1", "TimeStamp")
	f.SetCellValue("Event Banner 2", "B1", "Quantity")
	f.SetCellValue("Event Banner 2", "C1", "Name")
	f.SetCellValue("Event Banner 2", "D1", "Type")
	f.SetCellValue("Event Banner 2", "E1", "Rarity")
	f.NewSheet("Weapon Banner")
	f.SetCellValue("Weapon Banner", "A1", "TimeStamp")
	f.SetCellValue("Weapon Banner", "B1", "Quantity")
	f.SetCellValue("Weapon Banner", "C1", "Name")
	f.SetCellValue("Weapon Banner", "D1", "Type")
	f.SetCellValue("Weapon Banner", "E1", "Rarity")
	f.NewSheet("Standard Banner")
	f.SetCellValue("Standard Banner", "A1", "TimeStamp")
	f.SetCellValue("Standard Banner", "B1", "Quantity")
	f.SetCellValue("Standard Banner", "C1", "Name")
	f.SetCellValue("Standard Banner", "D1", "Type")
	f.SetCellValue("Standard Banner", "E1", "Rarity")
	f.SetColWidth("Event Banner 1", "A", "E", 20)
	f.SetColWidth("Event Banner 2", "A", "E", 20)
	f.SetColWidth("Weapon Banner", "A", "E", 20)
	f.SetColWidth("Standard Banner", "A", "E", 20)
	loopcounter := 2
	totalwish := 0
	currentBanner := "Event Banner 1"
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
		var GachaResponse ExampleGachaResponse
		errnew := json.Unmarshal(b, &GachaResponse)
		if errnew != nil {
			panic(errnew)
		}

		if GachaResponse.Data.List == nil {
			fmt.Println("NILL")
			break
		}
		for i, value := range GachaResponse.Data.List {
			count, _ := strconv.ParseInt(value.Count, 10, 64)
			rarity, _ := strconv.ParseInt(value.RankType, 10, 64)
			fmt.Println(value.Time + "\t" + value.Count + "\t" + value.Name + "\t" + value.ItemType + "\t" + value.RankType + "\t" + currentBanner)
			f.SetCellValue(currentBanner, fmt.Sprintf("A%v", strconv.Itoa(i+loopcounter)), value.Time)
			f.SetCellValue(currentBanner, fmt.Sprintf("B%v", strconv.Itoa(i+loopcounter)), count)
			f.SetCellValue(currentBanner, fmt.Sprintf("C%v", strconv.Itoa(i+loopcounter)), value.Name)
			f.SetCellValue(currentBanner, fmt.Sprintf("D%v", strconv.Itoa(i+loopcounter)), value.ItemType)
			f.SetCellValue(currentBanner, fmt.Sprintf("E%v", strconv.Itoa(i+loopcounter)), rarity)
			if i == 4 {

				q.Set("end_id", value.Id)
				loopcounter += 5
				totalwish += 5
			}
		}
		if len(GachaResponse.Data.List) < 4 {
			if currentBanner == "Event Banner 1" {
				q.Set("init_type", "302")
				q.Set("gacha_type", "302")
				q.Set("end_id", "0")
				loopcounter = 2
				f.SetCellValue(currentBanner, fmt.Sprintf("F%v", strconv.Itoa(1)), totalwish)
				totalwish = 0
				currentBanner = "Event Banner 2"
			} else if currentBanner == "Event Banner 2" {
				q.Set("init_type", "200")
				q.Set("gacha_type", "200")
				q.Set("end_id", "0")
				loopcounter = 2
				f.SetCellValue(currentBanner, fmt.Sprintf("F%v", strconv.Itoa(1)), totalwish)
				totalwish = 0
				currentBanner = "Standard Banner"
			} else if currentBanner == "Standard Banner" {
				f.SetCellValue(currentBanner, fmt.Sprintf("F%v", strconv.Itoa(1)), totalwish)
				break
			}
		}
	}
	if err := f.SaveAs("Genshin-Wishing.xlsx"); err != nil {
		fmt.Println(err)
	}
}
