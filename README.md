# HITS

![Hits](https://storage.googleapis.com/hit-counter/hits.png)
A simple way to see how many people have visited your website or github project.
<p align="center">
<a href="https://circleci.com/gh/gjbae1212/hit-counter"><img src="https://circleci.com/gh/gjbae1212/hit-counter.svg?style=svg"></a>
<a href="https://hits.seeyoufarm.com"/><img src="https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2Fgjbae1212%2Fhit-counter"/></a>
<a href="/LICENSE"><img src="https://img.shields.io/badge/license-GPL-blue.svg" alt="license" /></a>
<a href="https://goreportcard.com/report/github.com/gjbae1212/hit-counter"><img src="https://goreportcard.com/badge/github.com/gjbae1212/hit-counter" alt="Go Report Card" /></a> 
</p>

## Overview

[HITS](https://hits.seeyoufarm.com) provides Badge of SVG having format whether script or markdown.

If you will be put Badge on either your website or Github project, the paging count is increased Badge count when people do visit its site.    

And Badge involves paging count on both a day(from GMT) and a total(all).

[HITS](https://hits.seeyoufarm.com) will show Github projects of highest paging count.(TOP 10)

And then [HITS](https://hits.seeyoufarm.com) show currently visiting projects by users, using Websocket. 

[HITS](https://hits.seeyoufarm.com) made by gjbae1212 using golang, currently serving from google cloud.
 
## How to use
### How to generate a svg of badge 
You will generate Badge from edit form of url input in [HITS](https://hits.seeyoufarm.com/#badge).

![Hits](https://storage.googleapis.com/hit-counter/generate.png)

Or you could generate a badge by directly writing markdown or html.
```
# markdown
![Hits](https://hits.seeyoufarm.com/api/count/incr/badge.svg?url={your-website or github-project})]

# html
<img src="https://hits.seeyoufarm.com/api/count/incr/badge.svg?url={your-website or github-project}" alt="Hits" />
```

#### If you would like to only represent Badge, not to increase paging count, Edit as below url.
```
# markdown
![Hits](https://hits.seeyoufarm.com/api/count/keep/badge.svg?url={your-website or github-project})]

# html
<img src="https://hits.seeyoufarm.com/api/count/keep/badge.svg?url={your-website or github-project}" alt="Hits" />
```
 
#### If you'd like to change Badge title
Additional query string *title={your-change-badge-name}*
```
# example (increase)
![Hits](https://hits.seeyoufarm.com/api/count/incr/badge.svg?url={your-website or github-project}&title={your-change-badge-name})]

# example (not increase)
![Hits](https://hits.seeyoufarm.com/api/count/keep/badge.svg?url={your-website or github-project}&title={your-change-badge-name})]
```

## Features
- Support daily and total badge  
- Show a graph of your site about daily count of histories in recently 6 month
- Show ranks about github projects.
- Show stream which count of hits is counting in real time.

## How to run your machine
#### If you'd like to run HITS to your machine
- Install [redis](https://redis.io).
- Set environment ```export REDIS_ADDRS=REDIS_URL(ex: localhost:6379)```
- Exec ```$ go build && ./hit-counter```
##### Of course you should modify a little codes associated URL(as hits.seeyoufarm.com) in view/index.html and wasm/main.go 
      
## ETC
[HITS](https://hits.seeyoufarm.com) will increase paging count when getting a badge api is requested.

It will only be to increase paging count in DB(redis), but something as request information(ip, header, ... and so on) don't store in DB(redis).

For only protect abusing increasing a massive of request, part of request information(ip,user-agent) do convert to a value of hashing saved in local-cache and it is deleted after elapsed time.
 
## Inspiration
This project was inspired by [brentvollebregt](https://github.com/brentvollebregt/hit-counter) and [dwyl](https://github.com/dwyl/hits-nodejs/).

Added additional features like daily-count, rank, and so on.
 
## LICENSE
This project is following The GPL V3.0.
