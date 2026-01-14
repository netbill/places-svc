# PlaceDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Status** | **string** | The status of the place (e.g., active, inactive) | 
**Verified** | **bool** | Indicates if the place has been verified | 
**Point** | [**Point**](Point.md) |  | 
**Address** | **string** | The physical address of the place | 
**Name** | **string** | The name of the place | 
**Description** | Pointer to **string** | A brief description of the place | [optional] 
**Icon** | Pointer to **string** | A URL to an icon representing the place | [optional] 
**Banner** | Pointer to **string** | A URL to a banner image for the place | [optional] 
**Website** | Pointer to **string** | The website URL of the place | [optional] 
**Phone** | Pointer to **string** | The contact phone number of the place | [optional] 
**CreatedAt** | **time.Time** | The date and time when the place was created | 
**UpdatedAt** | **time.Time** | The date and time when the place was last updated | 

## Methods

### NewPlaceDataAttributes

`func NewPlaceDataAttributes(status string, verified bool, point Point, address string, name string, createdAt time.Time, updatedAt time.Time, ) *PlaceDataAttributes`

NewPlaceDataAttributes instantiates a new PlaceDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPlaceDataAttributesWithDefaults

`func NewPlaceDataAttributesWithDefaults() *PlaceDataAttributes`

NewPlaceDataAttributesWithDefaults instantiates a new PlaceDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStatus

`func (o *PlaceDataAttributes) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *PlaceDataAttributes) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *PlaceDataAttributes) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetVerified

`func (o *PlaceDataAttributes) GetVerified() bool`

GetVerified returns the Verified field if non-nil, zero value otherwise.

### GetVerifiedOk

`func (o *PlaceDataAttributes) GetVerifiedOk() (*bool, bool)`

GetVerifiedOk returns a tuple with the Verified field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVerified

`func (o *PlaceDataAttributes) SetVerified(v bool)`

SetVerified sets Verified field to given value.


### GetPoint

`func (o *PlaceDataAttributes) GetPoint() Point`

GetPoint returns the Point field if non-nil, zero value otherwise.

### GetPointOk

`func (o *PlaceDataAttributes) GetPointOk() (*Point, bool)`

GetPointOk returns a tuple with the Point field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoint

`func (o *PlaceDataAttributes) SetPoint(v Point)`

SetPoint sets Point field to given value.


### GetAddress

`func (o *PlaceDataAttributes) GetAddress() string`

GetAddress returns the Address field if non-nil, zero value otherwise.

### GetAddressOk

`func (o *PlaceDataAttributes) GetAddressOk() (*string, bool)`

GetAddressOk returns a tuple with the Address field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddress

`func (o *PlaceDataAttributes) SetAddress(v string)`

SetAddress sets Address field to given value.


### GetName

`func (o *PlaceDataAttributes) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *PlaceDataAttributes) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *PlaceDataAttributes) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *PlaceDataAttributes) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *PlaceDataAttributes) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *PlaceDataAttributes) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *PlaceDataAttributes) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetIcon

`func (o *PlaceDataAttributes) GetIcon() string`

GetIcon returns the Icon field if non-nil, zero value otherwise.

### GetIconOk

`func (o *PlaceDataAttributes) GetIconOk() (*string, bool)`

GetIconOk returns a tuple with the Icon field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIcon

`func (o *PlaceDataAttributes) SetIcon(v string)`

SetIcon sets Icon field to given value.

### HasIcon

`func (o *PlaceDataAttributes) HasIcon() bool`

HasIcon returns a boolean if a field has been set.

### GetBanner

`func (o *PlaceDataAttributes) GetBanner() string`

GetBanner returns the Banner field if non-nil, zero value otherwise.

### GetBannerOk

`func (o *PlaceDataAttributes) GetBannerOk() (*string, bool)`

GetBannerOk returns a tuple with the Banner field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBanner

`func (o *PlaceDataAttributes) SetBanner(v string)`

SetBanner sets Banner field to given value.

### HasBanner

`func (o *PlaceDataAttributes) HasBanner() bool`

HasBanner returns a boolean if a field has been set.

### GetWebsite

`func (o *PlaceDataAttributes) GetWebsite() string`

GetWebsite returns the Website field if non-nil, zero value otherwise.

### GetWebsiteOk

`func (o *PlaceDataAttributes) GetWebsiteOk() (*string, bool)`

GetWebsiteOk returns a tuple with the Website field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWebsite

`func (o *PlaceDataAttributes) SetWebsite(v string)`

SetWebsite sets Website field to given value.

### HasWebsite

`func (o *PlaceDataAttributes) HasWebsite() bool`

HasWebsite returns a boolean if a field has been set.

### GetPhone

`func (o *PlaceDataAttributes) GetPhone() string`

GetPhone returns the Phone field if non-nil, zero value otherwise.

### GetPhoneOk

`func (o *PlaceDataAttributes) GetPhoneOk() (*string, bool)`

GetPhoneOk returns a tuple with the Phone field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPhone

`func (o *PlaceDataAttributes) SetPhone(v string)`

SetPhone sets Phone field to given value.

### HasPhone

`func (o *PlaceDataAttributes) HasPhone() bool`

HasPhone returns a boolean if a field has been set.

### GetCreatedAt

`func (o *PlaceDataAttributes) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *PlaceDataAttributes) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *PlaceDataAttributes) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetUpdatedAt

`func (o *PlaceDataAttributes) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *PlaceDataAttributes) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *PlaceDataAttributes) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


