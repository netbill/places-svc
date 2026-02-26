# PlaceIncludedInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | [**uuid.UUID**](uuid.UUID.md) | place ID | 
**Type** | **string** |  | 
**Attributes** | [**PlaceClassDataAttributes**](PlaceClassDataAttributes.md) |  | 
**Relationships** | Pointer to [**PlaceClassDataRelationships**](PlaceClassDataRelationships.md) |  | [optional] 

## Methods

### NewPlaceIncludedInner

`func NewPlaceIncludedInner(id uuid.UUID, type_ string, attributes PlaceClassDataAttributes, ) *PlaceIncludedInner`

NewPlaceIncludedInner instantiates a new PlaceIncludedInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPlaceIncludedInnerWithDefaults

`func NewPlaceIncludedInnerWithDefaults() *PlaceIncludedInner`

NewPlaceIncludedInnerWithDefaults instantiates a new PlaceIncludedInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *PlaceIncludedInner) GetId() uuid.UUID`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *PlaceIncludedInner) GetIdOk() (*uuid.UUID, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *PlaceIncludedInner) SetId(v uuid.UUID)`

SetId sets Id field to given value.


### GetType

`func (o *PlaceIncludedInner) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *PlaceIncludedInner) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *PlaceIncludedInner) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *PlaceIncludedInner) GetAttributes() PlaceClassDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *PlaceIncludedInner) GetAttributesOk() (*PlaceClassDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *PlaceIncludedInner) SetAttributes(v PlaceClassDataAttributes)`

SetAttributes sets Attributes field to given value.


### GetRelationships

`func (o *PlaceIncludedInner) GetRelationships() PlaceClassDataRelationships`

GetRelationships returns the Relationships field if non-nil, zero value otherwise.

### GetRelationshipsOk

`func (o *PlaceIncludedInner) GetRelationshipsOk() (*PlaceClassDataRelationships, bool)`

GetRelationshipsOk returns a tuple with the Relationships field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRelationships

`func (o *PlaceIncludedInner) SetRelationships(v PlaceClassDataRelationships)`

SetRelationships sets Relationships field to given value.

### HasRelationships

`func (o *PlaceIncludedInner) HasRelationships() bool`

HasRelationships returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


