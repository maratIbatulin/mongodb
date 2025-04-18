// Package mongo provides types and utilities for working with MongoDB.
package mongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// D is a type alias for a BSON document. It is an ordered representation of a BSON
// document using a Go slice of key-value pairs.
type D primitive.D

// M is a type alias for an unordered BSON document, represented as a Go map with
// string keys.
type M primitive.M

// A is a type alias for a BSON array, represented as a Go slice.
type A primitive.A

// E is a type alias for a single element (key-value pair) in a BSON document.
type E primitive.E
