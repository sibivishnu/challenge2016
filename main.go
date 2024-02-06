package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

type City struct {
	CityCode, ProvinceCode, CountryCode string
	CityName, ProvinceName, CountryName string
}

type Permission struct {
	Include []string
	Exclude []string
}

type Distributor struct {
	Name        string
	Permissions Permission
	Parent      *Distributor
}

func LoadCities(filePath string) ([]City, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var cities []City
	for _, record := range records[1:] {
		cities = append(cities, City{
			CityCode:     record[0],
			ProvinceCode: record[1],
			CountryCode:  record[2],
			CityName:     record[3],
			ProvinceName: record[4],
			CountryName:  record[5],
		})
	}

	return cities, nil
}

func (d *Distributor) HasPermission(city City) bool {
	for _, exclude := range d.Permissions.Exclude {
		if strings.Contains(city.CityCode, exclude) || strings.Contains(city.ProvinceCode, exclude) || strings.Contains(city.CountryCode, exclude) {
			return false
		}
	}

	for _, include := range d.Permissions.Include {
		if strings.Contains(city.CityCode, include) || strings.Contains(city.ProvinceCode, include) || strings.Contains(city.CountryCode, include) {
			return true
		}
	}

	if d.Parent != nil {
		return d.Parent.HasPermission(city)
	}

	return false
}
func main() {
	cities, err := LoadCities("<YOUR-CSV-PATH>")
	if err != nil {
		fmt.Println("Error loading cities:", err)
		return
	}

	distributor1 := Distributor{
		Name: "DISTRIBUTOR1",
		Permissions: Permission{
			Include: []string{"IN", "US"},
			Exclude: []string{"KARNATAKA-INDIA", "CHENNAI-TAMILNADU-INDIA"},
		},
	}

	distributor2 := Distributor{
		Name: "DISTRIBUTOR2",
		Permissions: Permission{
			Include: []string{"IN"},
			Exclude: []string{"TAMILNADU-INDIA"},
		},
		Parent: &distributor1,
	}

	for _, city := range cities {

		if distributor2.HasPermission(city) {
			fmt.Printf("%s has permission to distribute in %s, %s, %s\n", distributor2.Name, city.CityName, city.ProvinceName, city.CountryName)
		} else {
			fmt.Printf("%s does not have permission to distribute in %s, %s, %s\n", distributor2.Name, city.CityName, city.ProvinceName, city.CountryName)
		}
	}
}
