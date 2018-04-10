package main

import (
    "net/http"
    "io/ioutil"
    "time"
    "log"
    "strings"
    "errors"
    "os"
    "html"
    "labix.org/v2/mgo"
)

func get_between (str string, start string, finish string) (string, error) {
    begin := strings.Index(str, start)
    if begin == -1 {
        return "", errors.New("No such start string match") 
    }
    begin = begin + len(start)
    end := strings.Index(str[begin:], finish) + begin
    if begin == -1 {
        return "", errors.New("No such end string match")
    }

    toReturn := html.UnescapeString(string(str[begin:end]))

    return toReturn, nil
}



func main() {
    var client = &http.Client {
        Timeout: time.Second * 10,
    }
    resp, err := client.Get("https://electronic.vegas/vegas-edm-event-calendar/")
    defer resp.Body.Close()
    if err != nil {
        log.Fatal(err)
    }
    
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
    
    plaintextBody := string(body)
    var events []map[string]string
    
    for _, line := range strings.Split(plaintextBody, "\n") {
        l := log.New(os.Stderr, "", 0)
        lineHasDate := strings.Contains(line, "wideeventDate")
        lineHasTitle := strings.Contains(line, "wideeventTitle")
        if lineHasDate {
            s, err := get_between(line, "<meta itemprop='startDate' content='", "'")
            if err != nil {
                log.Fatal("line has wideEventDate but couldn't extract the date")
            }
            l.Print(s)
        }
        if lineHasTitle {
            s, err := get_between(line, "itemprop='name'>", "</span>")
            if err != nil {
                log.Fatal("line has wideeventTitle but couldn't extract the title")
            }
            l.Print(s)
        }
    }

    session, err := mgo.Dial("mongodb://127.0.0.1:27017")    
    defer session.Close()
    go_event_col := session.DB("go_event_db").C("go_event_collection")
    
    log.Print(go_event_col.Count())
}

