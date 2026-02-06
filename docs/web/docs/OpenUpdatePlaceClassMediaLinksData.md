# OpenUpdatePlaceClassMediaLinksData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | [**uuid.UUID**](uuid.UUID.md) | Upload session id | 
**Type** | **string** |  | 
**Attributes** | [**OpenUpdatePlaceClassMediaLinksDataAttributes**](OpenUpdatePlaceClassMediaLinksDataAttributes.md) |  | 
**Relationships** | [**OpenUpdatePlaceClassMediaLinksDataRelationships**](OpenUpdatePlaceClassMediaLinksDataRelationships.md) |  | 

## Methods

### NewOpenUpdatePlaceClassMediaLinksData

`func NewOpenUpdatePlaceClassMediaLinksData(id uuid.UUID, type_ string, attributes OpenUpdatePlaceClassMediaLinksDataAttributes, relationships OpenUpdatePlaceClassMediaLinksDataRelationships, ) *OpenUpdatePlaceClassMediaLinksData`

NewOpenUpdatePlaceClassMediaLinksData instantiates a new OpenUpdatePlaceClassMediaLinksData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOpenUpdatePlaceClassMediaLinksDataWithDefaults

`func NewOpenUpdatePlaceClassMediaLinksDataWithDefaults() *OpenUpdatePlaceClassMediaLinksData`

NewOpenUpdatePlaceClassMediaLinksDataWithDefaults instantiates a new OpenUpdatePlaceClassMediaLinksData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *OpenUpdatePlaceClassMediaLinksData) GetId() uuid.UUID`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *OpenUpdatePlaceClassMediaLinksData) GetIdOk() (*uuid.UUID, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *OpenUpdatePlaceClassMediaLinksData) SetId(v uuid.UUID)`

SetId sets Id field to given value.


### GetType

`func (o *OpenUpdatePlaceClassMediaLinksData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *OpenUpdatePlaceClassMediaLinksData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *OpenUpdatePlaceClassMediaLinksData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *OpenUpdatePlaceClassMediaLinksData) GetAttributes() OpenUpdatePlaceClassMediaLinksDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *OpenUpdatePlaceClassMediaLinksData) GetAttributesOk() (*OpenUpdatePlaceClassMediaLinksDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *OpenUpdatePlaceClassMediaLinksData) SetAttributes(v OpenUpdatePlaceClassMediaLinksDataAttributes)`

SetAttributes sets Attributes field to given value.


### GetRelationships

`func (o *OpenUpdatePlaceClassMediaLinksData) GetRelationships() OpenUpdatePlaceClassMediaLinksDataRelationships`

GetRelationships returns the Relationships field if non-nil, zero value otherwise.

### GetRelationshipsOk

`func (o *OpenUpdatePlaceClassMediaLinksData) GetRelationshipsOk() (*OpenUpdatePlaceClassMediaLinksDataRelationships, bool)`

GetRelationshipsOk returns a tuple with the Relationships field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRelationships

`func (o *OpenUpdatePlaceClassMediaLinksData) SetRelationships(v OpenUpdatePlaceClassMediaLinksDataRelationships)`

SetRelationships sets Relationships field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


