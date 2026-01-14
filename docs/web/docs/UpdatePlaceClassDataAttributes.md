# UpdatePlaceClassDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ParentId** | Pointer to [**uuid.UUID**](uuid.UUID.md) | The ID of the parent class | [optional] 
**Code** | Pointer to **string** | The unique code of the class | [optional] 
**Name** | Pointer to **string** | The name of the class | [optional] 
**Description** | Pointer to **string** | A brief description of the class | [optional] 
**Icon** | Pointer to **string** | A URL to an icon representing the class | [optional] 

## Methods

### NewUpdatePlaceClassDataAttributes

`func NewUpdatePlaceClassDataAttributes() *UpdatePlaceClassDataAttributes`

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

### GetCode

`func (o *UpdatePlaceClassDataAttributes) GetCode() string`

GetCode returns the Code field if non-nil, zero value otherwise.

### GetCodeOk

`func (o *UpdatePlaceClassDataAttributes) GetCodeOk() (*string, bool)`

GetCodeOk returns a tuple with the Code field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCode

`func (o *UpdatePlaceClassDataAttributes) SetCode(v string)`

SetCode sets Code field to given value.

### HasCode

`func (o *UpdatePlaceClassDataAttributes) HasCode() bool`

HasCode returns a boolean if a field has been set.

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

### HasName

`func (o *UpdatePlaceClassDataAttributes) HasName() bool`

HasName returns a boolean if a field has been set.

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

### HasDescription

`func (o *UpdatePlaceClassDataAttributes) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetIcon

`func (o *UpdatePlaceClassDataAttributes) GetIcon() string`

GetIcon returns the Icon field if non-nil, zero value otherwise.

### GetIconOk

`func (o *UpdatePlaceClassDataAttributes) GetIconOk() (*string, bool)`

GetIconOk returns a tuple with the Icon field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIcon

`func (o *UpdatePlaceClassDataAttributes) SetIcon(v string)`

SetIcon sets Icon field to given value.

### HasIcon

`func (o *UpdatePlaceClassDataAttributes) HasIcon() bool`

HasIcon returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


