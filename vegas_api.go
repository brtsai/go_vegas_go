package main

import (
    "fmt"
    "net/http"
    "labix.org/v2/mgo"
    "labix.org/v2/mgo/bson"
//    "io/ioutil"
    "log"
)

func event_lookup(event string)  {
    session, err := mgo.Dial("mongodb://127.0.0.1:27017")
    defer session.Close()
    if err != nil {
        return
    }
    go_event_col := session.DB("go_event_db").C("go_event_collection")
    result := go_event_col.Find(bson.M{})
    log.Print("querying: " + event)
    log.Print(bson.M{"_id" : event})
    log.Print(result)    
}

func handler(w http.ResponseWriter, r *http.Request) {
    var topic string
    if len(r.URL.Path[1:]) > 0 {
        topic = r.URL.Path[1:]
        event_lookup(topic)
    } else {
        topic = "programming"
    }
    fmt.Fprintf(w, "Hi there, I love %s!", topic)
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}





