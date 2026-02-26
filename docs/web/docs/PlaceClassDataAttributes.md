# PlaceClassDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | The name of the class | 
**Description** | **string** | A brief description of the class | 
**IconKey** | Pointer to **string** | A URL to an icon representing the class | [optional] 
**Version** | **int32** | The version number of the class data, used for concurrency control | 
**CreatedAt** | **time.Time** | The date and time when the class was created | 
**UpdatedAt** | **time.Time** | The date and time when the class was last updated | 
**DeprecatedAt** | Pointer to **time.Time** | The date and time when the class was deprecated, if applicable | [optional] 

## Methods

### NewPlaceClassDataAttributes

`func NewPlaceClassDataAttributes(name string, description string, version int32, createdAt time.Time, updatedAt time.Time, ) *PlaceClassDataAttributes`

NewPlaceClassDataAttributes instantiates a new PlaceClassDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPlaceClassDataAttributesWithDefaults

`func NewPlaceClassDataAttributesWithDefaults() *PlaceClassDataAttributes`

NewPlaceClassDataAttributesWithDefaults instantiates a new PlaceClassDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *PlaceClassDataAttributes) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *PlaceClassDataAttributes) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *PlaceClassDataAttributes) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *PlaceClassDataAttributes) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *PlaceClassDataAttributes) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *PlaceClassDataAttributes) SetDescription(v string)`

SetDescription sets Description field to given value.


### GetIconKey

`func (o *PlaceClassDataAttributes) GetIconKey() string`

GetIconKey returns the IconKey field if non-nil, zero value otherwise.

### GetIconKeyOk

`func (o *PlaceClassDataAttributes) GetIconKeyOk() (*string, bool)`

GetIconKeyOk returns a tuple with the IconKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIconKey

`func (o *PlaceClassDataAttributes) SetIconKey(v string)`

SetIconKey sets IconKey field to given value.

### HasIconKey

`func (o *PlaceClassDataAttributes) HasIconKey() bool`

HasIconKey returns a boolean if a field has been set.

### GetVersion

`func (o *PlaceClassDataAttributes) GetVersion() int32`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *PlaceClassDataAttributes) GetVersionOk() (*int32, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *PlaceClassDataAttributes) SetVersion(v int32)`

SetVersion sets Version field to given value.


### GetCreatedAt

`func (o *PlaceClassDataAttributes) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *PlaceClassDataAttributes) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *PlaceClassDataAttributes) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetUpdatedAt

`func (o *PlaceClassDataAttributes) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *PlaceClassDataAttributes) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *PlaceClassDataAttributes) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.


### GetDeprecatedAt

`func (o *PlaceClassDataAttributes) GetDeprecatedAt() time.Time`

GetDeprecatedAt returns the DeprecatedAt field if non-nil, zero value otherwise.

### GetDeprecatedAtOk

`func (o *PlaceClassDataAttributes) GetDeprecatedAtOk() (*time.Time, bool)`

GetDeprecatedAtOk returns a tuple with the DeprecatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeprecatedAt

`func (o *PlaceClassDataAttributes) SetDeprecatedAt(v time.Time)`

SetDeprecatedAt sets DeprecatedAt field to given value.

### HasDeprecatedAt

`func (o *PlaceClassDataAttributes) HasDeprecatedAt() bool`

HasDeprecatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


