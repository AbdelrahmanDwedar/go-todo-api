# Go Todo API

Simple TODO API made using pure [Go](https://go.dev) with [SQLite].
This API uses built-in modules, for most of things, but the SQLite driver.

## Routes

### `/ping`

This route is for testing if the API is live or not.

### `/todo`

This route can show all lists of tasks.

example:

```json
{
  "lists": {
    "1": [
      {
        "id": 1,
        "title": "Get the garbage out",
        "description": "make sure all garbage has been taken out by tonight",
        "listID": 1
      },
      {
        "id": 2,
        "title": "Do the homework",
        "description": "Make sure to finish the homework before afternoon",
        "listID": 1
      }
    ],
    "2": [
      {
        "id": 1,
        "title": "Finish the new three tasks in jira",
        "description": "Finish the tasks assigned to me on jira for today in the work hours",
        "listID": 2
      }
    ]
  }
}
```

### `/todo/new`

A posting route to add a new general task (on the list number 0).

Post content should include the title and the description.

```json
{
  "title": "Help my mother in the house",
  "description": "I need to help my mother with cleaning the house tonight"
}
```

### `/todo/lists/{id}`

A dynamic route with the value `id` being dynamic to the number added instead if it.

This shows a specific list of tasks.

```json
{
  "1": [
    {
      "id": 1,
      "title": "Get the garbage out",
      "description": "make sure all garbage has been taken out by tonight",
      "listID": 1
    },
    {
      "id": 2,
      "title": "Do the homework",
      "description": "Make sure to finish the homework before afternoon",
      "listID": 1
    }
  ]
}
```

### `/todo/lists/new`

This is a post request route creating a new list to have more tasks.

Post content could be empty as the API will automatically create a new list by incrementing the current last list id.

### `/todo/lists/{id}/new`

Post request route to add new task to a specific list.

Post content should include the title and description only.

```json
{
    "title": "Fix the bug in routing",
    "description": "Fix the bug in routes on line 20"
}
```
