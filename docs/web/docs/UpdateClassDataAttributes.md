# UpdateClassDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ParentId** | Pointer to [**uuid.UUID**](uuid.UUID.md) | The ID of the parent place | [optional] 
**Code** | Pointer to **string** | The unique code of the class | [optional] 
**Name** | Pointer to **string** | The name of the class | [optional] 
**Description** | Pointer to **string** | A brief description of the class | [optional] 
**Icon** | Pointer to **string** | A URL to an icon representing the class | [optional] 

## Methods

### NewUpdateClassDataAttributes

`func NewUpdateClassDataAttributes() *UpdateClassDataAttributes`

NewUpdateClassDataAttributes instantiates a new UpdateClassDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateClassDataAttributesWithDefaults

`func NewUpdateClassDataAttributesWithDefaults() *UpdateClassDataAttributes`

NewUpdateClassDataAttributesWithDefaults instantiates a new UpdateClassDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetParentId

`func (o *UpdateClassDataAttributes) GetParentId() uuid.UUID`

GetParentId returns the ParentId field if non-nil, zero value otherwise.

### GetParentIdOk

`func (o *UpdateClassDataAttributes) GetParentIdOk() (*uuid.UUID, bool)`

GetParentIdOk returns a tuple with the ParentId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetParentId

`func (o *UpdateClassDataAttributes) SetParentId(v uuid.UUID)`

SetParentId sets ParentId field to given value.

### HasParentId

`func (o *UpdateClassDataAttributes) HasParentId() bool`

HasParentId returns a boolean if a field has been set.

### GetCode

`func (o *UpdateClassDataAttributes) GetCode() string`

GetCode returns the Code field if non-nil, zero value otherwise.

### GetCodeOk

`func (o *UpdateClassDataAttributes) GetCodeOk() (*string, bool)`

GetCodeOk returns a tuple with the Code field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCode

`func (o *UpdateClassDataAttributes) SetCode(v string)`

SetCode sets Code field to given value.

### HasCode

`func (o *UpdateClassDataAttributes) HasCode() bool`

HasCode returns a boolean if a field has been set.

### GetName

`func (o *UpdateClassDataAttributes) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *UpdateClassDataAttributes) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *UpdateClassDataAttributes) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *UpdateClassDataAttributes) HasName() bool`

HasName returns a boolean if a field has been set.

### GetDescription

`func (o *UpdateClassDataAttributes) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *UpdateClassDataAttributes) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *UpdateClassDataAttributes) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *UpdateClassDataAttributes) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetIcon

`func (o *UpdateClassDataAttributes) GetIcon() string`

GetIcon returns the Icon field if non-nil, zero value otherwise.

### GetIconOk

`func (o *UpdateClassDataAttributes) GetIconOk() (*string, bool)`

GetIconOk returns a tuple with the Icon field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIcon

`func (o *UpdateClassDataAttributes) SetIcon(v string)`

SetIcon sets Icon field to given value.

### HasIcon

`func (o *UpdateClassDataAttributes) HasIcon() bool`

HasIcon returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


