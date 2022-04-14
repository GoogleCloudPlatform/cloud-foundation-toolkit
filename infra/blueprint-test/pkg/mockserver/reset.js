var mockServerClient = require('mockserver-client').mockServerClient;
mockServerClient("localhost", 1080)
    .reset()
