var mockServerClient = require('mockserver-client').mockServerClient;
// $!request.headers['Session-Id'] returns an array of values because headers and queryStringParameters have multiple values
mockServerClient("localhost", 1080)
    // .reset()
    .mockAnyResponse({
        "httpRequest": {
            "method": "GET",
            "path": "/storage/v1/b",
            // "queryStringParameters": {
            //     "projection": [
            //     "noAcl"
            //     ],
            //     "project": [
            //     "clf-lz-pp2"
            //     ],
            //     "maxResults": [
            //     "1000"
            //     ],
            //     "fields": [
            //     "items/name,nextPageToken"
            //     ],
            //     "alt": [
            //     "json"
            //     ]
            // },
            // "headers": {
            //     "content-length": [
            //     "0"
            //     ],
            //     "accept-encoding": [
            //     "gzip, deflate"
            //     ],
            //     "accept": [
            //     "application/json"
            //     ],
            //     "Host": [
            //     "storage.googleapis.com"
            //     ],
            //     "Connection": [
            //     "keep-alive"
            //     ]
            // },
            // "keepAlive": true,
            // "secure": true,
            // "remoteAddress": "127.0.0.1"
            },
        "httpResponseTemplate": {
            "templateType": "MUSTACHE",
            "template": `{
                "statusCode": 200,
                "reasonPhrase": "OK",
                "headers": {
                },
                "body": {
                    "items": [
                    {
                        "name": "this-is-fake"
                    },
                    {
                        "name": "116961867251-us-central1-blueprint-config"
                    },
                    {
                        "name": "clf-lz-pp2-{{ request.method }}"
                    },
                    {
                        "name": "clf-lz-pp2-cool-dev"
                    },
                    {
                        "name": "clf-lz-pp2-my-cool-bucket-dev"
                    },
                    {
                        "name": "clf-lz-pp2-my-unique-bucket-dev"
                    },
                    {
                        "name": "clf-lz-pp2_blueprints"
                    }
                    ]
                }
            }`
        }
    }).then(
        function () {
            console.log("expectation created");
        },
        function (error) {
            console.log(error);
        }
    );
