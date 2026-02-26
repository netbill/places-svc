# UpdatePlaceClassDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ParentId** | Pointer to [**uuid.UUID**](uuid.UUID.md) | The ID of the parent class, if any | [optional] 
**Name** | **string** | The name of the class | 
**Description** | **string** | A brief description of the class | 
**IconKey** | Pointer to **string** | The S3 key for the class icon | [optional] 

## Methods

### NewUpdatePlaceClassDataAttributes

`func NewUpdatePlaceClassDataAttributes(name string, description string, ) *UpdatePlaceClassDataAttributes`

NewUpdatePlaceClassDataAttributes instantiates a new UpdatePlaceClassDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdatePlaceClassDataAttributesWithDefaults

`func NewUpdatePlaceClassDataAttributesWithDefaults() *UpdatePlaceClassDataAttributes`

NewUpdatePlaceClassDataAttributesWithDefaults instantiates a new UpdatePlaceClassDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetParentId

`func (o *UpdatePlaceClassDataAttributes) GetParentId() uuid.UUID`

GetParentId returns the ParentId field if non-nil, zero value otherwise.

### GetParentIdOk

`func (o *UpdatePlaceClassDataAttributes) GetParentIdOk() (*uuid.UUID, bool)`

GetParentIdOk returns a tuple with the ParentId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetParentId

`func (o *UpdatePlaceClassDataAttributes) SetParentId(v uuid.UUID)`

SetParentId sets ParentId field to given value.

### HasParentId

`func (o *UpdatePlaceClassDataAttributes) HasParentId() bool`

HasParentId returns a boolean if a field has been set.

### GetName

`func (o *UpdatePlaceClassDataAttributes) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *UpdatePlaceClassDataAttributes) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *UpdatePlaceClassDataAttributes) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *UpdatePlaceClassDataAttributes) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *UpdatePlaceClassDataAttributes) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *UpdatePlaceClassDataAttributes) SetDescription(v string)`

SetDescription sets Description field to given value.


### GetIconKey

`func (o *UpdatePlaceClassDataAttributes) GetIconKey() string`

GetIconKey returns the IconKey field if non-nil, zero value otherwise.

### GetIconKeyOk

`func (o *UpdatePlaceClassDataAttributes) GetIconKeyOk() (*string, bool)`

GetIconKeyOk returns a tuple with the IconKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIconKey

`func (o *UpdatePlaceClassDataAttributes) SetIconKey(v string)`

SetIconKey sets IconKey field to given value.

### HasIconKey

`func (o *UpdatePlaceClassDataAttributes) HasIconKey() bool`

HasIconKey returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


