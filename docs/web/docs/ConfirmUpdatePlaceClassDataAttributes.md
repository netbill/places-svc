# ConfirmUpdatePlaceClassDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ParentId** | Pointer to [**uuid.UUID**](uuid.UUID.md) | The ID of the parent class | [optional] 
**Code** | **string** | The unique code of the class | 
**Name** | **string** | The name of the class | 
**Description** | **string** | A brief description of the class | 
**DeleteIcon** | Pointer to **bool** | Flag to indicate if the place class banner should be deleted | [optional] 

## Methods

### NewConfirmUpdatePlaceClassDataAttributes

`func NewConfirmUpdatePlaceClassDataAttributes(code string, name string, description string, ) *ConfirmUpdatePlaceClassDataAttributes`

NewConfirmUpdatePlaceClassDataAttributes instantiates a new ConfirmUpdatePlaceClassDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfirmUpdatePlaceClassDataAttributesWithDefaults

`func NewConfirmUpdatePlaceClassDataAttributesWithDefaults() *ConfirmUpdatePlaceClassDataAttributes`

NewConfirmUpdatePlaceClassDataAttributesWithDefaults instantiates a new ConfirmUpdatePlaceClassDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetParentId

`func (o *ConfirmUpdatePlaceClassDataAttributes) GetParentId() uuid.UUID`

GetParentId returns the ParentId field if non-nil, zero value otherwise.

### GetParentIdOk

`func (o *ConfirmUpdatePlaceClassDataAttributes) GetParentIdOk() (*uuid.UUID, bool)`

GetParentIdOk returns a tuple with the ParentId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetParentId

`func (o *ConfirmUpdatePlaceClassDataAttributes) SetParentId(v uuid.UUID)`

SetParentId sets ParentId field to given value.

### HasParentId

`func (o *ConfirmUpdatePlaceClassDataAttributes) HasParentId() bool`

HasParentId returns a boolean if a field has been set.

### GetCode

`func (o *ConfirmUpdatePlaceClassDataAttributes) GetCode() string`

GetCode returns the Code field if non-nil, zero value otherwise.

### GetCodeOk

`func (o *ConfirmUpdatePlaceClassDataAttributes) GetCodeOk() (*string, bool)`

GetCodeOk returns a tuple with the Code field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCode

`func (o *ConfirmUpdatePlaceClassDataAttributes) SetCode(v string)`

SetCode sets Code field to given value.


### GetName

`func (o *ConfirmUpdatePlaceClassDataAttributes) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfirmUpdatePlaceClassDataAttributes) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfirmUpdatePlaceClassDataAttributes) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *ConfirmUpdatePlaceClassDataAttributes) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *ConfirmUpdatePlaceClassDataAttributes) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *ConfirmUpdatePlaceClassDataAttributes) SetDescription(v string)`

SetDescription sets Description field to given value.


### GetDeleteIcon

`func (o *ConfirmUpdatePlaceClassDataAttributes) GetDeleteIcon() bool`

GetDeleteIcon returns the DeleteIcon field if non-nil, zero value otherwise.

### GetDeleteIconOk

`func (o *ConfirmUpdatePlaceClassDataAttributes) GetDeleteIconOk() (*bool, bool)`

GetDeleteIconOk returns a tuple with the DeleteIcon field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeleteIcon

`func (o *ConfirmUpdatePlaceClassDataAttributes) SetDeleteIcon(v bool)`

SetDeleteIcon sets DeleteIcon field to given value.

### HasDeleteIcon

`func (o *ConfirmUpdatePlaceClassDataAttributes) HasDeleteIcon() bool`

HasDeleteIcon returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


