# OpenUpdatePlaceMediaLinksData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | [**uuid.UUID**](uuid.UUID.md) | Upload session id | 
**Type** | **string** |  | 
**Attributes** | [**OpenUpdatePlaceMediaLinksDataAttributes**](OpenUpdatePlaceMediaLinksDataAttributes.md) |  | 
**Relationships** | [**OpenUpdatePlaceMediaLinksDataRelationships**](OpenUpdatePlaceMediaLinksDataRelationships.md) |  | 

## Methods

### NewOpenUpdatePlaceMediaLinksData

`func NewOpenUpdatePlaceMediaLinksData(id uuid.UUID, type_ string, attributes OpenUpdatePlaceMediaLinksDataAttributes, relationships OpenUpdatePlaceMediaLinksDataRelationships, ) *OpenUpdatePlaceMediaLinksData`

NewOpenUpdatePlaceMediaLinksData instantiates a new OpenUpdatePlaceMediaLinksData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOpenUpdatePlaceMediaLinksDataWithDefaults

`func NewOpenUpdatePlaceMediaLinksDataWithDefaults() *OpenUpdatePlaceMediaLinksData`

NewOpenUpdatePlaceMediaLinksDataWithDefaults instantiates a new OpenUpdatePlaceMediaLinksData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *OpenUpdatePlaceMediaLinksData) GetId() uuid.UUID`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *OpenUpdatePlaceMediaLinksData) GetIdOk() (*uuid.UUID, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *OpenUpdatePlaceMediaLinksData) SetId(v uuid.UUID)`

SetId sets Id field to given value.


### GetType

`func (o *OpenUpdatePlaceMediaLinksData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *OpenUpdatePlaceMediaLinksData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *OpenUpdatePlaceMediaLinksData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *OpenUpdatePlaceMediaLinksData) GetAttributes() OpenUpdatePlaceMediaLinksDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *OpenUpdatePlaceMediaLinksData) GetAttributesOk() (*OpenUpdatePlaceMediaLinksDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *OpenUpdatePlaceMediaLinksData) SetAttributes(v OpenUpdatePlaceMediaLinksDataAttributes)`

SetAttributes sets Attributes field to given value.


### GetRelationships

`func (o *OpenUpdatePlaceMediaLinksData) GetRelationships() OpenUpdatePlaceMediaLinksDataRelationships`

GetRelationships returns the Relationships field if non-nil, zero value otherwise.

### GetRelationshipsOk

`func (o *OpenUpdatePlaceMediaLinksData) GetRelationshipsOk() (*OpenUpdatePlaceMediaLinksDataRelationships, bool)`

GetRelationshipsOk returns a tuple with the Relationships field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRelationships

`func (o *OpenUpdatePlaceMediaLinksData) SetRelationships(v OpenUpdatePlaceMediaLinksDataRelationships)`

SetRelationships sets Relationships field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


