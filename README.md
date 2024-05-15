# Barelly and HTTP Server

This is just a very simple HTTP server I've made to help me learn Go.
It's based on the http server challenge in codecrafters.io website.
At the moment the runs statically at port `4221` and only these endpoints are available for multiple tests:

- 'GET /': this is just to see if server is running
- 'GET /echo/{message}': this will return message back in a 200 OK response. This was done to test reading path and parsing response body.
- 'GET /user-agent': will return the 'User-Agent' header if present in request. This was done to test reading request headers and adding response headers.
- 'GET /files/{filename}': returns a file present server side. This was done to test bodies in binary instead of just text and file download.
- 'POST /files/{filename}': this will store a file with the given filename. The content will be what the request gives in the body.

## Next Features

This is the next features I'll try to work on:

- Good URL path parsing (query/path Params);
- Improved body parsing for POST/PUT requests;
- Support for multipart for proper file upload;
- Settings;
- CORS;
- Docker container;
