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

func main() {
    // read input file
    file, err := os.Open("input.csv")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    // parse input file into slice of tracks
    var tracks []Track
    lines := csv.NewReader(file)
    for i := 0; ; i++ {
        line, err := lines.Read()
        if err == io.EOF {
            break
        } else if err != nil {
            panic(err)
        } else if (line[1] == "Title" && line[2] == "Artist") {
            continue
        }
        track := Track{
            TrackID:  i,
            Name:     line[1],
            Artist:   line[2],
            Album:    line[3],
            Genre:    line[6],
            Location: fmt.Sprintf("file://localhost%s", strings.Trim(line[11], "\"")),
        }
        tracks = append(tracks, track)
    }

    // write tracks to output file
    output := Collection{Tracks: tracks}
    xmlOutput, err := xml.MarshalIndent(output, "", "  ")
    if err != nil {
        panic(err)
    }

    // write xmlOutput to output.xml
    outputFile, err := os.Create("output.xml")
    if err != nil {
        panic(err)
    }
    defer outputFile.Close()

    outputFile.Write([]byte(xml.Header))
    outputFile.Write(xmlOutput)
}
