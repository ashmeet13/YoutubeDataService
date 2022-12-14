# Youtube Data Service

A demo Go service that -
1. Collects the latest data for Youtube Videos from Youtube API. The defination of latest used here is videos that are uploaded after the service started.
2. Exposes an API to fiddle around with the data being collected in the background.
   1. Search API to fetch data on requested search parameters.
   2. Fetch API that returns data in reverse chronological order.


## Breakdown of the service

Broadly, the service has two main components - the `Worker` and the `Server`.

`Worker` is our async background job that every few seconds requests new data from the Youtube API and stores into a MongoDB instance.

`Server` exposes two APIs to fetch this data. 

- `GET /fetch` - This will return a UniqueID `userid` back that can be used to fetch the data in pages. You can set custom `userid` and `pagesize` by setting them in url parameters.

Example - `/fetch?userid=ashmeet&pagesize=3`

- `GET /fetch/<userid>/<pagenumber>` - This will return the data for the `userid` for the `pagenumber` with the `pagesize` that was mentioned in `GET /fetch`. If no pagesize was mentioned the default is 5

- `POST /search` - This takes in a JSON body with two components `Title` and `Description` to search the database against. Both the parameters cannot be empty at the same time. Search scans through the database to find documents with similar `Title` and `Description` and returns a list of the videos found. This endpoint uses MongoDB text indexes to perform text searches over the documents.

```json
{
    "Title" : "Title To Search",
    "Description" : "Description To Search",
}
```

## Why do I require a User?

Taking an analogy to a Facebook feed that shows us events in reverse chronological order i.e. the most
recent posts published at the top - the `GET /fetch/<userid>/<pagenumber>` API tries to replicate that.

Now, since we have a background worker that is going to keep adding videos into our database as and when it recieves them from the API, we don't want our feed to automatically update. This would cause our pages to have duplicate data when traversing over it.

Hence the backend controls the starting point of the `GET /fetch/<userid>/<pagenumber>` API with a timestamp which is recorded when `GET /fetch` call is made to register a user. Following calls by that user only include documents after the timestamp recorded.

The user can hit the endpoint `GET /fetch` once more to refresh the timestamp and page 1 for the user will now include latest records.

## MongoDB Indexes

The service uses multiple indexes to optimise for queries. There are two main collections `users` and `video_metadata`.

For `users` we have a simple index on `userid`

For `video_metadata` we have two sorted indexes - 
  1. `VideoID` Sorted Ascending - This is to optimise the search for duplicates in case Youtube API sends us any
  2. `PublishedAt` Sorted Descending - This is to optimise the fetch query since we fetch data in reverse chronological order.

We also have two text indexs on the `Title` and `Description` field to enable a naive version of fuzzy text search for the Search API.

## How to run the service?

If you are using docker, a simple `docker compose up` should do the work. This will start both the MongoDB and the service.

You would have to set the `YOUTUBE_API_KEYS` variable in `docker-compose.yml` file. This can be a comma seperated list of keys and the backend would cycle through these keys in case one key exceeds it's quota.

## Setting it up for development

1. Start up the MongoDB server
2. Set the credentials for MongoDB in the URL `MONGO_BASE_URL`
3. Set the API keys for Youtube under `YOUTUBE_API_KEYS`
4. Set the Youtube Query under `YOUTUBE_QUERY`

You can set these up in `devsetup/setup.sh` for quick setups later

Running the server is running `go run main.go`

You can find a postman collection under `devsetup` folder to help with structure of API Calls

## Possible Improvements

1. Having a cache for saving User data would be a nice to have to bring down lookup times. Two possible solutions -
   1. Implement an in memory store that gets updated after a Mongo read for the user.
   2. Implement a redis store for the same
