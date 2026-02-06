# UpdatePlaceDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Status** | **string** | The status of the place (e.g., active, inactive) | 
**Address** | **string** | The physical address of the place | 
**Name** | **string** | The name of the place | 
**Description** | Pointer to **string** | A brief description of the place | [optional] 
**Website** | Pointer to **string** | The website URL of the place | [optional] 
**Phone** | Pointer to **string** | The contact phone number of the place | [optional] 
**DeleteIcon** | Pointer to **bool** | Flag to indicate if the place&#39;s icon should be deleted | [optional] 
**DeleteBanner** | Pointer to **bool** | Flag to indicate if the place&#39;s banner should be deleted | [optional] 

## Methods

### NewUpdatePlaceDataAttributes

`func NewUpdatePlaceDataAttributes(status string, address string, name string, ) *UpdatePlaceDataAttributes`

NewUpdatePlaceDataAttributes instantiates a new UpdatePlaceDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdatePlaceDataAttributesWithDefaults

`func NewUpdatePlaceDataAttributesWithDefaults() *UpdatePlaceDataAttributes`

NewUpdatePlaceDataAttributesWithDefaults instantiates a new UpdatePlaceDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

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

### GetDeleteIcon

`func (o *UpdatePlaceDataAttributes) GetDeleteIcon() bool`

GetDeleteIcon returns the DeleteIcon field if non-nil, zero value otherwise.

### GetDeleteIconOk

`func (o *UpdatePlaceDataAttributes) GetDeleteIconOk() (*bool, bool)`

GetDeleteIconOk returns a tuple with the DeleteIcon field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeleteIcon

`func (o *UpdatePlaceDataAttributes) SetDeleteIcon(v bool)`

SetDeleteIcon sets DeleteIcon field to given value.

### HasDeleteIcon

`func (o *UpdatePlaceDataAttributes) HasDeleteIcon() bool`

HasDeleteIcon returns a boolean if a field has been set.

### GetDeleteBanner

`func (o *UpdatePlaceDataAttributes) GetDeleteBanner() bool`

GetDeleteBanner returns the DeleteBanner field if non-nil, zero value otherwise.

### GetDeleteBannerOk

`func (o *UpdatePlaceDataAttributes) GetDeleteBannerOk() (*bool, bool)`

GetDeleteBannerOk returns a tuple with the DeleteBanner field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeleteBanner

`func (o *UpdatePlaceDataAttributes) SetDeleteBanner(v bool)`

SetDeleteBanner sets DeleteBanner field to given value.

### HasDeleteBanner

`func (o *UpdatePlaceDataAttributes) HasDeleteBanner() bool`

HasDeleteBanner returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


