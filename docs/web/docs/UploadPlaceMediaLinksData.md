# UploadPlaceMediaLinksData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | [**uuid.UUID**](uuid.UUID.md) | place id | 
**Type** | **string** |  | 
**Attributes** | [**UploadPlaceMediaLinksDataAttributes**](UploadPlaceMediaLinksDataAttributes.md) |  | 
**Relationships** | [**UploadPlaceMediaLinksDataRelationships**](UploadPlaceMediaLinksDataRelationships.md) |  | 

## Methods

### NewUploadPlaceMediaLinksData

`func NewUploadPlaceMediaLinksData(id uuid.UUID, type_ string, attributes UploadPlaceMediaLinksDataAttributes, relationships UploadPlaceMediaLinksDataRelationships, ) *UploadPlaceMediaLinksData`

NewUploadPlaceMediaLinksData instantiates a new UploadPlaceMediaLinksData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUploadPlaceMediaLinksDataWithDefaults

`func NewUploadPlaceMediaLinksDataWithDefaults() *UploadPlaceMediaLinksData`

NewUploadPlaceMediaLinksDataWithDefaults instantiates a new UploadPlaceMediaLinksData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *UploadPlaceMediaLinksData) GetId() uuid.UUID`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *UploadPlaceMediaLinksData) GetIdOk() (*uuid.UUID, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *UploadPlaceMediaLinksData) SetId(v uuid.UUID)`

SetId sets Id field to given value.


### GetType

`func (o *UploadPlaceMediaLinksData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *UploadPlaceMediaLinksData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *UploadPlaceMediaLinksData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *UploadPlaceMediaLinksData) GetAttributes() UploadPlaceMediaLinksDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *UploadPlaceMediaLinksData) GetAttributesOk() (*UploadPlaceMediaLinksDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *UploadPlaceMediaLinksData) SetAttributes(v UploadPlaceMediaLinksDataAttributes)`

SetAttributes sets Attributes field to given value.


### GetRelationships

`func (o *UploadPlaceMediaLinksData) GetRelationships() UploadPlaceMediaLinksDataRelationships`

GetRelationships returns the Relationships field if non-nil, zero value otherwise.

### GetRelationshipsOk

`func (o *UploadPlaceMediaLinksData) GetRelationshipsOk() (*UploadPlaceMediaLinksDataRelationships, bool)`

GetRelationshipsOk returns a tuple with the Relationships field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRelationships

`func (o *UploadPlaceMediaLinksData) SetRelationships(v UploadPlaceMediaLinksDataRelationships)`

SetRelationships sets Relationships field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


