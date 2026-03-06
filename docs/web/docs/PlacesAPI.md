# \PlacesAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**PlacesSvcV1PlacesGet**](PlacesAPI.md#PlacesSvcV1PlacesGet) | **Get** /places-svc/v1/places/ | Get places list
[**PlacesSvcV1PlacesPlaceIdActivatePatch**](PlacesAPI.md#PlacesSvcV1PlacesPlaceIdActivatePatch) | **Patch** /places-svc/v1/places/{place_id}/activate | Activate place
[**PlacesSvcV1PlacesPlaceIdDeactivatePatch**](PlacesAPI.md#PlacesSvcV1PlacesPlaceIdDeactivatePatch) | **Patch** /places-svc/v1/places/{place_id}/deactivate | Deactivate place
[**PlacesSvcV1PlacesPlaceIdDelete**](PlacesAPI.md#PlacesSvcV1PlacesPlaceIdDelete) | **Delete** /places-svc/v1/places/{place_id} | Delete place
[**PlacesSvcV1PlacesPlaceIdGet**](PlacesAPI.md#PlacesSvcV1PlacesPlaceIdGet) | **Get** /places-svc/v1/places/{place_id} | Get place
[**PlacesSvcV1PlacesPlaceIdMediaDelete**](PlacesAPI.md#PlacesSvcV1PlacesPlaceIdMediaDelete) | **Delete** /places-svc/v1/places/{place_id}/media | Delete place uploaded media
[**PlacesSvcV1PlacesPlaceIdMediaPost**](PlacesAPI.md#PlacesSvcV1PlacesPlaceIdMediaPost) | **Post** /places-svc/v1/places/{place_id}/media | Create place upload media link
[**PlacesSvcV1PlacesPlaceIdPatch**](PlacesAPI.md#PlacesSvcV1PlacesPlaceIdPatch) | **Patch** /places-svc/v1/places/{place_id} | Update place
[**PlacesSvcV1PlacesPlaceIdVerifyDelete**](PlacesAPI.md#PlacesSvcV1PlacesPlaceIdVerifyDelete) | **Delete** /places-svc/v1/places/{place_id}/verify | Unverify place
[**PlacesSvcV1PlacesPlaceIdVerifyPatch**](PlacesAPI.md#PlacesSvcV1PlacesPlaceIdVerifyPatch) | **Patch** /places-svc/v1/places/{place_id}/verify | Verify place
[**PlacesSvcV1PlacesPost**](PlacesAPI.md#PlacesSvcV1PlacesPost) | **Post** /places-svc/v1/places/ | Create place



## PlacesSvcV1PlacesGet

> PlacesCollection PlacesSvcV1PlacesGet(ctx).Size(size).Page(page).OrganizationId(organizationId).Statuses(statuses).OrgStatus(orgStatus).Verified(verified).Text(text).ClassIds(classIds).IncludeParent(includeParent).IncludeChildren(includeChildren).Lon(lon).Lat(lat).Radius(radius).Include(include).Execute()

Get places list



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	size := int32(56) // int32 | Pagination limit. Cannot be greater than 1000. (optional)
	page := int32(56) // int32 | Pagination offset. (optional)
	organizationId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Filter places by organization ID. (optional)
	statuses := []string{"Statuses_example"} // []string | Filter places by place status. (optional)
	orgStatus := []string{"Inner_example"} // []string | Filter places by organization status. (optional)
	verified := true // bool | Filter places by verified flag. (optional)
	text := "text_example" // string | Full-text search query. Returns best-matching places. (optional)
	classIds := []uuid.UUID{"TODO"} // []uuid.UUID | Filter places by place class IDs. (optional)
	includeParent := true // bool | When filtering by class_ids, also include places whose class is a parent of one of the specified classes.  (optional)
	includeChildren := true // bool | When filtering by class_ids, also include places whose class is a child of one of the specified classes.  (optional)
	lon := float64(1.2) // float64 | Longitude for geo-proximity filter. Must be provided together with `lat` and `radius`.  (optional)
	lat := float64(1.2) // float64 | Latitude for geo-proximity filter. Must be provided together with `lon` and `radius`.  (optional)
	radius := int32(56) // int32 | Radius in meters for geo-proximity filter. Must be provided together with `lon` and `lat`.  (optional)
	include := []string{"Include_example"} // []string | Comma-separated list of related resources to include in the response. Supported values: `place_classes`, `organizations`.  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PlacesAPI.PlacesSvcV1PlacesGet(context.Background()).Size(size).Page(page).OrganizationId(organizationId).Statuses(statuses).OrgStatus(orgStatus).Verified(verified).Text(text).ClassIds(classIds).IncludeParent(includeParent).IncludeChildren(includeChildren).Lon(lon).Lat(lat).Radius(radius).Include(include).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlacesAPI.PlacesSvcV1PlacesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PlacesSvcV1PlacesGet`: PlacesCollection
	fmt.Fprintf(os.Stdout, "Response from `PlacesAPI.PlacesSvcV1PlacesGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPlacesSvcV1PlacesGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **size** | **int32** | Pagination limit. Cannot be greater than 1000. | 
 **page** | **int32** | Pagination offset. | 
 **organizationId** | **uuid.UUID** | Filter places by organization ID. | 
 **statuses** | **[]string** | Filter places by place status. | 
 **orgStatus** | **[]string** | Filter places by organization status. | 
 **verified** | **bool** | Filter places by verified flag. | 
 **text** | **string** | Full-text search query. Returns best-matching places. | 
 **classIds** | [**[]uuid.UUID**](uuid.UUID.md) | Filter places by place class IDs. | 
 **includeParent** | **bool** | When filtering by class_ids, also include places whose class is a parent of one of the specified classes.  | 
 **includeChildren** | **bool** | When filtering by class_ids, also include places whose class is a child of one of the specified classes.  | 
 **lon** | **float64** | Longitude for geo-proximity filter. Must be provided together with &#x60;lat&#x60; and &#x60;radius&#x60;.  | 
 **lat** | **float64** | Latitude for geo-proximity filter. Must be provided together with &#x60;lon&#x60; and &#x60;radius&#x60;.  | 
 **radius** | **int32** | Radius in meters for geo-proximity filter. Must be provided together with &#x60;lon&#x60; and &#x60;lat&#x60;.  | 
 **include** | **[]string** | Comma-separated list of related resources to include in the response. Supported values: &#x60;place_classes&#x60;, &#x60;organizations&#x60;.  | 

### Return type

[**PlacesCollection**](PlacesCollection.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PlacesSvcV1PlacesPlaceIdActivatePatch

> Place PlacesSvcV1PlacesPlaceIdActivatePatch(ctx, placeId).Execute()

Activate place



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	placeId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Place ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PlacesAPI.PlacesSvcV1PlacesPlaceIdActivatePatch(context.Background(), placeId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlacesAPI.PlacesSvcV1PlacesPlaceIdActivatePatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PlacesSvcV1PlacesPlaceIdActivatePatch`: Place
	fmt.Fprintf(os.Stdout, "Response from `PlacesAPI.PlacesSvcV1PlacesPlaceIdActivatePatch`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**placeId** | **uuid.UUID** | Place ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiPlacesSvcV1PlacesPlaceIdActivatePatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Place**](Place.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PlacesSvcV1PlacesPlaceIdDeactivatePatch

> Place PlacesSvcV1PlacesPlaceIdDeactivatePatch(ctx, placeId).Execute()

Deactivate place



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	placeId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Place ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PlacesAPI.PlacesSvcV1PlacesPlaceIdDeactivatePatch(context.Background(), placeId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlacesAPI.PlacesSvcV1PlacesPlaceIdDeactivatePatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PlacesSvcV1PlacesPlaceIdDeactivatePatch`: Place
	fmt.Fprintf(os.Stdout, "Response from `PlacesAPI.PlacesSvcV1PlacesPlaceIdDeactivatePatch`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**placeId** | **uuid.UUID** | Place ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiPlacesSvcV1PlacesPlaceIdDeactivatePatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Place**](Place.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PlacesSvcV1PlacesPlaceIdDelete

> PlacesSvcV1PlacesPlaceIdDelete(ctx, placeId).Execute()

Delete place



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	placeId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Place ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.PlacesAPI.PlacesSvcV1PlacesPlaceIdDelete(context.Background(), placeId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlacesAPI.PlacesSvcV1PlacesPlaceIdDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**placeId** | **uuid.UUID** | Place ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiPlacesSvcV1PlacesPlaceIdDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PlacesSvcV1PlacesPlaceIdGet

> Place PlacesSvcV1PlacesPlaceIdGet(ctx, placeId).Include(include).Execute()

Get place



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	placeId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Place ID
	include := []string{"Include_example"} // []string | List of related resources to include in the response. Supported values: `place_class`, `organization`.  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PlacesAPI.PlacesSvcV1PlacesPlaceIdGet(context.Background(), placeId).Include(include).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlacesAPI.PlacesSvcV1PlacesPlaceIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PlacesSvcV1PlacesPlaceIdGet`: Place
	fmt.Fprintf(os.Stdout, "Response from `PlacesAPI.PlacesSvcV1PlacesPlaceIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**placeId** | **uuid.UUID** | Place ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiPlacesSvcV1PlacesPlaceIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **include** | **[]string** | List of related resources to include in the response. Supported values: &#x60;place_class&#x60;, &#x60;organization&#x60;.  | 

### Return type

[**Place**](Place.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PlacesSvcV1PlacesPlaceIdMediaDelete

> PlacesSvcV1PlacesPlaceIdMediaDelete(ctx, placeId).DeleteUploadPlaceMedia(deleteUploadPlaceMedia).Execute()

Delete place uploaded media



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	placeId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Place ID
	deleteUploadPlaceMedia := *openapiclient.NewDeleteUploadPlaceMedia(*openapiclient.NewDeleteUploadPlaceMediaData("TODO", "Type_example", *openapiclient.NewDeleteUploadPlaceMediaDataAttributes())) // DeleteUploadPlaceMedia | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.PlacesAPI.PlacesSvcV1PlacesPlaceIdMediaDelete(context.Background(), placeId).DeleteUploadPlaceMedia(deleteUploadPlaceMedia).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlacesAPI.PlacesSvcV1PlacesPlaceIdMediaDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**placeId** | **uuid.UUID** | Place ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiPlacesSvcV1PlacesPlaceIdMediaDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **deleteUploadPlaceMedia** | [**DeleteUploadPlaceMedia**](DeleteUploadPlaceMedia.md) |  | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PlacesSvcV1PlacesPlaceIdMediaPost

> UploadPlaceMediaLinks PlacesSvcV1PlacesPlaceIdMediaPost(ctx, placeId).Execute()

Create place upload media link



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	placeId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Place ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PlacesAPI.PlacesSvcV1PlacesPlaceIdMediaPost(context.Background(), placeId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlacesAPI.PlacesSvcV1PlacesPlaceIdMediaPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PlacesSvcV1PlacesPlaceIdMediaPost`: UploadPlaceMediaLinks
	fmt.Fprintf(os.Stdout, "Response from `PlacesAPI.PlacesSvcV1PlacesPlaceIdMediaPost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**placeId** | **uuid.UUID** | Place ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiPlacesSvcV1PlacesPlaceIdMediaPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**UploadPlaceMediaLinks**](UploadPlaceMediaLinks.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PlacesSvcV1PlacesPlaceIdPatch

> Place PlacesSvcV1PlacesPlaceIdPatch(ctx, placeId).UpdatePlace(updatePlace).Execute()

Update place



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	placeId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Place ID
	updatePlace := *openapiclient.NewUpdatePlace(*openapiclient.NewUpdatePlaceData("TODO", "Type_example", *openapiclient.NewUpdatePlaceDataAttributes())) // UpdatePlace | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PlacesAPI.PlacesSvcV1PlacesPlaceIdPatch(context.Background(), placeId).UpdatePlace(updatePlace).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlacesAPI.PlacesSvcV1PlacesPlaceIdPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PlacesSvcV1PlacesPlaceIdPatch`: Place
	fmt.Fprintf(os.Stdout, "Response from `PlacesAPI.PlacesSvcV1PlacesPlaceIdPatch`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**placeId** | **uuid.UUID** | Place ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiPlacesSvcV1PlacesPlaceIdPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updatePlace** | [**UpdatePlace**](UpdatePlace.md) |  | 

### Return type

[**Place**](Place.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PlacesSvcV1PlacesPlaceIdVerifyDelete

> Place PlacesSvcV1PlacesPlaceIdVerifyDelete(ctx, placeId).Execute()

Unverify place



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	placeId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Place ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PlacesAPI.PlacesSvcV1PlacesPlaceIdVerifyDelete(context.Background(), placeId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlacesAPI.PlacesSvcV1PlacesPlaceIdVerifyDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PlacesSvcV1PlacesPlaceIdVerifyDelete`: Place
	fmt.Fprintf(os.Stdout, "Response from `PlacesAPI.PlacesSvcV1PlacesPlaceIdVerifyDelete`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**placeId** | **uuid.UUID** | Place ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiPlacesSvcV1PlacesPlaceIdVerifyDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Place**](Place.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PlacesSvcV1PlacesPlaceIdVerifyPatch

> Place PlacesSvcV1PlacesPlaceIdVerifyPatch(ctx, placeId).Execute()

Verify place



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	placeId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Place ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PlacesAPI.PlacesSvcV1PlacesPlaceIdVerifyPatch(context.Background(), placeId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlacesAPI.PlacesSvcV1PlacesPlaceIdVerifyPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PlacesSvcV1PlacesPlaceIdVerifyPatch`: Place
	fmt.Fprintf(os.Stdout, "Response from `PlacesAPI.PlacesSvcV1PlacesPlaceIdVerifyPatch`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**placeId** | **uuid.UUID** | Place ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiPlacesSvcV1PlacesPlaceIdVerifyPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Place**](Place.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PlacesSvcV1PlacesPost

> Place PlacesSvcV1PlacesPost(ctx).CreatePlace(createPlace).Execute()

Create place



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	createPlace := *openapiclient.NewCreatePlace(*openapiclient.NewCreatePlaceData("Type_example", *openapiclient.NewCreatePlaceDataAttributes("TODO", "TODO", *openapiclient.NewPoint(float64(123), float64(123)), "Address_example", "Name_example"))) // CreatePlace | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PlacesAPI.PlacesSvcV1PlacesPost(context.Background()).CreatePlace(createPlace).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlacesAPI.PlacesSvcV1PlacesPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PlacesSvcV1PlacesPost`: Place
	fmt.Fprintf(os.Stdout, "Response from `PlacesAPI.PlacesSvcV1PlacesPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPlacesSvcV1PlacesPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **createPlace** | [**CreatePlace**](CreatePlace.md) |  | 

### Return type

[**Place**](Place.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

