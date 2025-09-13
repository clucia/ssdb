# SSDB - Google Sheets Database Library
A Go library that provides a database-like interface for Google Sheets, allowing you to read, write, and manipulate spreadsheet data programmatically.

# Motivation
I am running a website that uses Google Sheets as a backend database.

This library allows my website to update the Sheets that are behind the database. There are still many admin tasks that I do manually in Google Sheets, but this library allows me to automate many of the repetitive tasks.

At some point, when all of the functions can be done through the website, I may migrate away from Google Sheets to a more traditional database. But for now, this works well.

# Status
This is used in one production system but should be considered alpha quality. Issues and contributions are welcome.

# Overview
SSDB treats Google Sheets as a database, providing structured access to spreadsheet data through familiar database-like operations. It supports reading, writing, updating, and querying data with automatic handling of Google Sheets API interactions.

SSDB keeps a local cache of the spreadsheet data to minimize API calls and improve performance. It provides utilities for managing ranges, batch updates, and type conversions. For really huge spreadheets this can use a lot of memory, but for most use cases it should be fine.

When you make updates, SSDB uses a transaction-like model to ensure data consistency. You create an Updater object, make your changes, and then call Sync() to apply them. If the underlying data has changed since you loaded it, the Sync() will fail to prevent overwriting changes made by others.

Sync() batches multiple updates into a single API call to write the data to improve efficiency and reduce the chance of hitting rate limits.

Sync() then reads the data back in a single query. Assuming the read-back succeeds, the the result is merged into the local cache. This ensures that the local cache is always in sync with the remote data.

# Features

## Database-like Interface: Access Google Sheets data using familiar database patterns
* Batch Operations: Efficient batch updates and reads
* Range Management: Flexible range selection and manipulation using A1 notation
* Type Conversion: Automatic handling of string, numeric, and boolean data types
* Specialized Modules:

## Sub-modules
* sslist: Append-only list operations
* sstable: Table-like data access with lookups
* sslog: Logging functionality
* ssaudit: Audit trail management
* sschat: Chat/conversation data handling

# Installation
bash
go get github.com/clucia/ssdb

# Quick Start
```go  
package main

import (
    "context"
    "log"
    
    "github.com/clucia/ssdb"
)

func main() {
    ctx := context.Background()
    
    // Load service account credentials
    credentials := []byte(`your-service-account-json`)
    spreadsheetID := "your-spreadsheet-id"
    
    // Open the database connection
    db, err := ssdb.Open(ctx, spreadsheetID, credentials)
    if err != nil {
        log.Fatal(err)
    }
    
    // Load spreadsheet data
    if err := db.Loader(ctx); err != nil {
        log.Fatal(err)
    }
    
    // Iterate through sheets
    db.SheetIter(func(sheetName string, sheet *ssdb.Sheet) {
        log.Printf("Sheet: %s", sheetName)
    })
}
```

# Core Components
## SSDB
The main database connection object that manages the Google Sheets API connection and provides access to spreadsheet data.

## Sheet, Row, Cell
Hierarchical data structures representing spreadsheet elements:

* Sheet: Represents a worksheet
* Row: Represents a row within a sheet
* Cell: Represents an individual cell

## DBRange
Handles range operations using A1 notation (e.g., "A1") and provides utilities for range manipulation.
## Updater
Manages batch updates to ensure data consistency and efficient API usage.

# Key Operations
## Reading Data
The table metaphor is very useful, so I'll demonstrate reading data in that way.

```go
    // Open the SSDB connection
	ssdbHandle, err = ssdb.Open(ctx, spreadsheetID, credentials)
	// Handle error

    // Load the cached copy of the spreadsheet
   	err = ssdbHandle.Loader(ctx)
	// Handle error

    // Open the config table
    configTableHandle, err = sstable.Open(ssdbHandle, "config")
	// Handle error

    // Lookup specific config values
    // This assumes that the config table has a row with "Maintenance Mode" in the first column
    // and has a column with "My AppName" in the first row.
   	maintMode, err = configTableHandle.Lookup("Maintenance Mode", "My App Name")
	// Handle error
```


## Writing Data
// This will show how to update the maintMode cell we read above:
```go
	updater := hState.ssdb.NewUpdater()
	updater.Update(maintMode.Range(), [][]any{{"TRUE"}})
    err = updater.Sync()
    // Handle error
```

## Search for data
```go
row := sheet.SearchV(true, func(r *ssdb.Row) bool {
    return r.GetCellN(0).GetString() == "target_value"
})
```

## Create ranges using A1 notation
```go
dbrange := db.NewDBRangeFromSymbolicRange("Sheet1!A1:B10")
```

## Convert between different range formats
```go
rangeString := dbrange.String()
gridRange := dbrange.gridRange
```


# Specialized Modules
## SSList - Append-only Operations
```go
import "github.com/clucia/ssdb/sslist"

list, err := sslist.Open(db, "LogSheet")
updater := db.NewUpdater()
data := [][]any{{"New", "Entry"}}
list.AppendBlank(updater, data)
updater.Sync()
```
## SSLog - Logging Functionality
```go
import "github.com/clucia/ssdb/sslog"

logger, err := sslog.Open(db, "LogSheet")
logger.LogWithData(ctx, "operation", "completed", time.Now())
```

## Range Format Support
The library supports various Google Sheets range formats:

* Single cells: A1, B5
* Ranges: A1:B10, C1:Z100
* With sheet names: Sheet1!A1:B10
* Quoted sheet names: 'My Sheet'!A1:B10
* Full rows/columns: 1:5, A:C

## Authentication
SSDB uses Google Service Account authentication. You'll need:

* A Google Cloud Project with Sheets API enabled
* A service account with appropriate permissions
* The service account JSON key file
* The spreadsheet shared with the service account email

## Error Handling
The library provides structured error handling:
```go
if err != nil {
    switch {
    case errors.Is(err, sslist.ErrSheetNotFound):
        // Handle sheet not found
    case errors.Is(err, sstable.ErrLookupFailed):
        // Handle lookup failure
    default:
        // Handle other errors
    }
}
```

## Performance Considerations

* Use batch operations (Updater) for multiple updates
* The library automatically handles API rate limits
* Data is cached locally after loading
* Use ReloadDBGet() to refresh cached data

## License
This library is designed for programmatic access to Google Sheets data and requires appropriate Google API credentials and permissions.

# MIT License
This project is licensed under the MIT License. See the LICENSE file for details.
