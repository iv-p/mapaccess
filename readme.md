# mapaccess

A small golang library to retrieve arbitrary keys from golang interface{}. Think of it as acessing JSON keys from arbitrary interfaces{}. It is heavily influenced by the golang temlpate engine.

![](https://github.com/iv-p/mapaccess/workflows/test/badge.svg)

## Installation

To install mapaccess just run the following command in your terminal
```
go get -u github.com/iv-p/mapaccess
```

## Usage

mapaccess exposes only one function, which takes a string key and a interface{}:
```go
result, err := mapaccess.Get("key.one[0].two", data)
```
The key should represent a JSON type location of the data in the interface{}. It is intended to work only with basic interfaces - only map[string]interface{} and []interface{}, it won't work if you pass arbitrary structure.

It is intended towards using with serializing arbitary JSON and getting an arbitrary key out of it like so:
```go
j := []byte(`{
    "id": "9b92b11b-b57f-4fa6-af5e-e35a290dc764",	
    "name": "John Doe",
    "friends": [
        {
            "name": "Jaime Mckinney"
        },
        {
            "name": "Evangeline Alvarado"
        },
        {
            "name": "Beth Cantrell"
        }
    ]
}`)
var deserialised interface{}
err := json.Unmarshal(j, &deserialised)
if err != nil {
    panic(err)
}

bestFriendName, err := mapaccess.Get("friends[0].name", deserialised)
if err != nil {
    panic(err)
}

myName, err := mapaccess.Get("name", deserialised)
if err != nil {
    panic(err)
}

log.Printf("My name is %s and my best friend's name is %s",
    myName, bestFriendName)
```

Running the above snippet produces the following output
```
My name is John Doe and my best friend's name is Jaime Mckinney
```

## Valid keys

There is a limitation of the alphabet for keys, it includes underscores, dashes and alphanumeric characters.