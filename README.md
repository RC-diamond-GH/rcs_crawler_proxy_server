# rcs_crawler_proxy_server

> A proxy server for web crawler

**rcs_crawler_proxy_server** 是一款辅助爬虫工作的、支持缓存的代理服务器。它可以根据 HTTP 请求报文的 **URL**、**请求体**、**请求方法** 等要素来缓存响应包，并在收到相同请求时，直接从缓存返回数据，从而有效减少对目标服务器的重复请求，提升爬虫效率。  

## 功能特性

1. **HTTP 缓存**  
   - 默认根据请求报文 (URL、请求体、请求方法) 进行缓存。  
   - 如果在请求报头中手动添加 `No-Proxy-Cache` 字段（值不限），则此次请求将直接向目标服务器发起请求并返回，不经过缓存。

2. **HTTPS 缓存 (MITM 代理方式)**  
   - 通过中间人 (MITM) 代理实现对 HTTPS 请求的缓存。  
   - 在使用前，需要保证客户端信任此代理服务所使用的证书（`server.crt`、`server.key`）。

3. **可使用外部代理服务器**  
   - 可配置多个外部代理服务器。  
   - 目前对多个外部代理的轮询、并发及负载均衡尚不成熟，需要谨慎使用。

4. **配置灵活，简单易用**  
   - 只需在 `config.json` 中对 Redis 缓存、代理服务端口、TLS 证书等进行配置即可。  
   - 配合 `docker-compose.yml` 和 `Dockerfile`，可轻松进行容器化部署。

## 快速开始

1. **安装并启动 Redis**  
   - 本代理服务器的缓存功能依赖 Redis。请先在本地或远程部署 Redis，并保证在 `config.json` 中填写的连接信息正确。

2. **克隆本仓库并编译**  
   ```bash
   git clone https://github.com/RC-diamond-GH/rcs_crawler_proxy_server.git
   cd rcs_crawler_proxy_server/src
   go build
   ```
   编译完成后会在当前目录生成可执行文件

3. **准备配置文件 `config.json`**  
   以下为一个示例配置：
   ```json
   {
       "Redis": {
           "Host": "redis:6379",
           "Password": "",
           "DB": 0
       },
       "Cache" :{
           "ExpireTime": "60"
       },
       "OuterProxy":[
           "http://10.52.111.143:8080"
       ],
       "ProxySettings":{
           "HttpPort": 8080,
           "HttpsPort": 8081,
           "TLSCert": "server.crt",
           "TLSKey": "server.key"
       }
   }
   ```
   - `ExpireTime` 表示缓存的有效期，单位为分钟。
   - 若要使用外部代理，请在 `OuterProxy` 中添加相应的代理地址。
   - 通过 `ProxySettings` 中的 `HttpPort` 和 `HttpsPort` 可以配置监听的端口；`TLSCert` 和 `TLSKey` 用于 HTTPS 缓存（MITM 代理）时的证书和私钥路径。

4. **启动服务**  
   ```bash
   ./rcs_crawler_proxy_server
   ```
   启动后，即可在 `config.json` 中指定的端口上使用代理服务。

## Docker 部署

1. **修改 `docker-compose.yml` 和 `Dockerfile`**  
   - 根据需要修改镜像名称、端口映射、Redis 连接等信息。

2. **构建并启动容器**  
   ```bash
   docker-compose up --build -d
   ```
   这样就可以使用容器化的方式来部署和运行本服务。
   > 注意：您可能需要通过修改 Dockerfile 中的其中一个 COPY 命令来保证您的证书被复制到容器中

## 注意事项

1. **HTTPS 缓存的实现方式**  
   - 通过 MITM 代理实现，需要客户端信任您的 CA 证书，否则在多数情况下会出现证书错误。

2. **外部代理服务器的负载均衡**  
   - 当前版本在对多代理的使用上处理不够成熟，可能存在轮询与负载均衡方面的问题。

3. **Redis 的使用**  
   - 若对服务性能或缓存持久化等有更高要求，请正确调整 Redis 的配置（如持久化策略），并在生产环境下确保其稳定性。

## 贡献与反馈

- 欢迎对该项目提出 Issues 或提交 Pull Requests。  
- 如果有使用过程中的问题或改进建议，欢迎进行交流。

## 开源许可证

本项目基于 [GNU GENERAL PUBLIC LICENSE v3](./LICENSE) 开源协议进行发布。

如需了解更多详细信息，请查阅 [GNU GPL v3](https://www.gnu.org/licenses/gpl-3.0.html)。
