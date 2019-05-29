# go-shortme

> A URL shortener service encodes URLs in base62 and store them in Redis.   
It has 3 features: shorten, info of shortened URL and redirect.

## API Endpoints

1. Shorten URL

> POST /api/shorten

**Params:**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `url`  | `string` | `Required. URL to shorten. e.g. https://www.example.com` |
| `expiration_in_minutes` | `int` | `Required. expiration of short link in minutes. e.g. value 0 represents ` |

**Response:**

<pre>
{
    "shortlink": "P"
}
</pre>

2. Get info of shortened URL

`GET /api/info`

**Params:**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `shortlink`  | `string` | `Required. id of shortened. e.g. P` |

**Response:**

<pre>
{
    "url": "https://www.example.com",
    "created_at": "2017-06-12 22:34:59.468660058 +0800 CST",
    "expiration_in_minutes": 60
}
</pre>

3. Redirect

`GET /:shortlink`

Expand the short link and return a **temporary redirect** HTTP status

