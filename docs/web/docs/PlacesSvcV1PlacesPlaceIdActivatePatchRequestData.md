# PlacesSvcV1PlacesPlaceIdActivatePatchRequestData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | [**uuid.UUID**](uuid.UUID.md) | The unique identifier of the place | 
**Type** | **string** |  | 
**Attributes** | [**PlacesSvcV1PlacesPlaceIdActivatePatchRequestDataAttributes**](PlacesSvcV1PlacesPlaceIdActivatePatchRequestDataAttributes.md) |  | 

## Methods

### NewPlacesSvcV1PlacesPlaceIdActivatePatchRequestData

`func NewPlacesSvcV1PlacesPlaceIdActivatePatchRequestData(id uuid.UUID, type_ string, attributes PlacesSvcV1PlacesPlaceIdActivatePatchRequestDataAttributes, ) *PlacesSvcV1PlacesPlaceIdActivatePatchRequestData`

NewPlacesSvcV1PlacesPlaceIdActivatePatchRequestData instantiates a new PlacesSvcV1PlacesPlaceIdActivatePatchRequestData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPlacesSvcV1PlacesPlaceIdActivatePatchRequestDataWithDefaults

`func NewPlacesSvcV1PlacesPlaceIdActivatePatchRequestDataWithDefaults() *PlacesSvcV1PlacesPlaceIdActivatePatchRequestData`

NewPlacesSvcV1PlacesPlaceIdActivatePatchRequestDataWithDefaults instantiates a new PlacesSvcV1PlacesPlaceIdActivatePatchRequestData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *PlacesSvcV1PlacesPlaceIdActivatePatchRequestData) GetId() uuid.UUID`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *PlacesSvcV1PlacesPlaceIdActivatePatchRequestData) GetIdOk() (*uuid.UUID, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *PlacesSvcV1PlacesPlaceIdActivatePatchRequestData) SetId(v uuid.UUID)`

SetId sets Id field to given value.


### GetType

`func (o *PlacesSvcV1PlacesPlaceIdActivatePatchRequestData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *PlacesSvcV1PlacesPlaceIdActivatePatchRequestData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *PlacesSvcV1PlacesPlaceIdActivatePatchRequestData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *PlacesSvcV1PlacesPlaceIdActivatePatchRequestData) GetAttributes() PlacesSvcV1PlacesPlaceIdActivatePatchRequestDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *PlacesSvcV1PlacesPlaceIdActivatePatchRequestData) GetAttributesOk() (*PlacesSvcV1PlacesPlaceIdActivatePatchRequestDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *PlacesSvcV1PlacesPlaceIdActivatePatchRequestData) SetAttributes(v PlacesSvcV1PlacesPlaceIdActivatePatchRequestDataAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


