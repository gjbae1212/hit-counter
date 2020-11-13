# HITS

![Hits](https://storage.googleapis.com/hit-counter/main.png)
A simple way to see how many people have visited your website or GitHub repo.
<p align="center">
<a href="https://circleci.com/gh/gjbae1212/hit-counter"><img src="https://circleci.com/gh/gjbae1212/hit-counter.svg?style=svg"></a>
<a href="https://hits.seeyoufarm.com"><img src="https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2Fgjbae1212%2Fhit-counter%2FREADME&count_bg=%2379C83D&title_bg=%23555555&icon=go.svg&icon_color=%2300ADD8&title=hits&edge_flat=false"/></a>
<a href="/LICENSE"><img src="https://img.shields.io/badge/license-GPL-blue.svg" alt="license" /></a>
<a href="https://goreportcard.com/report/github.com/gjbae1212/hit-counter"><img src="https://goreportcard.com/badge/github.com/gjbae1212/hit-counter" alt="Go Report Card" /></a> 
</p>

## Overview

[HITS](https://hits.seeyoufarm.com) provides the SVG badge presented **title** and **daily/total** page count.

If you embed the badge on either website or GitHub or Notion, every page hit will be counted from visitors.

The badge includes per day (from GMT) and the total (all) page count.

[HITS](https://hits.seeyoufarm.com) also shows the GitHub projects with the highest visitors. (TOP 10)

[HITS](https://hits.seeyoufarm.com) shows real-time page hits (using Websocket) of every project  or site that is using this service. 

[HITS](https://hits.seeyoufarm.com) was made by gjbae1212@gmail.com using Golang, WebAssembly (Wasm), HTML, currently serving from Google Cloud platform.
 
## How To Use
### How To Generate The Badge 
You generate the badge through [HITS](https://hits.seeyoufarm.com/#badge).

![Hits](https://storage.googleapis.com/hit-counter/gen.png)

## Features
- Displays daily and total page views on your page.  
- Support badge with customize style.
- Support badge free icon (https://simpleicons.org). 
- Show a graph of your site about daily count of histories in recently 6 month
- Show ranks about github projects.
- Show real-time stream.
      
## ETC
[HITS](https://hits.seeyoufarm.com) counts every page hit without storing sensitive information (IP, headers, etc.).  
To protect from abuse by massive requests, parts of request information are converted to hashing data in local-cache, and it deletes after the elapsed time.

Also, HITS does not use GitHub Traffic or Google Analytics data, it simply counts every page hit of your site or repo.
  
## LICENSE
This project is licensed under GPL V3.0.
