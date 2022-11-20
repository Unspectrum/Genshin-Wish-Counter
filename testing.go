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

	API_URL += "?authkey=" + match + "&" + q.Encode()
	response, err := http.Get(API_URL)
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
	f := excelize.NewFile()
	f.SetSheetName("Sheet1", "Event Banner")
	f.SetCellValue("Event Banner", "A1", "TimeStamp")
	f.SetCellValue("Event Banner", "B1", "Quantity")
	f.SetCellValue("Event Banner", "C1", "Name")
	f.SetCellValue("Event Banner", "D1", "Type")
	f.SetCellValue("Event Banner", "E1", "Rarity")
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
	f.SetColWidth("Event Banner", "A", "E", 20)
	for i, value := range GachaResponse.Data.List {
		count, _ := strconv.ParseInt(value.Count, 10, 64)
		rarity, _ := strconv.ParseInt(value.RankType, 10, 64)
		f.SetCellValue("Event Banner", fmt.Sprintf("A%v", strconv.Itoa(i+2)), value.Time)
		f.SetCellValue("Event Banner", fmt.Sprintf("B%v", strconv.Itoa(i+2)), count)
		f.SetCellValue("Event Banner", fmt.Sprintf("C%v", strconv.Itoa(i+2)), value.Name)
		f.SetCellValue("Event Banner", fmt.Sprintf("D%v", strconv.Itoa(i+2)), value.ItemType)
		f.SetCellValue("Event Banner", fmt.Sprintf("E%v", strconv.Itoa(i+2)), rarity)
	}
	// Create a new sheet.

	// Set value of a cell.
	// f.SetCellValue("Sheet2", "A1", "Hello world2.")
	// f.SetCellValue("Sheet1", "B1", 1000)
	// Set active sheet of the workbook.

	// Save spreadsheet by the given path.
	if err := f.SaveAs("Genshin-Wishing.xlsx"); err != nil {
		fmt.Println(err)
	}
}
