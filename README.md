# DDBMoose
Golang package for operations with AWS DynamoDB using SDK version 2.0. It facilitates data manipulation, using its structures to save or read your project data.

## Install

To use "ddbmoose" in your Go project, you must first install it. Use the `go get` command to do this:

    go get github.com/fenix-ds/ddbmoose

## Use

Here is an example of how to use "ddbmoose" to interact with DynamoDB:

### 1. Creating connection to the Aws DynamoDB table

```
package main

import (
    "fmt"
    "github.com/fenix-ds/ddbmoose"
)

func main() {
    //Create an instance of "ddbmoose" for a specific AWS region.
    ddb, err := ddbmoose.DdbMooseCreate("us-west-2")
    if err != nil {
        fmt.Println("Error creating a DdbMoose instance:", err)
        return
    }

    // Set the name of the table you want to interact with.
    err = ddb.SetTable("NameOfYourTable")
        
    if err != nil {
        fmt.Println("Error defining table name:", err)
        return
    }
}
```

### 2. Saved data from your structure

```
package main

type Person struct {
	Id   string `json:"id" dynamodbav:"id"`
	Name string `json:"name" dynamodbav:"name"`
}

import (
    "fmt"
    "github.com/fenix-ds/ddbmoose"
)

func main() {
    ...

    id := uuid.New().String()
	newPerson, failed := ddbMoose.Save(Person{
		Id:   id,
		Name: "Mauricio",
	})

	if failed != nil {
		fmt.Println(failed.Error())
	}

	fmt.Println(newPerson)
}
```

To update, just send the record as an integer, and it will update all the data

### 3. Listing data with filters

```
package main

import (
    "fmt"
    "github.com/fenix-ds/ddbmoose"
)

func main() {
    ...

    resultFilters, failed = ddbMoose.FindWithFilters(&[]DdbMooseFilter{
		{Field: "name", Operation: TfrContains, Value: "Mau", LogicalOperator: TloNone},
	})

	if failed != nil {
		fmt.Println(failed.Error())
	}

	fmt.Println(resultFilters)
}
```