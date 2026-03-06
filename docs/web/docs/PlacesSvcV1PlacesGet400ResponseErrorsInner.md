# PlacesSvcV1PlacesGet400ResponseErrorsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Title** | **string** | Title is a short, human-readable summary of the problem | 
**Detail** | Pointer to **string** | Detail is a human-readable explanation specific to this occurrence of the problem | [optional] 
**Status** | **int32** | Status is the HTTP status code applicable to this problem | 

## Methods

### NewPlacesSvcV1PlacesGet400ResponseErrorsInner

`func NewPlacesSvcV1PlacesGet400ResponseErrorsInner(title string, status int32, ) *PlacesSvcV1PlacesGet400ResponseErrorsInner`

NewPlacesSvcV1PlacesGet400ResponseErrorsInner instantiates a new PlacesSvcV1PlacesGet400ResponseErrorsInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPlacesSvcV1PlacesGet400ResponseErrorsInnerWithDefaults

`func NewPlacesSvcV1PlacesGet400ResponseErrorsInnerWithDefaults() *PlacesSvcV1PlacesGet400ResponseErrorsInner`

NewPlacesSvcV1PlacesGet400ResponseErrorsInnerWithDefaults instantiates a new PlacesSvcV1PlacesGet400ResponseErrorsInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTitle

`func (o *PlacesSvcV1PlacesGet400ResponseErrorsInner) GetTitle() string`

GetTitle returns the Title field if non-nil, zero value otherwise.

### GetTitleOk

`func (o *PlacesSvcV1PlacesGet400ResponseErrorsInner) GetTitleOk() (*string, bool)`

GetTitleOk returns a tuple with the Title field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTitle

`func (o *PlacesSvcV1PlacesGet400ResponseErrorsInner) SetTitle(v string)`

SetTitle sets Title field to given value.


### GetDetail

`func (o *PlacesSvcV1PlacesGet400ResponseErrorsInner) GetDetail() string`

GetDetail returns the Detail field if non-nil, zero value otherwise.

### GetDetailOk

`func (o *PlacesSvcV1PlacesGet400ResponseErrorsInner) GetDetailOk() (*string, bool)`

GetDetailOk returns a tuple with the Detail field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDetail

`func (o *PlacesSvcV1PlacesGet400ResponseErrorsInner) SetDetail(v string)`

SetDetail sets Detail field to given value.

### HasDetail

`func (o *PlacesSvcV1PlacesGet400ResponseErrorsInner) HasDetail() bool`

HasDetail returns a boolean if a field has been set.

### GetStatus

`func (o *PlacesSvcV1PlacesGet400ResponseErrorsInner) GetStatus() int32`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *PlacesSvcV1PlacesGet400ResponseErrorsInner) GetStatusOk() (*int32, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *PlacesSvcV1PlacesGet400ResponseErrorsInner) SetStatus(v int32)`

SetStatus sets Status field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


