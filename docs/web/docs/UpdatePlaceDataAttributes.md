# UpdatePlaceDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ClassId** | [**uuid.UUID**](uuid.UUID.md) | The class ID of the place (e.g., restaurant, park) | 
**Name** | **string** | The name of the place | 
**Address** | **string** | The address of the place | 
**Description** | Pointer to **string** | A brief description of the place | [optional] 
**IconKey** | Pointer to **string** | The S3 key for the place&#39;s icon image | [optional] 
**BannerKey** | Pointer to **string** | The S3 key for the place&#39;s banner image | [optional] 
**Website** | Pointer to **string** | The website URL of the place | [optional] 
**Phone** | Pointer to **string** | The contact phone number of the place | [optional] 

## Methods

### NewUpdatePlaceDataAttributes

`func NewUpdatePlaceDataAttributes(classId uuid.UUID, name string, address string, ) *UpdatePlaceDataAttributes`

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

### GetIconKey

`func (o *UpdatePlaceDataAttributes) GetIconKey() string`

GetIconKey returns the IconKey field if non-nil, zero value otherwise.

### GetIconKeyOk

`func (o *UpdatePlaceDataAttributes) GetIconKeyOk() (*string, bool)`

GetIconKeyOk returns a tuple with the IconKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIconKey

`func (o *UpdatePlaceDataAttributes) SetIconKey(v string)`

SetIconKey sets IconKey field to given value.

### HasIconKey

`func (o *UpdatePlaceDataAttributes) HasIconKey() bool`

HasIconKey returns a boolean if a field has been set.

### GetBannerKey

`func (o *UpdatePlaceDataAttributes) GetBannerKey() string`

GetBannerKey returns the BannerKey field if non-nil, zero value otherwise.

### GetBannerKeyOk

`func (o *UpdatePlaceDataAttributes) GetBannerKeyOk() (*string, bool)`

GetBannerKeyOk returns a tuple with the BannerKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBannerKey

`func (o *UpdatePlaceDataAttributes) SetBannerKey(v string)`

SetBannerKey sets BannerKey field to given value.

### HasBannerKey

`func (o *UpdatePlaceDataAttributes) HasBannerKey() bool`

HasBannerKey returns a boolean if a field has been set.

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


