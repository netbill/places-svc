# ConfirmUpdatePlaceDataAttributes

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

### NewConfirmUpdatePlaceDataAttributes

`func NewConfirmUpdatePlaceDataAttributes(status string, address string, name string, ) *ConfirmUpdatePlaceDataAttributes`

NewConfirmUpdatePlaceDataAttributes instantiates a new ConfirmUpdatePlaceDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfirmUpdatePlaceDataAttributesWithDefaults

`func NewConfirmUpdatePlaceDataAttributesWithDefaults() *ConfirmUpdatePlaceDataAttributes`

NewConfirmUpdatePlaceDataAttributesWithDefaults instantiates a new ConfirmUpdatePlaceDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStatus

`func (o *ConfirmUpdatePlaceDataAttributes) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ConfirmUpdatePlaceDataAttributes) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ConfirmUpdatePlaceDataAttributes) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetAddress

`func (o *ConfirmUpdatePlaceDataAttributes) GetAddress() string`

GetAddress returns the Address field if non-nil, zero value otherwise.

### GetAddressOk

`func (o *ConfirmUpdatePlaceDataAttributes) GetAddressOk() (*string, bool)`

GetAddressOk returns a tuple with the Address field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddress

`func (o *ConfirmUpdatePlaceDataAttributes) SetAddress(v string)`

SetAddress sets Address field to given value.


### GetName

`func (o *ConfirmUpdatePlaceDataAttributes) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfirmUpdatePlaceDataAttributes) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfirmUpdatePlaceDataAttributes) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *ConfirmUpdatePlaceDataAttributes) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *ConfirmUpdatePlaceDataAttributes) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *ConfirmUpdatePlaceDataAttributes) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *ConfirmUpdatePlaceDataAttributes) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetWebsite

`func (o *ConfirmUpdatePlaceDataAttributes) GetWebsite() string`

GetWebsite returns the Website field if non-nil, zero value otherwise.

### GetWebsiteOk

`func (o *ConfirmUpdatePlaceDataAttributes) GetWebsiteOk() (*string, bool)`

GetWebsiteOk returns a tuple with the Website field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWebsite

`func (o *ConfirmUpdatePlaceDataAttributes) SetWebsite(v string)`

SetWebsite sets Website field to given value.

### HasWebsite

`func (o *ConfirmUpdatePlaceDataAttributes) HasWebsite() bool`

HasWebsite returns a boolean if a field has been set.

### GetPhone

`func (o *ConfirmUpdatePlaceDataAttributes) GetPhone() string`

GetPhone returns the Phone field if non-nil, zero value otherwise.

### GetPhoneOk

`func (o *ConfirmUpdatePlaceDataAttributes) GetPhoneOk() (*string, bool)`

GetPhoneOk returns a tuple with the Phone field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPhone

`func (o *ConfirmUpdatePlaceDataAttributes) SetPhone(v string)`

SetPhone sets Phone field to given value.

### HasPhone

`func (o *ConfirmUpdatePlaceDataAttributes) HasPhone() bool`

HasPhone returns a boolean if a field has been set.

### GetDeleteIcon

`func (o *ConfirmUpdatePlaceDataAttributes) GetDeleteIcon() bool`

GetDeleteIcon returns the DeleteIcon field if non-nil, zero value otherwise.

### GetDeleteIconOk

`func (o *ConfirmUpdatePlaceDataAttributes) GetDeleteIconOk() (*bool, bool)`

GetDeleteIconOk returns a tuple with the DeleteIcon field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeleteIcon

`func (o *ConfirmUpdatePlaceDataAttributes) SetDeleteIcon(v bool)`

SetDeleteIcon sets DeleteIcon field to given value.

### HasDeleteIcon

`func (o *ConfirmUpdatePlaceDataAttributes) HasDeleteIcon() bool`

HasDeleteIcon returns a boolean if a field has been set.

### GetDeleteBanner

`func (o *ConfirmUpdatePlaceDataAttributes) GetDeleteBanner() bool`

GetDeleteBanner returns the DeleteBanner field if non-nil, zero value otherwise.

### GetDeleteBannerOk

`func (o *ConfirmUpdatePlaceDataAttributes) GetDeleteBannerOk() (*bool, bool)`

GetDeleteBannerOk returns a tuple with the DeleteBanner field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeleteBanner

`func (o *ConfirmUpdatePlaceDataAttributes) SetDeleteBanner(v bool)`

SetDeleteBanner sets DeleteBanner field to given value.

### HasDeleteBanner

`func (o *ConfirmUpdatePlaceDataAttributes) HasDeleteBanner() bool`

HasDeleteBanner returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


