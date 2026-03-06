# \PlaceClassesAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**PlacesSvcV1PlacesClassesGet**](PlaceClassesAPI.md#PlacesSvcV1PlacesClassesGet) | **Get** /places-svc/v1/places/classes/ | Get place classes list
[**PlacesSvcV1PlacesClassesPlaceClassIdDeprecatedDelete**](PlaceClassesAPI.md#PlacesSvcV1PlacesClassesPlaceClassIdDeprecatedDelete) | **Delete** /places-svc/v1/places/classes/{place_class_id}/deprecated | Undeprecate place class
[**PlacesSvcV1PlacesClassesPlaceClassIdDeprecatedPatch**](PlaceClassesAPI.md#PlacesSvcV1PlacesClassesPlaceClassIdDeprecatedPatch) | **Patch** /places-svc/v1/places/classes/{place_class_id}/deprecated | Deprecate place class
[**PlacesSvcV1PlacesClassesPlaceClassIdGet**](PlaceClassesAPI.md#PlacesSvcV1PlacesClassesPlaceClassIdGet) | **Get** /places-svc/v1/places/classes/{place_class_id} | Get place class
[**PlacesSvcV1PlacesClassesPlaceClassIdMediaDelete**](PlaceClassesAPI.md#PlacesSvcV1PlacesClassesPlaceClassIdMediaDelete) | **Delete** /places-svc/v1/places/classes/{place_class_id}/media | Delete place class uploaded media
[**PlacesSvcV1PlacesClassesPlaceClassIdMediaPost**](PlaceClassesAPI.md#PlacesSvcV1PlacesClassesPlaceClassIdMediaPost) | **Post** /places-svc/v1/places/classes/{place_class_id}/media | Create place class upload media link
[**PlacesSvcV1PlacesClassesPlaceClassIdPatch**](PlaceClassesAPI.md#PlacesSvcV1PlacesClassesPlaceClassIdPatch) | **Patch** /places-svc/v1/places/classes/{place_class_id} | Update place class
[**PlacesSvcV1PlacesClassesPost**](PlaceClassesAPI.md#PlacesSvcV1PlacesClassesPost) | **Post** /places-svc/v1/places/classes/ | Create place class



## PlacesSvcV1PlacesClassesGet

> PlaceClassesCollection PlacesSvcV1PlacesClassesGet(ctx).Size(size).Page(page).Text(text).ParentId(parentId).WithParents(withParents).WithChild(withChild).Deprecated(deprecated).Include(include).Execute()

Get place classes list



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
	size := int32(56) // int32 | Pagination limit. Cannot be greater than 100. (optional)
	page := int32(56) // int32 | Pagination offset. (optional)
	text := "text_example" // string | Full-text search query. Returns best-matching place classes. (optional)
	parentId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Filter place classes by parent class ID. (optional)
	withParents := true // bool | When filtering by parent_id, also include ancestor classes up the hierarchy. Requires `parent_id` to be set.  (optional)
	withChild := true // bool | When filtering by parent_id, also include descendant classes down the hierarchy. Requires `parent_id` to be set.  (optional)
	deprecated := true // bool | Filter place classes by deprecated flag. (optional)
	include := []string{"Include_example"} // []string | List of related resources to include in the response. Supported values: `parents`.  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PlaceClassesAPI.PlacesSvcV1PlacesClassesGet(context.Background()).Size(size).Page(page).Text(text).ParentId(parentId).WithParents(withParents).WithChild(withChild).Deprecated(deprecated).Include(include).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlaceClassesAPI.PlacesSvcV1PlacesClassesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PlacesSvcV1PlacesClassesGet`: PlaceClassesCollection
	fmt.Fprintf(os.Stdout, "Response from `PlaceClassesAPI.PlacesSvcV1PlacesClassesGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPlacesSvcV1PlacesClassesGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **size** | **int32** | Pagination limit. Cannot be greater than 100. | 
 **page** | **int32** | Pagination offset. | 
 **text** | **string** | Full-text search query. Returns best-matching place classes. | 
 **parentId** | **uuid.UUID** | Filter place classes by parent class ID. | 
 **withParents** | **bool** | When filtering by parent_id, also include ancestor classes up the hierarchy. Requires &#x60;parent_id&#x60; to be set.  | 
 **withChild** | **bool** | When filtering by parent_id, also include descendant classes down the hierarchy. Requires &#x60;parent_id&#x60; to be set.  | 
 **deprecated** | **bool** | Filter place classes by deprecated flag. | 
 **include** | **[]string** | List of related resources to include in the response. Supported values: &#x60;parents&#x60;.  | 

### Return type

[**PlaceClassesCollection**](PlaceClassesCollection.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PlacesSvcV1PlacesClassesPlaceClassIdDeprecatedDelete

> PlaceClass PlacesSvcV1PlacesClassesPlaceClassIdDeprecatedDelete(ctx, placeClassId).Execute()

Undeprecate place class



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
	placeClassId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Place class ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PlaceClassesAPI.PlacesSvcV1PlacesClassesPlaceClassIdDeprecatedDelete(context.Background(), placeClassId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlaceClassesAPI.PlacesSvcV1PlacesClassesPlaceClassIdDeprecatedDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PlacesSvcV1PlacesClassesPlaceClassIdDeprecatedDelete`: PlaceClass
	fmt.Fprintf(os.Stdout, "Response from `PlaceClassesAPI.PlacesSvcV1PlacesClassesPlaceClassIdDeprecatedDelete`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**placeClassId** | **uuid.UUID** | Place class ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiPlacesSvcV1PlacesClassesPlaceClassIdDeprecatedDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**PlaceClass**](PlaceClass.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PlacesSvcV1PlacesClassesPlaceClassIdDeprecatedPatch

> PlaceClass PlacesSvcV1PlacesClassesPlaceClassIdDeprecatedPatch(ctx, placeClassId).Execute()

Deprecate place class



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
	placeClassId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Place class ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PlaceClassesAPI.PlacesSvcV1PlacesClassesPlaceClassIdDeprecatedPatch(context.Background(), placeClassId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlaceClassesAPI.PlacesSvcV1PlacesClassesPlaceClassIdDeprecatedPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PlacesSvcV1PlacesClassesPlaceClassIdDeprecatedPatch`: PlaceClass
	fmt.Fprintf(os.Stdout, "Response from `PlaceClassesAPI.PlacesSvcV1PlacesClassesPlaceClassIdDeprecatedPatch`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**placeClassId** | **uuid.UUID** | Place class ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiPlacesSvcV1PlacesClassesPlaceClassIdDeprecatedPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**PlaceClass**](PlaceClass.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PlacesSvcV1PlacesClassesPlaceClassIdGet

> PlaceClass PlacesSvcV1PlacesClassesPlaceClassIdGet(ctx, placeClassId).Include(include).Execute()

Get place class



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
	placeClassId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Place class ID
	include := []string{"Include_example"} // []string | List of related resources to include in the response. Supported values: `parent`.  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PlaceClassesAPI.PlacesSvcV1PlacesClassesPlaceClassIdGet(context.Background(), placeClassId).Include(include).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlaceClassesAPI.PlacesSvcV1PlacesClassesPlaceClassIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PlacesSvcV1PlacesClassesPlaceClassIdGet`: PlaceClass
	fmt.Fprintf(os.Stdout, "Response from `PlaceClassesAPI.PlacesSvcV1PlacesClassesPlaceClassIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**placeClassId** | **uuid.UUID** | Place class ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiPlacesSvcV1PlacesClassesPlaceClassIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **include** | **[]string** | List of related resources to include in the response. Supported values: &#x60;parent&#x60;.  | 

### Return type

[**PlaceClass**](PlaceClass.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PlacesSvcV1PlacesClassesPlaceClassIdMediaDelete

> PlacesSvcV1PlacesClassesPlaceClassIdMediaDelete(ctx, placeClassId).DeleteUploadPlaceClassMedia(deleteUploadPlaceClassMedia).Execute()

Delete place class uploaded media



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
	placeClassId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Place class ID
	deleteUploadPlaceClassMedia := *openapiclient.NewDeleteUploadPlaceClassMedia(*openapiclient.NewDeleteUploadPlaceClassMediaData("TODO", "Type_example", *openapiclient.NewDeleteUploadPlaceClassMediaDataAttributes())) // DeleteUploadPlaceClassMedia | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.PlaceClassesAPI.PlacesSvcV1PlacesClassesPlaceClassIdMediaDelete(context.Background(), placeClassId).DeleteUploadPlaceClassMedia(deleteUploadPlaceClassMedia).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlaceClassesAPI.PlacesSvcV1PlacesClassesPlaceClassIdMediaDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**placeClassId** | **uuid.UUID** | Place class ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiPlacesSvcV1PlacesClassesPlaceClassIdMediaDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **deleteUploadPlaceClassMedia** | [**DeleteUploadPlaceClassMedia**](DeleteUploadPlaceClassMedia.md) |  | 

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


## PlacesSvcV1PlacesClassesPlaceClassIdMediaPost

> UploadPlaceClassMediaLinks PlacesSvcV1PlacesClassesPlaceClassIdMediaPost(ctx, placeClassId).Execute()

Create place class upload media link



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
	placeClassId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Place class ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PlaceClassesAPI.PlacesSvcV1PlacesClassesPlaceClassIdMediaPost(context.Background(), placeClassId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlaceClassesAPI.PlacesSvcV1PlacesClassesPlaceClassIdMediaPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PlacesSvcV1PlacesClassesPlaceClassIdMediaPost`: UploadPlaceClassMediaLinks
	fmt.Fprintf(os.Stdout, "Response from `PlaceClassesAPI.PlacesSvcV1PlacesClassesPlaceClassIdMediaPost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**placeClassId** | **uuid.UUID** | Place class ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiPlacesSvcV1PlacesClassesPlaceClassIdMediaPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**UploadPlaceClassMediaLinks**](UploadPlaceClassMediaLinks.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PlacesSvcV1PlacesClassesPlaceClassIdPatch

> PlaceClass PlacesSvcV1PlacesClassesPlaceClassIdPatch(ctx, placeClassId).UpdatePlaceClass(updatePlaceClass).Execute()

Update place class



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
	placeClassId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Place class ID
	updatePlaceClass := *openapiclient.NewUpdatePlaceClass(*openapiclient.NewUpdatePlaceClassData("TODO", "Type_example", *openapiclient.NewUpdatePlaceClassDataAttributes())) // UpdatePlaceClass | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PlaceClassesAPI.PlacesSvcV1PlacesClassesPlaceClassIdPatch(context.Background(), placeClassId).UpdatePlaceClass(updatePlaceClass).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlaceClassesAPI.PlacesSvcV1PlacesClassesPlaceClassIdPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PlacesSvcV1PlacesClassesPlaceClassIdPatch`: PlaceClass
	fmt.Fprintf(os.Stdout, "Response from `PlaceClassesAPI.PlacesSvcV1PlacesClassesPlaceClassIdPatch`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**placeClassId** | **uuid.UUID** | Place class ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiPlacesSvcV1PlacesClassesPlaceClassIdPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updatePlaceClass** | [**UpdatePlaceClass**](UpdatePlaceClass.md) |  | 

### Return type

[**PlaceClass**](PlaceClass.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PlacesSvcV1PlacesClassesPost

> PlaceClass PlacesSvcV1PlacesClassesPost(ctx).CreatePlaceClass(createPlaceClass).Execute()

Create place class



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
	createPlaceClass := *openapiclient.NewCreatePlaceClass(*openapiclient.NewCreatePlaceClassData("Type_example", *openapiclient.NewCreatePlaceClassDataAttributes("Name_example", "Description_example"))) // CreatePlaceClass | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PlaceClassesAPI.PlacesSvcV1PlacesClassesPost(context.Background()).CreatePlaceClass(createPlaceClass).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlaceClassesAPI.PlacesSvcV1PlacesClassesPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PlacesSvcV1PlacesClassesPost`: PlaceClass
	fmt.Fprintf(os.Stdout, "Response from `PlaceClassesAPI.PlacesSvcV1PlacesClassesPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPlacesSvcV1PlacesClassesPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **createPlaceClass** | [**CreatePlaceClass**](CreatePlaceClass.md) |  | 

### Return type

[**PlaceClass**](PlaceClass.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

