# PlaceClassData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | [**uuid.UUID**](uuid.UUID.md) | place ID | 
**Type** | **string** |  | 
**Attributes** | [**PlaceClassDataAttributes**](PlaceClassDataAttributes.md) |  | 
**Relationships** | Pointer to [**PlaceClassDataRelationships**](PlaceClassDataRelationships.md) |  | [optional] 

## Methods

### NewPlaceClassData

`func NewPlaceClassData(id uuid.UUID, type_ string, attributes PlaceClassDataAttributes, ) *PlaceClassData`

NewPlaceClassData instantiates a new PlaceClassData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPlaceClassDataWithDefaults

`func NewPlaceClassDataWithDefaults() *PlaceClassData`

NewPlaceClassDataWithDefaults instantiates a new PlaceClassData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *PlaceClassData) GetId() uuid.UUID`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *PlaceClassData) GetIdOk() (*uuid.UUID, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *PlaceClassData) SetId(v uuid.UUID)`

SetId sets Id field to given value.


### GetType

`func (o *PlaceClassData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *PlaceClassData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *PlaceClassData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *PlaceClassData) GetAttributes() PlaceClassDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *PlaceClassData) GetAttributesOk() (*PlaceClassDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *PlaceClassData) SetAttributes(v PlaceClassDataAttributes)`

SetAttributes sets Attributes field to given value.


### GetRelationships

`func (o *PlaceClassData) GetRelationships() PlaceClassDataRelationships`

GetRelationships returns the Relationships field if non-nil, zero value otherwise.

### GetRelationshipsOk

`func (o *PlaceClassData) GetRelationshipsOk() (*PlaceClassDataRelationships, bool)`

GetRelationshipsOk returns a tuple with the Relationships field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRelationships

`func (o *PlaceClassData) SetRelationships(v PlaceClassDataRelationships)`

SetRelationships sets Relationships field to given value.

### HasRelationships

`func (o *PlaceClassData) HasRelationships() bool`

HasRelationships returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


