# PlacesSvcV1PlacesGet400Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Errors** | [**[]PlacesSvcV1PlacesGet400ResponseErrorsInner**](PlacesSvcV1PlacesGet400ResponseErrorsInner.md) | Non empty array of errors occurred during request processing | 

## Methods

### NewPlacesSvcV1PlacesGet400Response

`func NewPlacesSvcV1PlacesGet400Response(errors []PlacesSvcV1PlacesGet400ResponseErrorsInner, ) *PlacesSvcV1PlacesGet400Response`

NewPlacesSvcV1PlacesGet400Response instantiates a new PlacesSvcV1PlacesGet400Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPlacesSvcV1PlacesGet400ResponseWithDefaults

`func NewPlacesSvcV1PlacesGet400ResponseWithDefaults() *PlacesSvcV1PlacesGet400Response`

NewPlacesSvcV1PlacesGet400ResponseWithDefaults instantiates a new PlacesSvcV1PlacesGet400Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetErrors

`func (o *PlacesSvcV1PlacesGet400Response) GetErrors() []PlacesSvcV1PlacesGet400ResponseErrorsInner`

GetErrors returns the Errors field if non-nil, zero value otherwise.

### GetErrorsOk

`func (o *PlacesSvcV1PlacesGet400Response) GetErrorsOk() (*[]PlacesSvcV1PlacesGet400ResponseErrorsInner, bool)`

GetErrorsOk returns a tuple with the Errors field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrors

`func (o *PlacesSvcV1PlacesGet400Response) SetErrors(v []PlacesSvcV1PlacesGet400ResponseErrorsInner)`

SetErrors sets Errors field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


