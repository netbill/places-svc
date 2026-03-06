# PlacesSvcV1PlacesPlaceIdVerifyPatchRequestData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | [**uuid.UUID**](uuid.UUID.md) | The unique identifier of the place | 
**Type** | **string** |  | 
**Attributes** | [**PlacesSvcV1PlacesPlaceIdVerifyPatchRequestDataAttributes**](PlacesSvcV1PlacesPlaceIdVerifyPatchRequestDataAttributes.md) |  | 

## Methods

### NewPlacesSvcV1PlacesPlaceIdVerifyPatchRequestData

`func NewPlacesSvcV1PlacesPlaceIdVerifyPatchRequestData(id uuid.UUID, type_ string, attributes PlacesSvcV1PlacesPlaceIdVerifyPatchRequestDataAttributes, ) *PlacesSvcV1PlacesPlaceIdVerifyPatchRequestData`

NewPlacesSvcV1PlacesPlaceIdVerifyPatchRequestData instantiates a new PlacesSvcV1PlacesPlaceIdVerifyPatchRequestData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPlacesSvcV1PlacesPlaceIdVerifyPatchRequestDataWithDefaults

`func NewPlacesSvcV1PlacesPlaceIdVerifyPatchRequestDataWithDefaults() *PlacesSvcV1PlacesPlaceIdVerifyPatchRequestData`

NewPlacesSvcV1PlacesPlaceIdVerifyPatchRequestDataWithDefaults instantiates a new PlacesSvcV1PlacesPlaceIdVerifyPatchRequestData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *PlacesSvcV1PlacesPlaceIdVerifyPatchRequestData) GetId() uuid.UUID`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *PlacesSvcV1PlacesPlaceIdVerifyPatchRequestData) GetIdOk() (*uuid.UUID, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *PlacesSvcV1PlacesPlaceIdVerifyPatchRequestData) SetId(v uuid.UUID)`

SetId sets Id field to given value.


### GetType

`func (o *PlacesSvcV1PlacesPlaceIdVerifyPatchRequestData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *PlacesSvcV1PlacesPlaceIdVerifyPatchRequestData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *PlacesSvcV1PlacesPlaceIdVerifyPatchRequestData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *PlacesSvcV1PlacesPlaceIdVerifyPatchRequestData) GetAttributes() PlacesSvcV1PlacesPlaceIdVerifyPatchRequestDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *PlacesSvcV1PlacesPlaceIdVerifyPatchRequestData) GetAttributesOk() (*PlacesSvcV1PlacesPlaceIdVerifyPatchRequestDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *PlacesSvcV1PlacesPlaceIdVerifyPatchRequestData) SetAttributes(v PlacesSvcV1PlacesPlaceIdVerifyPatchRequestDataAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


