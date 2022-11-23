package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
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
	logLocation = "\\AppData\\LocalLow\\miHoYo\\Genshin Impact\\output_log.txt"
	exp         = "(authkey=.+?game_biz=)"
)

var (
	ApiUrl = "https://hk4e-api-os.hoyoverse.com/event/gacha_info/api/getGachaLog"
)

var clear map[string]func() //create a map for storing clear funcs

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func main() {
	// apiHost := "hk4e-api-os.hoyoverse.com"
	dir, err := os.UserHomeDir()
	check(err)
	log := dir + logLocation
	lines, err := utils.OpenFileToStringArray(log, "\n")
	check(err)

	dataFileLocation := utils.GetDataFileLocation(lines)
	DataOpen, err := utils.OpenReadFileToString(dataFileLocation)
	if err != nil {
		panic(err)
	}

	dataSplit := strings.Split(DataOpen, "1/0")
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

	params, err := utils.ParseStructToJsonMap(ExampleGachaReq)
	if err != nil {
		panic(err)
	}
	q := utils.GenerateGetParameter(params)

	f := utils.NewExcel().ExcelFile()
	f.SetSheetName("Sheet1", "Event Banner")
	f.SetCellValue("Event Banner", "A1", "TimeStamp")
	f.SetCellValue("Event Banner", "B1", "Name")
	f.SetCellValue("Event Banner", "C1", "Type")
	f.SetCellValue("Event Banner", "D1", "Rarity")
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
	f.SetColWidth("Event Banner", "A", "G", 20)
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

	currentBanner := "Event Banner"
	DummyData := models.GachaDetail{
		RateOn:               false,
		Last5Stars:           "None",
		Last4Stars:           "None",
		CountAfterLast5Stars: 0,
		CountAfterLast4Stars: 0,
		Last5StarsFlag:       false,
		Last4StarsFlag:       false,
	}

	var MyCharacters []models.Inventory
	MyCharacters = append(MyCharacters,
		models.Inventory{
			ItemType: 5,
		},
		models.Inventory{
			ItemType: 4,
		},
	)

	var MyWeapons []models.Inventory
	MyWeapons = append(MyWeapons,
		models.Inventory{
			ItemType: 5,
		},
		models.Inventory{
			ItemType: 4,
		},
	)

	loopcounter := 2
	totalwish := 0
	EventCounter := 0
	WeaponCounter := 0
	StandardCounter := 0
	MainCounter := 0
	for {
		API_URL_EXEC := ApiUrl + "?authkey=" + match + "&" + q.Encode()
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
				if value.ItemType == "Character" {
					flags := 0
					temp := models.InventoryDetail{
						Name:     value.Name,
						Quantity: 1,
					}
					for IData, Data := range MyCharacters[0].ItemLists {
						if value.Name == Data.Name {
							MyCharacters[0].ItemLists[IData].Quantity += 1
							flags = 1
						}
					}
					if len(MyCharacters[0].ItemLists) == 0 || flags == 0 {
						MyCharacters[0].ItemLists = append(MyCharacters[0].ItemLists, temp)
					}
				} else if value.ItemType == "Weapon" {
					flags := 0
					temp := models.InventoryDetail{
						Name:     value.Name,
						Quantity: 1,
					}
					for IData, Data := range MyWeapons[0].ItemLists {
						if value.Name == Data.Name {
							MyWeapons[0].ItemLists[IData].Quantity += 1
							flags = 1
						}
					}
					if len(MyWeapons[0].ItemLists) == 0 || flags == 0 {
						MyWeapons[0].ItemLists = append(MyWeapons[0].ItemLists, temp)
					}
				}
			} else if rarity == 4 { //☆☆☆☆
				f.SetCellValue(currentBanner, fmt.Sprintf("D%v", strconv.Itoa(i+loopcounter)), "☆☆☆☆")
				errB4 = f.SetCellStyle(currentBanner, fmt.Sprintf("A%v", strconv.Itoa(i+loopcounter)), fmt.Sprintf("D%v", strconv.Itoa(i+loopcounter)), styleB4)
				if DummyData.Last4StarsFlag == false {
					DummyData.Last4Stars = value.Name
					DummyData.Last4StarsFlag = true
					DummyData.CountAfterLast4Stars = totalwish
				}
				if value.ItemType == "Character" {
					flags := 0
					temp := models.InventoryDetail{
						Name:     value.Name,
						Quantity: 1,
					}
					for IData, Data := range MyCharacters[1].ItemLists {
						if value.Name == Data.Name {
							MyCharacters[1].ItemLists[IData].Quantity += 1
							flags = 1
						}
					}
					if len(MyCharacters[1].ItemLists) == 0 || flags == 0 {
						MyCharacters[1].ItemLists = append(MyCharacters[1].ItemLists, temp)
					}
				} else if value.ItemType == "Weapon" {
					flags := 0
					temp := models.InventoryDetail{
						Name:     value.Name,
						Quantity: 1,
					}
					for IData, Data := range MyWeapons[1].ItemLists {
						if value.Name == Data.Name {
							MyWeapons[1].ItemLists[IData].Quantity += 1
							flags = 1
						}
					}
					if len(MyWeapons[1].ItemLists) == 0 || flags == 0 {
						MyWeapons[1].ItemLists = append(MyWeapons[1].ItemLists, temp)
					}
				}
			} else { //☆☆☆
				f.SetCellValue(currentBanner, fmt.Sprintf("D%v", strconv.Itoa(i+loopcounter)), "☆☆☆")
				errB3 = f.SetCellStyle(currentBanner, fmt.Sprintf("A%v", strconv.Itoa(i+loopcounter)), fmt.Sprintf("D%v", strconv.Itoa(i+loopcounter)), styleB3)
			}
			totalwish++
			MainCounter++
			if i == 19 {
				q.Set("end_id", value.Id)
				loopcounter += 20
			}
		}
		if len(GachaResponse.Data.List) < 20 {
			if currentBanner == "Standard Banner" || currentBanner == "Weapon Banner" {
				DummyData.RateOn = false
			}
			if currentBanner == "Event Banner" {
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
				EventCounter = totalwish
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
				WeaponCounter = totalwish
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
				StandardCounter = totalwish
			}
		}
		CallClear()
		fmt.Println(
			" _____ _   _  _____ ___________   _____ _________________" + "\n" +
				"/  ___| | | ||  ___|  ___| ___ \\ /  __ \\  _  | ___ \\ ___ \\" + "\n" +
				"\\ `--.| |_| || |__ | |__ | |_/ / | /  \\/ | | | |_/ / |_/ /" + "\n" +
				" `--. \\  _  ||  __||  __||  __/  | |   | | | |    /|  __/" + "\n" +
				"/\\__/ / | | || |___| |___| |     | \\__/\\ \\_/ / |\\ \\| |" + "\n" +
				"\\____/\\_| |_/\\____/\\____/\\_|      \\____/\\___/\\_| \\_\\_|",
		)
		fmt.Println("Please Kindly Wait While We Counting Your Genshin Wish OwO")
		fmt.Println("Wish Count: ", MainCounter)
		fmt.Println("Event Banner : ", EventCounter)
		fmt.Println("Weapon Banner : ", WeaponCounter)
		fmt.Println("Standard Banner : ", StandardCounter)
		if StandardCounter > 0 {
			fmt.Println("YAY We've Done Counting Your Genshin Wish OwO")
			fmt.Println("Pweese Pwess Enter To Continue UwU")
			fmt.Println("Bye Bye OwO")
			fmt.Println()
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			break
		}
	}

	f.NewSheet("Inventory")
	f.SetCellValue("Inventory", "A1", "Character Name")
	f.SetCellValue("Inventory", "B1", "Quantity")
	f.SetCellValue("Inventory", "C1", "Rarity")

	f.SetCellValue("Inventory", "E1", "Weapon Name")
	f.SetCellValue("Inventory", "F1", "Quantity")
	f.SetCellValue("Inventory", "G1", "Rarity")
	errHeader = f.SetCellStyle("Inventory", "A1", "C1", styleHeader)
	errHeader = f.SetCellStyle("Inventory", "E1", "G1", styleHeader)
	f.SetColWidth("Inventory", "A", "G", 20)
	loopcounter = 2
	for index, items := range MyCharacters[0].ItemLists {
		f.SetCellValue("Inventory", fmt.Sprintf("A%v", strconv.Itoa(index+loopcounter)), items.Name)
		f.SetCellValue("Inventory", fmt.Sprintf("B%v", strconv.Itoa(index+loopcounter)), items.Quantity)
		f.SetCellValue("Inventory", fmt.Sprintf("C%v", strconv.Itoa(index+loopcounter)), "☆☆☆☆☆")
		errB5 = f.SetCellStyle("Inventory", fmt.Sprintf("A%v", strconv.Itoa(index+loopcounter)), fmt.Sprintf("C%v", strconv.Itoa(index+loopcounter)), styleB5)
	}
	loopcounter += len(MyCharacters[0].ItemLists)
	for index, items := range MyCharacters[1].ItemLists {
		f.SetCellValue("Inventory", fmt.Sprintf("A%v", strconv.Itoa(index+loopcounter)), items.Name)
		f.SetCellValue("Inventory", fmt.Sprintf("B%v", strconv.Itoa(index+loopcounter)), items.Quantity)
		f.SetCellValue("Inventory", fmt.Sprintf("C%v", strconv.Itoa(index+loopcounter)), "☆☆☆☆")
		errB4 = f.SetCellStyle("Inventory", fmt.Sprintf("A%v", strconv.Itoa(index+loopcounter)), fmt.Sprintf("C%v", strconv.Itoa(index+loopcounter)), styleB4)
	}

	loopcounter = 2
	for index, items := range MyWeapons[0].ItemLists {
		f.SetCellValue("Inventory", fmt.Sprintf("E%v", strconv.Itoa(index+loopcounter)), items.Name)
		f.SetCellValue("Inventory", fmt.Sprintf("F%v", strconv.Itoa(index+loopcounter)), items.Quantity)
		f.SetCellValue("Inventory", fmt.Sprintf("G%v", strconv.Itoa(index+loopcounter)), "☆☆☆☆☆")
		errB5 = f.SetCellStyle("Inventory", fmt.Sprintf("E%v", strconv.Itoa(index+loopcounter)), fmt.Sprintf("G%v", strconv.Itoa(index+loopcounter)), styleB5)
	}
	loopcounter += len(MyWeapons[0].ItemLists)
	for index, items := range MyWeapons[1].ItemLists {
		f.SetCellValue("Inventory", fmt.Sprintf("E%v", strconv.Itoa(index+loopcounter)), items.Name)
		f.SetCellValue("Inventory", fmt.Sprintf("F%v", strconv.Itoa(index+loopcounter)), items.Quantity)
		f.SetCellValue("Inventory", fmt.Sprintf("G%v", strconv.Itoa(index+loopcounter)), "☆☆☆☆")
		errB4 = f.SetCellStyle("Inventory", fmt.Sprintf("E%v", strconv.Itoa(index+loopcounter)), fmt.Sprintf("G%v", strconv.Itoa(index+loopcounter)), styleB4)
	}

	if err := f.SaveAs("Genshin-Wishing.xlsx"); err != nil {
		fmt.Println(err)
	}
}
