# ClassDataRelationshipsParent

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | [**uuid.UUID**](uuid.UUID.md) | The ID of the parent class | 
**Type** | **string** |  | 

## Methods

### NewClassDataRelationshipsParent

`func NewClassDataRelationshipsParent(id uuid.UUID, type_ string, ) *ClassDataRelationshipsParent`

NewClassDataRelationshipsParent instantiates a new ClassDataRelationshipsParent object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewClassDataRelationshipsParentWithDefaults

`func NewClassDataRelationshipsParentWithDefaults() *ClassDataRelationshipsParent`

NewClassDataRelationshipsParentWithDefaults instantiates a new ClassDataRelationshipsParent object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *ClassDataRelationshipsParent) GetId() uuid.UUID`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ClassDataRelationshipsParent) GetIdOk() (*uuid.UUID, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ClassDataRelationshipsParent) SetId(v uuid.UUID)`

SetId sets Id field to given value.


### GetType

`func (o *ClassDataRelationshipsParent) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ClassDataRelationshipsParent) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ClassDataRelationshipsParent) SetType(v string)`

SetType sets Type field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


