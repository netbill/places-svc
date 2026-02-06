# ConfirmUpdatePlaceClassData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | [**uuid.UUID**](uuid.UUID.md) | place ID | 
**Type** | **string** |  | 
**Attributes** | [**ConfirmUpdatePlaceClassDataAttributes**](ConfirmUpdatePlaceClassDataAttributes.md) |  | 

## Methods

### NewConfirmUpdatePlaceClassData

`func NewConfirmUpdatePlaceClassData(id uuid.UUID, type_ string, attributes ConfirmUpdatePlaceClassDataAttributes, ) *ConfirmUpdatePlaceClassData`

NewConfirmUpdatePlaceClassData instantiates a new ConfirmUpdatePlaceClassData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfirmUpdatePlaceClassDataWithDefaults

`func NewConfirmUpdatePlaceClassDataWithDefaults() *ConfirmUpdatePlaceClassData`

NewConfirmUpdatePlaceClassDataWithDefaults instantiates a new ConfirmUpdatePlaceClassData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *ConfirmUpdatePlaceClassData) GetId() uuid.UUID`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ConfirmUpdatePlaceClassData) GetIdOk() (*uuid.UUID, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ConfirmUpdatePlaceClassData) SetId(v uuid.UUID)`

SetId sets Id field to given value.


### GetType

`func (o *ConfirmUpdatePlaceClassData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ConfirmUpdatePlaceClassData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ConfirmUpdatePlaceClassData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *ConfirmUpdatePlaceClassData) GetAttributes() ConfirmUpdatePlaceClassDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *ConfirmUpdatePlaceClassData) GetAttributesOk() (*ConfirmUpdatePlaceClassDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *ConfirmUpdatePlaceClassData) SetAttributes(v ConfirmUpdatePlaceClassDataAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


