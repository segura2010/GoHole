package parser

import (
    "fmt"
    "os"
    "bufio"
    "strings"
    "io"
    "io/ioutil"
    "net/http"

    //"GoHole/config"
    "GoHole/dnscache"
)

func ParseBlacklistFile(path string) (error){
    var err error = nil
    if path[0:4] == "http"{
        // download file, parse and delete
        path, err = downloadFile(path)
        if err != nil{
            return err
        }
        defer os.Remove(path)
    }

    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        // read file line by line
        line := scanner.Text()
        var parsedLine []string = nil

        // if starts with # it is a comment
        if line != "" && line[0:1] != "#" {
            parsedLine = strings.Split(line, "\t")
            if len(parsedLine) < 2{
                parsedLine = strings.Split(line, " ")
                if len(parsedLine) < 2{
                    parsedLine = append(parsedLine, parsedLine[0]) // it is not a hosts file, it just include a domain per line, so les's create a hosts like array
                    parsedLine[0] = "127.0.0.1" // block domain redirected to local address
                }
            }
            
            if parsedLine[1] != "localhost"{
                // clean domain and save it on block list
                parsedLine[1] = strings.Replace(parsedLine[1], " ", "", -1)
                parsedLine[1] = strings.Replace(parsedLine[1], "\t", "", -1)

                fmt.Printf("\nDomain %s blocked with %s", parsedLine[1], parsedLine[0])

                dnscache.AddDomainIPv4(parsedLine[1], parsedLine[0], 0)
                dnscache.AddDomainIPv6(parsedLine[1], "::1", 0) // by default ad lists doesn't include ipv6 block..
            }
        }
    }

    return nil
}

func downloadFile(url string) (string, error) {
    // Create temporal file
    tmpfile, err := ioutil.TempFile(os.TempDir(), "gohole")

    // Get the data
    resp, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    // Write the body to file
    _, err = io.Copy(tmpfile, resp.Body)
    if err != nil  {
        return "", err
    }

    return tmpfile.Name(), nil
}

func ParseBlacklistsListFile(path string) (error){
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        // read file line by line
        line := scanner.Text()

        // if starts with # it is a comment
        if line != "" && line[0:1] != "#" {
            ParseBlacklistFile(line)
        }
    }

    return nil
}

