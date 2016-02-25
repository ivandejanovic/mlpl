package cfg

import (
	"bufio"
	"fmt"
	"mlpl/types"
	"os"
)

func GetDefaultReserved() []types.ReservedWord {
	reserved := make([]types.ReservedWord, 0, 8)

	reserved = append(reserved, types.ReservedWord{types.IF, "if"})
	reserved = append(reserved, types.ReservedWord{types.THEN, "then"})
	reserved = append(reserved, types.ReservedWord{types.ELSE, "else"})
	reserved = append(reserved, types.ReservedWord{types.END, "end"})
	reserved = append(reserved, types.ReservedWord{types.REPEAT, "repeat"})
	reserved = append(reserved, types.ReservedWord{types.UNTIL, "until"})
	reserved = append(reserved, types.ReservedWord{types.READ, "read"})
	reserved = append(reserved, types.ReservedWord{types.WRITE, "write"})

	return reserved
}

func GetConfigReservedWords(configFile string) []types.ReservedWord {
	reserved := make([]types.ReservedWord, 0, 8)
	var localization []string
	const length = 8

	config, err := os.Open(configFile)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(config)
	for scanner.Scan() {
		localization = append(localization, scanner.Text())
	}
	
	config.Close()

	if len(localization) != length {
		fmt.Println("Configuration file must contain localizations for eight key word.")
	}

	reserved = append(reserved, types.ReservedWord{types.IF, localization[0]})
	reserved = append(reserved, types.ReservedWord{types.THEN, localization[1]})
	reserved = append(reserved, types.ReservedWord{types.ELSE, localization[2]})
	reserved = append(reserved, types.ReservedWord{types.END, localization[3]})
	reserved = append(reserved, types.ReservedWord{types.REPEAT, localization[4]})
	reserved = append(reserved, types.ReservedWord{types.UNTIL, localization[5]})
	reserved = append(reserved, types.ReservedWord{types.READ, localization[6]})
	reserved = append(reserved, types.ReservedWord{types.WRITE, localization[7]})

	return reserved
}
