package main

import (
    "encoding/csv"
    "encoding/xml"
    "fmt"
    "io"
    "os"
    "strings"
)

type Track struct {
    TrackID  int    `xml:"TrackID,attr"`
    Name     string `xml:"Name,attr"`
    Artist   string `xml:"Artist,attr"`
    Album    string `xml:"Album,attr"`
    Genre    string `xml:"Genre,attr"`
    Location string `xml:"Location,attr"`
}

type Collection struct {
    XMLName xml.Name `xml:"DJ_PLAYLISTS"`
    Tracks  []Track  `xml:"COLLECTION>TRACK"`
}

type CustomCSVReader struct {
    reader *csv.Reader
}

func NewCustomCSVReader(r io.Reader) *CustomCSVReader {
    data, err := io.ReadAll(r)
    if err != nil {
        panic(err)
    }

    reader := csv.NewReader(strings.NewReader(string(data)))
    reader.LazyQuotes = true
    reader.TrimLeadingSpace = true

    return &CustomCSVReader{reader: reader}
}

func (c *CustomCSVReader) Read() ([]string, error) {
    record, err := c.reader.Read()
    if err != nil {
        return nil, err
    }
    for i := range record {
        record[i] = strings.Trim(record[i], "\"")
    }
    return record, nil
}

func main() {
    // read input file
    file, err := os.Open("input.csv")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    // parse input file into a slice of tracks
    var tracks []Track
    reader := NewCustomCSVReader(file)
    lineCount := 0
    for {
        line, err := reader.Read()
        if err == io.EOF {
            break
        } else if err != nil {
            panic(err)
        }
        lineCount++

        // Skip the header line
        if lineCount == 1 {
            continue
        }

        if len(line) != 12 {
            fmt.Printf("LENGTH=[%d]", len(line))
            fmt.Printf("Record on line %d: wrong number of fields\n", lineCount)
            continue
        }

        track := Track{
            TrackID:  lineCount - 1,
            Name:     line[1],
            Artist:   line[2],
            Album:    line[3],
            Genre:    line[6],
            Location: fmt.Sprintf("file://localhost%s", line[11]),
        }
        tracks = append(tracks, track)
    }

    // write tracks to output file
    output := Collection{Tracks: tracks}
    xmlOutput, err := xml.MarshalIndent(output, "", "  ")
    if err != nil {
        panic(err)
    }


    // Create the output.xml file
    outputFile, err := os.Create("output.xml")
    if err != nil {
        panic(err)
    }
    defer outputFile.Close()

    // Write XML header and content to the output file
    outputFile.Write([]byte(xml.Header))
    outputFile.Write(xmlOutput)

    fmt.Println("Output.xml file created successfully.")
}
