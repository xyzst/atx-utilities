package main

import (
	"atx/posterior/utilities/internal/find-city-council-district/model"
	"atx/posterior/utilities/pkg/processor"
	"atx/posterior/utilities/pkg/util"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {

	var c = processor.CsvFile{}
	var proc = c.New("output.csv").(*processor.CsvFile)
	err := proc.Write([]string{"address", "city", "state", "zip_code", "district", "district_url", "confidence_score", "flagged?"})
	if err != nil {
		log.Fatal(err)
	}

	var client = &http.Client{}

	for {
		line, err := proc.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		var address = strings.ToUpper(line[0])
		var city = strings.ToUpper(line[1])
		var state = strings.ToUpper(line[2])
		var zip = strings.ToUpper(line[3])

		request, err := http.NewRequest("GET",
			"https://geo.austintexas.gov/arcgis/rest/services/Geocode/COA_Address_Locator/GeocodeServer/findAddressCandidates",
			nil,
		)

		if err != nil {
			log.Fatalf("❌ err: %v", err)
		}

		query := request.URL.Query()

		query.Add("f", "json")
		query.Add("Street", address)
		query.Add("City", city)
		query.Add("State", state)
		query.Add("ZIP", zip)
		query.Add("matchOutOfRange", "false")

		request.URL.RawQuery = query.Encode()

		request.Header.Set("Accept", "application/json;q=0.9,*/*;q=0.8")
		request.Header.Set("Accept-Language", "en-US,en;q=0.9")
		request.Header.Set("Accept-Encoding", "gzip, deflate, br")

		response, err := client.Do(request)

		if err != nil {
			log.Fatalf("err: %v", err)
		}

		decompressed, err := util.DecompressResponse(*response)
		if err != nil {
			log.Fatal(err)
		}

		body, err := io.ReadAll(decompressed)
		if err != nil {
			log.Fatalf("❌ err: unable transform serialized data: %v", err)
		}

		var matches model.Address

		err = json.Unmarshal(body, &matches)
		if err != nil {
			log.Fatalf("unable to unmarshal json: %v", err)
		}

		match, err := matches.GetLikelyCandidate()
		if err != nil {
			log.Printf("⚠️ warn: no match found for %s %s %s %s, will 🚩 address as invalid (confidence_score: 0) for manual review", address, city, state, zip)
			err = proc.Write([]string{address, city, state, zip, "", "", "0", "🚩"})
			if err != nil {
				log.Fatal(err)
			}
			continue
		}

		log.Printf("ℹ️ info: found %d address candidate(s), this address has the highest match potential: %v", len(matches.Candidates), match)

		// todo: refactor, abstract out

		request, err = http.NewRequest("POST",
			"https://geo.austintexas.gov/arcgis/rest/services/Shared/Property/MapServer/3/query",
			nil,
		)

		if err != nil {
			log.Fatalf("err: %v", err)
		}

		query = request.URL.Query()

		// values: pjson (pretty), json, html
		query.Add("f", "json")
		query.Add("outFields", "*")
		query.Add("geometry", strconv.FormatFloat(match.Location.X, 'f', 9, 64)+","+strconv.FormatFloat(match.Location.Y, 'f', 9, 64))
		query.Add("returnDistinctValues", "false")
		query.Add("geometryType", "esriGeometryPoint")

		request.URL.RawQuery = query.Encode()

		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Accept-Language", "en-US,en;q=0.9")
		request.Header.Set("Accept-Encoding", "gzip, deflate, br")

		response, err = client.Do(request)
		//defer response.Body.Close()

		if err != nil {
			log.Fatalf("❌ err: %v", err)
		}

		decompressed, err = util.DecompressResponse(*response)
		if err != nil {
			log.Fatal(err)
		}
		body, err = io.ReadAll(decompressed)
		if err != nil {
			log.Fatalf("❌ err: unable transform serialized data: %v", err)
		}

		var district model.District

		err = json.Unmarshal(body, &district)
		if err != nil {
			log.Fatalf("❌ err: unable to unmarshal json: %v", err)
		}

		district.Candidate = match

		if len(district.Features) == 0 {
			err = proc.Write([]string{address, city, state, zip, "-1", "N/A - address does not reside within a council district, likely within an unincorporated area or outside of Travis County", strconv.FormatFloat(district.Candidate.Score, 'f', 3, 64), "🚩"})
		} else {
			err = proc.Write([]string{address, city, state, zip, strconv.FormatInt(district.Features[0].Attributes.CouncilDistrict, 10), district.Features[0].Attributes.CouncilDistrictPath, strconv.FormatFloat(district.Candidate.Score, 'f', 3, 64), ""})
		}

		if err != nil {
			log.Fatal(err)
		}
	}
}
