# CreatePlaceClassDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ParentId** | Pointer to [**uuid.UUID**](uuid.UUID.md) | The ID of the parent class, if any | [optional] 
**Code** | **string** | The unique code of the class | 
**Name** | **string** | The name of the class | 
**Description** | **string** | A brief description of the class | 
**Icon** | Pointer to **string** | A URL to an icon representing the class | [optional] 

## Methods

### NewCreatePlaceClassDataAttributes

`func NewCreatePlaceClassDataAttributes(code string, name string, description string, ) *CreatePlaceClassDataAttributes`

NewCreatePlaceClassDataAttributes instantiates a new CreatePlaceClassDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreatePlaceClassDataAttributesWithDefaults

`func NewCreatePlaceClassDataAttributesWithDefaults() *CreatePlaceClassDataAttributes`

NewCreatePlaceClassDataAttributesWithDefaults instantiates a new CreatePlaceClassDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetParentId

`func (o *CreatePlaceClassDataAttributes) GetParentId() uuid.UUID`

GetParentId returns the ParentId field if non-nil, zero value otherwise.

### GetParentIdOk

`func (o *CreatePlaceClassDataAttributes) GetParentIdOk() (*uuid.UUID, bool)`

GetParentIdOk returns a tuple with the ParentId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetParentId

`func (o *CreatePlaceClassDataAttributes) SetParentId(v uuid.UUID)`

SetParentId sets ParentId field to given value.

### HasParentId

`func (o *CreatePlaceClassDataAttributes) HasParentId() bool`

HasParentId returns a boolean if a field has been set.

### GetCode

`func (o *CreatePlaceClassDataAttributes) GetCode() string`

GetCode returns the Code field if non-nil, zero value otherwise.

### GetCodeOk

`func (o *CreatePlaceClassDataAttributes) GetCodeOk() (*string, bool)`

GetCodeOk returns a tuple with the Code field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCode

`func (o *CreatePlaceClassDataAttributes) SetCode(v string)`

SetCode sets Code field to given value.


### GetName

`func (o *CreatePlaceClassDataAttributes) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CreatePlaceClassDataAttributes) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CreatePlaceClassDataAttributes) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *CreatePlaceClassDataAttributes) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *CreatePlaceClassDataAttributes) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *CreatePlaceClassDataAttributes) SetDescription(v string)`

SetDescription sets Description field to given value.


### GetIcon

`func (o *CreatePlaceClassDataAttributes) GetIcon() string`

GetIcon returns the Icon field if non-nil, zero value otherwise.

### GetIconOk

`func (o *CreatePlaceClassDataAttributes) GetIconOk() (*string, bool)`

GetIconOk returns a tuple with the Icon field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIcon

`func (o *CreatePlaceClassDataAttributes) SetIcon(v string)`

SetIcon sets Icon field to given value.

### HasIcon

`func (o *CreatePlaceClassDataAttributes) HasIcon() bool`

HasIcon returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


