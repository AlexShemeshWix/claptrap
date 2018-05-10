package claymore

import (
	"errors"
	"strings"
	"regexp"

	"encoding/json"
	"github.com/ghodss/yaml"
	"fmt"
)

type MinerEntry struct {
	MinerName   string
	Currency    string
	RunningTime string
	Hashrate    []string
	MiningPool  []string
}

func ParseString(RawData string) (retVal []MinerEntry, err error) {

	return retVal, err
}

func (this *MinerEntry) SetFieldForIndex(data string, index int) {
	if index == 0 {
		this.MinerName = data
	} else if index == 2 {
		this.RunningTime = data
	} else if index == 3 {
		this.Hashrate[0] = data
	} else if index == 4 {
		this.Hashrate[1] = data
	} else if index == 6 {
		minerpools := strings.Split(data, ";")
		this.MiningPool[0] = minerpools[0]
		if len(minerpools) >1 {
			this.MiningPool[1] = minerpools[1]
		}
	}
}

func parseMinerInfoFromRawData(data string) (retVal *MinerEntry, err error) {
	ok := strings.Contains(data, "<td>")
	if ok {
		retVal = &MinerEntry{Hashrate: make([]string, 2), MiningPool: make([]string, 2)}
		rp, err := regexp.Compile("(COLOR=(.*)>(.*)<)")
		if err != nil {
			return retVal, err
		}
		fieldCounter := 0
		fontEntries := strings.Split(data, "<FONT")
		for _, newValue := range (fontEntries) {
			values := rp.FindAllString(newValue, -1)
			for _, newValue1 := range (values) {
				tagDataStart := strings.Index(newValue1, ">") + 1
				tagDataEnd := strings.Index(newValue1, "<")

				retVal.SetFieldForIndex(string(newValue1[tagDataStart:tagDataEnd]), fieldCounter)
				fieldCounter++

			}
		}
		return retVal, nil
	}
	return retVal, errors.New("No miner info found")
}

func ObjectAsYAMLToString(obj interface{}) (retVal string) {
	var objectContent []byte
	var err error
	objectContent, err = json.Marshal(obj)
	objectasYaml, err := yaml.JSONToYAML(objectContent)
	if err != nil {
		print(err)
	}
	return "\n" + string(objectasYaml)
}

func SplitTable(tableText string) (retVal []MinerEntry){
	lines := strings.Split(tableText, "<tr>")
	for _, val := range (lines) {
		miner, err := parseMinerInfoFromRawData(val)
		if err == nil {
			fmt.Printf("Miner: %s, RunningTime: %s, HashRate1: %s for pool1:%s, HashRate2: %s, for pool2: %s" , miner.MinerName, miner.RunningTime, miner.Hashrate[0], miner.MiningPool[0],miner.Hashrate[1], miner.MiningPool[1])
			println("========================================")
			retVal = append(retVal,*miner)
		}
	}
	return retVal
}
