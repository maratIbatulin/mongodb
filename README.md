# MongoDB Wrapper

A lightweight and fully-featured wrapper for the official MongoDB Go driver that simplifies working with MongoDB in Go applications.

## Features

- Easy connection setup with fluent interface
- Simplified collection operations
- Support for transactions
- Rich query aggregation pipeline builder
- Improved error handling
- Context-aware operations

## Installation

```bash
go get github.com/maratIbatulin/mongodb
```

### Available Aggregation Stages

The library supports all MongoDB aggregation pipeline stages:

- `AddFields`: Add computed fields
- `Bucket`: Categorize documents into buckets
- `BucketAuto`: Automatically categorize documents into buckets
- `Count`: Count documents in the pipeline
- `Densify`: Create additional documents to fill gaps
- `Facet`: Process multiple aggregation pipelines
- `Fill`: Fill missing values in documents
- `GeoNear`: Find documents near a geographic location
- `GraphLookup`: Recursive lookup for graph-like relationships
- `Group`: Group documents by a specified expression
- `Limit`: Limit the number of documents
- `Lookup`: Perform a join with another collection
- `Match`: Filter documents
- `Project`: Reshape documents
- `Sample`: Randomly select documents
- `Set`: Set values in documents
- `Skip`: Skip documents
- `Sort`: Sort documents
- `Unwind`: Deconstruct arrays into separate documents
- And more...
