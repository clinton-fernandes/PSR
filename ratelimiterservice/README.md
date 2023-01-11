StartRateLimiterService() starts a ratelimiter service that listens on port 9090 that has /GetBooks endpoint.

rateLimit = 3
i.e /GetBooks can accept 3 consecutive requests every 5 seconds. After which "Exceeded no. of requests for: /GetBooks " is returned.
The ratelimiter is reset after an interval of 5 seconds.

InitializeRateLimiter()
initializes the RateLimiter object and also adds the list of endpoints to ApiDetails.

resetRequestRate()
is called after regular interval of 5 seconds. This will reset the no. of requests being hit at every available endpoint.

inc()
is called to increment the counter for the request endpoint.

MiddleWare()
will increment the counter for each time the endpoint is hit. It will return an error if the no. of hits have exceeded the limit.



