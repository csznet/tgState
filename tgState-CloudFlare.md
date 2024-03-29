进阶指南
==

Vercel
--

vercel默认域名部分地区会有阻断，建议添加自定义域名  
使用cname解析到```cname-china.vercel-dns.com```  

CloudFlare
--
**SSL证书访问**  
目的：解决开启SSL&Cloudflare CDN后重定向过多问题  
设置路径：域名->SSL/TLS->Overview  
设置为Full(strict)  
![SSL证书访问](https://img-static.csz.net/d/BQACAgUAAxkDAAMUZSV2Wggiieo9_XSgODTLhW6fg-UAAukLAAKrzjBVQ0hH_g6a9OUwBA)

**完全缓存图片**  
目的：加快访问速度，减少api请求次数  
设置路径：域名->Rules->Page Rules->Create Page Rule  
给```/d/*```设置缓存所有，如下  
![完全缓存图片](https://img-static.csz.net/d/BQACAgUAAxkDAAMVZSV2jVzUitEjGJz_GjZwprJ-nV8AAuoLAAKrzjBV7g9PeEBhKrkwBA)  
![完全缓存图片](https://img-static.csz.net/d/BQACAgUAAxkDAAMXZSV2yzXwcPXgwuRctimd5_EDWq8AAuwLAAKrzjBVuAQYIFm1Sv4wBA)  

**控制请求速率**  
目的：防止刷上传  
设置路径：域名->Security->WAF->Rate limiting rules  
建议给```/api```限制在10s不超过2次请求，如下  
![控制请求速率](https://img-static.csz.net/d/BQACAgUAAxkDAAMWZSV2nJe5fOA6DZsdez4DAAG_MWbEAALrCwACq84wVaOhPWnmR--HMAQ)  

**开启Always Online**  
目的：当服务宕机后，图片正常访问  
设置路径：域名->Caching->Configuration->Always Online
![开启Always Online](https://img-static.csz.net/d/BQACAgUAAxkDAAMYZSV23bs8YRmChLhrs2BLwUWsRZ4AAu4LAAKrzjBVHlJjirBp9hgwBA)