# PlaceClass

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | [**PlaceClassData**](PlaceClassData.md) |  | 
**Included** | Pointer to [**[]PlaceClassData**](PlaceClassData.md) | Included related resources (e.g., organization) | [optional] 
**Links** | Pointer to [**PaginationData**](PaginationData.md) |  | [optional] 

## Methods

### NewPlaceClass

`func NewPlaceClass(data PlaceClassData, ) *PlaceClass`

NewPlaceClass instantiates a new PlaceClass object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPlaceClassWithDefaults

`func NewPlaceClassWithDefaults() *PlaceClass`

NewPlaceClassWithDefaults instantiates a new PlaceClass object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *PlaceClass) GetData() PlaceClassData`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *PlaceClass) GetDataOk() (*PlaceClassData, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *PlaceClass) SetData(v PlaceClassData)`

SetData sets Data field to given value.


### GetIncluded

`func (o *PlaceClass) GetIncluded() []PlaceClassData`

GetIncluded returns the Included field if non-nil, zero value otherwise.

### GetIncludedOk

`func (o *PlaceClass) GetIncludedOk() (*[]PlaceClassData, bool)`

GetIncludedOk returns a tuple with the Included field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIncluded

`func (o *PlaceClass) SetIncluded(v []PlaceClassData)`

SetIncluded sets Included field to given value.

### HasIncluded

`func (o *PlaceClass) HasIncluded() bool`

HasIncluded returns a boolean if a field has been set.

### GetLinks

`func (o *PlaceClass) GetLinks() PaginationData`

GetLinks returns the Links field if non-nil, zero value otherwise.

### GetLinksOk

`func (o *PlaceClass) GetLinksOk() (*PaginationData, bool)`

GetLinksOk returns a tuple with the Links field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLinks

`func (o *PlaceClass) SetLinks(v PaginationData)`

SetLinks sets Links field to given value.

### HasLinks

`func (o *PlaceClass) HasLinks() bool`

HasLinks returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


