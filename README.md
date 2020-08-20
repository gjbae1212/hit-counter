# HITS

![Hits](https://storage.googleapis.com/hit-counter/main.png)
A simple way to see how many people have visited your website or github project.
<p align="center">
<a href="https://circleci.com/gh/gjbae1212/hit-counter"><img src="https://circleci.com/gh/gjbae1212/hit-counter.svg?style=svg"></a>
<a href="https://hits.seeyoufarm.com"/><img src="https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2Fgjbae1212%2Fhit-counter"/></a>
<a href="/LICENSE"><img src="https://img.shields.io/badge/license-GPL-blue.svg" alt="license" /></a>
<a href="https://goreportcard.com/report/github.com/gjbae1212/hit-counter"><img src="https://goreportcard.com/badge/github.com/gjbae1212/hit-counter" alt="Go Report Card" /></a> 
</p>

## Overview

[HITS](https://hits.seeyoufarm.com) provides the SVG badge presented **title** and **daily/total** page count.

If you will be put the badge on either website or github or notion and so on, Paging count is calculated when people do visit it.    

Badge includes a day(from GMT) and a total(all) page count.

[HITS](https://hits.seeyoufarm.com) can show github projects with highest paging count.(TOP 10)

[HITS](https://hits.seeyoufarm.com) can show realtime visiting projects using Websocket. 

[HITS](https://hits.seeyoufarm.com) have made by gjbae1212@gmail.com using golang, wasm, html, and so on, currently serving from google cloud platform.
 
## How to use
### How to generate a svg of badge 
You can generate badge through [HITS](https://hits.seeyoufarm.com/#badge).

![Hits](https://storage.googleapis.com/hit-counter/gen.png)

## Features
- Support daily and total badge  
- Support badge with customize style
- Show a graph of your site about daily count of histories in recently 6 month
- Show ranks about github projects.
- Show realtime stream.
      
## ETC
[HITS](https://hits.seeyoufarm.com) is calculated page count without store sensitive information(ip, header, ... and so on).  
For protect abuse by massive requests, parts of request information are converted to hashing data in local-cache, and it deletes after the elapsed time.
  
## LICENSE
This project is following The GPL V3.0.
