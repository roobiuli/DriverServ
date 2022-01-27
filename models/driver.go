package models

// I choose to use go vendor and have the models used by both services.
//Here should be placed future modules for this test project if needed.

// Driver Struct (Used by Driver Service)

type Driver struct {
	Id string `json:"id"`
	Name string `json:"name"`
}
