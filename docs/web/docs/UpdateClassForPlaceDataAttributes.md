# UpdateClassForPlaceDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ClassId** | [**uuid.UUID**](uuid.UUID.md) | The ID of the class this place belongs to | 

## Methods

### NewUpdateClassForPlaceDataAttributes

`func NewUpdateClassForPlaceDataAttributes(classId uuid.UUID, ) *UpdateClassForPlaceDataAttributes`

NewUpdateClassForPlaceDataAttributes instantiates a new UpdateClassForPlaceDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateClassForPlaceDataAttributesWithDefaults

`func NewUpdateClassForPlaceDataAttributesWithDefaults() *UpdateClassForPlaceDataAttributes`

NewUpdateClassForPlaceDataAttributesWithDefaults instantiates a new UpdateClassForPlaceDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetClassId

`func (o *UpdateClassForPlaceDataAttributes) GetClassId() uuid.UUID`

GetClassId returns the ClassId field if non-nil, zero value otherwise.

### GetClassIdOk

`func (o *UpdateClassForPlaceDataAttributes) GetClassIdOk() (*uuid.UUID, bool)`

GetClassIdOk returns a tuple with the ClassId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClassId

`func (o *UpdateClassForPlaceDataAttributes) SetClassId(v uuid.UUID)`

SetClassId sets ClassId field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


