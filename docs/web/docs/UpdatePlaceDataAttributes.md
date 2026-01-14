# UpdatePlaceDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ClassId** | Pointer to [**uuid.UUID**](uuid.UUID.md) | The ID of the class this place belongs to | [optional] 
**Status** | Pointer to **string** | The status of the place (e.g., active, inactive) | [optional] 
**Address** | Pointer to **string** | The physical address of the place | [optional] 
**Name** | Pointer to **string** | The name of the place | [optional] 
**Description** | Pointer to **string** | A brief description of the place | [optional] 
**Icon** | Pointer to **string** | A URL to an icon representing the place | [optional] 
**Banner** | Pointer to **string** | A URL to a banner image for the place | [optional] 
**Website** | Pointer to **string** | The website URL of the place | [optional] 
**Phone** | Pointer to **string** | The contact phone number of the place | [optional] 

## Methods

### NewUpdatePlaceDataAttributes

`func NewUpdatePlaceDataAttributes() *UpdatePlaceDataAttributes`

NewUpdatePlaceDataAttributes instantiates a new UpdatePlaceDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdatePlaceDataAttributesWithDefaults

`func NewUpdatePlaceDataAttributesWithDefaults() *UpdatePlaceDataAttributes`

NewUpdatePlaceDataAttributesWithDefaults instantiates a new UpdatePlaceDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetClassId

`func (o *UpdatePlaceDataAttributes) GetClassId() uuid.UUID`

GetClassId returns the ClassId field if non-nil, zero value otherwise.

### GetClassIdOk

`func (o *UpdatePlaceDataAttributes) GetClassIdOk() (*uuid.UUID, bool)`

GetClassIdOk returns a tuple with the ClassId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClassId

`func (o *UpdatePlaceDataAttributes) SetClassId(v uuid.UUID)`

SetClassId sets ClassId field to given value.

### HasClassId

`func (o *UpdatePlaceDataAttributes) HasClassId() bool`

HasClassId returns a boolean if a field has been set.

### GetStatus

`func (o *UpdatePlaceDataAttributes) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *UpdatePlaceDataAttributes) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *UpdatePlaceDataAttributes) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *UpdatePlaceDataAttributes) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetAddress

`func (o *UpdatePlaceDataAttributes) GetAddress() string`

GetAddress returns the Address field if non-nil, zero value otherwise.

### GetAddressOk

`func (o *UpdatePlaceDataAttributes) GetAddressOk() (*string, bool)`

GetAddressOk returns a tuple with the Address field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddress

`func (o *UpdatePlaceDataAttributes) SetAddress(v string)`

SetAddress sets Address field to given value.

### HasAddress

`func (o *UpdatePlaceDataAttributes) HasAddress() bool`

HasAddress returns a boolean if a field has been set.

### GetName

`func (o *UpdatePlaceDataAttributes) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *UpdatePlaceDataAttributes) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *UpdatePlaceDataAttributes) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *UpdatePlaceDataAttributes) HasName() bool`

HasName returns a boolean if a field has been set.

### GetDescription

`func (o *UpdatePlaceDataAttributes) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *UpdatePlaceDataAttributes) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *UpdatePlaceDataAttributes) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *UpdatePlaceDataAttributes) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetIcon

`func (o *UpdatePlaceDataAttributes) GetIcon() string`

GetIcon returns the Icon field if non-nil, zero value otherwise.

### GetIconOk

`func (o *UpdatePlaceDataAttributes) GetIconOk() (*string, bool)`

GetIconOk returns a tuple with the Icon field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIcon

`func (o *UpdatePlaceDataAttributes) SetIcon(v string)`

SetIcon sets Icon field to given value.

### HasIcon

`func (o *UpdatePlaceDataAttributes) HasIcon() bool`

HasIcon returns a boolean if a field has been set.

### GetBanner

`func (o *UpdatePlaceDataAttributes) GetBanner() string`

GetBanner returns the Banner field if non-nil, zero value otherwise.

### GetBannerOk

`func (o *UpdatePlaceDataAttributes) GetBannerOk() (*string, bool)`

GetBannerOk returns a tuple with the Banner field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBanner

`func (o *UpdatePlaceDataAttributes) SetBanner(v string)`

SetBanner sets Banner field to given value.

### HasBanner

`func (o *UpdatePlaceDataAttributes) HasBanner() bool`

HasBanner returns a boolean if a field has been set.

### GetWebsite

`func (o *UpdatePlaceDataAttributes) GetWebsite() string`

GetWebsite returns the Website field if non-nil, zero value otherwise.

### GetWebsiteOk

`func (o *UpdatePlaceDataAttributes) GetWebsiteOk() (*string, bool)`

GetWebsiteOk returns a tuple with the Website field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWebsite

`func (o *UpdatePlaceDataAttributes) SetWebsite(v string)`

SetWebsite sets Website field to given value.

### HasWebsite

`func (o *UpdatePlaceDataAttributes) HasWebsite() bool`

HasWebsite returns a boolean if a field has been set.

### GetPhone

`func (o *UpdatePlaceDataAttributes) GetPhone() string`

GetPhone returns the Phone field if non-nil, zero value otherwise.

### GetPhoneOk

`func (o *UpdatePlaceDataAttributes) GetPhoneOk() (*string, bool)`

GetPhoneOk returns a tuple with the Phone field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPhone

`func (o *UpdatePlaceDataAttributes) SetPhone(v string)`

SetPhone sets Phone field to given value.

### HasPhone

`func (o *UpdatePlaceDataAttributes) HasPhone() bool`

HasPhone returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


