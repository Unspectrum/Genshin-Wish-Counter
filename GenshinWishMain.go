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
	cp "github.com/nmrshll/go-cp"
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
	Inventory   = "Inventory"
)

var (
	ApiUrl  = "https://hk4e-api-os.hoyoverse.com/event/gacha_info/api/getGachaLog"
	Banners = []string{"Event Banner", "Weapon Banner", "Standard Banner"}
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
	currDir, err := os.Getwd()
	if err != nil {
		fmt.Println(currDir)
		panic("Fail to Get Current Directory")
	}
	projDir := currDir + "\\cache"
	_err := cp.CopyFile(dataFileLocation, projDir)
	if _err != nil {
		panic(_err)
	}
	DataOpen, err := utils.OpenReadFileToString(dataFileLocation)
	if err != nil {
		fmt.Print("Could Open Cache")
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

	f := utils.NewExcel("Genshin-Wishing.xlsx")
	f.ChangeSheetName("Sheet1", Banners[0])
	f.GenerateSheets(Banners[1:])
	f.SetCellValues(Banners, 'A', 1, []interface{}{"TimeStamp", "Name", "Type", "Rarity"})
	f.SetColWidth("Event Banner", "A", "G", 20)
	f.SetColWidth("Weapon Banner", "A", "G", 20)
	f.SetColWidth("Standard Banner", "A", "G", 20)

	styleB5 := f.MakeStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#FFB13F"}, Pattern: 1},
	})
	styleB4 := f.MakeStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#D28FD6"}, Pattern: 1},
	})
	styleB3 := f.MakeStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#4E7CFF"}, Pattern: 1},
	})
	styleHeader := f.MakeStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#B8E8FC"}, Pattern: 1},
	})

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
			continue
		}
		for i, value := range GachaResponse.Data.List {
			rarity, _ := strconv.ParseInt(value.RankType, 10, 64)
			// fmt.Println(value.Time + "\t" + value.Name + "\t" + value.ItemType + "\t" + value.RankType + "\t" + currentBanner)
			// fmt.Sprintf("A%v", strconv.Itoa(i+loopcounter))
			f.SetCellValues([]string{currentBanner}, 'A', i+loopcounter, []interface{}{value.Time})
			f.SetCellValues([]string{currentBanner}, 'B', i+loopcounter, []interface{}{value.Name})
			f.SetCellValues([]string{currentBanner}, 'C', i+loopcounter, []interface{}{value.ItemType})
			if rarity == 5 { //☆☆☆☆☆
				f.SetCellValues([]string{currentBanner}, 'D', i+loopcounter, []interface{}{"☆☆☆☆☆"})
				f.SetCellStyle(currentBanner, fmt.Sprintf("A%v", strconv.Itoa(i+loopcounter)), fmt.Sprintf("D%v", strconv.Itoa(i+loopcounter)), styleB5)
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
				f.SetCellValues([]string{currentBanner}, 'D', i+loopcounter, []interface{}{"☆☆☆☆"})
				f.SetCellStyle(currentBanner, fmt.Sprintf("A%v", strconv.Itoa(i+loopcounter)), fmt.Sprintf("D%v", strconv.Itoa(i+loopcounter)), styleB4)
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
				f.SetCellValues([]string{currentBanner}, 'D', i+loopcounter, []interface{}{"☆☆☆"})
				f.SetCellStyle(currentBanner, fmt.Sprintf("A%v", strconv.Itoa(i+loopcounter)), fmt.Sprintf("D%v", strconv.Itoa(i+loopcounter)), styleB3)
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
				f.SetCellStyle(currentBanner, "A1", "D1", styleHeader)
				f.SetCellValues([]string{currentBanner}, 'F', 1, []interface{}{"Total Wish:", totalwish})
				f.SetCellValues([]string{currentBanner}, 'F', 2, []interface{}{"Rate ON", DummyData.RateOn})
				f.SetCellStyle(currentBanner, "F1", "G2", styleHeader)
				f.SetCellValues([]string{currentBanner}, 'F', 3, []interface{}{"Last 5 Stars", DummyData.Last5Stars})
				f.SetCellValues([]string{currentBanner}, 'F', 4, []interface{}{"Wish Until Next 5 Stars Estimate", 90 - DummyData.CountAfterLast5Stars})
				f.SetCellStyle(currentBanner, "F3", "G4", styleB5)
				f.SetCellValues([]string{currentBanner}, 'F', 5, []interface{}{"Last 4 Stars", DummyData.Last4Stars})
				f.SetCellValues([]string{currentBanner}, 'F', 6, []interface{}{"Wish Until Next 4 Stars Estimate", 10 - DummyData.CountAfterLast4Stars})
				f.SetCellStyle(currentBanner, "F5", "G6", styleB4)

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
				f.SetCellStyle(currentBanner, "A1", "D1", styleHeader)
				f.SetCellValues([]string{currentBanner}, 'F', 1, []interface{}{"Total Wish:", totalwish})
				f.SetCellValues([]string{currentBanner}, 'F', 2, []interface{}{"Rate ON", DummyData.RateOn})
				f.SetCellStyle(currentBanner, "F1", "G2", styleHeader)
				f.SetCellValues([]string{currentBanner}, 'F', 3, []interface{}{"Last 5 Stars", DummyData.Last5Stars})
				f.SetCellValues([]string{currentBanner}, 'F', 4, []interface{}{"Wish Until Next 5 Stars Estimate", 90 - DummyData.CountAfterLast5Stars})
				f.SetCellStyle(currentBanner, "F3", "G4", styleB5)
				f.SetCellValues([]string{currentBanner}, 'F', 5, []interface{}{"Last 4 Stars", DummyData.Last4Stars})
				f.SetCellValues([]string{currentBanner}, 'F', 6, []interface{}{"Wish Until Next 4 Stars Estimate", 10 - DummyData.CountAfterLast4Stars})
				f.SetCellStyle(currentBanner, "F5", "G6", styleB4)
				WeaponCounter = totalwish
				totalwish = 0
				currentBanner = "Standard Banner"
				DummyData.Last4StarsFlag = false
				DummyData.Last5StarsFlag = false
			} else if currentBanner == "Standard Banner" {
				f.SetCellStyle(currentBanner, "A1", "D1", styleHeader)
				f.SetCellValues([]string{currentBanner}, 'F', 1, []interface{}{"Total Wish:", totalwish})
				f.SetCellValues([]string{currentBanner}, 'F', 2, []interface{}{"Rate ON", DummyData.RateOn})
				f.SetCellStyle(currentBanner, "F1", "G2", styleHeader)
				f.SetCellValues([]string{currentBanner}, 'F', 3, []interface{}{"Last 5 Stars", DummyData.Last5Stars})
				f.SetCellValues([]string{currentBanner}, 'F', 4, []interface{}{"Wish Until Next 5 Stars Estimate", 90 - DummyData.CountAfterLast5Stars})
				f.SetCellStyle(currentBanner, "F3", "G4", styleB5)
				f.SetCellValues([]string{currentBanner}, 'F', 5, []interface{}{"Last 4 Stars", DummyData.Last4Stars})
				f.SetCellValues([]string{currentBanner}, 'F', 6, []interface{}{"Wish Until Next 4 Stars Estimate", 10 - DummyData.CountAfterLast4Stars})
				f.SetCellStyle(currentBanner, "F5", "G6", styleB4)
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

	f.GenerateSheets([]string{Inventory})
	f.SetCellValues([]string{Inventory}, 'A', 1, []interface{}{"Character Name", "Quantity", "Rarity"})

	f.SetCellValues([]string{Inventory}, 'E', 1, []interface{}{"Weapon Name", "Quantity", "Rarity"})
	f.SetCellStyle("Inventory", "A1", "C1", styleHeader)
	f.SetCellStyle("Inventory", "E1", "G1", styleHeader)
	f.SetColWidth("Inventory", "A", "G", 20)
	loopcounter = 2
	for index, items := range MyCharacters[0].ItemLists {
		f.SetCellValues([]string{Inventory}, 'A', index+loopcounter, []interface{}{items.Name, items.Quantity, "☆☆☆☆☆"})
		f.SetCellStyle("Inventory", fmt.Sprintf("A%v", strconv.Itoa(index+loopcounter)), fmt.Sprintf("C%v", strconv.Itoa(index+loopcounter)), styleB5)
	}
	loopcounter += len(MyCharacters[0].ItemLists)
	for index, items := range MyCharacters[1].ItemLists {
		f.SetCellValues([]string{Inventory}, 'A', index+loopcounter, []interface{}{items.Name, items.Quantity, "☆☆☆☆"})
		f.SetCellStyle("Inventory", fmt.Sprintf("A%v", strconv.Itoa(index+loopcounter)), fmt.Sprintf("C%v", strconv.Itoa(index+loopcounter)), styleB4)
	}

	loopcounter = 2
	for index, items := range MyWeapons[0].ItemLists {
		f.SetCellValues([]string{Inventory}, 'E', index+loopcounter, []interface{}{items.Name, items.Quantity, "☆☆☆☆☆"})
		f.SetCellStyle("Inventory", fmt.Sprintf("E%v", strconv.Itoa(index+loopcounter)), fmt.Sprintf("G%v", strconv.Itoa(index+loopcounter)), styleB5)
	}
	loopcounter += len(MyWeapons[0].ItemLists)
	for index, items := range MyWeapons[1].ItemLists {
		f.SetCellValues([]string{Inventory}, 'E', index+loopcounter, []interface{}{items.Name, items.Quantity, "☆☆☆☆"})
		f.SetCellStyle("Inventory", fmt.Sprintf("E%v", strconv.Itoa(index+loopcounter)), fmt.Sprintf("G%v", strconv.Itoa(index+loopcounter)), styleB4)
	}

	err = f.SaveFile()
	if err != nil {
		panic(err)
	}
}
